package models

import (
	"time"
)

// File represents the metadata for uploaded files
type File struct {
	ID         uint      `gorm:"primaryKey"`
	FileName   string    `gorm:"not null"`
	UploadDate time.Time `gorm:"not null"`
	Size       int64     `gorm:"not null"`
	S3URL      string    `gorm:"not null"`
	UserID     uint      `gorm:"not null"` // Foreign key to associate with a user
}
