package cencori

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ChatService struct {
	client *Client
}

func (s *ChatService) Chat(ctx context.Context, params ChatParams) (*ChatResponse, error) {
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

		reader := bufio.NewReader(resp.Body)

		for {
			select {
			case <-ctx.Done():
				chunks <- StreamChunk{Err: ctx.Err()}
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				}
				chunks <- StreamChunk{Err: fmt.Errorf("stream read : %w", err)}
				return
			}

			line = strings.TrimSpace(line)

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			// Stream termination
			if data == "[DONE]" {
				return
			}

			// Try decode as normal chunk
			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				var apiErr APIError
				if err2 := json.Unmarshal([]byte(data), &apiErr); err2 != nil {
					chunks <- StreamChunk{Err: &apiErr}
					return
				}

				chunks <- StreamChunk{Err: fmt.Errorf("unmarshal chunk: %w", err)}
				return
			}

			chunks <- chunk
		}
	}()

	return chunks, nil
}
