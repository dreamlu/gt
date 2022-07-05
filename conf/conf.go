package conf

type Config struct {
	// different devMode yaml data
	YamlS []*Yaml
	// yaml project path
	path string
}

var (
	// global config
	config *Config
)

// 设计模式--单例模式[饿汉式]
func init() {
	config = NewConfig()
}

type value interface {
	int64 | int | string | bool | float64
}

func Get[T value](name string) (t T) {
	if v := config.Get(name); v != nil {
		return v.(T)
	}
	return
}

func UnmarshalField(name string, v any) {
	config.UnmarshalField(name, v)
}

func Unmarshal(v any) {
	config.Unmarshal(v)
}
