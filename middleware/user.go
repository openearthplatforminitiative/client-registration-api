package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.Header.Get("X-Preferred-Username")
		if username == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Not supported without user",
			})
			c.Abort()
			return
		}
		c.Set("user", username)
		c.Next()
	}
}
