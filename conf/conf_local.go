package conf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/dreamlu/gt/src/cons"
	"github.com/dreamlu/gt/src/util"
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
		// embed 模式：优先找 .enc
		path, encrypted := findEmbedPath(c.path)
		c.path = path

		f, err := globalEmbedFS.Open(c.path)
		if err != nil {
			panic(errors.New("embed fs: no yaml: " + c.path))
		}
		defer f.Close()

		if stat, err := fs.Stat(globalEmbedFS, c.path); err != nil || stat.IsDir() {
			panic(errors.New("embed fs: no yaml: " + c.path))
		}

		data, err := io.ReadAll(f)
		if err != nil {
			panic(errors.New("embed fs: read failed: " + c.path))
		}
		if encrypted {
			data, err = util.MaybeDecrypt(data)
			if err != nil {
				panic(errors.New("embed fs: decrypt failed: " + c.path))
			}
		}
		yaml.loadYamlFromReader(bytes.NewReader(data), c.path)

	} else {
		// 文件系统模式：优先找 .enc
		path, encrypted := findFSPath(c.path)
		c.path = path

		data, err := os.ReadFile(c.path)
		if err != nil {
			panic(fmt.Sprintf("can not read %s config, error: %v", c.path, err))
		}
		if encrypted {
			data, err = util.MaybeDecrypt(data)
			if err != nil {
				panic(fmt.Errorf("conf: decrypt %s failed: %w", c.path, err))
			}
		}
		yaml.loadYamlFromReader(bytes.NewReader(data), c.path)
	}

	c.YamlS = append(c.YamlS, yaml)
	return yaml
}

// findFSPath 明文优先，不存在再找 .enc
func findFSPath(path string) (string, bool) {
	if _, err := os.Stat(path); err == nil {
		return path, false // 明文存在，直接用
	}
	encPath := path + ".enc"
	if _, err := os.Stat(encPath); err == nil {
		return encPath, true // 降级到加密文件
	}
	return path, false // 都不存在，返回原路径让上层报错
}

// findEmbedPath 同理
func findEmbedPath(path string) (string, bool) {
	if _, err := fs.Stat(*globalEmbedFS, path); err == nil {
		return path, false
	}
	encPath := path + ".enc"
	if _, err := fs.Stat(*globalEmbedFS, encPath); err == nil {
		return encPath, true
	}
	return path, false
}
