package mock

import (
	"crypto/rand"
	"math/big"
)

const (
	Number  = "number"
	Char    = "char"
	Chinese = "zh"
	En      = "en"
	LowerEn = "en_lower"
	UpperEn = "en_upper"
)

func GetRand(lang string, length int) string {
	res := make([]rune, length)
	var model []rune
	switch lang {
	case Number:
		model = []rune("0123456789")
	case LowerEn:
		model = []rune("abcdefghijklmnopqrstuvwxyz")
	case UpperEn:
		model = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case Char:
		model = []rune("!@#~$%^&*()+|_")
	case En:
		model = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	case Chinese:
		model = []rune("的一国在人了有中是年和大业不为发会工经上地市要个产这出行作生家以成到日民来我部对进多全建他公开们场展时理新方主企资实学报制政济用同于法高长现本月定化加动合品重关机分力自外者区能设后就等体下万元社过前面")
	}
	for i := range res {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(model))))
		res[i] = model[int(index.Int64())]
	}
	return string(res)
}
