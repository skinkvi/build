package main

import (
	"log"
	"podbor/internal/app"
	"podbor/internal/config"
	"podbor/internal/logger"
	"podbor/internal/server"
	"podbor/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = logger.InitLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.SyncLogger()

	st, err := storage.New(config.Cfg, logger.Log)
	if err != nil {
		logger.Log.Fatal("Не удалось инициализировать хранилище: ", zap.Error(err))
	}
	defer st.Close()

	gin.SetMode(gin.ReleaseMode)

	appInstance := app.New(st, logger.Log, config.Cfg)

	srv := server.New(appInstance, config.Cfg, logger.Log)
	srv.Engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	if err := srv.Start(); err != nil {
		logger.Log.Fatal("Не удалось запустить сервер: ", zap.Error(err))
	}
}
