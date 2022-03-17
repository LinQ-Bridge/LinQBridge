package bridge

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	polyCommon "github.com/polynetwork/poly/common"
)

type TxParam struct {
	TxHash       common.Hash
	FromChainID  uint64
	FromContract common.Address
	ToChainID    uint64
	ToContract   common.Address
	Args         []byte
	Signatures   []byte
}

func (p *TxParam) GetTxHash() common.Hash {
	return p.TxHash
}

func (p *TxParam) ChainID() uint64 {
	return p.ToChainID
}

func (p *TxParam) SetSignatures(signatures []byte) {
	p.Signatures = signatures
}

func (p *TxParam) GetSignatures() []byte {
	return p.Signatures
}

func (p *TxParam) Serialize() []byte {
	sink := polyCommon.NewZeroCopySink(nil)
	sink.WriteHash(polyCommon.Uint256(p.TxHash))
	sink.WriteUint64(p.FromChainID)
	sink.WriteVarBytes(p.FromContract.Bytes())
	sink.WriteUint64(p.ToChainID)
	sink.WriteVarBytes(p.ToContract.Bytes())
	sink.WriteVarBytes(p.Args)
	return sink.Bytes()
}

func (p *TxParam) Hash() []byte {
	return crypto.Keccak256(p.Serialize())
}

func (p *TxParam) UnmarshalJSON(input []byte) error {
	type txInfo struct {
		TxHash       *common.Hash         `json:"tx_hash"`
		FromChainID  *math.HexOrDecimal64 `json:"from_chain_id"`
		FromContract *common.Address      `json:"from_contract"`
		ToChainID    *math.HexOrDecimal64 `json:"to_chain_id"`
		ToContract   *common.Address      `json:"to_contract"`
		Args         *hexutil.Bytes       `json:"Args"`
		Signatures   *hexutil.Bytes       `json:"Signatures"`
	}

	var dec txInfo
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.TxHash != nil {
		p.TxHash = *dec.TxHash
	}

	if dec.FromChainID == nil {
		return errors.New("missing required field 'FromChainID' for Genesis")
	}
	p.FromChainID = uint64(*dec.FromChainID)

	if dec.FromContract != nil {
		p.FromContract = *dec.FromContract
	}

	if dec.ToChainID == nil {
		return errors.New("missing required field 'ToChainID' for Genesis")
	}
	p.ToChainID = uint64(*dec.ToChainID)

	if dec.ToContract != nil {
		p.ToContract = *dec.ToContract
	}

	if dec.Args != nil {
		p.Args = *dec.Args
	}

	if dec.Signatures != nil {
		p.Signatures = *dec.Signatures
	}
	return nil
}
