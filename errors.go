package cencori

import "fmt"

type APIError struct {
	StatusCode int
	Code       string         `json:"code"`
	Message    string         `json:"error"`
	Details    map[string]any `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}
