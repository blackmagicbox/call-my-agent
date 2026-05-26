package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ---------------------------------------------------------------------------
// GeminiProvider.Complete
// ---------------------------------------------------------------------------

func newTestGemini(url string) *GeminiProvider {
	return &GeminiProvider{
		apiKey:  "test-key",
		model:   "gemini-3.5-flash",
		http:    &http.Client{},
		baseURL: url,
	}
}

func TestGeminiComplete_HappyPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json content-type, got %s", ct)
		}
		if key := r.Header.Get("X-LLM-API-KEY"); key != "test-key" {
			t.Errorf("expected api key header test-key, got %s", key)
		}

		resp := map[string]any{
			"candidates": []map[string]any{
				{"content": map[string]any{
					"parts": []map[string]string{{"text": "hello world"}},
				}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	g := newTestGemini(srv.URL)
	got, err := g.Complete(context.Background(), "system prompt", "user prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestGeminiComplete_RequestBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		// Verify system_instruction is present.
		si, ok := body["system_instruction"].(map[string]any)
		if !ok {
			t.Fatal("missing system_instruction in request body")
		}
		parts := si["parts"].([]any)
		if parts[0].(map[string]any)["text"] != "be helpful" {
			t.Error("system_instruction text mismatch")
		}

		// Verify contents.
		contents := body["contents"].([]any)
		cParts := contents[0].(map[string]any)["parts"].([]any)
		if cParts[0].(map[string]any)["text"] != "say hi" {
			t.Error("contents text mismatch")
		}

		resp := map[string]any{
			"candidates": []map[string]any{
				{"content": map[string]any{
					"parts": []map[string]string{{"text": "hi"}},
				}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	g := newTestGemini(srv.URL)
	_, err := g.Complete(context.Background(), "be helpful", "say hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGeminiComplete_NoCandidates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"candidates": []any{}})
	}))
	defer srv.Close()

	g := newTestGemini(srv.URL)
	_, err := g.Complete(context.Background(), "sys", "usr")
	if err == nil {
		t.Fatal("expected error for empty candidates, got nil")
	}
}

func TestGeminiComplete_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	g := newTestGemini(srv.URL)
	_, err := g.Complete(context.Background(), "sys", "usr")
	if err == nil {
		t.Fatal("expected error for invalid JSON response, got nil")
	}
}

func TestGeminiComplete_CancelledContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server")
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	g := newTestGemini(srv.URL)
	_, err := g.Complete(ctx, "sys", "usr")
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

// ---------------------------------------------------------------------------
// GeminiProvider satisfies the Provider interface
// ---------------------------------------------------------------------------

func TestGeminiProvider_ImplementsProvider(t *testing.T) {
	var _ Provider = (*GeminiProvider)(nil)
}

// ---------------------------------------------------------------------------
// FromConfig
// ---------------------------------------------------------------------------

func TestFromConfig_MissingProvider(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "")
	t.Setenv("LLM_API_KEY", "key")

	_, err := FromConfig()
	if err == nil {
		t.Fatal("expected error when LLM_PROVIDER is empty")
	}
}

func TestFromConfig_UnknownProvider(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "openai")
	t.Setenv("LLM_API_KEY", "key")

	_, err := FromConfig()
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestFromConfig_GeminiMissingKey(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "gemini")
	t.Setenv("LLM_API_KEY", "")

	_, err := FromConfig()
	if err == nil {
		t.Fatal("expected error when api key is missing for gemini")
	}
}

func TestFromConfig_GeminiHappyPath(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "gemini")
	t.Setenv("LLM_API_KEY", "test-key")
	t.Setenv("LLM_MODEL", "gemini-3.5-flash")

	p, err := FromConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}

	gp, ok := p.(*GeminiProvider)
	if !ok {
		t.Fatal("expected *GeminiProvider")
	}
	if gp.apiKey != "test-key" {
		t.Errorf("expected apiKey %q, got %q", "test-key", gp.apiKey)
	}
	if gp.model != "gemini-3.5-flash" {
		t.Errorf("expected model %q, got %q", "gemini-3.5-flash", gp.model)
	}
}
