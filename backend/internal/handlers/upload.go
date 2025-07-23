package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadResponse represents upload response
type UploadResponse struct {
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// UploadImage handles image upload
func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}
	defer file.Close()

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file too large. Maximum size is 5MB"})
		return
	}

	// Validate file type
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
	contentType := header.Header.Get("Content-Type")
	isAllowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format. Allowed formats: JPEG, PNG, GIF, WebP"})
		return
	}

	// Create upload directory if it doesn't exist
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Create file
	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Return response
	response := UploadResponse{
		URL:        fmt.Sprintf("/uploads/images/%s", filename),
		Filename:   filename,
		Size:       header.Size,
		MimeType:   contentType,
		UploadedAt: time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"data":    response,
	})
}

// UploadFile handles general file upload
func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 10MB"})
		return
	}

	// Validate file type
	allowedExtensions := []string{".pdf", ".doc", ".docx", ".txt", ".zip", ".rar", ".mp4", ".mp3"}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	isAllowed := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Create upload directory if it doesn't exist
	uploadDir := "uploads/files"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Create file
	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Return response
	response := UploadResponse{
		URL:        fmt.Sprintf("/uploads/files/%s", filename),
		Filename:   filename,
		Size:       header.Size,
		MimeType:   header.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"data":    response,
	})
}

// DeleteFile handles file deletion
func DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Validate filename to prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	// Try to delete from images directory first
	imagePath := filepath.Join("uploads/images", filename)
	if _, err := os.Stat(imagePath); err == nil {
		if err := os.Remove(imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Image file deleted successfully"})
		return
	}

	// Try to delete from files directory
	filePath := filepath.Join("uploads/files", filename)
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
}

// GetUploadStats handles getting upload statistics
func GetUploadStats(c *gin.Context) {
	// Count files in images directory
	imageCount := 0
	imageSize := int64(0)
	if entries, err := os.ReadDir("uploads/images"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				imageCount++
				if info, err := entry.Info(); err == nil {
					imageSize += info.Size()
				}
			}
		}
	}

	// Count files in files directory
	fileCount := 0
	fileSize := int64(0)
	if entries, err := os.ReadDir("uploads/files"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				fileCount++
				if info, err := entry.Info(); err == nil {
					fileSize += info.Size()
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"images": gin.H{
			"count": imageCount,
			"size":  imageSize,
		},
		"files": gin.H{
			"count": fileCount,
			"size":  fileSize,
		},
		"total": gin.H{
			"count": imageCount + fileCount,
			"size":  imageSize + fileSize,
		},
	})
} 