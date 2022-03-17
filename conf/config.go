package conf

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"

	"land-bridge/utils"
)

var (
	BSC      uint64
	KLAYTN   uint64
	PLATON   uint64
	ETHEREUM uint64
)

func initChainID(env string) {
	switch env {
	case "testnet":
		ETHEREUM = 5
		BSC = 97
		KLAYTN = 1001
		PLATON = 210309
	default:
		ETHEREUM = 1
		BSC = 56
		KLAYTN = 8217
		PLATON = 100
		//PLATON = 210425
	}
}

func NewConfig(filePath string) *Config {
	buf, err := utils.ReadFile(filePath)
	if err != nil {
		logs.Error("NewServiceConfig: failed, err: %s", err)
		return nil
	}
	config := &Config{}
	err = json.Unmarshal(buf, config)
	if err != nil {
		logs.Error("NewServiceConfig: failed, err: %s", err)
		return nil
	}
	initChainID(config.RunMode)
	return config
}

func (cc *ChainListenConfig) GetNodesURL() []string {
	urls := make([]string, 0)
	for _, node := range cc.Nodes {
		urls = append(urls, node.URL)
	}
	return urls
}
