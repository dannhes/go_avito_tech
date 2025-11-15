package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func Init() {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	L = l
}

func Sync() {
	_ = L.Sync()
}
