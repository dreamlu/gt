// package gt

package gt

import (
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	var config = Configger()
	log.Println("config read test: ", config.GetString("app.port"))
}

// can not read privilege field
type dbas struct {
	MaxIdleConn int
	MaxOpenConn int
	User        string
	Password    string
	host        string
	name        string
}

func TestConfig_GetStruct(t *testing.T) {
	dba := &dbas{}
	Configger().GetStruct("app.db", dba)
	t.Log(dba)
}
