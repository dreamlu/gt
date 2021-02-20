package daemon

import (
	"fmt"
	"github.com/dreamlu/gt/tool/type/time"
	"testing"
	time2 "time"
)

func TestNewDaemon(t *testing.T) {

	for i := 0; i < 15; i++ {
		Daemoner().
			AddTask(
				Func(f),
			)
		t.Log(i)
	}
	for i := 0; i < 5; i++ {
		Daemoner().
			AddTask(
				Func(f2),
				Time(time.ParseCTime("2020-07-08 16:00:00")),
			)
		t.Log(i)
	}

	time2.Sleep(25 * time2.Second)
}

func f() {
	fmt.Println("f func")
	time2.Sleep(3 * time2.Second)
}

func f2() {
	fmt.Println("f2 func")
	time2.Sleep(3 * time2.Second)
}
