package api

import (
	"bytes"
	"context"
	"file-sharing-system/models"
	"file-sharing-system/services"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

var s3Client *s3.Client

func init() {
	awsCfg, err := loadAWSCredentials()
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS credentials: %v", err))
	}
	s3Client = s3.NewFromConfig(awsCfg)
}

func loadAWSCredentials() (aws.Config, error) {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	if awsAccessKeyID == "" || awsSecretAccessKey == "" || awsRegion == "" {
		return aws.Config{}, fmt.Errorf("AWS credentials or region not set in environment")
	}

	return aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, ""),
		Region:      awsRegion,
	}, nil
}

func UploadHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// Read file content into a buffer
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), "uploaded-file")
	s3Bucket := os.Getenv("S3_BUCKET_NAME")

	// Upload the file to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s3Bucket,
		Key:    &fileName,
		Body:   buf,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file to S3: %v", err)})
		return
	}

	// Save metadata to PostgreSQL
	fileMetadata := models.File{
		Name:       fileName,
		UploadDate: time.Now(),
		Size:       int64(buf.Len()), // Convert int to int64
		S3URL:      fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3Bucket, os.Getenv("AWS_REGION"), fileName),
	}
	if err := services.SaveFileMetadata(fileMetadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file metadata: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_url": fileMetadata.S3URL})
}
