package linq

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"

	"land-bridge/network/consensus"
	"land-bridge/network/linq/downloader"
	"land-bridge/network/linq/fetcher"
	"land-bridge/network/p2p"
	"land-bridge/network/p2p/enode"
	"land-bridge/network/utils"
)

type handler struct {
	acceptTxs uint32 // Flag whether we're considered synchronised (enables transaction processing)

	networkID uint64
	maxPeers  int

	checkpointNumber uint64      // Block number for the sync progress validator to cross reference
	checkpointHash   common.Hash // Block hash for the sync progress validator to cross reference

	downloader *downloader.Downloader
	fetcher    *fetcher.Fetcher
	peers      *peerSet

	engine consensus.Engine

	blockStore *Store

	eventMux      *event.TypeMux
	minedBlockSub *event.TypeMuxSubscription

	quitSync chan struct{}

	chainSync *chainSyncer
	wg        sync.WaitGroup
	peerWG    sync.WaitGroup
}

func (h *handler) Enqueue(id string, block *utils.CBlock) {
	h.fetcher.Enqueue(id, block)
}

func (h *handler) FindPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	m := make(map[common.Address]consensus.Peer)
	for _, p := range h.peers.Peers() {
		pubKey := p.Node().Pubkey()
		addr := crypto.PubkeyToAddress(*pubKey)
		if targets[addr] {
			m[addr] = p
		}
	}
	return m
}

type Handler func(peer *Peer) error

type Backend interface {
	Chain() *Store

	AcceptTxs() bool
	// RunPeer is invoked when a peer joins on the `eth` protocol. The handler
	// should do any peer maintenance work, handshakes and validations. If all
	// is passed, control should be given back to the `handler` to process the
	// inbound messages going forward.
	RunPeer(peer *Peer, handler Handler) error

	// PeerInfo retrieves all known `eth` information about a peer.
	PeerInfo(id enode.ID) interface{}

	// Handle is a callback to be invoked when a data packet is received from
	// the remote peer. Only packets not consumed by the protocol handler will
	// be forwarded to the backend.
	Handle(peer *Peer, packet Packet) error

	Engine() consensus.Engine
}

func newHandler(engine consensus.Engine, blockStore *Store, mux *event.TypeMux) (*handler, error) {
	h := &handler{
		peers:      newPeerSet(),
		quitSync:   make(chan struct{}),
		engine:     engine,
		blockStore: blockStore,
		eventMux:   mux,
	}

	if handler, ok := h.engine.(consensus.Handler); ok {
		logs.Debug("handler SetBroadcaster")
		handler.SetBroadcaster(h)
	}

	h.downloader = downloader.New(h.checkpointNumber, blockStore.db, h.eventMux, blockStore, h.removePeer)

	// Construct the fetcher (short sync)
	validator := func(header *utils.CBlock) error {
		return engine.VerifyHeader(blockStore, header, true)
	}
	heighter := func() uint64 {
		return blockStore.CurrentBlock().Height
	}
	inserter := func(blocks utils.Blocks) (int, error) {
		n, err := h.blockStore.InsertChain(blocks)
		if err == nil {
			atomic.StoreUint32(&h.acceptTxs, 1) // Mark initial sync done on any fetcher import
		}
		return n, err
	}
	h.fetcher = fetcher.New(blockStore.GetBlockByHash, validator, h.BroadcastBlock, heighter, inserter, h.unregisterPeer)

	h.chainSync = newChainSyncer(h)

	return h, nil
}

func (h *handler) runEthPeer(peer *Peer, handler Handler) error {
	h.peerWG.Add(1)
	defer h.peerWG.Done()
	var (
		head   = h.blockStore.CurrentBlock()
		hash   = head.Hash()
		height = head.Number()
	)

	if err := peer.Handshake(h.networkID, height, hash); err != nil {
		logs.Debug("Ethereum handshake failed", "err", err)
		return err
	}

	logs.Debug("Ethereum peer connected", "name", peer.Name(), "id", peer.id)

	// Register the peer locally
	if err := h.peers.registerPeer(peer); err != nil {
		logs.Error("Ethereum peer registration failed", "err", err)
		return err
	}
	defer h.unregisterPeer(peer.ID())

	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	p := h.peers.peer(peer.ID())
	if p == nil {
		return errors.New("peer dropped during handling")
	}

	if err := h.downloader.RegisterPeer(peer.id, int(peer.version), peer); err != nil {
		return err
	}

	return handler(peer)
}

// Handle is invoked whenever an `eth` connection is made that successfully passes
// the protocol handshake. This method will keep processing messages until the
// connection is torn down.
func Handle(backend Backend, peer *Peer) error {
	logs.Trace("new peer begin to loop handle msg")
	for {
		if err := handleMessage(backend, peer); err != nil {
			logs.Debug("Message handling failed in `eth`", "err", err)
			return err
		}
	}
}

