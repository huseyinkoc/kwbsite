package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCommentHandler yeni bir yorum oluşturur
func CreateCommentHandler(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.ID = primitive.NewObjectID()

	// Yorumu oluştur
	result, err := services.CreateComment(c.Request.Context(), &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment created", "result": result})
}

// GetCommentsByPostIDHandler bir gönderiye ait yorumları döndürür
func GetCommentsByPostIDHandler(c *gin.Context) {
	postID := c.Param("postID")

	// PostID'nin ObjectID'ye dönüşümünü kontrol et
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Post ID"})
		return
	}

	// Yorumları getir
	comments, err := services.GetCommentsByPostID(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// AddReplyHandler bir yoruma cevap ekler
func AddReplyHandler(c *gin.Context) {
	commentID := c.Param("commentID")

	// CommentID'nin ObjectID'ye dönüşümünü kontrol et
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Comment ID"})
		return
	}

	var reply models.Comment
	if err := c.ShouldBindJSON(&reply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply.ID = primitive.NewObjectID()

	// Cevap yorumu oluştur
	replyResult, err := services.CreateComment(c.Request.Context(), &reply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cevap olarak ekle
	err = services.AddReply(c.Request.Context(), objectID, replyResult.InsertedID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply added", "result": replyResult})
}

// LikeCommentHandler bir yorumu beğenmek için kullanılır
func LikeCommentHandler(c *gin.Context) {
	commentID := c.Param("commentID")

	// CommentID'nin ObjectID'ye dönüşümünü kontrol et
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Comment ID"})
		return
	}

	// Yorumu beğen
	if err := services.LikeComment(c.Request.Context(), objectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment liked"})
}
