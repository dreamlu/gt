package conf

import "sync"

type Config struct {
	// different devMode yaml data
	YamlS []*Yaml
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

type value interface {
	int | string | bool | float64
}

func Get[T value](name string) (t T) {
	if v := Configger().Get(name); v != nil {
		return v.(T)
	}
	return
}

func UnmarshalField(name string, v any) {
	Configger().UnmarshalField(name, v)
}

func Unmarshal(v any) {
	Configger().Unmarshal(v)
}
