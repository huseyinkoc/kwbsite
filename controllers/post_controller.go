package controllers

import (
	"admin-panel/models"   // Post modelini import ettik
	"admin-panel/services" // Post servislerini import ettik
	"admin-panel/utils"
	"net/http"
	"time"

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
		if localization, ok := post.Localizations["en"]; ok {
			post.Slug = utils.GenerateSlug(localization.Title) // İngilizce başlık üzerinden slug oluşturuluyor
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "English title is required for slug generation"})
			return
		}
	}

	// Varsayılan durum
	if post.Status == "" {
		post.Status = "draft"
	}

	// Zaman damgaları
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	// Veritabanına kaydet
	result, err := services.CreatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, result)
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

	// Tüm gönderileri alın
	posts, err := services.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	// Dil bazlı gönderileri filtreleyin
	localizedPosts := []map[string]interface{}{}
	for _, post := range posts {
		// Belirtilen dilde içerik var mı kontrol edin
		if localizedContent, ok := post.Localizations[lang]; ok {
			localizedPosts = append(localizedPosts, map[string]interface{}{
				"id":         post.ID.Hex(), // MongoDB ObjectID'yi stringe çevirin
				"slug":       post.Slug,
				"title":      localizedContent.Title,
				"content":    localizedContent.Content,
				"status":     post.Status,
				"categories": post.CategoryIDs, // Kategorileri döndürün
			})
		}
	}

	// Eğer hiçbir içerik yoksa uygun bir yanıt döndürün
	if len(localizedPosts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No posts found for the specified language"})
		return
	}

	c.JSON(http.StatusOK, localizedPosts)
}
