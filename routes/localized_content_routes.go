package routes

import (
	"admin-panel/controllers"

	"github.com/gin-gonic/gin"
)

func LocalizedContentRoutes(router *gin.Engine) {
	content := router.Group("/content")
	{
		content.POST("/", controllers.CreateLocalizedContentHandler)
		content.GET("/:id", controllers.GetLocalizedContentHandler)
	}
}
