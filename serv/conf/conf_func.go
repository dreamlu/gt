// package gt

package conf

import (
	"errors"
	"fmt"
	cons2 "github.com/dreamlu/gt/tool/cons"
	"github.com/imdario/mergo"
	"strings"
)

func DevMode(field string) {
	cons2.DefaultDevMode = field
}

// NewConfig new Config
// load all devMode yaml data
func NewConfig(params ...string) *Config {

	// default param
	path := cons2.ConfPath
	if len(params) > 0 {
		path = params[0]
	}
	path = newConf(path)
	c := &Config{
		YamlS: make(map[string]*Yaml, 2),
		path:  path,
	}
	// devMode
	devMode := c.getDevMode()

	// try
	devModePath := fmt.Sprintf("-%s", devMode)
	ss := strings.Split(path, ".")
	if len(ss) > 1 {
		devModePath += "."
	}
	c.path = strings.Join(ss, devModePath)
	// load data
	yaml := c.loadYaml()

	// add yamlS data
	c.YamlS[devMode] = yaml
	return c
}

// dir: default is conf/
// change dir to abs /xxx/conf/
// new abs conf dir
func newConf(dir string) string {
	return ProjectPath() + dir
}

// find yaml dev mode
// default devMode is app.yaml
// use 'app' as the map key
func (c *Config) getDevMode() (devMode string) {
	yaml := c.loadYaml()

	// add yamlS data
	c.YamlS["app"] = yaml

	if yaml.data == nil {
		panic(errors.New("no yaml: " + c.path))
	}
	return yaml.Get(cons2.DefaultDevMode).(string)
}

// load dev mode data
func (c *Config) loadYaml() *Yaml {
	yaml := &Yaml{}
	err := yaml.loadYaml(c.path)
	if err != nil {
		panic(errors.New("no yaml: " + c.path))
	}
	return yaml
}

// Get yaml data
func (c *Config) Get(name string) interface{} {
	for _, v := range c.YamlS {
		if value := v.Get(name); value != nil {
			return value
		}
	}
	return nil
}

func (c *Config) GetString(name string) string {
	for _, v := range c.YamlS {
		if value := v.GetString(name); value != "" {
			return value
		}
	}
	return ""
}

func (c *Config) GetInt(name string) int {
	for _, v := range c.YamlS {
		if value := v.GetInt(name); value != 0 {
			return value
		}
	}
	return 0
}

func (c *Config) GetBool(name string) bool {
	for _, v := range c.YamlS {
		if value := v.GetBool(name); value != false {
			return value
		}
	}
	return false
}

// GetStruct yaml to struct
// only support Accessible Field
func (c *Config) GetStruct(name string, s interface{}) {
	for _, v := range c.YamlS {
		v.Unmarshal(v.Get(name), s)
	}
}

func (c *Config) Unmarshal(s interface{}) {
	var t interface{}
	for _, v := range c.YamlS {
		v.Unmarshal(v.data, s)
		if t != nil {
			err := mergo.Merge(t, s, mergo.WithOverride)
			if err != nil {
				panic(err)
			}
		}
		t = s
	}
	s = t
}
