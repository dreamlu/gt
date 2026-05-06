package conf

import (
	"errors"
	"io/fs"

	"github.com/dreamlu/gt/src/cons"
)

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

func (c *Config) loadYaml() *Yaml {
	yaml := &Yaml{}

	if globalEmbedFS != nil {
		f, err := globalEmbedFS.Open(c.path)
		if err != nil {
			panic(errors.New("embed fs: no yaml: " + c.path))
		}
		defer f.Close()

		if stat, err := fs.Stat(globalEmbedFS, c.path); err != nil || stat.IsDir() {
			panic(errors.New("embed fs: no yaml: " + c.path))
		}

		yaml.loadYamlFromReader(f, c.path)
	} else {
		yaml.loadYaml(c.path)
	}

	c.YamlS = append(c.YamlS, yaml)
	return yaml
}
