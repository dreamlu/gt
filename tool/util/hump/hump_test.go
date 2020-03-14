package hump

import "testing"

func TestHumpToLine(t *testing.T) {
	t.Log(HumpToLine("ABTest"))
	t.Log(LineToHump("a_b_test"))
}