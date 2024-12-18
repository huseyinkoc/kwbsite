package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTagHandler(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := services.CreateTag(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag created successfully"})
}

func GetAllTagsHandler(c *gin.Context) {
	tags, err := services.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}
