package ghttp

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestGet(t *testing.T) {
	//b := NewRequest(GET, "https://github.com/dreamlu").Exec()
	//t.Log(string(b.data))

	//b = NewRequest("GET", "https://github.com/dreamlu/gt/search").
	//	AddParam("q", "gt").Exec()
	//t.Log(string(b.data))

	b := NewRequest("GET", "https://dev.thuwater.com/api/data/dataSumRange?et=2022-08-29 23:59:59&placeId=16100&st=2022-08-28 00:00:00&type=1").
		AddHeader("X_TH_TOKEN", "eyJhbGciOiJIUzM4NCJ9.eyJwbGFjZXMiOiIxNjA5OSwxNjEwMCwxNjEwMSwxNjEwMiwxNjEwMywxNjEwNCwxNjEwOCwxNjExMywxNjExNCwxNjExNSwxNjExNiwxNjExNywxNjExOCwxNjExOSwxNjE4OCwxNjE4OSwxNjE5MCwxNjE5MSwxNjE5MiwxNjE5MywxNjE5NCwxNjE5NSwxNjE5NiwxNjE5NywxNjE5OCwxNjE5OSwxNjIwMCwxNjIwMSwxNjIwMiwxNjIwMywxNjIwNCwxNjIwNSwxNjIwNiwxNjIwNywxNjIwOSwxNjIxMCwxNjIzMSwxNjIzMiwxNjIzMywxNjIzNCwxNjIzNSIsInByb2plY3RzIjoiMjAyMjAwMDEiLCJ1c2VybmFtZSI6ImRhdGEyMDIyMDAwMSIsImlhdCI6MTY2MTc1Njc1MiwiZXhwIjoxNjYxNzYyMDkyfQ.pREjAlhSKGLubEvAhYiIPEprvTiRxoja0I-iNisNcpolTsjj-hKyHagL4GUWE6IO").
		AddHeader("", "").
		Exec()
	t.Log(string(b.data))
}

func TestPostJSON(t *testing.T) {
	r := NewRequest(POST, "https://github.com/dreamlu")
	type Search struct {
		Q string `json:"q"`
	}
	b, _ := json.Marshal(Search{Q: "gt"})
	r.SetBody(bytes.NewReader(b))
	res := r.Exec()
	t.Log(res.String())
}

func TestUpload(t *testing.T) {
	r := NewRequest(POST, "http://192.168.10.11/upload")
	//r.SetContentType(ContentTypeForm)
	_ = r.AddFile("file", "test.txt", "test1.txt")
	res := r.Exec()
	t.Log(res.String())
}
