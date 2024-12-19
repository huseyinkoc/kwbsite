package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRoleHandler handles role creation
func CreateRoleHandler(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kimin oluşturduğunu al
	createdBy, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	role.CreatedBy = createdBy.(string)
	role.UpdatedBy = createdBy.(string)
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	// Veritabanına ekle
	_, err := services.CreateRole(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role created successfully", "role": role})
}

// GetAllRolesHandler retrieves all roles
func GetAllRolesHandler(c *gin.Context) {
	roles, err := services.GetAllRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve roles"})
		fmt.Println("Database error during role retrieval:", err)
		return
	}

	c.JSON(http.StatusOK, roles)
}

// UpdateRoleHandler handles role updates
func UpdateRoleHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedBy, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	update["updated_by"] = updatedBy.(string)
	update["updated_at"] = time.Now()

	updatedCount, err := services.UpdateRole(c.Request.Context(), objectID, bson.M(update))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	if updatedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// DeleteRoleHandler handles role deletion
func DeleteRoleHandler(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		fmt.Println("Database error during role deletion:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
