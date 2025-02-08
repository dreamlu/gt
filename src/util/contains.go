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

func ContainsField[T comparable](s []T, sep []T, field string) bool {
	for _, v := range s {
		for _, se := range sep {
			if reflect.Field(v, field) == reflect.Field(se, field) {
				return true
			}
		}
	}
	return false
}

func ContainsDiffField[T1 comparable, T2 comparable](s []T1, sep []T2, sField, sepField string) bool {
	for _, v := range s {
		for _, se := range sep {
			if reflect.Field(v, sField) == reflect.Field(se, sepField) {
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

func ContainDiffField[T comparable](dst []T, src any, dstField, srcField string) bool {
	for _, d := range dst {
		if reflect.Field(d, dstField) == reflect.Field(src, srcField) {
			return true
		}
	}
	return false
}
