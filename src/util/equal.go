package util

import "encoding/json"

func Equal(src, dst any) bool {
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
