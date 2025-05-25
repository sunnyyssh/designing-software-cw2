package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/model"
)

type AnalysisRepository struct {
	db *pgxpool.Pool
}

func NewAnalysisRepository(db *pgxpool.Pool) *AnalysisRepository {
	return &AnalysisRepository{
		db: db,
	}
}

func (r *AnalysisRepository) Store(ctx context.Context, a *model.AnalysisResult) error {
	q := `
	INSERT INTO file_analysis (file_id, plagiated_ids, word_count, symbol_count, word_cloud_url)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, q, a.FileID, a.PlagiatedIDs, a.WordCount, a.SymbolCount, a.WordCloudURL)
	if err != nil {
		return err
	}

	return nil
}

// Must return nil, nil if nothing found
func (r *AnalysisRepository) Get(ctx context.Context, id int64) (*model.AnalysisResult, error) {
	q := `
	SELECT plagiated_ids, word_count, symbol_count, word_cloud_url
	FROM file_analysis
	WHERE file_id = $1`

	var (
		plagiatedIDs []int64
		wordCount    int
		symbolCount  int
		wordCloudURL string
	)

	err := r.db.QueryRow(ctx, q, id).Scan(&plagiatedIDs, &wordCount, &symbolCount, &wordCloudURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.AnalysisResult{
		FileID:       id,
		PlagiatedIDs: plagiatedIDs,
		WordCount:    wordCount,
		SymbolCount:  symbolCount,
		WordCloudURL: wordCloudURL,
	}, nil
}
