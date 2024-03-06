package reflect

import (
	"reflect"
)

// Set data field value
// field type must the same as value
// can set basic or ptr value
func Set(data any, field string, value any) {
	var (
		v  = TrueValueOf(data)
		fd = v.FieldByName(field)
	)
	if IsPtr(fd) { // include ptr value or nil
		fd.Set(reflect.New(fd.Type().Elem()))
		fd.Elem().Set(TrueValueOf(value)) // like *a = b
		// or use following code
		//vv := reflect.ValueOf(value)
		//if !IsPtr(vv) {
		//	// perfect! nice!
		//	// like: *a = b
		//	fd.Elem().Set(reflect.ValueOf(value))
		//	return
		//}
		//fd.Set(vv)
		return
	}
	v.FieldByName(field).Set(TrueValueOf(value))
}

// SetByIndex data field index value
// type and field must same
func SetByIndex(data any, index int, value any) {
	TrueValueOf(data).Field(index).Set(TrueValueOf(value))
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
func ToSlice(slice any) []any {
	v := TrueValueOf(slice)
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
