package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	"land-bridge/network/bridge"
)

var emptyHash = common.Hash{}

type Blocks []*Block

func (s Blocks) ToCBlock() []*CBlock {
	var res []*CBlock
	for _, block := range s {
		res = append(res, block.ToCBlock())
	}
	return res
}

type Block struct {
	ParentHash common.Hash    `json:"parentHash"`
	BlockHash  common.Hash    `json:"blockHash"`
	Coinbase   common.Address `json:"coinbase"`
	Height     uint64         `json:"height"`
	Time       uint64         `json:"timestamp"`
	Type       TxType         `json:"type"`
	SrcTx      TxInfo         `json:"srcTx"`

	// For consensus(lbft: vote)
	CBytes LBFTBytes `json:"cBytes"`

	// These fields are not required to hash
	DstTx     TxInfo `json:"dstTx"`
	ExtraData []byte `json:"extraData"`
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	ParentHash common.Hash
	BlockHash  common.Hash
	Coinbase   common.Address
	Height     uint64
	Time       uint64
	Type       TxType
	SrcTx      TxInfo
	CBytes     LBFTBytes
	DstTx      TxInfo
	ExtraData  []byte
}

func (b *Block) UnmarshalJSON(input []byte) error {
	type block struct {
		ParentHash *common.Hash         `json:"parentHash"`
		BlockHash  *common.Hash         `json:"blockHash"`
		Coinbase   *common.Address      `json:"coinbase"`
		Height     *math.HexOrDecimal64 `json:"height"`
		Time       *math.HexOrDecimal64 `json:"timestamp"`
		Type       *math.HexOrDecimal64 `json:"type"`
		SrcTx      *TxInfo              `json:"srcTx"`
		CBytes     *LBFTBytes           `json:"cBytes"`
		DstTx      *TxInfo              `json:"dstTx"`
		ExtraData  *hexutil.Bytes       `json:"extraData"`
		TxParam    *bridge.TxParam      `json:"txParam"`
	}

	var dec block
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.ParentHash != nil {
		b.ParentHash = *dec.ParentHash
	}
	if dec.BlockHash != nil {
		b.BlockHash = *dec.BlockHash
	}
	if dec.Coinbase != nil {
		b.Coinbase = *dec.Coinbase
	}
	if dec.Height == nil {
		return errors.New("missing required field 'Height' for Genesis")
	}
	b.Height = uint64(*dec.Height)

	if dec.Time == nil {
		return errors.New("missing required field 'Time' for Genesis")
	}
	b.Time = uint64(*dec.Time)

	if dec.Type != nil {
		b.Type = TxType(uint8(*dec.Type))
	}

	if dec.SrcTx != nil {
		b.SrcTx = *dec.SrcTx
	}

	if dec.CBytes != nil {
		b.CBytes = *dec.CBytes
	}

	if dec.DstTx != nil {
		b.DstTx = *dec.DstTx
	}

	if dec.ExtraData != nil {
		b.ExtraData = *dec.ExtraData
	}

	return nil
}

func (b *Block) Number() *big.Int {
	return big.NewInt(int64(b.Height))
}

func (b *Block) Hash() common.Hash {
	if b.BlockHash != emptyHash {
		return b.BlockHash
	} else {
		b.BlockHash = b.proposedBlockHash()
		return b.BlockHash
	}
}

func (b *Block) TxHash() common.Hash {
	return common.BytesToHash(b.SrcTx.TxHash)
}

func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		ParentHash: b.ParentHash,
		BlockHash:  b.BlockHash,
		Coinbase:   b.Coinbase,
		Height:     b.Height,
		Time:       b.Time,
		Type:       b.Type,
		SrcTx:      b.SrcTx,
		CBytes:     b.CBytes,
		DstTx:      b.DstTx,
		ExtraData:  b.ExtraData,
	})
}

func (b *Block) ToCBlock() *CBlock {
	return &CBlock{
		ParentHash: b.ParentHash,
		BlockHash:  b.BlockHash,
		Coinbase:   b.Coinbase,
		Height:     b.Height,
		Time:       b.Time,
		Type:       b.Type,
		SrcTx:      b.SrcTx,
		CBytes:     b.CBytes,
		DstTx:      b.DstTx,
		ExtraData:  b.ExtraData,
	}
}

