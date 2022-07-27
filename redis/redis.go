// package gt

package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

// Redis is RDS struct
type Redis struct {
	redisClient *redis.Client
}

type Options struct {
	redis.Options
}

// NewRedis new redis
func NewRedis(option *Options) *Redis {
	return &Redis{
		redisClient: redis.NewClient(&option.Options),
	}
}

// Close pool
func (r *Redis) Close() error {
	return r.redisClient.Close()
}

// Do command
func (r *Redis) Do(args ...any) *redis.Cmd {
	return r.redisClient.Do(context.TODO(), args...)
}

func (r *Redis) Set(key any, value any) *redis.Cmd {
	return r.Do("SET", key, value)
}

func (r *Redis) Get(key any) *redis.Cmd {
	// get one connection from pool
	return r.Do("GET", key)
}

func (r *Redis) Keys(keys any) *redis.Cmd {
	// get one connection from pool
	return r.Do("KEYS", keys)
}

func (r *Redis) Delete(key any) *redis.Cmd {
	return r.Do("DEL", key)
}

// ExpireKey for key
func (r *Redis) ExpireKey(key any, seconds int64) *redis.Cmd {
	return r.Do("EXPIRE", key, seconds)
}

// sugar

func (r *Redis) SetMarshal(key any, data any) error {

	// change key to string
	keyS, err := json.Marshal(key)
	if err != nil {
		return err
	}

	// can not store struct data
	// change data to string
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// set string data
	err = r.Set(keyS, bs).Err()
	if err != nil {
		return err
	}
	//if value.Time != 0 {
	//	err = r.ExpireKey(keyS, value.Time).Err()
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (r *Redis) GetMarshal(key any, v any) error {

	// change key to string
	keyS, err := json.Marshal(key)
	if err != nil {
		return err
	}

	// data
	res := r.Get(keyS).Val()
	if res == nil {
		return nil
	}

	// string to struct data
	return json.Unmarshal([]byte(res.(string)), v)
}

// DeleteM delete like all *key* values
func (r *Redis) DeleteM(key any) error {

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
	res := r.Keys(buf.Bytes()).Val()
	if res != nil {
		for _, v := range res.([]any) {
			err = r.Delete(v).Err()
			if err != nil {
				return err
			}
		}

	}
	return nil
}
