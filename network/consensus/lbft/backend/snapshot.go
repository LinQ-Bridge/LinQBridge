package backend

import (
	"bytes"
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"

	"land-bridge/models"
	"land-bridge/network/consensus/lbft"
	"land-bridge/network/consensus/lbft/validator"
	"land-bridge/network/utils"
)

// Vote represents a single vote that an authorized validator made to modify the
// list of authorizations.
type Vote struct {
	Validator common.Address `json:"validator"` // Authorized validator that cast this vote
	Block     uint64         `json:"block"`     // Block number the vote was cast in (expire old votes)
	Address   common.Address `json:"address"`   // Account being voted on to change its authorization
	Authorize bool           `json:"authorize"` // Whether to authorize or deauthorize the voted account
}

// Tally is a simple vote tally to keep the current score of votes. Votes that
// go against the proposal aren't counted since it's equivalent to not voting.
type Tally struct {
	Authorize bool `json:"authorize"` // Whether the vote it about authorizing or kicking someone
	Votes     int  `json:"votes"`     // Number of votes until now wanting to pass the proposal
}

type Snapshot struct {
	Epoch uint64 // The number of blocks after which to checkpoint and reset the pending votes

	Number uint64                   // Block number where the snapshot was created
	Hash   common.Hash              // Block hash where the snapshot was created
	Votes  []*Vote                  // List of votes cast in chronological order
	Tally  map[common.Address]Tally // Current vote tally to avoid recalculating
	ValSet lbft.ValidatorSet        // Set of authorized validators at this moment
}

// newSnapshot create a new snapshot with the specified startup parameters. This
// method does not initialize the set of recent validators, so only ever use if for
// the genesis block.
func newSnapshot(epoch uint64, number uint64, hash common.Hash, valSet lbft.ValidatorSet) *Snapshot {
	snap := &Snapshot{
		Epoch:  epoch,
		Number: number,
		Hash:   hash,
		ValSet: valSet,
		Tally:  make(map[common.Address]Tally),
	}
	return snap
}

// loadSnapshot loads an existing snapshot from the database.
func loadSnapshot(epoch uint64, db *gorm.DB, hash common.Hash) (*Snapshot, error) {
	ssModel := &models.Snapshot{}
	if err := db.Where("hash = ?", hash.Hex()).Find(ssModel).Error; err != nil {
		logs.Error("GORM snapshot err: ", err)
		return nil, err
	}
	snap := new(Snapshot)
	if err := json.Unmarshal([]byte(ssModel.Bytes), snap); err != nil {
		return nil, err
	}
	snap.Epoch = epoch

	return snap, nil
}

// store inserts the snapshot into the database.
func (s *Snapshot) store(db *gorm.DB) error {
	blob, err := json.Marshal(s)
	if err != nil {
		return err
	}

	ssModel := &models.Snapshot{}
	if err := db.Model(&models.Snapshot{}).Where("hash = ?", s.Hash.Hex()).First(&ssModel).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	ssModel.Hash = s.Hash.Hex()
	ssModel.Bytes = string(blob)

	return db.Save(ssModel).Error
}

// copy creates a deep copy of the snapshot, though not the individual votes.
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		Epoch:  s.Epoch,
		Number: s.Number,
		Hash:   s.Hash,
		ValSet: s.ValSet.Copy(),
		Votes:  make([]*Vote, len(s.Votes)),
		Tally:  make(map[common.Address]Tally),
	}

	for address, tally := range s.Tally {
		cpy.Tally[address] = tally
	}
	copy(cpy.Votes, s.Votes)

	return cpy
}

// checkVote return whether it's a valid vote
func (s *Snapshot) checkVote(address common.Address, authorize bool) bool {
	_, validator := s.ValSet.GetByAddress(address)
	return (validator != nil && !authorize) || (validator == nil && authorize)
}

// cast adds a new vote into the tally.
func (s *Snapshot) cast(address common.Address, authorize bool) bool {
	// Ensure the vote is meaningful
	if !s.checkVote(address, authorize) {
		return false
	}
	logs.Debug("apply cast pass checkVote")
	// Cast the vote into an existing or new tally
	if old, ok := s.Tally[address]; ok {
		old.Votes++
		s.Tally[address] = old
	} else {
		s.Tally[address] = Tally{Authorize: authorize, Votes: 1}
		logs.Debug("Tally add ", address.Hash())
	}
	return true
}

// uncast removes a previously cast vote from the tally.
func (s *Snapshot) uncast(address common.Address, authorize bool) bool {
	// If there's no tally, it's a dangling vote, just drop
	tally, ok := s.Tally[address]
	if !ok {
		return false
	}
	// Ensure we only revert counted votes
	if tally.Authorize != authorize {
		return false
	}
	// Otherwise revert the vote
	if tally.Votes > 1 {
		tally.Votes--
		s.Tally[address] = tally
	} else {
		delete(s.Tally, address)
	}
	return true
}

