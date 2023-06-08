package reflect

import (
	"reflect"
)

func IsImplements(data any, i any) bool {
	return reflect.TypeOf(data).Implements(TrueTypeof(i))
}

func Call(data any, method string, args ...any) bool {
	value := reflect.ValueOf(data)
	f := value.MethodByName(method)
	if !f.IsValid() {
		return false
	}
	var params []reflect.Value
	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}
	f.Call(params)
	return true
}
