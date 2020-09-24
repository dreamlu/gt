// package gt

package gt

import (
	"fmt"
	time2 "github.com/dreamlu/gt/tool/type/time"
	os2 "github.com/dreamlu/gt/tool/util/os"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// log
type Log struct {
	*logrus.Logger
	LogWriter
}

// LogWriter log writer interface
type LogWriter interface {
	Println(v ...interface{})
}

// log level
const (
	Debug = "debug" // default level
	Info  = "info"
	Warn  = "warn"
	Error = "error"
)

var (
	onceLog sync.Once
	// global log
	l *Log
)

// one single log
func Logger() *Log {

	onceLog.Do(func() {
		l = NewLog()
	})

	return l
}

// new log
func NewLog() *Log {

	lgr := logrus.New()
	lgr.SetFormatter(&myFormatter{})
	log := &Log{}
	log.Logger = lgr

	return log
}

// new output file log
func (l *Log) FileLog(logPath string, logFileName string, maxNum uint, rotationTime time.Duration) {

	if !os2.Exists(logPath) {
		_ = os.Mkdir(logPath, os.ModePerm)
	}

	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		path.Join(logPath, "%Y-%m-%d-"+logFileName),
		rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		//rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
		rotatelogs.WithRotationCount(maxNum),      // 维持的最近日志文件数量
	)
	if err != nil {
		l.Errorf("日志文件系统配置错误. %+v", errors.WithStack(err))
	}

	level := Configger().GetString("app.log.level")
	wm := lfshook.WriterMap{}
	switch level {
	case Debug, "":
		wm[logrus.DebugLevel] = writer
		fallthrough
	case Info:
		wm[logrus.InfoLevel] = writer
		fallthrough
	case Warn:
		wm[logrus.WarnLevel] = writer
		fallthrough
	case Error:
		wm[logrus.ErrorLevel] = writer
		wm[logrus.FatalLevel] = writer
		wm[logrus.PanicLevel] = writer
	}
	lfHook := lfshook.NewHook(wm, l.Formatter)
	l.Hooks.Add(lfHook)
}

// Default file log
// maintain 7 days data, every 24 hour split file
func (l *Log) DefaultFileLog() {

	var (
		//projectPath, _ = os.Getwd()
		projectName = Configger().GetString("app.db.name") // use db name replace
	)
	l.FileLog("log/", projectName+".log", 7, time2.Day)
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

// print error line
func (l *Log) ErrorLine(args ...interface{}) {
	l.Error(fmt.Sprintf("%+v\n", args))
}
