package controllers

import (
	"context"
	"net/http"
	"time"

	"admin-panel/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Sliderları Listeleme
func GetSliders(c *gin.Context) {
	lang := c.Query("lang")
	var sliders []models.Slider
	collection := database.MongoDB.Collection("sliders")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var slider models.Slider
		cursor.Decode(&slider)

		if translation, exists := slider.Translations[lang]; exists {
			slider.Translations = map[string]struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			}{
				lang: translation,
			}
		}
		sliders = append(sliders, slider)
	}
	c.JSON(http.StatusOK, sliders)
}

// Yeni Slider Ekleme
func CreateSlider(c *gin.Context) {
	var slider models.Slider
	if err := c.ShouldBindJSON(&slider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slider.CreatedAt = time.Now()
	slider.UpdatedAt = time.Now()
	collection := database.MongoDB.Collection("sliders")

	_, err := collection.InsertOne(context.TODO(), slider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Slider created successfully"})
}

// Slider Güncelleme
func UpdateSlider(c *gin.Context) {
	id := c.Param("id")
	var updatedData models.Slider

	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedData.UpdatedAt = time.Now()

	collection := database.MongoDB.Collection("sliders")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedData}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Slider updated successfully"})
}

// Slider Silme
func DeleteSlider(c *gin.Context) {
	id := c.Param("id")

	collection := database.MongoDB.Collection("sliders")
	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Slider deleted successfully"})
}
