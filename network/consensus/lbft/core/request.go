package core

import (
	"github.com/beego/beego/v2/core/logs"

	"land-bridge/network/consensus"
)

func (c *core) handleRequest(request *consensus.Request) error {
	if err := c.checkRequestMsg(request); err != nil {
		if err == errInvalidMessage {
			logs.Warn("invalid request")
			return err
		}
		logs.Warn("unexpected request", "err", err, "number", request.Proposal.Number(), "hash", request.Proposal.Hash())
		return err
	}
	logs.Trace("handleRequest", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.current.pendingRequest = request

	if c.state == StateAcceptRequest {
		c.sendPreprepare(request)
	}
	return nil
}

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
func (c *core) checkRequestMsg(request *consensus.Request) error {
	if request == nil || request.Proposal == nil {
		return errInvalidMessage
	}

	if c := c.current.sequence.Cmp(request.Proposal.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(request *consensus.Request) {
	logs.Trace("Store future request", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(request, float32(-request.Proposal.Number().Int64()))
}
