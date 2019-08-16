// package der

package der

import (
	"log"
	"testing"
)

var config = &Config{}

func init() {
	config.NewConfig()
}

func TestConfig(t *testing.T) {
	log.Println("config read test: ", config.GetString("app.port"))
}
