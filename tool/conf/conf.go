package conf

import "sync"

// config
type Config struct {
	// different devMode yaml data
	YamlS map[string]*Yaml
	// yaml dir
	dir string
}

var (
	onceConfig sync.Once
	// global log
	config *Config
)

// single config
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

// yaml to struct
// only support Accessible Field
func GetStruct(name string, s interface{}) {
	Configger().GetStruct(name, s)
}
