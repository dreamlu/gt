package reflect

import (
	"encoding/json"
	"errors"
	"reflect"
)

// reflect value via field name
func GetDataByFieldName(data interface{}, filedName string) (interface{}, error) {
	typ := reflect.TypeOf(data)

	switch typ.Kind() {
	case reflect.Ptr, reflect.Chan, reflect.Map, reflect.Array, reflect.Slice:
		v := reflect.ValueOf(data).Elem()
		f := v.FieldByName(filedName)
		return f.Interface(), nil
	case reflect.Struct:
		v := reflect.ValueOf(data)
		f := v.FieldByName(filedName)
		return f.Interface(), nil
	}
	return nil, errors.New(filedName + "not exit")
}

// arr must array data
// array struct data to []interface
func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// return struct string name
func StructToString(st interface{}) string {
	typ := reflect.TypeOf(st)
	return typ.Name()
}

// struct to map
func ToMap(st interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	bs, _ := json.Marshal(st)
	_ = json.Unmarshal(bs, &m)
	return m
}
