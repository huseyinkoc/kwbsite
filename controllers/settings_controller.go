package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSettingsHandler(c *gin.Context) {
	settings, err := services.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func UpdateSettingsHandler(c *gin.Context) {
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBy := c.GetString("username") // Kullanıcı bilgisi JWT'den alınabilir
	if err := services.UpdateSettings(update, updatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

func GetSocialMediaLinksHandler(c *gin.Context) {
	links, err := services.GetSocialMediaLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve social media links"})
		return
	}

	c.JSON(http.StatusOK, links)
}

func UpdateSocialMediaLinksHandler(c *gin.Context) {
	var links map[string]models.SocialMedia
	if err := c.ShouldBindJSON(&links); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBy := c.GetString("username") // Kullanıcı bilgisi
	if err := services.UpdateSocialMediaLinks(links, updatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media links"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Social media links updated successfully"})
}
