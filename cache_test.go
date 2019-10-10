// package der

package der

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var r = RedisManager{}

// cache test
var cache, _ = NewCache(new(RedisManager))

func init() {
	// init redis
	_ = r.NewCache()
	// init cache
	//_ = cache.NewCache()

}

// redis method set test
func TestRedis(t *testing.T) {
	err := r.Rc.Set("test", "testValue").Err()
	log.Println("set err:", err)
	value := r.Rc.Get("test")
	reqRes, _ := value.Result()
	log.Println("value", reqRes)
}

// user model
var user = User{
	ID:   1,
	Name: "test",
	//Createtime: JsonDate(time.Now()),
}

// set and get interface value
func TestCacheRedis(t *testing.T) {
	// data
	data := CacheModel{
		Time: 50 * CacheMinute,
		Data: user,
	}

	// key can use user.ID,user.Name,user
	// because it can be interface
	// set
	err := cache.Set(user, data)
	log.Println("set err: ", err)

	// get
	reply, _ := cache.Get(user)
	log.Println("user data :", reply.Data)

}

// check or delete cache
func TestCacheCheckDelRedis(t *testing.T) {
	// check
	//err := cache.Check(user.ID)
	//log.Println("check: ", err)

	// del
	//err := cache.Delete(user.ID)
	//log.Println("delete: ", err)

	// del *

	//err := cache.Delete("1*")
	//log.Println("delete: ", err)

	// del more
	err := cache.DeleteMore(user)
	log.Println("delete: ", err)
}

// cookie test
func TestCookie(t *testing.T) {

	recorder := httptest.NewRecorder()

	// Drop a cookie on the recorder.
	http.SetCookie(recorder, &http.Cookie{Name: "test", Value: "test"})

	// Copy the Cookie over to a new Request
	request := &http.Request{Header: http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}}

	// Extract the dropped cookie from the request.
	cookie, _ := request.Cookie("test")
	log.Println(cookie)

}
