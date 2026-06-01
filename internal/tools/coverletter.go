package tools

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/blackmagicbox/call-my-agent/internal/llm"
	"github.com/blackmagicbox/call-my-agent/internal/profile"
	"github.com/blackmagicbox/call-my-agent/internal/prompts"
	"github.com/mark3labs/mcp-go/mcp"
)

func CoverLetterTool() mcp.Tool {
	return mcp.NewTool(
		"generate_cover_letter",
		mcp.WithDescription("Generates a tailored cover letter for a job listing based on the candidate profile"),
		mcp.WithString("job_json", mcp.Required(), mcp.Description("JSON-encoded job listing")),
		mcp.WithString("tone", mcp.Description("Tone of the cover letter (e.g. professional, friendly, enthusiastic). Defaults to professional.")),
	)
}

func HandleCoverLetter(provider llm.Provider, p profile.CandidateProfile) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("missing arguments"), nil
		}
		jobJSON, ok := args["job_json"].(string)
		if !ok {
			return mcp.NewToolResultError("missing job_json"), nil
		}
		tone := "professional"
		if t, ok := args["tone"].(string); ok && strings.TrimSpace(t) != "" {
			tone = t
		}

		profileBytes, err := json.Marshal(p)
		if err != nil {
			return mcp.NewToolResultError("failed to serialize profile"), nil
		}

		data := prompts.CoverLetterData{
			ProfileJSON: string(profileBytes),
			JobJSON:     jobJSON,
			Tone:        tone,
		}
		systemPrompt, err := prompts.CoverLetterSystem(data)
		if err != nil {
			return mcp.NewToolResultError("failed to render system prompt"), nil
		}
		userPrompt, err := prompts.CoverLetterUser(data)
		if err != nil {
			return mcp.NewToolResultError("failed to render user prompt"), nil
		}

		raw, err := provider.Complete(ctx, systemPrompt, userPrompt)
		if err != nil {
			return mcp.NewToolResultError("LLM call failed"), nil
		}
		return mcp.NewToolResultText(strings.TrimSpace(raw)), nil
	}
}
