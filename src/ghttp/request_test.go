package ghttp

import (
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
	r.SetJsonBody(Search{Q: "gt"})
	res := r.Exec()
	t.Log(res.String())
}

func TestPostForm(t *testing.T) {
	r := NewRequest(POST, "https://www.beqege.cc/search.php?keyword=t")
	r.SetContentType(ContentTypeFormUrl)
	r.SetHeader("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	r.SetForm("keyword", "圣墟")
	r.SetHeader("cookie", "cf_clearance=8f7e8826a9fca60e000b254aa1c62eb66ea95e0d-1667200369-0-150; __cf_bm=5Cm7Fkpqq9.1MnLU2BazDmowXaqnhikXmo01Sdt3zvM-1667205564-0-AdSmEFZDmispSpenxu3ozQhp2WaM7fT6qNTGeiPR4/R/h14+fc/6d0sHdrnuedOcTs5ADq8P3vlOefRTbw6Nq2WyCsamhMZr0qaWIzObr0ieRp1qc1V0yFA03zvwfBb9aQ==")
	res := r.Exec()
	t.Log(res.String())
}

func TestUpload(t *testing.T) {
	r := NewRequest(POST, "http://192.168.10.11/upload")
	//r.SetContentType(ContentTypeForm)
	_ = r.SetFile("file", "test.txt", "test1.txt")
	res := r.Exec()
	t.Log(res.String())
}

func TestStatus(t *testing.T) {
	r := NewRequest(POST, "403 url")
	res := r.Exec()
	if err := res.Error(); err != nil {
		t.Log(err)
		t.Log(res.String())
		return
	}
	t.Log(res.String())
}
