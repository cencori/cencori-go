package cencori

import (
	"context"
	"errors"
)

// WebService provides first-party web retrieval, extraction, crawling, and search.
type WebService struct {
	client *Client
}

type WebFetchParams struct {
	URL       string `json:"url"`
	MaxBytes  int    `json:"maxBytes,omitempty"`
	TimeoutMs int    `json:"timeoutMs,omitempty"`
}

type WebResponseHeaders struct {
	CacheControl string `json:"cacheControl,omitempty"`
	ETag         string `json:"etag,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

type WebFetchResponse struct {
	URL         string             `json:"url"`
	FinalURL    string             `json:"finalUrl"`
	StatusCode  int                `json:"statusCode"`
	MimeType    string             `json:"mimeType"`
	Body        string             `json:"body"`
	Bytes       int                `json:"bytes"`
	ContentHash string             `json:"contentHash"`
	RetrievedAt string             `json:"retrievedAt"`
	Headers     WebResponseHeaders `json:"headers"`
	Untrusted   bool               `json:"untrusted"`
}

type WebLink struct {
	URL      string   `json:"url"`
	Text     string   `json:"text"`
	Rel      []string `json:"rel"`
	Internal bool     `json:"internal"`
}

type WebEvidenceSpan struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type WebExtractResponse struct {
	URL           string            `json:"url"`
	CanonicalURL  string            `json:"canonicalUrl"`
	Title         string            `json:"title"`
	Description   string            `json:"description,omitempty"`
	Language      string            `json:"language,omitempty"`
	Content       string            `json:"content"`
	ContentHash   string            `json:"contentHash"`
	MimeType      string            `json:"mimeType"`
	StatusCode    int               `json:"statusCode"`
	RetrievedAt   string            `json:"retrievedAt"`
	PublishedAt   string            `json:"publishedAt,omitempty"`
	ModifiedAt    string            `json:"modifiedAt,omitempty"`
	Links         []WebLink         `json:"links"`
	EvidenceSpans []WebEvidenceSpan `json:"evidenceSpans"`
	Metadata      map[string]string `json:"metadata"`
	Untrusted     bool              `json:"untrusted"`
}

type WebSearchParams struct {
	Query     string `json:"query"`
	Limit     int    `json:"limit,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Freshness string `json:"freshness,omitempty"`
}

type WebSearchEvidence struct {
	Quote       string `json:"quote"`
	ContentHash string `json:"contentHash"`
	RetrievedAt string `json:"retrievedAt"`
}

type WebSearchHit struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	URL          string            `json:"url"`
	CanonicalURL string            `json:"canonicalUrl"`
	Snippet      string            `json:"snippet"`
	Score        float64           `json:"score"`
	ContentHash  string            `json:"contentHash"`
	RetrievedAt  string            `json:"retrievedAt"`
	PublishedAt  string            `json:"publishedAt,omitempty"`
	Evidence     WebSearchEvidence `json:"evidence"`
}

type WebSearchResponse struct {
	Query        string         `json:"query"`
	Results      []WebSearchHit `json:"results"`
	Count        int            `json:"count"`
	SearchEngine string         `json:"searchEngine"`
}

type WebCrawlParams struct {
	Seeds      []string `json:"seeds"`
	MaxPages   int      `json:"maxPages,omitempty"`
	MaxDepth   int      `json:"maxDepth,omitempty"`
	SameOrigin *bool    `json:"sameOrigin,omitempty"`
}

type WebCrawlPage struct {
	URL        string `json:"url"`
	Status     string `json:"status"`
	DocumentID string `json:"documentId,omitempty"`
	Error      string `json:"error,omitempty"`
}

type WebCrawlResponse struct {
	Pages      []WebCrawlPage `json:"pages"`
	Indexed    int            `json:"indexed"`
	Failed     int            `json:"failed"`
	Discovered int            `json:"discovered"`
}

func (s *WebService) Fetch(ctx context.Context, params *WebFetchParams) (*WebFetchResponse, error) {
	if params == nil || params.URL == "" {
		return nil, errors.New("web.Fetch requires URL")
	}
	return doRequest[WebFetchParams, WebFetchResponse](s.client, ctx, "POST", "/api/v1/web/fetch", params)
}

func (s *WebService) Extract(ctx context.Context, params *WebFetchParams) (*WebExtractResponse, error) {
	if params == nil || params.URL == "" {
		return nil, errors.New("web.Extract requires URL")
	}
	return doRequest[WebFetchParams, WebExtractResponse](s.client, ctx, "POST", "/api/v1/web/extract", params)
}

func (s *WebService) Search(ctx context.Context, params *WebSearchParams) (*WebSearchResponse, error) {
	if params == nil || params.Query == "" {
		return nil, errors.New("web.Search requires Query")
	}
	return doRequest[WebSearchParams, WebSearchResponse](s.client, ctx, "POST", "/api/v1/web/search", params)
}

func (s *WebService) Crawl(ctx context.Context, params *WebCrawlParams) (*WebCrawlResponse, error) {
	if params == nil || len(params.Seeds) == 0 {
		return nil, errors.New("web.Crawl requires at least one seed URL")
	}
	return doRequest[WebCrawlParams, WebCrawlResponse](s.client, ctx, "POST", "/api/v1/web/crawl", params)
}
