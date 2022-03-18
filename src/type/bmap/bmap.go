package bmap

import (
	"encoding/json"
)

type BMap map[string]any

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v BMap) Get(key string) any {
	if v == nil {
		return ""
	}
	return v[key]
}

// Pop return Get value and Del key
func (v BMap) Pop(key string) any {
	s := v.Get(key)
	v.Del(key)
	return s
}

// Set sets the key to value. It replaces any existing
// values.
func (v BMap) Set(key string, value any) BMap {
	v[key] = value
	return v
}

// Del deletes the values associated with key.
func (v BMap) Del(key string) BMap {
	delete(v, key)
	return v
}

// Marshal BMap to v
func (v BMap) Marshal(value any) error {
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

// ToBMap struct/slice... to BMap
// v must be allowed
func ToBMap(v any) (values BMap) {
	values = NewBMap()
	bs, _ := json.Marshal(v)
	_ = json.Unmarshal(bs, &values)
	return
}

func NewBMap() BMap {
	return BMap{}
}
