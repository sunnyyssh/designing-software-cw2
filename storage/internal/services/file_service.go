package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/sunnyyssh/designing-software-cw2/storage/internal/model"
)

type FileMetaRepository interface {
	Store(ctx context.Context, meta *model.FileMeta) (*model.FileMeta, error)
	Get(ctx context.Context, id int64) (*model.FileMeta, error)
	ListByHash(ctx context.Context, hash model.MD5) ([]model.FileMeta, error)
}

type S3Config struct {
	Bucket    string
	UrlPrefix string
}

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type FileService struct {
	repo     FileMetaRepository
	s3Cfg    *S3Config
	s3Client S3Client
}

func NewFileService(repo FileMetaRepository, s3Cfg *S3Config, s3Client S3Client) *FileService {
	return &FileService{
		repo:     repo,
		s3Client: s3Client,
		s3Cfg:    s3Cfg,
	}
}

func (s *FileService) Upload(ctx context.Context, text string) (*model.FileMeta, error) {
	sha256Hash := sha256.Sum256([]byte(text))
	s3Key := uuid.NewString()
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:         aws.String(s.s3Cfg.Bucket),
		Key:            aws.String(s3Key),
		Body:           bytes.NewReader([]byte(text)),
		ContentType:    aws.String("text/plain"),
		ChecksumSHA256: aws.String(base64.StdEncoding.EncodeToString(sha256Hash[:])),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put object to s3: %w", err)
	}

	md5Hash := md5.Sum([]byte(text))

	meta := &model.FileMeta{
		MD5Hash: model.MD5(base64.StdEncoding.EncodeToString(md5Hash[:])),
		S3Key:   s3Key,
		URL:     s.s3Cfg.UrlPrefix + s3Key,
	}

	meta, err = s.repo.Store(ctx, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (s *FileService) Get(ctx context.Context, id int64) (*model.FileMeta, error) {
	return s.repo.Get(ctx, id)
}

func (s *FileService) ListByHash(ctx context.Context, hash model.MD5) ([]model.FileMeta, error) {
	return s.repo.ListByHash(ctx, hash)
}
