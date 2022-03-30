package mock

import (
	"github.com/bxcodec/faker/v3"
	"log"
	"reflect"
)

var (
	//Sets the random size for slices and maps.
	randomSize = 21
)

// Mock mock data
func Mock(data any) {
	_ = faker.SetRandomMapAndSliceSize(randomSize)
	faker.SetStringLang(faker.LangCHI)
	//CustomGenerator()
	err := faker.FakeData(data)
	if err != nil {
		log.Println("[mock data Error]:", err)
	}
}

func AddProvider(key string, f func(v reflect.Value) (interface{}, error)) {
	_ = faker.AddProvider(key, f)
}

func SetRandomSize(size int) {
	randomSize = size
}

// no effect, there is a bug for faker
//// CustomGenerator ...
//func CustomGenerator() {
//	faker.AddProvider("cjon", func() faker.TaggedFunction {
//		return func(v reflect.Value) (any, error) {
//			return "danger-ranger", nil
//		}
//	}())
//}
