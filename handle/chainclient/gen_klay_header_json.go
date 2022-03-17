package chainclient

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// MarshalJSON marshals as JSON.
func (h KlayHeader) MarshalJSON() ([]byte, error) {
	return nil, errors.New("not support")
}

// UnmarshalJSON unmarshals from JSON.
func (h *KlayHeader) UnmarshalJSON(input []byte) error {
	type KlayHeader struct {
		ParentHash  *common.Hash    `json:"parentHash"       gencodec:"required"`
		Rewardbase  *common.Address `json:"reward"           gencodec:"required"`
		Root        *common.Hash    `json:"stateRoot"        gencodec:"required"`
		TxHash      *common.Hash    `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash *common.Hash    `json:"receiptsRoot"     gencodec:"required"`
		BlockScore  *hexutil.Big    `json:"blockScore"       gencodec:"required"`
		Number      *hexutil.Big    `json:"number"           gencodec:"required"`
		GasUsed     *hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
		Time        *hexutil.Big    `json:"timestamp"        gencodec:"required"`
		TimeFoS     *hexutil.Uint   `json:"timestampFoS"     gencodec:"required"`
		Extra       *hexutil.Bytes  `json:"extraData"        gencodec:"required"`
		Governance  *hexutil.Bytes  `json:"governanceData"        gencodec:"required"`
		Vote        *hexutil.Bytes  `json:"voteData,omitempty"`
	}
	var dec KlayHeader
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for KlayHeader")
	}
	h.ParentHash = *dec.ParentHash
	if dec.Rewardbase == nil {
		return errors.New("missing required field 'reward' for KlayHeader")
	}
	h.Rewardbase = *dec.Rewardbase
	if dec.Root == nil {
		return errors.New("missing required field 'stateRoot' for KlayHeader")
	}
	h.Root = *dec.Root
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionsRoot' for KlayHeader")
	}
	h.TxHash = *dec.TxHash
	if dec.ReceiptHash == nil {
		return errors.New("missing required field 'receiptsRoot' for KlayHeader")
	}
	h.ReceiptHash = *dec.ReceiptHash
	if dec.BlockScore == nil {
		return errors.New("missing required field 'blockScore' for KlayHeader")
	}
	h.BlockScore = (*big.Int)(dec.BlockScore)
	if dec.Number == nil {
		return errors.New("missing required field 'number' for KlayHeader")
	}
	h.Number = (*big.Int)(dec.Number)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for KlayHeader")
	}
	h.GasUsed = uint64(*dec.GasUsed)
	if dec.Time == nil {
		return errors.New("missing required field 'timestamp' for KlayHeader")
	}
	h.Time = (*big.Int)(dec.Time)
	if dec.TimeFoS == nil {
		return errors.New("missing required field 'timestampFoS' for KlayHeader")
	}
	h.TimeFoS = uint8(*dec.TimeFoS)
	if dec.Extra == nil {
		return errors.New("missing required field 'extraData' for KlayHeader")
	}
	h.Extra = *dec.Extra
	if dec.Governance == nil {
		return errors.New("missing required field 'governanceData' for KlayHeader")
	}
	h.Governance = *dec.Governance
	if dec.Vote != nil {
		h.Vote = *dec.Vote
	}
	return nil
}
