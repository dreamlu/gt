package der

import (
	"os"
	"testing"
	"time"
)

var projectPath, _ = os.Getwd()

func init()  {
	NewFileLog(projectPath + "/test/log/", "go-tool.log", 3*time.Second, time.Second)
}

func TestNewFileLog(t *testing.T) {

	Log.Info("项目路径", projectPath)
	for {
		time.Sleep(1 * time.Second)
		Log.Error("测试")
	}
}

func TestDefaultDevModeLog(t *testing.T) {
	DefaultDevModeLog()
}