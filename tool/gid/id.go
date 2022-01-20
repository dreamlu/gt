package gid

import (
	"github.com/bwmarrin/snowflake"
)

// NewID id generate
// workId is >= 2, default 1 use by gt
func NewID(workId int64) (snowflake.ID, error) {

	node, err := snowflake.NewNode(workId)
	if err != nil {
		return 0, err
	}
	return node.Generate(), nil
}
