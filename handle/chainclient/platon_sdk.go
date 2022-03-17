package chainclient

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// BlockNonce is an 81-byte vrf proof containing random numbers
// Used to verify the block when receiving the block
type BlockNonce [81]byte

// EncodeNonce converts the given byte to a block nonce.
func EncodeNonce(v []byte) BlockNonce {
	var n BlockNonce
	copy(n[:], v)
	return n
}

func (n BlockNonce) Bytes() []byte {
	return n[:]
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

//go:generate gencodec -type Header -field-override headerMarshaling -out gen_Platon_header_json.go

// PlatonHeader represents a block header in the Ethereum core.
type PlatonHeader struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	//Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Number   *big.Int   `json:"number"           gencodec:"required"`
	GasLimit uint64     `json:"gasLimit"         gencodec:"required"`
	GasUsed  uint64     `json:"gasUsed"          gencodec:"required"`
	Time     uint64     `json:"timestamp"        gencodec:"required"`
	Extra    []byte     `json:"extraData"        gencodec:"required"`
	Nonce    BlockNonce `json:"nonce"            gencodec:"required"`

	// caches
	sealHash  atomic.Value `json:"-" rlp:"-"`
	hash      atomic.Value `json:"-" rlp:"-"`
	publicKey atomic.Value `json:"-" rlp:"-"`
}

type PlatonSdk struct {
	rpcClient *rpc.Client
	rawClient *ethclient.Client
	url       string
}

func NewPlatonSdk(url string) (*PlatonSdk, error) {
	rpcClient, err := rpc.Dial(url)
	if rpcClient == nil || err != nil {
		return nil, fmt.Errorf("rpc client works error(%s)", err)
	}
	rawClient, err := ethclient.Dial(url)
	if rawClient == nil || err != nil {
		return nil, fmt.Errorf("raw client works error(%s)", err)
	}
	return &PlatonSdk{
		rpcClient: rpcClient,
		rawClient: rawClient,
		url:       url,
	}, nil
}

func (gs *PlatonSdk) GetClient() *ethclient.Client {
	return gs.rawClient
}

func (gs *PlatonSdk) GetCurrentBlockHeight() (uint64, error) {
	var result hexutil.Big
	err := gs.rpcClient.CallContext(context.Background(), &result, "eth_blockNumber")
	for err != nil {
		return 0, err
	}
	return (*big.Int)(&result).Uint64(), err
}

func (gs *PlatonSdk) GetHeaderByNumber(number uint64) (*PlatonHeader, error) {
	var header *PlatonHeader
	err := gs.rpcClient.CallContext(context.Background(), &header, "eth_getBlockByNumber", hexutil.EncodeBig(big.NewInt(int64(number))), true)
	return header, err
}

func (gs *PlatonSdk) GetTransactionByHash(hash common.Hash) (*types.Transaction, error) {
	tx, _, err := gs.rawClient.TransactionByHash(context.Background(), hash)
	for err != nil {
		return nil, err
	}
	return tx, err
}

func (gs *PlatonSdk) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	receipt, err := gs.rawClient.TransactionReceipt(context.Background(), hash)
	for err != nil {
		return nil, err
	}
	return receipt, nil
}
