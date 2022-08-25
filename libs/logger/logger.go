package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log     *zap.Logger
	Console *zap.Logger
)

func New() {
	fileLogger := &lumberjack.Logger{
		Filename: "./logs/logging.log",
		MaxSize:  50,
		MaxAge:   7,
	}

	cfg := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(cfg)
	consoleEncoder := zapcore.NewJSONEncoder(cfg)

	Log = zap.New(zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileLogger), zap.LevelEnablerFunc(func(_ zapcore.Level) bool {
			return true
		})),
	))

	Console = zap.New(zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	))
}
