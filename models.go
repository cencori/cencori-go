package cencori

type ChatParams struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int      `json:"maxTokens,omitempty"`
	Stream      bool      `json:"stream"`
	UserID      *string   `json:"userId,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Content      string  `json:"content"`
	Model        string  `json:"model"`
	Provider     string  `json:"provider"`
	Usage        Usage   `json:"usage"`
	CostUSD      float64 `json:"cost_usd"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamChunk struct {
	Delta        string `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
	Err          error  `json"-"`
}
