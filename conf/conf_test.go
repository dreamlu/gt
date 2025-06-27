// package gt

package conf

import (
	"os"
	"testing"
	"time"
)

// can not read privilege field
type dba struct {
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Name        string `yaml:"name"`
	Log         bool   `yaml:"log"`
}

type conf struct {
	DevMode    string
	Port       string
	RemotePort string   `yaml:"remote-port"`
	Addr       []string `yaml:"addr"`
	Db         dba      `yaml:"db"`
}

func TestConfig_GetStruct(t *testing.T) {
	dba := &dba{}
	UnmarshalField("app.db", dba)
	t.Log(dba)
}

func TestConfigCus(t *testing.T) {
	DevMode("devMode")
	cf := NewConfig("conf/main.yml")
	//dba := &dbas{}
	cof := &conf{}
	//cf.UnmarshalField("db", dba)
	//t.Log(dba)
	cf.Unmarshal(cof)
	t.Log(cof)
	t.Log(cf.Get("db.name"))
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

func TestRemoteConfig(t *testing.T) {

	remote := &Remote{
		Provider: "consul",
		Endpoint: "http://192.168.10.168:8500",
		Path:     []string{"config/common.yaml"},
	}
	cfg := EmptyConfigger()
	cfg.AddRemoteConfig(remote)
	t.Log(cfg.Get("port"))
	t.Log(cfg.Get("test.remote-port"))

	cof := &conf{}
	DevMode("devMode")
	OverrideRemote(false)
	//OverrideRemote(false)
	cf := NewConfig("conf/main.yml")
	cf.AddRemoteConfig(remote)
	cf.WatchRemoteConfig()
	cf.Unmarshal(cof)
	t.Log(cof)
	for {
		t.Log(cf.Get("remote-port"))
		time.Sleep(3 * time.Second)
	}
}

func TestArray(t *testing.T) {
	type A struct {
		B int    `yaml:"b"`
		C string `yaml:"c"`
		D []struct {
			B int    `yaml:"b"`
			C string `yaml:"c"`
		}
		E []string
		F [][]string
	}
	var as []A
	UnmarshalField("app.a", &as)
	t.Log(as)

	a := GetSlice[string]("app.a.e")
	t.Log(a)
}
