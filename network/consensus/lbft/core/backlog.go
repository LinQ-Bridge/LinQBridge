package core

import (
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"

	"land-bridge/network/consensus"
	"land-bridge/network/consensus/lbft"
)

var (
	// msgPriority is defined for calculating processing priority to speedup consensus
	// msgPreprepare > msgCommit > msgPrepare
	msgPriority = map[uint64]int{
		MsgPreprepare: 1,
		MsgCommit:     2,
		MsgPrepare:    3,
	}
)

// checkMessage checks the message state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the message view is larger than current view
// return errOldMessage if the message view is smaller than current view
func (c *core) checkMessage(msgCode uint64, view *consensus.View) error {
	if view == nil || view.Sequence == nil || view.Round == nil {
		return errInvalidMessage
	}

	if msgCode == MsgRoundChange {
		if view.Sequence.Cmp(c.currentView().Sequence) > 0 {
			logs.Error("errFutureMessage view", view.Sequence.Uint64(), " CS:", c.currentView().Sequence.Uint64())
			return errFutureMessage
		} else if view.Cmp(c.currentView()) < 0 {
			return errOldMessage
		}
		return nil
	}

	if view.Cmp(c.currentView()) > 0 {
		return errFutureMessage
	}

	if view.Cmp(c.currentView()) < 0 {
		return errOldMessage
	}

	if c.waitingForRoundChange {
		return errFutureMessage
	}

	// StateAcceptRequest only accepts msgPreprepare
	// other messages are future messages
	if c.state == StateAcceptRequest {
		if msgCode > MsgPreprepare {
			return errFutureMessage
		}
		return nil
	}
	// For states(StatePreprepared, StatePrepared, StateCommitted),
	// can accept all message types if processing with same view
	return nil
}

func (c *core) processBacklog() {
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	for srcAddress, backlog := range c.backlogs {
		if backlog == nil {
			continue
		}
		_, src := c.valSet.GetByAddress(srcAddress)
		if src == nil {
			// validator is not available
			delete(c.backlogs, srcAddress)
			continue
		}

		isFuture := false

		// We stop processing if
		//   1. backlog is empty
		//   2. The first message in queue is a future message
		for !(backlog.Empty() || isFuture) {
			m, prio := backlog.Pop()
			msg := m.(*Message)
			var view *consensus.View
			switch msg.Code {
			case MsgPreprepare:
				var m *consensus.Preprepare
				err := msg.Decode(&m)
				if err == nil {
					view = m.View
				}
				// for msgRoundChange, msgPrepare and msgCommit cases
			default:
				var sub *consensus.Subject
				err := msg.Decode(&sub)
				if err == nil {
					view = sub.View
				}
			}
			if view == nil {
				logs.Debug("Nil view", "msg", msg)
				continue
			}
			// Push back if it's a future message
			err := c.checkMessage(msg.Code, view)
			if err != nil {
				if err == errFutureMessage {
					logs.Trace("Stop processing backlog", "msg", msg)
					backlog.Push(msg, prio)
					isFuture = true
					break
				}
				logs.Trace("Skip the backlog event", "msg", msg, "err", err)
				continue
			}
			logs.Trace("Post backlog event", "msg", msg)

			go c.sendEvent(backlogEvent{
				src: src,
				msg: msg,
			})
		}
	}
}

func (c *core) storeBacklog(msg *Message, src lbft.Validator) {
	if src.Address() == c.Address() {
		logs.Warn("Backlog from self")
		return
	}

	logs.Trace("Store future message")

	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	logs.Debug("Retrieving backlog queue", "for", src.Address(), "backlogs_size", len(c.backlogs))
	backlog := c.backlogs[src.Address()]
	if backlog == nil {
		backlog = prque.New()
	}
	switch msg.Code {
	case MsgPreprepare:
		var p *consensus.Preprepare
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.View))
		}
		// for msgRoundChange, msgPrepare and msgCommit cases
	default:
		var p *consensus.Subject
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.View))
		}
	}
	c.backlogs[src.Address()] = backlog
}

func toPriority(msgCode uint64, view *consensus.View) float32 {
	if msgCode == MsgRoundChange {
		// For msgRoundChange, set the message priority based on its sequence
		return -float32(view.Sequence.Uint64() * 1000)
	}
	// FIXME: round will be reset as 0 while new sequence
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Sequence limits the range of round is from 0 to 99
	return -float32(view.Sequence.Uint64()*1000 + view.Round.Uint64()*10 + uint64(msgPriority[msgCode]))
}
