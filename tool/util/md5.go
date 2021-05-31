package util

import (
	"crypto/md5"
	"fmt"
)

// Md5
func Md5(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}
