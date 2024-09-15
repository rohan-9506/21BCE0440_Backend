package services

import (
	"file-sharing-system/models"
)

// SaveFileMetadata saves file metadata to the database
func SaveFileMetadata(file models.File) error {
	db := models.GetDB()
	return db.Create(&file).Error
}
