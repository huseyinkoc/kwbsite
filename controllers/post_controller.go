package controllers

import (
	"admin-panel/models"   // Post modelini import ettik
	"admin-panel/services" // Post servislerini import ettik
	"admin-panel/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CreatePostHandler handles creating a new post
func CreatePostHandler(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Varsayılan slug oluşturma
	if post.Slug == "" {
		post.Slug = utils.GenerateSlug(post.Localizations["en"].Title) // İngilizce başlık üzerinden slug oluşturuluyor
	}

	_, err := services.CreatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

// GetAllPostsHandler handles retrieving all posts
func GetAllPostsHandler(c *gin.Context) {
	posts, err := services.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetFilteredPostsHandler(c *gin.Context) {
	// Filtreleri al
	category := c.Query("category")
	tag := c.Query("tag")
	status := c.Query("status")

	// Filtreleme kriterlerini oluştur
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	if tag != "" {
		filter["tags"] = tag
	}
	if status != "" {
		filter["status"] = status
	}

	// Filtrelenmiş postları al
	posts, err := services.GetFilteredPosts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPostsByLanguageHandler(c *gin.Context) {
	lang := c.Param("lang") // URL parametresinden dil kodu alınır (örnek: "en", "tr")
	posts, err := services.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	localizedPosts := []map[string]interface{}{}
	for _, post := range posts {
		if localizedContent, ok := post.Localizations[lang]; ok {
			localizedPosts = append(localizedPosts, map[string]interface{}{
				"id":       post.ID,
				"slug":     post.Slug,
				"title":    localizedContent.Title,
				"content":  localizedContent.Content,
				"status":   post.Status,
				"category": post.CategoryIDs,
			})
		}
	}

	c.JSON(http.StatusOK, localizedPosts)
}
