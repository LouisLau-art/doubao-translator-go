package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DoubaoClient handles communication with Doubao API
type DoubaoClient struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

// NewDoubaoClient creates a new Doubao API client
func NewDoubaoClient(apiKey, apiURL string) *DoubaoClient {
	return &DoubaoClient{
		apiKey: apiKey,
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request structures
type TranslateRequest struct {
	Text   string `json:"text" binding:"required"`
	Source string `json:"source"`
	Target string `json:"target" binding:"required"`
}

type DoubaoRequest struct {
	Model string               `json:"model"`
	Input []DoubaoInputMessage `json:"input"`
}

type DoubaoInputMessage struct {
	Role    string          `json:"role"`
	Content []DoubaoContent `json:"content"`
}

type DoubaoContent struct {
	Type               string           `json:"type"`
	Text               string           `json:"text"`
	TranslationOptions *TranslationOpts `json:"translation_options,omitempty"`
}

type TranslationOpts struct {
	SourceLanguage string `json:"source_language,omitempty"`
	TargetLanguage string `json:"target_language"`
}

type DoubaoResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type DoubaoNewResponse struct {
	Output []struct {
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Status string `json:"status"`
}

// Translate sends a translation request to Doubao API
func (c *DoubaoClient) Translate(text, source, target string) (string, error) {
	opts := &TranslationOpts{
		TargetLanguage: target,
	}
	if source != "" {
		opts.SourceLanguage = source
	}

	reqBody := DoubaoRequest{
		Model: "doubao-seed-translation-250915",
		Input: []DoubaoInputMessage{
			{
				Role: "user",
				Content: []DoubaoContent{
					{
						Type:               "input_text",
						Text:               text,
						TranslationOptions: opts,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	// Try new format first
	var newResult DoubaoNewResponse
	if err := json.Unmarshal(body, &newResult); err == nil && newResult.Status == "completed" {
		for _, output := range newResult.Output {
			if output.Type == "message" && output.Role == "assistant" {
				for _, content := range output.Content {
					if content.Type == "output_text" {
						return content.Text, nil
					}
				}
			}
		}
		return "", fmt.Errorf("no output_text in new format")
	}

	// Fallback to old format
	var oldResult DoubaoResponse
	if err := json.Unmarshal(body, &oldResult); err == nil {
		if len(oldResult.Choices) > 0 {
			return oldResult.Choices[0].Message.Content, nil
		}
	}

	return "", fmt.Errorf("unable to parse API response")
}