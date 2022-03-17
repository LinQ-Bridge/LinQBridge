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

type GethSdk struct {
	rpcClient *rpc.Client
	rawClient *ethclient.Client
	url       string
}

func NewGethSdk(url string) (*GethSdk, error) {
	rpcClient, err := rpc.Dial(url)
	if rpcClient == nil || err != nil {
		return nil, fmt.Errorf("rpc client works error(%s)", err)
	}
	rawClient, err := ethclient.Dial(url)
	if rawClient == nil || err != nil {
		return nil, fmt.Errorf("raw client works error(%s)", err)
	}
	return &GethSdk{
		rpcClient: rpcClient,
		rawClient: rawClient,
		url:       url,
	}, nil
}

func (gs *GethSdk) GetClient() *ethclient.Client {
	return gs.rawClient
}

func (gs *GethSdk) GetCurrentBlockHeight() (uint64, error) {
	var result hexutil.Big
	err := gs.rpcClient.CallContext(context.Background(), &result, "eth_blockNumber")
	for err != nil {
		return 0, err
	}
	return (*big.Int)(&result).Uint64(), err
}

func (gs *GethSdk) GetHeaderByNumber(number uint64) (*types.Header, error) {
	header := &types.Header{}
	err := gs.rpcClient.CallContext(context.Background(), header, "eth_getBlockByNumber", hexutil.EncodeBig(big.NewInt(int64(number))), false)
	return header, err
}

func (gs *GethSdk) GetTransactionByHash(hash common.Hash) (*types.Transaction, error) {
	tx, _, err := gs.rawClient.TransactionByHash(context.Background(), hash)
	for err != nil {
		return nil, err
	}
	return tx, err
}

func (gs *GethSdk) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	receipt, err := gs.rawClient.TransactionReceipt(context.Background(), hash)
	for err != nil {
		return nil, err
	}
	return receipt, nil
}
