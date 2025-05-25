package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/errs"
	"github.com/sunnyyssh/designing-software-cw2/file-storage/internal/model"
)

type FileMetaRepository struct {
	db *pgxpool.Pool
}

func (r *FileMetaRepository) Store(ctx context.Context, meta *model.FileMeta) (*model.FileMeta, error) {
	q := `
	INSERT INTO file_meta (md5_hash, url)
	VALUES ($1, $2)
	RETURNING id`

	if err := r.db.QueryRow(ctx, q, meta.MD5Hash, meta.URL).Scan(&meta.ID); err != nil {
		return nil, fmt.Errorf("failed to insert file meta: %w", err)
	}

	return meta, nil
}

func (r *FileMetaRepository) Get(ctx context.Context, id int64) (*model.FileMeta, error) {
	q := `
	SELECT md5_hash, url 
	FROM file_meta
	WHERE id = $1`

	var (
		hash string
		url  string
	)

	if err := r.db.QueryRow(ctx, q, id).Scan(&hash, &url); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("file meta with id %d is not found", id)
		}
		return nil, fmt.Errorf("failed to insert file meta: %w", err)
	}

	return &model.FileMeta{
		ID:      id,
		MD5Hash: model.MD5(hash),
		URL:     url,
	}, nil
}

func (r *FileMetaRepository) ListByHash(ctx context.Context, hash model.MD5) ([]model.FileMeta, error) {
	q := `
	SELECT id, url 
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
			id  int64
			url string
		)
		if err := rows.Scan(&id, &url); err != nil {
			return nil, fmt.Errorf("failed to insert file meta: %w", err)
		}
		metas = append(metas, model.FileMeta{
			ID:      id,
			MD5Hash: hash,
			URL:     url,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows returned err: %w", err)
	}

	return metas, nil
}
