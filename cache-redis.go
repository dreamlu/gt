// package der

package der

import (
	"bytes"
	"encoding/json"
)

// impl cache manager
// redis cache
// interface key, interface value
type RedisManager struct {
	// do nothing else
	Rc *ConnPool
}

// toMe: wait the future complete it
// new cache by redis
// other cacher maybe have this too
func (r *RedisManager) NewCache(args ...interface{}) error {

	config := Configger()

	var (
		Host         = config.GetString("app.redis.host")
		Password     = config.GetString("app.redis.password")
		Database     = config.GetInt("app.redis.database")
		poolSize     = config.GetInt("app.redis.poolSize")     // max number of connections
		MinIdleConns = config.GetInt("app.redis.minIdleConns") // 最大的空闲连接数
	)

	// default value

	r.Rc = InitRedisPool(Host, Password, Database, poolSize, MinIdleConns)
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
