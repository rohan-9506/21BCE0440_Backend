package services

// import (
// 	"bytes"
// 	"fmt"
// 	"mime/multipart"
// 	"net/http"
// 	"os"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// )

// // UploadFileToS3 uploads a file to the S3 bucket
// func UploadFileToS3(file multipart.File, fileName string) (string, int64, error) {
// 	// Load AWS credentials and region from environment variables
// 	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
// 	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
// 	region := os.Getenv("AWS_REGION")
// 	bucket := os.Getenv("S3_BUCKET_NAME")

// 	// Ensure that necessary environment variables are present
// 	if accessKey == "" || secretKey == "" || region == "" || bucket == "" {
// 		return "", 0, fmt.Errorf("missing AWS configuration in environment variables")
// 	}

// 	// Load the AWS session with credentials
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(region),
// 		Credentials: credentials.NewStaticCredentials(
// 			accessKey, // AWS Access Key
// 			secretKey, // AWS Secret Key
// 			""),       // Token (not needed for now)
// 	})
// 	if err != nil {
// 		return "", 0, fmt.Errorf("failed to create AWS session: %v", err)
// 	}

// 	// Create an S3 service client
// 	svc := s3.New(sess)

// 	// Read file content into a buffer
// 	buffer := new(bytes.Buffer)
// 	fileSize, err := buffer.ReadFrom(file)
// 	if err != nil {
// 		return "", 0, fmt.Errorf("failed to read file: %v", err)
// 	}

// 	// Set up the S3 upload parameters
// 	params := &s3.PutObjectInput{
// 		Bucket:      aws.String(bucket),   // S3 Bucket Name from the environment
// 		Key:         aws.String(fileName), // The file name
// 		Body:        bytes.NewReader(buffer.Bytes()),
// 		ContentType: aws.String(http.DetectContentType(buffer.Bytes())),
// 	}

// 	// Upload the file to S3
// 	_, err = svc.PutObject(params)
// 	if err != nil {
// 		return "", 0, fmt.Errorf("failed to upload file to S3: %v", err)
// 	}

// 	// Construct the S3 file URL
// 	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, fileName)

// 	return fileURL, fileSize, nil
// }
