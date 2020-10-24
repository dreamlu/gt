package http

import (
	"testing"
)

func TestGet(t *testing.T) {
	b, _ := Get("https://github.com/dreamlu")
	t.Log(string(b))
}
