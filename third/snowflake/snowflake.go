package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/dreamlu/gt/third/log"
)

// NewID id generate
// workId is >= 2, default 1 use by gt
func NewID(workId int64) snowflake.ID {

	node, err := snowflake.NewNode(workId)
	if err != nil {
		log.Error(err)
		return 0
	}
	return node.Generate()
}
