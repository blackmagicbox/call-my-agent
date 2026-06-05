package tools

import (
	"context"
	"encoding/json"
	"time"

	"github.com/blackmagicbox/call-my-agent/internal/db"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
)

func SaveJobTool() mcp.Tool {
	return mcp.NewTool(
		"save_job",
		mcp.WithDescription("Saves a job listing and its evaluation result to the local tracking database"),
		mcp.WithString("job_json", mcp.Required(), mcp.Description("JSON-encoded job listing")),
		mcp.WithString("evaluation_json", mcp.Required(), mcp.Description("JSON-encoded evaluation result from evaluate_job")),
		mcp.WithString("status", mcp.Description("Job status: to_evaluate, to_apply, applied, rejected, archived. Defaults to to_apply.")),
	)
}

func HandleSaveJob(database *db.DB) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("missing arguments"), nil
		}
		jobJSON, ok := args["job_json"].(string)
		if !ok {
			return mcp.NewToolResultError("missing job_json"), nil
		}
		evalJSON, ok := args["evaluation_json"].(string)
		if !ok {
			return mcp.NewToolResultError("missing evaluation_json"), nil
		}
		status := "to_apply"
		if s, ok := args["status"].(string); ok && s != "" {
			status = s
		}

		var job map[string]any
		if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
			return mcp.NewToolResultError("invalid job_json"), nil
		}
		var eval map[string]any
		if err := json.Unmarshal([]byte(evalJSON), &eval); err != nil {
			return mcp.NewToolResultError("invalid evaluation_json"), nil
		}

		id, _ := job["id"].(string)
		if id == "" {
			id = uuid.New().String()
		}

		matchedSkills, _ := json.Marshal(eval["matched_skills"])
		skillGaps, _ := json.Marshal(eval["skill_gaps"])
		redFlags, _ := json.Marshal(eval["red_flags_hit"])
		fitScore, _ := eval["fit_score"].(float64)

		err := database.SaveJob(db.SaveJobParams{
			ID:             id,
			Title:          str(job["title"]),
			Company:        str(job["company"]),
			Location:       str(job["location"]),
			Remote:         str(job["remote"]),
			URL:            str(job["url"]),
			Status:         status,
			FitScore:       int(fitScore),
			Recommendation: str(eval["recommendation"]),
			MatchedSkills:  string(matchedSkills),
			SkillGaps:      string(skillGaps),
			RedFlagsHit:    string(redFlags),
			Reasoning:      str(eval["reasoning"]),
			JobJSON:        jobJSON,
			SeenAt:         time.Now().UTC().Format(time.RFC3339),
		})
		if err != nil {
			return mcp.NewToolResultError("failed to save job: " + err.Error()), nil
		}

		out, err := json.Marshal(map[string]any{"saved": true, "id": id})
		if err != nil {
			return mcp.NewToolResultError("failed to serialise response"), nil
		}
		return mcp.NewToolResultText(string(out)), nil
	}
}

func ListJobsTool() mcp.Tool {
	return mcp.NewTool(
		"list_jobs",
		mcp.WithDescription("Lists tracked jobs, optionally filtered by status"),
		mcp.WithString("status", mcp.Description("Filter by status: to_evaluate, to_apply, applied, rejected, archived. Omit for all.")),
	)
}

func HandleListJobs(database *db.DB) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		status := ""
		if args, ok := req.Params.Arguments.(map[string]any); ok {
			status, _ = args["status"].(string)
		}

		jobs, err := database.ListJobs(status)
		if err != nil {
			return mcp.NewToolResultError("failed to list jobs: " + err.Error()), nil
		}

		out, err := json.Marshal(jobs)
		if err != nil {
			return mcp.NewToolResultError("failed to serialise jobs"), nil
		}
		return mcp.NewToolResultText(string(out)), nil
	}
}

func GetApplicationStatusTool() mcp.Tool {
	return mcp.NewTool(
		"get_application_status",
		mcp.WithDescription("Returns a summary of your job application pipeline"),
	)
}

func HandleGetApplicationStatus(database *db.DB) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		summary, err := database.GetStatusSummary()
		if err != nil {
			return mcp.NewToolResultError("failed to get summary: " + err.Error()), nil
		}

		out, err := json.Marshal(summary)
		if err != nil {
			return mcp.NewToolResultError("failed to serialise summary"), nil
		}
		return mcp.NewToolResultText(string(out)), nil
	}
}

// str safely casts an any to string.
func str(v any) string {
	s, _ := v.(string)
	return s
}
