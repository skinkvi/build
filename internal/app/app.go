package app

import (
	"podbor/internal/ai"
	"podbor/internal/config"
	imageprocessor "podbor/internal/imagepocessor"
	"podbor/internal/storage"

	"go.uber.org/zap"
)

type App struct {
	Storage        *storage.Storage
	Logger         *zap.Logger
	Config         *config.Config
	AIService      *ai.AIService
	ImageProcessor *imageprocessor.ImageProcessor
}

func New(storage *storage.Storage, logger *zap.Logger, config *config.Config) *App {
	aiService := ai.NewAIService(logger, config)
	imageProcessor := imageprocessor.NewImageProcessor(logger, config)
	return &App{
		Storage:        storage,
		Logger:         logger,
		Config:         config,
		AIService:      aiService,
		ImageProcessor: imageProcessor,
	}
}
