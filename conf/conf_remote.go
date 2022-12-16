package conf

import (
	"github.com/dreamlu/gt/src/cons"
	"time"
)

type Remote struct {
	Provider string   `yaml:"provider"`
	Endpoint string   `yaml:"endpoint"`
	Path     []string `yaml:"path"`
}

func (c *Config) AddRemoteConfig(remote *Remote) {
	for _, path := range remote.Path {
		c.addRemoteConfig(remote.Provider, remote.Endpoint, path)
	}
}

func (c *Config) WatchRemoteConfig() {
	for _, yaml := range c.YamlS {
		if !yaml.isRemote {
			continue
		}
		err := yaml.WatchRemoteConfigOnChannel()
		if err != nil {
			panic(err)
		}
	}
}

type watchFunc func() bool

func (c *Config) WatchRemoteConfigFunc(second int64, watchFunc watchFunc) {
	c.WatchRemoteConfig()
	go func() {
		for {
			time.Sleep(time.Duration(second) * time.Second)
			if watchFunc() {
				break
			}
		}
	}()
}

func (c *Config) addRemoteConfig(provider, endpoint string, path string) *Yaml {
	yaml := c.loadRemoteConfig(provider, endpoint, path)
	// add yamlS data
	if cons.ConfOverride {
		c.YamlS = append([]*Yaml{yaml}, c.YamlS...)
	} else {
		c.YamlS = append(c.YamlS, yaml)
	}
	return yaml
}

func (c *Config) loadRemoteConfig(provider, endpoint string, path string) *Yaml {
	yaml := &Yaml{}
	yaml.loadRemoteYaml(provider, endpoint, path)
	return yaml
}
