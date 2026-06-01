package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/blackmagicbox/call-my-agent/internal/llm"
	"github.com/blackmagicbox/call-my-agent/internal/profile"
	"github.com/blackmagicbox/call-my-agent/internal/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Resolve the profile.json path relative to the executable.
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to resolve executable path: %v", err)
	}
	profilePath := filepath.Join(filepath.Dir(exePath), "data", "profile.json")

	// Load and parse the candidate profile.
	data, err := os.ReadFile(profilePath)
	if err != nil {
		log.Fatalf("failed to load profile: %v", err)
	}
	var p profile.CandidateProfile
	if err := json.Unmarshal(data, &p); err != nil {
		log.Fatalf("failed to parse profile: %v", err)
	}

	// Initialize the MCP server.
	s := server.NewMCPServer(
		"call-my-agent",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Register tools.
	s.AddTool(
		mcp.NewTool(
			"get_candidate_profile",
			mcp.WithDescription(
				"Returns the candidate's profile information, including skills, experience, education, and preferences. This information can be used to match the candidate with suitable job opportunities."),
		),
		tools.HandleGetCandidateProfile(p),
	)
	provider, err := llm.FromConfig()
	if err != nil {
		log.Fatalf("failed to load llm: %v", err)
	}
	s.AddTool(tools.EvaluateJobTool(), tools.HandleEvaluateJob(provider, p))
	s.AddTool(tools.CoverLetterTool(), tools.HandleCoverLetter(provider, p))

	// Start server.
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
