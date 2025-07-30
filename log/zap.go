package log

import (
	"fmt"
	. "github.com/dreamlu/gt/src/type/time"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Zap = new(_zap)

type _zap struct{}

func (z *_zap) GetEncoder(isConsole bool) zapcore.Encoder {
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig(isConsole))
}

func (z *_zap) GetEncoderConfig(isConsole bool) zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:    customColorLevelEncoder,
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if isConsole {
		encoderConfig.EncodeLevel = customColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder // 纯文本
	}
	return encoderConfig
}

func (z *_zap) GetEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc, isConsole bool) zapcore.Core {
	writer, err := FileRotatelogs.GetWriteSyncer(l.String())
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}

	return zapcore.NewCore(z.GetEncoder(isConsole), writer, level)
}

func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format(LayoutN))
}

func (z *_zap) GetZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 10)

	for level := z.ZapLevel(confLogLevel); level <= zapcore.ErrorLevel; level++ {
		levelEnabler := z.GetLevelPriority(level)

		switch logIn {
		case InTerminal:
			// 只输出终端
			consoleCore := zapcore.NewCore(z.GetEncoder(true), zapcore.AddSync(os.Stdout), levelEnabler)
			cores = append(cores, consoleCore)
		case InFile:
			// 只输出文件
			fileWriter, err := FileRotatelogs.GetWriteSyncer(level.String())
			if err == nil {
				fileCore := zapcore.NewCore(z.GetEncoder(false), fileWriter, levelEnabler)
				cores = append(cores, fileCore)
			}
		case InAll:
			// 终端 + 文件
			consoleCore := zapcore.NewCore(z.GetEncoder(true), zapcore.AddSync(os.Stdout), levelEnabler)
			cores = append(cores, consoleCore)

			fileWriter, err := FileRotatelogs.GetWriteSyncer(level.String())
			if err == nil {
				fileCore := zapcore.NewCore(z.GetEncoder(false), fileWriter, levelEnabler)
				cores = append(cores, fileCore)
			}
		}
	}

	// success level 同理处理
	levelEnabler := z.GetLevelPriority(SuccessZapLevel)
	switch logIn {
	case InTerminal:
		consoleCore := zapcore.NewCore(z.GetEncoder(true), zapcore.AddSync(os.Stdout), levelEnabler)
		cores = append(cores, consoleCore)
	case InFile:
		fileWriter, err := FileRotatelogs.GetWriteSyncer(SuccessZapLevel.String())
		if err == nil {
			fileCore := zapcore.NewCore(z.GetEncoder(false), fileWriter, levelEnabler)
			cores = append(cores, fileCore)
		}
	case InAll:
		consoleCore := zapcore.NewCore(z.GetEncoder(true), zapcore.AddSync(os.Stdout), levelEnabler)
		cores = append(cores, consoleCore)

		fileWriter, err := FileRotatelogs.GetWriteSyncer(SuccessZapLevel.String())
		if err == nil {
			fileCore := zapcore.NewCore(z.GetEncoder(false), fileWriter, levelEnabler)
			cores = append(cores, fileCore)
		}
	}

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
