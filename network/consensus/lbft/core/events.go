package core

import (
	"land-bridge/network/consensus/lbft"
)

type backlogEvent struct {
	src lbft.Validator
	msg *Message
}

type timeoutEvent struct{}
