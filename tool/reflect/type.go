package reflect

import (
	"fmt"
	"reflect"
)

type Kind interface {
	Kind() reflect.Kind
}

func TrueTypeof(v interface{}) reflect.Type {
	return TrueType(reflect.TypeOf(v))
}

func TrueType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func TrueValueOf(v interface{}) reflect.Value {
	return TrueValue(reflect.ValueOf(v))
}

func TrueValue(typ reflect.Value) reflect.Value {
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

func IsStruct(typ Kind) bool {
	if typ.Kind() == reflect.Struct {
		return true
	}
	return false
}
