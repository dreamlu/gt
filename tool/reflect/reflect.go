package reflect

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

// 典型reflect 操作
// 返回对应字段值
func GetDataByFieldName(data interface{}, filedName string) (interface{}, error) {
	typ := reflect.TypeOf(data)
	//log.Println(typ.String())

	// 指针类型判断
	switch typ.Kind() {
	case reflect.Ptr, reflect.Chan, reflect.Map, reflect.Array, reflect.Slice:
		v := reflect.ValueOf(data).Elem()
		f := v.FieldByName(filedName)
		//判断字段是否存在
		// 字段不催在判断错误
		//if !v.IsValid() {
		//	return nil, errors.New(filedName + "字段不存在")
		//}
		return f.Interface(), nil
	case reflect.Struct:
		// struct type like User
		v := reflect.ValueOf(data)
		f := v.FieldByName(filedName)
		//判断字段是否存在
		// 字段不催在判断错误
		return f.Interface(), nil
	}
	return nil, errors.New("id不存在")
}

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		log.Println(arr, "[ERROR]: arr not slice")
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
