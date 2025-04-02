package ghttp

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	requestMsg  []byte
	responseMsg []byte
	*http.Response
	data  []byte
	error error
}

func (m *Response) StatusCode() int {
	return m.Response.StatusCode
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

func (m *Response) Json() (json.RawMessage, error) {
	return m.data, m.error
}

func (m *Response) MustJson() json.RawMessage {
	return m.data
}

func (m *Response) Unmarshal(v any) error {
	if m.error != nil {
		return m.error
	}
	return json.Unmarshal(m.data, v)
}

func (m *Response) RequestMsg() string {
	return string(m.requestMsg)
}

func (m *Response) ResponseMsg() string {
	return string(m.responseMsg)
}
