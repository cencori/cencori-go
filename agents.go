package cencori

import (
	"context"
	"fmt"
)

// AgentsService provides methods for managing AI agents and their API keys.
type AgentsService struct {
	client *Client
}

// List retrieves all agents for the project.
func (s *AgentsService) List(ctx context.Context) (*AgentListResponse, error) {
	return doRequest[any, AgentListResponse](s.client, ctx, "GET", "/v1/agents", nil)
}

// Create creates a new AI agent.
func (s *AgentsService) Create(ctx context.Context, params *CreateAgentParams) (*Agent, error) {
	return doRequest[CreateAgentParams, Agent](s.client, ctx, "POST", "/v1/agents", params)
}

// Get retrieves an agent by ID.
func (s *AgentsService) Get(ctx context.Context, agentID string) (*Agent, error) {
	path := fmt.Sprintf("/v1/agents/%s", agentID)
	return doRequest[any, Agent](s.client, ctx, "GET", path, nil)
}

// UpdateConfig updates an agent's configuration.
func (s *AgentsService) UpdateConfig(ctx context.Context, agentID string, params *UpdateAgentParams) (*Agent, error) {
	path := fmt.Sprintf("/v1/agents/%s", agentID)
	return doRequest[UpdateAgentParams, Agent](s.client, ctx, "PATCH", path, params)
}

// Delete deletes an agent by ID.
func (s *AgentsService) Delete(ctx context.Context, agentID string) error {
	path := fmt.Sprintf("/v1/agents/%s", agentID)
	_, err := doRequest[any, any](s.client, ctx, "DELETE", path, nil)
	return err
}

// CreateKey creates an API key for an agent.
func (s *AgentsService) CreateKey(ctx context.Context, agentID string, params *CreateAgentKeyParams) (*AgentKey, error) {
	path := fmt.Sprintf("/v1/agents/%s/keys", agentID)
	return doRequest[CreateAgentKeyParams, AgentKey](s.client, ctx, "POST", path, params)
}
