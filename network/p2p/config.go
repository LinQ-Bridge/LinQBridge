package p2p

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common/mclock"

	"land-bridge/network/p2p/enode"
	"land-bridge/network/p2p/netutil"
)

type Config struct {
	PrivateKey     *ecdsa.PrivateKey
	MaxPeers       int
	ListenAddr     string
	Name           string
	BootstrapNodes []*enode.Node
	Protocols      []Protocol
	NetRestrict    *netutil.Netlist
	clock          mclock.Clock
	StaticNodes    []*enode.Node
}
