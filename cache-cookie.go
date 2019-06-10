// @author  dreamlu
package der

import (
	"errors"
	"net/http"
)

// impl cache manager
// cookie cache
// interface key, interface value
type CookieManager struct {
	// http
	Writer  http.ResponseWriter
	Request *http.Request
}

// new cache by redis
// other cacher maybe have this too
func (c *CookieManager) NewCache(args ...interface{}) error {
	c.Writer = args[0].(http.ResponseWriter)
	c.Request = args[1].(*http.Request)
	// err
	return nil
}

// cookie value must be string value
func (c *CookieManager) Set(key interface{}, value CacheModel) error {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     key.(string),
		Value:    value.Data.(string),
		//Domain:   "*",
		MaxAge:   int(value.Time),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})
	return nil
}

func (c *CookieManager) Get(key interface{}) (reply CacheModel, err error) {

	cookie, err := c.Request.Cookie(key.(string))
	if err != nil {
		return
	}

	reply.Data = cookie.Value
	reply.Time = int64(cookie.MaxAge)
	return
}

func (c *CookieManager) Delete(key interface{}) error {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:   key.(string),
		Domain: "*",
		MaxAge: -1,
	})
	return nil
}

// cookie not support this!
func (c *CookieManager) DeleteMore(key interface{}) error {

	return errors.New("cookie not support DeleteMore method")
}

func (c *CookieManager) Check(key interface{}) error {

	_, err := c.Request.Cookie(key.(string))
	if err != nil {
		return err
	}

	return nil
}
