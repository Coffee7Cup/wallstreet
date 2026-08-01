package logs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	CategoryEngine = "engine"
)

var Log *zap.Logger

// filteredConsoleCore only allows Error level and above, or Warn level logs with the "engine" category.
type filteredConsoleCore struct {
	zapcore.Core
}

func (c *filteredConsoleCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if ent.Level >= zapcore.ErrorLevel {
		return c.Core.Write(ent, fields)
	}
	if ent.Level == zapcore.WarnLevel {
		for _, f := range fields {
			if f.Key == "category" && f.String == CategoryEngine {
				return c.Core.Write(ent, fields)
			}
		}
	}
	return nil
}

func (c *filteredConsoleCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if ent.Level >= zapcore.ErrorLevel {
		return c.Core.Check(ent, ce)
	}
	if ent.Level == zapcore.WarnLevel {
		return ce.AddCore(ent, c)
	}
	return ce
}

func Init() {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "ts",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/engine.log",
		MaxSize:    50, // MB
		MaxBackups: 7,
		MaxAge:     14, // days
		Compress:   true,
	})

	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			fileWriter,
			zapcore.InfoLevel, // file: info+
		),
		&filteredConsoleCore{
			Core: zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderCfg),
				consoleWriter,
				zapcore.WarnLevel, // console: Warn+ (filtered further by filteredConsoleCore)
			),
		},
	)

	Log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}
