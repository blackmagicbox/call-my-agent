package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/blackmagicbox/call-my-agent/internal/llm"
	"github.com/mark3labs/mcp-go/mcp"
)

func makeEvaluateRequest(jobJSON string) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{
		"job_json": jobJSON,
	}
	return req
}

const validJobJSON = `{"title":"SRE","company":"Acme","location":"Berlin","description":"Run things"}`

const validLLMResponse = `{"fit_score":85,"seniority_match":"yes","matched_skills":["Go"],"skill_gaps":[],"red_flags_hit":[],"recommendation":"apply","reasoning":"Good fit."}`

// TestHandleEvaluateJob_HappyPath verifies the handler returns valid JSON on success.
func TestHandleEvaluateJob_HappyPath(t *testing.T) {
	mock := &llm.MockProvider{Response: validLLMResponse}
	handler := HandleEvaluateJob(mock, testProfile)

	result, err := handler(context.Background(), makeEvaluateRequest(validJobJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got tool error: %v", result.Content)
	}
	text := result.Content[0].(mcp.TextContent).Text
	if text != validLLMResponse {
		t.Errorf("unexpected response: %s", text)
	}
}

// TestHandleEvaluateJob_StripsJSONFences verifies markdown fences are removed.
func TestHandleEvaluateJob_StripsJSONFences(t *testing.T) {
	fenced := "```json\n" + validLLMResponse + "\n```"
	mock := &llm.MockProvider{Response: fenced}
	handler := HandleEvaluateJob(mock, testProfile)

	result, err := handler(context.Background(), makeEvaluateRequest(validJobJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := result.Content[0].(mcp.TextContent).Text
	if text != validLLMResponse {
		t.Errorf("fences not stripped, got: %s", text)
	}
}

// TestHandleEvaluateJob_MissingArguments verifies a tool error is returned when
// arguments are nil/wrong type.
func TestHandleEvaluateJob_MissingArguments(t *testing.T) {
	mock := &llm.MockProvider{}
	handler := HandleEvaluateJob(mock, testProfile)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = nil

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected tool error for nil arguments")
	}
}

// TestHandleEvaluateJob_MissingJobJSON verifies a tool error is returned when
// job_json key is absent.
func TestHandleEvaluateJob_MissingJobJSON(t *testing.T) {
	mock := &llm.MockProvider{}
	handler := HandleEvaluateJob(mock, testProfile)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected tool error for missing job_json")
	}
}

// TestHandleEvaluateJob_LLMError verifies a tool error is returned when the
// provider fails, and the Go error is nil (server stays alive).
func TestHandleEvaluateJob_LLMError(t *testing.T) {
	mock := &llm.MockProvider{Err: errors.New("rate limited")}
	handler := HandleEvaluateJob(mock, testProfile)

	result, err := handler(context.Background(), makeEvaluateRequest(validJobJSON))
	if err != nil {
		t.Fatalf("expected nil Go error, got: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected tool error when LLM fails")
	}
}

// TestHandleEvaluateJob_ProfileInjectedIntoPrompt verifies the profile JSON
// is passed through to the LLM prompt.
func TestHandleEvaluateJob_ProfileInjectedIntoPrompt(t *testing.T) {
	mock := &llm.MockProvider{Response: validLLMResponse}
	handler := HandleEvaluateJob(mock, testProfile)

	handler(context.Background(), makeEvaluateRequest(validJobJSON))

	if mock.LastUser == "" {
		t.Fatal("expected LLM to be called with a user prompt")
	}
	if mock.LastSystem == "" {
		t.Fatal("expected LLM to be called with a system prompt")
	}
}
