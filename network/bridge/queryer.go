package bridge

import (
	"math/big"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"land-bridge/conf"
	"land-bridge/contracts/nftquery"
)

type Queryer struct {
	urls        map[uint64]string
	nftQueries  map[uint64]string
	lockProxies map[uint64]string
}

func NewBridgeQueryer(cfg *conf.Config) *Queryer {
	var (
		urls        = make(map[uint64]string)
		nftQueries  = make(map[uint64]string)
		lockProxies = make(map[uint64]string)
	)
	for _, c := range cfg.Chains {
		urls[c.ChainID] = c.Nodes[0].URL
		nftQueries[c.ChainID] = c.NFTQueryContract
		lockProxies[c.ChainID] = c.NFTProxyContract
	}
	return &Queryer{urls, nftQueries, lockProxies}
}

func (bq *Queryer) GetTokenURIByAssetWithID(chainID uint64, asset string, tokenID *big.Int) (string, error) {
	var (
		client *ethclient.Client
		err    error
	)
	for {
		client, err = ethclient.Dial(bq.urls[chainID])
		if client != nil {
			break
		}
		if err != nil {
			logs.Warn("error when get client", err)
		}
	}

	query, err := nftquery.NewPolyNFTQuery(common.HexToAddress(bq.nftQueries[chainID]), client)
	if err != nil {
		logs.Error("error when get queryer contract instance", err)
		return "", err
	}
	_, tokenURI, err := query.GetAndCheckTokenUrl(nil, common.HexToAddress(asset), common.HexToAddress(bq.lockProxies[chainID]), tokenID)
	return tokenURI, err
}
