package controllers

import (
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ToggleMaintenanceMode(c *gin.Context) {
	var input struct {
		Enable  bool              `json:"enable"`
		Message map[string]string `json:"message"` // Çok dilli mesaj
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Güncellemeleri settings üzerinden yap
	update := map[string]interface{}{
		"maintenance_mode": input.Enable,
		"maintenance_msg":  input.Message,
	}

	updatedBy := c.GetString("username") // Kullanıcı bilgisi
	if err := services.UpdateSettings(update, updatedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance mode"})
		return
	}

	status := "disabled"
	if input.Enable {
		status = "enabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Maintenance mode " + status,
	})
}
