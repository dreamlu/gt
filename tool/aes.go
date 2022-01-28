// package gt

package tool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// defaultKey can be assigned
const defaultKey = "github.com/dreamlu/gt dreamlu123"

type Aes struct {
	key string
}

func NewAes(key ...string) *Aes {
	var as Aes
	as.key = defaultKey
	if len(key) != 0 {
		as.key = key[0]
	}
	return &as
}

func (as *Aes) IsAes(data string) bool {
	defer func() {
		recover()
	}()
	as.AesDe(data)
	return true
}

func (as *Aes) AesEn(data string) string {
	origData := []byte(data)
	k := []byte(as.key)

	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)
}

func (as *Aes) AesDe(data string) string {
	origData, _ := base64.StdEncoding.DecodeString(data)
	k := []byte(as.key)

	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	orig := make([]byte, len(origData))
	blockMode.CryptBlocks(orig, origData)
	orig = pKCS7UnPadding(orig)
	return string(orig)
}

func pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
