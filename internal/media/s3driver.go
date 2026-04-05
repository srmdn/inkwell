package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Driver stores files in an S3-compatible bucket (AWS S3, Cloudflare R2, NevaObjects, MinIO).
type S3Driver struct {
	client    *s3.Client
	bucket    string
	publicURL string // base URL for public file access, e.g. https://s3.nevaobjects.id/my-bucket
}

// NewS3Driver creates an S3Driver using static credentials.
// endpoint: S3-compatible API endpoint (e.g. https://s3.nevaobjects.id)
// bucket:   bucket name
// region:   region string; use "auto" for providers that don't require one
// publicURL: base URL prepended to keys for PublicURL(), e.g. https://s3.nevaobjects.id/my-bucket
func NewS3Driver(endpoint, bucket, region, accessKey, secretKey, publicURL string) *S3Driver {
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       region,
		Credentials:  creds,
		UsePathStyle: true, // required for most S3-compatible providers
	})

	return &S3Driver{
		client:    client,
		bucket:    bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
	}
}

func (d *S3Driver) Upload(filename string, r io.Reader, contentType string) (string, int64, error) {
	safe := sanitizeFilename(filename)
	key := uuid.New().String() + "-" + safe

	// Read into memory to get size (S3 PutObject needs Content-Length for streaming readers).
	body, err := io.ReadAll(r)
	if err != nil {
		return "", 0, fmt.Errorf("reading upload body: %w", err)
	}

	_, err = d.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(d.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(body),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(body))),
	})
	if err != nil {
		return "", 0, fmt.Errorf("s3 put object: %w", err)
	}

	return key, int64(len(body)), nil
}

func (d *S3Driver) Delete(key string) error {
	if strings.Contains(key, "/") || strings.Contains(key, "..") {
		return fmt.Errorf("invalid key")
	}

	_, err := d.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 delete object: %w", err)
	}
	return nil
}

func (d *S3Driver) PublicURL(key string) string {
	return d.publicURL + "/" + key
}
