// package gt

package conf

import (
	"github.com/dreamlu/gt/src/cons"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/util"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
)

type Yaml struct {
	*viper.Viper
	isRemote bool
}

// load the default app.yaml data
func (c *Yaml) loadYaml(path string) {

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(cons.Yaml)
	err := v.ReadInConfig()
	if err != nil {
		panic("can not read " + path + " config")
	}
	c.Viper = v
}

// load the default app.yaml data
func (c *Yaml) loadRemoteYaml(provider, endpoint, path string) {

	v := viper.New()
	err := v.AddRemoteProvider(provider, endpoint, path)
	if err != nil {
		panic(err)
	}
	v.SetConfigType(cons.Yaml)
	err = v.ReadRemoteConfig()
	if err != nil {
		panic("can not read " + path + " config")
	}
	c.Viper = v
	c.isRemote = true
}

func (c *Yaml) Get(name string) any {
	return c.Viper.Get(name)
}

func (c *Yaml) Unmarshal(v any) {
	mp := c.AllSettings()
	c.yamlUnmarshal(mp, v)
}

func (c *Yaml) UnmarshalKey(key string, v any) {
	if mp := c.Get(key); mp != nil {
		c.yamlUnmarshal(mp.(map[string]any), v)
	}
}

func (c *Yaml) yamlUnmarshal(viper map[string]any, v any) {
	var (
		typ = mr.TrueTypeof(v)
		val = mr.TrueValueOf(v)
	)
	c.yamlViper(viper, typ, val)
	bs, err := yaml.Marshal(viper)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bs, v)
	if err != nil {
		panic(err)
	}
}

func (c *Yaml) yamlViper(viper map[string]any, typ reflect.Type, v reflect.Value) {
	var (
		tn = typ.NumField()
	)
	for i := 0; i < tn; i++ {
		field := typ.Field(i)
		val := v.Field(i)
		if field.Anonymous {
			c.yamlViper(viper, field.Type, val)
			continue
		}
		tv := field.Tag.Get(cons.Yaml)
		key := tv
		if tv == "" {
			tv = strings.ToLower(field.Name)
			key = util.HumpToLine(field.Name)
		}
		tv = strings.ToLower(tv)
		if vs := viper[tv]; vs != nil {
			switch field.Type.Kind() {
			case reflect.Struct:
				c.yamlViper(vs.(map[string]any), field.Type, val)
			default:
				delete(viper, tv)
				viper[key] = vs
			}
		}
	}
}
