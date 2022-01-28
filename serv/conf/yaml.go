// package gt

package conf

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Yaml struct {
	// yaml data
	data map[interface{}]interface{}
}

// load the default app.yaml data
func (c *Yaml) loadYaml(path string) error {

	yamlS, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return readErr
	}
	// yaml解析的时候c.data如果没有被初始化，会自动为你做初始化
	err := yaml.Unmarshal(yamlS, &c.data)
	if err != nil {
		return errors.New("can not parse " + path + " config")
	}
	return nil
}

func (c *Yaml) Get(name string) interface{} {
	path := strings.Split(name, ".")
	data := c.data
	for key, value := range path {
		v, ok := data[value]
		if !ok {
			break
		}
		if (key + 1) == len(path) {
			return v
		}
		// print yaml v3 problem!
		// use go-yaml v2 replace
		//log.Println(name, "&&",reflect.TypeOf(v).String())
		if reflect.TypeOf(v).String() == "map[interface {}]interface {}" {
			data = v.(map[interface{}]interface{})
		}
	}
	return nil
}

func (c *Yaml) GetString(name string) string {
	value := c.Get(name)
	switch value := value.(type) {
	case string:
		return value
	case bool, float64, int:
		return fmt.Sprint(value)
	default:
		return ""
	}
}

func (c *Yaml) GetInt(name string) int {
	value := c.Get(name)
	switch value := value.(type) {
	case string:
		i, err := strconv.Atoi(value)
		log.Println("[YAML type error]: ", err)
		return i
	case int:
		return value
	case bool:
		if value {
			return 1
		}
		return 0
	case float64:
		return int(value)
	default:
		return 0
	}
}

func (c *Yaml) GetBool(name string) bool {
	value := c.Get(name)
	switch value := value.(type) {
	case string:
		str, err := strconv.ParseBool(value)
		log.Println("[YAML type error]: ", err)
		return str
	case int:
		if value != 0 {
			return true
		}
		return false
	case bool:
		return value
	case float64:
		if value != 0.0 {
			return true
		}
		return false
	default:
		return false
	}
}

func (c *Yaml) Unmarshal(data interface{}, s interface{}) {

	b, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(b, s)
	if err != nil {
		panic(err)
	}
}
