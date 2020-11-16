package http

import (
	"testing"
)

func TestGet(t *testing.T) {
	b := NewRequest("GET", "https://github.com/dreamlu").Exec()
	t.Log(string(b.data))
}
