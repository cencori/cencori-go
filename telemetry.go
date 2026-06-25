package cencori

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// TelemetryService provides methods for reporting web traffic to the Cencori dashboard.
type TelemetryService struct {
	client *Client
}

// ReportWebRequest reports a web request to the Cencori dashboard.
// This is fire-and-forget — it never blocks or returns application-level errors.
func (s *TelemetryService) ReportWebRequest(ctx context.Context, payload *WebTelemetryPayload) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		s.client.BaseURL+"/api/v1/telemetry/web",
		bytes.NewReader(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}

// ReportWebRequestAsync is an async variant that launches the request in a goroutine.
func (s *TelemetryService) ReportWebRequestAsync(payload *WebTelemetryPayload) {
	go s.ReportWebRequest(context.Background(), payload)
}
