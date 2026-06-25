package cencori

import "time"

// Shared Components.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Stats struct {
	TotalRequests int        `json:"total_requests"`
	TotalCostUSD  float64    `json:"total_cost_usd"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
}

// ── Chat Models ──

type MessageContentType string

const (
	MessageContentTypeText     MessageContentType = "text"
	MessageContentTypeImageURL MessageContentType = "image_url"
	MessageContentTypeImage    MessageContentType = "image"
)

type MessageContent struct {
	Type     MessageContentType `json:"type"`
	Text     string             `json:"text,omitempty"`
	ImageURL *ImageURLContent   `json:"image_url,omitempty"`
	Image    string             `json:"image,omitempty"`
}

type ImageURLContent struct {
	URL    string  `json:"url"`
	Detail *string `json:"detail,omitempty"`
}

type Message struct {
	Role       string            `json:"role"`
	Content    string            `json:"content,omitempty"`
	Contents   []MessageContent  `json:"contents,omitempty"`
	ToolCallID *string           `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall        `json:"tool_calls,omitempty"`
}

type ChatParams struct {
	Model       string          `json:"model"`
	Messages    []Message       `json:"messages"`
	Temperature *float64        `json:"temperature,omitempty"`
	MaxTokens   *int            `json:"max_tokens,omitempty"`
	TopP        *float64        `json:"top_p,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
	User        *string         `json:"user,omitempty"`
	Tools       []ToolDefinition `json:"tools,omitempty"`
	ToolChoice  *ToolChoice     `json:"tool_choice,omitempty"`
	Prompt      *PromptReference `json:"prompt,omitempty"`
}

type PromptReference struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables,omitempty"`
}

