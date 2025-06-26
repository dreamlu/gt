// package gt

package conf

import (
	"fmt"
	"github.com/dreamlu/gt/src/cons"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/util"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
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
	raw := c.Get(key)
	if raw != nil {
		c.yamlUnmarshal(raw, v)
	}
}

func (c *Yaml) yamlUnmarshal(viper any, v any) {
	out := mr.TrueValueOf(v)
	switch data := viper.(type) {
	case map[string]any:
		c.yamlViper(data, out.Type(), out)
	case []any:
		sliceElemType := out.Type().Elem()
		slice := reflect.MakeSlice(out.Type(), len(data), len(data))

		for i, item := range data {
			var target reflect.Value
			if sliceElemType.Kind() == reflect.Ptr {
				target = reflect.New(sliceElemType.Elem()) // *T
				c.yamlViper(item.(map[string]any), sliceElemType.Elem(), target.Elem())
				slice.Index(i).Set(target)
			} else {
				target = reflect.New(sliceElemType).Elem() // T
				c.yamlViper(item.(map[string]any), sliceElemType, target)
				slice.Index(i).Set(target)
			}
		}
		out.Set(slice)
	default:
		panic("not supported yaml struct type")
	}
}

func (c *Yaml) yamlViper(viper map[string]any, typ reflect.Type, val reflect.Value) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		if field.Anonymous {
			c.yamlViper(viper, field.Type, fieldVal)
			continue
		}

		tag := field.Tag.Get(cons.Yaml)
		var key string
		if tag != "" && tag != "-" {
			key = tag
		} else {
			key = util.HumpToLine(field.Name)
		}

		rawVal, ok := viper[strings.ToLower(key)]
		if !ok {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			if m, ok := rawVal.(map[string]any); ok {
				c.yamlViper(m, field.Type, fieldVal)
			}
		case reflect.Ptr:
			if m, ok := rawVal.(map[string]any); ok {
				ptr := reflect.New(field.Type.Elem())
				c.yamlViper(m, field.Type.Elem(), ptr.Elem())
				fieldVal.Set(ptr)
			}
		case reflect.Slice:
			if arr, ok := rawVal.([]any); ok {
				elemType := field.Type.Elem()
				slice := reflect.MakeSlice(field.Type, len(arr), len(arr))
				for j, item := range arr {
					elem := reflect.New(elemType).Elem()
					if m, ok := item.(map[string]any); ok {
						c.yamlViper(m, elemType, elem)
					} else {
						elem.Set(reflect.ValueOf(item).Convert(elemType))
					}
					slice.Index(j).Set(elem)
				}
				fieldVal.Set(slice)
			}
		default:
			fieldVal.Set(reflect.ValueOf(rawVal).Convert(field.Type))
		}
	}
}
