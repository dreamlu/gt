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
	var ss []any
	for _, arg := range args {
		ss = append(ss, fmt.Sprintf("%+v", arg))
	}
	GetLog().Error(ss...)
}

func Success(args ...any) {
	GetLog().Log(SuccessZapLevel, args...)
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

func Successf(format string, args ...any) {
	GetLog().Log(SuccessZapLevel, fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
	GetLog().Error(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...any) {
	GetLog().Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
	GetLog().Warn(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...any) {
	GetLog().Debug(fmt.Sprintf(format, args...))
}
