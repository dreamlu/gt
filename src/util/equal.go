package util

import (
	"encoding/json"
	"fmt"
	"github.com/dreamlu/gt/src/reflect"
)

func Equals[T comparable](s []T, sep ...T) bool {
	for _, v := range s {
		for _, se := range sep {
			if Equal(v, se) {
				return true
			}
		}
	}
	return false
}

func EqualJson(src, dst any) bool {
	srcBs, err := json.Marshal(src)
	if err != nil {
		return false
	}
	dstBs, err := json.Marshal(dst)
	if err != nil {
		return false
	}
	if string(srcBs) == string(dstBs) {
		return true
	}
	return false
}

func Equal(src, dst any) bool {
	if fmt.Sprint(reflect.TrueValueOf(src).Interface()) == fmt.Sprint(reflect.TrueValueOf(dst).Interface()) {
		return true
	}
	return false
}
