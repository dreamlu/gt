package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Md5 md5
func Md5(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

// Sha256 sha256
func Sha256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
