package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Key prefixes inside the single app bucket (see BUCKET). Matches console folders:
// character-images/, class images/ (seed art), maps/{backgrounds,tokens}/, monsters/, bulletin/, storeowners/.
const (
	s3PrefixCharacterImages = "character-images/"
	s3PrefixMapBackground   = "maps/backgrounds/"
	s3PrefixMonsters        = "monsters/"
	s3PrefixBulletin        = "bulletin/"
	s3PrefixStoreowners     = "storeowners/"
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

	bucket := os.Getenv("BUCKET")
	if bucket == "" {
		log.Warn().Msg("BUCKET environment variable is not set")
	}

	return &S3Service{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		region: cfg.Region,
	}, nil
}

func (s *S3Service) publicObjectURL(key string) string {
	region := s.region
	if region == "" {
		region = "us-east-1"
	}
	// Encode each path segment so keys like "class images/foo.jpg" work in browsers.
	segs := strings.Split(key, "/")
	for i := range segs {
		segs[i] = url.PathEscape(segs[i])
	}
	encKey := strings.Join(segs, "/")
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, region, encKey)
}

func (s *S3Service) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if s.bucket == "" {
		return "", fmt.Errorf("BUCKET not set")
	}

	ext := filepath.Ext(fileHeader.Filename)
	key := fmt.Sprintf("%s%s%s", s3PrefixCharacterImages, uuid.New().String(), ext)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return s.publicObjectURL(key), nil
}

// UploadMapBackgroundImage stores a battle map background under maps/backgrounds/ in the app bucket.
func (s *S3Service) UploadMapBackgroundImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if s.bucket == "" {
		return "", fmt.Errorf("BUCKET not set")
	}

	ext := filepath.Ext(fileHeader.Filename)
	key := fmt.Sprintf("%s%s%s", s3PrefixMapBackground, uuid.New().String(), ext)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload map background to S3: %v", err)
	}

	return s.publicObjectURL(key), nil
}

// UploadStream uploads an object to S3 from an io.Reader, generating a unique key.
func (s *S3Service) UploadStream(ctx context.Context, reader io.Reader, filename, contentType string) (string, error) {
	if s.bucket == "" {
		return "", fmt.Errorf("BUCKET not set")
	}

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("%s%s%s", s3PrefixMonsters, uuid.New().String(), ext)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload stream to S3: %v", err)
	}

	return s.publicObjectURL(key), nil
}

// StoreownerPortraitURL returns the public HTTPS URL for a vendor portrait uploaded as
// storeowners/{displayName}.png (same basename as seed Name). Empty string if BUCKET is unset.
func (s *S3Service) StoreownerPortraitURL(displayName string) string {
	if s == nil || s.bucket == "" || strings.TrimSpace(displayName) == "" {
		return ""
	}
	key := s3PrefixStoreowners + strings.TrimSpace(displayName) + ".png"
	return s.publicObjectURL(key)
}

// UploadBulletinPDF uploads a PDF under bulletin/ in the app bucket.
func (s *S3Service) UploadBulletinPDF(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if s.bucket == "" {
		return "", fmt.Errorf("BUCKET not set")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		return "", fmt.Errorf("invalid file type: expected application/pdf")
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".pdf" {
		return "", fmt.Errorf("invalid file extension: expected .pdf")
	}

	key := fmt.Sprintf("%s%s.pdf", s3PrefixBulletin, uuid.New().String())

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload bulletin PDF to S3: %v", err)
	}

	return s.publicObjectURL(key), nil
}
