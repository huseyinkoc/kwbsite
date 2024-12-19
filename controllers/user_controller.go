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

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Şifre kontrolü
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	// Rolleri kontrol et
	if len(user.Roles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Roles cannot be empty"})
		return
	}

	// Şifreyi hashle
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// FullName oluştur
	user.FullName = fmt.Sprintf("%s %s", user.Name, user.Surname)

	// Veritabanına ekle
	_, err = services.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
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

	// Şifre güncelleniyorsa hashle
	if password, ok := update["password"].(string); ok && password != "" {
		hashedPassword, err := services.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		update["password"] = hashedPassword
	} else {
		delete(update, "password")
	}

	// Roller güncelleniyorsa kontrol et
	if roles, ok := update["roles"]; ok {
		userRole, _ := c.Get("role")
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to change roles"})
			return
		}
		update["roles"] = roles
	}

	// FullName güncelleniyorsa
	if name, ok := update["name"].(string); ok {
		update["name"] = name
	}
	if surname, ok := update["surname"].(string); ok {
		update["surname"] = surname
	}
	if name, nameOk := update["name"].(string); nameOk {
		if surname, surnameOk := update["surname"].(string); surnameOk {
			update["full_name"] = fmt.Sprintf("%s %s", name, surname)
		}
	}

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
