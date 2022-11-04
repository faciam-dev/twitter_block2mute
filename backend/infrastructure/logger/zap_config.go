package logger

import (
	"strings"

	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapConfig(config config.Config) zap.Config {

	level := zap.NewAtomicLevel()
	switch strings.ToLower(config.Logger.LoggerLevel) {
	case "info":
		level.SetLevel(zapcore.InfoLevel)
	case "debug":
		level.SetLevel(zapcore.DebugLevel)
	case "warn":
		level.SetLevel(zapcore.WarnLevel)
	default:
		level.SetLevel(zapcore.ErrorLevel)

	}

	myConfig := zap.Config{
		Level:    level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      config.Logger.LoggerOutputPaths,
		ErrorOutputPaths: config.Logger.LoggerErrorOutputPaths,
	}

	return myConfig
}
