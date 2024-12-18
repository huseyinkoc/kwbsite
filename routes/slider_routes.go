package routes

import (
	"admin-panel/controllers"
	"admin-panel/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterSliderRoutes(router *gin.Engine) {
	sliders := router.Group("/admin/sliders")
	sliders.Use(middlewares.AuthMiddleware())
	sliders.Use(middlewares.AuthorizeRolesMiddleware("admin", "editor"))
	{
		sliders.GET("/", controllers.GetSliders)         // Slider listeleme
		sliders.POST("/", controllers.CreateSlider)      // Yeni slider ekleme
		sliders.PUT("/:id", controllers.UpdateSlider)    // Slider güncelleme
		sliders.DELETE("/:id", controllers.DeleteSlider) // Slider silme
	}
}
