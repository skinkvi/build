package server

import (
	"context"
	"podbor/internal/app"
	"podbor/internal/config"
	"podbor/internal/handlers"
	"podbor/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	app    *app.App
	Engine *gin.Engine
	logger *zap.Logger
}

func New(app *app.App, cfg *config.Config, logger *zap.Logger) *Server {
	engine := gin.New()
	engine.Use(middleware.GinZapLogger(logger))
	engine.Use(gin.Recovery())

	handlers := handlers.New(app, logger)

	handlers.RegisterRoutes(engine)

	return &Server{
		app:    app,
		Engine: engine,
		logger: logger,
	}
}

func (s *Server) Start() error {
	addr := ":" + s.app.Config.Server.Port
	s.logger.Info("Server started", zap.String("addr", addr))
	return s.Engine.Run(addr)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Сервер останавливается")

	return nil
}
