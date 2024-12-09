package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"podbor/ai"
	"podbor/db"
	"podbor/models"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
)

func ChoiceHandler(c *gin.Context) {
	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.HTML(http.StatusBadRequest, "choice.html", gin.H{"Error": "Не удалось распарсить форму"})
		log.Printf("Failed to parse form: %v", err)
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusBadRequest, "choice.html", gin.H{"Error": "Изображение не загружено"})
		log.Printf("No image uploaded: %v", err)
		return
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "choice.html", gin.H{"Error": "Не удалось открыть изображение"})
		log.Printf("Failed to open image: %v", err)
		return
	}
	defer file.Close()

	// Save the image
	imagePath, err := saveImage(file, fileHeader)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "choice.html", gin.H{"Error": "Не удалось сохранить изображение"})
		log.Printf("Failed to save image: %v", err)
		return
	}

	log.Printf("Image saved at: %s", imagePath)

	// Run AI to get a descriptive response
	promptText := "Пожалуйста, опишите изображение и перечислите используемые материалы. Ответьте на русском языке."
	responseText, err := ai.RunAI(imagePath, promptText)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "choice.html", gin.H{"Error": "Ошибка анализа изображения AI"})
		log.Printf("Failed to run AI: %v", err)
		return
	}

	log.Printf("AI response: %s", responseText)

	// Extract keywords from AI response
	keywords := extractKeywords(responseText)
	if len(keywords) == 0 {
		c.HTML(http.StatusInternalServerError, "choice.html", gin.H{"Error": "Не удалось извлечь ключевые слова из ответа AI"})
		log.Println("No keywords extracted from AI response")
		return
	}

	log.Printf("Extracted keywords: %v", keywords)

	// Get materials from the database based on keywords
	materials, err := getMaterialsByKeywords(keywords)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "choice.html", gin.H{"Error": "Ошибка получения материалов из базы данных"})
		log.Printf("Failed to get materials: %v", err)
		return
	}

	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(responseText), &buf); err != nil {
		log.Printf("Ошибка при конвертации Markdown: %v", err)
		// Используем простой текст, если конвертация не удалась
		responseText = "Ошибка при обработке ответа от AI."
	}

	// Render the template with materials
	c.HTML(http.StatusOK, "choice.html", gin.H{
		"Materials":  materials,
		"AIResponse": template.HTML(buf.String()),
	})
}

// Функция для получения материалов из базы данных на основе ключевых слов
func getMaterialsByKeywords(keywords []string) ([]models.Material, error) {
	var materials []models.Material

	if len(keywords) == 0 {
		return materials, nil
	}

	dbQuery := db.DB.Model(&models.Material{})
	// Используем первое ключевое слово для начала запроса
	dbQuery = dbQuery.Where("name ILIKE ?", "%"+keywords[0]+"%")

	// Добавляем остальные ключевые слова
	for _, kw := range keywords[1:] {
		dbQuery = dbQuery.Or("name ILIKE ?", "%"+kw+"%")
	}

	// Выполняем запрос
	if err := dbQuery.Find(&materials).Error; err != nil {
		log.Printf("Ошибка запроса к базе данных: %v", err)
		return nil, err
	}

	return materials, nil
}
