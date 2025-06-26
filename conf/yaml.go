// package gt

package conf

import (
	"fmt"
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
	isRemote                 bool
	provider, endpoint, path string
}

// load the default app.yaml data
func (c *Yaml) loadYaml(path string) {

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(cons.Yaml)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("can not read %s config, error: %v", path, err))
	}
	c.Viper = v
	c.path = path
}

// load the default app.yaml data
func (c *Yaml) loadRemoteYaml(provider, endpoint, path string) {

	c.Viper = viper.New()
	c.Viper.SetConfigType(cons.Yaml)
	err := c.AddRemoteProvider(provider, endpoint, path)
	if err != nil {
		panic(err)
	}
	err = c.Viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	c.isRemote = true
	c.provider = provider
	c.endpoint = endpoint
	c.path = path
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
		c.yamlUnmarshal(mp, v)
	}
}

func (c *Yaml) yamlUnmarshal(raw any, v any) {
	typ := mr.TrueTypeof(v)
	val := mr.TrueValueOf(v)

	switch val.Kind() {
	case reflect.Struct:
		// 单个结构体
		if m, ok := raw.(map[string]any); ok {
			c.yamlViper(m, typ, val)
		}
	case reflect.Slice:
		// 结构体切片或指针切片
		if arr, ok := raw.([]any); ok {
			elemType := typ.Elem()
			slice := reflect.MakeSlice(typ, len(arr), len(arr))
			for i, item := range arr {
				if itemMap, ok := item.(map[string]any); ok {
					if elemType.Kind() == reflect.Ptr {
						ptr := reflect.New(elemType.Elem())
						c.yamlViper(itemMap, elemType.Elem(), ptr.Elem())
						slice.Index(i).Set(ptr)
					} else {
						elem := reflect.New(elemType).Elem()
						c.yamlViper(itemMap, elemType, elem)
						slice.Index(i).Set(elem)
					}
				}
			}
			val.Set(slice)
		}
	}

	// 最终用 yaml 来赋值（支持 yaml tag）
	bs, err := yaml.Marshal(raw)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bs, v); err != nil {
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
			switch data := vs.(type) {
			case map[string]any:
				c.yamlViper(data, field.Type, val)

			case []map[string]any:
				for _, d := range data {
					c.yamlViper(d, field.Type, val)
				}
			default:
				delete(viper, tv)
				viper[key] = vs
			}

		}
	}
}
