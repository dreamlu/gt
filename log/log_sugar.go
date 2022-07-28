package log

import "fmt"

var l *Log

// gt tool log, only use by gt
func init() {
	l = NewLog()
}

// GetLog get once log
func GetLog() *Log {
	return l
}

func Error(args ...any) {
	l.Error(fmt.Sprintf("%+v\n", args))
}

func Info(args ...any) {
	l.Info(args...)
}

func Warn(args ...any) {
	l.Warn(args...)
}
