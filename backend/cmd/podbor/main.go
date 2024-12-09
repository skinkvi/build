package main

import (
	"log"
	"net/http"
	"os"
	"podbor/db"
	"podbor/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file")
	}
}

func main() {
	// Инициализация базы данных
	db.InitDB()

	// Создание нового рутера Gin с стандартными middleware (Logger и Recovery)
	router := gin.Default()

	// Настройка хранилища сессий (используем cookie store)
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options(sessions.Options{
		HttpOnly: true,
		Secure:   false, // Установите true в продакшене с HTTPS
		MaxAge:   3600,  // Срок действия сессии в секундах
		Path:     "/",
	})
	router.Use(sessions.Sessions("mysession", store))

	// Загрузка HTML шаблонов
	router.LoadHTMLGlob("frontend/*.html")

	// Обслуживание статических файлов
	router.Static("/styles", "./frontend/styles")
	router.Static("/img", "./frontend/styles/img")
	router.Static("/uploads", "./uploads") // Обслуживание загруженных изображений

	// Обслуживание favicon.ico
	router.StaticFile("/favicon.ico", "./frontend/styles/img/favicon.ico")

	// Определение маршрутов
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// Защищённый маршрут /repair
	router.GET("/repair", handlers.AuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "repair.html", nil)
	})

	// Маршрут для отображения choice.html
	router.GET("/choice", handlers.AuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "choice.html", nil)
	})

	// Маршрут для обработки отправки формы из choice.html
	router.POST("/choice", handlers.AuthRequired(), handlers.ChoiceHandler)

	router.POST("/fix", handlers.AuthRequired(), handlers.FixHandler)

	// Маршруты аутентификации
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/login", handlers.LoginHandler)
		auth.GET("/logout", handlers.LogoutHandler)
	}

	// Маршруты для входа и регистрации через HTML-файлы
	router.GET("/registration.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "registration.html", nil)
	})

	router.GET("/enter.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "enter.html", nil)
	})

	// Запуск сервера на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
