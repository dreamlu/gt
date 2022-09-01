// package gt

package conf

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/src/gos"
	"github.com/imdario/mergo"
	"strings"
)

type Config struct {
	// different devMode yaml data
	YamlS []*Yaml
	// yaml project path
	path string
}

func DevMode(field string) {
	cons.DefaultDevMode = field
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
	if !gos.Exists(c.path) {
		return c
	}

	// load data
	c.loadYaml()
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
	if devModeI := yaml.Get(cons.DefaultDevMode); devModeI != nil {
		return devModeI.(string)
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
	// add yamlS data
	c.YamlS = append(c.YamlS, yaml)
	return yaml
}

// Get yaml data
// find the first data, must different from app.yaml
func (c *Config) Get(name string) (value any) {
	for _, v := range c.YamlS {
		if t := v.Get(name); t != nil {
			value = t
		}
	}
	return value
}

// UnmarshalField yaml to struct
// only support Accessible Field
func (c *Config) UnmarshalField(name string, s any) {
	for _, v := range c.YamlS {
		v.Unmarshal(v.Get(name), s)
	}
}

func (c *Config) Unmarshal(s any) {
	var t any
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
