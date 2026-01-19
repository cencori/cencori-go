package cencori

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatParams struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int      `json:"maxTokens,omitempty"`
	Stream      bool      `json:"stream"`
	UserID      *string   `json:"userId,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Content      string  `json:"content"`
	Model        string  `json:"model"`
	Provider     string  `json:"provider"`
	Usage        Usage   `json:"usage"`
	CostUSD      float64 `json:"cost_usd"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamChunk struct {
	Delta        string `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
}

func (c *Client) Chat(ctx context.Context, params ChatParams) (*ChatResponse, error) {
	params.Stream = false
	return doRequest[ChatParams, ChatResponse](c, ctx, "POST", "/api/ai/chat", &params)
}

func (c *Client) ChatStream(ctx context.Context, params ChatParams) (<-chan StreamChunk, error) {
	params.Stream = true

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.BaseURL+"/api/ai/chat",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", c.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Code:       "READ_ERROR",
				Message:    fmt.Sprintf("failed to read response body: %v", err),
			}
		}
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Code:       "UNKNOWN",
				Message:    string(body),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	chunks := make(chan StreamChunk)

	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)

		for {
			var chunk StreamChunk
			if err := decoder.Decode(&chunk); err != nil {
				return
			}

			select {
			case <-ctx.Done():
				return
			case chunks <- chunk:
			}
		}
	}()

	return chunks, nil
}
