package ghttp

import (
	"encoding/json"
	"net/http"
)

//Response Http请求返回内容
type Response struct {
	*http.Response
	data  []byte
	error error
}

func (m *Response) Error() error {
	return m.error
}

func (m *Response) Bytes() ([]byte, error) {
	return m.data, m.error
}

func (m *Response) MustBytes() []byte {
	return m.data
}

func (m *Response) String() (string, error) {
	return string(m.data), m.error
}

func (m *Response) MustString() string {
	return string(m.data)
}

func (m *Response) Unmarshal(v any) error {
	if m.error != nil {
		return m.error
	}
	return json.Unmarshal(m.data, v)
}
