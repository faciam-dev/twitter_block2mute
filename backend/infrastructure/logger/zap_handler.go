package logger

import (
	"strings"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapHandler struct {
	LoggerLevel            string
	LoggerOutputPaths      []string
	LoggerErrorOutputPaths []string
	logger                 *zap.SugaredLogger
}

func NewZapHandler(config *config.Config) handler.LoggerHandler {
	return newZapHandler(ZapHandler{
		LoggerLevel:            config.Logger.LoggerLevel,
		LoggerOutputPaths:      config.Logger.LoggerOutputPaths,
		LoggerErrorOutputPaths: config.Logger.LoggerErrorOutputPaths,
	})
}

func newZapHandler(z ZapHandler) handler.LoggerHandler {
	zapHandler := new(ZapHandler)

	level := zap.NewAtomicLevel()
	switch strings.ToLower(z.LoggerLevel) {
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
		OutputPaths:      z.LoggerOutputPaths,
		ErrorOutputPaths: z.LoggerErrorOutputPaths,
	}

	logger, _ := myConfig.Build()
	//logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	zapHandler.logger = sugar

	// for zap
	return zapHandler
}

// debug
func (z *ZapHandler) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

func (z *ZapHandler) Debugf(template string, args ...interface{}) {
	z.logger.Debugf(template, args...)
}

func (z *ZapHandler) Debugw(msg string, keysAndValues ...interface{}) {
	z.logger.Debugw(msg, keysAndValues...)
}

// info
func (z *ZapHandler) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z *ZapHandler) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}

func (z *ZapHandler) Infow(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

// warning
func (z *ZapHandler) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

func (z *ZapHandler) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}

func (z *ZapHandler) Warnw(msg string, keysAndValues ...interface{}) {
	z.logger.Warnw(msg, keysAndValues...)
}

// error
func (z *ZapHandler) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z *ZapHandler) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}

func (z *ZapHandler) Errorw(msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, keysAndValues...)
}

// panic
func (z *ZapHandler) Panic(args ...interface{}) {
	z.logger.Panic(args...)
}

func (z *ZapHandler) Panicf(template string, args ...interface{}) {
	z.logger.Panicf(template, args...)
}

func (z *ZapHandler) Panicw(msg string, keysAndValues ...interface{}) {
	z.logger.Panicw(msg, keysAndValues...)
}

// fatal
func (z *ZapHandler) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

func (z *ZapHandler) Fatalf(template string, args ...interface{}) {
	z.logger.Fatalf(template, args...)
}

func (z *ZapHandler) Fatalw(msg string, keysAndValues ...interface{}) {
	z.logger.Fatalw(msg, keysAndValues...)
}
