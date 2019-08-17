// package der

package der

import (
	time2 "github.com/dreamlu/go-tool/tool/type/time"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

// log
type Log struct {
	Log *logrus.Logger
}

// new log
func (l *Log) NewLog() {

	l.Log = logrus.New()
	l.Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// new output file log
func (l *Log) NewFileLog(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {

	l.NewLog()

	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		l.Log.Errorf("日志文件系统配置错误. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	l.Log.Hooks.Add(lfHook)
}

// Default file log
// maintain 7 days data, every 24 hour split file
func (l *Log) DefaultFileLog(logPath string, logFileName string) {

	l.NewFileLog(logPath, logFileName, time2.Week, time2.Day)
}

// dev/prod/.. mode
// dev mode not output file
// other mode output your project/log/projectName.log
func (l *Log) DefaultDevModeLog() {
	config := &Config{}
	config.NewConfig()
	devMode := config.GetString("devMode")
	if devMode == "dev" {
		l.NewLog()
	} else {
		var projectPath, _ = os.Getwd()
		var pns = strings.Split(projectPath, "/")
		var projectName = pns[len(pns)-1]
		l.DefaultFileLog(projectPath+"/log/", projectName+".log")
	}
}
