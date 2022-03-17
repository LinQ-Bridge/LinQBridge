package linq

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/beego/beego/v2/core/logs"
	"github.com/urfave/cli"

	"land-bridge/conf"
	"land-bridge/handle/listener"
	"land-bridge/network"
)

var (
	configFile string
)

var LinQFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "config",
		Usage:       "server config file `<path>`",
		Value:       "./conf/config_devnet.json",
		Destination: &configFile,
	},
	network.NodeKeyFileFlag,
}

func StartServer(ctx *cli.Context) {
	for {
		startServer(ctx)
		sig := waitSignal()
		stopServer()
		if sig != syscall.SIGHUP {
			break
		} else {
			continue
		}
	}
}

func startServer(ctx *cli.Context) {
	config := conf.NewConfig(configFile)
	if config == nil {
		logs.Error("startServer - read config failed!")
		return
	}

	if args := ctx.Args(); len(args) > 0 {
		logs.Error(fmt.Errorf("invalid command: %q", args[0]))
	}

	listener.StartCrossChainListen(config.Chains, config.DBConfig)
	network.StartNetWork(ctx, config)
}

func waitSignal() os.Signal {
	exit := make(chan os.Signal, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sc)
	go func() {
		for sig := range sc {
			logs.Info("cross chain listen received signal:(%s).", sig.String())
			exit <- sig
			close(exit)
			break
		}
	}()
	sig := <-exit
	return sig
}

func stopServer() {
	listener.StopCrossChainListen()
}
