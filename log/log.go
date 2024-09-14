// package log

package log

import (
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/src/cons"
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

// InitProfile init log conf replace default
// if you only use log print, not call this func
func InitProfile() {
	confLogLevel = conf.Get[string](cons.ConfLogLevel)
	confLogDirector = conf.Get[string](cons.ConfLogDirector)
	confLogMaxAge = conf.Get[int](cons.ConfLogMaxAge)
	logIn = conf.Get[string](cons.ConfLogLogIn)
}
