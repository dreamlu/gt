package reflect

import (
	"errors"
	"reflect"
)

// Set data field value
// type and field must same
func Set(data any, field string, value any) {
	TrueValueOf(data).FieldByName(field).Set(TrueValueOf(value))
}

// Field reflect value via field name
func Field(data any, field string) (any, error) {
	typ := TrueValueOf(data)

	switch typ.Kind() {
	case reflect.Ptr, reflect.Chan, reflect.Map, reflect.Array, reflect.Slice:
		v := reflect.ValueOf(data).Elem()
		f := v.FieldByName(field)
		return f.Interface(), nil
	case reflect.Struct:
		v := reflect.ValueOf(data)
		f := v.FieldByName(field)
		return f.Interface(), nil
	}
	return nil, errors.New(field + "not exit")
}

// ToSlice arr must array data
// array struct data to []interface
func ToSlice(arr any) []any {
	v := TrueValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil
	}
	l := v.Len()
	ret := make([]any, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// Name return struct string name
func Name(v any) string {
	return TrueTypeof(v).Name()
}
