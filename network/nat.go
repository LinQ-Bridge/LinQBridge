package network

import (
	"crypto/ecdsa"
	"os"
	"os/signal"
	"syscall"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"land-bridge/conf"
	"land-bridge/network/bridge"
	"land-bridge/network/linq"
	"land-bridge/network/linq/txblock"
	"land-bridge/network/node"
	"land-bridge/utils"
)

const (
	clientIdentifier = "land-bridge" // Client identifier to advertise over the network
)

func StartNetWork(ctx *cli.Context, conf *conf.Config) error {
	stack, _ := MakeFullNode(ctx, conf)

	defer stack.Close()
	startNode(stack)

	stack.Wait()
	return nil
}

func dbinit(dbCfg *conf.DBConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dbCfg.User+":"+dbCfg.Password+"@tcp("+dbCfg.URL+")/"+
		dbCfg.Scheme+"?charset=utf8"), &gorm.Config{Logger: Nologger{}})

	if err != nil {
		panic(err)
	}

	return db
}

func startNode(stack *node.Node) {
	if err := stack.Start(); err != nil {
		utils.Fatalf("Error starting protocol stack: %v", err)
	}

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)

		<-sigc
		logs.Info("Got interrupt, shutting down...")
		go stack.Close()
		for i := 5; i > 0; i-- {
			<-sigc
			if i > 1 {
				logs.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
	}()
}

func MakeFullNode(ctx *cli.Context, conf *conf.Config) (*node.Node, *linq.LinQ) {
	db := dbinit(conf.DBConfig)
	stack, cfg := makeConfigNode(ctx, conf)
	privkey := cfg.Node.NodeKey()
	privStr := crypto.PubkeyToAddress(*privkey.Public().(*ecdsa.PublicKey)).Hex()

	stack.SetDB(db)
	stack.Bridge = bridge.NewBridge(db, conf, privkey)
	stack.Pool = txblock.NewBlockPool()

	Linq := RegisterPeerService(stack, privStr)

	return stack, Linq
}

func RegisterPeerService(stack *node.Node, privStr string) *linq.LinQ {
	land, err := linq.New(stack, privStr)
	if err != nil {
		utils.Fatalf("Failed to register the Ethereum service: %v", err)
	}
	return land
}

func makeConfigNode(ctx *cli.Context, conf *conf.Config) (*node.Node, Config) {
	cfg := Config{
		Node: defaultNodeConfig(),
	}

	SetNodeConfig(ctx, &cfg.Node, conf)
	stack, err := node.NewNode(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}

	return stack, cfg
}

type Config struct {
	Node node.Config
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier

	return cfg
}
