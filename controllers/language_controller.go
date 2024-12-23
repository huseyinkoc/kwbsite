package controllers

import (
	"admin-panel/configs"
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateLanguageHandler(c *gin.Context) {
	var language models.Language
	if err := c.ShouldBindJSON(&language); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateLanguage(&language); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language created successfully"})
}

func GetLanguagesHandler(c *gin.Context) {
	// Aktif dil parametresini al
	activeLang := c.DefaultQuery("lang", "en") // Varsayılan dil "en"

	// Tüm dilleri al
	languages, err := services.GetLanguagesWithActiveAndDefault(activeLang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve languages"})
		return
	}

	// Varsayılan ve aktif dil bilgilerini ekle
	response := gin.H{
		"active_language":  activeLang,
		"default_language": configs.LanguageConfig.DefaultLanguage, // Varsayılan dili config'den çek
		"languages":        languages,
	}

	// Yanıt döndür
	c.JSON(http.StatusOK, response)
}

func UpdateLanguageHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language ID"})
		return
	}

	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateLanguage(objectID, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated successfully"})
}

func DeleteLanguageHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language ID"})
		return
	}

	if err := services.DeleteLanguage(objectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language deleted successfully"})
}