// apply creates a new authorization snapshot by applying the given headers to
// the original one.
func (s *Snapshot) apply(blocks []*utils.Block) (*Snapshot, error) {
	// Allow passing in no headers for cleaner code
	if len(blocks) == 0 {
		return s, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < len(blocks)-1; i++ {
		if blocks[i+1].Height != blocks[i].Height+1 {
			return nil, errInvalidVotingChain
		}
	}
	if blocks[0].Height != s.Number+1 {
		return nil, errInvalidVotingChain
	}
	// Iterate through the headers and create a new snapshot
	snap := s.copy()

	for _, block := range blocks {
		// Remove any votes on checkpoint blocks
		number := block.Height
		if number%s.Epoch == 0 {
			snap.Votes = nil
			snap.Tally = make(map[common.Address]Tally)
		}
		// Resolve the authorization key and check against validators
		validator, err := ecrecover(block)
		if err != nil {
			return nil, err
		}
		if _, v := snap.ValSet.GetByAddress(validator); v == nil {
			return nil, errUnauthorized
		}
		candidate := common.BytesToAddress(block.CBytes[:20])
		// Header authorized, discard any previous votes from the validator
		for i, vote := range snap.Votes {
			if vote.Validator == validator && vote.Address == candidate {
				// Uncast the vote from the cached tally
				snap.uncast(vote.Address, vote.Authorize)

				// Uncast the vote from the chronological list
				snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)
				break // only one vote allowed
			}
		}
		// Tally up the new vote from the validator
		var authorize bool

		switch {
		case bytes.Compare(block.CBytes[20:], nonceAuthVote) == 0:
			authorize = true
		case bytes.Compare(block.CBytes[20:], nonceDropVote) == 0:
			authorize = false
		default:
			return nil, errInvalidVote
		}

		if snap.cast(candidate, authorize) {
			snap.Votes = append(snap.Votes, &Vote{
				Validator: validator,
				Block:     number,
				Address:   candidate,
				Authorize: authorize,
			})
		}

		// If the vote passed, update the list of validators
		if tally := snap.Tally[candidate]; tally.Votes > snap.ValSet.Size()/2 {
			logs.Debug("snapshot Tally", snap.Tally[candidate].Authorize)
			if tally.Authorize {
				snap.ValSet.AddValidator(candidate)
			} else {
				snap.ValSet.RemoveValidator(candidate)

				// Discard any previous votes the deauthorized validator cast
				for i := 0; i < len(snap.Votes); i++ {
					if snap.Votes[i].Validator == candidate {
						// Uncast the vote from the cached tally
						snap.uncast(snap.Votes[i].Address, snap.Votes[i].Authorize)

						// Uncast the vote from the chronological list
						snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)

						i--
					}
				}
			}
			// Discard any previous votes around the just changed account
			for i := 0; i < len(snap.Votes); i++ {
				if snap.Votes[i].Address == candidate {
					snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)
					i--
				}
			}
			delete(snap.Tally, candidate)
		}
	}
	snap.Number += uint64(len(blocks))
	snap.Hash = blocks[len(blocks)-1].Hash()

	return snap, nil
}

// validators retrieves the list of authorized validators in ascending order.
func (s *Snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, s.ValSet.Size())
	for _, validator := range s.ValSet.List() {
		validators = append(validators, validator.Address())
	}
	for i := 0; i < len(validators); i++ {
		for j := i + 1; j < len(validators); j++ {
			if bytes.Compare(validators[i][:], validators[j][:]) > 0 {
				validators[i], validators[j] = validators[j], validators[i]
			}
		}
	}
	return validators
}

type snapshotJSON struct {
	Epoch  uint64                   `json:"epoch"`
	Number uint64                   `json:"number"`
	Hash   common.Hash              `json:"hash"`
	Votes  []*Vote                  `json:"votes"`
	Tally  map[common.Address]Tally `json:"tally"`

	// for validator set
	Validators []common.Address    `json:"validators"`
	Policy     lbft.ProposerPolicy `json:"policy"`
}

func (s *Snapshot) toJSONStruct() *snapshotJSON {
	return &snapshotJSON{
		Epoch:      s.Epoch,
		Number:     s.Number,
		Hash:       s.Hash,
		Votes:      s.Votes,
		Tally:      s.Tally,
		Validators: s.validators(),
		Policy:     s.ValSet.Policy(),
	}
}

// Unmarshal from a json byte array
func (s *Snapshot) UnmarshalJSON(b []byte) error {
	var j snapshotJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	s.Epoch = j.Epoch
	s.Number = j.Number
	s.Hash = j.Hash
	s.Votes = j.Votes
	s.Tally = j.Tally
	s.ValSet = validator.NewSet(j.Validators, j.Policy)
	return nil
}

// Marshal to a json byte array
func (s *Snapshot) MarshalJSON() ([]byte, error) {
	j := s.toJSONStruct()
	return json.Marshal(j)
}
