package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct{}

func (r *fileRotatelogs) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	switch logIn {
	case InFile:
		fileWriter, err := r.GetFileWriter(level)
		return zapcore.AddSync(fileWriter), err
	case InConsole:
		return zapcore.AddSync(os.Stdout), nil
	default:
		fileWriter, err := r.GetFileWriter(level)
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
}

func (r *fileRotatelogs) GetFileWriter(level string) (*rotatelogs.RotateLogs, error) {
	fileWriter, err := rotatelogs.New(
		path.Join(confLogDirector, "%Y-%m-%d", level+".log"),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(confLogMaxAge)*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	return fileWriter, err
}
