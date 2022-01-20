package mock

import (
	"github.com/bxcodec/faker/v3"
	"log"
)

// Mock mock data
func Mock(data interface{}) {
	//CustomGenerator()
	err := faker.FakeData(data)
	if err != nil {
		log.Println("[mock data Error]:", err)
	}
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
