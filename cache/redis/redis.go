// package gt

package redis

import (
	"github.com/go-redis/redis"
	"log"
)

// ConnPool is RDS struct
type ConnPool struct {
	redisDB *redis.Client
}

type Options func(*redis.Options)

// InitRedisPool func init RDS fd
func InitRedisPool(options ...Options) *ConnPool {
	r := &ConnPool{}
	option := &redis.Options{}
	for _, o := range options {
		o(option)
	}
	r.redisDB = redis.NewClient(option)
	//r.redisDB.Ping()
	return r
}

// Close pool
func (p *ConnPool) Close() {
	err := p.redisDB.Close()
	if err != nil {
		log.Println("[Redis Error]: ", err)
	}
}

// Do commands
func (p *ConnPool) Do(args ...interface{}) *redis.Cmd {
	// close problem
	//defer p.Close()
	return p.redisDB.Do(args[:]...)
}

// Set
func (p *ConnPool) Set(key interface{}, value interface{}) *redis.Cmd {
	//defer p.Close()
	return p.Do("SET", key, value)
}

// Get
func (p *ConnPool) Get(key interface{}) *redis.Cmd {
	// get one connection from pool
	//defer p.Close()
	return p.Do("GET", key)
}

// keys
func (p *ConnPool) Keys(keys interface{}) *redis.Cmd {
	// get one connection from pool
	//defer p.Close()
	return p.Do("KEYS", keys)
}

// DelKey for key
func (p *ConnPool) Delete(key interface{}) *redis.Cmd {
	//defer p.Close()
	return p.Do("DEL", key)
}

// ExpireKey for key
func (p *ConnPool) ExpireKey(key interface{}, seconds int64) *redis.Cmd {
	//defer p.Close()
	return p.Do("EXPIRE", key, seconds)
}
