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

// ChatService provides methods for managing chat-related operations.
// It uses a Client to communicate with the chat API endpoints.
type ChatService struct {
	client *Client
}

// Create sends a chat request to the AI service and returns the response.
// It disables streaming and makes a synchronous request to the /api/ai/chat endpoint.
// The context can be used to cancel the request or set a timeout.
// It returns a ChatResponse on success or an error if the request fails.
func (s *ChatService) Create(ctx context.Context, params *ChatParams) (*ChatResponse, error) {
	params.Stream = false
	return doRequest[ChatParams, ChatResponse](s.client, ctx, "POST", "/api/ai/chat", params)
}

// Completions is a convenience method that wraps Create for simple text completions.
// It takes a single prompt string and converts it to a chat message internally.
// This is useful for quick, single-turn completions without managing conversation history.
func (s *ChatService) Completions(ctx context.Context, params CompletionParams) (*ChatResponse, error) {
	chatParams := &ChatParams{
		Model:       params.Model,
		Temperature: params.Temperature,
		MaxTokens:   params.MaxTokens,
		Messages: []Message{
			{Role: "user", Content: params.Prompt},
		},
	}
	return s.Create(ctx, chatParams)
}

// Embeddings generates vector embeddings for the given input text(s).
// Input can be a single string or a slice of strings.
// Returns an EmbeddingResponse containing the embeddings and token usage.
func (s *ChatService) Embeddings(ctx context.Context, params EmbeddingParams) (*EmbeddingResponse, error) {
	return doRequest[EmbeddingParams, EmbeddingResponse](s.client, ctx, "POST", "/api/v1/embeddings", &params)
}

// Stream sends a chat request with streaming enabled and returns a channel that receives
// chat response chunks as they arrive from the server. The stream continues until the server
// sends a "[DONE]" message or an error occurs. The context can be used to cancel the stream.
// If the context is cancelled, it simply closes.
// The returned channel will be closed when the stream ends or an error occurs.
func (s *ChatService) Stream(ctx context.Context, params *ChatParams) (<-chan StreamChunk, error) {
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
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.httpClient.Do(req) //nolint:bodyclose // Body is closed by the streaming goroutine
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		err := handleError(resp)
		resp.Body.Close() //nolint:errcheck // Closing the response body; error can be ignored here.

		return nil, err
	}

	chunks := make(chan StreamChunk)

	go func() {
		defer close(chunks)
		defer resp.Body.Close() //nolint:errcheck // Closing the response body; error can be ignored here.

		done := make(chan struct{})
		defer close(done)

		go func() {
			select {
			case <-ctx.Done():
				resp.Body.Close() //nolint:errcheck // Closing the response body; error can be ignored here.
			case <-done:
				return
			}
		}()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')

			if err != nil {
				if ctx.Err() != nil {
					return
				}

				if err == io.EOF {
					return
				}

				chunks <- StreamChunk{Err: fmt.Errorf("stream read: %w", err)}
				return
			}

			line = strings.TrimSpace(line)

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			if data == "[DONE]" {
				return
			}

			if strings.Contains(data, "\"error\":") {
				var apiErr APIError
				if err := json.Unmarshal([]byte(data), &apiErr); err == nil {
					apiErr.fillSentinel()
					chunks <- StreamChunk{Err: &apiErr}
					return
				}
			}

			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				chunks <- StreamChunk{Err: fmt.Errorf("unmarshal chunk: %w", err)}
				return
			}

			// Extract delta content for convenience
			if len(chunk.Choices) > 0 {
				chunk.Delta = chunk.Choices[0].Delta.Content
			}

			chunks <- chunk
		}
	}()

	return chunks, nil
}

// GenerateObject generates structured output matching a JSON schema.
// Uses function calling to enforce the schema on the model output.
func (s *ChatService) GenerateObject(ctx context.Context, params *GenerateObjectRequest) (*GenerateObjectResponse, error) {
	messages := params.Messages
	if messages == nil {
		messages = []Message{
			{Role: "user", Content: params.Prompt},
		}
	}

	schemaName := params.SchemaName
	if schemaName == "" {
		schemaName = "generate_object"
	}

	chatParams := &ChatParams{
		Model:     params.Model,
		Messages:  messages,
		MaxTokens: params.MaxTokens,
		Tools: []ToolDefinition{
			{
				Type: "function",
				Function: ToolFunction{
					Name:        schemaName,
					Description: params.SchemaDescription,
					Parameters:  params.Schema,
				},
			},
		},
		ToolChoice: ToolChoiceFunctionName(schemaName),
	}
	if params.Temperature != nil {
		chatParams.Temperature = params.Temperature
	}

	resp, err := doRequest[ChatParams, chatResponseRaw](s.client, ctx, "POST", "/api/ai/chat", chatParams)
	if err != nil {
		return nil, err
	}

	var toolCall *ToolCall
	if len(resp.ToolCalls) > 0 {
		toolCall = &resp.ToolCalls[0]
	} else if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		toolCall = &resp.Choices[0].Message.ToolCalls[0]
	}

	if toolCall == nil {
		return nil, fmt.Errorf("model did not return structured output")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse structured output: %w", err)
	}

	return &GenerateObjectResponse{
		Object: parsed,
		Usage:  resp.Usage,
	}, nil
}

