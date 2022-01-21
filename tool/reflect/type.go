package reflect

import (
	"fmt"
	"reflect"
)

func TrueTypeof(v interface{}) (typ reflect.Type) {
	typ = reflect.TypeOf(v)
	return TrueType(typ)
}

func TrueType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func TrueTypeofValue(v interface{}) (typ reflect.Type, i interface{}) {
	typ = reflect.TypeOf(v)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		v = reflect.ValueOf(v).Elem().Interface()
	}
	return typ, v
}

func Path(typ reflect.Type, path ...string) string {
	return fmt.Sprintf("%s%s_%s", typ.PkgPath(), typ.Name(), path)
}

func IsStruct(typ reflect.Type) bool {
	if typ.Kind() == reflect.Struct {
		return true
	}
	return false
}
