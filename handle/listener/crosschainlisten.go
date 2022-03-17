package listener

import (
	"fmt"
	"math"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"land-bridge/conf"
	"land-bridge/handle/dao"
	"land-bridge/handle/listener/gethlisten"
	"land-bridge/handle/listener/klaylisten"
	"land-bridge/handle/listener/platonlisten"
	"land-bridge/models"
)

var chainListens [12]*ChainListen

func StartCrossChainListen(cfg []*conf.ChainListenConfig, dbCfg *conf.DBConfig) {
	dao := dao.NewBridgeDao(dbCfg)
	if dao == nil {
		panic("sql server is invalid")
	}
	for i, clc := range cfg {
		core := NewChainListenCore(clc)
		if core == nil {
			panic(fmt.Sprintf("chain %s(%d) is not supported", clc.ChainName, clc.ChainID))
		}
		chainListen := NewChainListen(core, dao)
		chainListen.Start()
		chainListens[i] = chainListen
	}
}

func StopCrossChainListen() {
	for _, chainListen := range chainListens {
		if chainListen != nil {
			chainListen.Stop()
		}
	}
}

type ChainListenCore interface {
	GetChainName() string
	GetChainID() uint64
	GetChainListenSlot() uint64
	GetBatchSize() uint64
	GetDefer() uint64
	GetLatestHeight() (uint64, error)
	HandleNewBlock(height uint64) ([]*models.WrapperTransaction, []*models.SrcTransaction, []*models.DstTransaction, int, int, error)
}

func NewChainListenCore(clc *conf.ChainListenConfig) ChainListenCore {
	switch clc.ChainID {
	case conf.ETHEREUM:
		return gethlisten.NewGethChainListen(clc)
	case conf.BSC:
		return gethlisten.NewGethChainListen(clc)
	case conf.KLAYTN:
		return klaylisten.NewKlayChainListen(clc)
	case conf.PLATON:
		return platonlisten.NewPlatonChainListen(clc)
	default:
		return nil
	}
}

type ChainListen struct {
	core   ChainListenCore
	db     *dao.BridgeDao
	height uint64
	exit   chan bool
}

func NewChainListen(core ChainListenCore, db *dao.BridgeDao) *ChainListen {
	return &ChainListen{
		core: core,
		db:   db,
		exit: make(chan bool, 0),
	}
}

func (cl *ChainListen) Start() {
	logs.Info("start listen: %s(%d)", cl.core.GetChainName(), cl.core.GetChainID())
	go cl.ListenChain()
}

func (cl *ChainListen) Stop() {
	cl.exit <- true
}

func (cl *ChainListen) ListenChain() {
	for {
		exit := cl.listenChain()
		if exit {
			close(cl.exit)
			break
		}
		time.Sleep(time.Second * 3)
	}
}

func (cl *ChainListen) HandleNewBlock(height uint64) (w []*models.WrapperTransaction, s []*models.SrcTransaction, d []*models.DstTransaction, err error) {
	chainID := cl.core.GetChainID()
	var locks, unlocks int
	defer func() {
		if r := recover(); r != nil {
			logs.Error("Possible inconsistent chain %d height %d wrapper %d/%d src %d/%d dst %d/%d", chainID, height, len(w), locks, len(s), locks, len(d), unlocks, "error:", r)
		}
	}()
	for c := 3; c > 0; c-- {
		w, s, d, locks, unlocks, err = cl.core.HandleNewBlock(height)
		if err != nil {
			return
		}
		if locks == len(s) && unlocks == len(d) {
			return
		}
		if c > 1 {
			logs.Warn("Possible missing events for chain %d height %d", chainID, height)
			time.Sleep(time.Second * 5)
		}
	}

	return
}

func (cl *ChainListen) listenChain() (exit bool) {
	defer func() {
		if r := recover(); r != nil {
			//logs.Error("service start, recover info: %s", r)
			exit = false
		}
	}()
	chain, err := cl.db.GetChain(cl.core.GetChainID())
	if err != nil {
		panic(err)
	}
	lastHeight, err := cl.core.GetLatestHeight()
	if err != nil || lastHeight == 0 {
		panic(err)
	}
	if chain.Height == 0 {
		chain.Height = lastHeight
	}
	cl.db.UpdateChain(chain)
	if cl.height != 0 {
		chain.Height = cl.height
	}
	timedelay := time.Second
	ticker := time.NewTimer(timedelay)
	for {
		select {
		case <-ticker.C:
			height, err := cl.core.GetLatestHeight()
			if err != nil || height == 0 || height == math.MaxUint64 {
				logs.Error("listenChain - cannot get chain %s height, err: %s", cl.core.GetChainName(), err)
				ticker.Reset(timedelay)
				continue
			}
			if chain.Height >= height-cl.core.GetDefer() {
				ticker.Reset(timedelay)
				continue
			}
			//logs.Info("ListenChain - chain %s latest height is %d, listen height: %d", cl.core.GetChainName(), height, chain.Height)

			for chain.Height < height-cl.core.GetDefer() {
				batchSize := cl.core.GetBatchSize()
				if batchSize == 0 {
					batchSize = 1
				}
				if batchSize > height-chain.Height-cl.core.GetDefer() {
					batchSize = height - chain.Height - cl.core.GetDefer()
				}
				ch := make(chan bool, batchSize)
				for i := uint64(1); i <= batchSize; i++ {
					go func(height uint64) {
						select {
						case ch <- cl.listenChainHeight(height):
						case <-time.NewTicker(3 * time.Second).C:
							logs.Error("listen time ticker may too short")
							ch <- false
						}
					}(chain.Height + i)
				}
				allTaskSuccess := true
				for j := 0; j < int(batchSize); j++ {
					ok := <-ch
					if !ok {
						allTaskSuccess = false
					}
				}
				close(ch)
				if !allTaskSuccess {
					break
				}

				chain.Height += batchSize
				if err := cl.db.UpdateChain(chain); err != nil {
					logs.Error("UpdateChain [chainID:%d, height:%d] err %v", chain.ChainID, chain.Height, err)
					chain.Height -= batchSize
				}
			}
			ticker.Reset(timedelay)
		case <-cl.exit:
			logs.Info("cross chain listen exit, chain: %s(%d)", cl.core.GetChainName(), cl.core.GetChainID())
			exit = true
			return
		}
	}
}

func (cl *ChainListen) listenChainHeight(height uint64) bool {
	wrapperTransactions, srcTransactions, dstTransactions, err := cl.HandleNewBlock(height)
	if err != nil {
		logs.Error("HandleNewBlock %d err: %v", height, err)
		return false
	}
	err = cl.db.UpdateEvents(wrapperTransactions, srcTransactions, dstTransactions)
	if err != nil {
		logs.Error("UpdateEvents on block %d err: %v", height, err)
		return false
	}
	return true
}
