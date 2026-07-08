package cencori

import (
	"context"
	"errors"
)

// VisionService provides methods for image analysis, description, OCR, and classification
// via the Cencori Vision API.
type VisionService struct {
	client *Client
}

// VisionParams is the request payload for any Vision endpoint.
// Provide exactly one of ImageURL or ImageBase64 (with MimeType).
type VisionParams struct {
	// ImageURL is an http(s):// or data: URL pointing to the image.
	ImageURL string `json:"image_url,omitempty"`
	// ImageBase64 is the raw base64 encoding of the image bytes. Requires MimeType.
	ImageBase64 string `json:"image_base64,omitempty"`
	// MimeType is required when using ImageBase64 (e.g. "image/png").
	MimeType string `json:"mime_type,omitempty"`
	// Prompt overrides the default task prompt. Only meaningful for Analyze.
	Prompt string `json:"prompt,omitempty"`
	// Model overrides the default (gpt-4o-mini). See GET /api/ai/vision for the list.
	Model string `json:"model,omitempty"`
	// MaxTokens caps the response length.
	MaxTokens int `json:"max_tokens,omitempty"`
	// Temperature is the sampling temperature (0–2).
	Temperature float64 `json:"temperature,omitempty"`
	// ResponseFormat is "text" or "json".
	ResponseFormat string `json:"response_format,omitempty"`
}

// VisionUsage reports token usage for a Vision request.
type VisionUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}

// VisionCost reports provider and Cencori-billed cost.
type VisionCost struct {
	ProviderCostUsd    float64 `json:"providerCostUsd"`
	CencoriChargeUsd   float64 `json:"cencoriChargeUsd"`
	MarkupPercentage   float64 `json:"markupPercentage"`
}

// VisionResponse is the response from the general /api/ai/vision endpoint.
type VisionResponse struct {
	Analysis string      `json:"analysis"`
	Model    string      `json:"model"`
	Provider string      `json:"provider"`
	Usage    VisionUsage `json:"usage"`
	Cost     VisionCost  `json:"cost"`
}

// VisionDescribeResponse is the response from /api/ai/vision/describe.
type VisionDescribeResponse struct {
	Description string      `json:"description"`
	Model       string      `json:"model"`
	Provider    string      `json:"provider"`
	Usage       VisionUsage `json:"usage"`
	Cost        VisionCost  `json:"cost"`
}

// VisionOCRResponse is the response from /api/ai/vision/ocr.
type VisionOCRResponse struct {
	Text     string      `json:"text"`
	Model    string      `json:"model"`
	Provider string      `json:"provider"`
	Usage    VisionUsage `json:"usage"`
	Cost     VisionCost  `json:"cost"`
}

// VisionClassifyResponse is the response from /api/ai/vision/classify.
// Classification is a map when the model returns valid JSON, otherwise the raw string is in Raw.
type VisionClassifyResponse struct {
	Classification map[string]any `json:"classification"`
	Raw            string         `json:"raw"`
	Model          string         `json:"model"`
	Provider       string         `json:"provider"`
	Usage          VisionUsage    `json:"usage"`
	Cost           VisionCost     `json:"cost"`
}

// Analyze runs a general vision analysis with an optional custom prompt.
func (s *VisionService) Analyze(ctx context.Context, params *VisionParams) (*VisionResponse, error) {
	if err := validateVisionParams(params); err != nil {
		return nil, err
	}
	return doRequest[VisionParams, VisionResponse](s.client, ctx, "POST", "/api/ai/vision", params)
}

// Describe returns a rich description of the image.
func (s *VisionService) Describe(ctx context.Context, params *VisionParams) (*VisionDescribeResponse, error) {
	if err := validateVisionParams(params); err != nil {
		return nil, err
	}
	return doRequest[VisionParams, VisionDescribeResponse](s.client, ctx, "POST", "/api/ai/vision/describe", params)
}

// OCR extracts all text visible in the image.
func (s *VisionService) OCR(ctx context.Context, params *VisionParams) (*VisionOCRResponse, error) {
	if err := validateVisionParams(params); err != nil {
		return nil, err
	}
	return doRequest[VisionParams, VisionOCRResponse](s.client, ctx, "POST", "/api/ai/vision/ocr", params)
}

// Classify returns structured tags, categories, and objects for the image.
func (s *VisionService) Classify(ctx context.Context, params *VisionParams) (*VisionClassifyResponse, error) {
	if err := validateVisionParams(params); err != nil {
		return nil, err
	}
	return doRequest[VisionParams, VisionClassifyResponse](s.client, ctx, "POST", "/api/ai/vision/classify", params)
}

func validateVisionParams(params *VisionParams) error {
	if params == nil {
		return errors.New("vision params required")
	}
	if params.ImageURL == "" && params.ImageBase64 == "" {
		return errors.New("vision request requires ImageURL or ImageBase64")
	}
	return nil
}
