package reflect

import "reflect"

// New returns a Value representing a pointer to a new zero value
func New(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return reflect.New(t).Interface()
}

// NewArray returns a []Value representing a pointer to a new zero value
func NewArray(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	//log.Println(reflect.MakeSlice(reflect.SliceOf(t), 0, 0))
	return reflect.New(reflect.SliceOf(t)).Interface()
}
