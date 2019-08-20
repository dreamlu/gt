package id

import (
	"github.com/bwmarrin/snowflake"
)

// id generate
// workId is >= 2, default 1 use by go-tool
func NewID(workId int64) (snowflake.ID, error) {

	node, err := snowflake.NewNode(workId)
	if err != nil {
		return 0, err
	}
	return node.Generate(), nil
}
