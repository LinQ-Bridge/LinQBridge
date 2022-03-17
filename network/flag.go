package network

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"

	"land-bridge/conf"
	"land-bridge/network/node"
	"land-bridge/network/p2p"
	"land-bridge/network/p2p/enode"
	"land-bridge/utils"
)

var (
	NodeKeyFileFlag = cli.StringFlag{
		Name:  "nodekey",
		Usage: "p2p node key file",
	}
)

// SetNodeConfig applies node-related command line flags to the config.
func SetNodeConfig(ctx *cli.Context, cfg *node.Config, conf *conf.Config) {
	SetP2PConfig(ctx, &cfg.P2P, conf)
}

func SetP2PConfig(ctx *cli.Context, cfg *p2p.Config, conf *conf.Config) {
	setNodeKey(ctx, cfg)
	setListenAddress(cfg, conf)
	setBootstrapNodes(cfg, conf)
}

func setNodeKey(ctx *cli.Context, cfg *p2p.Config) {
	var (
		file = ctx.GlobalString(NodeKeyFileFlag.Name)
		key  *ecdsa.PrivateKey
		err  error
	)

	if file != "" {
		if key, err = crypto.LoadECDSA(file); err != nil {
			utils.Fatalf("Option %q: %v", NodeKeyFileFlag.Name, err)
		}
		cfg.PrivateKey = key
	}
}

func setListenAddress(cfg *p2p.Config, conf *conf.Config) {
	if conf.LinQConfig.Port > 0 {
		cfg.ListenAddr = fmt.Sprintf("%s:%d", conf.LinQConfig.Addr, conf.LinQConfig.Port)
		logs.Info("P2P addr set :", cfg.ListenAddr)
	}
}

func setBootstrapNodes(cfg *p2p.Config, conf *conf.Config) {
	urls := conf.LinQConfig.DefaultBootNodes
	switch {
	case cfg.BootstrapNodes != nil:
		return // already set, don't apply defaults.
	}

	cfg.BootstrapNodes = make([]*enode.Node, 0, len(urls))
	for _, url := range urls {
		if url != "" {
			node, err := enode.Parse(enode.ValidSchemes, url)
			if err != nil {
				logs.Warn("Bootstrap URL invalid", "enode", url, "err", err)
				continue
			}
			cfg.BootstrapNodes = append(cfg.BootstrapNodes, node)
		}
	}
}
