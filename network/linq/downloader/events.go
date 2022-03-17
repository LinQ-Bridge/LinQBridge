package downloader

import "land-bridge/network/utils"

type DoneEvent struct {
	Latest *utils.Block
}
type StartEvent struct{}
type FailedEvent struct{ Err error }