// handleMessage is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func handleMessage(backend Backend, peer *Peer) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := peer.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > maxMessageSize {
		return fmt.Errorf("%w: %v > %v", errMsgTooLarge, msg.Size, maxMessageSize)
	}
	defer msg.Discard()

	if handler, ok := backend.Engine().(consensus.Handler); ok {
		pubKey := peer.Node().Pubkey()
		addr := crypto.PubkeyToAddress(*pubKey)
		handled, err := handler.HandleMsg(addr, msg)
		if handled {
			if err != nil {
				logs.Error("handleMessage HandleMsg error ", err)
			}
			return err
		}
	} else {
		logs.Debug("handleMessage Engine type error ")
	}

	var handlers = msglist

	if handler := handlers[msg.Code]; handler != nil {
		return handler(backend, msg, peer)
	}
	return fmt.Errorf("%w: %v", errInvalidMsgCode, msg.Code)
}

// unregisterPeer removes a peer from the downloader, fetchers and main peer set.
func (h *handler) unregisterPeer(id string) {
	// Abort if the peer does not exist
	peer := h.peers.peer(id)
	if peer == nil {
		logs.Error("Ethereum peer removal failed", "err", errPeerNotRegistered)
		return
	}

	h.downloader.UnregisterPeer(id)

	if err := h.peers.unregisterPeer(id); err != nil {
		logs.Error("Ethereum peer removal failed", "err", err)
	}
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (h *handler) Start(maxPeers int) {
	logs.Trace("land handler start")
	h.maxPeers = maxPeers

	// broadcast mined blocks
	h.wg.Add(1)
	h.minedBlockSub = h.eventMux.Subscribe(NewMinedBlockEvent{})
	go h.minedBroadcastLoop()

	// start sync handlers
	h.wg.Add(1)
	go h.chainSync.loop()
}

// minedBroadcastLoop sends mined blocks to connected peers.
func (h *handler) minedBroadcastLoop() {
	defer h.wg.Done()

	for obj := range h.minedBlockSub.Chan() {
		if ev, ok := obj.Data.(NewMinedBlockEvent); ok {
			h.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			h.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}

func (h *handler) Stop() {
	h.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	// Quit chainSync and txsync64.
	// After this is done, no new peers will be accepted.
	close(h.quitSync)
	h.wg.Wait()

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to h.peers yet
	// will exit when they try to register.
	h.peers.close()
	h.peerWG.Wait()

	logs.Info("Ethereum protocol stopped")
}

// BroadcastBlock will either propagate a block to a subset of it's peers, or
// will only announce it's availability (depending what's requested).
func (h *handler) BroadcastBlock(block *utils.Block, propagate bool) {
	hash := block.Hash()
	peers := h.peers.PeersWithoutBlock(hash)

	if propagate {
		// Send the block to a subset of our peers
		transfer := peers[:int(math.Sqrt(float64(len(peers))))]
		for _, peer := range transfer {
			peer.AsyncSendNewBlock(block)
		}
		logs.Trace("Propagated block", "hash", hash, "height", block.Height, "recipients", len(transfer))
		return
	}

	// Otherwise if the block is indeed in out own chain, announce it
	if h.blockStore.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.AsyncSendNewBlockHash(block)
		}
		logs.Trace("Announced block", "hash", hash, "height", block.Height, "recipients", len(peers))
	}
}

func (h *handler) removePeer(id string) {
	// Short circuit if the peer was already removed
	peer := h.peers.peer(id)
	if peer == nil {
		return
	}
	logs.Debug("Removing Ethereum peer", "peer", id)

	// Unregister the peer from the downloader and Ethereum peer set
	h.downloader.UnregisterPeer(id)
	if err := h.peers.Unregister(id); err != nil {
		logs.Error("Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

// doSync synchronizes the local blockchain with a remote peer.
func (h *handler) doSync(op *chainSyncOp) error {

	// Run the sync cycle, and disable snap sync if we're past the pivot block
	err := h.downloader.Synchronise(op.peer.ID(), op.head, op.height)
	if err != nil {
		return err
	}
	// If we've successfully finished a sync cycle and passed any required checkpoint,
	// enable accepting transactions from the network.
	head := h.blockStore.CurrentBlock()
	if head.NumberU64() >= h.checkpointNumber {
		// Checkpoint passed, sanity check the timestamp to have a fallback mechanism
		// for non-checkpointed (number = 0) private networks.
		if head.Time >= uint64(time.Now().AddDate(0, -1, 0).Unix()) {
			atomic.StoreUint32(&h.acceptTxs, 1)
		}
	}
	if head.NumberU64() > 0 {
		// We've completed a sync cycle, notify all peers of new state. This path is
		// essential in star-topology networks where a gateway node needs to notify
		// all its out-of-date peers of the availability of a new block. This failure
		// scenario will most often crop up in private and hackathon networks with
		// degenerate connectivity, but it should be healthy for the mainnet too to
		// more reliably update peers or the local TD state.
		h.BroadcastBlock(head, false)
	}
	return nil
}

// NodeInfo represents a short summary of the `eth` sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Height    uint64      `json:"height"`    // Total difficulty of the host's blockchain
	Head      common.Hash `json:"head"`      // Hex hash of the host's best owned block
	Consensus string      `json:"consensus"` // Consensus mechanism in use
}

// nodeInfo retrieves some `eth` protocol metadata about the running host node.
func nodeInfo(chain *Store) *NodeInfo {
	head := chain.CurrentBlock()
	return &NodeInfo{
		Height: head.Height,
		Head:   head.Hash(),
	}
}
