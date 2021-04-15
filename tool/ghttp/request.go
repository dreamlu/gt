package ghttp

import (
	"bytes"
	"github.com/dreamlu/gt/tool/type/cmap"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
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
	ContentTypeForm = "application/x-www-form-urlencoded"
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
	file    *file
}

type file struct {
	name     string
	filename string
	path     string
}

//NewRequest 新的Request指针
func NewRequest(method, urlString string) *Request {
	var r = &Request{}
	r.method = strings.ToUpper(method)
	r.url = urlString
	r.params = cmap.NewCMap()
	r.header = http.Header{}
	r.Client = &http.Client{
		Timeout: time.Second * 10,
	}
	r.SetContentType(ContentTypeForm)
	return r
}

//SetContentType 设定Content-Type
func (m *Request) SetContentType(contentType string) *Request {
	m.SetHeader("Content-Type", contentType)
	return m
}

//AddHeader 增加Header头
func (m *Request) AddHeader(key, value string) *Request {
	m.header.Add(key, value)
	return m
}

//SetHeader 设定Header头
func (m *Request) SetHeader(key, value string) *Request {
	m.header.Set(key, value)
	return m
}

//SetHeaders 设定Header头
func (m *Request) SetHeaders(header http.Header) *Request {
	m.header = header
	return m
}

//SetBody 设定POST内容
func (m *Request) SetBody(body io.Reader) *Request {
	m.body = body
	m.params = nil
	return m
}

//AddParam 增加Get请求参数
func (m *Request) AddParam(key, value string) *Request {
	m.params.Add(key, value)
	m.body = nil
	return m
}

//SetParam 设定Get请求参数
func (m *Request) SetParam(key, value string) *Request {
	m.params.Set(key, value)
	m.body = nil
	return m
}

// SetStructParams struct to Params
func (m *Request) SetStructParams(v interface{}) *Request {

	m.params = cmap.StructToCMap(v)
	return m
}

//SetParams 设定Get请求参数
func (m *Request) SetParams(params cmap.CMap) *Request {
	m.params = params
	return m
}

//AddFile 增加文件
func (m *Request) AddFile(name, filename, path string) *Request {
	m.file = &file{name, filename, path}
	return m
}

//RemoveFile 移除文件
func (m *Request) RemoveFile() *Request {
	m.file = nil
	return m
}

//AddCookie 添加COOKIE
func (m *Request) AddCookie(cookie *http.Cookie) *Request {
	m.cookies = append(m.cookies, cookie)
	return m
}

//Exec 发送HTTP请求
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
	} else if m.file != nil {
		uploadFile, err := os.Open(m.file.path)
		if err != nil {
			return &Response{nil, nil, err}
		}
		defer uploadFile.Close()

		bodyByte := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyByte)
		part, err := writer.CreateFormFile(m.file.name, m.file.filename)
		if err != nil {
			return &Response{nil, nil, err}
		}
		_, err = io.Copy(part, uploadFile)
		if err != nil {
			return &Response{nil, nil, err}
		}

		for key, values := range m.params {
			for _, value := range values {
				writer.WriteField(key, value)
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
		return &Response{nil, nil, err}
	}

	data, err := ioutil.ReadAll(resp.Body)
	return &Response{resp, data, err}
}
