package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Clients(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"all_clients": "Ok",
	})
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
