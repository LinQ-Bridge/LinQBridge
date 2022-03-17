package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"land-bridge/conf"
	"land-bridge/models"
	"land-bridge/network"
)

type BridgeDao struct {
	dbCfg *conf.DBConfig
	db    *gorm.DB
}

func NewBridgeDao(dbCfg *conf.DBConfig) *BridgeDao {
	dao := &BridgeDao{
		dbCfg: dbCfg,
	}
	Logger := logger.Default
	if dbCfg.Debug == true {
		Logger = Logger.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(dbCfg.User+":"+dbCfg.Password+"@tcp("+dbCfg.URL+")/"+
		dbCfg.Scheme+"?charset=utf8"), &gorm.Config{Logger: network.Nologger{}})
	if err != nil {
		panic(err)
	}
	dao.db = db
	return dao
}

func (dao *BridgeDao) UpdateEvents(wrapperTransactions []*models.WrapperTransaction, srcTransactions []*models.SrcTransaction, dstTransactions []*models.DstTransaction) error {
	tx := dao.db.Begin()
	if wrapperTransactions != nil && len(wrapperTransactions) > 0 {
		res := tx.Save(wrapperTransactions)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}
	if srcTransactions != nil && len(srcTransactions) > 0 {
		res := tx.Save(srcTransactions)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	if dstTransactions != nil && len(dstTransactions) > 0 {
		res := tx.Save(dstTransactions)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}
	tx.Commit()
	return nil
}

func (dao *BridgeDao) GetChain(chainID uint64) (*models.Chain, error) {
	chain := new(models.Chain)
	res := dao.db.Where("chain_id = ?", chainID).First(chain)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("no record")
	}
	chain.HeightSwap = 0
	return chain, nil
}

func (dao *BridgeDao) UpdateChain(chain *models.Chain) error {
	if chain == nil {
		return fmt.Errorf("no value!\n")
	}

	res := dao.db.Updates(chain)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no update!\n")
	}
	return nil
}

func (dao *BridgeDao) AddChains(chain []*models.Chain, chainFees []*models.ChainFee) error {
	if chain == nil || len(chain) == 0 {
		return nil
	}
	res := dao.db.Create(chain)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("add chain failed!\n")
	}
	if chainFees == nil || len(chainFees) == 0 {
		return nil
	}
	res = dao.db.Create(chainFees)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("add chain fee failed!\n")
	}
	return nil
}
