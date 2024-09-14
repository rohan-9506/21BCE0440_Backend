package services

import (
	"file-sharing-system/models"
	"time"
)

// FileMetadata represents the metadata to be stored
type FileMetadata struct {
	FileName   string
	UploadDate time.Time
	Size       int64
	S3URL      string
	UserID     uint
}

// SaveFileMetadata stores the file metadata in the database
func SaveFileMetadata(fileMetadata FileMetadata) error {
	file := models.File{
		FileName:   fileMetadata.FileName,
		UploadDate: fileMetadata.UploadDate,
		Size:       fileMetadata.Size,
		S3URL:      fileMetadata.S3URL,
		UserID:     fileMetadata.UserID, // Assuming user is already authenticated
	}
	return models.GetDB().Create(&file).Error
}
