package backend

import (
	"crypto/ecdsa"
	"gorm.io/gorm"
	"math/big"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	lru "github.com/hashicorp/golang-lru"

	"land-bridge/network/bridge"
	"land-bridge/network/consensus"
	"land-bridge/network/consensus/lbft"
	"land-bridge/network/consensus/lbft/core"
	"land-bridge/network/consensus/lbft/validator"
	"land-bridge/network/linq/txblock"
	"land-bridge/network/utils"
)

const (
	// fetcherID is the ID indicates the block is from LBFT engine
	fetcherID = "lbft"
)

// New creates an Ethereum backend for LBFT core engine.
func New(config *lbft.Config, privateKey *ecdsa.PrivateKey, db *gorm.DB, bridge *bridge.Bridge, pool *txblock.BlockPool) *backend {
	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)

	sb := &backend{
		config:            config,
		consensusEventMux: new(event.TypeMux),
		privateKey:        privateKey,
		address:           crypto.PubkeyToAddress(privateKey.PublicKey),
		db:                db,
		bridge:            bridge,
		pool:              pool,
		commitCh:          make(chan *utils.Block, 1),
		recents:           recents,
		candidates:        make(map[common.Address]bool),
		coreStarted:       false,
		recentMessages:    recentMessages,
		knownMessages:     knownMessages,
	}

	return sb
}

type backend struct {
	config *lbft.Config
	bridge *bridge.Bridge
	pool   *txblock.BlockPool

	privateKey *ecdsa.PrivateKey
	address    common.Address

	chain        consensus.ChainReader
	currentBlock func() *utils.Block

	core consensus.Core

	consensusEventMux *event.TypeMux

	db *gorm.DB

	commitCh          chan *utils.Block
	proposedBlockHash common.Hash
	sealMu            sync.Mutex
	coreStarted       bool
	coreMu            sync.RWMutex

	// Current list of candidates we are pushing
	candidates map[common.Address]bool
	// Protects the signer fields
	candidatesLock sync.RWMutex

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	recents        *lru.ARCCache
	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages
}

func (sb *backend) Address() common.Address {
	return sb.address
}

func (sb *backend) Bridge() *bridge.Bridge {
	return sb.bridge
}

func (sb *backend) Pool() *txblock.BlockPool {
	return sb.pool
}

func (sb *backend) Validators(proposal consensus.Proposal) lbft.ValidatorSet {
	return sb.getValidators(proposal.Number().Uint64(), proposal.Hash())
}

func (sb *backend) EventMux() *event.TypeMux {
	return sb.consensusEventMux
}

func (sb *backend) Broadcast(valSet lbft.ValidatorSet, payload []byte) error {
	// send to others
	err := sb.Gossip(valSet, payload)
	if err != nil {
		logs.Error("Broadcast Gossip error", err)
	}
	// send to self
	msg := consensus.MessageEvent{
		Payload: payload,
	}
	go sb.consensusEventMux.Post(msg)
	return nil
}

func (sb *backend) Gossip(valSet lbft.ValidatorSet, payload []byte) error {
	hash := lbft.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	targets := make(map[common.Address]bool)
	for _, val := range valSet.List() {
		if val.Address() != sb.Address() {
			targets[val.Address()] = true
		}
	}

	if sb.broadcaster != nil && len(targets) > 0 {
		ps := sb.broadcaster.FindPeers(targets)
		for addr, p := range ps {
			ms, ok := sb.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)
			go p.Send(LbftMsg, payload)
		}
	}
	return nil
}

func (sb *backend) Commit(proposal consensus.Proposal, seals [][]byte) error {
	// Check if the proposal is a valid block
	block := &utils.CBlock{}
	block, ok := proposal.(*utils.CBlock)
	if !ok {
		logs.Error("Invalid proposal, %v", proposal)
		return errInvalidProposal
	}

	// Append seals into extra-data
	err := writeCommittedSeals(block, seals)
	if err != nil {
		return err
	}

	sb.pool.Delete(proposal.TxHash().Hex()[2:])

	// - if the proposed and committed blocks are the same, send the proposed hash
	//   to commit channel, which is being watched inside the engine.Seal() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	b := block.ToBlock()
	if sb.proposedBlockHash == b.Hash() {
		// feed block hash to Seal() and wait the Seal() result
		sb.commitCh <- b
		return nil
	}

	if sb.broadcaster != nil {
		logs.Trace("follower node add block into fetcher queue hash", block.Hash().Hex())
		sb.broadcaster.Enqueue(fetcherID, block)
	}

	return nil
}

