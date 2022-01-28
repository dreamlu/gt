package conf

import "sync"

type Config struct {
	// different devMode yaml data
	YamlS map[string]*Yaml
	// yaml project path
	path string
}

var (
	onceConfig sync.Once
	// global log
	config *Config
)

// Configger single config
func Configger(params ...string) *Config {

	onceConfig.Do(func() {
		config = NewConfig(params[:]...)
	})
	return config
}

func GetString(name string) string {
	return Configger().GetString(name)
}

func GetInt(name string) int {
	return Configger().GetInt(name)
}

func GetBool(name string) bool {
	return Configger().GetBool(name)
}

func GetStruct(name string, v interface{}) {
	Configger().GetStruct(name, v)
}

func Unmarshal(v interface{}) {
	Configger().Unmarshal(v)
}
