package core

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/network/consensus"
	"land-bridge/network/consensus/lbft"
)

// Start implements core.Engine.Start
func (c *core) Start() error {
	// Start a new round from last sequence + 1
	c.startNewRound(common.Big0)

	// Tests will handle events itself, so we have to make subscribeEvents()
	// be able to call in test.
	c.subscribeEvents()
	c.handlerWg.Add(1)
	go c.handleEvents()

	return nil
}

// Stop implements core.Engine.Stop
func (c *core) Stop() error {
	c.stopTimer()
	c.unsubscribeEvents()

	c.handlerWg.Wait()
	return nil
}

// Subscribe both internal and external events
func (c *core) subscribeEvents() {
	c.events = c.backend.EventMux().Subscribe(
		// external events
		consensus.RequestEvent{},
		consensus.MessageEvent{},
		// internal events
		backlogEvent{},
	)
	c.timeoutSub = c.backend.EventMux().Subscribe(
		timeoutEvent{},
	)
	c.finalCommittedSub = c.backend.EventMux().Subscribe(
		consensus.FinalCommittedEvent{},
	)
}

// Unsubscribe all events
func (c *core) unsubscribeEvents() {
	c.events.Unsubscribe()
	c.timeoutSub.Unsubscribe()
	c.finalCommittedSub.Unsubscribe()
}

func (c *core) handleEvents() {
	// Clear state
	defer func() {
		c.current = nil
		c.handlerWg.Done()
	}()

	for {
		select {
		case event, ok := <-c.events.Chan():
			if !ok {
				return
			}

			// A real event arrived, process interesting content
			switch ev := event.Data.(type) {
			case consensus.RequestEvent:
				c.newRoundChangeTimer(true)
				r := &consensus.Request{
					Proposal: ev.Proposal,
				}
				err := c.handleRequest(r)
				if err == errFutureMessage {
					c.storeRequestMsg(r)
				}
			case consensus.MessageEvent:
				if err := c.handleMsg(ev.Payload); err == nil {
					c.backend.Gossip(c.valSet, ev.Payload)
				}
			case backlogEvent:
				// No need to check signature for internal messages
				if err := c.handleCheckedMsg(ev.msg, ev.src); err == nil {
					p, err := ev.msg.Payload()
					if err != nil {
						logs.Warn("Get message payload failed", "err", err)
						continue
					}
					c.backend.Gossip(c.valSet, p)
				}
			}
		case _, ok := <-c.timeoutSub.Chan():
			if !ok {
				return
			}
			logs.Error("lbft consensus timeout")
			c.handleTimeoutMsg()
		case event, ok := <-c.finalCommittedSub.Chan():
			if !ok {
				return
			}
			switch event.Data.(type) {
			case consensus.FinalCommittedEvent:
				c.handleFinalCommitted()
			}
		}
	}
}

// sendEvent sends events to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.EventMux().Post(ev)
}

func (c *core) handleMsg(payload []byte) error {
	// Decode message and check its signature
	msg := new(Message)
	if err := msg.FromPayload(payload, c.validateFn); err != nil {
		logs.Error("Failed to decode message from payload", "err", err)
		return err
	}

	// Only accept message if the address is valid
	_, src := c.valSet.GetByAddress(msg.Address)
	if src == nil {
		logs.Error("Invalid address in message", "msg", msg)
		return lbft.ErrUnauthorizedAddress
	}

	return c.handleCheckedMsg(msg, src)
}

func (c *core) handleCheckedMsg(msg *Message, src lbft.Validator) error {
	// Store the message if it's a future message
	testBacklog := func(err error) error {
		if err == errFutureMessage {
			c.storeBacklog(msg, src)
		}
		return err
	}

	switch msg.Code {
	case MsgPreprepare:
		err := c.handlePreprepare(msg, src)
		return testBacklog(err)
	case MsgPrepare:
		err := c.handlePrepare(msg, src)
		return testBacklog(err)
	case MsgCommit:
		err := c.handleCommit(msg, src)
		return testBacklog(err)
	case MsgRoundChange:
		err := c.handleRoundChange(msg, src)
		return testBacklog(err)
	default:
		logs.Error("Invalid message", "msg", msg)
	}

	return errInvalidMessage
}

func (c *core) handleTimeoutMsg() {
	// If we're not waiting for round change yet, we can try to catch up
	// the max round with F+1 round change message. We only need to catch up
	// if the max round is larger than current round.
	if !c.waitingForRoundChange {
		maxRound := c.roundChangeSet.MaxRound(c.valSet.F() + 1)
		if maxRound != nil && maxRound.Cmp(c.current.Round()) > 0 {
			c.sendRoundChange(maxRound)
			return
		}
	}

	lastProposal, _ := c.backend.LastProposal()
	if lastProposal != nil && lastProposal.Number().Cmp(c.current.Sequence()) >= 0 {
		logs.Trace("round change timeout, catch up latest sequence", "number", lastProposal.Number().Uint64())
		c.startNewRound(common.Big0)
	} else {
		c.sendNextRoundChange()
	}
}
