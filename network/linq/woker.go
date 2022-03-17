package linq

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"land-bridge/network/bridge"
	"land-bridge/network/consensus"
	"land-bridge/network/linq/txblock"
	"land-bridge/network/utils"
)

const (
	staleThreshold = 7
)

// newWorkReq represents a request for new sealing work submitting with relative interrupt notifier.
type newWorkReq struct {
	timestamp int64
	tx        *txblock.TxInfo
}

type task struct {
	block    *utils.CBlock
	createAt time.Time
}

type worker struct {
	engine   consensus.Engine
	coinbase common.Address

	chain  *Store
	bridge *bridge.Bridge
	pool   *txblock.BlockPool

	// Channels
	mux         *event.TypeMux
	taskCh      chan *task
	startCh     chan struct{}
	exitCh      chan struct{}
	newWorkCh   chan *newWorkReq
	resultCh    chan *utils.Block
	chainHeadCh chan ChainHeadEvent

	mu           sync.RWMutex // The lock used to protect the coinbase and extra fields
	pendingMu    sync.RWMutex
	pendingTasks map[common.Hash]*task

	// atomic status counters
	running int32 // The indicator whether the consensus engine is running or not.
	newTxs  int32 // New arrival transaction count since last sealing work submitting.

	snapshotMu    sync.RWMutex // The lock used to protect the block snapshot and state snapshot
	snapshotBlock *types.Block
}

func newWorker(chain *Store, bridge *bridge.Bridge, pool *txblock.BlockPool, mux *event.TypeMux, engine consensus.Engine, coinbase common.Address) *worker {
	worker := &worker{
		engine:       engine,
		coinbase:     coinbase,
		chain:        chain,
		bridge:       bridge,
		pool:         pool,
		mux:          mux,
		newWorkCh:    make(chan *newWorkReq),
		resultCh:     make(chan *utils.Block, 10),
		taskCh:       make(chan *task),
		exitCh:       make(chan struct{}),
		startCh:      make(chan struct{}, 1),
		pendingTasks: make(map[common.Hash]*task),
		chainHeadCh:  make(chan ChainHeadEvent),
	}

	go worker.mainLoop()
	go worker.listenLoop()
	go worker.newWorkLoop()
	go worker.taskLoop()
	go worker.resultLoop()

	return worker
}

// isRunning returns an indicator whether worker is running or not.
func (w *worker) isRunning() bool {
	return atomic.LoadInt32(&w.running) == 1
}

func (w *worker) start() {
	atomic.StoreInt32(&w.running, 1)
	if engine, ok := w.engine.(consensus.LBFT); ok {
		engine.Start(w.chain, w.chain.CurrentBlock)
	}
	w.startCh <- struct{}{}
}

func (w *worker) mainLoop() {
	for {
		select {
		case req := <-w.newWorkCh:
			w.commitNewWork(req.timestamp, req.tx)
		case <-w.exitCh:
			return
		}
	}
}

func (w *worker) listenLoop() {
	td := time.Second
	t := time.NewTimer(td)
	for {
		select {
		case <-t.C:
			txs, err := w.chain.PendingTxs()
			if err != nil {
				logs.Error("listenLoop error", err)
				t.Reset(td)
				continue
			}
			for _, tx := range txs {
				if w.chain.CheckTxHash(tx.Hash, tx.SrcChainID) {
					if err := w.bridge.PendingWrapperSkip(tx); err != nil {
						logs.Error("add skip transaction to Pending")
						continue
					}
				} else {
					var sign *bridge.TxParam

					if sign, err = w.bridge.BridgeMakeTx(tx); err != nil {
						logs.Error("Failed to sign transaction", err)
						t.Reset(td)
						continue
					}

					if err := w.bridge.PendingWrapper(tx); err != nil {
						logs.Error("add transaction to Pending")
						t.Reset(td)
						continue
					}

					w.pool.Push(tx.Hash, &txblock.TxInfo{
						W:       tx,
						TxParam: sign,
					})

				}
			}
			t.Reset(td)
		case <-w.exitCh:
			return
		}
	}
}

func (w *worker) taskLoop() {
	var (
		stopCh chan struct{}
		prev   common.Hash
	)

	// interrupt aborts the in-flight sealing task.
	interrupt := func() {
		if stopCh != nil {
			close(stopCh)
			stopCh = nil
		}
	}

	for {
		select {
		case task := <-w.taskCh:
			SealHash := w.engine.SealHash(task.block.ToBlock())
			if SealHash == prev {
				continue
			}
			interrupt()
			stopCh, prev = make(chan struct{}), SealHash

			w.pendingMu.Lock()
			w.pendingTasks[SealHash] = task
			w.pendingMu.Unlock()

			// consensus begin
			if err := w.engine.Seal(w.chain, task.block, w.resultCh, stopCh); err != nil {
				logs.Warn("Block sealing failed, err: %v", err)
			}
		case <-w.exitCh:
			interrupt()
			return
		}
	}
}