type ChatResponse struct {
	ID           string `json:"id"`
	Object       string `json:"object"`
	Created      int64  `json:"created"`
	Model        string `json:"model"`
	Content      string `json:"content,omitempty"`
	FinishReason string `json:"finish_reason,omitempty"`
	ToolCalls    []ToolCall `json:"toolCalls,omitempty"`
	Choices      []struct {
		Index        int     `json:"index"`
		Message      struct {
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices,omitempty"`
	Usage Usage `json:"usage"`
}

// ── Tool / Function Calling Models ──

type ToolDefinition struct {
	Type     string             `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Parameters  map[string]interface{}   `json:"parameters,omitempty"`
}

type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolChoice struct {
	Type     string              `json:"type,omitempty"`
	Function *ToolChoiceFunction `json:"function,omitempty"`
}

type ToolChoiceFunction struct {
	Name string `json:"name"`
}

func ToolChoiceAuto() *ToolChoice {
	return &ToolChoice{Type: "auto"}
}

func ToolChoiceNone() *ToolChoice {
	return &ToolChoice{Type: "none"}
}

func ToolChoiceRequired() *ToolChoice {
	return &ToolChoice{Type: "required"}
}

func ToolChoiceFunctionName(name string) *ToolChoice {
	return &ToolChoice{
		Type:     "function",
		Function: &ToolChoiceFunction{Name: name},
	}
}

// ── Completions Models ──

type CompletionParams struct {
	Prompt      string   `json:"prompt"`
	Model       string   `json:"model"`
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
}

// ── Embedding Models ──

type EmbeddingParams struct {
	Input interface{} `json:"input"` // string or []string
	Model string      `json:"model"`
}

type EmbeddingUsage struct {
	TotalTokens int `json:"total_tokens"`
}

type EmbeddingData struct {
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Model  string          `json:"model"`
	Data   []EmbeddingData `json:"data"`
	Usage  EmbeddingUsage  `json:"usage"`
	Object string          `json:"object"`
}

// ── Stream Response ──

type StreamChunk struct {
	ID      string         `json:"id,omitempty"`
	Object  string         `json:"object,omitempty"`
	Created int64          `json:"created,omitempty"`
	Model   string         `json:"model,omitempty"`
	Choices []StreamChoice `json:"choices,omitempty"`
	Delta   string         `json:"delta,omitempty"`
	Err     error          `json:"-"`
}

type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type StreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// ── Generate Object (Structured Output) Models ──

type GenerateObjectRequest struct {
	Model             string                 `json:"model"`
	Prompt            string                 `json:"prompt,omitempty"`
	Messages          []Message              `json:"messages,omitempty"`
	Schema            map[string]interface{} `json:"schema"`
	SchemaName        string                 `json:"schemaName,omitempty"`
	SchemaDescription string                 `json:"schemaDescription,omitempty"`
	Temperature       *float64               `json:"temperature,omitempty"`
	MaxTokens         *int                   `json:"max_tokens,omitempty"`
}

type GenerateObjectResponse struct {
	Object  map[string]interface{} `json:"object"`
	Usage   Usage                  `json:"usage"`
}

// ── Image Generation Models ──

type ImageGenerationRequest struct {
	Prompt         string   `json:"prompt"`
	Model          string   `json:"model"`
	N              *int     `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"`
	Quality        string   `json:"quality,omitempty"`
	Style          string   `json:"style,omitempty"`
	ResponseFormat string   `json:"responseFormat,omitempty"`
}

type GeneratedImage struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revisedPrompt,omitempty"`
}

type ImageGenerationResponse struct {
	Images   []GeneratedImage `json:"images"`
	Model    string           `json:"model"`
	Provider string           `json:"provider"`
}

// ── RAG Models ──

type RagRequest struct {
	Model           string    `json:"model"`
	Messages        []Message `json:"messages"`
	Namespace       string    `json:"namespace"`
	Temperature     *float64  `json:"temperature,omitempty"`
	MaxTokens       *int      `json:"max_tokens,omitempty"`
	Limit           int       `json:"limit,omitempty"`
	Threshold       float64   `json:"threshold,omitempty"`
	IncludeSources  bool      `json:"include_sources,omitempty"`
	Stream          bool      `json:"stream,omitempty"`
}

type RagSource struct {
	Content    string                 `json:"content"`
	Metadata   map[string]interface{} `json:"metadata"`
	Similarity float64                `json:"similarity"`
}

type RagResponse struct {
	Message   Message     `json:"message"`
	Model     string      `json:"model"`
	Provider  string      `json:"provider"`
	Usage     Usage       `json:"usage"`
	Sources   []RagSource `json:"sources,omitempty"`
	LatencyMs int         `json:"latency_ms"`
}

type RagStreamChunk struct {
	Type         string      `json:"type"`
	Delta        string      `json:"delta,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
	Sources      []RagSource `json:"sources,omitempty"`
	Err          error       `json:"-"`
}

// ── Responses API Models ──

type ResponseInputItem struct {
	Type      string                 `json:"type"`
	Role      string                 `json:"role,omitempty"`
	Content   string                 `json:"content,omitempty"`
	CallID    string                 `json:"call_id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Arguments string                 `json:"arguments,omitempty"`
	Output    string                 `json:"output,omitempty"`
	Status    string                 `json:"status,omitempty"`
	Filename  string                 `json:"filename,omitempty"`
	MimeType  string                 `json:"mime_type,omitempty"`
}

type ResponsesTool struct {
	Type               string                 `json:"type"`
	SearchContextSize  string                 `json:"search_context_size,omitempty"`
	MaxNumResults      *int                   `json:"max_num_results,omitempty"`
	Filters            map[string]interface{} `json:"filters,omitempty"`
	Function           *ToolFunction          `json:"function,omitempty"`
	UserLocation       map[string]interface{} `json:"user_location,omitempty"`
}

type ResponsesRequest struct {
	Model              string                 `json:"model"`
	Input              interface{}            `json:"input"` // string or []ResponseInputItem
	Instructions       string                 `json:"instructions,omitempty"`
	Tools              []ResponsesTool        `json:"tools,omitempty"`
	ToolChoice         interface{}            `json:"tool_choice,omitempty"`
	Temperature        *float64               `json:"temperature,omitempty"`
	MaxOutputTokens    *int                   `json:"max_output_tokens,omitempty"`
	TopP               *float64               `json:"top_p,omitempty"`
	Store              *bool                  `json:"store,omitempty"`
	Metadata           map[string]string      `json:"metadata,omitempty"`
	PreviousResponseID string                 `json:"previous_response_id,omitempty"`
	ParallelToolCalls  *bool                  `json:"parallel_tool_calls,omitempty"`
	Truncation         string                 `json:"truncation,omitempty"`
	ResponseFormat     map[string]interface{} `json:"response_format,omitempty"`
	Include            []string               `json:"include,omitempty"`
	Stream             bool                   `json:"stream,omitempty"`
	User               string                 `json:"user,omitempty"`
}

type ResponsesResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Output  []ResponsesOutputItem  `json:"output"`
	Usage   ResponsesUsage         `json:"usage"`
	Status  string                 `json:"status"`
	Metadata map[string]string     `json:"metadata,omitempty"`
}

