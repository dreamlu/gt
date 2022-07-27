// package gt

package redis

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	r    *Redis
	once sync.Once
)

// OpenRedis open once redis client
func OpenRedis(option *Options) {
	once.Do(func() {
		r = NewRedis(option)
	})
}

func GetRedis() *Redis {
	return r
}

// Close pool
func Close() error {
	return r.Close()
}

// Do command
func Do(args ...any) *redis.Cmd {
	return r.Do(args...)
}

func Set(key any, value any) *redis.Cmd {
	return r.Set(key, value)
}

func Get(key any) *redis.Cmd {
	return r.Get(key)
}

func Keys(keys any) *redis.Cmd {
	return r.Keys(keys)
}

func Delete(key any) *redis.Cmd {
	return r.Delete(key)
}

// ExpireKey for key
func ExpireKey(key any, seconds int64) *redis.Cmd {
	return r.ExpireKey(key, seconds)
}

func SetMarshal(key any, data any) error {
	return r.SetMarshal(key, data)
}

func GetMarshal(key any, v any) error {
	return r.GetMarshal(key, v)
}

// DeleteM delete like all *key* values
func DeleteM(key any) error {
	return r.DeleteM(key)
}
