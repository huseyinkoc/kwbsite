package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func PageRoutes(router *gin.Engine) {
	pages := router.Group("/admin/pages")
	pages.Use(middlewares.AuthMiddleware())                  // JWT kontrolü
	pages.Use(middlewares.AuthorizeRoles("admin", "editor")) // Roller
	{
		pages.POST("/create", controllers.CreatePageHandler)
		pages.GET("/", controllers.GetAllPagesHandler)
		pages.PUT("/:id", controllers.UpdatePageHandler)
		pages.DELETE("/:id", controllers.DeletePageHandler)
	}
}
