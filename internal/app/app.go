package app

import (
	"podbor/internal/config"
	"podbor/internal/storage"

	"go.uber.org/zap"
)

type App struct {
	Storage *storage.Storage
	Logger  *zap.Logger
	Config  *config.Config
}

func New(storage *storage.Storage, logger *zap.Logger, config *config.Config) *App {
	return &App{
		Storage: storage,
		Logger:  logger,
		Config:  config,
	}
}
