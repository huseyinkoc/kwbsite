package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {
	posts := router.Group("/admin/posts")
	posts.Use(middlewares.AuthMiddleware())
	posts.Use(middlewares.AuthorizeRolesMiddleware("admin", "editor"))
	{
		posts.POST("/create", controllers.CreatePostHandler)
		posts.GET("/", controllers.GetAllPostsHandler)
		posts.GET("/filter", controllers.GetFilteredPostsHandler) // Filtrelenmiş postlar
	}
}
