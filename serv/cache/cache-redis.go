// package gt

package cache

import (
	"bytes"
	"encoding/json"
	redis2 "github.com/dreamlu/gt/serv/cache/redis"
	"github.com/dreamlu/gt/serv/conf"
	"github.com/dreamlu/gt/serv/log"
	"github.com/dreamlu/gt/tool/cons"
	"github.com/go-redis/redis/v8"
)

// RedisManager impl cache manager
// redis cache
// interface key, interface value
type RedisManager struct {
	// do nothing else
	Rc *redis2.ConnPool
}

type redisOptions struct {
	// host:port address.
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	// Database to be selected after connecting to the server.
	DB int `yaml:"DB"`
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int `yaml:"poolSize"`

	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `yaml:"minIdleConns"`
}

func (r *RedisManager) NewCache() error {

	// read config
	r.Rc = redis2.InitRedisPool(
		func(options *redis.Options) {
			var opt redisOptions
			conf.GetStruct(cons.ConfRedis, &opt)
			options.Addr = opt.Addr
			options.Username = opt.Username
			options.Password = opt.Password
			options.DB = opt.DB
			options.PoolSize = opt.PoolSize
			options.MinIdleConns = opt.MinIdleConns
		})
	return nil
}

func (r *RedisManager) Set(key interface{}, value CacheModel) error {

	// change key to string
	keyS, err := json.Marshal(key)
	if err != nil {
		return err
	}

	// can not store struct data
	// change data to string
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// set string data
	err = r.Rc.Set(keyS, data).Err()
	if err != nil {
		return err
	}
	if value.Time != 0 {
		err = r.Rc.ExpireKey(keyS, value.Time).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisManager) Get(key interface{}) (CacheModel, error) {

	var reply CacheModel

	// change key to string
	keyS, err := json.Marshal(key)
	if err != nil {
		return reply, err
	}

	// data
	res := r.Rc.Get(keyS).Val()
	if res == nil {
		return reply, nil
	}

	// string to struct data
	err = json.Unmarshal([]byte(res.(string)), &reply)
	if err != nil {
		return reply, err
	}

	return reply, nil
}

func (r *RedisManager) Delete(key interface{}) error {

	// change key to string
	keyS, err := json.Marshal(key)
	if err != nil {
		return err
	}

	return r.Rc.Delete(keyS).Err()
}

func (r *RedisManager) DeleteMore(key interface{}) error {

	// change key to string
	keyS, err := json.Marshal(key)
	if err != nil {
		return err
	}

	var (
		buf bytes.Buffer
	)
	buf.WriteString("*")
	buf.Write(keyS)
	buf.WriteString("*")

	// keys
	res := r.Rc.Keys(buf.Bytes()).Val()
	if res != nil {
		for _, v := range res.([]interface{}) {
			err := r.Rc.Delete(v).Err()
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (r *RedisManager) Check(key interface{}) error {

	var reply CacheModel

	// change key to string
	keyS, err := json.Marshal(key)
	if err != nil {
		return err
	}

	// data
	res := r.Rc.Get(keyS).Val()

	// string to struct data
	err = json.Unmarshal([]byte(res.(string)), &reply)
	if err != nil {
		return err
	}

	return r.Rc.ExpireKey(keyS, reply.Time).Err()
}

func (r *RedisManager) ExpireKey(key interface{}, seconds int64) bool {
	b, err := r.Rc.ExpireKey(key, seconds).Bool()
	if err != nil {
		log.Error(err)
	}
	return b
}
