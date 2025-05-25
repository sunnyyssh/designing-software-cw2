package main

import (
	"context"
	"log"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/config"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/external"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/repository"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/rest"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/services"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load("/etc/analysis/config.yaml")
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
	awsS3Client := s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	db, err := pgxpool.New(ctx, cfg.PGConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewAnalysisRepository(db)
	storageClient := external.NewStorageClient(cfg.StorageURL)
	s3Client := external.NewS3Client(
		&external.S3Config{
			FileBucket:     cfg.S3.FileBucket,
			ImageBucket:    cfg.S3.ImageBucket,
			ImageURLPrefix: cfg.S3.ImageURLPrefix,
		},
		awsS3Client,
	)
	quickchartClient := external.NewQuickchartClient()

	svc := services.NewAnalysisService(
		repo,
		storageClient,
		s3Client,
		quickchartClient,
	)

	handler := rest.NewAnalysisService(svc)

	r := gin.New()

	r.Use(gin.Logger())

	r.GET("/analyze/:id", handler.Analyze)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
