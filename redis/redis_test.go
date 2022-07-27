// package gt

package redis

import (
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/src/type/time"
	"testing"
)

type User struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Createtime time.CTime `json:"createtime"`
}

func init() {
	var opt Options
	conf.UnmarshalField(cons.ConfRedis, &opt)
	OpenRedis(&opt)
}

// redis method set test
func TestRedis(t *testing.T) {
	err := Set("test", "testValue").Err()
	if err != nil {
		t.Fatal(err)
	}
	value := Get("test")
	reqRes, _ := value.Result()
	t.Log("value", reqRes)
}

// user model
var user = User{
	ID:   1,
	Name: "test",
	//Createtime: JsonDate(time.Now()),
}

func TestMarshal(t *testing.T) {
	err := SetMarshal("user", user)
	if err != nil {
		t.Fatal(err)
	}
	var data User
	err = GetMarshal("user", &data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("user", data)
}
