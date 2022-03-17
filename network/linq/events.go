package linq

import (
	"land-bridge/network/utils"
)

type ChainHeadEvent struct{ Block *utils.Block }

// NewMinedBlockEvent is posted when a block has been imported.
type NewMinedBlockEvent struct{ Block *utils.Block }
