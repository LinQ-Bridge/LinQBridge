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

type KlayClient struct {
	client       *KlaySdk
	latestHeight uint64
}

func NewKlayClient(url string) *KlayClient {
	sdk, err := NewKlaySdk(url)
	if err != nil || sdk == nil {
		panic(err)
	}

	//height, err := sdk.GetCurrentBlockHeight()
	//if err != nil {
	//	return nil
	//}

	return &KlayClient{
		client:       sdk,
		latestHeight: 0,
	}
}

type KlaySdkPro struct {
	clients       map[string]*KlayClient
	selectionSlot uint64
	id            uint64
	mu            sync.Mutex
}

func NewKlaySdkPro(urls []string, slot uint64, id uint64) *KlaySdkPro {
	clients := make(map[string]*KlayClient)
	for _, url := range urls {
		clients[url] = NewKlayClient(url)
	}
	gsp := &KlaySdkPro{clients: clients, selectionSlot: slot, id: id}
	go gsp.NodeSelection()
	return gsp
}

func (gsp *KlaySdkPro) NodeSelection() {
	for {
		gsp.nodeSelection()
	}
}

func (gsp *KlaySdkPro) nodeSelection() {
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

func (gsp *KlaySdkPro) selection() {
	for url, klayClient := range gsp.clients {
		height, err := klayClient.client.GetCurrentBlockHeight()
		if err != nil || height == math.MaxUint64 {
			logs.Error("get current block height err: %v, url: %s", err, url)
			height = 1
		}
		gsp.mu.Lock()
		klayClient.latestHeight = height - 1
		gsp.mu.Unlock()
	}
}

func (gsp *KlaySdkPro) GetClient() *ethclient.Client {
	clients := gsp.GetLatest()
	if clients == nil {
		return nil
	}
	return clients.client.GetClient()
}

func (gsp *KlaySdkPro) GetLatest() *KlayClient {
	gsp.mu.Lock()
	defer gsp.mu.Unlock()
	height := uint64(0)
	var latestClient *KlayClient
	for _, client := range gsp.clients {
		if client != nil && client.latestHeight > height {
			height = client.latestHeight
			latestClient = client
		}
	}
	return latestClient
}

func (gsp *KlaySdkPro) GetLatestHeight() (uint64, error) {
	client := gsp.GetLatest()
	if client == nil {
		return 0, fmt.Errorf("all node is not working")
	}
	return client.latestHeight, nil
}

func (gsp *KlaySdkPro) GetHeaderByNumber(number uint64) (*KlayHeader, error) {
	klayClient := gsp.GetLatest()
	if klayClient == nil {
		return nil, fmt.Errorf("all node is not working")
	}

	for klayClient != nil {
		header, err := klayClient.client.GetHeaderByNumber(number)
		if err != nil {
			klayClient.latestHeight = 0
			klayClient = gsp.GetLatest()
		} else {
			return header, nil
		}
	}
	return nil, fmt.Errorf("all node is not working")
}

func (gsp *KlaySdkPro) GetTransactionByHash(hash common.Hash) (*types.Transaction, error) {
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

func (gsp *KlaySdkPro) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
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
