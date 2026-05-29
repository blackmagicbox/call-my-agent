package tools

import (
	"context"
	"encoding/json"

	"github.com/blackmagicbox/call-my-agent/internal/llm"
	"github.com/blackmagicbox/call-my-agent/internal/profile"
	"github.com/blackmagicbox/call-my-agent/internal/prompts"
	"github.com/mark3labs/mcp-go/mcp"
)

func EvaluateJobTool() mcp.Tool {
	return mcp.NewTool(
		"evaluate_job",
		mcp.WithDescription("Evaluates a job listing against the candidate profile and returns a fit score, skill gaps, and recommendation"),
		mcp.WithString("job_json", mcp.Required(), mcp.Description("JSON-encoded job listing")),
	)
}

func HandleEvaluateJob(provider llm.Provider, p profile.CandidateProfile) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the job listing JSON from the tool arguments
		args, ok := req.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Missing arguments"), nil
		}
		jobJSON, ok := args["job_json"].(string)
		if !ok {
			return mcp.NewToolResultError("Missing job_json"), nil
		}

		// Serialize profile to Json
		profileBytes, err := json.Marshal(p)
		if err != nil {
			return mcp.NewToolResultError("Failed to serialize profile"), nil
		}

		// Render the user prompt from template
		userPrompt, err := prompts.EvaluateJobUser(prompts.EvaluateJobData{
			ProfileJSON: string(profileBytes),
			JobJSON:     jobJSON,
		})
		if err != nil {
			return mcp.NewToolResultError("failed to render prompt"), nil
		}
		// Call the llm
		raw, err := provider.Complete(ctx, prompts.EvaluateJobSystem(), userPrompt)
		if err != nil {
			return mcp.NewToolResultError("LLM call failed"), nil
		}
		// Strip Markdown fences and return
		return mcp.NewToolResultText(StripJSONFences(raw)), nil
	}
}
