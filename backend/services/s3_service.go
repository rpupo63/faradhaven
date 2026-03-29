package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type S3Service struct {
	client               *s3.Client
	characterImagesBucket string
	bulletinBucket        string
	region                string
}

func NewS3Service() (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	characterImagesBucket := os.Getenv("CHARACTER_IMAGES_BUCKET")
	if characterImagesBucket == "" {
		log.Warn().Msg("CHARACTER_IMAGES_BUCKET environment variable is not set")
	}
	bulletinBucket := os.Getenv("BULLETIN_BUCKET")
	if bulletinBucket == "" {
		log.Warn().Msg("BULLETIN_BUCKET environment variable is not set")
	}

	return &S3Service{
		client:                s3.NewFromConfig(cfg),
		characterImagesBucket: characterImagesBucket,
		bulletinBucket:        bulletinBucket,
		region:                cfg.Region,
	}, nil
}

func (s *S3Service) publicObjectURL(bucket, key string) string {
	region := s.region
	if region == "" {
		region = "us-east-1"
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}

func (s *S3Service) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if s.characterImagesBucket == "" {
		return "", fmt.Errorf("CHARACTER_IMAGES_BUCKET not set")
	}

	ext := filepath.Ext(fileHeader.Filename)
	key := fmt.Sprintf("characters/%s%s", uuid.New().String(), ext)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.characterImagesBucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return s.publicObjectURL(s.characterImagesBucket, key), nil
}

// UploadStream uploads an object to S3 from an io.Reader, generating a unique key.
func (s *S3Service) UploadStream(ctx context.Context, reader io.Reader, filename, contentType string) (string, error) {
	if s.characterImagesBucket == "" {
		return "", fmt.Errorf("CHARACTER_IMAGES_BUCKET not set")
	}

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("monsters/%s%s", uuid.New().String(), ext) // Use "monsters" prefix

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.characterImagesBucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload stream to S3: %v", err)
	}

	return s.publicObjectURL(s.characterImagesBucket, key), nil
}

// UploadBulletinPDF uploads a PDF to the bulletin bucket.
func (s *S3Service) UploadBulletinPDF(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if s.bulletinBucket == "" {
		return "", fmt.Errorf("BULLETIN_BUCKET not set")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		return "", fmt.Errorf("invalid file type: expected application/pdf")
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".pdf" {
		return "", fmt.Errorf("invalid file extension: expected .pdf")
	}

	key := fmt.Sprintf("bulletin/%s.pdf", uuid.New().String())

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bulletinBucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload bulletin PDF to S3: %v", err)
	}

	return s.publicObjectURL(s.bulletinBucket, key), nil
}
