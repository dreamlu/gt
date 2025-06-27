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
		config = NewConfig(params...)
	})
	return config
}

func EmptyConfigger() *Config {
	onceConfig.Do(func() {
		config = NewConfig("")
	})
	return config
}

type value interface {
	int8 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | string | bool | float32 | float64
}

func Get[T value](name string) (t T) {
	if v := Configger().Get(name); v != nil {
		return v.(T)
	}
	return
}

func GetSlice[T value](name string) (t []T) {
	if v := Configger().Get(name); v != nil {
		slice := v.([]any)
		for _, s := range slice {
			t = append(t, s.(T))
		}
		return t
	}
	return
}

func UnmarshalField(name string, v any) {
	Configger().UnmarshalField(name, v)
}

func Unmarshal(v any) {
	Configger().Unmarshal(v)
}
