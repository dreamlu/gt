package log

import (
	"fmt"
	"sync"
)

var (
	onceLog sync.Once
	l       *Log
)

// gt tool log, only use by gt
//func init() {
//	l = NewLog()
//}

// GetLog get once log, only as log
// maybe use InitProfile to init some param
func GetLog() *Log {
	onceLog.Do(func() {
		l = NewLog()
	})
	return l
}

func Error(args ...any) {
	GetLog().Error(fmt.Sprintf("%+v", args...))
}

func Info(args ...any) {
	GetLog().Info(args...)
}

func Warn(args ...any) {
	GetLog().Warn(args...)
}

func Debug(args ...any) {
	GetLog().Debug(args...)
}
