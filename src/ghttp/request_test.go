package ghttp

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestGet(t *testing.T) {
	b := NewRequest(GET, "https://github.com/dreamlu").Exec()
	t.Log(string(b.data))

	b = NewRequest("GET", "https://github.com/dreamlu/gt/search").
		AddParam("q", "gt").Exec()
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
