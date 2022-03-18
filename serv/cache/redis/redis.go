// package gt

package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
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
func (p *ConnPool) Do(args ...any) *redis.Cmd {
	// close problem
	//defer p.Close()
	return p.redisDB.Do(context.TODO(), args...)
}

func (p *ConnPool) Set(key any, value any) *redis.Cmd {
	//defer p.Close()
	return p.Do("SET", key, value)
}

func (p *ConnPool) Get(key any) *redis.Cmd {
	// get one connection from pool
	//defer p.Close()
	return p.Do("GET", key)
}

func (p *ConnPool) Keys(keys any) *redis.Cmd {
	// get one connection from pool
	//defer p.Close()
	return p.Do("KEYS", keys)
}

func (p *ConnPool) Delete(key any) *redis.Cmd {
	//defer p.Close()
	return p.Do("DEL", key)
}

// ExpireKey for key
func (p *ConnPool) ExpireKey(key any, seconds int64) *redis.Cmd {
	//defer p.Close()
	return p.Do("EXPIRE", key, seconds)
}