type chatResponseRaw struct {
	ID        string `json:"id"`
	Model     string `json:"model"`
	Content   string `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"toolCalls,omitempty"`
	Choices   []struct {
		Message struct {
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices,omitempty"`
	Usage Usage `json:"usage"`
}

// GenerateImage generates images from a text prompt.
func (s *ChatService) GenerateImage(ctx context.Context, params *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	if params.Model == "" {
		params.Model = "dall-e-3"
	}
	return doRequest[ImageGenerationRequest, ImageGenerationResponse](s.client, ctx, "POST", "/api/ai/images/generate", params)
}

// Rag performs a RAG (Retrieval-Augmented Generation) chat request.
// It searches the memory namespace for relevant context and includes it automatically.
func (s *ChatService) Rag(ctx context.Context, params *RagRequest) (*RagResponse, error) {
	params.Stream = false
	return doRequest[RagRequest, RagResponse](s.client, ctx, "POST", "/api/ai/rag", params)
}

// RagStream streams a RAG response with automatic memory context.
// Returns a channel that yields RagStreamChunk values as they arrive.
func (s *ChatService) RagStream(ctx context.Context, params *RagRequest) (<-chan RagStreamChunk, error) {
	params.Stream = true

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		s.client.BaseURL+"/api/ai/rag",
		bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		err := handleError(resp)
		resp.Body.Close()
		return nil, err
	}

	chunks := make(chan RagStreamChunk)

	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		done := make(chan struct{})
		defer close(done)

		go func() {
			select {
			case <-ctx.Done():
				resp.Body.Close()
			case <-done:
			}
		}()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if ctx.Err() != nil || err == io.EOF {
					return
				}
				chunks <- RagStreamChunk{Err: fmt.Errorf("rag stream read: %w", err)}
				return
			}

			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}

			var chunk RagStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				chunks <- RagStreamChunk{Err: fmt.Errorf("unmarshal rag chunk: %w", err)}
				return
			}
			chunks <- chunk
		}
	}()

	return chunks, nil
}

// Responses sends a request to the OpenAI-compatible Responses API.
// Supports built-in tools: web_search_preview, file_search, code_interpreter.
func (s *ChatService) Responses(ctx context.Context, params *ResponsesRequest) (*ResponsesResponse, error) {
	params.Stream = false
	return doRequest[ResponsesRequest, ResponsesResponse](s.client, ctx, "POST", "/v1/responses", params)
}

// ResponsesStream streams responses from the Responses API via SSE.
// Returns a channel that yields ResponsesStreamEvent values.
func (s *ChatService) ResponsesStream(ctx context.Context, params *ResponsesRequest) (<-chan ResponsesStreamEvent, error) {
	params.Stream = true

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		s.client.BaseURL+"/v1/responses",
		bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		err := handleError(resp)
		resp.Body.Close()
		return nil, err
	}

	events := make(chan ResponsesStreamEvent)

	go func() {
		defer close(events)
		defer resp.Body.Close()

		done := make(chan struct{})
		defer close(done)

		go func() {
			select {
			case <-ctx.Done():
				resp.Body.Close()
			case <-done:
			}
		}()

		reader := bufio.NewReader(resp.Body)
		var eventType string

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if ctx.Err() != nil || err == io.EOF {
					return
				}
				events <- ResponsesStreamEvent{Err: fmt.Errorf("responses stream read: %w", err)}
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				eventType = ""
				continue
			}
			if strings.HasPrefix(line, "event: ") {
				eventType = strings.TrimPrefix(line, "event: ")
				continue
			}
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "" {
					continue
				}
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(data), &parsed); err != nil {
					continue
				}
				events <- ResponsesStreamEvent{Type: eventType, Data: parsed}
				eventType = ""
			}
		}
	}()

	return events, nil
}
