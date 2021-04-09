// package gt

package cache

import (
	"encoding/json"
	"github.com/dreamlu/gt/tool/log"
)

// data model
type CacheModel struct {
	// seconds
	Time int64 `json:"time,omitempty"`
	// data
	Data interface{} `json:"data,omitempty"`
}

// Deprecated
// c.Data to v
func (c CacheModel) Struct(v interface{}) error {
	b, err := json.Marshal(c.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// Unmarshal support Struct/Array
// c.Data to v
func (c CacheModel) Unmarshal(v interface{}) error {
	b, err := json.Marshal(c.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// cache manager
type Cache interface {
	// init cache
	NewCache() error
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
	// expire key time
	ExpireKey(key interface{}, seconds int64) bool
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
// the second param is confDir
func NewCache(params ...interface{}) (cache Cache) {

	// default set
	if len(params) == 0 {
		cache = new(RedisManager)
		err := cache.NewCache()
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
	// init
	cache = params[0].(Cache)
	err := cache.NewCache()
	if err != nil {
		log.Error(err.Error())
	}
	return
}
