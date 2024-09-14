package api

import (
	"file-sharing-system/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file uploads
func UploadHandler(c *gin.Context) {
	// Parse the form to retrieve the file
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	// Create a channel to handle concurrency
	done := make(chan struct {
		url  string
		size int64
		err  error
	}, 1)

	// Handle file upload concurrently
	go func() {
		url, size, err := services.UploadFileToS3(file, fileHeader.Filename)
		done <- struct {
			url  string
			size int64
			err  error
		}{url, size, err}
	}()

	// Wait for the upload to complete
	result := <-done
	if result.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Save file metadata in PostgreSQL
	fileMetadata := services.FileMetadata{
		FileName:   fileHeader.Filename,
		UploadDate: time.Now(),
		Size:       result.size,
		S3URL:      result.url,
	}

	if err := services.SaveFileMetadata(fileMetadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": result.url})
}
