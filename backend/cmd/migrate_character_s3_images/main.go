// Copies character portrait objects into the unified app bucket under character-images/
// (same key prefix as services.S3Service / UploadFile).
//
// Typical use: legacy bucket characters/* -> BUCKET/character-images/*
// If your source keys already use character-images/, set -source-prefix character-images/
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_character_s3_images -dry-run
//	go run ./cmd/migrate_character_s3_images
//	go run ./cmd/migrate_character_s3_images -delete-source
//
// Optional env: CHARACTER_IMAGES_MIGRATE_FROM, CHARACTER_IMAGES_MIGRATE_TO (defaults: BUCKET for destination)
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/rpupo63/faradhaven/backend/internal/bootstrap"
)

const (
	defaultSourcePrefix = "characters/"
	// Must match services.s3PrefixCharacterImages (character-images/).
	defaultDestPrefix = "character-images/"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("info: no .env in cwd: %v", err)
	}
	bootstrap.LoadEnv()

	dryRun := flag.Bool("dry-run", false, "list planned copies only; no S3 writes")
	deleteSource := flag.Bool("delete-source", false, "delete each source object after a successful copy (ignored with -dry-run)")
	sourceBucket := flag.String("source-bucket", getenvDefault("CHARACTER_IMAGES_MIGRATE_FROM", "faradhaven-character-images"), "source S3 bucket")
	destBucket := flag.String("dest-bucket", getenvDefault("CHARACTER_IMAGES_MIGRATE_TO", getenvDefault("BUCKET", "faradhaven-images")), "destination S3 bucket (default: BUCKET from .env)")
	sourcePrefix := flag.String("source-prefix", defaultSourcePrefix, "key prefix under source bucket (trailing slash recommended)")
	destPrefix := flag.String("dest-prefix", defaultDestPrefix, "key prefix under destination bucket (trailing slash recommended)")
	flag.Parse()

	if *sourcePrefix != "" && !strings.HasSuffix(*sourcePrefix, "/") {
		log.Fatal("source-prefix should end with /")
	}
	if *destPrefix != "" && !strings.HasSuffix(*destPrefix, "/") {
		log.Fatal("dest-prefix should end with /")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}
	client := s3.NewFromConfig(cfg)

	keys, err := listKeys(ctx, client, *sourceBucket, *sourcePrefix)
	if err != nil {
		log.Fatalf("list source: %v", err)
	}
	if len(keys) == 0 {
		log.Printf("no objects under s3://%s/%s", *sourceBucket, *sourcePrefix)
		return
	}
	log.Printf("found %d object(s) to copy", len(keys))

	var copied, skipped, deleted int
	var copyErrs, delErrs int

	for _, srcKey := range keys {
		if strings.HasSuffix(srcKey, "/") {
			continue
		}
		destKey := mapKey(srcKey, *sourcePrefix, *destPrefix)
		if destKey == "" {
			skipped++
			log.Printf("skip (unexpected key shape): %s", srcKey)
			continue
		}

		if *dryRun {
			fmt.Printf("would copy s3://%s/%s -> s3://%s/%s\n", *sourceBucket, srcKey, *destBucket, destKey)
			continue
		}

		_, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(*destBucket),
			Key:        aws.String(destKey),
			CopySource: aws.String(copySourceHeader(*sourceBucket, srcKey)),
		})
		if err != nil {
			copyErrs++
			log.Printf("copy error %s: %v", srcKey, err)
			continue
		}
		copied++
		log.Printf("copied %s -> %s", srcKey, destKey)

		if *deleteSource {
			_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(*sourceBucket),
				Key:    aws.String(srcKey),
			})
			if err != nil {
				delErrs++
				log.Printf("delete source error %s: %v", srcKey, err)
				continue
			}
			deleted++
		}
	}

	if *dryRun {
		return
	}
	log.Printf("done: copied=%d skipped=%d deleted=%d copy_errors=%d delete_errors=%d", copied, skipped, deleted, copyErrs, delErrs)
	if copyErrs > 0 || delErrs > 0 {
		os.Exit(1)
	}
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func listKeys(ctx context.Context, client *s3.Client, bucket, prefix string) ([]string, error) {
	var keys []string
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, obj := range out.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
	}
	return keys, nil
}

func mapKey(srcKey, sourcePrefix, destPrefix string) string {
	if !strings.HasPrefix(srcKey, sourcePrefix) {
		return ""
	}
	rest := strings.TrimPrefix(srcKey, sourcePrefix)
	if rest == "" {
		return ""
	}
	return destPrefix + rest
}

// copySourceHeader builds the x-amz-copy-source value: bucket/key with key URL-encoded per S3 rules.
func copySourceHeader(bucket, key string) string {
	escaped := strings.ReplaceAll(url.PathEscape(key), "+", "%20")
	return fmt.Sprintf("%s/%s", bucket, escaped)
}
