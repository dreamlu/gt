// package gt

package conf

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/src/cons"
	"github.com/dreamlu/gt/src/gos"
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

func OverrideRemote(override bool) {
	cons.ConfOverride = override
}

// NewConfig new Config
// load all devMode yaml data
func NewConfig(params ...string) *Config {

	// default param
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
	if !gos.Exists(c.path) {
		return c
	}

	// load data
	c.loadYaml()
	return c
}

// find yaml dev mode
// default devMode is app.yaml
// use 'app' as the map key
func (c *Config) getDevMode() (devMode string) {
	yaml := c.loadYaml()
	if yaml.Viper == nil {
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
	yaml.loadYaml(c.path)
	// add yamlS data
	c.YamlS = append(c.YamlS, yaml)
	return yaml
}

func (c *Config) AddRemoteConfig(provider, endpoint, path string) *Yaml {
	yaml := &Yaml{}
	yaml.loadRemoteYaml(provider, endpoint, path)
	// add yamlS data
	if cons.ConfOverride {
		c.YamlS = append([]*Yaml{yaml}, c.YamlS...)
	} else {
		c.YamlS = append(c.YamlS, yaml)
	}
	return yaml
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
// only support Accessible Field
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
