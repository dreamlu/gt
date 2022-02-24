package mock

import (
	"github.com/bxcodec/faker/v3"
	"log"
)

var (
	//Sets the random size for slices and maps.
	randomSize = 21
)

// Mock mock data
func Mock(data interface{}) {
	_ = faker.SetRandomMapAndSliceSize(randomSize)
	//CustomGenerator()
	err := faker.FakeData(data)
	if err != nil {
		log.Println("[mock data Error]:", err)
	}
}

func SetRandomSize(size int) {
	randomSize = size
}

// no effect, there is a bug for faker
//// CustomGenerator ...
//func CustomGenerator() {
//	faker.AddProvider("cjon", func() faker.TaggedFunction {
//		return func(v reflect.Value) (interface{}, error) {
//			return "danger-ranger", nil
//		}
//	}())
//}