type ResponsesOutputItem struct {
	ID      string                    `json:"id"`
	Type    string                    `json:"type"`
	Status  string                    `json:"status,omitempty"`
	Role    string                    `json:"role,omitempty"`
	Content []ResponseContentPart     `json:"content,omitempty"`
	CallID  string                    `json:"call_id,omitempty"`
	Name    string                    `json:"name,omitempty"`
	Arguments string                  `json:"arguments,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

type ResponseContentPart struct {
	Type        string        `json:"type"`
	Text        string        `json:"text,omitempty"`
	Annotations []interface{} `json:"annotations,omitempty"`
}

type ResponsesUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type ResponsesStreamEvent struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
	Err  error                  `json:"-"`
}

// --- Project Models ---.
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Visibility  string    `json:"visibility"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Stats       *Stats    `json:"stats,omitempty"` // Included in "Get" only
}

type CreateProjectParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Visibility  string `json:"visibility,omitempty"` // "public" | "private"
}

// --- API Key Models ---.

type APIKey struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Prefix      string     `json:"prefix,omitempty"`
	Key         string     `json:"key,omitempty"` // Only present on Create
	Environment string     `json:"environment"`
	CreatedAt   time.Time  `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	UsageCount  int        `json:"usage_count,omitempty"`
}

type CreateAPIKeyParams struct {
	Name        string `json:"name"`
	Environment string `json:"environment"` // "production" | "test"
}

type KeyUsageStats struct {
	KeyID           string         `json:"key_id"`
	TotalRequests   int            `json:"total_requests"`
	TotalCostUSD    float64        `json:"total_cost_usd"`
	LastUsedAt      time.Time      `json:"last_used_at"`
	RequestsByDay   []DailyStat    `json:"requests_by_day"`
	RequestsByModel map[string]int `json:"requests_by_model"`
}

type DailyStat struct {
	Date    string  `json:"date"`
	Count   int     `json:"count"`
	CostUSD float64 `json:"cost_usd"`
}

// Metric Models.
type MetricsResponse struct {
	Period    string               `json:"period"`
	StartDate time.Time            `json:"start_date"`
	EndDate   time.Time            `json:"end_date"`
	Requests  RequestMetrics       `json:"requests"`
	Cost      CostMetrics          `json:"cost"`
	Tokens    TokenMetrics         `json:"tokens"`
	Latency   LatencyMetrics       `json:"latency"`
	Providers map[string]Breakdown `json:"providers"`
	Models    map[string]Breakdown `json:"models"`
}

type RequestMetrics struct {
	Total       int     `json:"total"`
	Success     int     `json:"success"`
	Error       int     `json:"error"`
	Filtered    int     `json:"filtered"`
	SuccessRate float64 `json:"success_rate"`
}

type CostMetrics struct {
	TotalUSD             float64 `json:"total_usd"`
	AveragePerRequestUSD float64 `json:"average_per_request_usd"`
}

type TokenMetrics struct {
	Prompt     int `json:"prompt"`
	Completion int `json:"completion"`
	Total      int `json:"total"`
}

type LatencyMetrics struct {
	AvgMS int `json:"avg_ms"`
	P50MS int `json:"p50_ms"`
	P90MS int `json:"p90_ms"`
	P99MS int `json:"p99_ms"`
}

type Breakdown struct {
	Requests int     `json:"requests"`
	CostUSD  float64 `json:"cost_usd"`
}

// ── Agent Models ──

type AgentConfig struct {
	Model       string   `json:"model"`
	SystemPrompt *string `json:"system_prompt,omitempty"`
	Tools       []string `json:"tools,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
}

type Agent struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	IsActive    bool        `json:"is_active"`
	ShadowMode  bool        `json:"shadow_mode"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   *string     `json:"updated_at,omitempty"`
	Config      AgentConfig `json:"config"`
}

type AgentListItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
	ShadowMode  bool    `json:"shadow_mode"`
	CreatedAt   string  `json:"created_at"`
}

type CreateAgentParams struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Config      *AgentConfig `json:"config,omitempty"`
}

type UpdateAgentParams struct {
	Name        *string       `json:"name,omitempty"`
	Description *string       `json:"description,omitempty"`
	IsActive    *bool         `json:"is_active,omitempty"`
	ShadowMode  *bool         `json:"shadow_mode,omitempty"`
	Config      *AgentConfig  `json:"config,omitempty"`
}

type AgentKey struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	KeyPrefix   string  `json:"key_prefix"`
	FullKey     *string `json:"full_key,omitempty"`
	Environment string  `json:"environment"`
	KeyType     string  `json:"key_type"`
	AgentID     string  `json:"agent_id"`
	CreatedAt   string  `json:"created_at"`
}

type CreateAgentKeyParams struct {
	Name            string   `json:"name,omitempty"`
	Environment     string   `json:"environment,omitempty"`
	KeyType         string   `json:"key_type,omitempty"`
	AllowedDomains  []string `json:"allowed_domains,omitempty"`
}

type AgentListResponse struct {
	Data []AgentListItem `json:"data"`
}

// ── Memory Models ──

type MemoryNamespace struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description,omitempty"`
	EmbeddingModel string                 `json:"embeddingModel"`
	Dimensions     int                    `json:"dimensions"`
	Metadata       map[string]interface{} `json:"metadata"`
	MemoryCount    *int                   `json:"memoryCount,omitempty"`
	CreatedAt      string                 `json:"createdAt"`
}

