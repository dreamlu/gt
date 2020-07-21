package gt

import (
	"github.com/dreamlu/gt/tool/type/errors"
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

func TestErrorLine(t *testing.T) {
	e := errors.New("origin error")
	err2 := errors.Wrap(e, "new err")
	err3 := errors.Wrap(err2, "new err2")
	t.Log(err3)
	//s := fmt.Sprintf("%+v\n", err3)
	Logger().Error(err3)
	Logger().ErrorLine(err3)
}
