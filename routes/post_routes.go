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
		posts.POST("/create", middlewares.CSRFMiddleware(), middlewares.ModulePermissionMiddleware("posts", "create"), controllers.CreatePostHandler)
		posts.GET("/", middlewares.ModulePermissionMiddleware("posts", "read"), controllers.GetAllPostsHandler)
		posts.GET("/filter", controllers.GetFilteredPostsHandler)       // Filtrelenmiş postlar
		posts.GET("/lang/:lang", controllers.GetPostsByLanguageHandler) // Dil bazlı içerik listeleme
	}
}
