package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"

	"land-bridge/cmd/linq"
	"land-bridge/cmd/tools"
	"land-bridge/conf"
)

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Name = "LinQ"
	app.Usage = "MultiChain-NFT-Bridge Service"
	app.Action = linq.StartServer
	app.Version = "1.0.0"
	app.Copyright = "Copyright 2022 The LinQ Authors"
	app.Flags = linq.LinQFlags
	app.Commands = []cli.Command{
		tools.CMD,
	}
	app.Before = func(context *cli.Context) error {
		fmt.Print(conf.LOGO)
		fmt.Println()
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	app.After = func(context *cli.Context) error {
		fmt.Println()
		fmt.Println("Thank You For Use LinQ.")
		fmt.Println()
		return nil
	}

	return app
}

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
