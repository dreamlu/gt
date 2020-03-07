package util

import (
	"log"
	"testing"
)

func TestAesEn(t *testing.T) {

	log.Println("[加密测试]:", AesEn("admin"))
	log.Println("[解密测试]:", AesDe("sPa0sTmDf6gasS9tHvIqKw=="))
}
