package util

import "github.com/dreamlu/gt/src/reflect"

// RemoveDuplicate remove duplicated slice or pointer slice
func RemoveDuplicate[T comparable](s []T) []T {
	var (
		result []T
		temp   = map[any]struct{}{}
	)
	for _, item := range s {
		key := reflect.TrueValueOf(item).Interface()
		if _, ok := temp[key]; !ok {
			temp[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Remove[T comparable](s []T, sep ...T) (res []T) {
	for _, v := range s {
		if !Equals(sep, v) {
			res = append(res, v)
		}
	}
	return
}

func RemoveIndex[T comparable](s []T, i *int) {
	s = append(s[:*i], s[*i+1:]...)
	*i--
}
