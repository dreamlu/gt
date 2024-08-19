// package log

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log log
type Log struct {
	*zap.SugaredLogger
	LogWriter
}

// LogWriter log writer interface
type LogWriter interface {
	Println(v ...any)
}

// NewLog new log
func NewLog() *Log {
	cores := Zap.GetZapCores()
	lgr := zap.New(zapcore.NewTee(cores...))
	lgr.WithOptions(zap.AddCaller())
	log := &Log{}
	log.SugaredLogger = lgr.Sugar()
	return log
}
