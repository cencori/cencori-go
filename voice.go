package cencori

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// VoiceService provides text-to-speech and speech-to-text across providers
// via the Cencori Voice API. The provider is inferred from the model.
type VoiceService struct {
	client *Client
}

// SpeakParams is the request payload for text-to-speech.
type SpeakParams struct {
	// Input is the text to synthesize (max 4096 characters).
	Input string `json:"input"`
	// Model id — provider is inferred from it. Default: tts-1 (OpenAI).
	Model string `json:"model,omitempty"`
	// Voice id/name. Provider-specific; falls back to the model default.
	Voice string `json:"voice,omitempty"`
	// Provider is optional — usually inferred from Model.
	Provider string `json:"provider,omitempty"`
	// ResponseFormat: mp3 (default), opus, aac, flac, wav, pcm.
	ResponseFormat string `json:"response_format,omitempty"`
	// Speed is playback speed 0.25–4.0 (OpenAI only).
	Speed float64 `json:"speed,omitempty"`
	// Language code (Spitch: en, yo, ha, ig, am).
	Language string `json:"language,omitempty"`
}

// SpeakResult holds synthesized audio and the resolved provider.
type SpeakResult struct {
	Audio       []byte
	ContentType string
	Provider    string
}

// TranscribeParams configures a speech-to-text request.
type TranscribeParams struct {
	// Model id — provider is inferred from it. Default: whisper-1 (OpenAI).
	Model string
	// Provider is optional — usually inferred from Model.
	Provider string
	// Language hint (Spitch uses it as the target language).
	Language string
	// Prompt is optional context (OpenAI only).
	Prompt string
	// Temperature 0–1.
	Temperature float64
	// Diarize requests speaker labels where the provider supports it.
	Diarize bool
	// ResponseFormat: json (default), text, srt, vtt, verbose_json.
	ResponseFormat string
	// Filename hint used for format detection. Default: audio.mp3.
	Filename string
}

// TranscriptSegment is one contiguous span of transcript, optionally attributed.
type TranscriptSegment struct {
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
	Text    string  `json:"text"`
	Speaker string  `json:"speaker,omitempty"`
}

// TranscriptWord is a single word with timing and optional speaker.
type TranscriptWord struct {
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
	Word    string  `json:"word"`
	Speaker string  `json:"speaker,omitempty"`
}

// TranscribeResult is the decoded transcription response.
type TranscribeResult struct {
	Text     string              `json:"text"`
	Language string              `json:"language,omitempty"`
	Duration float64             `json:"duration,omitempty"`
	Provider string              `json:"provider,omitempty"`
	Model    string              `json:"model,omitempty"`
	Segments []TranscriptSegment `json:"segments,omitempty"`
	Words    []TranscriptWord    `json:"words,omitempty"`
}

// VoiceModels lists available TTS and STT models.
type VoiceModels struct {
	TTS []map[string]any `json:"tts"`
	STT []map[string]any `json:"stt"`
}

// Speak synthesizes speech and returns the raw audio bytes plus the resolved
// provider. Provider is inferred from Model (default tts-1).
func (s *VoiceService) Speak(ctx context.Context, params *SpeakParams) (*SpeakResult, error) {
	if params == nil || params.Input == "" {
		return nil, errors.New("voice: Speak requires non-empty Input")
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", s.client.BaseURL+"/api/ai/audio/speech", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}
	audio, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("read audio: %w", err)
	}
	return &SpeakResult{
		Audio:       audio,
		ContentType: resp.Header.Get("Content-Type"),
		Provider:    resp.Header.Get("X-Provider"),
	}, nil
}

// Transcribe converts audio to text. Provider is inferred from Model
// (default whisper-1). audio is the raw file bytes.
func (s *VoiceService) Transcribe(ctx context.Context, audio []byte, params *TranscribeParams) (*TranscribeResult, error) {
	if len(audio) == 0 {
		return nil, errors.New("voice: Transcribe requires audio bytes")
	}
	if params == nil {
		params = &TranscribeParams{}
	}
	format := params.ResponseFormat
	if format == "" {
		format = "json"
	}
	filename := params.Filename
	if filename == "" {
		filename = "audio.mp3"
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(audio); err != nil {
		return nil, fmt.Errorf("write audio: %w", err)
	}
	fields := map[string]string{"response_format": format}
	if params.Model != "" {
		fields["model"] = params.Model
	}
	if params.Provider != "" {
		fields["provider"] = params.Provider
	}
	if params.Language != "" {
		fields["language"] = params.Language
	}
	if params.Prompt != "" {
		fields["prompt"] = params.Prompt
	}
	if params.Temperature != 0 {
		fields["temperature"] = strconv.FormatFloat(params.Temperature, 'f', -1, 64)
	}
	if params.Diarize {
		fields["diarize"] = "true"
	}
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("write field %s: %w", k, err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.client.BaseURL+"/api/ai/audio/transcriptions", &buf)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	// text/srt/vtt come back as plain text, not JSON.
	if format == "text" || format == "srt" || format == "vtt" {
		return &TranscribeResult{Text: string(body), Provider: resp.Header.Get("X-Provider")}, nil
	}
	var result TranscribeResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// Diarize transcribes with speaker labels. Use a diarization-capable model
// (nova-3, assemblyai-universal).
func (s *VoiceService) Diarize(ctx context.Context, audio []byte, params *TranscribeParams) (*TranscribeResult, error) {
	if params == nil {
		params = &TranscribeParams{}
	}
	params.Diarize = true
	params.ResponseFormat = "verbose_json"
	return s.Transcribe(ctx, audio, params)
}

// ListModels returns the available TTS and STT voice models.
func (s *VoiceService) ListModels(ctx context.Context) (*VoiceModels, error) {
	tts, err := s.getModels(ctx, "/api/ai/audio/speech")
	if err != nil {
		return nil, err
	}
	stt, err := s.getModels(ctx, "/api/ai/audio/transcriptions")
	if err != nil {
		return nil, err
	}
	return &VoiceModels{TTS: tts, STT: stt}, nil
}

func (s *VoiceService) getModels(ctx context.Context, path string) ([]map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.client.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}
	var payload struct {
		Models []map[string]any `json:"models"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Models, nil
}
