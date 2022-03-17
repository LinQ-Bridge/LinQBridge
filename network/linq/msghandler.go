package linq

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/network/utils"
)

const (
	// softResponseLimit is the target maximum size of replies to data retrievals.
	softResponseLimit = 2 * 1024 * 1024

	// estHeaderSize is the approximate size of an RLP encoded block header.
	estHeaderSize = 500

	// maxHeadersServe is the maximum number of block headers to serve. This number
	// is there to limit the number of disk lookups.
	maxHeadersServe = 1024
)

type msgHandler func(backend Backend, msg Decoder, peer *Peer) error

const (
	StatusMsg          = 0x00
	NewBlockHashesMsg  = 0x01
	GetBlockHeadersMsg = 0x03
	BlockHeadersMsg    = 0x04
	NewBlockMsg        = 0x07
)

var msglist = map[uint64]msgHandler{
	NewBlockHashesMsg:  handleNewBlockhashes,
	NewBlockMsg:        handleNewBlock,
	GetBlockHeadersMsg: handleGetBlockHeaders,
	BlockHeadersMsg:    handleBlockHeaders,
}

func handleNewBlockhashes(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of new block announcements just arrived
	ann := new(NewBlockHashesPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Mark the hashes as present at the remote node
	for _, block := range *ann {
		peer.MarkBlock(block.Hash)
	}
	// Deliver them all to the backend for queuing
	return backend.Handle(peer, ann)
}

func handleNewBlock(backend Backend, msg Decoder, peer *Peer) error {
	ann := new(NewBlockPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}

	if err := ann.sanityCheck(); err != nil {
		return err
	}

	// Mark the peer as owning the block and schedule it for import
	peer.MarkBlock(ann.Block.Hash())

	return backend.Handle(peer, ann)
}

func handleBlockHeaders(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of headers arrived to one of our previous requests
	headers := new(BlockHeadersPacket66)
	if err := msg.Decode(headers); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	return backend.Handle(peer, headers)
}

func handleGetBlockHeaders(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the complex header query
	var query GetBlockHeadersPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetBlockHeadersQuery(backend.Chain(), query.GetBlockHeadersPacket, peer)
	blocks := utils.Blocks(response)
	return peer.ReplyBlockHeaders(query.RequestID, blocks.ToCBlock())
}

// ServiceGetBlockHeadersQuery assembles the response to a header query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetBlockHeadersQuery(chain *Store, query *GetBlockHeadersPacket, peer *Peer) []*utils.Block {
	hashMode := query.Origin.Hash != (common.Hash{})
	first := true
	maxNonCanonical := uint64(100)

	// Gather headers until the fetch or network limits is reached
	var (
		bytes   common.StorageSize
		headers []*utils.Block
		unknown bool
		lookups int
	)
	for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit &&
		len(headers) < maxHeadersServe && lookups < 2*maxHeadersServe {
		lookups++
		// Retrieve the next header satisfying the query
		var origin *utils.Block
		if hashMode {
			if first {
				first = false
				origin = chain.GetBlockByHash(query.Origin.Hash)
				if origin != nil {
					query.Origin.Number = origin.NumberU64()
				}
			} else {
				origin = chain.GetBlockByHash(query.Origin.Hash)
			}
		} else {
			origin = chain.GetBlockByNumber(query.Origin.Number)
		}
		if origin == nil {
			break
		}
		headers = append(headers, origin)
		bytes += estHeaderSize

		// Advance to the next header of the query
		switch {
		case hashMode && query.Reverse:
			// Hash based traversal towards the genesis block
			ancestor := query.Skip + 1
			if ancestor == 0 {
				unknown = true
			} else {
				query.Origin.Hash, query.Origin.Number = chain.GetAncestor(query.Origin.Hash, query.Origin.Number, ancestor, &maxNonCanonical)
				unknown = query.Origin.Hash == common.Hash{}
			}
		case hashMode && !query.Reverse:
			// Hash based traversal towards the leaf block
			var (
				current = origin.NumberU64()
				next    = current + query.Skip + 1
			)
			if next <= current {
				logs.Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", peer.Peer.Name())
				unknown = true
			} else {
				if header := chain.GetBlockByNumber(next); header != nil {
					nextHash := header.Hash()
					expOldHash, _ := chain.GetAncestor(nextHash, next, query.Skip+1, &maxNonCanonical)
					if expOldHash == query.Origin.Hash {
						query.Origin.Hash, query.Origin.Number = nextHash, next
					} else {
						unknown = true
					}
				} else {
					unknown = true
				}
			}
		case query.Reverse:
			// Number based traversal towards the genesis block
			if query.Origin.Number >= query.Skip+1 {
				query.Origin.Number -= query.Skip + 1
			} else {
				unknown = true
			}

		case !query.Reverse:
			// Number based traversal towards the leaf block
			query.Origin.Number += query.Skip + 1
		}
	}
	return headers
}
