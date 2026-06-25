package cencori

import (
	"context"
	"fmt"
)

// ComputeService provides serverless functions and GPU access.
// 🚧 Coming Soon
type ComputeService struct {
	client *Client
}

// Run executes a serverless function.
// 🚧 Coming Soon — currently returns an error.
func (s *ComputeService) Run(ctx context.Context, functionName string, input map[string]interface{}) error {
	return fmt.Errorf("cencori: Compute is coming soon")
}

// WorkflowService provides AI pipelines and orchestration.
// 🚧 Coming Soon
type WorkflowService struct {
	client *Client
}

// Trigger triggers a workflow pipeline.
// 🚧 Coming Soon — currently returns an error.
func (s *WorkflowService) Trigger(ctx context.Context, pipelineID string, data map[string]interface{}) error {
	return fmt.Errorf("cencori: Workflow is coming soon")
}

// StorageService provides vector database, knowledge base, and RAG storage.
// 🚧 Coming Soon
type StorageService struct {
	client *Client
	Vectors *VectorsSubmodule
}

// VectorsSubmodule provides vector search and upsert operations.
// 🚧 Coming Soon
type VectorsSubmodule struct {
	client *Client
}

// Search performs a vector search.
// 🚧 Coming Soon — currently returns an error.
func (s *VectorsSubmodule) Search(ctx context.Context, query string, opts map[string]interface{}) error {
	return fmt.Errorf("cencori: Vector search is coming soon")
}

// Upsert upserts vectors into the store.
// 🚧 Coming Soon — currently returns an error.
func (s *VectorsSubmodule) Upsert(ctx context.Context, vectors []map[string]interface{}) error {
	return fmt.Errorf("cencori: Vector upsert is coming soon")
}
