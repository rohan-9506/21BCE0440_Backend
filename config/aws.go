package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3Client *s3.S3

func InitS3() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Adjust to your region
	}))
	S3Client = s3.New(sess)
}
