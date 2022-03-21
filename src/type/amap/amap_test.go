package amap

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

	bm := ToAMap(test)
	t.Log(bm)

	var m = make(map[string]string)
	m["A"] = "A"
	bm = ToAMap(m)
	t.Log(bm)
}
