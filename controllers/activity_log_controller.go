package controllers

import (
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetActivityLogsHandler(c *gin.Context) {
	// Filtreleme ve limit parametrelerini al
	limit := 100
	filter := map[string]interface{}{}

	logs, err := services.GetActivityLogs(filter, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
