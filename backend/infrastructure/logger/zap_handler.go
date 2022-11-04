package logger

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"go.uber.org/zap"
)

type ZapHandler struct {
	config *config.Config
	logger *zap.SugaredLogger
}

func NewZapHandler(config *config.Config) handler.LoggerHandler {
	return newZapHandler(ZapHandler{
		config: config,
	})
}

func newZapHandler(z ZapHandler) handler.LoggerHandler {
	zapHandler := new(ZapHandler)

	myConfig := NewZapConfig(*z.config)

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
