package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openearthplatforminitiative/client-registration-api/config"
	"github.com/openearthplatforminitiative/client-registration-api/handlers"
	"github.com/openearthplatforminitiative/client-registration-api/middleware"
)

func InitRoutes(router *gin.Engine, client *config.Config) {
	router.GET("/health", handlers.Health)
	router.GET("/ready", handlers.Ready)
	router.GET("/dump", handlers.Dump)

	cc := &handlers.ClientsHandler{Keycloak: config.AppSettings.GetKeycloakClient()}

	api := router.Group("/clients", middleware.UserRequired(), middleware.ErrorHandler())
	api.GET("/", cc.Clients)
	api.GET("/:id", cc.Client)
	api.POST("/", cc.AddClient)
	api.PUT("/:id", cc.UpdateClient)
	api.DELETE("/:id", cc.DeleteClient)
}
