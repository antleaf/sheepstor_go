package internal

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func NewZapSugarLogger(debugging bool) (*zap.SugaredLogger, error) {
	var zapLogger *zap.Logger
	var err error
	zapLogger, err = NewZapLogger(debugging)
	return zapLogger.Sugar(), err
}

func NewZapLogger(debugging bool) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
	}
	if debugging == true {
		level = zapcore.DebugLevel
		encoderConfig = zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:   "level",
			// TimeKey:      "timestamp",
			EncodeTime:   zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000"),
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			// StacktraceKey: "stacktrace",
		}
	}
	zapConfig := zap.Config{
		Encoding:      "console",
		Level:         zap.NewAtomicLevelAt(level),
		OutputPaths:   []string{"stdout"},
		EncoderConfig: encoderConfig,
	}
	return zapConfig.Build()
}
