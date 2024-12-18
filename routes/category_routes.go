package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categories := router.Group("/admin/categories")
	categories.Use(middlewares.AuthMiddleware())
	categories.Use(middlewares.AuthorizeRoles("admin", "editor"))
	{
		categories.POST("/create", controllers.CreateCategoryHandler)
		categories.GET("/", controllers.GetAllCategoriesHandler)
	}
}
