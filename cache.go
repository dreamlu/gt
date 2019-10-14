// package gt

package gt

// data model
type CacheModel struct {
	// seconds
	Time int64 `json:"time"`
	// data
	Data interface{} `json:"data"`
}

// cache manager
type CacheManager interface {
	// init cache
	NewCache(args ...interface{}) error
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
	CacheMinute = 60
	CacheHour   = 60 * CacheMinute
	CacheDay    = 24 * CacheHour
	CacheWeek   = 7 * CacheDay
)

// cache sugar
func NewCache(cache CacheManager, args ...interface{}) (CacheManager, error) {
	err := cache.NewCache(args)
	return cache, err
}
