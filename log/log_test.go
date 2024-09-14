package log

import (
	"github.com/dreamlu/gt/src/type/errors"
	"testing"
	"time"
)

func init() {
	InitProfile()
	GetLog()
}

func TestNewFileLog(t *testing.T) {

	i := 0
	for {
		i++
		if i > 3 {
			break
		}
		Info("number info:", "this is info")
		time.Sleep(1 * time.Second)
		Error("this is error")
	}
	t.Log("log over")
}

func TestErrorLine(t *testing.T) {
	e := errors.New("origin error")
	err2 := errors.Wrap(e, "new err")
	err3 := errors.Wrap(err2, "new err2")
	t.Log(err3)
	//s := fmt.Sprintf("%+v\n", err3)
	Error(err3)
	//GetLog().Error(err3)
}
