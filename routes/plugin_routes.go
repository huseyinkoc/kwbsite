package routes

import (
	"admin-panel/controllers"

	"github.com/gin-gonic/gin"
)

func PluginRoutes(router *gin.Engine) {
	plugins := router.Group("/plugins")
	{
		plugins.POST("/", controllers.CreatePluginHandler)
		plugins.GET("/", controllers.GetAllPluginsHandler)
		plugins.PUT("/:id", controllers.UpdatePluginHandler)
		plugins.DELETE("/:id", controllers.DeletePluginHandler)
		plugins.POST("/upload", controllers.UploadPluginHandler) // Dosya yükleme

	}
}
