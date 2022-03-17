package linq

import (
	"sync"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"

	"land-bridge/network/bridge"
	"land-bridge/network/consensus"
	"land-bridge/network/consensus/lbft"
	"land-bridge/network/consensus/lbft/backend"
	"land-bridge/network/node"
	"land-bridge/network/p2p"
	"land-bridge/network/p2p/enode"
)

// lbft consensus Protocol variables are optionally set in addition to the "eth" protocol variables (eth/protocol.go).
var lbftConsensusProtocolName = ""

// ProtocolVersions are the supported versions of the linq consensus protocol (first is primary), e.g. []uint{LBFT64, LBFT99, LBFT100}.
var lbftConsensusProtocolVersions []uint

// protocol Length describe the number of messages support by the protocol/version map[uint]uint64{LBFT64: 18, LBFT99: 18, LBFT100: 18}
var lbftConsensusProtocolLengths map[uint]uint64

type LinQ struct {
	eventMux *event.TypeMux

	blockStore *Store
	bridge     *bridge.Bridge
	engine     consensus.Engine
	worker     *worker

	handler           *handler
	ethDialCandidates enode.Iterator
	p2pServer         *p2p.Server

	lock sync.RWMutex
}

func New(stack *node.Node, privStr string) (*LinQ, error) {
	lq := &LinQ{
		p2pServer:  stack.Server(),
		blockStore: NewLinQStore(stack.GetDB(), stack.Pool),
		bridge:     stack.Bridge,
		eventMux:   stack.EventMux(),
	}

	lq.engine = CreateConsensusEngine(stack)

	lbftProtocol := lq.engine.Protocol()
	// set the linq specific consensus devp2p subprotocol, eth subprotocol remains set to protocolName as in upstream geth.
	lbftConsensusProtocolName = lbftProtocol.Name
	lbftConsensusProtocolVersions = lbftProtocol.Versions
	lbftConsensusProtocolLengths = lbftProtocol.Lengths

	var err error
	lq.handler, err = newHandler(lq.engine, lq.blockStore, lq.eventMux)
	if err != nil {
		return nil, err
	}

	lq.worker = newWorker(lq.blockStore, stack.Bridge, stack.Pool, stack.EventMux(), lq.engine, common.HexToAddress(privStr))
	lq.blockStore.SetHeadCh(lq.worker.chainHeadCh)

	stack.RegisterProtocols(lq.Protocols())
	stack.RegisterLifecycle(lq)

	return lq, nil
}

func (lq *LinQ) Start() error {
	logs.Info("linq work start")
	go lq.worker.start()
	maxPeers := lq.p2pServer.MaxPeers
	go lq.handler.Start(maxPeers)
	return nil
}

func (lq *LinQ) Stop() error {
	// Stop all the peer-related stuff first.
	lq.ethDialCandidates.Close()

	lq.handler.Stop()

	return nil
}

// Protocols returns all the currently configured
// network protocols to start.
func (lq *LinQ) Protocols() []p2p.Protocol {
	var protos []p2p.Protocol

	//consensus Protocol
	lbftprotos := ConsensusProtocols((*linqHandler)(lq.handler))
	protos = append(protos, lbftprotos...)

	return protos
}

func CreateConsensusEngine(stack *node.Node) consensus.Engine {
	return backend.New(lbft.DefaultConfig, stack.GetNodeKey(), stack.GetDB(), stack.Bridge, stack.Pool)
}
