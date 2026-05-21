package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/blackmagicbox/call-my-agent/internal/profile"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 1. Read the local profile.json
	data, err := os.ReadFile("data/profile.json")
	if err != nil {
		log.Fatalf("failed to load profile: %v", err)
	}

	// 2. Create the profile variable
	var p profile.CandidateProfile

	// 3. Parse the contents of the profile.json to `p`
	if err := json.Unmarshal(data, &p); err != nil {
		log.Fatalf("failed to parse profile: %v", err)
	}

	// 4. Initialize the MCP server
	s := server.NewMCPServer(
		"call-my-agent",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// 5. Create the Tool
	tool := mcp.NewTool(
		"get_candidate_profile",
		mcp.WithDescription(
			"Returns the candidate's profile information, including skills, experience, education, and preferences. This information can be used to match the candidate with suitable job opportunities."),
	)

	// 6. Register the Tool with the server
	s.AddTool(
		tool,
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := json.Marshal(p)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(out)), nil
		})

	// 7. Start server
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
