package external

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type QuickchartClient struct{}

func NewQuickchartClient() *QuickchartClient {
	return &QuickchartClient{}
}

func (*QuickchartClient) GetWordCloud(ctx context.Context, text string) ([]byte, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://quickchart.io/wordcloud?text=%s", url.QueryEscape(text)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get word cloud: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get word cloud: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
