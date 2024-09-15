package models

import (
	"time"

	"gorm.io/gorm"
)

// File represents the file metadata in the database
type File struct {
	gorm.Model
	Name       string    `json:"name"`
	UploadDate time.Time `json:"upload_date"`
	Size       int64     `json:"size"` // Changed to int64 to match the size of the file
	S3URL      string    `json:"s3_url"`
}

// SaveFileMetadata saves file metadata to the database
func SaveFileMetadata(file File) error {
	db := GetDB()
	return db.Create(&file).Error
}
