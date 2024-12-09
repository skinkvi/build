package handlers

import (
	"log"
	"net/http"
	"podbor/db"
	"podbor/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Структура для регистрации
type RegisterInput struct {
	Username        string `form:"username" binding:"required"`
	Email           string `form:"email" binding:"required,email"`
	Password        string `form:"password" binding:"required,min=6"`
	PasswordConfirm string `form:"password_confirm" binding:"required,min=6"`
}

// Структура для входа
type LoginInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

// Обработчик регистрации
func RegisterHandler(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Ошибка связывания данных регистрации: %v", err)
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": "Некорректные данные. Пожалуйста, заполните форму правильно."})
		return
	}

	// Проверка совпадения паролей
	if input.Password != input.PasswordConfirm {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": "Пароли не совпадают."})
		return
	}

	// Проверка существования пользователя
	var existingUser models.User
	if err := db.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusConflict, "registration.html", gin.H{"error": "Пользователь с таким email уже существует."})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		c.HTML(http.StatusInternalServerError, "registration.html", gin.H{"error": "Не удалось обработать пароль. Попробуйте позже."})
		return
	}

	// Создание нового пользователя
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		c.HTML(http.StatusInternalServerError, "registration.html", gin.H{"error": "Не удалось зарегистрировать пользователя. Попробуйте позже."})
		return
	}

	// Создание сессии
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Printf("Ошибка при сохранении сессии: %v", err)
		c.HTML(http.StatusInternalServerError, "registration.html", gin.H{"error": "Не удалось сохранить сессию. Попробуйте позже."})
		return
	}

	// Перенаправление на главную страницу после успешной регистрации
	c.Redirect(http.StatusSeeOther, "/")
	return
}

// Обработчик входа
func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Ошибка связывания данных входа: %v", err)
		c.HTML(http.StatusBadRequest, "enter.html", gin.H{"error": "Некорректные данные. Пожалуйста, заполните форму правильно."})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("Пользователь не найден: %s", input.Email)
		c.HTML(http.StatusUnauthorized, "enter.html", gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		log.Printf("Неверный пароль для пользователя: %s", input.Email)
		c.HTML(http.StatusUnauthorized, "enter.html", gin.H{"error": "Неверный email или пароль"})
		return
	}

	log.Printf("Пользователь успешно вошёл: %s", user.Username)

	// Создание сессии
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Printf("Ошибка при сохранении сессии: %v", err)
		c.HTML(http.StatusInternalServerError, "enter.html", gin.H{"error": "Не удалось сохранить сессию. Попробуйте позже."})
		return
	}

	log.Printf("Сессия создана для пользователя ID: %d", user.ID)

	// Перенаправление на главную страницу после успешного входа
	c.Redirect(http.StatusSeeOther, "/")
	return
}

// Обработчик выхода
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Printf("Ошибка при очистке сессии: %v", err)
		c.HTML(http.StatusInternalServerError, "enter.html", gin.H{"error": "Не удалось очистить сессию"})
		return
	}
	log.Println("Пользователь успешно вышел")
	// Перенаправление на страницу входа после выхода
	c.Redirect(http.StatusSeeOther, "/enter.html")
	return
}

// Обработчики для отображения страниц
func EnterPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "enter.html", nil)
}

func RegistrationPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", nil)
}
