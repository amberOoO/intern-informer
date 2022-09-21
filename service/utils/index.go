package utils

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapSingleton *zap.Logger
var loadZapOnce sync.Once

func GetZapLogger(fields ...zapcore.Field) *zap.Logger {
	loadZapOnce.Do(func() {
		zapSingleton, _ = zap.NewProduction()
	})
	return zapSingleton.With(fields...)
}
