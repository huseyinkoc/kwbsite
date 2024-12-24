package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func SettingsRoutes(router *gin.Engine) {
	settings := router.Group("/settings")
	settings.Use(middlewares.MaintenanceMiddleware()) // Bakım modu kontrolü
	settings.Use(middlewares.AuthMiddleware())        // Yetkilendirme
	{
		settings.GET("/", controllers.GetSettingsHandler)
		settings.PUT("/", controllers.UpdateSettingsHandler)
	}
}
