package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func ContactRoutes(router *gin.Engine) {
	contacts := router.Group("/admin/contact")
	contacts.Use(middlewares.AuthMiddleware())
	contacts.Use(middlewares.AuthorizeRolesMiddleware("admin"))
	{
		contacts.POST("/", middlewares.CSRFMiddleware(), controllers.CreateContactMessageHandler)
		contacts.GET("/", controllers.GetAllContactMessagesHandler)
		contacts.PUT("/:id", middlewares.CSRFMiddleware(), controllers.UpdateContactMessageStatusHandler)
		contacts.DELETE("/:id", middlewares.CSRFMiddleware(), controllers.DeleteContactMessageHandler)
	}
}
