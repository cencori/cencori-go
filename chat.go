package cencori

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ChatService struct {
	client *Client
}

func (s *ChatService) Create(ctx context.Context, params ChatParams) (*ChatResponse, error) {
	params.Stream = false
	return doRequest[ChatParams, ChatResponse](s.client, ctx, "POST", "/api/ai/chat", &params)
}

func (s *ChatService) Stream(ctx context.Context, params ChatParams) (<-chan StreamChunk, error) {
	params.Stream = true

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.client.BaseURL+"/api/ai/chat",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.ApiKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, handleError(resp)
	}

	chunks := make(chan StreamChunk)

	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				chunks <- StreamChunk{Err: fmt.Errorf("unmarshal chunk : %w", err)}
				return
			}

			select {
			case <-ctx.Done():
				chunks <- StreamChunk{Err: ctx.Err()}
				return
			case chunks <- chunk:
			}
		}

		if err := scanner.Err(); err != nil {
			chunks <- StreamChunk{Err: fmt.Errorf("stream scan: %w", err)}
		}

	}()

	return chunks, nil
}
