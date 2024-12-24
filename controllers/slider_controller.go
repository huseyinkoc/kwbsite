package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSliderHandler(c *gin.Context) {
	var slider models.Slider
	if err := c.ShouldBindJSON(&slider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateSlider(&slider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create slider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slider created successfully"})
}

func GetSlidersHandler(c *gin.Context) {
	sliders, err := services.GetSliders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sliders"})
		return
	}

	c.JSON(http.StatusOK, sliders)
}

func UpdateSliderHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slider ID"})
		return
	}

	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateSlider(objectID, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update slider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slider updated successfully"})
}

func DeleteSliderHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slider ID"})
		return
	}

	if err := services.DeleteSlider(objectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete slider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slider deleted successfully"})
}
