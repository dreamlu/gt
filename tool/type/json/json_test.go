package json

import "testing"

func TestCJSON_Struct(t *testing.T) {
	type Test struct {
		Name string
	}
	var test []Test
	cj := CJSON("[{\"name\":\"static/file/1244456632704831488.jpg\"}]")
	err := cj.Struct(&test)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(test)
}
