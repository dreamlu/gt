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

	type Config struct {
		PrivateKey      string             `yaml:"private_key"`       // 加密存储的私钥字符串
		MinProfitUSD    float64            `yaml:"min_profit_usd"`    // 执行套利的最小利润阈值SOL
		MaxFee          float64            `yaml:"max_fee"`           // 优先级费用SOL
		MaxDailyTrades  int                `yaml:"max_daily_trades"`  // 每日允许的最大交易次数
		IntervalSeconds int                `yaml:"interval_seconds"`  // 套利机会检查的时间间隔
		MonitorPairs    [][]string         `yaml:"monitor_pairs"`     // 需要监控的交易对列表
		RPCTimeout      int                `yaml:"rpc_timeout"`       // Solana RPC请求的超时时间
		HTTPTimeout     int                `yaml:"http_timeout"`      // HTTP请求的超时时间
		SlippageBps     int                `yaml:"slippage_bps"`      // 交易滑点设置(基点)
		MaxTradeUSDMap  map[string]float64 `yaml:"max_trade_usd_map"` // usd -> 最大交易金额映射
		Tip             float64            `yaml:"tip"`               // 小费usd
		JupiterAddr     string             `yaml:"jupiter_addr"`      // jupiterAddr
		JitoAddr        string             `yaml:"jito_addr"`         // jitoAddr
		RpcAddr         string             `yaml:"rpc_addr"`          // Solana RPC地址
		WsAddr          string             `yaml:"ws_addr"`           // Solana WS地址
	}
	Cfg := &Config{}
	UnmarshalField("app.solana", Cfg)
	t.Log(Cfg)
}
