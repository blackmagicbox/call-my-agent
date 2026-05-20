package db_test

import (
	"testing"

	"github.com/blackmagicbox/call-my-agent/internal/db"
	"github.com/blackmagicbox/call-my-agent/internal/job"
)

func newTestDB(t *testing.T) *db.DB {
	t.Helper()
	d, err := db.New(":memory:")
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

func TestSaveJob_Insert(t *testing.T) {
	d := newTestDB(t)

	j := &job.Job{
		URL:       "https://example.com/job/1",
		Title:     "SRE",
		Company:   "Acme",
		FitScore:  85,
		FitReason: "matches infra background",
	}
	if err := d.SaveJob(j); err != nil {
		t.Fatalf("save: %v", err)
	}
	if j.ID == 0 {
		t.Fatal("expected non-zero ID after insert")
	}
	if j.Status != job.StatusToApply {
		t.Errorf("default status: got %q, want %q", j.Status, job.StatusToApply)
	}
}

func TestSaveJob_Upsert(t *testing.T) {
	d := newTestDB(t)

	j := &job.Job{URL: "https://example.com/job/2", FitScore: 50}
	if err := d.SaveJob(j); err != nil {
		t.Fatalf("first save: %v", err)
	}

	j.FitScore = 90
	j.Status = job.StatusApplied
	if err := d.SaveJob(j); err != nil {
		t.Fatalf("second save: %v", err)
	}

	jobs, err := d.ListJobs("")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].FitScore != 90 {
		t.Errorf("fit score: got %d, want 90", jobs[0].FitScore)
	}
	if jobs[0].Status != job.StatusApplied {
		t.Errorf("status: got %q, want %q", jobs[0].Status, job.StatusApplied)
	}
}

func TestListJobs_FilterByStatus(t *testing.T) {
	d := newTestDB(t)

	jobs := []*job.Job{
		{URL: "https://example.com/1", Status: job.StatusToApply},
		{URL: "https://example.com/2", Status: job.StatusApplied},
		{URL: "https://example.com/3", Status: job.StatusApplied},
		{URL: "https://example.com/4", Status: job.StatusRejected},
	}
	for _, j := range jobs {
		if err := d.SaveJob(j); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	applied, err := d.ListJobs("applied")
	if err != nil {
		t.Fatalf("list applied: %v", err)
	}
	if len(applied) != 2 {
		t.Errorf("applied count: got %d, want 2", len(applied))
	}

	all, err := d.ListJobs("")
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(all) != 4 {
		t.Errorf("all count: got %d, want 4", len(all))
	}
}

func TestApplicationStats(t *testing.T) {
	d := newTestDB(t)

	for _, j := range []*job.Job{
		{URL: "https://example.com/1", Status: job.StatusToApply},
		{URL: "https://example.com/2", Status: job.StatusApplied},
		{URL: "https://example.com/3", Status: job.StatusApplied},
		{URL: "https://example.com/4", Status: job.StatusRejected},
		{URL: "https://example.com/5", Status: job.StatusArchived},
	} {
		if err := d.SaveJob(j); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	stats, err := d.ApplicationStats()
	if err != nil {
		t.Fatalf("stats: %v", err)
	}

	if stats.Total != 5 {
		t.Errorf("total: got %d, want 5", stats.Total)
	}
	if stats.ToApply != 1 {
		t.Errorf("to_apply: got %d, want 1", stats.ToApply)
	}
	if stats.Applied != 2 {
		t.Errorf("applied: got %d, want 2", stats.Applied)
	}
	if stats.Rejected != 1 {
		t.Errorf("rejected: got %d, want 1", stats.Rejected)
	}
	if stats.Archived != 1 {
		t.Errorf("archived: got %d, want 1", stats.Archived)
	}
}
