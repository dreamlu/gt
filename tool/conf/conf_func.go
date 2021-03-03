// package gt

package conf

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/tool/file/file_func"
	"github.com/dreamlu/gt/tool/util/cons"
	"log"
)

// new Config
// load all devMode yaml data
func NewConfig(params ...string) *Config {

	// default param
	confDir := cons.ConfDir
	if len(params) > 0 {
		confDir = params[0]
	}
	confDir = newConf(confDir)
	c := &Config{
		YamlS: make(map[string]*Yaml, 2),
		dir:   confDir,
	}
	// devMode
	devMode := c.getDevMode()
	// load data
	yaml := c.loadYaml(fmt.Sprintf("%sapp-%s.yaml", c.dir, devMode))

	// add yamlS data
	c.YamlS[devMode] = yaml
	return c
}

// dir: default is conf/
// change dir to abs /xxx/conf/
// new abs conf dir
func newConf(dir string) string {
	return file_func.ProjectPath() + dir
}

// find yaml dev mode
// default devMode is app.yaml
// use 'app' as the map key
func (c *Config) getDevMode() (devMode string) {
	yaml := c.loadYaml(c.dir + "app.yaml")

	// add yamlS data
	c.YamlS["app"] = yaml

	if yaml.data == nil {
		panic(errors.New("[gt] no app.yaml, path:" + c.dir))
	}

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
