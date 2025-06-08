package util

import (
	"testing"
)

type Contain struct {
	ID   int64
	Name string
}

type Contain2 struct {
	ID    int64
	Name2 string
}

func TestContainField(t *testing.T) {
	var cs = []Contain{{
		ID:   1,
		Name: "name1",
	}, {
		ID:   3,
		Name: "name3",
	}}
	var c1 = Contain{
		ID:   1,
		Name: "name1",
	}
	var c3 = []Contain2{{
		ID:    1,
		Name2: "name1",
	}}
	var c2 = c1
	c2.Name = "name2"
	t.Log(Contains(cs, c1))
	t.Log(ContainField(cs, c1, "ID"))
	t.Log(Contains(cs, c2))
	t.Log(ContainField(cs, c2, "Name"))
	t.Log(ContainsDiffField[Contain, Contain2](cs, c3, "Name", "Name2"))
	t.Log(ContainsField(cs, cs, "Name"))
}

func TestFloat64(t *testing.T) {
	t.Log(Float64(3.888, 2))
}

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
	t.Log(Remove(ss, "b", "c"))
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

func TestEqual(t *testing.T) {
	t.Log(Equal(Contain{}, Contain2{}))
	t.Log(EqualJson(Contain{}, Contain2{}))
}
