// Package util 提供 AES-GCM 对称加解密能力。
//
// 使用方式：
//
//  1. 在 main 的 init() 中注册密钥（32字节，AES-256）：
//     crypto.RegisterEncryptKey([]byte("your-32-byte-key-here-----------"))
//
//  2. 加密：
//     enc, err := crypto.Encrypt(plaintext)
//
//  3. 解密（未注册密钥时原样返回，兼容明文文件）：
//     plain, err := crypto.MaybeDecrypt(data)
package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"sync"
)

var (
	mu         sync.RWMutex
	encryptKey []byte
)

// RegisterEncryptKey 注册全局加密密钥，必须是 16/24/32 字节（对应 AES-128/192/256）。
// 全局只需调用一次，通常在 main 的 init() 最前面执行。
func RegisterEncryptKey(key []byte) {
	mu.Lock()
	defer mu.Unlock()
	cp := make([]byte, len(key))
	copy(cp, key)
	encryptKey = cp
}

// HasKey 判断是否已注册密钥。
func HasKey() bool {
	mu.RLock()
	defer mu.RUnlock()
	return len(encryptKey) > 0
}

// Encrypt 使用已注册的密钥对数据进行 AES-GCM 加密。
// 返回格式：nonce(12B) + ciphertext + tag(16B)。
func Encrypt(plaintext []byte) ([]byte, error) {
	mu.RLock()
	key := encryptKey
	mu.RUnlock()

	if len(key) == 0 {
		return nil, errors.New("crypto: encrypt key not registered")
	}
	return EncryptWithKey(key, plaintext)
}

// MaybeDecrypt 若已注册密钥则解密，否则原样返回（兼容明文文件读取场景）。
func MaybeDecrypt(data []byte) ([]byte, error) {
	mu.RLock()
	key := encryptKey
	mu.RUnlock()

	if len(key) == 0 {
		return data, nil
	}
	return DecryptWithKey(key, data)
}

// EncryptWithKey 使用指定 key 加密，不依赖全局状态。
func EncryptWithKey(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// DecryptWithKey 使用指定 key 解密，不依赖全局状态。
func DecryptWithKey(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("crypto: encrypted data too short")
	}
	return gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
}
