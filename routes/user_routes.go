package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	users := router.Group("/admin/users")
	users.Use(middlewares.AuthMiddleware())        // JWT Middleware
	users.Use(middlewares.AuthorizeRoles("admin")) // Sadece adminler erişebilir
	{
		users.POST("/create", controllers.CreateUserHandler)
		users.GET("/", controllers.GetAllUsersHandler)
		users.PUT("/:id", controllers.UpdateUserHandler)
		users.DELETE("/:id", controllers.DeleteUserHandler)
	}
}
