package bridge

import (
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/models"
)

func ConstructTx(tx *models.WrapperTransaction, transfer *models.SrcTransfer, proxyAddr string, tokenURI string) *TxParam {
	args := TxArgs{
		toAssetHash: common.HexToAddress(transfer.DstAsset).Bytes(),
		toAddress:   common.HexToAddress(transfer.DstUser).Bytes(),
		tokenID:     *transfer.TokenID,
		tokenURI:    []byte(tokenURI),
	}
	argsBytes := args.Serialize()
	return &TxParam{
		TxHash:       common.HexToHash(tx.Hash),
		FromChainID:  tx.SrcChainID,
		FromContract: common.HexToAddress(transfer.Asset),
		ToChainID:    tx.DstChainID,
		ToContract:   common.HexToAddress(proxyAddr),
		Args:         argsBytes,
	}
}
