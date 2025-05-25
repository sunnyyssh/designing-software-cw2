package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunnyyssh/designing-software-cw2/storage/internal/errs"
	"github.com/sunnyyssh/designing-software-cw2/storage/internal/model"
)

type FileMetaRepository struct {
	db *pgxpool.Pool
}

func NewFileMetaRepository(db *pgxpool.Pool) *FileMetaRepository {
	return &FileMetaRepository{
		db: db,
	}
}

func (r *FileMetaRepository) Store(ctx context.Context, meta *model.FileMeta) (*model.FileMeta, error) {
	q := `
	INSERT INTO file_meta (md5_hash, s3_key, url)
	VALUES ($1, $2, $3)
	RETURNING id`

	if err := r.db.QueryRow(ctx, q, meta.MD5Hash, meta.S3Key, meta.URL).Scan(&meta.ID); err != nil {
		return nil, fmt.Errorf("failed to insert file meta: %w", err)
	}

	return meta, nil
}

func (r *FileMetaRepository) Get(ctx context.Context, id int64) (*model.FileMeta, error) {
	q := `
	SELECT md5_hash, s3_key, url 
	FROM file_meta
	WHERE id = $1`

	var (
		hash  string
		s3Key string
		url   string
	)

	if err := r.db.QueryRow(ctx, q, id).Scan(&hash, &s3Key, &url); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("file meta with id %d is not found", id)
		}
		return nil, fmt.Errorf("failed to get file meta: %w", err)
	}

	return &model.FileMeta{
		ID:      id,
		MD5Hash: model.MD5(hash),
		S3Key:   s3Key,
		URL:     url,
	}, nil
}

func (r *FileMetaRepository) ListByHash(ctx context.Context, hash model.MD5) ([]model.FileMeta, error) {
	q := `
	SELECT id, s3_key, url 
	FROM file_meta
	WHERE md5_hash = $1`

	rows, err := r.db.Query(ctx, q, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to query file metas: %w", err)
	}
	defer rows.Close()

	metas := make([]model.FileMeta, 0)

	for rows.Next() {
		var (
			id    int64
			s3Key string
			url   string
		)
		if err := rows.Scan(&id, &s3Key, &url); err != nil {
			return nil, fmt.Errorf("failed to get file meta: %w", err)
		}
		metas = append(metas, model.FileMeta{
			ID:      id,
			MD5Hash: hash,
			S3Key:   s3Key,
			URL:     url,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows returned err: %w", err)
	}

	return metas, nil
}
