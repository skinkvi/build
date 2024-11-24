package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	Log = logger
	return nil
}

func SyncLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
