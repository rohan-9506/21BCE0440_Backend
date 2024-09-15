package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"file-sharing-system/models"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func StartFileCleanupWorker(s3Client *s3.Client) {
	ticker := time.NewTicker(24 * time.Hour) // Run every 24 hours
	for {
		select {
		case <-ticker.C:
			deleteExpiredFiles(s3Client)
		}
	}
}

func deleteExpiredFiles(s3Client *s3.Client) {
	// Example: delete files older than 30 days
	expirationDate := time.Now().AddDate(0, 0, -30)

	var files []models.File
	if err := models.GetDB().Where("upload_date < ?", expirationDate).Find(&files).Error; err != nil {
		log.Printf("Failed to retrieve files for cleanup: %v", err)
		return
	}

	for _, file := range files {
		_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: &file.S3URL, // Update with correct bucket
			Key:    &file.Name,
		})
		if err != nil {
			log.Printf("Failed to delete file from S3: %v", err)
		}

		if err := models.GetDB().Delete(&file).Error; err != nil {
			log.Printf("Failed to delete file metadata from DB: %v", err)
		}
	}

	fmt.Println("Expired files deleted successfully")
}
