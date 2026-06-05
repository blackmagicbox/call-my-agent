package db_test

import (
	"testing"

	"github.com/blackmagicbox/call-my-agent/internal/db"
)

func openTestDB(t *testing.T) *db.DB {
	t.Helper()
	d, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

func testParams(id, status string) db.SaveJobParams {
	return db.SaveJobParams{
		ID:             id,
		Title:          "SRE L5",
		Company:        "Acme",
		Location:       "Berlin",
		Remote:         "hybrid",
		URL:            "https://example.com/job/" + id,
		Status:         status,
		FitScore:       85,
		Recommendation: "apply",
		MatchedSkills:  `["Go","Kubernetes"]`,
		SkillGaps:      `[]`,
		RedFlagsHit:    `[]`,
		Reasoning:      "strong match",
		JobJSON:        `{"id":"` + id + `"}`,
		SeenAt:         "2026-01-01T00:00:00Z",
	}
}

func TestSaveJob_HappyPath(t *testing.T) {
	d := openTestDB(t)

	if err := d.SaveJob(testParams("job-1", "to_apply")); err != nil {
		t.Fatalf("SaveJob: %v", err)
	}

	jobs, err := d.ListJobs("to_apply")
	if err != nil {
		t.Fatalf("ListJobs: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("want 1 job, got %d", len(jobs))
	}
	if jobs[0].ID != "job-1" {
		t.Errorf("want id job-1, got %s", jobs[0].ID)
	}
	if jobs[0].FitScore != 85 {
		t.Errorf("want fit_score 85, got %d", jobs[0].FitScore)
	}
}

func TestSaveJob_Upsert(t *testing.T) {
	d := openTestDB(t)

	if err := d.SaveJob(testParams("job-1", "to_apply")); err != nil {
		t.Fatalf("first SaveJob: %v", err)
	}

	updated := testParams("job-1", "applied")
	updated.FitScore = 90
	if err := d.SaveJob(updated); err != nil {
		t.Fatalf("second SaveJob: %v", err)
	}

	jobs, err := d.ListJobs("")
	if err != nil {
		t.Fatalf("ListJobs: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("want 1 job after upsert, got %d", len(jobs))
	}
	if jobs[0].Status != "applied" {
		t.Errorf("want status applied, got %s", jobs[0].Status)
	}
	if jobs[0].FitScore != 90 {
		t.Errorf("want fit_score 90, got %d", jobs[0].FitScore)
	}
}

func TestListJobs_FilterByStatus(t *testing.T) {
	d := openTestDB(t)

	if err := d.SaveJob(testParams("job-1", "to_apply")); err != nil {
		t.Fatalf("SaveJob job-1: %v", err)
	}
	if err := d.SaveJob(testParams("job-2", "rejected")); err != nil {
		t.Fatalf("SaveJob job-2: %v", err)
	}

	jobs, err := d.ListJobs("to_apply")
	if err != nil {
		t.Fatalf("ListJobs: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("want 1 job, got %d", len(jobs))
	}
	if jobs[0].ID != "job-1" {
		t.Errorf("want job-1, got %s", jobs[0].ID)
	}
}

func TestListJobs_NoFilter(t *testing.T) {
	d := openTestDB(t)

	if err := d.SaveJob(testParams("job-1", "to_apply")); err != nil {
		t.Fatalf("SaveJob job-1: %v", err)
	}
	if err := d.SaveJob(testParams("job-2", "rejected")); err != nil {
		t.Fatalf("SaveJob job-2: %v", err)
	}

	jobs, err := d.ListJobs("")
	if err != nil {
		t.Fatalf("ListJobs: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("want 2 jobs, got %d", len(jobs))
	}
}

func TestGetStatusSummary_Counts(t *testing.T) {
	d := openTestDB(t)

	if err := d.SaveJob(testParams("job-1", "to_apply")); err != nil {
		t.Fatalf("SaveJob: %v", err)
	}
	if err := d.SaveJob(testParams("job-2", "to_apply")); err != nil {
		t.Fatalf("SaveJob: %v", err)
	}
	if err := d.SaveJob(testParams("job-3", "applied")); err != nil {
		t.Fatalf("SaveJob: %v", err)
	}

	s, err := d.GetStatusSummary()
	if err != nil {
		t.Fatalf("GetStatusSummary: %v", err)
	}
	if s.ToApply != 2 {
		t.Errorf("want to_apply=2, got %d", s.ToApply)
	}
	if s.Applied != 1 {
		t.Errorf("want applied=1, got %d", s.Applied)
	}
	if s.Total != 3 {
		t.Errorf("want total=3, got %d", s.Total)
	}
}

func TestGetStatusSummary_EmptyDB(t *testing.T) {
	d := openTestDB(t)

	s, err := d.GetStatusSummary()
	if err != nil {
		t.Fatalf("GetStatusSummary on empty DB: %v", err)
	}
	if s.Total != 0 {
		t.Errorf("want total=0, got %d", s.Total)
	}
}
