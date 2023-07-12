package reflect

import (
	"errors"
	"reflect"
)

func IsImplements(data any, i any) bool {
	return reflect.TypeOf(data).Implements(TrueTypeof(i))
}

func Call(data any, method string, args ...any) error {
	value := reflect.ValueOf(data)
	f := value.MethodByName(method)
	if !f.IsValid() {
		return errors.New("method not found")
	}
	var params []reflect.Value
	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}
	result := f.Call(params)
	if len(result) > 0 {
		ri := result[0].Interface()
		if ri != nil {
			return ri.(error)
		}
	}
	return nil
}
