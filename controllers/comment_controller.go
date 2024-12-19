package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"log"
	"net/http"
	"strconv"
	"time"

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
/*
func GetCommentsByPostIDHandler(c *gin.Context) {
	postID := c.Param("postID")

	// PostID'yi ObjectID'ye çevir
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	// Yorumları getir
	comments, err := services.GetCommentsByPostID(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments", "details": err.Error()})
		return
	}

	// Cevapları yorumlara bağla
	commentMap := make(map[primitive.ObjectID]*models.Comment)
	var rootComments []models.Comment

	for i := range comments {
		comment := &comments[i]
		commentMap[comment.ID] = comment

		// Ana yorumları ve cevapları ayır
		if comment.ParentID == nil {
			// Ana yorum
			rootComments = append(rootComments, *comment)
		} else {
			// Cevapları ana yoruma ekle
			parent, exists := commentMap[*comment.ParentID]
			if exists {
				parent.Replies = append(parent.Replies, comment.ID)
			}
		}
	}

	// Reaksiyonları ve cevap detaylarını döndür
	response := make([]map[string]interface{}, 0)
	for _, root := range rootComments {
		// Ana yorum detaylarını ekle
		commentData := map[string]interface{}{
			"id":        root.ID.Hex(),
			"content":   root.Content,
			"reactions": root.Reactions,
			"replies":   []map[string]interface{}{},
		}

		// Cevapların detaylarını ekle
		for _, replyID := range root.Replies {
			if reply, exists := commentMap[replyID]; exists {
				replyData := map[string]interface{}{
					"id":        reply.ID.Hex(),
					"content":   reply.Content,
					"reactions": reply.Reactions,
				}
				commentData["replies"] = append(commentData["replies"].([]map[string]interface{}), replyData)
			}
		}

		response = append(response, commentData)
	}

	c.JSON(http.StatusOK, gin.H{"comments": response})
}
*/

func GetCommentsByPostIDHandler(c *gin.Context) {
	postID := c.Param("postID")

	// PostID'yi ObjectID'ye çevir
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	// Sayfa ve limit parametrelerini al
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Sayfalama için başlangıç ve limit hesapla
	skip := (page - 1) * limit

	// Yorumları getir
	comments, err := services.GetCommentsByPostIDWithPagination(c.Request.Context(), objectID, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit, "comments": comments})
}

// AddReplyHandler bir yoruma cevap ekler
func AddReplyHandler(c *gin.Context) {
	commentID := c.Param("commentID")

	// CommentID'yi ObjectID'ye çevir
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id"})
		return
	}

	var reply models.Comment
	if err := c.ShouldBindJSON(&reply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	reply.ID = primitive.NewObjectID()
	reply.ParentID = &objectID // ParentID olarak ana yorumun ID'sini belirle
	reply.CreatedAt = time.Now()
	reply.UpdatedAt = time.Now()

	// Cevap yorumunu oluştur
	replyResult, err := services.CreateComment(c.Request.Context(), &reply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply", "details": err.Error()})
		return
	}

	// InsertedID'yi al ve logla veya başka bir işlemde kullan
	replyID, ok := replyResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("InsertedID is not a valid ObjectID: %v", replyResult.InsertedID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse inserted ID"})
		return
	}

	// Ana yoruma bu cevabı ekle
	err = services.AddReply(c.Request.Context(), objectID, replyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reply", "details": err.Error()})
		return
	}

	// Ana yorumu bul ve bildirim oluştur
	parentComment, err := services.FetchCommentByID(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch parent comment", "details": err.Error()})
		return
	}

	err = services.CreateNotification(c.Request.Context(), parentComment.UserID, "Yorumunuza bir yanıt geldi.")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply added successfully", "reply_id": replyID})
}

func AddReactionHandler(c *gin.Context) {
	commentID := c.Param("commentID")
	reaction := c.Query("reaction") // İfade parametresi (örneğin: 😊, 😡, ❤️)

	// CommentID'yi ObjectID'ye çevir
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id"})
		return
	}

	// Reaction kontrolü
	if reaction == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reaction is required"})
		return
	}

	// Reaksiyon ekle
	err = services.AddReaction(c.Request.Context(), objectID, reaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction added successfully"})
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

func DeleteCommentHandler(c *gin.Context) {
	commentID := c.Param("commentID")

	// CommentID'yi ObjectID'ye çevir
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id"})
		return
	}

	// Yorum sil
	err = services.DeleteComment(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func UpdateCommentHandler(c *gin.Context) {
	commentID := c.Param("commentID")

	// CommentID'yi ObjectID'ye çevir
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id"})
		return
	}

	var updatedData struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Yorum güncelle
	err = services.UpdateComment(c.Request.Context(), objectID, updatedData.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}
