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
		Name string `json:"name" gorm:"column:tb_name" gt:"field:name2"`
		A    int
		B    *int
		C    *int
		D    int
	}

	var i = 3
	cm := StructToCMap(&Name{
		Name: "test",
		A:    1,
		B:    nil,
		C:    &i,
	})
	t.Log(cm)
	t.Log(cm.Encode())
}
