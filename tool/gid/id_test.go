package gid

import (
	"log"
	"testing"
)

func TestId(t *testing.T) {
	id, err := NewID(1)
	if err != nil {
		log.Print(err)
		return
	}
	t.Log(id.String())
}
