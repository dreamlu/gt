// package gt

package cache

import (
	"github.com/dreamlu/gt"
)

// data model
type CacheModel struct {
	// seconds
	Time int64 `json:"time"`
	// data
	Data interface{} `json:"data"`
}

// cache manager
type Cache interface {
	// init cache
	NewCache(params ...interface{}) error
	// operate method
	// set value
	// if time != 0 set it
	Set(key interface{}, value CacheModel) error
	// get value
	Get(key interface{}) (CacheModel, error)
	// delete value
	Delete(key interface{}) error
	// more del
	// key will become *key*
	DeleteMore(key interface{}) error
	// check value
	// flush the time
	Check(key interface{}) error
}

// time for cache unit
// unit: second
const (
	CacheSecond = 1
	CacheMinute = 60
	CacheHour   = 60 * CacheMinute
	CacheDay    = 24 * CacheHour
	CacheWeek   = 7 * CacheDay
)

// cache sugar
// the first param is Cache
// the second param is confiDir
func NewCache(params ...interface{}) (cache Cache) {

	// default set
	if len(params) == 0 {
		cache = new(RedisManager)
		err := cache.NewCache()
		if err != nil {
			gt.Logger().Error(err.Error())
		}
		return
	}
	// init
	cache = params[0].(Cache)
	err := cache.NewCache(params[1:]...)
	if err != nil {
		gt.Logger().Error(err.Error())
	}
	return
}
