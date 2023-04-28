package util

import "github.com/dreamlu/gt/src/reflect"

func Contains[T comparable](s []T, sep ...T) bool {
	for _, v := range s {
		for _, se := range sep {
			if Equal(v, se) {
				return true
			}
		}
	}
	return false
}

func ContainField[T comparable](dst []T, src any, fieldName string) bool {
	for _, d := range dst {
		if reflect.Field(d, fieldName) == reflect.Field(src, fieldName) {
			return true
		}
	}
	return false
}
