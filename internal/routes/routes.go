package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	router.GET("/health", Health)
	router.GET("/ready", Ready)
	router.GET("/dump", Dump)

	api := router.Group("/clients")
	api.GET("/", Clients)
	api.GET("/:id", Client)
	api.POST("/", AddClient)
	api.PUT("/:id", UpdateClient)
	api.DELETE("/:id", DeleteClient)
}
