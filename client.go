package cencori

import (
	"errors"
	"net/http"
	"time"
)

type ClientOptions struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

func WithAPIKey(apiKey string) Option {
	return func(c *ClientOptions) { c.APIKey = apiKey }
}

func WithBaseURL(baseURL string) Option {
	return func(c *ClientOptions) { c.BaseURL = baseURL }
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *ClientOptions) { c.Timeout = timeout }
}

type Client struct {
	APIKey     string
	BaseURL    string
	httpClient *http.Client

	Chat      *ChatService
	Vision    *VisionService
	Voice     *VoiceService
	Documents *DocumentsService
	Projects  *ProjectsService
	APIKeys   *APIKeysService
	Metrics   *MetricsService
	Agents    *AgentsService
	Memory    *MemoryService
	Sessions  *SessionsService
	Telemetry *TelemetryService
	Compute   *ComputeService
	Workflow  *WorkflowService
	Storage   *StorageService
}

type Option func(*ClientOptions)

func NewClient(opts ...Option) (*Client, error) {
	config := &ClientOptions{
		BaseURL: "https://cencori.com",
		Timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(config)
	}

	if config.APIKey == "" {
		return nil, errors.New("you need a valid API Key to use this client")
	}

	c := &Client{
		APIKey:  config.APIKey,
		BaseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}

	c.Chat = &ChatService{client: c}
	c.Vision = &VisionService{client: c}
	c.Voice = &VoiceService{client: c}
	c.Documents = &DocumentsService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.Metrics = &MetricsService{client: c}
	c.Agents = &AgentsService{client: c}
	c.Memory = &MemoryService{client: c}
	c.Sessions = &SessionsService{client: c}
	c.Telemetry = &TelemetryService{client: c}
	c.Compute = &ComputeService{client: c}
	c.Workflow = &WorkflowService{client: c}
	c.Storage = &StorageService{client: c, Vectors: &VectorsSubmodule{client: c}}

	return c, nil
}
