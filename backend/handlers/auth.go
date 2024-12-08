package handlers

import (
	"net/http"
	"podbor/db"
	"podbor/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Структура для регистрации
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Обработчик регистрации
func RegisterHandler(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": err.Error()})
		return
	}

	// Проверка совпадения паролей
	passwordConfirm := c.PostForm("password_confirm")
	if input.Password != passwordConfirm {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": "Пароли не совпадают"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "registration.html", gin.H{"error": "Ошибка сервера"})
		return
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": "Пользователь с таким именем или email уже существует"})
		return
	}

	// Создание сессии
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "registration.html", gin.H{"error": "Не удалось сохранить сессию"})
		return
	}

	// Перенаправление на главную страницу после успешной регистрации
	c.Redirect(http.StatusSeeOther, "/")
}

// Структура для входа
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Обработчик входа
func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "enter.html", gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "enter.html", gin.H{"error": "Неверный email ��ли пароль"})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "enter.html", gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Создание сессии
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "enter.html", gin.H{"error": "Не удалось сохранить сессию"})
		return
	}

	// Перенаправление на главную страницу после успешного входа
	c.Redirect(http.StatusSeeOther, "/")
}

// Обработчик выхода
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "/login", gin.H{"error": "Не удалось очистить сессию"})
		return
	}
	// Перенаправление на страницу входа после выхода
	c.Redirect(http.StatusSeeOther, "/enter.html")
}
