package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePageHandler handles creating a new page
func CreatePageHandler(c *gin.Context) {
	var page models.Page

	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := services.CreatePage(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page created successfully"})
}

// GetAllPagesHandler handles retrieving all pages
func GetAllPagesHandler(c *gin.Context) {
	pages, err := services.GetAllPages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pages"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// UpdatePageHandler handles updating an existing page
func UpdatePageHandler(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = services.UpdatePage(id, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page updated successfully"})
}

// DeletePageHandler handles deleting a page
func DeletePageHandler(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	_, err = services.DeletePage(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
