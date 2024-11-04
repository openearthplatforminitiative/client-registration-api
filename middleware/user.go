package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.Header.Get("X-Auth-Request-Preferred-Username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Not supported without user",
			})
			c.Abort()
			return
		}
		c.Set("user", username)
		c.Next()
	}
}
