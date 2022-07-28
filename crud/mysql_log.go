package crud

import (
	"context"
	"fmt"
	"github.com/dreamlu/gt/log"
	"github.com/dreamlu/gt/src/cons/color"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"runtime"
	"time"
)

type Config struct {
	SlowThreshold time.Duration
	Colorful      bool
	LogLevel      gormLog.LogLevel
}

func defaultLog() *log.Log {
	return log.GetLog()
}

func newMysqlLog(config Config) gormLog.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful && runtime.GOOS == "linux" {
		infoStr = color.Reset
		warnStr = color.BlueBold + "%s\n" + color.Reset
		errStr = color.Magenta + "%s\n" + color.Reset
		traceStr = color.Reset + color.Yellow + "[%.3fms] " + color.BlueBold + "[rows:%v]" + color.Reset + " %s"
		traceWarnStr = color.Yellow + "%s\n" + color.Reset + color.RedBold + "[%.3fms] " + color.Yellow + "[rows:%v]" + color.Magenta + " %s" + color.Reset
		traceErrStr = color.RedBold + "%s\n" + color.Reset + color.Yellow + "[%.3fms] " + color.BlueBold + "[rows:%v]" + color.Reset + " %s"
	}

	return &logger{
		Log:          defaultLog(),
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logger struct {
	*log.Log
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= gormLog.Info {
		l.Infof(l.infoStr+msg, append([]any{}, data...)...)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= gormLog.Warn {
		l.Warnf(l.warnStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= gormLog.Error {
		l.Errorf(l.errStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= gormLog.Error:
			sql, rows := fc()
			if rows == -1 {
				l.Errorf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Errorf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLog.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.Warnf(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Warnf(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= gormLog.Info:
			sql, rows := fc()
			if rows == -1 {
				l.Infof(l.traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Infof(l.traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

type traceRecorder struct {
	gormLog.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

func (l traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface, BeginAt: time.Now()}
}

func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}
