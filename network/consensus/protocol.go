package consensus

import (
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/network/utils"
)

// Constants to match up protocol versions and messages
// LBFT/99 was added to accommodate new eth/64 handshake status data with fork id
// this is for backward compatibility which allows a mixed old/new lbft node network
// lbft/64 will continue using old status data as eth/63
const (
	LinQ99 = 99
)

var (
	LinQProtocol = Protocol{
		Name:     "linQ",
		Versions: []uint{LinQ99},
		Lengths:  map[uint]uint64{LinQ99: 18},
	}
)

// Protocol defines the protocol of the consensus
type Protocol struct {
	// Official short name of the protocol used during capability negotiation.
	Name string
	// Supported versions of the eth protocol (first is primary).
	Versions []uint
	// Number of implemented message corresponding to different protocol versions.
	Lengths map[uint]uint64
}

// Broadcaster defines the interface to enqueue blocks to fetcher and find peer
type Broadcaster interface {
	// Enqueue add a block into fetcher queue
	Enqueue(id string, block *utils.CBlock)
	// FindPeers retrives peers by addresses
	FindPeers(map[common.Address]bool) map[common.Address]Peer
}

// Peer defines the interface to communicate with peer
type Peer interface {
	// Send sends the message to this peer
	Send(msgcode uint64, data interface{}) error
}
