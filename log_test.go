package gt

import (
	"os"
	"testing"
	"time"
)

var projectPath, _ = os.Getwd()

//var myLog = Logger()

func init() {
	Logger().FileLog(projectPath+"/test/log/", "go-tool.log", 3, time.Minute)
}

func TestNewFileLog(t *testing.T) {

	for {
		Logger().Info("[debug]")
		time.Sleep(1 * time.Second)
		Logger().Error("[error]")
	}
}
