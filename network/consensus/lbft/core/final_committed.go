package core

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
)

func (c *core) handleFinalCommitted() error {
	logs.Trace("Received a final committed proposal")

	c.startNewRound(common.Big0)
	return nil
}
