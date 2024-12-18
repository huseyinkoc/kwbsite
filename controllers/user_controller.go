package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserHandler handles user creation
func CreateUserHandler(c *gin.Context) {
	var user models.User

	// JSON verisini bind et
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("JSON binding error:", err)
		return
	}

	// Şifre alanını kontrol et
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	// Şifreyi hashle
	hashedPassword, errHP := services.HashPassword(user.Password)
	if errHP != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Kullanıcıyı veritabanına ekle
	_, err := services.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		fmt.Println("Database error:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// GetAllUsersHandler handles retrieving all users
func GetAllUsersHandler(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		fmt.Println("Database error during user retrieval:", err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUserHandler handles user updates
func UpdateUserHandler(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Şifre güncellemesi varsa hashle
	if password, ok := update["password"].(string); ok && password != "" {
		hashedPassword, errHP := services.HashPassword(password)
		if errHP != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		update["password"] = hashedPassword
	} else {
		delete(update, "password") // Boş şifre güncellemesini önle
	}

	// Rol güncellemesi varsa ve güncelleyen kullanıcı admin değilse engelle
	if _, ok := update["role"]; ok {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to change roles"})
			return
		}
	}

	// Kullanıcıyı güncelle
	_, err = services.UpdateUser(id, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUserHandler handles user deletion
func DeleteUserHandler(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		fmt.Println("Invalid ID format:", err)
		return
	}

	_, err = services.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		fmt.Println("Database error during deletion:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
