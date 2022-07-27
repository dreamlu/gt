package log

import (
	"sync"
)

var (
	onceLog sync.Once
	// global log
	l *Log
)

// OpenLog open once log
func OpenLog(option *Options) {
	onceLog.Do(func() {
		l = NewLog(option)
	})
}

// GetLog one single log
func GetLog() *Log {
	return l
}

// Error sugar
func Error(args ...any) {
	l.ErrorLine(args...)
}

func Info(args ...any) {
	l.Info(args...)
}

func Warn(args ...any) {
	l.Warn(args...)
}
