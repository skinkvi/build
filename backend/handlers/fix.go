package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"podbor/ai"
	"podbor/db"
	"podbor/models"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
)

func FixHandler(c *gin.Context) {
	// Парсим форму
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 МБ макс память
		c.HTML(http.StatusBadRequest, "repair.html", gin.H{"error": "Не удалось распарсить форму"})
		log.Printf("Не удалось распарсить форму: %v", err)
		return
	}

	promptText := c.PostForm("description")
	log.Printf("Полученный текст описания: %s", promptText)

	if strings.TrimSpace(promptText) == "" {
		c.HTML(http.StatusBadRequest, "repair.html", gin.H{"error": "Описание не должно быть пустым"})
		log.Println("Описание пустое")
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusBadRequest, "repair.html", gin.H{"error": "Не удалось получить изображение"})
		log.Printf("Не удалось получить изображение: %v", err)
		return
	}

	// Открываем загруженный файл
	file, err := fileHeader.Open()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "repair.html", gin.H{"error": "Не удалось открыть загруженное изображение"})
		log.Printf("Не удалось открыть загруженное изображение: %v", err)
		return
	}
	defer file.Close()

	// Сохраняем загруженное изображение
	imagePath, err := saveImage(file, fileHeader)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "repair.html", gin.H{"error": "Не удалось сохранить изображение"})
		log.Printf("Не удалось сохранить изображение: %v", err)
		return
	}

	log.Printf("Изображение сохранено по пути: %s", imagePath)

	// Запуск AI для получения ответа
	responseText, err := ai.RunAI(imagePath, promptText+"\n\nПожалуйста, ответьте на русском языке, распишите решение по 5 пунктам и перечислите необходимые материалы.")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "repair.html", gin.H{"error": err.Error()})
		log.Printf("Не удалось запустить AI: %v", err)
		return
	}

	log.Printf("Ответ от AI: %s", responseText)

	// Извлечение необходимых материалов из базы данных
	requiredMaterials, err := getRequiredMaterials(promptText)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "repair.html", gin.H{"error": "Не удалось извлечь материалы из базы данных"})
		log.Printf("Не удалось извлечь материалы из базы данных: %v", err)
		return
	}

	log.Printf("Материалы для отображения: %v", requiredMaterials)

	// Конвертация ответа AI из Markdown в HTML
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(responseText), &buf); err != nil {
		log.Printf("Ошибка при конвертации Markdown: %v", err)
		// Используем простой текст, если конвертация не удалась
		responseText = "Ошибка при обработке ответа от AI."
	}

	// Подготовка данных для шаблона
	data := gin.H{
		"ai_response":        template.HTML(buf.String()), // Помечаем как безопасный HTML
		"required_materials": requiredMaterials,
	}

	// Рендеринг шаблона с ответом AI и материалами
	c.HTML(http.StatusOK, "repair.html", data)
}

// saveImage сохраняет загруженное изображение в директорию /uploads и возвращает путь к изображению.
func saveImage(file io.Reader, header *multipart.FileHeader) (string, error) {
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Printf("Не удалось создать директорию uploads: %v", err)
		return "", err
	}

	// Генерируем уникальное имя файла с помощью временной метки
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	imagePath := filepath.Join(uploadsDir, filename)

	outFile, err := os.Create(imagePath)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		return "", err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		log.Printf("Ошибка при сохранении файла: %v", err)
		return "", err
	}

	return imagePath, nil
}

// getRequiredMaterials извлекает материалы на основе текста запроса, используя ILIKE для поиска подстрок.
func getRequiredMaterials(prompt string) ([]models.Material, error) {
	var materials []models.Material

	// Извлекаем ключевые слова из запроса
	keywords := extractKeywords(prompt)
	log.Printf("Используемые ключевые слова для запроса к базе данных: %v", keywords)
	if len(keywords) == 0 {
		log.Printf("Не удалось извлечь ключевые слова из запроса")
		return materials, nil // Возвращаем пустой срез, если нет ключевых слов
	}

	// Создаем начальный запрос к модели Material
	dbQuery := db.DB.Model(&models.Material{})

	// Используем первый ключ для инициализации WHERE
	dbQuery = dbQuery.Where("name ILIKE ?", "%"+keywords[0]+"%")

	// Для остальных ключей добавляем OR
	for _, kw := range keywords[1:] {
		dbQuery = dbQuery.Or("name ILIKE ?", "%"+kw+"%")
	}

	// Выполняем запрос
	if err := dbQuery.Find(&materials).Error; err != nil {
		log.Printf("Ошибка при выполнении запроса к базе данных: %v", err)
		return nil, err
	}

	log.Printf("Найдено материалов: %d", len(materials))
	return materials, nil
}

// extractKeywords извлекает значимые ключевые слова из запроса.
func extractKeywords(prompt string) []string {
	words := strings.FieldsFunc(prompt, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	var keywords []string
	stopWords := map[string]bool{
		"порвался": true,
		"помоги":   true,
		"починить": true,
		// Добавьте другие общие слова
	}
	for _, word := range words {
		cleanWord := strings.ToLower(strings.Trim(word, ".,!?"))
		if len(cleanWord) > 2 && !stopWords[cleanWord] {
			keywords = append(keywords, cleanWord)
		}
	}
	log.Printf("Извлечённые ключевые слова: %v", keywords)
	return keywords
}
