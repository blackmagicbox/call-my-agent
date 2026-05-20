package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/blackmagicbox/call-my-agent/internal/job"
	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	d := &DB{conn: conn}
	if err := d.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return d, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) migrate() error {
	_, err := d.conn.Exec(`CREATE TABLE IF NOT EXISTS jobs (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		url         TEXT    NOT NULL UNIQUE,
		title       TEXT    NOT NULL DEFAULT '',
		company     TEXT    NOT NULL DEFAULT '',
		description TEXT    NOT NULL DEFAULT '',
		fit_score   INTEGER NOT NULL DEFAULT 0,
		fit_reason  TEXT    NOT NULL DEFAULT '',
		status      TEXT    NOT NULL DEFAULT 'to_apply',
		created_at  DATETIME NOT NULL,
		updated_at  DATETIME NOT NULL
	)`)
	return err
}

// SaveJob inserts or updates a job record keyed on URL.
func (d *DB) SaveJob(j *job.Job) error {
	now := time.Now().UTC()
	if j.CreatedAt.IsZero() {
		j.CreatedAt = now
	}
	j.UpdatedAt = now

	if j.Status == "" {
		j.Status = job.StatusToApply
	}

	res, err := d.conn.Exec(`
		INSERT INTO jobs (url, title, company, description, fit_score, fit_reason, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(url) DO UPDATE SET
			title       = excluded.title,
			company     = excluded.company,
			description = excluded.description,
			fit_score   = excluded.fit_score,
			fit_reason  = excluded.fit_reason,
			status      = excluded.status,
			updated_at  = excluded.updated_at`,
		j.URL, j.Title, j.Company, j.Description,
		j.FitScore, j.FitReason, j.Status,
		j.CreatedAt, j.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save job: %w", err)
	}
	if j.ID == 0 {
		j.ID, _ = res.LastInsertId()
	}
	return nil
}

// ListJobs returns all tracked jobs, optionally filtered by status.
// Pass an empty string to return all.
func (d *DB) ListJobs(status string) ([]*job.Job, error) {
	var (
		rows *sql.Rows
		err  error
	)
	const base = `SELECT id, url, title, company, description, fit_score, fit_reason, status, created_at, updated_at FROM jobs`
	if status != "" {
		rows, err = d.conn.Query(base+` WHERE status = ? ORDER BY created_at DESC`, status)
	} else {
		rows, err = d.conn.Query(base + ` ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*job.Job
	for rows.Next() {
		j := &job.Job{}
		if err := rows.Scan(
			&j.ID, &j.URL, &j.Title, &j.Company, &j.Description,
			&j.FitScore, &j.FitReason, &j.Status,
			&j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

type ApplicationStats struct {
	Total    int
	ToApply  int
	Applied  int
	Rejected int
	Archived int
}

// ApplicationStats returns a count of jobs in each status.
func (d *DB) ApplicationStats() (*ApplicationStats, error) {
	row := d.conn.QueryRow(`
		SELECT
			COUNT(*),
			SUM(CASE WHEN status = 'to_apply'  THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'applied'   THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'rejected'  THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'archived'  THEN 1 ELSE 0 END)
		FROM jobs`)
	var s ApplicationStats
	if err := row.Scan(&s.Total, &s.ToApply, &s.Applied, &s.Rejected, &s.Archived); err != nil {
		return nil, fmt.Errorf("application stats: %w", err)
	}
	return &s, nil
}
