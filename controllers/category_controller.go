package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"admin-panel/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCategoryHandler(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Slug kontrolü
	if category.Slug == nil {
		category.Slug = make(map[string]string) // Slug map'ini başlat
	}

	for lang, localization := range category.Localizations {
		if localization.Title != "" {
			category.Slug[lang] = utils.GenerateSlug(localization.Title) // Her dil için slug oluştur
		} else if lang == "en" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "English title is required for slug generation"})
			return
		}
	}

	_, err := services.CreateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category created successfully"})
}

func GetAllCategoriesHandler(c *gin.Context) {
	categories, err := services.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
