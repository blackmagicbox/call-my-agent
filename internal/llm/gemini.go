package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GeminiProvider struct {
	apiKey  string
	model   string
	http    *http.Client
	baseURL string // overridden in tests; defaults to Gemini production endpoint
}

func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}
func (g *GeminiProvider) Complete(ctx context.Context, system string, user string) (string, error) {
	url := g.baseURL
	if url == "" {
		url = "https://generativelanguage.googleapis.com/v1beta/models/" + g.model + ":generateContent"
	}

	body, err := json.Marshal(map[string]any{
		"system_instruction": map[string]any{
			"parts": []map[string]string{{"text": system}},
		},
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": user}}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("content-type", "application/json")
	request.Header.Add("X-LLM-API-KEY", g.apiKey)

	response, err := g.http.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to complete request: %w", err)
	}
	defer response.Body.Close()

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("failed to complete Gemini model")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
