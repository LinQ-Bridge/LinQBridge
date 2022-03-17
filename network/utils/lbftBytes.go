package utils

import (
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TxType uint8

const (
	NFT_CROSS TxType = iota
)

var (
	LBFTExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	LBFTExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	// ErrInvalidLBFTBlockExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidLBFTBlockExtra = errors.New("invalid LBFT header extra-data")
	LBFTBytesType            = reflect.TypeOf(LBFTBytes{})
)

type LBFTBytes [22]byte

func (by *LBFTBytes) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(LBFTBytesType, input, by[:])
}
