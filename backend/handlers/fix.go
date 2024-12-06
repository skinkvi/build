package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"podbor/ai"
	"time"

	"github.com/gin-gonic/gin"
)

func FixHandler(c *gin.Context) {
	// Парсим форму
	err := c.Request.ParseMultipartForm(10 << 20) // 10 МБ макс. памяти
	if err != nil {
		c.HTML(http.StatusBadRequest, "repair.html", gin.H{
			"error": "Не удалось распарсить форму",
		})
		return
	}

	promptText := c.PostForm("prompt_text")
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusBadRequest, "repair.html", gin.H{
			"error": "Не удалось получить изображение",
		})
		return
	}

	// Открываем загруженный файл
	file, err := fileHeader.Open()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "repair.html", gin.H{
			"error": "Не удалось открыть загруженное изображение",
		})
		return
	}
	defer file.Close()

	// Сохраняем загруженное изображение
	imagePath, err := saveImage(file, fileHeader)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "repair.html", gin.H{
			"error": "Не удалось сохранить изображение",
		})
		return
	}

	// Передаем путь к изображению и текст запроса в функцию AI
	responseText, err := ai.RunAI(imagePath, promptText)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "repair.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Отображаем шаблон с ответом AI
	c.HTML(http.StatusOK, "repair.html", gin.H{
		"ai_response": responseText,
	})
}

func saveImage(file io.Reader, header *multipart.FileHeader) (string, error) {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return "", err
	}

	// Генерация уникального имени файла с помощью временной метки
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	imagePath := filepath.Join("uploads", filename)

	outFile, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}
