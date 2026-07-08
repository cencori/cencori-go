package cencori

import (
	"context"
	"errors"
)

// DocumentsService provides methods for extracting text from PDFs and images,
// summarizing, and answering questions about documents.
type DocumentsService struct {
	client *Client
}

// DocumentParams is the request payload for any Documents endpoint.
// Provide one of DocumentURL or DocumentBase64+MimeType.
type DocumentParams struct {
	DocumentURL    string  `json:"document_url,omitempty"`
	DocumentBase64 string  `json:"document_base64,omitempty"`
	MimeType       string  `json:"mime_type,omitempty"`
	Filename       string  `json:"filename,omitempty"`
	Prompt         string  `json:"prompt,omitempty"`
	Model          string  `json:"model,omitempty"`
	MaxTokens      int     `json:"max_tokens,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
	// Question is only used by Query.
	Question string `json:"question,omitempty"`
}

type DocumentUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}

type DocumentCost struct {
	ProviderCostUsd  float64 `json:"providerCostUsd"`
	CencoriChargeUsd float64 `json:"cencoriChargeUsd"`
	MarkupPercentage float64 `json:"markupPercentage"`
}

// DocumentExtractResponse is the response from /api/ai/documents/extract.
type DocumentExtractResponse struct {
	Text      string         `json:"text"`
	PageCount int            `json:"pageCount"`
	Kind      string         `json:"kind"`
	Method    string         `json:"method"`
	Model     string         `json:"model,omitempty"`
	Provider  string         `json:"provider,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Usage     *DocumentUsage `json:"usage,omitempty"`
	Cost      *DocumentCost  `json:"cost,omitempty"`
}

// DocumentSummarizeResponse is the response from /api/ai/documents/summarize.
type DocumentSummarizeResponse struct {
	Summary       string        `json:"summary"`
	Text          string        `json:"text"`
	PageCount     int           `json:"pageCount"`
	Model         string        `json:"model"`
	Provider      string        `json:"provider"`
	ExtractMethod string        `json:"extractMethod"`
	Usage         DocumentUsage `json:"usage"`
	Cost          DocumentCost  `json:"cost"`
}

// DocumentQueryResponse is the response from /api/ai/documents/query.
type DocumentQueryResponse struct {
	Answer        string        `json:"answer"`
	Question      string        `json:"question"`
	PageCount     int           `json:"pageCount"`
	Model         string        `json:"model"`
	Provider      string        `json:"provider"`
	ExtractMethod string        `json:"extractMethod"`
	Usage         DocumentUsage `json:"usage"`
	Cost          DocumentCost  `json:"cost"`
}

// Extract extracts text from a PDF (native) or image (via Vision OCR).
func (s *DocumentsService) Extract(ctx context.Context, params *DocumentParams) (*DocumentExtractResponse, error) {
	if err := validateDocumentParams(params); err != nil {
		return nil, err
	}
	return doRequest[DocumentParams, DocumentExtractResponse](s.client, ctx, "POST", "/api/ai/documents/extract", params)
}

// Summarize extracts and then summarizes the document.
func (s *DocumentsService) Summarize(ctx context.Context, params *DocumentParams) (*DocumentSummarizeResponse, error) {
	if err := validateDocumentParams(params); err != nil {
		return nil, err
	}
	return doRequest[DocumentParams, DocumentSummarizeResponse](s.client, ctx, "POST", "/api/ai/documents/summarize", params)
}

// Query extracts and then answers a question about the document.
// params.Question is required.
func (s *DocumentsService) Query(ctx context.Context, params *DocumentParams) (*DocumentQueryResponse, error) {
	if err := validateDocumentParams(params); err != nil {
		return nil, err
	}
	if params.Question == "" {
		return nil, errors.New("documents.Query requires Question")
	}
	return doRequest[DocumentParams, DocumentQueryResponse](s.client, ctx, "POST", "/api/ai/documents/query", params)
}

func validateDocumentParams(params *DocumentParams) error {
	if params == nil {
		return errors.New("documents params required")
	}
	if params.DocumentURL == "" && params.DocumentBase64 == "" {
		return errors.New("documents request requires DocumentURL or DocumentBase64")
	}
	return nil
}
