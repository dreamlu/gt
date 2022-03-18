// package gt

package conf

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strings"
)

type Yaml struct {
	// yaml data
	data map[any]any
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

func (c *Yaml) Get(name string) any {
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
			data = v.(map[any]any)
		}
	}
	return nil
}

func (c *Yaml) Unmarshal(data any, s any) {

	b, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(b, s)
	if err != nil {
		panic(err)
	}
}
