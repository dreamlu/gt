package bmap

import (
	"testing"
)

type Test struct {
	A string `json:"a"`
	B int    `json:"b"`
	C []int  `json:"c"`
}

func TestStructToBMap(t *testing.T) {

	var test = Test{
		A: "A",
		B: 1,
		C: []int{3, 4},
	}

	bm := ToBMap(test)
	t.Log(bm)

	var m = make(map[string]string)
	m["A"] = "A"
	bm = ToBMap(m)
	t.Log(bm)
}
