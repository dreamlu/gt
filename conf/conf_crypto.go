package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dreamlu/gt/src/util"
)

// GenerateEncryptedConfig 同之前逻辑，加密调用改为 util.Encrypt
func (c *Config) GenerateEncryptedConfig() error {
	if !util.HasKey() {
		return errors.New("conf: encrypt key not registered")
	}
	for _, y := range c.YamlS {
		if y.isRemote || y.path == "" || strings.HasSuffix(y.path, ".enc") {
			continue
		}
		data, err := os.ReadFile(y.path)
		if err != nil {
			continue
		}
		enc, err := util.Encrypt(data)
		if err != nil {
			return fmt.Errorf("conf: encrypt %s failed: %w", y.path, err)
		}
		if err := os.WriteFile(y.path+".enc", enc, 0600); err != nil {
			return fmt.Errorf("conf: write %s failed: %w", y.path+".enc", err)
		}
	}
	return nil
}
