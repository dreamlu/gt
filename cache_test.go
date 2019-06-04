// @author  dreamlu
package der

import (
	"log"
	"testing"
)

var r = RedisManager{}
// cache test
var cache CacheManager = new(RedisManager)

func init()  {
	// init redis
	_ = r.NewCache()
	// init cache
	_ = cache.NewCache()
}

// redis method set test
func TestRedis(t *testing.T) {
	err := r.Rc.Set("test", "testValue").Err()
	log.Println("set err:", err)
	value := r.Rc.Get("test")
	reqRes,_ := value.Result()
	log.Println("value",reqRes)
}

// user model
var user = User{
	ID:   1,
	Name: "test",
	//Createtime: JsonDate(time.Now()),
}

// set and get interface value
func TestCache(t *testing.T) {
	// data
	data := CacheModel{
		Time: 50,
		Data: user,
	}

	// key can use user.ID,user.Name,user
	// because it can be interface
	// set
	err := cache.Set(user, data)
	log.Println("set err: ", err)

	// get
	reply,_ := cache.Get(user)
	log.Println("user data :", reply.Data)

}

// check or delete cache
func  TestCacheCheckDel(t *testing.T)  {
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
