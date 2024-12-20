package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateContactMessageHandler(c *gin.Context) {
	var message models.ContactMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdMessage, err := services.CreateContactMessage(c.Request.Context(), &message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save contact message", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact message sent successfully", "data": createdMessage})
}

func GetAllContactMessagesHandler(c *gin.Context) {
	messages, err := services.GetAllContactMessages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact messages", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

func UpdateContactMessageStatusHandler(c *gin.Context) {
	id := c.Param("id")
	status := c.Query("status")

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	resolvedBy, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := services.UpdateContactMessageStatus(c.Request.Context(), id, status, resolvedBy.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact message status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact message status updated successfully"})
}

func DeleteContactMessageHandler(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteContactMessage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact message", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact message deleted successfully"})
}
