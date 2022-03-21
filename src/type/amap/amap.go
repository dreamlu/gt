package amap

import (
	"encoding/json"
)

type AMap map[string]string

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v AMap) Get(key string) string {
	if v == nil {
		return ""
	}
	return v[key]
}

// Pop return Get value and Del key
func (v AMap) Pop(key string) string {
	s := v.Get(key)
	v.Del(key)
	return s
}

// Set sets the key to value. It replaces string existing
// values.
func (v AMap) Set(key string, value string) AMap {
	v[key] = value
	return v
}

// Del deletes the values associated with key.
func (v AMap) Del(key string) AMap {
	delete(v, key)
	return v
}

// Marshal AMap to v
func (v AMap) Marshal(value any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, value)
	if err != nil {
		return err
	}
	return nil
}

// ToAMap struct/slice... to AMap
// v must be allowed
func ToAMap(v any) (values AMap) {
	values = NewAMap()
	bs, _ := json.Marshal(v)
	_ = json.Unmarshal(bs, &values)
	return
}

func NewAMap() AMap {
	return AMap{}
}
