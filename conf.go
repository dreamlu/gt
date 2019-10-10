// package der

package der

import (
	"log"
	"sync"
)

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
func Configger() *Config {

	onceConfig.Do(func() {
		config = NewConfig()
	})
	return config
}

// new Config
// load all devMode yaml data
func NewConfig() *Config {

	c := &Config{}
	// init param
	c.dir = "conf/"
	c.YamlS = make(map[string]*Yaml, 2)

	// devMode
	devMode := c.getDevMode()
	// load data
	yaml := c.loadYaml(c.dir + "app-" + devMode + ".yaml")

	// add yamlS data
	c.YamlS[devMode] = yaml
	return c
}

// find yaml dev mode
// default devMode is app.yaml
// use 'app' as the map key
func (c *Config) getDevMode() (devMode string) {
	yaml := c.loadYaml(c.dir + "app.yaml")

	// add yamlS data
	c.YamlS["app"] = yaml

	return yaml.Get("app.devMode").(string)
}

// load dev mode data
func (c *Config) loadYaml(path string) *Yaml {
	yaml := &Yaml{}
	err := yaml.loadYaml(path)
	if err != nil {
		log.Println("[CONFIG ERROR]: ", err)
	}
	return yaml
}

// get yaml data
// find the first data, must different from app.yaml
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

// yaml to struct
// only support Accessible Field
func (c *Config) GetStruct(name string, s interface{}) {
	for _, v := range c.YamlS {
		v.GetStruct(name, s)
	}
}
