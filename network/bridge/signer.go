package bridge

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer struct {
	priv *ecdsa.PrivateKey
	addr *common.Address
}

func NewSigner(priv *ecdsa.PrivateKey, addr *common.Address) *Signer {
	return &Signer{
		priv: priv,
		addr: addr,
	}
}

func (val *Signer) Sign(tx *TxParam) ([]byte, error) {
	hash := tx.Hash()
	return crypto.Sign(hash, val.priv)
}
