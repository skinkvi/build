package handlers

import (
	"net/http"
	"podbor/internal/app"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	app    *app.App
	logger *zap.Logger
}

func New(app *app.App, logger *zap.Logger) *Handlers {
	return &Handlers{
		app:    app,
		logger: logger,
	}
}

func (h *Handlers) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.HealthCheck)
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	h.logger.Info("Health check")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
