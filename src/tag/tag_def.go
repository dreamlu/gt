package tag

import (
	"fmt"
	"github.com/dreamlu/gt/src/type/amap"
	"reflect"
)

type GtField struct {
	Field string
	Type  string
}

// GtTag A GT tag
type GtTag struct {
	Name  string
	Value string
}

type GtTags struct {
	origin string
	GtTags []*GtTag
}

func (g GtTags) ToAMap() amap.AMap {
	var ap = amap.NewAMap()
	for _, f := range g.GtTags {
		ap.Set(f.Name, f.Value)
	}
	return ap
}

type TFunc[T any] func(p reflect.Type) GF[T]
type FFunc[T any] func(p reflect.StructField) GF[T]

type GF[T any] map[GtField]T

func NewGF[T any]() GF[T] {
	return GF[T]{}
}

func (g GF[T]) zero() T {
	var zero T // zero
	return zero
}

func (g GF[T]) Get(key string) T {
	for k, value := range g {
		if k.Field == key {
			return value
		}
	}
	return g.zero()
}

// Pop return Get value and Del key
func (g GF[T]) Pop(key string) T {
	s := g.Get(key)
	g.Del(key)
	return s
}

// Top return pop one value
func (g GF[T]) Top() T {
	for _, value := range g {
		return value
	}
	return g.zero()
}

// Set sets the key to value. It replaces string existing
// values.
func (g GF[T]) Set(key GtField, value T) GF[T] {
	g[key] = value
	return g
}

// Del deletes the values associated with key.
func (g GF[T]) Del(key string) GF[T] {
	for k := range g {
		if k.Field == key {
			delete(g, k)
		}
	}
	return g
}

func (g GF[T]) ToAMap() amap.AMap {
	var am = amap.NewAMap()
	for k, v := range g {
		am.Set(k.Field, fmt.Sprint(v))
	}
	return am
}

func (g GF[T]) Parse(f func(T) T) GF[T] {
	for k, v := range g {
		g.Set(k, f(v))
	}
	return g
}

// ToGtTag eg: GF[string] to GF[GtTag]
func (g GF[T]) ToGtTag() GF[GtTags] {
	var (
		n = NewGF[GtTags]()
	)
	for k, v := range g {
		n.Set(k, GtTags{origin: fmt.Sprint(v)})
	}
	return n
}
