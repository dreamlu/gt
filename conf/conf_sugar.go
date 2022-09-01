package conf

import "sync"

// 设计模式--单例模式[懒汉式]
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
	int64 | int | string | bool | float64
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
