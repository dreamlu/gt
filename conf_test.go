// package der

package der

import (
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	var config = Configger()
	log.Println("config read test: ", config.GetString("app.port"))
}
