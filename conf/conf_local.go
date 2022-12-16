package conf

import (
	"errors"
	"github.com/dreamlu/gt/src/cons"
)

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
