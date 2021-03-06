package core

import (
	"reflect"

	"github.com/beego/beego/v2/core/logs"

	"land-bridge/network/consensus"
	"land-bridge/network/consensus/lbft"
)

func (c *core) sendPrepare() {
	sub := c.current.Subject()
	encodedSubject, err := Encode(sub)
	if err != nil {
		logs.Error("Failed to encode", "subject", sub)
		return
	}

	c.broadcast(&Message{
		Code: MsgPrepare,
		Msg:  encodedSubject,
	})
}

func (c *core) handlePrepare(msg *Message, src lbft.Validator) error {
	// Decode PREPARE message
	var prepare *consensus.Subject
	err := msg.Decode(&prepare)
	if err != nil {
		return errFailedDecodePrepare
	}

	if err := c.checkMessage(MsgPrepare, prepare.View); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrepare and checkMessage implies it is processing on the locked block since it was verified in the Preprepared state.
	if err := c.verifyPrepare(prepare, src); err != nil {
		return err
	}

	c.acceptPrepare(msg, src)

	// Change to Prepared state if we've received enough PREPARE messages or it is locked
	// and we are in earlier state before Prepared state.
	if ((c.current.IsHashLocked() && prepare.Digest == c.current.GetLockedHash()) || c.current.GetPrepareOrCommitSize() >= c.QuorumSize()) &&
		c.state.Cmp(StatePrepared) < 0 {
		c.current.LockHash()
		c.setState(StatePrepared)
		c.sendCommit()
	}

	return nil
}

// verifyPrepare verifies if the received PREPARE message is equivalent to our subject
func (c *core) verifyPrepare(prepare *consensus.Subject, src lbft.Validator) error {
	sub := c.current.Subject()
	if !reflect.DeepEqual(prepare, sub) {
		logs.Warn("Inconsistent subjects between PREPARE and proposal", "expected", sub, "got", prepare)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptPrepare(msg *Message, src lbft.Validator) error {
	// Add the PREPARE message to current round state
	if err := c.current.Prepares.Add(msg); err != nil {
		logs.Error("Failed to add PREPARE message to round state", "msg", msg, "err", err)
		return err
	}

	return nil
}
