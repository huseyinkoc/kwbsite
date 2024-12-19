package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func TagRoutes(router *gin.Engine) {
	tags := router.Group("/admin/tags")
	tags.Use(middlewares.AuthMiddleware())
	tags.Use(middlewares.AuthorizeRolesMiddleware("admin", "editor"))
	{
		tags.POST("/create", middlewares.CSRFMiddleware(), controllers.CreateTagHandler)
		tags.GET("/", controllers.GetAllTagsHandler)
	}
}
