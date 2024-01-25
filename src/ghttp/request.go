package ghttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/dreamlu/gt/src/file/fs"
	"github.com/dreamlu/gt/src/type/cmap"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const (
	GET                 = http.MethodGet
	POST                = http.MethodPost
	DELETE              = http.MethodDelete
	PUT                 = http.MethodPut
	PATCH               = http.MethodPatch
	HEAD                = http.MethodHead
	OPTIONS             = http.MethodOptions
	ContentTypeJSON     = "application/json"
	ContentTypeFormUrl  = "application/x-www-form-urlencoded"
	ContentTypeFormData = "multipart/form-data"
)

// Request Http Request
type Request struct {
	url       string
	method    string
	header    http.Header
	urlValues cmap.CMap // only get/head/delete
	forms     cmap.CMap // post/put/patch form-data
	body      io.Reader
	Client    *http.Client
	cookies   []*http.Cookie
	f         file
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
	r.urlValues = cmap.NewCMap()
	r.forms = cmap.NewCMap()
	r.header = http.Header{}
	r.Client = &http.Client{}
	r.SetContentType(ContentTypeJSON)
	return r
}

func (m *Request) SetTimeout(timeout time.Duration) *Request {
	m.Client.Timeout = timeout
	return m
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

func (m *Request) SetHeaders(headers any) *Request {
	switch headers.(type) {
	case http.Header:
		m.header = headers.(http.Header)
	case cmap.CMap:
		for k, v := range headers.(cmap.CMap) {
			m.SetHeader(k, v[0])
		}
	}
	return m
}

func (m *Request) SetBody(body io.Reader) *Request {
	m.body = body
	return m
}

func (m *Request) SetJsonBody(v any) *Request {
	if v == nil {
		v = cmap.NewCMap()
	}
	bs, _ := json.Marshal(v)
	m.body = bytes.NewReader(bs)
	return m
}

func (m *Request) SetUrlValue(key, value string) *Request {
	m.urlValues.Set(key, value)
	return m
}

// SetUrlValues struct to Params
func (m *Request) SetUrlValues(v any) *Request {
	m.urlValues = cmap.StructToCMap(v)
	return m
}

func (m *Request) SetForm(key, value string) *Request {
	m.forms.Set(key, value)
	return m
}

// SetForms struct to Params
func (m *Request) SetForms(v any) *Request {
	m.forms = cmap.StructToCMap(v)
	return m
}

func (m *Request) SetFile(field, fileName, path string) (err error) {
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
	var (
		req      *http.Request
		err      error
		body     io.Reader
		rawQuery string
	)

	// url params
	if len(m.urlValues) > 0 {
		rawQuery = m.urlValues.Encode()
	}
	// json/form
	if m.body != nil {
		body = m.body
	}
	// form-data
	if m.f.File != nil || len(m.forms) > 0 {
		bs := &bytes.Buffer{}
		writer := multipart.NewWriter(bs)
		for key := range m.forms {
			_ = writer.WriteField(key, m.forms.Get(key))
		}
		if m.f.File != nil {
			f, _ := writer.CreateFormFile(m.f.field, m.f.Name())
			_, _ = io.Copy(f, m.f)
			err = writer.Close()
			if err != nil {
				return &Response{nil, nil, err}
			}
			m.SetContentType(writer.FormDataContentType())
		}
		body = bs
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

	// skip verify ssl https
	tc := &tls.Config{
		InsecureSkipVerify: true,
	}
	m.Client.Transport = &http.Transport{
		TLSClientConfig: tc,
	}

	resp, err := m.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return &Response{resp, nil, err}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{resp, data, err}
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New("http status: " + resp.Status)
	}
	return &Response{resp, data, err}
}
