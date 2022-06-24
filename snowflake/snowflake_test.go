package snowflake

import (
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID(1)
	t.Log(id.String())
}
