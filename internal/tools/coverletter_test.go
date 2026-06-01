package tools

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/blackmagicbox/call-my-agent/internal/llm"
	"github.com/mark3labs/mcp-go/mcp"
)

func makeCoverLetterRequest(jobJSON, tone string) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	args := map[string]any{"job_json": jobJSON}
	if tone != "" {
		args["tone"] = tone
	}
	req.Params.Arguments = args
	return req
}

const validCoverLetter = `Dear Hiring Manager,

I am excited to apply for the SRE role at Acme Corp.`

func TestHandleCoverLetter_HappyPath(t *testing.T) {
	mock := &llm.MockProvider{Response: validCoverLetter}
	handler := HandleCoverLetter(mock, testProfile)

	result, err := handler(context.Background(), makeCoverLetterRequest(validJobJSON, "professional"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got tool error: %v", result.Content)
	}
	text := result.Content[0].(mcp.TextContent).Text
	if text != strings.TrimSpace(validCoverLetter) {
		t.Errorf("unexpected response: %s", text)
	}
}

// TestHandleCoverLetter_DefaultsToProfessionalTone verifies that when tone is
// omitted, "professional" is passed to the LLM prompts.
func TestHandleCoverLetter_DefaultsToProfessionalTone(t *testing.T) {
	mock := &llm.MockProvider{Response: validCoverLetter}
	handler := HandleCoverLetter(mock, testProfile)

	handler(context.Background(), makeCoverLetterRequest(validJobJSON, ""))

	if !strings.Contains(mock.LastSystem, "professional") {
		t.Errorf("expected default tone 'professional' in system prompt, got: %s", mock.LastSystem)
	}
}

// TestHandleCoverLetter_CustomTone verifies that a caller-supplied tone is
// forwarded to the LLM prompts.
func TestHandleCoverLetter_CustomTone(t *testing.T) {
	mock := &llm.MockProvider{Response: validCoverLetter}
	handler := HandleCoverLetter(mock, testProfile)

	handler(context.Background(), makeCoverLetterRequest(validJobJSON, "enthusiastic"))

	if !strings.Contains(mock.LastSystem, "enthusiastic") {
		t.Errorf("expected tone 'enthusiastic' in system prompt, got: %s", mock.LastSystem)
	}
}

func TestHandleCoverLetter_MissingArguments(t *testing.T) {
	mock := &llm.MockProvider{}
	handler := HandleCoverLetter(mock, testProfile)

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

func TestHandleCoverLetter_MissingJobJSON(t *testing.T) {
	mock := &llm.MockProvider{}
	handler := HandleCoverLetter(mock, testProfile)

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

// TestHandleCoverLetter_LLMError verifies a tool error is returned when the
// provider fails, and the Go error is nil (server stays alive).
func TestHandleCoverLetter_LLMError(t *testing.T) {
	mock := &llm.MockProvider{Err: errors.New("quota exceeded")}
	handler := HandleCoverLetter(mock, testProfile)

	result, err := handler(context.Background(), makeCoverLetterRequest(validJobJSON, ""))
	if err != nil {
		t.Fatalf("expected nil Go error, got: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected tool error when LLM fails")
	}
}

// TestHandleCoverLetter_ProfileAndJobInjectedIntoPrompt verifies that both
// the profile and job JSON are passed through to the LLM.
func TestHandleCoverLetter_ProfileAndJobInjectedIntoPrompt(t *testing.T) {
	mock := &llm.MockProvider{Response: validCoverLetter}
	handler := HandleCoverLetter(mock, testProfile)

	handler(context.Background(), makeCoverLetterRequest(validJobJSON, ""))

	if mock.LastUser == "" {
		t.Fatal("expected LLM to be called with a user prompt")
	}
	if mock.LastSystem == "" {
		t.Fatal("expected LLM to be called with a system prompt")
	}
	if !strings.Contains(mock.LastUser, validJobJSON) {
		t.Errorf("expected job JSON in user prompt, got: %s", mock.LastUser)
	}
}
