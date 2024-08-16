package reflect

import "reflect"

// New returns a Value representing a pointer to a new zero value
func New(v any) any {
	return reflect.New(TrueTypeof(v)).Interface()
}

// NewArray returns a []Value representing a pointer to a new zero value
func NewArray(v any) any {
	//log.Println(reflect.MakeSlice(reflect.SliceOf(t), 0, 0))
	return reflect.New(reflect.SliceOf(TrueTypeof(v))).Interface()
}

func IsZero(v any) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func IsNil(v any) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).IsNil()
}
