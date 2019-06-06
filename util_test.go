package der

import (
	"log"
	"testing"
)

func TestAesEn(t *testing.T) {

	log.Println("[加密测试]:", AesEn("123456"))
	log.Println("[解密测试]:", AesDe("lIEbR7cEp2U10gtM0j8dCg=="))
}
