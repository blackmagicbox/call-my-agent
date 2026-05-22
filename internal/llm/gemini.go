package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeminiProvider struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewGemini(apiKey, model string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (g *GeminiProvider) Complete(ctx context.Context, system, user string) (string, error) {
	// Assign the API endpoint
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", g.model)

	// Create the body with the instructions and the data for the Generation
	body, err := json.Marshal(map[string]interface{}{
		"system_instruction": map[string]interface{}{
			"parts": []map[string]string{{
				"text": system,
			}},
			"contents": map[string]interface{}{
				"parts": []map[string]string{{
					"text": user,
				}},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal content: %w", err)
	}

	// Create request to the API
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type and pass the API key
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)

	// Request Gemini api
	resp, err := g.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
