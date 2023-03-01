package util

import (
	"testing"
)

func TestAesEn(t *testing.T) {

	var as = NewAes()
	t.Log("[aesEn]:", as.AesEn("admin"))
	t.Log("[aesDe]:", as.AesDe("sPa0sTmDf6gasS9tHvIqKw=="))
	t.Log(as.IsAes("13242trergf"))
	t.Log(as.IsAes("sPa0sTmDf6gasS9tHvIqKw=="))
}

func TestRemove(t *testing.T) {
	ss := []string{"a", "b", "c", "a", "b"}
	t.Log(RemoveDuplicate(ss))
	t.Log(Remove(ss, "b"))
	type S struct {
		A string
		B string
	}
	s := []*S{
		{
			A: "a",
			B: "b",
		},
		{
			A: "a",
			B: "b",
		},
		{
			A: "a",
			B: "c",
		},
	}
	res := RemoveDuplicate(s)
	t.Log(res)
	res = Remove(s, &S{
		A: "a",
		B: "b",
	})
	t.Log(res)
}

func TestHumpToLine(t *testing.T) {
	t.Log(HumpToLine("ABTest"))
	t.Log(LineToHump("a_b_test"))
	t.Log(HumpToLine("ID"))
	t.Log(HumpToLine("AbC"))
}
