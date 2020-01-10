// package gt

package test

import (
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/cache/redis"
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
	testConfDir = "../../conf/"
	r           = redis.RedisManager{}
	ce, _       = cache.NewCache(new(redis.RedisManager), testConfDir)
)

func init() {
	// init redis
	_ = r.NewCache(testConfDir)
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
	data := cache.CacheModel{
		Time: 50 * cache.CacheMinute,
		Data: user,
	}

	// key can use user.ID,user.Name,user
	// because it can be interface
	// set
	err := ce.Set(user, data)
	log.Println("set err: ", err)

	// get
	reply, _ := ce.Get(user)
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
	err := ce.DeleteMore(user)
	log.Println("delete: ", err)
}
