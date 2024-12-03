package handlers

import (
	"fmt"
	"net/http"
	"podbor/internal/app"
	imageprocessor "podbor/internal/imagepocessor"

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
	router.POST("/fix", h.FixIssue)
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	h.logger.Info("Health check")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handlers) FixIssue(c *gin.Context) {
	// Получение текстового описания проблемы
	description := c.PostForm("description")

	// Получение изображения от пользователя
	file, err := c.FormFile("image")
	if err != nil {
		h.logger.Error("Ошибка при получении файла", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить изображение"})
		return
	}

	// Сохранение изображения на сервере
	imagePath := fmt.Sprintf("uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		h.logger.Error("Ошибка при сохранении файла", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить изображение"})
		return
	}

	// Создание экземпляра ImageProcessor
	imageProcessor := imageprocessor.NewImageProcessor(h.logger, h.app.Config)

	// Анализ изображения
	imageAnalysisResult, err := imageProcessor.AnalyzeImage(imagePath)
	if err != nil {
		h.logger.Error("Ошибка при анализе изображения", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось проанализировать изображение"})
		return
	}

	// Передача данных в AIService для получения ответа
	aiResponse, err := h.app.AIService.ProcessIssue(description, imageAnalysisResult)
	if err != nil {
		h.logger.Error("Ошибка при обработке AI", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обработать запрос"})
		return
	}

	// Поиск материалов в базе данных
	materials, err := h.app.Storage.FindMaterials(aiResponse.Materials)
	if err != nil {
		h.logger.Error("Ошибка при поиске материалов", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось найти материалы"})
		return
	}

	// Отправка ответа пользователю
	c.JSON(http.StatusOK, gin.H{
		"instructions": aiResponse.Instructions,
		"materials":    materials,
	})
}
