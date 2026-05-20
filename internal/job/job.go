package job

import "time"

type Status string

const (
	StatusToApply  Status = "to_apply"
	StatusApplied  Status = "applied"
	StatusRejected Status = "rejected"
	StatusArchived Status = "archived"
)

type Job struct {
	ID          int64
	URL         string
	Title       string
	Company     string
	Description string
	FitScore    int
	FitReason   string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
