package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/blackmagicbox/call-my-agent/internal/db"
	"github.com/mark3labs/mcp-go/mcp"
)

func openTrackingDB(t *testing.T) *db.DB {
	t.Helper()
	d, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

const validTrackingJobJSON = `{"id":"job-abc","title":"SRE L5","company":"Acme","location":"Berlin","remote":"hybrid","url":"https://example.com"}`
const validEvalJSON = `{"fit_score":85,"recommendation":"apply","matched_skills":["Go"],"skill_gaps":[],"red_flags_hit":[],"reasoning":"strong match"}`

func makeSaveJobRequest(jobJSON, evalJSON, status string) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	args := map[string]any{
		"job_json":        jobJSON,
		"evaluation_json": evalJSON,
	}
	if status != "" {
		args["status"] = status
	}
	req.Params.Arguments = args
	return req
}

// TestHandleSaveJob_HappyPath verifies a valid save returns saved=true and the job ID.
func TestHandleSaveJob_HappyPath(t *testing.T) {
	d := openTrackingDB(t)
	handler := HandleSaveJob(d)

	result, err := handler(context.Background(), makeSaveJobRequest(validTrackingJobJSON, validEvalJSON, "to_apply"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got tool error: %v", result.Content)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var out map[string]any
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if out["saved"] != true {
		t.Errorf("want saved=true, got %v", out["saved"])
	}
	if out["id"] != "job-abc" {
		t.Errorf("want id=job-abc, got %v", out["id"])
	}
}

// TestHandleSaveJob_GeneratesIDWhenMissing verifies a UUID is assigned when job has no ID.
func TestHandleSaveJob_GeneratesIDWhenMissing(t *testing.T) {
	d := openTrackingDB(t)
	handler := HandleSaveJob(d)

	noIDJob := `{"title":"SRE L5","company":"Acme"}`
	result, err := handler(context.Background(), makeSaveJobRequest(noIDJob, validEvalJSON, ""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got tool error: %v", result.Content)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var out map[string]any
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	id, _ := out["id"].(string)
	if id == "" {
		t.Error("want a generated ID, got empty string")
	}
}

// TestHandleSaveJob_MissingArguments verifies a tool error is returned for nil args.
func TestHandleSaveJob_MissingArguments(t *testing.T) {
	d := openTrackingDB(t)
	handler := HandleSaveJob(d)

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

// TestHandleSaveJob_InvalidJobJSON verifies a tool error is returned for malformed JSON.
func TestHandleSaveJob_InvalidJobJSON(t *testing.T) {
	d := openTrackingDB(t)
	handler := HandleSaveJob(d)

	result, err := handler(context.Background(), makeSaveJobRequest("not-json", validEvalJSON, ""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected tool error for invalid job_json")
	}
}

// TestHandleListJobs_ReturnsJSON verifies list returns a JSON array after a save.
func TestHandleListJobs_ReturnsJSON(t *testing.T) {
	d := openTrackingDB(t)

	_, _ = HandleSaveJob(d)(context.Background(), makeSaveJobRequest(validTrackingJobJSON, validEvalJSON, "to_apply"))

	handler := HandleListJobs(d)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"status": "to_apply"}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got tool error: %v", result.Content)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var jobs []map[string]any
	if err := json.Unmarshal([]byte(text), &jobs); err != nil {
		t.Fatalf("response is not a JSON array: %v\ncontent: %s", err, text)
	}
	if len(jobs) != 1 {
		t.Fatalf("want 1 job, got %d", len(jobs))
	}
}

// TestHandleListJobs_FilterByStatus verifies only matching status jobs are returned.
func TestHandleListJobs_FilterByStatus(t *testing.T) {
	d := openTrackingDB(t)

	job2 := `{"id":"job-2","title":"SRE","company":"Beta"}`
	_, _ = HandleSaveJob(d)(context.Background(), makeSaveJobRequest(validTrackingJobJSON, validEvalJSON, "to_apply"))
	_, _ = HandleSaveJob(d)(context.Background(), makeSaveJobRequest(job2, validEvalJSON, "rejected"))

	handler := HandleListJobs(d)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"status": "to_apply"}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var jobs []map[string]any
	if err := json.Unmarshal([]byte(text), &jobs); err != nil {
		t.Fatalf("response is not a JSON array: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("want 1 job with status to_apply, got %d", len(jobs))
	}
}

// TestHandleGetApplicationStatus_ReturnsCounts verifies summary counts after inserts.
func TestHandleGetApplicationStatus_ReturnsCounts(t *testing.T) {
	d := openTrackingDB(t)

	job2 := `{"id":"job-2","title":"SRE","company":"Beta"}`
	job3 := `{"id":"job-3","title":"SRE","company":"Gamma"}`
	_, _ = HandleSaveJob(d)(context.Background(), makeSaveJobRequest(validTrackingJobJSON, validEvalJSON, "to_apply"))
	_, _ = HandleSaveJob(d)(context.Background(), makeSaveJobRequest(job2, validEvalJSON, "to_apply"))
	_, _ = HandleSaveJob(d)(context.Background(), makeSaveJobRequest(job3, validEvalJSON, "applied"))

	handler := HandleGetApplicationStatus(d)
	result, err := handler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got tool error: %v", result.Content)
	}

	text := result.Content[0].(mcp.TextContent).Text
	var summary map[string]any
	if err := json.Unmarshal([]byte(text), &summary); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if summary["to_apply"].(float64) != 2 {
		t.Errorf("want to_apply=2, got %v", summary["to_apply"])
	}
	if summary["applied"].(float64) != 1 {
		t.Errorf("want applied=1, got %v", summary["applied"])
	}
	if summary["total"].(float64) != 3 {
		t.Errorf("want total=3, got %v", summary["total"])
	}
}
