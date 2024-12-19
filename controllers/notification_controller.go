package controllers

import (
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetNotificationsHandler kullanıcıya ait bildirimleri döndürür
func GetNotificationsHandler(c *gin.Context) {
	// Kullanıcı kimliği (örneğin, bir middleware tarafından ayarlanmış olabilir)
	userID := c.GetString("userID") // Middleware ile ayarlanmalı
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Kullanıcı ID'sini ObjectID'ye çevir
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Bildirimleri getir
	notifications, err := services.FetchNotificationsByUserID(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}
