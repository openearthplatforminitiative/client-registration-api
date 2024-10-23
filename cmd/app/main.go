package main

import (
	"github.com/gin-gonic/gin"
	"github.com/openearthplatforminitiative/client-registration-api/internal/config"
	"github.com/openearthplatforminitiative/client-registration-api/internal/routes"
	"log"
)

func init() {
	config.Setup()
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	router := gin.Default()
	routes.InitRoutes(router)

	log.Println("Starting server on", config.AppSettings.GetServerBindAddress())
	err := router.Run(config.AppSettings.GetServerBindAddress())
	if err != nil {
		log.Println("Failed to start server")
		return
	}
}
