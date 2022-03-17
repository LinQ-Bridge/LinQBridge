package core

import (
	"reflect"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/network/consensus"
	"land-bridge/network/consensus/lbft"
)

func (c *core) sendCommit() {
	sub := c.current.Subject()
	c.broadcastCommit(sub)
}

func (c *core) sendCommitForOldBlock(view *consensus.View, digest common.Hash) {
	sub := &consensus.Subject{
		View:   view,
		Digest: digest,
	}
	c.broadcastCommit(sub)
}

func (c *core) broadcastCommit(sub *consensus.Subject) {
	encodedSubject, err := Encode(sub)
	if err != nil {
		logs.Error("Failed to encode", "subject", sub)
		return
	}

	c.broadcast(&Message{
		Code: MsgCommit,
		Msg:  encodedSubject,
	})
}

func (c *core) handleCommit(msg *Message, src lbft.Validator) error {
	// Decode COMMIT message
	var commit *consensus.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}

	if err := c.checkMessage(MsgCommit, commit.View); err != nil {
		return err
	}

	if err := c.verifyCommit(commit, src); err != nil {
		return err
	}

	c.acceptCommit(msg, src)

	// Commit the proposal once we have enough COMMIT messages and we are not in the Committed state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.current.Commits.Size() >= c.QuorumSize() && c.state.Cmp(StateCommitted) < 0 {
		bridge := c.backend.Bridge()
		if c.IsProposer() {
			signatures := make(map[common.Address][]byte)
			for addr, message := range c.current.Commits.messages {
				signatures[addr] = message.HashSign
			}
			go bridge.BridgeToChainB(c.current.Proposal().TxHash(), signatures)
		} else {
			go bridge.UpdateWrapper(c.current.Proposal().TxHash())
		}

		// Still need to call LockHash here since state can skip Prepared state and jump directly to the Committed state.
		c.current.LockHash()
		c.commit()
	}

	return nil
}

// verifyCommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyCommit(commit *consensus.Subject, src lbft.Validator) error {
	sub := c.current.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logs.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptCommit(msg *Message, src lbft.Validator) error {
	// Add the COMMIT message to current round state
	if err := c.current.Commits.Add(msg); err != nil {
		logs.Error("Failed to record commit message", "msg", msg, "err", err)
		return err
	}

	return nil
}
