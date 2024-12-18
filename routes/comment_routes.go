package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterCommentRoutes tüm yorum rotalarını ayarlar
func RegisterCommentRoutes(router *gin.Engine) {

	commentGroup := router.Group("/admin/comments")
	commentGroup.Use(middlewares.AuthMiddleware())
	commentGroup.Use(middlewares.AuthorizeRolesMiddleware("admin", "editor"))
	{
		commentGroup.POST("/", controllers.CreateCommentHandler)
		commentGroup.GET("/post/:postID", controllers.GetCommentsByPostIDHandler)
		commentGroup.POST("/:commentID/reply", controllers.AddReplyHandler)
		commentGroup.POST("/:commentID/like", controllers.LikeCommentHandler)
	}
}
