// package gt

package cache

import (
	"encoding/json"
	"github.com/dreamlu/gt/third/log"
)

// CacheModel data model
type CacheModel struct {
	// seconds
	Time int64 `json:"time,omitempty"`
	// data
	Data any `json:"data,omitempty"`
}

// Unmarshal support Struct/Array
// c.Data to v
func (c CacheModel) Unmarshal(v any) error {
	b, err := json.Marshal(c.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// Cache manager
type Cache interface {
	// Init init cache
	Init() error
	// operate method

	// Set value
	// if time != 0 set it
	Set(key any, value CacheModel) error
	// Get value
	Get(key any) (CacheModel, error)
	// Delete value
	Delete(key any) error
	// DeleteMore more del
	// key will become *key*
	DeleteMore(key any) error
	// Check value
	// flush the time
	Check(key any) error
	// ExpireKey expire key time
	ExpireKey(key any, seconds int64) bool
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

// NewCache cache sugar
// the first param is Cache
// the second param is confDir
func NewCache(params ...any) (cache Cache) {

	// default set
	if len(params) == 0 {
		cache = new(RedisManager)
		err := cache.Init()
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
	// init
	cache = params[0].(Cache)
	err := cache.Init()
	if err != nil {
		log.Error(err.Error())
	}
	return
}

// NewRedis new your custom db crud
func NewRedis(redis *RedisManager) (cache Cache) {

	cache = redis
	err := cache.Init()
	if err != nil {
		log.Error(err.Error())
	}
	return
}
