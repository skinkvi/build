package main

import (
	"log"
	"net/http"
	"podbor/db"
	"podbor/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла")
	}
}

func main() {
	db.InitDB()

	router := gin.Default()

	// Загружаем HTML-шаблоны из правильной директории
	router.LoadHTMLGlob("frontend/*.html")

	// Обслуживаем статические файлы
	router.Static("/styles", "frontend/styles")
	router.Static("/img", "frontend/styles/img")

	// Маршрут для главной страницы
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// Маршрут для формы ремонта
	router.GET("/repair", func(c *gin.Context) {
		c.HTML(http.StatusOK, "repair.html", nil)
	})

	router.POST("/fix", handlers.FixHandler)

	// Запускаем сервер
	router.Run(":8080")
}
