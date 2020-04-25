package gt

import (
	"os"
	"testing"
	"time"
)

var projectPath, _ = os.Getwd()

func init() {
	Logger().FileLog(projectPath+"/test/log/", "gt.log", 3, time.Minute)
}

func TestNewFileLog(t *testing.T) {

	i := 0
	for {
		i++
		if i > 3 {
			break
		}
		Logger().Info("this is info")
		time.Sleep(1 * time.Second)
		Logger().Error("this is error")
	}
	t.Log("log over")
}
