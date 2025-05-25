package services

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/model"
)

type AnalysisRepository interface {
	Store(context.Context, *model.AnalysisResult) error
	// Must return nil, nil if nothing found
	Get(context.Context, int64) (*model.AnalysisResult, error)
}

type StorageClient interface {
	Get(context.Context, int64) (*model.FileMeta, error)
	ListByHash(context.Context, model.MD5) ([]model.FileMeta, error)
}

type S3Client interface {
	GetFileText(_ context.Context, key string) (string, error)
	StoreWordCloud(_ context.Context, mimeType string, data []byte) (url string, err error)
}

type QuickchartClient interface {
	GetWordCloud(_ context.Context, text string) ([]byte, error)
}

type AnalysisService struct {
	repo             AnalysisRepository
	storageClient    StorageClient
	s3Client         S3Client
	quickchartClient QuickchartClient
}

func NewAnalysisService(
	repo AnalysisRepository,
	storageClient StorageClient,
	s3Client S3Client,
	quickchartClient QuickchartClient,
) *AnalysisService {
	return &AnalysisService{
		repo:             repo,
		storageClient:    storageClient,
		s3Client:         s3Client,
		quickchartClient: quickchartClient,
	}
}

func (s *AnalysisService) Analyze(ctx context.Context, id int64) (*model.AnalysisResult, error) {
	cached, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil
	}

	return s.analyze(ctx, id)
}

func (s *AnalysisService) analyze(ctx context.Context, id int64) (*model.AnalysisResult, error) {
	file, err := s.storageClient.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	plagiatedIDs, err := s.fetchPlagiatedIDs(ctx, id, file.MD5Hash)
	if err != nil {
		return nil, err
	}

	text, err := s.s3Client.GetFileText(ctx, file.S3Key)
	if err != nil {
		return nil, err
	}

	image, err := s.quickchartClient.GetWordCloud(ctx, text)
	if err != nil {
		return nil, err
	}

	imageURL, err := s.s3Client.StoreWordCloud(ctx, "image/svg+xml", image)
	if err != nil {
		return nil, err
	}

	return &model.AnalysisResult{
		FileID:       id,
		PlagiatedIDs: plagiatedIDs,
		WordCount:    len(strings.Fields(text)),
		SymbolCount:  utf8.RuneCountInString(text),
		WordCloudURL: imageURL,
	}, nil
}

func (s *AnalysisService) fetchPlagiatedIDs(ctx context.Context, id int64, hash model.MD5) ([]int64, error) {
	metas, err := s.storageClient.ListByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	var res []int64
	for _, meta := range metas {
		if meta.ID != id {
			res = append(res, meta.ID)
		}
	}
	return res, nil
}
