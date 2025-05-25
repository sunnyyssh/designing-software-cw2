package main

import (
	"context"
	"log"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/config"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/repository"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/rest"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/services"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("/etc/storage/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	s3Cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.S3.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.S3.AccessKeyId,
				cfg.S3.SecretAccessKey,
				"",
			),
		),
		awsconfig.WithBaseEndpoint(cfg.S3.EndpointURL),
	)
	if err != nil {
		log.Fatalf("loading default config failed: %s", err)
	}
	s3Client := s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	db, err := pgxpool.New(ctx, cfg.PGConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewFileMetaRepository(db)

	svc := services.NewFileService(
		repo,
		&services.S3Config{
			Bucket:    cfg.S3.Bucket,
			UrlPrefix: cfg.S3.URLPrefix,
		},
		s3Client,
	)

	handler := rest.NewFileHandler(svc)

	r := gin.New()

	r.Use(gin.Logger())

	r.POST("/file", handler.Upload)
	r.GET("/file/:id", handler.Get)
	r.GET("/file/hash", handler.ListByHash)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
