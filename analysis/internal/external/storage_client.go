package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/errs"
	"github.com/sunnyyssh/designing-software-cw2/analysis/internal/model"
)

type StorageClient struct {
	url string
}

func NewStorageClient(url string) *StorageClient {
	return &StorageClient{
		url: url,
	}
}

func (c *StorageClient) Get(ctx context.Context, id int64) (*model.FileMeta, error) {
	resp, err := http.Get(
		fmt.Sprintf("%s/file/%d", c.url, id),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, errs.NotFound("file with id %d not found in storage", id)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get file meta from storage API: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := new(model.FileMeta)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file meta: %w", err)
	}

	return res, nil
}

func (c *StorageClient) ListByHash(ctx context.Context, hash model.MD5) ([]model.FileMeta, error) {
	resp, err := http.Get(
		fmt.Sprintf("%s/file/hash?hash=%s", c.url, url.QueryEscape(string(hash))),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get file meta from storage API: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res []model.FileMeta
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file meta: %w", err)
	}

	return res, nil
}
