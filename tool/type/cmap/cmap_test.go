package cmap

import "testing"

func TestCMap_Struct(t *testing.T) {
	type Te struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
	var te Te
	var param = CMap{}
	param.Add("name", "tea")
	param.Add("id", "1")
	err := param.Struct(&te)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(te)

	t.Log(CMap{}.Add("id", "1").Set("test", "2"))
}

func TestStructToMap(t *testing.T) {
	type Name struct {
		Name string `json:"name"`
		A    int
		B    int
		D    int
		C    int
	}

	cm := StructToCMap(&Name{
		Name: "test",
		A:    1,
	})
	t.Log(cm)
	t.Log(cm.Encode())
}
