package ghttp

import (
	"bytes"
	"github.com/dreamlu/gt/src/file/fs"
	"github.com/dreamlu/gt/src/type/cmap"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const (
	GET             = http.MethodGet
	POST            = http.MethodPost
	DELETE          = http.MethodDelete
	PUT             = http.MethodPut
	PATCH           = http.MethodPatch
	HEAD            = http.MethodHead
	OPTIONS         = http.MethodOptions
	ContentTypeJSON = "application/json"
	ContentTypeUrl  = "application/x-www-form-urlencoded"
	ContentTypeForm = "multipart/form-data"
)

//Request Http Request
type Request struct {
	url     string
	method  string
	header  http.Header
	params  cmap.CMap // only get/head/delete
	body    io.Reader
	Client  *http.Client
	cookies []*http.Cookie
	f       file
}

type file struct {
	*fs.File
	field string
}

// NewRequest new request
func NewRequest(method, urlString string) *Request {
	var r = &Request{}
	r.method = strings.ToUpper(method)
	r.url = urlString
	r.params = cmap.NewCMap()
	r.header = http.Header{}
	r.Client = &http.Client{
		Timeout: time.Second * 10,
	}
	r.SetContentType(ContentTypeJSON)
	return r
}

func (m *Request) SetContentType(contentType string) *Request {
	m.SetHeader("Content-Type", contentType)
	return m
}

func (m *Request) AddHeader(key, value string) *Request {
	m.header.Add(key, value)
	return m
}

func (m *Request) SetHeader(key, value string) *Request {
	m.header.Set(key, value)
	return m
}

func (m *Request) SetHeaders(header http.Header) *Request {
	m.header = header
	return m
}

func (m *Request) SetBody(body io.Reader) *Request {
	m.body = body
	m.params = nil
	return m
}

func (m *Request) AddParam(key, value string) *Request {
	m.params.Add(key, value)
	m.body = nil
	return m
}

func (m *Request) SetParam(key, value string) *Request {
	m.params.Set(key, value)
	m.body = nil
	return m
}

// SetStructParams struct to Params
func (m *Request) SetStructParams(v any) *Request {

	m.params = cmap.StructToCMap(v)
	return m
}

// SetParams Get params
func (m *Request) SetParams(params cmap.CMap) *Request {
	m.params = params
	return m
}

func (m *Request) AddFile(field, fileName, path string) (err error) {
	m.f.File, err = fs.OpenFile(path)
	if err != nil {
		return
	}
	m.f.field = field
	m.f.SetName(fileName)
	return
}

func (m *Request) RemoveFile() *Request {
	m.f.File = nil
	return m
}

func (m *Request) AddCookie(cookie *http.Cookie) *Request {
	m.cookies = append(m.cookies, cookie)
	return m
}

func (m *Request) Exec() *Response {
	var req *http.Request
	var err error
	var body io.Reader
	var rawQuery string

	if len(m.params) > 0 {
		rawQuery = m.params.Encode()
	}

	if m.body != nil {
		body = m.body
	} else if m.f.File != nil {
		bodyByte := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyByte)

		f, _ := writer.CreateFormFile(m.f.field, m.f.Name())
		_, _ = io.Copy(f, m.f)

		for key, values := range m.params {
			for _, value := range values {
				_ = writer.WriteField(key, value)
			}
		}

		err = writer.Close()
		if err != nil {
			return &Response{nil, nil, err}
		}

		m.SetContentType(writer.FormDataContentType())
		body = bodyByte
	} else if m.params != nil {
		body = strings.NewReader(m.params.Encode())
	}

	req, err = http.NewRequest(m.method, m.url, body)
	if err != nil {
		return &Response{nil, nil, err}
	}
	if len(rawQuery) > 0 {
		req.URL.RawQuery = rawQuery
	}
	req.Header = m.header

	for _, cookie := range m.cookies {
		req.AddCookie(cookie)
	}

	resp, err := m.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return &Response{resp, nil, err}
	}

	data, err := ioutil.ReadAll(resp.Body)
	return &Response{resp, data, err}
}
