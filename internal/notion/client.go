package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"telenotion/internal/common"
)

const notionBaseURL = "https://api.notion.com/v1/"

type NotionClient struct {
	*common.Client
}

func NewNotionClient(token string, httpClient *http.Client) *NotionClient {
	return &NotionClient{
		Client: common.NewClient(token, httpClient),
	}
}

type DatabaseResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

func (n *NotionClient) GetDatabase(ctx context.Context, databaseID string) (*DatabaseResponse, error) {
	url := fmt.Sprintf("%sdatabases/%s", notionBaseURL, databaseID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Authorization", "Bearer "+n.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("notion API error (%d): %s", resp.StatusCode, string(b))
	}

	var db DatabaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&db); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &db, nil
}
