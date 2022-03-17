// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"land-bridge/network/p2p"
	"land-bridge/network/utils"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during block and/or uncle verification.
type ChainReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentBlock retrieves the current block from the local chain.
	CurrentBlock() *utils.Block

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *utils.Block

	// GetBlockByNumber retrieves a block from the database by number.
	GetBlockByNumber(number uint64) *utils.Block

	// GetBlockByHash retrieves a block from the database by its hash.
	GetBlockByHash(hash common.Hash) *utils.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the block's coinbase if a consensus
	// engine is based on signatures.
	Author(block *utils.Block) (common.Address, error)

	// VerifyHeader checks whether a block conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainReader, block *utils.CBlock, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainReader, headers []*utils.CBlock, seals []bool) (chan<- struct{}, <-chan error)

	// VerifySeal checks whether the crypto seal on a block is valid according to
	// the consensus rules of the given engine.
	VerifySeal(chain ChainReader, block *utils.Block) error

	// Prepare initializes the consensus fields of a block according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainReader, block *utils.CBlock) error

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	//
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	Seal(chain ChainReader, block *utils.CBlock, results chan<- *utils.Block, stop <-chan struct{}) error

	// SealHash returns the hash of a block prior to it being sealed.
	SealHash(block *utils.Block) common.Hash

	// Protocol returns the protocol for this consensus
	Protocol() Protocol

	// Close terminates any background threads maintained by the consensus engine.
	Close() error
}

// Handler should be implemented is the consensus needs to handle and send peer's message
type Handler interface {
	// NewBlock handles a new head block comes
	NewBlock() error

	// HandleMsg handles a message from peer
	HandleMsg(address common.Address, data p2p.Msg) (bool, error)

	// SetBroadcaster sets the broadcaster to send message to peers
	SetBroadcaster(Broadcaster)
}

// LBFT is a consensus engine to avoid byzantine failure
type LBFT interface {
	Engine

	// Start starts the engine
	Start(chain ChainReader, currentBlock func() *utils.Block) error

	// Stop stops the engine
	Stop() error
}
