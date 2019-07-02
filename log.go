// @author  dreamlu
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

// new log
func NewLog() *logrus.Logger {

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return log
}

// new output file log
func NewFileLog(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *logrus.Logger {

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("日志文件系统配置错误. %+v", errors.WithStack(err))
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
	log.Hooks.Add(lfHook)
	return log
}

// Default file log
// maintain 7 days data, every 24 hour split file
func DefaultFileLog(logPath string, logFileName string) *logrus.Logger {

	return NewFileLog(logPath, logFileName, time2.Week, time2.Day)
}

// dev/prod/.. mode
// dev mode not output file
// other mode output your project/log/projectName.log
func DefaultDevModeLog() *logrus.Logger {
	devMode := GetConfigValue("devMode")
	if devMode == "dev" {
		return NewLog()
	} else {
		var projectPath, _ = os.Getwd()
		var pns = strings.Split(projectPath, "/")
		var projectName = pns[len(pns)-1]
		return DefaultFileLog(projectPath+"/log/", projectName+".log")
	}
}
