package util

import (
	"testing"
)

func TestAesEn(t *testing.T) {

	t.Log("[aesEn]:", AesEn("admin"))
	t.Log("[aesDe]:", AesDe("sPa0sTmDf6gasS9tHvIqKw=="))
	t.Log(IsAes("13242trergf"))
	t.Log(IsAes("sPa0sTmDf6gasS9tHvIqKw=="))
}
