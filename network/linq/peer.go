package linq

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/network/p2p"
	"land-bridge/network/utils"
)

const (
	maxKnownBlocks = 1024 // Maximum block hashes to keep in the known list (prevent DOS)

	// maxQueuedProps is the maximum number of block propagations to queue up before
	// dropping broadcasts. There's not much point in queueing stale blocks, so a few
	// that might cover uncles should be enough.
	maxQueuedProps = 4

	// maxQueuedAnns is the maximum number of block announcements to queue up before
	// dropping broadcasts. Similarly to block propagations, there's no point to queue
	// above some healthy uncle limit, so use that.
	maxQueuedAnns = 4

	handshakeTimeout = 5 * time.Second
)

// propEvent is a block propagation, waiting for its turn in the broadcast queue.
type propEvent struct {
	block *utils.Block
}

type Peer struct {
	id string // Unique ID for the peer, cached

	*p2p.Peer                   // The embedded P2P package peer
	rw        p2p.MsgReadWriter // Input/output streams for snap

	head    common.Hash // Latest advertised head block hash
	height  *big.Int    // Latest advertised head block height
	version uint        // Protocol version negotiated

	lock sync.RWMutex // Mutex protecting the internal fields

	knownTxs    mapset.Set // Set of transaction hashes known to be known by this peer
	knownBlocks mapset.Set // Set of block hashes known to be known by this peer

	queuedProps chan *propEvent   // Queue of blocks to broadcast to the peer
	queuedAnns  chan *utils.Block // Queue of blocks to announce to the peer
	term        chan struct{}     // Termination channel to stop the broadcasters
}

func (p *Peer) Head() (hash common.Hash, td *big.Int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	copy(hash[:], p.head[:])
	return hash, p.height
}

// SetHead updates the head hash and total difficulty of the peer.
func (p *Peer) SetHead(hash common.Hash, td *big.Int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.head[:], hash[:])
	p.height.Set(td)
}

func (p *Peer) RequestHeadersByHash(origin common.Hash, amount int, skip int, reverse bool) error {
	logs.Debug("Fetching batch of headers", "count", amount, "fromhash", origin, "skip", skip, "reverse", reverse)
	id := rand.Uint64()
	req := &GetBlockHeadersPacket66{
		RequestID: id,
		GetBlockHeadersPacket: &GetBlockHeadersPacket{
			Origin:  HashOrNumber{Hash: origin},
			Amount:  uint64(amount),
			Skip:    uint64(skip),
			Reverse: reverse,
		},
	}
	return p2p.Send(p.rw, GetBlockHeadersMsg, &req)
}

func (p *Peer) RequestHeadersByNumber(origin uint64, amount int, skip int, reverse bool) error {
	logs.Debug("Fetching batch of headers", "count", amount, "fromnum", origin, "skip", skip, "reverse", reverse)
	id := rand.Uint64()
	req := &GetBlockHeadersPacket66{
		RequestID: id,
		GetBlockHeadersPacket: &GetBlockHeadersPacket{
			Origin:  HashOrNumber{Number: origin},
			Amount:  uint64(amount),
			Skip:    uint64(skip),
			Reverse: reverse,
		},
	}
	return p2p.Send(p.rw, GetBlockHeadersMsg, &req)
}

func (p *Peer) Send(msgcode uint64, data interface{}) error {
	return p2p.Send(p.rw, msgcode, data)
}

