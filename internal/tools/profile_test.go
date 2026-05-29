package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blackmagicbox/call-my-agent/internal/profile"
	"github.com/mark3labs/mcp-go/mcp"
)

// testProfile is a minimal but complete CandidateProfile used across tests.
var testProfile = profile.CandidateProfile{
	Name:     "Alex Johnson",
	Title:    "Senior Software Engineer",
	Location: "Amsterdam, Netherlands",
	Summary:  "8+ years building distributed systems in Go.",
	Skills:   []string{"Go", "Python", "Docker", "Kubernetes"},
	Preferences: profile.Preferences{
		TargetRoles:      []string{"Senior Software Engineer", "Staff Engineer"},
		TargetLevels:     []string{"Senior", "Staff"},
		Locations:        []string{"Amsterdam", "Remote"},
		RemotePreference: "preferred",
	},
}

// TestHandleGetCandidateProfile_HappyPath verifies that the handler returns a
// non-nil, non-error result containing valid JSON.
func TestHandleGetCandidateProfile_HappyPath(t *testing.T) {
	handler := HandleGetCandidateProfile(testProfile)

	result, err := handler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.IsError {
		t.Fatal("expected IsError=false")
	}
	if len(result.Content) == 0 {
		t.Fatal("expected at least one content item")
	}
}

// TestHandleGetCandidateProfile_ValidJSON verifies that the content text is
// valid, parseable JSON.
func TestHandleGetCandidateProfile_ValidJSON(t *testing.T) {
	handler := HandleGetCandidateProfile(testProfile)

	result, err := handler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var parsed map[string]any
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		t.Fatalf("result content is not valid JSON: %v\ncontent: %s", err, text)
	}
}

// TestHandleGetCandidateProfile_RequiredFields verifies that all required top-level
// fields are present in the JSON output.
func TestHandleGetCandidateProfile_RequiredFields(t *testing.T) {
	handler := HandleGetCandidateProfile(testProfile)

	result, err := handler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var parsed map[string]any
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		t.Fatalf("result content is not valid JSON: %v", err)
	}

	required := []string{"name", "title", "location", "skills", "preferences"}
	for _, field := range required {
		if _, ok := parsed[field]; !ok {
			t.Errorf("missing required field: %q", field)
		}
	}
}

// TestHandleGetCandidateProfile_SkillsNonEmpty verifies that the skills field
// is a non-empty list.
func TestHandleGetCandidateProfile_SkillsNonEmpty(t *testing.T) {
	handler := HandleGetCandidateProfile(testProfile)

	result, err := handler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var parsed map[string]any
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		t.Fatalf("result content is not valid JSON: %v", err)
	}

	skills, ok := parsed["skills"].([]any)
	if !ok {
		t.Fatal("expected skills to be a JSON array")
	}
	if len(skills) == 0 {
		t.Error("expected skills to be non-empty")
	}
}

// TestHandleGetCandidateProfile_PreferencesFields verifies that the preferences
// object contains target_roles and target_levels.
func TestHandleGetCandidateProfile_PreferencesFields(t *testing.T) {
	handler := HandleGetCandidateProfile(testProfile)

	result, err := handler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var parsed map[string]any
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		t.Fatalf("result content is not valid JSON: %v", err)
	}

	prefs, ok := parsed["preferences"].(map[string]any)
	if !ok {
		t.Fatal("expected preferences to be a JSON object")
	}
	for _, key := range []string{"target_roles", "target_levels"} {
		if _, ok := prefs[key]; !ok {
			t.Errorf("missing preferences field: %q", key)
		}
	}
}
