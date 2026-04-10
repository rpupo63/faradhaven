// Package seedmedia holds shared URLs for static seed assets in the app S3 bucket.
package seedmedia

import (
	"net/url"
	"strings"
)

// ClassRaceImagePrefix is the S3 key prefix for class/race artwork (folder name includes a space).
const ClassRaceImagePrefix = "class images/"

// BaseHost is the virtual-hosted–style origin for the app bucket (no trailing path).
// Keep region/bucket in sync with deployment BUCKET + AWS_REGION.
const BaseHost = "https://faradhaven-images.s3.us-east-2.amazonaws.com"

// URL returns the public HTTPS URL for a file under class images/ (e.g. "human.jpg" -> .../class%20images/human.jpg).
func URL(filename string) string {
	key := ClassRaceImagePrefix + strings.TrimPrefix(filename, "/")
	segs := strings.Split(key, "/")
	for i := range segs {
		segs[i] = url.PathEscape(segs[i])
	}
	return BaseHost + "/" + strings.Join(segs, "/")
}
