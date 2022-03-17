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

type PlatonClient struct {
	client       *PlatonSdk
	latestHeight uint64
}

func NewPlatonClient(url string) *PlatonClient {
	sdk, err := NewPlatonSdk(url)
	if err != nil || sdk == nil {
		panic(err)
	}
	return &PlatonClient{
		client:       sdk,
		latestHeight: 0,
	}
}

type PlatonSdkPro struct {
	clients       map[string]*PlatonClient
	selectionSlot uint64
	id            uint64
	mu            sync.Mutex
}

func NewPlatonSdkPro(urls []string, slot uint64, id uint64) *PlatonSdkPro {
	clients := make(map[string]*PlatonClient)
	for _, url := range urls {
		clients[url] = NewPlatonClient(url)
	}
	gsp := &PlatonSdkPro{clients: clients, selectionSlot: slot, id: id}
	go gsp.NodeSelection()
	return gsp
}

func (gsp *PlatonSdkPro) NodeSelection() {
	for {
		gsp.nodeSelection()
	}
}

func (gsp *PlatonSdkPro) nodeSelection() {
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

func (gsp *PlatonSdkPro) selection() {
	for url, PlatonClient := range gsp.clients {
		height, err := PlatonClient.client.GetCurrentBlockHeight()
		if err != nil || height == math.MaxUint64 {
			logs.Error("get current block height err: %v, url: %s", err, url)
			height = 1
		}
		gsp.mu.Lock()
		PlatonClient.latestHeight = height - 1
		gsp.mu.Unlock()
	}
}

func (gsp *PlatonSdkPro) GetClient() *ethclient.Client {
	clients := gsp.GetLatest()
	if clients == nil {
		return nil
	}
	return clients.client.GetClient()
}

func (gsp *PlatonSdkPro) GetLatest() *PlatonClient {
	gsp.mu.Lock()
	defer gsp.mu.Unlock()
	height := uint64(0)
	var latestClient *PlatonClient
	for _, client := range gsp.clients {
		if client != nil && client.latestHeight > height {
			height = client.latestHeight
			latestClient = client
		}
	}
	return latestClient
}

func (gsp *PlatonSdkPro) GetLatestHeight() (uint64, error) {
	client := gsp.GetLatest()
	if client == nil {
		return 0, fmt.Errorf("all node is not working")
	}
	return client.latestHeight, nil
}

func (gsp *PlatonSdkPro) GetHeaderByNumber(number uint64) (*PlatonHeader, error) {
	PlatonClient := gsp.GetLatest()
	if PlatonClient == nil {
		return nil, fmt.Errorf("all node is not working")
	}

	for PlatonClient != nil {
		header, err := PlatonClient.client.GetHeaderByNumber(number)
		if err != nil {
			PlatonClient.latestHeight = 0
			PlatonClient = gsp.GetLatest()
		} else {
			return header, nil
		}
	}
	return nil, fmt.Errorf("all node is not working")
}

func (gsp *PlatonSdkPro) GetTransactionByHash(hash common.Hash) (*types.Transaction, error) {
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

func (gsp *PlatonSdkPro) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
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
