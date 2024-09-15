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
	Size       int64     `json:"size"`
	S3URL      string    `json:"s3_url"`
}

// SaveFileMetadata saves file metadata to the database
func SaveFileMetadata(file File) error {
	db := GetDB()
	return db.Create(&file).Error
}

// SearchFiles searches for files based on metadata
func SearchFiles(name string, uploadDate time.Time) ([]File, error) {
	db := GetDB()
	var files []File

	query := db.Model(&File{})
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if !uploadDate.IsZero() {
		query = query.Where("upload_date = ?", uploadDate)
	}

	err := query.Find(&files).Error
	return files, err
}
