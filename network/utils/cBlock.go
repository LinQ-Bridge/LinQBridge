package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rlp"

	"land-bridge/network/bridge"
)

type CBlocks []*CBlock

func (s CBlocks) ToBlock() []*Block {
	var res []*Block
	for _, block := range s {
		res = append(res, block.ToBlock())
	}
	return res
}

type CBlock struct {
	ParentHash common.Hash    `json:"parent_hash"`
	BlockHash  common.Hash    `json:"block_hash"`
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

	TxParam bridge.TxParam `json:"txParam"`
}

func (b *CBlock) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extcblock{
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
		TxParam:    b.TxParam,
	})
}

func (b *CBlock) DecodeRLP(s *rlp.Stream) error {
	var eb extcblock

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
	b.TxParam = eb.TxParam
	return nil
}

func (b *CBlock) String() string {
	return fmt.Sprintf("{CBlock: %v}", *b)
}

func (b *CBlock) NumberU64() uint64 {
	return b.Height
}

// SanityCheck checks a few basic things -- these checks are way beyond what
// any 'sane' production values should hold, and can mainly be used to prevent
// that the unbounded fields are stuffed with junk data to add processing
// overhead
func (b *CBlock) SanityCheck() error {
	if eLen := len(b.ExtraData); eLen > 100*1024 {
		return fmt.Errorf("too large block extradata: size %d", eLen)
	}
	return nil
}

// "external" block encoding. used for eth protocol, etc.
type extcblock struct {
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
	TxParam    bridge.TxParam
}

func (b *CBlock) Number() *big.Int {
	return big.NewInt(int64(b.Height))
}

func (b *CBlock) Hash() common.Hash {
	if b.BlockHash != emptyHash {
		return b.BlockHash
	} else {
		b.BlockHash = b.ToBlock().Hash()
		return b.BlockHash
	}
}

func (b *CBlock) TxHash() common.Hash {
	return common.BytesToHash(b.SrcTx.TxHash)
}

func (b *CBlock) GetTxParam() bridge.TxParam {
	return b.TxParam
}

func (b *CBlock) ToBlock() *Block {
	return &Block{
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

func (b *CBlock) UnmarshalJSON(input []byte) error {
	type block struct {
		ParentHash *common.Hash         `json:"parent_hash"`
		BlockHash  *common.Hash         `json:"block_hash"`
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

	if dec.TxParam != nil {
		b.TxParam = *dec.TxParam
	}

	return nil
}

func (b *CBlock) LBFTBlockExtra() (*LBFTExtra, error) {
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
