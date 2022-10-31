package reflect

import (
	"reflect"
)

// Set data field value
// type and field must same
func Set(data any, field string, value any) {
	TrueValueOf(data).FieldByName(field).Set(TrueValueOf(value))
}

// Field reflect value via field name
// field must exist
func Field(data any, field string) any {
	return TrueValueOf(data).FieldByName(field).Interface()
}

// TrueField return Field true value, eg: Ptr value
func TrueField(data any, field string) any {
	return TrueValueOf(Field(data, field)).Interface()
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
