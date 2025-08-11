package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct {
	once sync.Once
}

func (r *fileRotatelogs) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {

	r.once.Do(func() {
		r.cleanOldLogs(confLogDirector, confLogMaxAge)
	})

	fileWriter, err := r.GetFileWriter(level)
	return zapcore.AddSync(fileWriter), err
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

// cleanOldLogs run on the start
func (r *fileRotatelogs) cleanOldLogs(dir string, maxAgeDays int) {
	cutoff := time.Now().Add(-time.Duration(maxAgeDays) * 24 * time.Hour)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if t, err := time.Parse("2006-01-02", info.Name()); err == nil {
				if t.Before(cutoff) {
					_ = os.RemoveAll(path)
				}
			}
		}
		return nil
	})
}
