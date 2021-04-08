package bmap

import (
	"encoding/json"
	"fmt"
	"github.com/dreamlu/gt/tool/util/tag"
	"reflect"
)

type BMap map[string]interface{}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v BMap) Get(key string) interface{} {
	if v == nil {
		return ""
	}
	return v[key]
}

// return Get value and Del key
func (v BMap) Pop(key string) interface{} {
	s := v.Get(key)
	v.Del(key)
	return s
}

// Set sets the key to value. It replaces any existing
// values.
func (v BMap) Set(key string, value interface{}) BMap {
	v[key] = value
	return v
}

// Del deletes the values associated with key.
func (v BMap) Del(key string) BMap {
	delete(v, key)
	return v
}

// BMap to struct data
// value like
// type Te struct {
//		Name string `json:"name"` // must string type
//		ID   string `json:"id"` // must string type
//	}
func (v BMap) Struct(value interface{}) error {
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

// struct to BMap, maybe use Encode
// ps: the BMap key->value , value will be string type
func StructToBMap(v interface{}) (values BMap) {
	values = NewBMap()
	el := reflect.ValueOf(v)
	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}
	iVal := el
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		fi := typ.Field(i)
		values.Set(tag.GetFieldTag(fi), fmt.Sprint(iVal.Field(i)))
	}
	return
}

// new BMap
func NewBMap() BMap {
	return BMap{}
}
