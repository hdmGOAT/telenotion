package notion

import (
	"bytes"
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


func (c *NotionClient) GetToDos(
	ctx context.Context,
	databaseID string,
	filter any, 
) ([]Todo, error) {
	reqBody := map[string]any{}
	if filter != nil {
		reqBody["filter"] = filter
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal filter: %w", err)
	}

	url := fmt.Sprintf("%sdata_sources/%s/query", notionBaseURL, databaseID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Notion-Version", "2022-06-28") 
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("notion query failed: %s – %s", resp.Status, string(body))
	}

	var dbResp DatabaseQueryResponse
	if decode_err := json.NewDecoder(resp.Body).Decode(&dbResp); decode_err != nil {
		return nil, fmt.Errorf("decode: %w", decode_err)
	}

	todos, err := MapNotionToTodos(dbResp)
	if err != nil {
		return nil, fmt.Errorf("map: %w", err)
	}

	return todos, nil
}
