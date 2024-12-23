package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePluginHandler(c *gin.Context) {
	var plugin models.Plugin
	if err := c.ShouldBindJSON(&plugin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreatePlugin(&plugin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plugin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin created successfully"})
}

func GetAllPluginsHandler(c *gin.Context) {
	plugins, err := services.GetAllPlugins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plugins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plugins": plugins})
}

func UpdatePluginHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plugin ID"})
		return
	}

	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdatePlugin(objectID, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plugin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin updated successfully"})
}

func DeletePluginHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plugin ID"})
		return
	}

	if err := services.DeletePlugin(objectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plugin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin deleted successfully"})
}

func UploadPluginHandler(c *gin.Context) {
	// Dosya yükleme
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file"})
		return
	}

	// Dosya adını kontrol et
	if !strings.HasSuffix(file.Filename, ".so") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format. Only .so files are allowed"})
		return
	}

	// Yükleme dizini
	uploadPath := "./plugins/"
	if err := c.SaveUploadedFile(file, uploadPath+file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	pluginPath := uploadPath + file.Filename
	expectedHash := "expected_sha256_hash_value"
	isValid, err := services.VerifyPluginIntegrity(pluginPath, expectedHash)
	if err != nil || !isValid {
		log.Println("Plugin integrity check failed")
		return
	}
	log.Println("Plugin is valid")

	c.JSON(http.StatusOK, gin.H{"message": "Plugin uploaded successfully", "file": file.Filename})
}

func ListUploadedPluginsHandler(c *gin.Context) {
	plugins, err := services.ListUploadedPlugins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list plugins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plugins": plugins})
}

func EnablePluginHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plugin ID"})
		return
	}

	if err := services.EnablePlugin(objectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable plugin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plugin enabled successfully"})
}
