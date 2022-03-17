package linq

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/network/consensus"
	"land-bridge/network/p2p/enode"
	"land-bridge/network/utils"
)

type linqHandler handler

func (l *linqHandler) RunPeer(peer *Peer, hand Handler) error {
	return (*handler)(l).runEthPeer(peer, hand)
}

func (l *linqHandler) PeerInfo(id enode.ID) interface{} {
	if p := l.peers.peer(id.String()); p != nil {
		return nil
	}
	return nil
}

func (l *linqHandler) Chain() *Store {
	return l.blockStore
}

func (l *linqHandler) AcceptTxs() bool {
	return atomic.LoadUint32(&l.acceptTxs) == 1
}

func (l *linqHandler) Handle(peer *Peer, packet Packet) error {
	// Consume any broadcasts and announces, forwarding the rest to the downloader
	switch packet := packet.(type) {
	case *NewBlockHashesPacket:
		hashes, numbers := packet.Unpack()
		return l.handleBlockAnnounces(peer, hashes, numbers)
	case *NewBlockPacket:
		return l.handleBlockBroadcast(peer, packet.Block, packet.TD)
	case *BlockHeadersPacket66:
		return l.handleBlockHeaderst(peer, packet.BlockHeadersPacket)
	default:
		return fmt.Errorf("unexpected eth packet type: %T", packet)
	}
}

func (l *linqHandler) Engine() consensus.Engine {
	return l.engine
}

// handleBlockAnnounces is invoked from a peer's message handler when it transmits a
// batch of block announcements for the local node to process.
func (l *linqHandler) handleBlockAnnounces(peer *Peer, hashes []common.Hash, numbers []uint64) error {
	var (
		unknownHashes  = make([]common.Hash, 0, len(hashes))
		unknownNumbers = make([]uint64, 0, len(numbers))
	)
	for i := 0; i < len(hashes); i++ {
		if !l.blockStore.HasBlock(hashes[i], numbers[i]) {
			unknownHashes = append(unknownHashes, hashes[i])
			unknownNumbers = append(unknownNumbers, numbers[i])
		}
	}
	for i := 0; i < len(unknownHashes); i++ {
		l.fetcher.Notify(peer.ID(), unknownHashes[i], unknownNumbers[i], time.Now(), peer.RequestOneHeader)
	}
	return nil
}

// handleBlockBroadcast is invoked from a peer's message handler when it transmits a
// block broadcast for the local node to process.
func (l *linqHandler) handleBlockBroadcast(peer *Peer, block *utils.Block, td *big.Int) error {
	// Schedule the block for import
	l.fetcher.Enqueue(peer.ID(), block.ToCBlock())

	var (
		trueHead = block.ParentHash
		trueTD   = td
	)
	// Update the peer's total difficulty if better than the previous
	if _, td := peer.Head(); trueTD.Cmp(td) > 0 {
		peer.SetHead(trueHead, trueTD)
		l.chainSync.handlePeerEvent(peer)
	}
	logs.Trace("NewBlockMsg get", "hash", block.Hash().Hex(), "height", block.Height)
	return nil
}

func (l *linqHandler) handleBlockHeaderst(peer *Peer, blocks []*utils.CBlock) error {

	filter := len(blocks) == 1
	if filter {
		blocks = l.fetcher.FilterHeaders(peer.id, blocks, time.Now())
	}

	if len(blocks) > 0 || !filter {
		err := l.downloader.DeliverHeaders(peer.id, blocks)
		if err != nil {
			logs.Debug("Failed to deliver headers", "err", err)
		}
	}

	return nil
}
