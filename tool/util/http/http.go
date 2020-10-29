package http

import (
	"github.com/dreamlu/gt/tool/log"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	GET             = "GET"
	POST            = "POST"
	DELETE          = "DELETE"
	PUT             = "PUT"
	PATCH           = "PATCH"
	HEAD            = "HEAD"
	OPTIONS         = "OPTIONS"
	ContentTypeJSON = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

type Request struct {
	*http.Request
}

func NewRequest(method, url string, body io.Reader) *Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Error(req)
		return nil
	}
	return &Request{
		req,
	}
}

func (r *Request) Do() (b []byte, err error) {
	res, err := http.DefaultClient.Do(r.Request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	return
}

// request sugar

// Get url res
func Get(url string) (b []byte, err error) {

	b, err = NewRequest(GET, url, nil).Do()
	return
}

// Post
func Post(url string, contentType string, body io.Reader) (b []byte, err error) {

	r := NewRequest(POST, url, body)
	r.Header.Set("Content-Type", contentType)
	b, err = r.Do()
	return
}

// DELETE
func Delete(url string) (b []byte, err error) {
	b, err = NewRequest(DELETE, url, nil).Do()
	return
}
