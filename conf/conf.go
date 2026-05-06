// Package conf 配置加载包，支持文件系统和 embed 两种模式。
//
// # 默认模式（文件系统）
//
// 不做任何额外操作，conf 会自动从文件系统查找并加载 yaml 配置文件，行为与原来完全一致。
//
//	conf.NewConfig() // 自动查找 conf/app.yaml
//
// # Embed 模式（将配置文件打包进二进制）
//
// 适用于需要单文件部署、不依赖外部配置文件的场景。
//
// 使用步骤：
//
// 1. 在和配置文件同级目录新建一个 Go 文件，声明 embed.FS 并在 init() 中注册：
//
//	package conf
//
//	import (
//	    "embed"
//	    gtconf "github.com/dreamlu/gt/conf"
//	)
//
//	//go:embed app.yaml app-dev.yaml
//	var FS embed.FS
//
//	func init() {
//	    gtconf.RegisterEmbedFS(FS)
//	}
//
// 2. 在最早被初始化、且会触发 conf 读取的包里，空导入上面的包，确保 embed 注册先于配置读取执行：
//
//	import _ "solana/conf" // 保证 RegisterEmbedFS 先于 conf.Configger 执行
//
// 3. 其余业务代码无需任何改动，conf.NewConfig / conf.Get / conf.UnmarshalField 等接口行为不变。
//
// # 注意事项
//
//   - //go:embed 路径不能包含 '..'，配置文件必须在当前目录或子目录下。
//   - embed 文件声明的包与 gt/conf 包名可能冲突，import 时注意加别名。
//   - embed 模式下 RegisterEmbedFS 必须在第一次调用 conf.Get/Configger 之前完成，
//     否则会 panic。通过空导入可以保证这一点。
package conf

import (
	"embed"
	"fmt"
	"strings"

	"github.com/dreamlu/gt/src/cons"
)

// Global embedFS, nil means to use the file system (default behavior)
var globalEmbedFS *embed.FS

// RegisterEmbedFS registers embed.FS. After calling, NewConfig automatically reads from binary
// Just call it once at the front of main's init() or main()
func RegisterEmbedFS(efs embed.FS) {
	globalEmbedFS = &efs
}

type Config struct {
	// different devMode yaml data
	YamlS []*Yaml
	// yaml project path
	path string
}

func DevMode(field string) {
	cons.DefaultDevMode = field
}

func OverrideRemote(override bool) {
	cons.ConfOverride = override
}

// NewConfig new Config
// load all devMode yaml data
// NewConfig new Config
func NewConfig(params ...string) *Config {
	var (
		path = cons.ConfPath
		c    = &Config{}
	)
	if len(params) > 0 {
		path = params[0]
		if path == "" {
			return c
		}
	}

	// 统一走 newPath，embed 分支在内部处理
	path = newPath(path)
	c.path = path

	devMode := c.getDevMode()
	if devMode == "" {
		return c
	}

	devModePath := fmt.Sprintf("-%s", devMode)
	ss := strings.Split(path, ".")
	if len(ss) > 1 {
		devModePath += "."
	}
	c.path = strings.Join(ss, devModePath)

	// 统一走 existsPath，embed 分支在内部处理
	if !existsPath(c.path) {
		return c
	}

	c.loadYaml()
	return c
}

// Get name
func (c *Config) Get(name string) (value any) {
	for _, v := range c.YamlS {
		if t := v.Get(name); t != nil {
			value = t
		}
	}
	return value
}

// UnmarshalField yaml to struct
func (c *Config) UnmarshalField(name string, s any) {
	for _, v := range c.YamlS {
		v.UnmarshalKey(name, s)
	}
}

func (c *Config) Unmarshal(s any) {
	for _, v := range c.YamlS {
		v.Unmarshal(s)
	}
}
