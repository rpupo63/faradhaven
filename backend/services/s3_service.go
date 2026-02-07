package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type S3Service struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Service() (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		log.Warn().Msg("AWS_S3_BUCKET environment variable is not set")
	}

	return &S3Service{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		region: cfg.Region,
	}, nil
}

func (s *S3Service) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if s.bucket == "" {
		return "", fmt.Errorf("AWS_S3_BUCKET not set")
	}

	ext := filepath.Ext(fileHeader.Filename)
	key := fmt.Sprintf("characters/%s%s", uuid.New().String(), ext)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Construct URL
	// If region is empty, it might be a problem for constructing the URL accurately
	// default to us-east-1 if needed or rely on config
	region := s.region
	if region == "" {
		region = "us-east-1"
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, region, key)
	return url, nil
}
