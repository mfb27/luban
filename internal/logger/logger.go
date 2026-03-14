package logger

import (
	"strings"

	"github.com/mfb27/luban/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LogConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	_ = level.Set(strings.ToLower(cfg.Level))

	zcfg := zap.NewProductionConfig()
	zcfg.Level = zap.NewAtomicLevelAt(level)
	zcfg.EncoderConfig.TimeKey = "ts"
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zcfg.Build()
}

