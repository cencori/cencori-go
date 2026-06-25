package cencori

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// SessionsService provides methods for managing durable execution sessions for AI agents.
type SessionsService struct {
	client *Client
}

func (s *SessionsService) request(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	url := s.client.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, handleError(resp)
	}

	return io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
}

// Create creates a new session.
func (s *SessionsService) Create(ctx context.Context, params *CreateSessionParams) (*Session, error) {
	data, err := s.request(ctx, "POST", "/v1/sessions", params)
	if err != nil {
		return nil, err
	}
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("decode session: %w", err)
	}
	return &session, nil
}

// List retrieves sessions with optional filtering.
func (s *SessionsService) List(ctx context.Context, params *SessionListParams) (*PaginatedSessionsResponse, error) {
	path := "/v1/sessions"
	if params != nil {
		var qs []string
		if params.Page != nil {
			qs = append(qs, fmt.Sprintf("page=%d", *params.Page))
		}
		if params.Limit != nil {
			qs = append(qs, fmt.Sprintf("limit=%d", *params.Limit))
		}
		if params.Status != "" {
			qs = append(qs, fmt.Sprintf("status=%s", params.Status))
		}
		if params.AgentID != "" {
			qs = append(qs, fmt.Sprintf("agent_id=%s", params.AgentID))
		}
		if len(qs) > 0 {
			path += "?" + strings.Join(qs, "&")
		}
	}
	return doRequest[any, PaginatedSessionsResponse](s.client, ctx, "GET", path, nil)
}

// Get retrieves a session by ID.
func (s *SessionsService) Get(ctx context.Context, sessionID string) (*Session, error) {
	path := fmt.Sprintf("/v1/sessions/%s", sessionID)
	return doRequest[any, Session](s.client, ctx, "GET", path, nil)
}

// Delete deletes a session by ID.
func (s *SessionsService) Delete(ctx context.Context, sessionID string) (*SessionDeleteResponse, error) {
	path := fmt.Sprintf("/v1/sessions/%s", sessionID)
	return doRequest[any, SessionDeleteResponse](s.client, ctx, "DELETE", path, nil)
}

// SubmitTurn submits a turn in a session and returns the raw HTTP response.
// For streaming, use SubmitTurnStream.
func (s *SessionsService) SubmitTurn(ctx context.Context, sessionID string, params *TurnParams) (*http.Response, error) {
	url := s.client.BaseURL + fmt.Sprintf("/v1/sessions/%s/turns", sessionID)

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal turn: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err := handleError(resp)
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}

// SubmitTurnStream submits a turn and returns the response body as a stream.
func (s *SessionsService) SubmitTurnStream(ctx context.Context, sessionID string, params *TurnParams) (io.ReadCloser, error) {
	resp, err := s.SubmitTurn(ctx, sessionID, params)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// GetEvents retrieves events for a session.
func (s *SessionsService) GetEvents(ctx context.Context, sessionID string, params *SessionEventsParams) (*PaginatedEventsResponse, error) {
	path := fmt.Sprintf("/v1/sessions/%s/events", sessionID)
	if params != nil {
		var qs []string
		if params.Page != nil {
			qs = append(qs, fmt.Sprintf("page=%d", *params.Page))
		}
		if params.Limit != nil {
			qs = append(qs, fmt.Sprintf("limit=%d", *params.Limit))
		}
		if params.TurnNumber != nil {
			qs = append(qs, fmt.Sprintf("turn_number=%d", *params.TurnNumber))
		}
		if len(qs) > 0 {
			path += "?" + strings.Join(qs, "&")
		}
	}
	return doRequest[any, PaginatedEventsResponse](s.client, ctx, "GET", path, nil)
}

// Approve approves a pending action in a session and returns the raw response.
// For streaming, use ApproveStream.
func (s *SessionsService) Approve(ctx context.Context, sessionID string, params *ApproveRejectParams) (*http.Response, error) {
	url := s.client.BaseURL + fmt.Sprintf("/v1/sessions/%s/approve", sessionID)

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal approve: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err := handleError(resp)
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}

// ApproveStream approves an action and returns the response body as a stream.
func (s *SessionsService) ApproveStream(ctx context.Context, sessionID string, params *ApproveRejectParams) (io.ReadCloser, error) {
	resp, err := s.Approve(ctx, sessionID, params)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Reject rejects a pending action in a session.
func (s *SessionsService) Reject(ctx context.Context, sessionID string, params *ApproveRejectParams) (*ApproveResponse, error) {
	path := fmt.Sprintf("/v1/sessions/%s/reject", sessionID)
	return doRequest[ApproveRejectParams, ApproveResponse](s.client, ctx, "POST", path, params)
}

// Response types for sessions.

type PaginatedSessionsResponse struct {
	Data       []Session  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type PaginatedEventsResponse struct {
	Data       []SessionEvent `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

type SessionEventsParams struct {
	Page       *int `json:"page,omitempty"`
	Limit      *int `json:"limit,omitempty"`
	TurnNumber *int `json:"turn_number,omitempty"`
}

type SessionDeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
