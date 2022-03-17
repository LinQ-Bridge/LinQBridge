package backend

import (
	"bytes"
	"errors"
	"math/rand"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"

	"land-bridge/network/consensus"
	"land-bridge/network/consensus/lbft"
	"land-bridge/network/consensus/lbft/core"
	"land-bridge/network/consensus/lbft/validator"
	"land-bridge/network/utils"
)

const (
	checkpointInterval = 4   // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128 // Number of recent vote snapshots to keep in memory
	inmemoryPeers      = 40
	inmemoryMessages   = 1024
)

var (
	// errInvalidProposal is returned when a prposal is malformed.
	errInvalidProposal = errors.New("invalid proposal")
	// errInvalidSignature is returned when given signature is not signed by given
	// address.
	errInvalidSignature = errors.New("invalid signature")
	// errUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")

	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")
	// errInvalidVote is returned if a nonce value is something else that the two
	// allowed constants of 0x00..0 or 0xff..f.
	errInvalidVote = errors.New("vote nonce not 0x00..0 or 0xff..f")
	// errInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	errInvalidCommittedSeals = errors.New("invalid committed seals")
	// errEmptyCommittedSeals is returned if the field of committed seals is zero.
	errEmptyCommittedSeals = errors.New("zero committed seals")
	// errInvalidTimestamp is returned if the timestamp of a block is lower than the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")
)

var (
	now = time.Now

	nonceAuthVote = hexutil.MustDecode("0x0001") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000") // Magic nonce number to vote on removing a validator.

	inmemoryAddresses  = 20 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)
)

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the block's coinbase if a consensus
// engine is based on signatures.
func (sb *backend) Author(block *utils.Block) (common.Address, error) {
	return ecrecover(block)
}

