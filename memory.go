package cencori

import (
	"context"
	"fmt"
)

// MemoryService provides methods for vector storage, semantic search, and RAG namespaces.
type MemoryService struct {
	client *Client
}

// CreateNamespace creates a new memory namespace for storing vector memories.
func (s *MemoryService) CreateNamespace(ctx context.Context, options *CreateNamespaceOptions) (*MemoryNamespace, error) {
	return doRequest[CreateNamespaceOptions, MemoryNamespace](s.client, ctx, "POST", "/api/memory/namespaces", options)
}

// ListNamespaces lists all memory namespaces for the project.
func (s *MemoryService) ListNamespaces(ctx context.Context) ([]MemoryNamespace, error) {
	resp, err := doRequest[any, ListNamespacesResponse](s.client, ctx, "GET", "/api/memory/namespaces", nil)
	if err != nil {
		return nil, err
	}
	return resp.Namespaces, nil
}

// Store stores a memory in a namespace.
func (s *MemoryService) Store(ctx context.Context, options *StoreMemoryOptions) (*Memory, error) {
	return doRequest[StoreMemoryOptions, Memory](s.client, ctx, "POST", "/api/memory/store", options)
}

// Search performs semantic search across memories in a namespace.
func (s *MemoryService) Search(ctx context.Context, options *SearchMemoryOptions) (*SearchResult, error) {
	return doRequest[SearchMemoryOptions, SearchResult](s.client, ctx, "POST", "/api/memory/search", options)
}

// Get retrieves a memory by ID.
func (s *MemoryService) Get(ctx context.Context, id string) (*Memory, error) {
	path := fmt.Sprintf("/api/memory/%s", id)
	return doRequest[any, Memory](s.client, ctx, "GET", path, nil)
}

// Delete deletes a memory by ID.
func (s *MemoryService) Delete(ctx context.Context, id string) (*MemoryDeleteResponse, error) {
	path := fmt.Sprintf("/api/memory/%s", id)
	return doRequest[any, MemoryDeleteResponse](s.client, ctx, "DELETE", path, nil)
}

// DeleteByFilter deletes all memories in a namespace matching a filter.
// This is a client-side implementation that searches then deletes.
func (s *MemoryService) DeleteByFilter(ctx context.Context, namespace string, filter map[string]interface{}) (*MemoryDeleteCount, error) {
	searchResult, err := s.Search(ctx, &SearchMemoryOptions{
		Namespace: namespace,
		Query:     "*",
		Limit:     intPtr(1000),
		Threshold: float64Ptr(0),
		Filter:    filter,
	})
	if err != nil {
		return nil, err
	}

	deleted := 0
	for _, mem := range searchResult.Results {
		_, err := s.Delete(ctx, mem.ID)
		if err != nil {
			return nil, err
		}
		deleted++
	}

	return &MemoryDeleteCount{Deleted: deleted}, nil
}

// StoreBatch stores multiple memories in a namespace.
// This is a client-side implementation that calls Store for each item.
func (s *MemoryService) StoreBatch(ctx context.Context, namespace string, items []StoreBatchItem) ([]Memory, error) {
	results := make([]Memory, 0, len(items))
	for _, item := range items {
		mem, err := s.Store(ctx, &StoreMemoryOptions{
			Namespace: namespace,
			Content:   item.Content,
			Metadata:  item.Metadata,
		})
		if err != nil {
			return nil, err
		}
		results = append(results, *mem)
	}
	return results, nil
}

// Convenience helpers.

type StoreBatchItem struct {
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type MemoryDeleteResponse struct {
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
}

type MemoryDeleteCount struct {
	Deleted int `json:"deleted"`
}

func intPtr(i int) *int { return &i }

func float64Ptr(f float64) *float64 { return &f }
