package tmap

import (
	"strings"
	"testing"
)

type Test struct {
	A string `json:"a"`
	B int    `json:"b"`
	C []int  `json:"c"`
}

func TestStructToTMap(t *testing.T) {

	var test = Test{
		A: "A",
		B: 1,
		C: []int{3, 4},
	}

	bm := ToTMap[string, any](test)
	t.Log(bm)

	bi := NewTMap[int, Test]()
	bi.Set(1, test)
	bi.Set(2, test)
	t.Log(bi)

	bis := Set(1, test).Set(2, test)
	t.Log(bis)

	var m = make(map[string]string)
	m["A"] = "A"
	bms := ToTMap[string, string](m)
	t.Log(bms)
	t.Log(strings.Index(bms.Get("A"), "A"))

	var mt = make(map[string]Test)
	mt["A"] = test
	bmt := TMap[string, Test]{}
	bmt.Set("A", test)
	t.Log(bmt.Get("A").C)

	bmt = ToTMap[string, Test](mt)
	t.Log(bmt)
	t.Log(test.A)
	t.Log(bmt.Get("A").A)
}