func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock

	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.ParentHash = eb.ParentHash
	b.BlockHash = eb.BlockHash
	b.Coinbase = eb.Coinbase
	b.Height = eb.Height
	b.Time = eb.Time
	b.Type = eb.Type
	b.SrcTx = eb.SrcTx
	b.CBytes = eb.CBytes
	b.DstTx = eb.DstTx
	b.ExtraData = eb.ExtraData
	return nil
}

func (b *Block) String() string {
	return fmt.Sprintf("{Block: %v}", *b)
}

func (b *Block) LBFTBlockExtra() (*LBFTExtra, error) {
	if len(b.ExtraData) < LBFTExtraVanity {
		return nil, ErrInvalidLBFTBlockExtra
	}

	var le *LBFTExtra

	err := rlp.DecodeBytes(b.ExtraData[LBFTExtraVanity:], &le)
	if err != nil {
		return nil, err
	}
	return le, nil
}

var headerSize = common.StorageSize(reflect.TypeOf(Block{}).Size())

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (b *Block) Size() common.StorageSize {
	return headerSize + common.StorageSize(len(b.ExtraData)+(b.Number().BitLen())/8)
}

func (b *Block) NumberU64() uint64 {
	return b.Height
}

// SanityCheck checks a few basic things -- these checks are way beyond what
// any 'sane' production values should hold, and can mainly be used to prevent
// that the unbounded fields are stuffed with junk data to add processing
// overhead
func (b *Block) SanityCheck() error {
	if eLen := len(b.ExtraData); eLen > 100*1024 {
		return fmt.Errorf("too large block extradata: size %d", eLen)
	}
	return nil
}

func (b *Block) proposedBlockHash() common.Hash {
	if istanbulHeader := LBFTBlockForEncode(b, true); istanbulHeader != nil {
		return rlpHash(istanbulHeader)
	} else {
		return common.Hash{}
	}
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type TxInfo struct {
	ChainID uint64 `json:"chain_id"`
	TxHash  []byte `json:"tx_hash"`
}

func (t *TxInfo) UnmarshalJSON(input []byte) error {
	type txInfo struct {
		ChainID *math.HexOrDecimal64 `json:"chain_id"`
		TxHash  *hexutil.Bytes       `json:"tx_hash"`
	}

	var dec txInfo
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ChainID != nil {
		t.ChainID = uint64(*dec.ChainID)
	}

	if dec.TxHash != nil {
		t.TxHash = *dec.TxHash
	}

	return nil
}

type BlockEvent struct{ Block *Block }

func LBFTBlockForEncode(b *Block, keepSeal bool) *Block {
	newBlock := *b
	newBlock.DstTx = TxInfo{}
	newBlock.BlockHash = common.Hash{}
	newBlock.ExtraData = make([]byte, len(b.ExtraData))
	copy(newBlock.ExtraData, b.ExtraData)

	extra, err := b.LBFTBlockExtra()
	if err != nil {
		logs.Error("LBFTBlockForEncode error:", err)
		return nil
	}
	if !keepSeal {
		extra.Seal = []byte{}
	}
	extra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(extra)
	if err != nil {
		return nil
	}

	newBlock.ExtraData = append(newBlock.ExtraData[:LBFTExtraVanity], payload...)

	return &newBlock
}

type LBFTExtra struct {
	Validators    []common.Address
	Seal          []byte
	CommittedSeal [][]byte
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (le *LBFTExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		le.Validators,
		le.Seal,
		le.CommittedSeal,
	})
}

// DecodeRLP implements rlp.Decoder, and load the lbft fields from a RLP stream.
func (le *LBFTExtra) DecodeRLP(s *rlp.Stream) error {
	var lbftExtra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
	}
	if err := s.Decode(&lbftExtra); err != nil {
		return err
	}
	le.Validators, le.Seal, le.CommittedSeal = lbftExtra.Validators, lbftExtra.Seal, lbftExtra.CommittedSeal
	return nil
}

func LBFTBlockExtra(block *CBlock) (*LBFTExtra, error) {
	if len(block.ExtraData) < LBFTExtraVanity {
		return nil, ErrInvalidLBFTBlockExtra
	}

	var le *LBFTExtra
	err := rlp.DecodeBytes(block.ExtraData[LBFTExtraVanity:], &le)
	if err != nil {
		return nil, err
	}
	return le, nil
}
