// package gt

package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dreamlu/gt"
	redis2 "github.com/dreamlu/gt/cache/redis"
	"github.com/go-redis/redis"
)

// impl cache manager
// redis cache
// interface key, interface value
type RedisManager struct {
	// do nothing else
	Rc *redis2.ConnPool
}

// new cache by redis
// other cache maybe like this
func (r *RedisManager) NewCache(params ...interface{}) error {

	var config *gt.Config
	if len(params) > 0 {
		config = gt.Configger(params[0].(string))
	} else {
		config = gt.Configger()
	}

	// read config
	r.Rc = redis2.InitRedisPool(
		func(options *redis.Options) {
			config.GetStruct("app.redis", options)
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
		return reply, errors.New("no such cache key:" + string(keyS))
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

	return r.Rc.ExpireKey(keyS, reply.Time*60).Err()
}
