package lbft

import "math/big"

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	Sticky
)

type Config struct {
	RequestTimeout         uint64         `toml:",omitempty"` // The timeout for each LBFT round in milliseconds.
	BlockPeriod            uint64         `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second
	ProposerPolicy         ProposerPolicy `toml:",omitempty"` // The policy for proposer selection
	Epoch                  uint64         `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes
	Ceil2Nby3Block         *big.Int       `toml:",omitempty"` // Number of confirmations required to move from one state to next [2F + 1 to Ceil(2N/3)]
	AllowedFutureBlockTime uint64         `toml:",omitempty"` // Max time (in seconds) from current time allowed for blocks, before they're considered future blocks
}

var DefaultConfig = &Config{
	RequestTimeout:         2000,
	BlockPeriod:            1,
	ProposerPolicy:         RoundRobin,
	Epoch:                  30000,
	Ceil2Nby3Block:         big.NewInt(0),
	AllowedFutureBlockTime: 0,
}
