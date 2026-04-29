// Copies all objects from the photos-for-apps bucket into the app bucket under the
// "class images/" prefix (matches seedmedia URLs for class/race artwork).
//
// List/delete use the source region; CopyObject uses the destination region (required when
// buckets differ in region).
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_photos_for_apps -dry-run
//	go run ./cmd/migrate_photos_for_apps
//	go run ./cmd/migrate_photos_for_apps -delete-source
//
// Optional env: PHOTOS_FOR_APPS_MIGRATE_SOURCE_BUCKET, PHOTOS_FOR_APPS_MIGRATE_DEST_BUCKET
// (defaults to BUCKET), PHOTOS_FOR_APPS_REGION (source, default us-east-2), AWS_REGION (destination).
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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("info: no .env in cwd: %v", err)
	}
	bootstrap.LoadEnv()

	dryRun := flag.Bool("dry-run", false, "list planned copies only; no S3 writes")
	deleteSource := flag.Bool("delete-source", false, "delete each source object after a successful copy (ignored with -dry-run)")
	sourceBucket := flag.String("source-bucket", getenvDefault("PHOTOS_FOR_APPS_MIGRATE_SOURCE_BUCKET", "photos-for-apps"), "source S3 bucket")
	destBucket := flag.String("dest-bucket", getenvDefault("PHOTOS_FOR_APPS_MIGRATE_DEST_BUCKET", getenvDefault("BUCKET", "faradhaven-images")), "destination S3 bucket (default: BUCKET from .env)")
	sourcePrefix := flag.String("source-prefix", "", "optional key prefix under source bucket (empty = entire bucket; use trailing / for a folder prefix)")
	destPrefix := flag.String("dest-prefix", "class images/", "destination prefix under dest bucket (must match seedmedia.ClassRaceImagePrefix)")
	sourceRegion := flag.String("source-region", getenvDefault("PHOTOS_FOR_APPS_REGION", "us-east-2"), "AWS region of the source bucket")
	destRegion := flag.String("dest-region", "", "AWS region of the destination bucket (default: AWS_REGION / config)")
	flag.Parse()

	if *sourcePrefix != "" && !strings.HasSuffix(*sourcePrefix, "/") {
		log.Fatal("source-prefix must be empty or end with /")
	}
	if *destPrefix != "" && !strings.HasSuffix(*destPrefix, "/") {
		log.Fatal("dest-prefix should end with /")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}

	dstReg := *destRegion
	if dstReg == "" {
		dstReg = cfg.Region
	}
	if dstReg == "" {
		dstReg = getenvDefault("AWS_REGION", "us-east-1")
	}

	srcClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = *sourceRegion
	})
	dstClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = dstReg
	})

	keys, err := listKeys(ctx, srcClient, *sourceBucket, *sourcePrefix)
	if err != nil {
		log.Fatalf("list source: %v", err)
	}
	if len(keys) == 0 {
		log.Printf("no objects under s3://%s/%s", *sourceBucket, *sourcePrefix)
		return
	}
	log.Printf("found %d object(s) to copy (source region %s, dest region %s)", len(keys), *sourceRegion, dstReg)

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

		_, err := dstClient.CopyObject(ctx, &s3.CopyObjectInput{
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
			_, err := srcClient.DeleteObject(ctx, &s3.DeleteObjectInput{
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
	if sourcePrefix != "" && !strings.HasPrefix(srcKey, sourcePrefix) {
		return ""
	}
	rest := strings.TrimPrefix(srcKey, sourcePrefix)
	if rest == "" {
		return ""
	}
	return destPrefix + rest
}

func copySourceHeader(bucket, key string) string {
	escaped := strings.ReplaceAll(url.PathEscape(key), "+", "%20")
	return fmt.Sprintf("%s/%s", bucket, escaped)
}