func (p *Peer) AsyncSendNewBlockHash(block *utils.Block) {
	select {
	case p.queuedAnns <- block:
		// Mark all the block hash as known, but ensure we don't overflow our limits
		p.knownBlocks.Add(block.Hash())
		for p.knownBlocks.Cardinality() >= maxKnownBlocks {
			p.knownBlocks.Pop()
		}
	default:
		logs.Debug("Dropping block announcement", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// NewPeer create a wrapper for a network connection and negotiated  protocol
// version.
func NewPeer(version uint, p *p2p.Peer, rw p2p.MsgReadWriter) *Peer {
	peer := &Peer{
		id:          p.ID().String(),
		Peer:        p,
		rw:          rw,
		version:     version,
		knownBlocks: mapset.NewSet(),

		queuedProps: make(chan *propEvent, maxQueuedProps),
		queuedAnns:  make(chan *utils.Block, maxQueuedAnns),
		term:        make(chan struct{}),
	}
	return peer
}

// Close signals the broadcast goroutine to terminate. Only ever call this if
// you created the peer yourself via NewPeer. Otherwise let whoever created it
// clean it up!
func (p *Peer) Close() {
	close(p.term)
}

// ID retrieves the peer's unique identifier.
func (p *Peer) ID() string {
	return p.id
}

// Handshake executes the eth protocol handshake, negotiating version number,
// network IDs, difficulties, head and genesis blocks.
func (p *Peer) Handshake(network uint64, height *big.Int, block common.Hash) error {
	errc := make(chan error, 2)

	var status StatusPacket // safe to read after two values have been received from errc

	go func() {
		logs.Info("handshake send StatusMsg type", reflect.TypeOf(p.rw))
		errc <- p2p.Send(p.rw, StatusMsg, &StatusPacket{
			ProtocolVersion: uint32(p.version),
			NetworkID:       network,
			Height:          height,
			Head:            block,
		})
	}()
	go func() {
		logs.Info("handshake read StatusMsg")
		errc <- p.readStatus(network, &status)
	}()
	timeout := time.NewTimer(handshakeTimeout)
	defer timeout.Stop()

	for i := 0; i < 2; i++ {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return p2p.DiscReadTimeout
		}
	}

	p.height, p.head = status.Height, status.Head

	logs.Info("pass handshake peer", p.id, p.Name())

	return nil
}

// readStatus reads the remote handshake message.
func (p *Peer) readStatus(network uint64, status *StatusPacket) error {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Code != StatusMsg {
		return fmt.Errorf("%w: first msg has code %x (!= %x)", errNoStatusMsg, msg.Code, StatusMsg)
	}
	if msg.Size > maxMessageSize {
		return fmt.Errorf("%w: %v > %v", errMsgTooLarge, msg.Size, maxMessageSize)
	}
	// Decode the handshake and make sure everything matches
	if err := msg.Decode(&status); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if status.NetworkID != network {
		return fmt.Errorf("%w: %d (!= %d)", errNetworkIDMismatch, status.NetworkID, network)
	}
	if uint(status.ProtocolVersion) != p.version {
		return fmt.Errorf("%w: %d (!= %d)", errProtocolVersionMismatch, status.ProtocolVersion, p.version)
	}

	return nil
}

func (p *Peer) MarkBlock(hash common.Hash) {
	// If we reached the memory allowance, drop a previously known block hash
	for p.knownBlocks.Cardinality() >= maxKnownBlocks {
		p.knownBlocks.Pop()
	}
	p.knownBlocks.Add(hash)
}

// AsyncSendNewBlock queues an entire block for propagation to a remote peer. If
// the peer's broadcast queue is full, the event is silently dropped.
func (p *Peer) AsyncSendNewBlock(block *utils.Block) {
	select {
	case p.queuedProps <- &propEvent{block: block}:
		// Mark all the block hash as known, but ensure we don't overflow our limits
		p.knownBlocks.Add(block.Hash())
		for p.knownBlocks.Cardinality() >= maxKnownBlocks {
			p.knownBlocks.Pop()
		}
	default:
		logs.Debug("Dropping block propagation", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// RequestOneHeader is a wrapper around the header query functions to fetch a
// single header. It is used solely by the fetcher.
func (p *Peer) RequestOneHeader(hash common.Hash) error {
	logs.Debug("Fetching single header", "hash", hash)
	id := rand.Uint64()
	req := &GetBlockHeadersPacket66{
		RequestID: id,
		GetBlockHeadersPacket: &GetBlockHeadersPacket{
			Origin:  HashOrNumber{Hash: hash},
			Amount:  uint64(1),
			Skip:    uint64(0),
			Reverse: false,
		},
	}
	return p2p.Send(p.rw, GetBlockHeadersMsg, &req)

}

// ReplyBlockHeaders is the eth/66 version of SendBlockHeaders.
func (p *Peer) ReplyBlockHeaders(id uint64, headers []*utils.CBlock) error {
	return p2p.Send(p.rw, BlockHeadersMsg, BlockHeadersPacket66{
		RequestID:          id,
		BlockHeadersPacket: headers,
	})
}

func (p *Peer) broadcast() {
	for {
		select {
		case prop := <-p.queuedProps:
			if err := p.SendNewBlock(prop.block, prop.block.Number()); err != nil {
				logs.Error("Propagated block SendNewBlock", "number", prop.block.Number(), "hash", prop.block.Hash(), "err", err)
				return
			}
			logs.Trace("Propagated block", "number", prop.block.Number(), "hash", prop.block.Hash())

		case block := <-p.queuedAnns:
			if err := p.SendNewBlockHashes([]common.Hash{block.Hash()}, []uint64{block.NumberU64()}); err != nil {
				logs.Error("Announced block SendNewBlockHashes", "number", block.Number(), "hash", block.Hash(), "err", err)
				return
			}
			logs.Trace("Announced block", "number", block.Number(), "hash", block.Hash())

		case <-p.term:
			return
		}
	}
}

// SendNewBlockHashes announces the availability of a number of blocks through
// a hash notification.
func (p *Peer) SendNewBlockHashes(hashes []common.Hash, numbers []uint64) error {
	// Mark all the block hashes as known, but ensure we don't overflow our limits
	for _, hash := range hashes {
		p.knownBlocks.Add(hash)
	}
	for p.knownBlocks.Cardinality() >= maxKnownBlocks {
		p.knownBlocks.Pop()
	}
	request := make(newBlockHashesData, len(hashes))
	for i := 0; i < len(hashes); i++ {
		request[i].Hash = hashes[i]
		request[i].Number = numbers[i]
	}
	return p2p.Send(p.rw, NewBlockHashesMsg, request)
}

// SendNewBlock propagates an entire block to a remote peer.
func (p *Peer) SendNewBlock(block *utils.Block, td *big.Int) error {
	// Mark all the block hash as known, but ensure we don't overflow our limits
	p.knownBlocks.Add(block.Hash())
	for p.knownBlocks.Cardinality() >= maxKnownBlocks {
		p.knownBlocks.Pop()
	}
	logs.Trace("NewBlockMsg send", "hash", block.Hash().Hex())

	return p2p.Send(p.rw, NewBlockMsg, []interface{}{block, td})
}
