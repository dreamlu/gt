package ghttp

import (
	"testing"
)

func TestGet(t *testing.T) {
	b := NewRequest("GET", "https://github.com/dreamlu").Exec()
	t.Log(string(b.data))

	b = NewRequest("GET", "https://github.com/dreamlu/gt/search").
		AddParam("q", "gt").Exec()
	t.Log(string(b.data))
}
