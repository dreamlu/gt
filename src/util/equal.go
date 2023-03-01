package util

import (
	"encoding/json"
	"github.com/dreamlu/gt/src/reflect"
)

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
	if reflect.TrueValueOf(src).Interface() != reflect.TrueValueOf(dst).Interface() {
		return true
	}
	return false
}
