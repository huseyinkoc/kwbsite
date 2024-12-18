package controllers

import (
	"admin-panel/models"
	"admin-panel/services"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadMediaHandler handles media file uploads
func UploadMediaHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is uploaded"})
		return
	}

	// Dosyayı uploads klasörüne kaydet
	uploadPath := "uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, os.ModePerm)
	}

	filePath := filepath.Join(uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Medya kaydı ekle
	media := models.Media{
		FileName: file.Filename,
		FilePath: filePath,
		FileType: filepath.Ext(file.Filename),
	}
	_, err = services.SaveMediaRecord(media)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": filePath})
}

// GetAllMediaHandler handles listing all media files
func GetAllMediaHandler(c *gin.Context) {
	media, err := services.GetAllMedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve media files"})
		return
	}

	c.JSON(http.StatusOK, media)
}

func DeleteMediaHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	// Veritabanından medya kaydını sil
	result, err := services.DeleteMedia(id)
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	// Opsiyonel: Dosya sisteminden de sil
	filePath := "uploads/" + idParam // Örnek dosya yolu, MongoDB'den dosya yolunu da çekebilirsiniz.
	if err := os.Remove(filePath); err != nil {
		fmt.Println("Warning: File not found or failed to delete:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

func GetMediaDetailHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	media, err := services.GetMediaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	c.JSON(http.StatusOK, media)
}

func GetFilteredMediaHandler(c *gin.Context) {
	// Sorgu parametrelerini al
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Filtreleme kriterlerini oluştur
	filter := bson.M{}
	if fileName != "" {
		filter["file_name"] = bson.M{"$regex": fileName, "$options": "i"} // Case-insensitive arama
	}
	if fileType != "" {
		filter["file_type"] = fileType
	}
	if startDate != "" && endDate != "" {
		start, err1 := time.Parse("2006-01-02", startDate)
		end, err2 := time.Parse("2006-01-02", endDate)
		if err1 == nil && err2 == nil {
			filter["uploaded_at"] = bson.M{
				"$gte": start.Unix(),
				"$lte": end.Unix(),
			}
		}
	}

	// Medya dosyalarını getir
	media, err := services.GetFilteredMedia(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve media files"})
		return
	}

	c.JSON(http.StatusOK, media)
}
