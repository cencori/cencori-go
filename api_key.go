package cencori

import (
	"context"
	"fmt"
)

// APIKeysService provides methods for managing api keys.
// It uses a Client to communicate with the api-keys API endpoints.
type APIKeysService struct {
	client *Client
}

// List retrieves all API keys for a given project and environment.
// It takes a context, projectID, and env as parameters and returns a slice of APIKey objects.
// Returns an error if the request fails.
func (s *APIKeysService) List(ctx context.Context, projectID, env string) ([]APIKey, error) {
	path := fmt.Sprintf("/api/projects/%s/api-keys?environment=%s", projectID, env)

	type response struct {
		Keys []APIKey `json:"keys"`
	}

	res, err := doRequest[any, response](s.client, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

// Create creates a new API key for the specified project.
// It takes a context, project ID, and API key parameters, then returns the created API key or an error.
func (s *APIKeysService) Create(ctx context.Context, projectID string, params CreateAPIKeyParams) (*APIKey, error) {
	path := fmt.Sprintf("/api/projects/%s/api-keys", projectID)
	return doRequest[CreateAPIKeyParams, APIKey](s.client, ctx, "POST", path, &params)
}

// Revoke deletes an API key
// It sends a DELETE request to the api-keyAPI endpoint with the given projectID and key.
// Returns an error if the request fails.
func (s *APIKeysService) Revoke(ctx context.Context, projectID, keyID string) error {
	path := fmt.Sprintf("/api/projects/%s/api-keys/%s", projectID, keyID)
	_, err := doRequest[any, any](s.client, ctx, "DELETE", path, nil)
	return err
}

// GetStats retrieve usage statistics for a specific API key.
// It takes a context, project ID, and API key parameters, then returns the created API key or an error.
func (s *APIKeysService) GetStats(ctx context.Context, projectID, keyID string) (*KeyUsageStats, error) {
	path := fmt.Sprintf("/api/projects/%s/api-keys/%s/stats", projectID, keyID)
	return doRequest[any, KeyUsageStats](s.client, ctx, "GET", path, nil)
}
