// package gt

package gt

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/tool/file/file_func"
	"github.com/dreamlu/gt/tool/util/cons"
	"log"
	"os"
	"runtime"
	"strings"
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
func Configger(params ...string) *Config {

	onceConfig.Do(func() {
		config = NewConfig(params[:]...)
	})
	return config
}

// new Config
// load all devMode yaml data
func NewConfig(params ...string) *Config {

	// default param
	confDir := cons.ConfDir
	if len(params) > 0 {
		confDir = params[0]
	}
	confDir = newConf(confDir)
	//confDir = copyConf(confDir)
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
	fileDir, _ := os.Getwd()
	GOPATH := os.Getenv("GOPATH")
	ss := strings.Split(fileDir, GOPATH)
	if len(ss) > 1 {

		// default linux/mac os
		spChar := "/"
		if runtime.GOOS == "windows" {
			spChar = "\\"
		}
		ss2 := strings.Split(ss[1], spChar)

		newDir := GOPATH
		for i := 1; i < len(ss2); i++ {
			newDir += spChar + ss2[i]
			tmpDir := newDir + spChar + dir
			if file_func.Exists(tmpDir) {
				return tmpDir
			}
		}
	}

	// 返回默认
	return dir
}

// copy file to home .gt dir
//func copyConf(dir string) string {
//
//	gtDir := os.Getenv("HOME") + "/.gt/"
//
//	confDir := gtDir + time2.Now().Format("2006-01-02-15-04-05") + "/conf/"
//	err := os.MkdirAll(confDir, os.ModePerm)
//	println(err)
//	file_func.CopyDir(dir, confDir)
//	return confDir
//}

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
