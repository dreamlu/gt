package tmap

import (
	"encoding/json"
)

type TI interface {
	any
}

type TMap[K comparable, V TI] map[K]V

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v TMap[K, V]) Get(key K) V {
	if v == nil {
		var zero V // zero
		return zero
	}
	return v[key]
}

// Pop return Get value and Del key
func (v TMap[K, V]) Pop(key K) V {
	s := v.Get(key)
	v.Del(key)
	return s
}

// Set sets the key to value. It replaces any existing
// values.
func (v TMap[K, V]) Set(key K, value V) TMap[K, V] {
	v[key] = value
	return v
}

// Del deletes the values associated with key.
func (v TMap[K, V]) Del(key K) TMap[K, V] {
	delete(v, key)
	return v
}

// Marshal TMap to v
func (v TMap[K, V]) Marshal(value V) error {
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

// ToTMap struct/slice... to TMap
// v must be allowed
func ToTMap[K comparable, V TI](v any) TMap[K, V] {
	values := NewTMap[K, V]()
	bs, _ := json.Marshal(v)
	_ = json.Unmarshal(bs, &values)
	return values
}

func NewTMap[K comparable, V TI]() TMap[K, V] {
	return TMap[K, V]{}
}
