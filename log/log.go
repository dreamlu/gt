// package log

package log

import (
	"fmt"
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/crud/dep/cons"
	"github.com/dreamlu/gt/src/gos"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

// Log log
type Log struct {
	*logrus.Logger
	LogWriter
}

type Options struct {
	LogPath      string
	LogFileName  string
	MaxNum       uint
	RotationTime time.Duration
}

// LogWriter log writer interface
type LogWriter interface {
	Println(v ...any)
}

// log level
const (
	DebugLevel = "debug" // default level
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

// NewLog new log
func NewLog(option *Options) *Log {
	lgr := logrus.New()
	lgr.SetFormatter(&myFormatter{})
	log := &Log{}
	log.Logger = lgr
	log.InitLog(option)
	return log
}

// InitLog init log config
func (l *Log) InitLog(options *Options) {
	if !gos.Exists(options.LogPath) {
		_ = os.Mkdir(options.LogPath, os.ModePerm)
	}

	baseLogPath := path.Join(options.LogPath, options.LogFileName)
	writer, err := rotatelogs.New(
		path.Join(options.LogPath, "%Y-%m-%d-"+options.LogFileName),
		rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		//rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(options.RotationTime), // 日志切割时间间隔
		rotatelogs.WithRotationCount(options.MaxNum),      // 维持的最近日志文件数量
	)
	if err != nil {
		l.Errorf("日志文件系统配置错误. %+v", errors.WithStack(err))
	}

	level := conf.Get[string](cons.ConfLogLevel)
	wm := lfshook.WriterMap{}
	switch level {
	case DebugLevel, "":
		wm[logrus.DebugLevel] = writer
		fallthrough
	case InfoLevel:
		wm[logrus.InfoLevel] = writer
		fallthrough
	case WarnLevel:
		wm[logrus.WarnLevel] = writer
		fallthrough
	case ErrorLevel:
		wm[logrus.ErrorLevel] = writer
		wm[logrus.FatalLevel] = writer
		wm[logrus.PanicLevel] = writer
	}
	lfHook := lfshook.NewHook(wm, l.Formatter)
	l.Hooks.Add(lfHook)
}

// gt log formatter
type myFormatter struct{}

func (s *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	msg := ""
	if entry.Level == logrus.ErrorLevel {
		msg = fmt.Sprintf("\u001B[36;31m[%s] [%s] %s\u001B[0m\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	} else {
		msg = fmt.Sprintf("\u001B[33m[%s]\u001B[0m \u001B[36;1m[%s]\u001B[0m %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	}
	return []byte(msg), nil
}

// ErrorLine print error line
func (l *Log) ErrorLine(args ...any) {
	l.Error(fmt.Sprintf("%+v\n", args))
}
