package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	menus := router.Group("/admin/menus")
	menus.Use(middlewares.AuthMiddleware()) // Kullanıcı giriş kontrolü
	{
		menus.POST("/", controllers.CreateMenuHandler) // Menü oluşturma
		menus.GET("/", controllers.GetMenusHandler)    // Yetkilere göre menüleri getirme
		menus.PUT("/:id", controllers.UpdateMenuHandler)
		menus.DELETE("/:id", controllers.DeleteMenuHandler)
	}
}
