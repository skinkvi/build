package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequired проверяет, аутентифицирован ли пользователь
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			// Перенаправление на страницу входа
			c.Redirect(http.StatusSeeOther, "/enter.html")
			c.Abort()
			return
		}
		// Можно дополнительно проверить существование пользователя в БД
		c.Next()
	}
}
