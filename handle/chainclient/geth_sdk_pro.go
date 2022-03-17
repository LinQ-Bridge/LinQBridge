package chainclient

import (
	"fmt"
	"math"
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GethClient struct {
	client       *GethSdk
	latestHeight uint64
}

func NewGethClient(url string) *GethClient {
	sdk, err := NewGethSdk(url)
	if err != nil || sdk == nil {
		panic(err)
	}
	return &GethClient{
		client:       sdk,
		latestHeight: 0,
	}
}

type GethSdkPro struct {
	clients       map[string]*GethClient
	selectionSlot uint64
	id            uint64
	mu            sync.Mutex
}

func NewGethSdkPro(urls []string, slot uint64, id uint64) *GethSdkPro {
	clients := make(map[string]*GethClient)
	for _, url := range urls {
		clients[url] = NewGethClient(url)
	}
	gsp := &GethSdkPro{clients: clients, selectionSlot: slot, id: id}
	go gsp.NodeSelection()
	return gsp
}

func (gsp *GethSdkPro) NodeSelection() {
	for {
		gsp.nodeSelection()
	}
}

func (gsp *GethSdkPro) nodeSelection() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("node selection, recover info: %s", string(debug.Stack()))
		}
	}()
	logs.Debug("node selection of chain : %d", gsp.id)
	ticker := time.NewTicker(time.Second * time.Duration(gsp.selectionSlot))
	for {
		select {
		case <-ticker.C:
			gsp.selection()
		}
	}
}

func (gsp *GethSdkPro) selection() {
	for url, gethClient := range gsp.clients {
		height, err := gethClient.client.GetCurrentBlockHeight()
		if err != nil || height == math.MaxUint64 {
			logs.Error("get current block height err: %v, url: %s", err, url)
			height = 1
		}
		gsp.mu.Lock()
		gethClient.latestHeight = height - 1
		gsp.mu.Unlock()
	}
}

func (gsp *GethSdkPro) GetClient() *ethclient.Client {
	clients := gsp.GetLatest()
	if clients == nil {
		return nil
	}
	return clients.client.GetClient()
}

func (gsp *GethSdkPro) GetLatest() *GethClient {
	gsp.mu.Lock()
	defer gsp.mu.Unlock()
	height := uint64(0)
	var latestClient *GethClient
	for _, client := range gsp.clients {
		if client != nil && client.latestHeight > height {
			height = client.latestHeight
			latestClient = client
		}
	}
	return latestClient
}

func (gsp *GethSdkPro) GetLatestHeight() (uint64, error) {
	client := gsp.GetLatest()
	if client == nil {
		return 0, fmt.Errorf("all node is not working")
	}
	return client.latestHeight, nil
}

func (gsp *GethSdkPro) GetHeaderByNumber(number uint64) (*types.Header, error) {
	gethClient := gsp.GetLatest()
	if gethClient == nil {
		return nil, fmt.Errorf("all node is not working")
	}

	for gethClient != nil {
		header, err := gethClient.client.GetHeaderByNumber(number)
		if err != nil {
			gethClient.latestHeight = 0
			gethClient = gsp.GetLatest()
		} else {
			return header, nil
		}
	}
	return nil, fmt.Errorf("all node is not working")
}

func (gsp *GethSdkPro) GetTransactionByHash(hash common.Hash) (*types.Transaction, error) {
	clients := gsp.GetLatest()
	if clients == nil {
		return nil, fmt.Errorf("all node is not working")
	}

	for clients != nil {
		tx, err := clients.client.GetTransactionByHash(hash)
		if err != nil {
			clients.latestHeight = 0
			clients = gsp.GetLatest()
		} else {
			return tx, nil
		}
	}
	return nil, fmt.Errorf("all node is not working")
}

func (gsp *GethSdkPro) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	clients := gsp.GetLatest()
	if clients == nil {
		return nil, fmt.Errorf("all node is not working")
	}

	for clients != nil {
		receipt, err := clients.client.GetTransactionReceipt(hash)
		if err != nil {
			clients.latestHeight = 0
			clients = gsp.GetLatest()
		} else {
			return receipt, nil
		}
	}
	return nil, fmt.Errorf("all node is not working")
}
