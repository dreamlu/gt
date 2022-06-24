// package gt

package conf

import (
	"os"
	"testing"
)

// can not read privilege field
type dbas struct {
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	host        string `yaml:"host"`
	name        string `yaml:"name"`
}

type conf struct {
	DevMode string `yaml:"devMode"`
	Port    string `yaml:"port"`
}

func TestConfig_GetStruct(t *testing.T) {
	dba := &dbas{}
	UnmarshalField("app.db", dba)
	t.Log(dba)
}

func TestConfigCus(t *testing.T) {
	dba := &dbas{}
	cof := &conf{}
	DevMode("devMode")
	cf := NewConfig("conf/main.yml")
	cf.UnmarshalField("db", dba)
	t.Log(dba)
	cf.Unmarshal(cof)
	t.Log(cof)
	t.Log(cf.Get("db.user"))
}

func TestConfigger(t *testing.T) {

	dir, _ := os.Getwd()
	t.Log(dir)
	t.Log(os.Getenv("GOPATH"))
	mode := Get[string]("app.devMode")
	t.Log(mode)
	port := Get[int]("app.port")
	t.Log(port)
}