type CreateNamespaceOptions struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description,omitempty"`
	EmbeddingModel  string                 `json:"embeddingModel,omitempty"`
	Dimensions      *int                   `json:"dimensions,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type Memory struct {
	ID         string                 `json:"id"`
	Namespace  string                 `json:"namespace"`
	Content    string                 `json:"content"`
	Metadata   map[string]interface{} `json:"metadata"`
	Similarity *float64               `json:"similarity,omitempty"`
	ExpiresAt  *string                `json:"expiresAt,omitempty"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  *string                `json:"updatedAt,omitempty"`
}

type StoreMemoryOptions struct {
	Namespace string                 `json:"namespace"`
	Content   string                 `json:"content"`
	Embedding []float64              `json:"embedding,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ExpiresAt *string                `json:"expiresAt,omitempty"`
}

type SearchMemoryOptions struct {
	Namespace string                 `json:"namespace"`
	Query     string                 `json:"query"`
	Limit     *int                   `json:"limit,omitempty"`
	Threshold *float64               `json:"threshold,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
}

type SearchResult struct {
	Results   []Memory `json:"results"`
	Query     string   `json:"query"`
	Namespace string   `json:"namespace"`
	Count     int      `json:"count"`
	LatencyMs int      `json:"latencyMs"`
}

type ListNamespacesResponse struct {
	Namespaces []MemoryNamespace `json:"namespaces"`
}

// ── Session Models ──

type Session struct {
	ID        string                 `json:"id"`
	Status    string                 `json:"status"`
	TurnCount int                    `json:"turn_count"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	AgentID   *string                `json:"agent_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
	TotalCost float64                `json:"total_cost"`
}

type SessionEvent struct {
	ID         string                 `json:"id"`
	SessionID  string                 `json:"session_id"`
	TurnNumber int                    `json:"turn_number"`
	Sequence   int                    `json:"sequence"`
	EventType  string                 `json:"event_type"`
	Payload    map[string]interface{} `json:"payload"`
	CreatedAt  string                 `json:"created_at"`
}

type CreateSessionParams struct {
	AgentID  *string                `json:"agent_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type TurnParams struct {
	Input            interface{}               `json:"input"` // string or []map[string]interface{}
	Tools            []map[string]interface{}   `json:"tools,omitempty"`
	Instructions     string                     `json:"instructions,omitempty"`
	AgentID          string                     `json:"agent_id,omitempty"`
	Model            string                     `json:"model,omitempty"`
	Temperature      *float64                   `json:"temperature,omitempty"`
	MaxOutputTokens  *int                       `json:"max_output_tokens,omitempty"`
	ToolChoice       interface{}                `json:"tool_choice,omitempty"`
	ResponseFormat   map[string]interface{}     `json:"response_format,omitempty"`
	User             string                     `json:"user,omitempty"`
	PauseOnToolCalls *bool                      `json:"pause_on_tool_calls,omitempty"`
}

type ApproveRejectParams struct {
	ActionID    string              `json:"action_id"`
	ToolResults []ToolActionResult  `json:"tool_results,omitempty"`
}

type ToolActionResult struct {
	ActionID string `json:"action_id"`
	Output   string `json:"output"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type SessionListParams struct {
	Page    *int    `json:"page,omitempty"`
	Limit   *int    `json:"limit,omitempty"`
	Status  string  `json:"status,omitempty"`
	AgentID string  `json:"agent_id,omitempty"`
}

type ApproveResponse struct {
	ID         string `json:"id"`
	ActionID   string `json:"action_id"`
	Resolution string `json:"resolution"`
	Status     string `json:"status"`
}

// ── Telemetry Models ──

type WebTelemetryPayload struct {
	Host        string  `json:"host"`
	Method      string  `json:"method"`
	Path        string  `json:"path"`
	StatusCode  int     `json:"statusCode"`
	RequestID   string  `json:"requestId,omitempty"`
	QueryString string  `json:"queryString,omitempty"`
	Message     string  `json:"message,omitempty"`
	UserAgent   string  `json:"userAgent,omitempty"`
	Referer     string  `json:"referer,omitempty"`
	IPAddress   string  `json:"ipAddress,omitempty"`
	CountryCode string  `json:"countryCode,omitempty"`
	LatencyMs   *int    `json:"latencyMs,omitempty"`
}