// newWorkLoop is a standalone goroutine to submit new mining work upon received events.
func (w *worker) newWorkLoop() {
	var (
		timestamp int64 // timestamp for each round of mining.
	)
	recommit := 3 * time.Second
	timer := time.NewTimer(recommit)
	<-timer.C

	// commit aborts in-flight transaction execution with given signal and resubmits a new one.
	commit := func() {
		if tx := w.pool.First(); tx != nil {
			atomic.StoreInt32(&w.newTxs, 0)
			logs.Trace("pool get", tx.W.Hash)
			timestamp = time.Now().Unix()
			w.newWorkCh <- &newWorkReq{timestamp: timestamp, tx: tx}
		} else {
			atomic.StoreInt32(&w.newTxs, 1)
		}
		timer.Reset(recommit)
	}
	// clearPending cleans the stale pending tasks.
	clearPending := func(number uint64) {
		w.pendingMu.Lock()
		for h, t := range w.pendingTasks {
			if t.block.NumberU64()+staleThreshold <= number {
				delete(w.pendingTasks, h)
			}
		}
		w.pendingMu.Unlock()
	}

	for {
		select {
		case <-w.startCh:
			clearPending(w.chain.CurrentBlock().NumberU64())
			commit()
		case block := <-w.chainHeadCh:
			// find txs
			if h, ok := w.engine.(consensus.Handler); ok {
				h.NewBlock()
			} else {
				panic("work engine error")
			}
			clearPending(block.Block.NumberU64())
			commit()
		case <-timer.C:
			if w.isRunning() {
				if atomic.LoadInt32(&w.newTxs) == 0 {
					timer.Reset(recommit)
					continue
				}
				commit()
			}
		case <-w.exitCh:
			return
		}
	}
}

// commitNewWork generates several new sealing tasks based on the parent block.
func (w *worker) commitNewWork(timestamp int64, tx *txblock.TxInfo) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	parent := w.chain.CurrentBlock()

	if parent.Time >= uint64(timestamp) {
		timestamp = int64(parent.Time + 1)
	}

	height := parent.Height
	block := &utils.CBlock{
		ParentHash: parent.BlockHash,
		Height:     height + 1,
		Time:       uint64(timestamp),
	}

	// Only set the coinbase if our consensus engine is running (avoid spurious block rewards)
	if w.isRunning() {
		if w.coinbase == (common.Address{}) {
			logs.Error("Refusing to mine without etherbase")
			return
		}
		block.Coinbase = w.coinbase
	}

	if err := w.engine.Prepare(w.chain, block); err != nil {
		logs.Error("Failed to prepare header for mining, err: %v", err)
		return
	}

	logs.Trace("fetch transaction: ", tx.W.Hash)

	w.commitTransaction(tx, block)
	w.commit(block)
}

func (w *worker) commitTransaction(tx *txblock.TxInfo, block *utils.CBlock) {
	block.SrcTx.ChainID = tx.W.SrcChainID
	block.SrcTx.TxHash = common.FromHex(tx.W.Hash)
	block.TxParam = *tx.TxParam
	block.Type = utils.NFT_CROSS
}

func (w *worker) commit(block *utils.CBlock) error {
	if w.isRunning() {
		select {
		case w.taskCh <- &task{block: block, createAt: time.Now()}:
			logs.Info("Commit new work, number: %d, sealHash: %s", block.Height, w.engine.SealHash(block.ToBlock()))
		case <-w.exitCh:
			logs.Info("Worker has exited")
		}
	}
	return nil
}

// resultLoop is a standalone goroutine to handle sealing result submitting
// and flush relative data to the database.
func (w *worker) resultLoop() {
	for {
		select {
		case block := <-w.resultCh:

			// Short circuit when receiving empty result.
			if block == nil {
				continue
			}

			// Short circuit when receiving duplicate result caused by resubmitting.
			if w.chain.HasBlock(block.Hash(), block.NumberU64()) {
				continue
			}
			var (
				sealhash = w.engine.SealHash(block)
				hash     = block.Hash()
			)
			w.pendingMu.RLock()
			_, exist := w.pendingTasks[sealhash]
			w.pendingMu.RUnlock()
			if !exist {
				logs.Error("Block found but no relative pending task", "number", block.Number(), "sealhash", sealhash, "hash", hash)
				continue
			}

			blocks := []*utils.Block{
				block,
			}
			num, err := w.chain.InsertChain(blocks)
			if err != nil {
				return
			}

			logs.Info("Successfully sealed new block", "number", block.Number(), "sealhash", sealhash, "hash", hash, "success", num)

			// Broadcast the block and announce chain insertion event
			w.mux.Post(NewMinedBlockEvent{Block: block})

		case <-w.exitCh:
			return
		}
	}
}
