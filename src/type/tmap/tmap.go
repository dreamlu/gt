package tmap

import (
	"encoding/json"
)

type TI interface {
	any
}

type TMap[T TI] map[string]T

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v TMap[T]) Get(key string) T {
	if v == nil {
		var zero T // zero
		return zero
	}
	return v[key]
}

// Pop return Get value and Del key
func (v TMap[T]) Pop(key string) T {
	s := v.Get(key)
	v.Del(key)
	return s
}

// Set sets the key to value. It replaces any existing
// values.
func (v TMap[T]) Set(key string, value T) TMap[T] {
	v[key] = value
	return v
}

// Del deletes the values associated with key.
func (v TMap[T]) Del(key string) TMap[T] {
	delete(v, key)
	return v
}

// Marshal TMap to v
func (v TMap[T]) Marshal(value T) error {
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
func ToTMap[T TI](v any) TMap[T] {
	values := NewTMap[T]()
	bs, _ := json.Marshal(v)
	_ = json.Unmarshal(bs, &values)
	return values
}

func NewTMap[T TI]() TMap[T] {
	return TMap[T]{}
}
