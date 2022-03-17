package linq

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"gorm.io/gorm"

	"land-bridge/constant"
	"land-bridge/models"
	"land-bridge/network/linq/txblock"
	"land-bridge/network/utils"
)

func NewLinQStore(db *gorm.DB, pool *txblock.BlockPool) *Store {
	ls := &Store{
		db:   db,
		pool: pool,
	}

	ls.UpdateCurrentBlock()

	return ls
}

type Store struct {
	currentBlock atomic.Value
	db           *gorm.DB
	pool         *txblock.BlockPool

	// procInterrupt must be atomically called
	procInterrupt int32          // interrupt signaler for block processing
	wg            sync.WaitGroup //
	chainmu       sync.RWMutex   // blockchain insertion lock

	ChainHeadCh chan<- ChainHeadEvent
}

func (ls *Store) SetHeadCh(ch chan<- ChainHeadEvent) {
	ls.ChainHeadCh = ch
}

func (ls *Store) HasBlock(hash common.Hash, number uint64) bool {
	if ls.GetBlock(hash, number) != nil {
		return true
	}
	return false
}

func (ls *Store) InsertChain(blocks utils.Blocks) (int, error) {
	if len(blocks) == 0 {
		return 0, nil
	}

	var block, prev *utils.Block

	for i := 1; i < len(blocks); i++ {
		block = blocks[i]
		prev = blocks[i-1]
		if block.NumberU64() != prev.NumberU64()+1 || block.ParentHash != prev.Hash() {
			return 0, fmt.Errorf("non contiguous insert: item %d is #%d [%x…], item %d is #%d [%x…] (parent [%x…])", i-1, prev.NumberU64(),
				prev.Hash().Bytes()[:4], i, block.NumberU64(), block.Hash().Bytes()[:4], block.ParentHash.Bytes()[:4])
		}
	}

	ls.wg.Add(1)
	ls.chainmu.Lock()
	rblock, n, err := ls.insertChain(blocks, true)
	ls.chainmu.Unlock()
	ls.wg.Done()

	if n > 0 {
		ls.ChainHeadCh <- ChainHeadEvent{
			rblock,
		}
	}

	return n, err
}

func (ls *Store) insertChain(chain utils.Blocks, verifySeals bool) (*utils.Block, int, error) {

	if atomic.LoadInt32(&ls.procInterrupt) == 1 {
		return nil, 0, nil
	}

	sort.Slice(chain, func(i, j int) bool {
		return chain[i].Height < chain[j].Height
	})

	n := 0

	tx := ls.db.Begin()
	for _, block := range chain {

		ls.pool.Delete(common.Bytes2Hex(block.SrcTx.TxHash))

		if ls.CheckBlock(block.BlockHash, block.Height) {
			continue
		}

		n++

		model := &models.Block{}

		buf := new(bytes.Buffer)
		err := rlp.Encode(buf, &block)
		if err != nil {
			tx.Rollback()
			logs.Error("insertChain encode err:", err)
			return nil, 0, err
		}
		model.Bytes = common.Bytes2Hex(buf.Bytes())
		model.Height = block.Height
		model.BlockHash = block.BlockHash.Hex()

		if err := tx.Create(model).Error; err != nil {
			tx.Rollback()
			logs.Error("insertChain create err:", err)
			return nil, 0, err
		}

		txhash := &models.TxHashHistory{}
		txhash.ChainID = block.SrcTx.ChainID
		txhash.TxHash = common.Bytes2Hex(block.SrcTx.TxHash)

		if err := tx.Create(txhash).Error; err != nil {
			tx.Rollback()
			logs.Error("insertTxHash create err:", err)
			return nil, 0, err
		}
	}
	tx.Commit()
	ls.UpdateCurrentBlock()

	return chain[len(chain)-1], n, nil
}

func (ls *Store) Rollback(chain []common.Hash) {
	ls.chainmu.Lock()
	defer ls.chainmu.Unlock()

	for i := len(chain) - 1; i >= 0; i-- {
		hash := chain[i]
		if currentBlock := ls.CurrentBlock(); currentBlock.Hash() == hash {
			newBlock := ls.GetBlock(currentBlock.ParentHash, currentBlock.NumberU64()-1)
			logs.Error("currentBlock Rollback", currentBlock.NumberU64()-1)
			ls.currentBlock.Store(newBlock)
		}
	}
}

// Config retrieves the blockchain's chain configuration.
func (ls *Store) Config() *params.ChainConfig {
	return nil
}

// CurrentBlock retrieves the current block from the local chain.
func (ls *Store) CurrentBlock() *utils.Block {
	return ls.currentBlock.Load().(*utils.Block)
}

