package cencori

import (
	"errors"
	"net/http"
	"time"
)

type ClientOptions struct {
	ApiKey  string
	BaseURL string
	Timeout time.Duration
}

func WithApiKey(apiKey string) Option {
	return func(c *ClientOptions) { c.ApiKey = apiKey }
}

func WithBaseURL(baseURL string) Option {
	return func(c *ClientOptions) { c.BaseURL = baseURL }
}

type Client struct {
	ApiKey     string
	BaseURL    string
	httpClient *http.Client

	Chat *ChatService
}

type Option func(*ClientOptions)

func NewClient(opts ...Option) (*Client, error) {
	config := &ClientOptions{
		BaseURL: "https://cencori.com",
		Timeout: 30,
	}
	for _, opt := range opts {
		opt(config)
	}

	if config.ApiKey == "" {
		return nil, errors.New("You need a valid API Key to use this client")
	}

	c := &Client{
		ApiKey:  config.ApiKey,
		BaseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}

	c.Chat = &ChatService{client: c}

	return c, nil
}
