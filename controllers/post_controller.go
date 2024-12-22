package controllers

import (
	"admin-panel/helpers"
	"admin-panel/models"
	"admin-panel/services"
	"admin-panel/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePostHandler handles creating a new post
func CreatePostHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hasPermission, err := helpers.HasModulePermission(c.Request.Context(), role.(string), "posts", "create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions", "details": err.Error()})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create posts"})
		return
	}

	var input struct {
		Localizations map[string]models.LocalizedField `json:"localizations" binding:"required"`
		Status        string                           `json:"status"`
		CategoryIDs   []primitive.ObjectID             `json:"category_ids"`
		TagIDs        []primitive.ObjectID             `json:"tag_ids"`
		PublishDate   *primitive.DateTime              `json:"publish_date"`
		MetaTags      map[string]models.MetaTag        `json:"meta_tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Slug oluşturma
	for lang, localization := range input.Localizations {
		if localization.Slug == "" {
			localization.Slug = utils.GenerateSlug(localization.Title)
			input.Localizations[lang] = localization
		}
	}

	// Varsayılan durum
	if input.Status == "" {
		input.Status = "draft"
	}

	// Yeni Post oluşturma
	post := models.Post{
		ID:            primitive.NewObjectID(),
		Localizations: input.Localizations,
		Status:        input.Status,
		CategoryIDs:   input.CategoryIDs,
		TagIDs:        input.TagIDs,
		PublishDate:   input.PublishDate,
		MetaTags:      input.MetaTags,
		CreatedAt:     primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
	}

	// Veritabanına kaydet
	if err := services.CreatePost(c.Request.Context(), &post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully", "post": post})
}

// GetAllPostsHandler handles retrieving all posts
func GetAllPostsHandler(c *gin.Context) {
	posts, err := services.GetAllPosts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetPostByIDHandler handles retrieving a single post by its ID
func GetPostByIDHandler(c *gin.Context) {
	id := c.Param("id")
	lang := c.DefaultQuery("lang", "en")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := services.GetPostByID(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	localizedContent, exists := post.Localizations[lang]
	if !exists {
		localizedContent = post.Localizations["en"] // Varsayılan dil
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         post.ID.Hex(),
		"slug":       localizedContent.Slug,
		"title":      localizedContent.Title,
		"content":    localizedContent.Content,
		"status":     post.Status,
		"categories": post.CategoryIDs,
		"tags":       post.TagIDs,
		"meta_tags":  post.MetaTags[lang],
	})
}

// UpdatePostHandler handles updating an existing post
func UpdatePostHandler(c *gin.Context) {
	id := c.Param("id")
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hasPermission, err := helpers.HasModulePermission(c.Request.Context(), role.(string), "posts", "update")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions", "details": err.Error()})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update posts"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var input struct {
		Localizations map[string]models.LocalizedField `json:"localizations"`
		Status        string                           `json:"status"`
		CategoryIDs   []primitive.ObjectID             `json:"category_ids"`
		TagIDs        []primitive.ObjectID             `json:"tag_ids"`
		PublishDate   *primitive.DateTime              `json:"publish_date"`
		MetaTags      map[string]models.MetaTag        `json:"meta_tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mevcut postu al
	post, err := services.GetPostByID(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	// Alanları güncelle
	if input.Localizations != nil {
		for lang, localization := range input.Localizations {
			if localization.Slug == "" {
				localization.Slug = utils.GenerateSlug(localization.Title)
			}
			post.Localizations[lang] = localization
		}
	}
	if input.Status != "" {
		post.Status = input.Status
	}
	if input.CategoryIDs != nil {
		post.CategoryIDs = input.CategoryIDs
	}
	if input.TagIDs != nil {
		post.TagIDs = input.TagIDs
	}
	if input.PublishDate != nil {
		post.PublishDate = input.PublishDate
	}
	if input.MetaTags != nil {
		post.MetaTags = input.MetaTags
	}

	post.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Veritabanında güncelle
	if err := services.UpdatePost(c.Request.Context(), post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully", "post": post})
}
