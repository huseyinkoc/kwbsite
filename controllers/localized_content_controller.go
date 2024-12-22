package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateLocalizedContentHandler handles creating new localized content
func CreateLocalizedContentHandler(c *gin.Context) {
	var input models.LocalizedContent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateLocalizedContent(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content created successfully"})
}

// GetLocalizedContentHandler handles retrieving localized content by ID and language
func GetLocalizedContentHandler(c *gin.Context) {
	id := c.Param("id")
	lang := c.DefaultQuery("lang", "en")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	content, err := services.GetLocalizedContent(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch content"})
		return
	}

	translation := content.Translations[lang]
	c.JSON(http.StatusOK, gin.H{"content": translation})
}
