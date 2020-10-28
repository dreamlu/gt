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