// VerifyHeader checks whether a block conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *backend) VerifyHeader(chain consensus.ChainReader, block *utils.CBlock, seal bool) error {
	return sb.verifyHeader(chain, block, nil)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *backend) verifyHeader(chain consensus.ChainReader, block *utils.CBlock, parents []*utils.CBlock) error {
	if block.Number() == nil {
		return errUnknownBlock
	}

	// Don't waste time checking blocks from the future (adjusting for allowed threshold)
	adjustedTimeNow := now().Add(time.Duration(sb.config.AllowedFutureBlockTime) * time.Second).Unix()
	if block.Time > uint64(adjustedTimeNow) {
		return lbft.ErrFutureBlock
	}

	parentsBlock := utils.CBlocks(parents).ToBlock()
	return sb.verifyCascadingFields(chain, block.ToBlock(), parentsBlock)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (sb *backend) verifyCascadingFields(chain consensus.ChainReader, header *utils.Block, parents []*utils.Block) error {
	// The genesis block is the always valid dead-end
	number := header.NumberU64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *utils.Block
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetBlock(header.ParentHash, number-1)
	}
	if parent == nil || parent.NumberU64() != number-1 || parent.Hash() != header.ParentHash {
		return lbft.ErrUnknownAncestor
	}
	if parent.Time+sb.config.BlockPeriod > header.Time {
		return errInvalidTimestamp
	}

	// Verify validators in extraData. Validators in snapshot and extraData should be the same.
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}
	validators := make([]byte, len(snap.validators())*common.AddressLength)
	for i, validator := range snap.validators() {
		copy(validators[i*common.AddressLength:], validator[:])
	}
	if err := sb.verifySigner(chain, header, parents); err != nil {
		return err
	}

	return sb.verifyCommittedSeals(chain, header, parents)
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (sb *backend) verifyCommittedSeals(chain consensus.ChainReader, header *utils.Block, parents []*utils.Block) error {
	number := header.NumberU64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	cheader := header.ToCBlock()

	// Ensure that the extra data format is satisfied
	extra, err := utils.LBFTBlockExtra(cheader)
	if err != nil {
		return err
	}

	// The length of Committed seals should be larger than 0
	if len(extra.CommittedSeal) == 0 {
		return errEmptyCommittedSeals
	}

	validators := snap.ValSet.Copy()
	// Check whether the committed seals are generated by parent's validators
	validSeal := 0
	committers, err := sb.Signers(cheader)
	if err != nil {
		return err
	}
	for _, addr := range committers {
		if validators.RemoveValidator(addr) {
			validSeal++
			continue
		}
		return errInvalidCommittedSeals
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= snap.ValSet.F() {
		return errInvalidCommittedSeals
	}

	return nil
}

// Signers extracts all the addresses who have signed the given header
// It will extract for each seal who signed it, regardless of if the seal is
// repeated
func (sb *backend) Signers(header *utils.CBlock) ([]common.Address, error) {
	extra, err := utils.LBFTBlockExtra(header)
	if err != nil {
		return []common.Address{}, err
	}

	var addrs []common.Address
	proposalSeal := core.PrepareCommittedSeal(header.Hash())

	// 1. Get committed seals from current header
	for _, seal := range extra.CommittedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := lbft.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			logs.Error("not a valid address", "err", err)
			return nil, errInvalidSignature
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *backend) VerifyHeaders(chain consensus.ChainReader, blocks []*utils.CBlock, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(blocks))
	go func() {
		for i, header := range blocks {
			err := sb.verifyHeader(chain, header, blocks[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifySeal checks whether the crypto seal on a block is valid according to
// the consensus rules of the given engine.
func (sb *backend) VerifySeal(chain consensus.ChainReader, block *utils.Block) error {
	// get parent header and ensure the signer is in parent's validator set
	number := block.Height
	if number == 0 {
		return errUnknownBlock
	}

	return sb.verifySigner(chain, block, nil)
}

// Prepare initializes the consensus fields of a block according to the
// rules of a particular engine. The changes are executed inline.
func (sb *backend) Prepare(chain consensus.ChainReader, block *utils.CBlock) error {
	// unused fields, force to set to empty
	block.Coinbase = common.Address{}

	// copy the parent extra data as the header extra data
	number := block.Height
	parent := chain.GetBlock(block.ParentHash, number-1)
	logs.Info("Prepare GetBlock", block.ParentHash, number-1)
	if parent == nil {
		return lbft.ErrUnknownAncestor
	}

	// Assemble the voting snapshot
	snap, err := sb.snapshot(chain, number-1, block.ParentHash, nil)
	if err != nil {
		return err
	}

	// get valid candidate list
	sb.candidatesLock.RLock()
	var addresses []common.Address
	var authorizes []bool
	for address, authorize := range sb.candidates {
		if snap.checkVote(address, authorize) {
			addresses = append(addresses, address)
			authorizes = append(authorizes, authorize)
		}
	}
	sb.candidatesLock.RUnlock()

	// pick one of the candidates randomly
	if len(addresses) > 0 {
		index := rand.Intn(len(addresses))
		// add validator voting in coinbase
		block.Coinbase = addresses[index]
		if authorizes[index] {
			copy(block.CBytes[20:], nonceAuthVote)
		} else {
			copy(block.CBytes[20:], nonceDropVote)
		}
	}

	// add validators in snapshot to extraData's validators section
	extra, err := prepareExtra(block, snap.validators())
	if err != nil {
		return err
	}
	block.ExtraData = extra

	// set header's timestamp
	block.Time = parent.Time + sb.config.BlockPeriod
	if block.Time < uint64(time.Now().Unix()) {
		block.Time = uint64(time.Now().Unix())
	}
	return nil
}

// Seal generates a new sealing request for the given input block and pushes
// the result into the given channel.
//
// Note, the method returns immediately and will send the result async. More
// than one result may also be returned depending on the consensus algorithm.
func (sb *backend) Seal(chain consensus.ChainReader, block *utils.CBlock, results chan<- *utils.Block, stop <-chan struct{}) error {
	number := block.Height
	snap, err := sb.snapshot(chain, number-1, block.ParentHash, nil)
	if err != nil {
		return err
	}
	if _, v := snap.ValSet.GetByAddress(sb.address); v == nil {
		return errUnauthorized
	}

	block, err = sb.updateBlock(block)

	if err != nil {
		return err
	}

	delay := time.Unix(int64(block.Time), 0).Sub(now())

	go func() {
		// wait for the timestamp of header, use this to adjust the block period
		select {
		case <-time.After(delay):
		case <-stop:
			results <- nil
			return
		}

		// get the proposed block hash and clear it if the seal() is completed.
		sb.sealMu.Lock()
		sb.proposedBlockHash = block.Hash()

		defer func() {
			sb.proposedBlockHash = common.Hash{}
			sb.sealMu.Unlock()
		}()

		// post block into LBFT engine
		go sb.EventMux().Post(consensus.RequestEvent{
			Proposal: block,
		})

		for {
			select {
			case result := <-sb.commitCh:
				// if the block hash and the hash from channel are the same,
				// return the result. Otherwise, keep waiting the next hash.
				if result != nil && block.Hash() == result.Hash() {
					results <- result
					return
				}
			case <-stop:
				results <- nil
				return
			}
		}
	}()
	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (sb *backend) SealHash(block *utils.Block) common.Hash {
	return sigHash(block)
}

// update timestamp and signature of the block based on its number of transactions
func (sb *backend) updateBlock(block *utils.CBlock) (*utils.CBlock, error) {
	b := block.ToBlock()
	// sign the hash
	seal, err := sb.Sign(sigHash(b).Bytes())

	if err != nil {
		return nil, err
	}

	err = writeSeal(b, seal)
	if err != nil {
		return nil, err
	}

	b.Hash()

	newblock := b.ToCBlock()

	newblock.TxParam = block.TxParam

	return newblock, nil
}

// Protocol returns the protocol for this consensus
func (sb *backend) Protocol() consensus.Protocol {
	return consensus.LinQProtocol
}

// Close terminates any background threads maintained by the consensus engine.
func (sb *backend) Close() error {
	return nil
}

// snapshot retrieves the authorization snapshot at a given point in time.
func (sb *backend) snapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*utils.Block) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		blocks []*utils.Block
		snap   *Snapshot
	)

	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := sb.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(sb.config.Epoch, sb.db, hash); err == nil {
				logs.Trace("Loaded voting snapshot form disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}

		// If we're at block zero, make a snapshot
		if number == 0 {
			genesis := chain.GetBlockByNumber(0).ToCBlock()
			if err := sb.VerifyHeader(chain, genesis, false); err != nil {
				return nil, err
			}
			lbftExtra, err := utils.LBFTBlockExtra(genesis)
			if err != nil {
				return nil, err
			}
			snap = newSnapshot(sb.config.Epoch, 0, genesis.Hash(), validator.NewSet(lbftExtra.Validators, sb.config.ProposerPolicy))
			if err := snap.store(sb.db); err != nil {
				return nil, err
			}
			logs.Trace("Stored genesis voting snapshot to disk")
			break
		}

		// No snapshot for this header, gather the header and move backward
		var block *utils.Block
		if len(parents) > 0 {
			// If we have explicit parents, pick from there (enforced)
			block = parents[len(parents)-1]
			if block.Hash() != hash || block.Height != number {
				return nil, lbft.ErrUnknownAncestor
			}
			parents = parents[:len(parents)-1]
		} else {
			// No explicit parents (or no more left), reach out to the database
			block = chain.GetBlock(hash, number)
			if block == nil {
				return nil, lbft.ErrUnknownAncestor
			}
		}

		blocks = append(blocks, block)
		number, hash = number-1, block.ParentHash
	}
	// Previous snapshot found, apply any pending blocks on top of it
	for i := 0; i < len(blocks)/2; i++ {
		blocks[i], blocks[len(blocks)-1-i] = blocks[len(blocks)-1-i], blocks[i]
	}
	snap, err := snap.apply(blocks)
	if err != nil {
		logs.Error("snapshot apply error ", err)
		return nil, err
	}
	sb.recents.Add(snap.Hash, snap)

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(blocks) > 0 {
		if err = snap.store(sb.db); err != nil {
			logs.Error("snapshot store error ", err)
			return nil, err
		}
		logs.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}
	return snap, err
}

// sigHash returns the hash which is used as input for the LBFT
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(block *utils.Block) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, utils.LBFTBlockForEncode(block, false))

	hasher.Sum(hash[:0])
	return hash
}

// verifySigner checks whether the signer is in parent's validator set
func (sb *backend) verifySigner(chain consensus.ChainReader, block *utils.Block, parents []*utils.Block) error {
	// Verifying the genesis block is not supported
	number := block.Height
	if number == 0 {
		return errUnknownBlock
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := sb.snapshot(chain, number-1, block.ParentHash, parents)
	if err != nil {
		return err
	}

	// resolve the authorization key and check against signers
	signer, err := ecrecover(block)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := snap.ValSet.GetByAddress(signer); v == nil {
		return errUnauthorized
	}
	return nil
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(block *utils.Block) (common.Address, error) {
	hash := block.Hash()
	if addr, ok := recentAddresses.Get(hash); ok {
		return addr.(common.Address), nil
	}

	// Retrieve the signature from the header extra-data
	extra, err := block.LBFTBlockExtra()
	if err != nil {
		return common.Address{}, err
	}

	addr, err := lbft.GetSignatureAddress(sigHash(block).Bytes(), extra.Seal)
	if err != nil {
		return addr, err
	}
	recentAddresses.Add(hash, addr)
	return addr, nil
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(block *utils.CBlock, vals []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough LBFTExtraVanity bytes.
	if len(block.ExtraData) < utils.LBFTExtraVanity {
		block.ExtraData = append(block.ExtraData, bytes.Repeat([]byte{0x00}, utils.LBFTExtraVanity-len(block.ExtraData))...)
	}
	buf.Write(block.ExtraData[:utils.LBFTExtraVanity])

	ist := &utils.LBFTExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(block *utils.Block, seal []byte) error {
	if len(seal)%utils.LBFTExtraSeal != 0 {
		return errInvalidSignature
	}

	lbftExtra, err := block.LBFTBlockExtra()
	if err != nil {
		return err
	}

	lbftExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&lbftExtra)
	if err != nil {
		return err
	}

	block.ExtraData = append(block.ExtraData[:utils.LBFTExtraVanity], payload...)
	return nil
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
func writeCommittedSeals(block *utils.CBlock, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return errInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != utils.LBFTExtraSeal {
			return errInvalidCommittedSeals
		}
	}

	lbftExtra, err := block.LBFTBlockExtra()
	if err != nil {
		return err
	}

	lbftExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(lbftExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&lbftExtra)
	if err != nil {
		return err
	}

	block.ExtraData = append(block.ExtraData[:utils.LBFTExtraVanity], payload...)
	return nil
}
