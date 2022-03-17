package linq

import (
	"github.com/beego/beego/v2/core/logs"

	"land-bridge/network/p2p"
	"land-bridge/network/p2p/enode"
)

func ConsensusProtocols(backend Backend) []p2p.Protocol {
	protos := make([]p2p.Protocol, len(lbftConsensusProtocolVersions))
	for i, vsn := range lbftConsensusProtocolVersions {
		length, ok := lbftConsensusProtocolLengths[vsn]
		if !ok {
			panic("makeProtocol for unknown version")
		}
		lp := makeLbftProtocol(backend, lbftConsensusProtocolName, vsn, length)
		protos[i] = lp
	}
	return protos
}

func makeLbftProtocol(backend Backend, protoName string, version uint, length uint64) p2p.Protocol {
	logs.Debug("registering a legacy protocol ", "protoName", protoName, "version", version)
	return p2p.Protocol{
		Name:    protoName,
		Version: version,
		Length:  length,
		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
			peer := NewPeer(version, p, rw)
			defer peer.Close()
			return backend.RunPeer(peer, func(peer *Peer) error {
				return Handle(backend, peer)
			})
		},
		NodeInfo: func() interface{} {
			node := nodeInfo(backend.Chain())
			node.Consensus = "lbft"
			return node
		},
		PeerInfo: func(id enode.ID) interface{} {
			return backend.PeerInfo(id)
		},
	}
}
