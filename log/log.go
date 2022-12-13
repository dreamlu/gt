// package log

package log

import (
	"fmt"
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/src/cons"
	"github.com/dreamlu/gt/src/cons/color"
	"github.com/dreamlu/gt/src/gos"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
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

type logFormatter struct{}

func (s *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var (
		msg  string
		info = "[%s] [%s] %s\n"
	)
	switch runtime.GOOS {
	case "linux":
		if entry.Level == logrus.ErrorLevel {
			info = color.RedBold + "[%s] [%s] %s\n" + color.Reset
		} else {
			info = color.Yellow + "[%s] " + color.Cyan + "[%s] " + color.Reset + "%s\n"
		}
	}
	msg = fmt.Sprintf(info, timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

// log level
const (
	DebugLevel = "debug" // default level
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

// NewLog new log
func NewLog() *Log {
	lgr := logrus.New()
	lgr.SetFormatter(&logFormatter{})
	log := &Log{}
	log.Logger = lgr
	return log
}

// ConfigLog init log config
// store log to file
// options:
// &log.Options{
//		LogPath:      "log/",
//		LogFileName:  "app.log",
//		MaxNum:       uint(7),
//		RotationTime: ctime.Day,
//	}
func (l *Log) ConfigLog(options *Options) {
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
