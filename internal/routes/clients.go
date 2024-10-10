package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
)

func Clients(c *gin.Context) {
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

func Client(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"client": "Ok",
	})
}

func AddClient(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"add_client": "Ok",
	})
}

func UpdateClient(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"update_client": "Ok",
	})
}

func DeleteClient(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"delete_client": "Ok",
	})
}
