// package gt

package cache

import (
	"github.com/dreamlu/gt/tool/type/time"
	"log"
	"testing"
)

type User struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Createtime time.CTime `json:"createtime"`
}

var (
	r  = RedisManager{}
	ce = NewCache()
)

func init() {
	// init redis
	_ = r.NewCache()
}

func TestCacheExpireKey(t *testing.T) {
	t.Log(ce.ExpireKey("test1", CacheHour))
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
	err := NewCache().Set(user, data)
	t.Log("set err: ", err)

	// get
	var user2 User
	reply, _ := ce.Get(user)
	t.Log(reply.Struct(&user2))
	t.Log("user data :", user2)

	var ar = []string{"test1", "test2"}
	data.Data = ar
	err = NewCache().Set("arr", data)
	t.Log("set err: ", err)

	// get
	ar = []string{}
	reply, _ = ce.Get("arr")
	t.Log(reply.Unmarshal(&ar))
	t.Log("user data :", user2)
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
	err := ce.DeleteMore(user)
	log.Println("delete: ", err)
}
