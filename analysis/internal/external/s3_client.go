package external

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Config struct {
	FileBucket     string
	ImageBucket    string
	ImageURLPrefix string
}

type S3Client struct {
	cfg      *S3Config
	s3Client *s3.Client
}

func NewS3Client(cfg *S3Config, s3Client *s3.Client) *S3Client {
	return &S3Client{
		s3Client: s3Client,
		cfg:      cfg,
	}
}

func (c *S3Client) GetFileText(ctx context.Context, key string) (string, error) {
	obj, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.FileBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}

	defer obj.Body.Close()
	b, err := io.ReadAll(obj.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *S3Client) StoreWordCloud(ctx context.Context, mimeType string, data []byte) (string, error) {
	sha256Hash := sha256.Sum256(data)
	s3Key := uuid.NewString()

	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:         aws.String(c.cfg.ImageBucket),
		Key:            aws.String(s3Key),
		Body:           bytes.NewReader(data),
		ContentType:    aws.String(mimeType),
		ChecksumSHA256: aws.String(base64.StdEncoding.EncodeToString(sha256Hash[:])),
	})
	if err != nil {
		return "", fmt.Errorf("failed to put object to s3: %w", err)
	}

	return c.cfg.ImageURLPrefix + s3Key, nil
}
