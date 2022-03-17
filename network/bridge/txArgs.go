package bridge

import (
	"github.com/beego/beego/v2/core/logs"
	polyCommon "github.com/polynetwork/poly/common"

	"land-bridge/models"
)

type TxArgs struct {
	toAssetHash []byte
	toAddress   []byte
	tokenID     models.BigInt
	tokenURI    []byte
}

func (a *TxArgs) Serialize() []byte {
	tokenIDBytes := a.tokenID.Bytes()
	if len(tokenIDBytes) == 0 || len(tokenIDBytes) > 32 {
		logs.Error("wrong tokenID")
	}
	var hash polyCommon.Uint256
	copy(hash[32-len(tokenIDBytes):], tokenIDBytes)
	reverse := polyCommon.ToArrayReverse(hash.ToArray())
	hashReverse, err := polyCommon.Uint256ParseFromBytes(reverse)
	if err != nil {
		logs.Error("err when decode bytes")
	}

	sink := polyCommon.NewZeroCopySink(nil)
	sink.WriteVarBytes(a.toAssetHash)
	sink.WriteVarBytes(a.toAddress)
	sink.WriteHash(hashReverse)
	sink.WriteVarBytes(a.tokenURI)
	return sink.Bytes()
}