func (sb *backend) Verify(proposal consensus.Proposal) (time.Duration, error) {
	// Check if the proposal is a valid block
	block := &utils.CBlock{}
	block, ok := proposal.(*utils.CBlock)
	if !ok {
		return 0, errInvalidProposal
	}

	// verify the header of proposed block
	err := sb.VerifyHeader(sb.chain, block, false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == errEmptyCommittedSeals {
		return 0, nil
	} else if err == lbft.ErrFutureBlock {
		return time.Unix(int64(block.Time), 0).Sub(now()), lbft.ErrFutureBlock
	}
	return 0, err
}

func (sb *backend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, sb.privateKey)
}

func (sb *backend) CheckSignature(data []byte, addr common.Address, sig []byte) error {
	signer, err := lbft.GetSignatureAddress(data, sig)
	if err != nil {
		//log.Error("Failed to get signer address", "err", err)
		return err
	}

	// Compare derived addresses
	if signer != addr {
		return errInvalidSignature
	}
	return nil
}

func (sb *backend) LastProposal() (consensus.Proposal, common.Address) {
	block := sb.currentBlock()
	var proposer common.Address
	if block.Number().Cmp(common.Big0) > 0 {
		var err error
		proposer, err = sb.Author(block)
		if err != nil {
			logs.Error("Failed to get block proposer block", block.Height, "err", err)
			return nil, common.Address{}
		}
	}

	// Return header only block here since we don't need block body
	return block.ToCBlock(), proposer
}

func (sb *backend) HasPropsal(hash common.Hash, number *big.Int) bool {
	return sb.chain.GetBlock(hash, number.Uint64()) != nil
}

func (sb *backend) GetProposer(number uint64) common.Address {
	if number == 0 {
		return common.Address{}
	}

	if h := sb.chain.GetBlockByNumber(number); h != nil {
		a, _ := sb.Author(h)
		return a
	}
	return common.Address{}
}

func (sb *backend) ParentValidators(proposal consensus.Proposal) lbft.ValidatorSet {
	if block, ok := proposal.(*utils.CBlock); ok {
		return sb.getValidators(block.Number().Uint64()-1, block.ParentHash)
	}
	return validator.NewSet(nil, sb.config.ProposerPolicy)
}

func (sb *backend) getValidators(number uint64, hash common.Hash) lbft.ValidatorSet {
	snap, err := sb.snapshot(sb.chain, number, hash, nil)
	if err != nil {
		return validator.NewSet(nil, sb.config.ProposerPolicy)
	}
	return snap.ValSet
}

func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *utils.Block) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return lbft.ErrStartedEngine
	}

	// clear previous data
	sb.proposedBlockHash = common.Hash{}
	if sb.commitCh != nil {
		close(sb.commitCh)
	}
	sb.commitCh = make(chan *utils.Block, 1)

	sb.chain = chain
	sb.currentBlock = currentBlock

	err := sb.startLBFT()

	if err != nil {
		return err
	}

	sb.coreStarted = true

	return nil
}

func (sb *backend) Stop() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if !sb.coreStarted {
		return lbft.ErrStoppedEngine
	}
	if err := sb.core.Stop(); err != nil {
		return err
	}
	sb.coreStarted = false
	return nil
}

func (sb *backend) startLBFT() error {
	logs.Info("BFT: activate LBFT")
	logs.Trace("BFT: set ProposerPolicy sorter to ValidatorSortByStringFun")

	sb.core = core.New(sb, sb.config)
	if err := sb.core.Start(); err != nil {
		logs.Error("BFT: failed to activate LBFT", "err", err)
		return err
	}

	return nil
}
