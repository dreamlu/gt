package json

import "testing"

func TestCJSON_Struct(t *testing.T) {
	type Test struct {
		Name string
	}
	var test []Test
	cj := CJSON("[{\"name\":\"static/file/1244456632704831488.jpg\"}]")
	err := cj.Unmarshal(&test)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(test)
}

func TestCJSON_Array(t *testing.T) {
	var test []int
	cj := CJSON("[1,2,3]")
	err := cj.Unmarshal(&test)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(test)
}
