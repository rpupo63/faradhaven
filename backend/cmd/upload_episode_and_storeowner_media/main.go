// Uploads local episode-images/ and storeowners/ to the app S3 bucket under
// episode-images/... and storeowners/... (episode subfolders are preserved).
//
// Usage (from backend/):
//
//	go run ./cmd/upload_episode_and_storeowner_media -dry-run
//	go run ./cmd/upload_episode_and_storeowner_media
//
// Env: BUCKET (required), AWS credentials/region as usual for the SDK.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
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

	dryRun := flag.Bool("dry-run", false, "print planned uploads only; no S3 writes")
	episodeDir := flag.String("episode-dir", "episode-images", "local directory containing episode subfolders (e.g. episode-4/)")
	storeownersDir := flag.String("storeowners-dir", "storeowners", "local directory of storeowner images")
	episodePrefix := flag.String("episode-prefix", "episode-images/", "S3 key prefix for episode images (trailing slash)")
	storeownersPrefix := flag.String("storeowners-prefix", "storeowners/", "S3 key prefix for storeowner images (trailing slash)")
	skipEpisode := flag.Bool("skip-episode", false, "do not upload episode-images")
	skipStoreowners := flag.Bool("skip-storeowners", false, "do not upload storeowners")
	flag.Parse()

	if *episodePrefix != "" && !strings.HasSuffix(*episodePrefix, "/") {
		log.Fatal("episode-prefix must end with /")
	}
	if *storeownersPrefix != "" && !strings.HasSuffix(*storeownersPrefix, "/") {
		log.Fatal("storeowners-prefix must end with /")
	}

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

	var total int
	if !*skipEpisode {
		n, err := uploadDir(ctx, client, bucket, *episodeDir, *episodePrefix, *dryRun)
		if err != nil {
			log.Fatalf("episode-images: %v", err)
		}
		total += n
	}
	if !*skipStoreowners {
		n, err := uploadDir(ctx, client, bucket, *storeownersDir, *storeownersPrefix, *dryRun)
		if err != nil {
			log.Fatalf("storeowners: %v", err)
		}
		total += n
	}

	if *dryRun {
		log.Printf("dry-run: would upload %d file(s)", total)
		return
	}
	log.Printf("done: uploaded %d file(s) to s3://%s", total, bucket)
}

func uploadDir(ctx context.Context, client *s3.Client, bucket, localRoot, keyPrefix string, dryRun bool) (int, error) {
	st, err := os.Stat(localRoot)
	if err != nil {
		return 0, fmt.Errorf("stat %q: %w", localRoot, err)
	}
	if !st.IsDir() {
		return 0, fmt.Errorf("%q is not a directory", localRoot)
	}

	absRoot, err := filepath.Abs(localRoot)
	if err != nil {
		return 0, err
	}

	var count int
	err = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") {
			return nil
		}

		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		key := keyPrefix + rel

		if dryRun {
			log.Printf("would put s3://%s/%s <- %s", bucket, key, path)
			count++
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		ct := contentType(path)
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			Body:        f,
			ContentType: aws.String(ct),
		})
		if err != nil {
			return fmt.Errorf("put %s: %w", key, err)
		}
		log.Printf("uploaded s3://%s/%s", bucket, key)
		count++
		return nil
	})
	if err != nil {
		return count, err
	}
	return count, nil
}

func contentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}
	return "application/octet-stream"
}
