// Creates placeholder objects under maps/ in BUCKET so the AWS console shows folder structure.
// S3 has no real directories; empty keys like maps/backgrounds/.keep establish prefixes.
//
// Usage (from backend/):
//
//	go run ./cmd/create_s3_map_folders -dry-run
//	go run ./cmd/create_s3_map_folders
package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
)

// Placeholder keys to create (battle map–related image prefixes).
var mapImagePlaceholders = []string{
	"maps/.keep",
	"maps/backgrounds/.keep",
	"maps/tokens/.keep",
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("info: no .env in cwd: %v", err)
	}
	bootstrap.LoadEnv()

	dryRun := flag.Bool("dry-run", false, "print keys only; no S3 writes")
	flag.Parse()

	bucket := os.Getenv("BUCKET")
	if bucket == "" {
		log.Fatal("BUCKET is not set (add it to backend/.env)")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}
	client := s3.NewFromConfig(cfg)

	for _, key := range mapImagePlaceholders {
		if *dryRun {
			log.Printf("would put s3://%s/%s", bucket, key)
			continue
		}
		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader(nil),
		})
		if err != nil {
			log.Fatalf("put %s: %v", key, err)
		}
		log.Printf("created s3://%s/%s", bucket, key)
	}
	if *dryRun {
		return
	}
	log.Printf("done: %d placeholder object(s)", len(mapImagePlaceholders))
}
