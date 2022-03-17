package tools

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var ConfigCMD = cli.Command{
	Name:   "config",
	Usage:  "Generate blank config-json file",
	Action: config,
}

func config(ctx *cli.Context) {
	fmt.Print("LinQ config.json Building.")
	fileName := "config.json"
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		return
	}

	fmt.Print(".")

	_, err = f.Write([]byte(configTemplate))
	if err != nil {
		return
	}

	fmt.Println(".")

	fmt.Println("LinQ config.json Generated successfully.")
}

const (
	configTemplate = `{
  "RunMode": "testnet",
  "DBConfig": {
    "Debug": false,
    "URL": "",
    "Scheme": "",
    "User": "",
    "Password": ""
  },
  "LinQConfig": {
    "DefaultBootNodes": [],
    "Addr": "0.0.0.0",
    "Port": 30303
  },
  "Chains": [
    {
      "ChainName": "Klaytn",
      "ChainID": 1001,
      "ListenSlot": 5,
      "BatchSize": 5,
      "defer": 5,
      "Nodes": [
        {
          "Url": ""
        }
      ],
      "CCMContract": "",
      "NFTProxyContract": "",
      "NFTWrapperContract": "",
      "NFTQueryContract": ""
    }
  ]
}`
)
