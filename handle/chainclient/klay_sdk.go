package chainclient

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

//const (
//	// BloomByteLength represents the number of bytes used in a header log bloom.
//	BloomByteLength = 256
//
//	// BloomBitLength represents the number of bits used in a header log bloom.
//	BloomBitLength = 8 * BloomByteLength
//)
//
//// Bloom represents a 2048 bit bloom filter.
//type Bloom [BloomByteLength]byte

//go:generate gencodec -type Header -field-override headerMarshaling -out gen_klay_header_json.go

// KlayHeader represents a block header in the Klaytn core.
type KlayHeader struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	Rewardbase  common.Address `json:"reward"           gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	//Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	BlockScore *big.Int `json:"blockScore"       gencodec:"required"`
	Number     *big.Int `json:"number"           gencodec:"required"`
	GasUsed    uint64   `json:"gasUsed"          gencodec:"required"`
	Time       *big.Int `json:"timestamp"        gencodec:"required"`
	// TimeFoS represents a fraction of a second since `Time`.
	TimeFoS    uint8  `json:"timestampFoS"     gencodec:"required"`
	Extra      []byte `json:"extraData"        gencodec:"required"`
	Governance []byte `json:"governanceData"        gencodec:"required"`
	Vote       []byte `json:"voteData,omitempty"`
}

type KlaySdk struct {
	rpcClient *rpc.Client
	rawClient *ethclient.Client
	url       string
}

func NewKlaySdk(url string) (*KlaySdk, error) {
	rpcClient, err := rpc.Dial(url)
	if rpcClient == nil || err != nil {
		return nil, fmt.Errorf("rpc client works error(%s)", err)
	}
	rawClient, err := ethclient.Dial(url)
	if rawClient == nil || err != nil {
		return nil, fmt.Errorf("raw client works error(%s)", err)
	}
	return &KlaySdk{
		rpcClient: rpcClient,
		rawClient: rawClient,
		url:       url,
	}, nil
}

func (gs *KlaySdk) GetClient() *ethclient.Client {
	return gs.rawClient
}

func (gs *KlaySdk) GetCurrentBlockHeight() (uint64, error) {
	var result hexutil.Big
	err := gs.rpcClient.CallContext(context.Background(), &result, "klay_blockNumber")
	for err != nil {
		return 0, err
	}
	return (*big.Int)(&result).Uint64(), err
}

func (gs *KlaySdk) GetHeaderByNumber(number uint64) (*KlayHeader, error) {
	var header *KlayHeader
	err := gs.rpcClient.CallContext(context.Background(), &header, "klay_getBlockByNumber", hexutil.EncodeBig(big.NewInt(int64(number))), true)
	return header, err
}

func (gs *KlaySdk) GetTransactionByHash(hash common.Hash) (*types.Transaction, error) {
	tx, _, err := gs.rawClient.TransactionByHash(context.Background(), hash)
	for err != nil {
		return nil, err
	}
	return tx, err
}

func (gs *KlaySdk) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	receipt, err := gs.rawClient.TransactionReceipt(context.Background(), hash)
	for err != nil {
		return nil, err
	}
	return receipt, nil
}
