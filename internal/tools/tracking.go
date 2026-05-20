package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blackmagicbox/call-my-agent/internal/db"
	"github.com/blackmagicbox/call-my-agent/internal/job"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTracking(s *server.MCPServer, d *db.DB) {
	s.AddTool(saveJobTool(), saveJobHandler(d))
	s.AddTool(listJobsTool(), listJobsHandler(d))
	s.AddTool(applicationStatusTool(), applicationStatusHandler(d))
}

func saveJobTool() mcp.Tool {
	return mcp.NewTool("save_job",
		mcp.WithDescription("Persist a job listing and its evaluation result to the local tracker."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("Canonical URL of the job listing"),
		),
		mcp.WithString("title",
			mcp.Description("Job title"),
		),
		mcp.WithString("company",
			mcp.Description("Company name"),
		),
		mcp.WithString("description",
			mcp.Description("Raw job description text"),
		),
		mcp.WithNumber("fit_score",
			mcp.Description("Fitness score 0–100 from the evaluation"),
		),
		mcp.WithString("fit_reason",
			mcp.Description("One-paragraph explanation of the fitness score"),
		),
		mcp.WithString("status",
			mcp.Description("Tracking status: to_apply | applied | rejected | archived (default: to_apply)"),
		),
	)
}

func saveJobHandler(d *db.DB) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := req.GetString("url", "")
		if url == "" {
			return mcp.NewToolResultError("url is required"), nil
		}

		status := job.Status(req.GetString("status", string(job.StatusToApply)))
		switch status {
		case job.StatusToApply, job.StatusApplied, job.StatusRejected, job.StatusArchived:
		default:
			return mcp.NewToolResultError(fmt.Sprintf("invalid status %q: must be one of to_apply, applied, rejected, archived", status)), nil
		}

		j := &job.Job{
			URL:         url,
			Title:       req.GetString("title", ""),
			Company:     req.GetString("company", ""),
			Description: req.GetString("description", ""),
			FitScore:    req.GetInt("fit_score", 0),
			FitReason:   req.GetString("fit_reason", ""),
			Status:      status,
		}

		if err := d.SaveJob(j); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("save job: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Job saved (id=%d, status=%s)", j.ID, j.Status)), nil
	}
}

func listJobsTool() mcp.Tool {
	return mcp.NewTool("list_jobs",
		mcp.WithDescription("Return all tracked jobs, optionally filtered by status."),
		mcp.WithString("status",
			mcp.Description("Filter by status: to_apply | applied | rejected | archived. Omit to return all."),
		),
	)
}

func listJobsHandler(d *db.DB) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		status := req.GetString("status", "")

		jobs, err := d.ListJobs(status)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("list jobs: %v", err)), nil
		}

		if len(jobs) == 0 {
			return mcp.NewToolResultText("No jobs tracked yet."), nil
		}

		type row struct {
			ID       int64  `json:"id"`
			URL      string `json:"url"`
			Title    string `json:"title"`
			Company  string `json:"company"`
			FitScore int    `json:"fit_score"`
			Status   string `json:"status"`
		}
		rows := make([]row, len(jobs))
		for i, j := range jobs {
			rows[i] = row{
				ID:       j.ID,
				URL:      j.URL,
				Title:    j.Title,
				Company:  j.Company,
				FitScore: j.FitScore,
				Status:   string(j.Status),
			}
		}

		out, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("marshal: %v", err)), nil
		}
		return mcp.NewToolResultText(string(out)), nil
	}
}

func applicationStatusTool() mcp.Tool {
	return mcp.NewTool("get_application_status",
		mcp.WithDescription("Return a summary of tracked jobs broken down by status."),
	)
}

func applicationStatusHandler(d *db.DB) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stats, err := d.ApplicationStats()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("stats: %v", err)), nil
		}

		type result struct {
			Total    int `json:"total"`
			ToApply  int `json:"to_apply"`
			Applied  int `json:"applied"`
			Rejected int `json:"rejected"`
			Archived int `json:"archived"`
		}
		out, _ := json.MarshalIndent(result{
			Total:    stats.Total,
			ToApply:  stats.ToApply,
			Applied:  stats.Applied,
			Rejected: stats.Rejected,
			Archived: stats.Archived,
		}, "", "  ")
		return mcp.NewToolResultText(string(out)), nil
	}
}
