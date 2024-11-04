package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}

func Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Ready",
	})
}

func Dump(c *gin.Context) {
	dump, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("%q", dump)

	c.String(http.StatusOK, string(dump))
}
