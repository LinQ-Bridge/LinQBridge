package utils

import (
	"encoding/hex"
	"strings"

	"land-bridge/contracts/nftlp"
	"land-bridge/contracts/nftwrap"
	"land-bridge/models"
)

const (
	Crosschainlock   = "CrossChainLockEvent"
	Crosschainunlock = "CrossChainUnlockEvent"
	Lock             = "LockEvent"
	Unlock           = "UnlockEvent"
)

func WrapLockEvent2WrapTx(evt *nftwrap.PolyNFTWrapperPolyWrapperLock) *models.WrapperTransaction {
	return &models.WrapperTransaction{
		Hash:         evt.Raw.TxHash.String()[2:],
		User:         strings.ToLower(evt.Sender.String()[2:]),
		DstChainID:   evt.ToChainId,
		DstUser:      strings.ToLower(evt.ToAddress.String()[2:]),
		FeeTokenHash: strings.ToLower(evt.FeeToken.String()[2:]),
		FeeAmount:    models.NewBigInt(evt.Fee),
		ServerID:     evt.Id.Uint64(),
		BlockHeight:  evt.Raw.BlockNumber,
		Standard:     models.TokenTypeErc721,
	}
}

func WrapSpeedUpEvent2WrapTx(evt *nftwrap.PolyNFTWrapperPolyWrapperSpeedUp) *models.WrapperTransaction {
	return &models.WrapperTransaction{
		Hash:         evt.TxHash.String(),
		User:         evt.Sender.String(),
		FeeTokenHash: evt.FeeToken.String(),
		FeeAmount:    models.NewBigInt(evt.Efee),
		Standard:     models.TokenTypeErc721,
	}
}

func ConvertLockProxyEvent(evt *nftlp.PolyNFTLockProxyLockEvent) *models.ProxyLockEvent {
	return &models.ProxyLockEvent{
		Method:        Lock,
		TxHash:        evt.Raw.TxHash.String()[2:],
		FromAddress:   evt.FromAddress.String()[2:],
		FromAssetHash: strings.ToLower(evt.FromAssetHash.String()[2:]),
		ToChainID:     uint32(evt.ToChainId),
		ToAssetHash:   hex.EncodeToString(evt.ToAssetHash),
		ToAddress:     hex.EncodeToString(evt.ToAddress),
		TokenID:       evt.TokenId,
	}
}

func ConvertUnlockProxyEvent(evt *nftlp.PolyNFTLockProxyUnlockEvent) *models.ProxyUnlockEvent {
	return &models.ProxyUnlockEvent{
		Method:      Unlock,
		TxHash:      evt.Raw.TxHash.String()[2:],
		ToAssetHash: strings.ToLower(evt.ToAssetHash.String()[2:]),
		ToAddress:   strings.ToLower(evt.ToAddress.String()[2:]),
		TokenID:     evt.TokenId,
	}
}
