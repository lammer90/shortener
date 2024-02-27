package logger

import (
	"go.uber.org/zap"
)

// Log основной рабочий логгер
var Log *zap.Logger = zap.NewNop()

// InitLogger активировать работу логгера
func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}
