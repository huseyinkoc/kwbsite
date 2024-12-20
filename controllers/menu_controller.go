package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateMenuHandler handles menu creation
func CreateMenuHandler(c *gin.Context) {
	var menu models.Menu

	// Bind JSON payload to menu model
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set metadata
	createdBy, exists := c.Get("username")
	if exists {
		menu.CreatedBy = createdBy.(string)
		menu.UpdatedBy = createdBy.(string)
	}

	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()

	// Create menu in database
	createdMenu, err := services.CreateMenu(c.Request.Context(), &menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu created successfully", "menu": createdMenu})
}

// GetMenusHandler retrieves menus based on type and user roles
func GetMenusHandler(c *gin.Context) {
	menuType := c.Query("type") // frontend or backend
	if menuType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menu type is required"})
		return
	}

	// Retrieve user roles from context
	roles, exists := c.Get("roles")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch menus based on roles and type
	menus, err := services.GetMenusByRoles(c.Request.Context(), menuType, roles.([]string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menus", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"menus": menus})
}

// UpdateMenuHandler updates a menu by its ID
func UpdateMenuHandler(c *gin.Context) {
	id := c.Param("id")

	var update bson.M
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set updated metadata
	updatedBy, exists := c.Get("username")
	if exists {
		update["updated_by"] = updatedBy.(string)
	}
	update["updated_at"] = time.Now()

	// Update menu in database
	updatedMenu, err := services.UpdateMenu(c.Request.Context(), id, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu updated successfully", "menu": updatedMenu})
}

// DeleteMenuHandler deletes a menu by its ID
func DeleteMenuHandler(c *gin.Context) {
	id := c.Param("id")

	// Delete menu from database
	if err := services.DeleteMenu(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}
