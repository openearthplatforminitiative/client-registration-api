package main

import (
	"github.com/gin-gonic/gin"
	"github.com/openearthplatforminitiative/client-registration-api/config"
	"github.com/openearthplatforminitiative/client-registration-api/routes"
	"log"
)

func init() {
	config.Setup()
	gin.SetMode(gin.ReleaseMode)
}

func main() {

	router := gin.Default()
	routes.InitRoutes(router, config.AppSettings)

	log.Println("Starting server on", config.AppSettings.GetServerBindAddress())
	err := router.Run(config.AppSettings.GetServerBindAddress())
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		return
	}
}
