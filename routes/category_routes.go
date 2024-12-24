package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categories := router.Group("/admin/categories")
	categories.Use(middlewares.MaintenanceMiddleware()) // Bakım modu kontrolü
	categories.Use(middlewares.AuthMiddleware())
	//categories.Use(middlewares.AuthorizeRolesMiddleware("admin", "editor"))
	{
		categories.POST("/create", middlewares.CSRFMiddleware(), middlewares.AuthorizeRolesMiddleware("admin"), controllers.CreateCategoryHandler)
		categories.GET("/", middlewares.CSRFMiddleware(), middlewares.AuthorizeRolesMiddleware("admin", "editor"), controllers.GetAllCategoriesHandler)
	}
}
