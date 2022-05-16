// package gt

package conf

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/tool/util/cons"
	"github.com/dreamlu/gt/tool/util/gos"
	"github.com/imdario/mergo"
	"strings"
)

var (
	defaultDevMode = "app.devMode"
)

func DevMode(field string) {
	defaultDevMode = field
}

// NewConfig new Config
// load all devMode yaml data
func NewConfig(params ...string) *Config {

	// default param
	path := cons.ConfPath
	if len(params) > 0 {
		path = params[0]
	}
	path = newConf(path)
	c := &Config{
		//YamlS: make([]*Yaml, 2),
		path: path,
	}

	devModePath := fmt.Sprintf("-%s", c.getDevMode())
	ss := strings.Split(path, ".")
	if len(ss) > 1 {
		devModePath += "."
	}
	c.path = strings.Join(ss, devModePath)
	if !gos.Exists(c.path) {
		return c
	}

	// load data
	yaml := c.loadYaml()

	// add yamlS data
	c.YamlS = append(c.YamlS, yaml)
	return c
}

// rPath: relative path
// aPath: absolute path
// return config path
func newConf(rPath string) string {
	aPath := ProjectPath() + rPath
	if gos.Exists(aPath) {
		return aPath
	}
	return rPath
}

// find yaml dev mode
// default devMode is app.yaml
// use 'app' as the map key
func (c *Config) getDevMode() (devMode string) {
	yaml := c.loadYaml()
	if yaml.data == nil {
		panic(errors.New("no yaml: " + c.path))
	}

	// add yamlS data
	c.YamlS = append(c.YamlS, yaml)

	if res := yaml.Get(defaultDevMode); res != nil {
		return res.(string)
	}
	return ""
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
// find the first data, must different from app.yaml
func (c *Config) Get(name string) (value interface{}) {
	for _, v := range c.YamlS {
		value = v.Get(name)
	}
	return value
}

func (c *Config) GetString(name string) (value string) {
	for _, v := range c.YamlS {
		value = v.GetString(name)
	}
	return value
}

func (c *Config) GetInt(name string) (value int) {
	for _, v := range c.YamlS {
		value = v.GetInt(name)
	}
	return value
}

func (c *Config) GetBool(name string) (value bool) {
	for _, v := range c.YamlS {
		value = v.GetBool(name)
	}
	return value
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
