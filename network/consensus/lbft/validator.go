package lbft

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Validator interface {
	// Address returns address
	Address() common.Address

	// String representation of Validator
	String() string
}

// ----------------------------------------------------------------------------

type Validators []Validator

func (slice Validators) Len() int {
	return len(slice)
}

func (slice Validators) Less(i, j int) bool {
	return strings.Compare(slice[i].String(), slice[j].String()) < 0
}

func (slice Validators) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ----------------------------------------------------------------------------

type ValidatorSet interface {
	// CalcProposer Calculate the proposer
	CalcProposer(lastProposer common.Address, round uint64)
	// Size Return the validator size
	Size() int
	// List Return the validator array
	List() []Validator
	// GetByIndex Get validator by index
	GetByIndex(i uint64) Validator
	// GetByAddress Get validator by given address
	GetByAddress(addr common.Address) (int, Validator)
	// GetProposer Get current proposer
	GetProposer() Validator
	// IsProposer Check whether the validator with given address is a proposer
	IsProposer(address common.Address) bool
	// AddValidator Add validator
	AddValidator(address common.Address) bool
	// RemoveValidator Remove validator
	RemoveValidator(address common.Address) bool
	// Copy Copy validator set
	Copy() ValidatorSet
	// F Get the maximum number of faulty nodes
	F() int
	// Policy Get proposer policy
	Policy() ProposerPolicy
}

// ----------------------------------------------------------------------------

type ProposalSelector func(ValidatorSet, common.Address, uint64) Validator
