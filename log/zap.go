package log

import (
	"fmt"
	. "github.com/dreamlu/gt/src/type/time"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Zap = new(_zap)

type _zap struct{}

func (z *_zap) GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
}

func (z *_zap) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customColorLevelEncoder,
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

func (z *_zap) GetEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer, err := FileRotatelogs.GetWriteSyncer(l.String())
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}

	return zapcore.NewCore(z.GetEncoder(), writer, level)
}

func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format(LayoutN))
}

func (z *_zap) GetZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 4)
	for level := z.ZapLevel(confLogLevel); level <= zapcore.ErrorLevel; level++ {
		cores = append(cores, z.GetEncoderCore(level, z.GetLevelPriority(level)))
	}
	cores = append(cores, z.GetEncoderCore(SuccessZapLevel, z.GetLevelPriority(SuccessZapLevel)))
	return cores
}

func (z *_zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case SuccessZapLevel:
		return func(level zapcore.Level) bool {
			return level == SuccessZapLevel
		}
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool {
			return level == zapcore.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool {
			return level == zapcore.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool {
			return level == zapcore.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool {
			return level == zapcore.ErrorLevel
		}
	default:
		return func(level zapcore.Level) bool {
			return level == zapcore.DebugLevel
		}
	}
}

func (z *_zap) ZapLevel(level string) zapcore.Level {
	switch level {
	case SuccessLevel:
		return SuccessZapLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

const (
	SuccessZapLevel = zapcore.Level(-2)
)

func customColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case SuccessZapLevel:
		enc.AppendString("\033[92msuccess\033[0m")
	default:
		zapcore.LowercaseColorLevelEncoder(level, enc)
	}
}
