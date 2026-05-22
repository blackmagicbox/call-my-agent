package tools

import (
	"context"
	"encoding/json"

	"github.com/blackmagicbox/call-my-agent/internal/profile"
	"github.com/mark3labs/mcp-go/mcp"
)

// HandleGetCandidateProfile returns an MCP tool handler that serialises the
// given CandidateProfile as JSON text. Keeping the profile as a parameter
// (rather than a global) makes the handler easy to test without touching the
// filesystem.
func HandleGetCandidateProfile(p profile.CandidateProfile) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		out, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(string(out)), nil
	}
}
