package log

import (
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/src/cons"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct{}

func (r *fileRotatelogs) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	fileWriter, err := rotatelogs.New(
		path.Join(conf.Get[string](cons.ConfLogDirector), "%Y-%m-%d", level+".log"),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(conf.Get[int](cons.ConfLogMaxAge))*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	switch conf.Get[string]("app.log.log-in") {
	case InFile:
		return zapcore.AddSync(fileWriter), err
	case InConsole:
		return zapcore.AddSync(os.Stdout), err
	default:
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
}
