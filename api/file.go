package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"file-sharing-system/models"
	"file-sharing-system/services"
	"file-sharing-system/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var s3Client *s3.Client
var redisClient *redis.Client

func init() {
	awsCfg, err := loadAWSCredentials()
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS credentials: %v", err))
	}
	s3Client = s3.NewFromConfig(awsCfg)
	InitRedis()
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

func InitRedis() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatalf("REDIS_ADDR not set in environment")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
}

func UploadHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), "uploaded-file")
	s3Bucket := os.Getenv("S3_BUCKET_NAME")

	// Encrypt the file content
	encryptedData, err := utils.Encrypt(buf.Bytes())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt file"})
		return
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s3Bucket,
		Key:    &fileName,
		Body:   bytes.NewReader(encryptedData),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file to S3: %v", err)})
		return
	}

	fileMetadata := models.File{
		Name:       fileName,
		UploadDate: time.Now(),
		Size:       int64(buf.Len()),
		S3URL:      fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3Bucket, os.Getenv("AWS_REGION"), fileName),
	}
	if err := services.SaveFileMetadata(fileMetadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file metadata: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_url": fileMetadata.S3URL})
}

func GetFilesHandler(c *gin.Context) {
	userIDStr := c.Query("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cacheKey := "files:" + strconv.Itoa(userID)
	cachedFilesJSON, err := redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedFilesJSON != "" {
		var cachedFiles []models.File
		if err := json.Unmarshal([]byte(cachedFilesJSON), &cachedFiles); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal cached files"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"files": cachedFiles})
		return
	}

	var files []models.File
	if err := models.GetDB().Where("user_id = ?", userID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	filesJSON, err := json.Marshal(files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal files"})
		return
	}

	redisClient.Set(context.Background(), cacheKey, filesJSON, 5*time.Minute) // Cache with 5 minutes expiry

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func UpdateFileMetadataHandler(c *gin.Context) {
	fileID := c.Param("file_id")
	var updatedFile models.File
	if err := c.BindJSON(&updatedFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Assuming file metadata can be updated, and no need for UserID in this context
	if err := models.GetDB().Model(&models.File{}).Where("id = ?", fileID).Updates(updatedFile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file metadata"})
		return
	}

	// Invalidate cache
	cacheKey := "files:" + strconv.Itoa(int(updatedFile.ID))
	redisClient.Del(context.Background(), cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "File metadata updated"})
}

func ShareFileHandler(c *gin.Context) {
	fileID := c.Param("file_id")

	var file models.File
	if err := models.GetDB().Where("id = ?", fileID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("S3_BUCKET_NAME"), os.Getenv("AWS_REGION"), file.S3URL)

	c.JSON(http.StatusOK, gin.H{"public_url": publicURL})
}

func SearchFilesHandler(c *gin.Context) {
	name := c.Query("name")
	uploadDateStr := c.Query("upload_date")

	var uploadDate time.Time
	if uploadDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", uploadDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload date format"})
			return
		}
		uploadDate = parsedDate
	}

	files, err := models.SearchFiles(name, uploadDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