func (ls *Store) UpdateCurrentBlock() {
	block := &models.Block{}
	tx := ls.db.Last(block)
	if tx.RowsAffected == 1 {
		var b *utils.Block
		if len(block.Bytes) > 0 {
			err := rlp.DecodeBytes(common.Hex2Bytes(block.Bytes), &b)
			if err != nil {
				panic("Error bytes decode " + err.Error())
			}
		}
		ls.currentBlock.Store(b)
		logs.Trace("Update Current Block, now is", b.Height, b.Hash().Hex())
	} else {
		panic("Table 'contract.blocks' doesn't exist or empty rows")
	}
}

// GetBlock retrieves a block from the database by hash and number.
func (ls *Store) GetBlock(hash common.Hash, number uint64) *utils.Block {
	block := new(models.Block)
	tx := ls.db.Where("block_hash = ? and height = ?", hash.Hex(), number).First(&block)
	if tx.Error != nil && tx.RowsAffected > 0 {
		return nil
	}
	var b *utils.Block
	err := rlp.DecodeBytes(common.Hex2Bytes(block.Bytes), &b)
	if err != nil {
		//logs.Error("block rlp change error ", err, "block height: ", block.Height)
		return nil
	}
	return b
}

func (ls *Store) CheckBlock(hash common.Hash, number uint64) bool {
	block := new(models.Block)
	tx := ls.db.Where("block_hash = ? and height = ?", hash.Hex(), number).First(&block)

	return tx.Error == nil
}

// GetBlockByNumber retrieves a block from the database by number.
func (ls *Store) GetBlockByNumber(number uint64) *utils.Block {
	block := &models.Block{Height: number}
	if ls.db.Where(block).First(&block).Error != nil {
		return nil
	}
	var b *utils.Block
	err := rlp.DecodeBytes(common.Hex2Bytes(block.Bytes), &b)
	if err != nil {
		return nil
	}
	return b
}

// GetBlockByHash retrieves a block from the database by its hash.
func (ls *Store) GetBlockByHash(hash common.Hash) *utils.Block {
	block := new(models.Block)
	if ls.db.Model(&block).Where("block_hash = ?", hash.Hex()).First(&block).Error != nil {
		return nil
	}
	var b *utils.Block
	err := rlp.DecodeBytes(common.Hex2Bytes(block.Bytes), &b)
	if err != nil {
		return nil
	}
	return b
}

// PendingTxs retrieves the pendingTx to deal.
func (ls *Store) PendingTxs() ([]*models.WrapperTransaction, error) {
	var wrapperTransaction []*models.WrapperTransaction
	err := ls.db.Where("status = ?", constant.STATE_SOURCE_DONE).Find(&wrapperTransaction).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return wrapperTransaction, nil
}

func (ls *Store) CheckTxHash(hashStr string, chainID uint64) bool {
	txHashHistory := &models.TxHashHistory{}
	err := ls.db.Model(txHashHistory).Where("tx_hash = ? and chain_id = ?", hashStr, chainID).First(txHashHistory).Error
	return err == nil
}

// StopInsert interrupts all insertion methods, causing them to return
// errInsertionInterrupted as soon as possible. Insertion is permanently disabled after
// calling this method.
func (ls *Store) StopInsert() {
	atomic.StoreInt32(&ls.procInterrupt, 1)
}

// insertStopped returns true after StopInsert has been called.
func (ls *Store) insertStopped() bool {
	return atomic.LoadInt32(&ls.procInterrupt) == 1
}

// GetAncestor retrieves the Nth ancestor of a given block. It assumes that either the given block or
// a close ancestor of it is canonical. maxNonCanonical points to a downwards counter limiting the
// number of blocks to be individually checked before we reach the canonical chain.
//
// Note: ancestor == 0 returns the same block, 1 returns its parent and so on.
func (ls *Store) GetAncestor(hash common.Hash, number, ancestor uint64, maxNonCanonical *uint64) (common.Hash, uint64) {
	if ancestor > number {
		return common.Hash{}, 0
	}
	if ancestor == 1 {
		// in this case it is cheaper to just read the header
		if header := ls.GetBlock(hash, number); header != nil {
			return header.ParentHash, number - 1
		}
		return common.Hash{}, 0
	}
	for ancestor != 0 {
		if *maxNonCanonical == 0 {
			return common.Hash{}, 0
		}
		*maxNonCanonical--
		ancestor--
		header := ls.GetBlock(hash, number)
		if header == nil {
			return common.Hash{}, 0
		}
		hash = header.ParentHash
		number--
	}

	return hash, number
}
