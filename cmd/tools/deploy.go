package tools

import (
	"fmt"

	"github.com/urfave/cli"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"land-bridge/conf"
	"land-bridge/handle/dao"
	"land-bridge/models"
	"land-bridge/network"
)

var (
	configPath, genesisPath string
)

var DeployCMD = cli.Command{
	Name:  "init",
	Usage: "Server will init or update db when true",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Usage:       "Server config file `<path>`",
			Value:       "./conf/config_devnet.json",
			Destination: &configPath,
		},
		cli.StringFlag{
			Name:        "genesis",
			Usage:       "Server genesis.json file `<path>`",
			Value:       "./conf/genesis.json",
			Destination: &genesisPath,
		},
	},
	Action: deploy,
}

func deploy(ctx *cli.Context) {
	fmt.Println("LinQ Tool Start Initialization Procedure.")
	cfg := conf.NewConfig(configPath)
	dbCfg := cfg.DBConfig
	db, err := gorm.Open(mysql.Open(dbCfg.User+":"+dbCfg.Password+"@tcp("+dbCfg.URL+")/"+
		dbCfg.Scheme+"?charset=utf8"), &gorm.Config{Logger: network.Nologger{}})
	if err != nil {
		panic(err)
	}
	fmt.Println("Linked Mysql Database Successfully.")
	err = db.AutoMigrate(
		&models.ChainFee{},
		&models.Chain{},
		&models.DstTransaction{},
		&models.DstSwap{},
		&models.DstTransfer{},
		&models.NFTProfile{},
		&models.PriceMarket{},
		&models.SrcTransaction{},
		&models.SrcSwap{},
		&models.SrcTransfer{},
		&models.TimeStatistic{},
		&models.TokenBasic{},
		&models.TokenMap{},
		&models.Token{},
		&models.WrapperTransaction{},
		&models.ErrorTransaction{},
		&models.Block{},
		&models.Snapshot{},
		&models.TxHashHistory{},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Table Information Generated Successfully.")

	dao := dao.NewBridgeDao(cfg.DBConfig)
	if dao == nil {
		panic("server is invalid")
	}
	var chains []*models.Chain
	for _, c := range cfg.Chains {
		chain := &models.Chain{ChainID: c.ChainID, Name: c.ChainName, BackwardBlockNumber: c.BatchSize}
		chains = append(chains, chain)
	}
	dao.AddChains(chains, nil)

	fmt.Println("Chain Information Insert Successfully.")

	initGenesis(genesisPath, db)

	fmt.Println("Genesis Block Information Insert Successfully.")

	fmt.Println("LinQ Initialization Successfully.")
}
