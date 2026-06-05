package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

type SaveJobParams struct {
	ID             string
	Title          string
	Company        string
	Location       string
	Remote         string
	URL            string
	Status         string
	FitScore       int
	Recommendation string
	MatchedSkills  string // JSON array string
	SkillGaps      string // JSON array string
	RedFlagsHit    string // JSON array string
	Reasoning      string
	JobJSON        string
	SeenAt         string // RFC3339
}

type JobRow struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Company        string `json:"company"`
	Location       string `json:"location"`
	URL            string `json:"url"`
	Status         string `json:"status"`
	FitScore       int    `json:"fit_score"`
	Recommendation string `json:"recommendation"`
	Reasoning      string `json:"reasoning"`
	SeenAt         string `json:"seen_at"`
	UpdatedAt      string `json:"updated_at"`
}

type StatusSummary struct {
	ToEvaluate int `json:"to_evaluate"`
	ToApply    int `json:"to_apply"`
	Applied    int `json:"applied"`
	Rejected   int `json:"rejected"`
	Archived   int `json:"archived"`
	Total      int `json:"total"`
}

// Open opens (or creates) the SQLite database at path.
// Use ":memory:" in tests.
func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
			id              TEXT PRIMARY KEY,
			title           TEXT NOT NULL,
			company         TEXT NOT NULL,
			location        TEXT,
			remote          TEXT,
			url             TEXT,
			status          TEXT NOT NULL DEFAULT 'to_evaluate',
			fit_score       INTEGER,
			recommendation  TEXT,
			matched_skills  TEXT,
			skill_gaps      TEXT,
			red_flags_hit   TEXT,
			reasoning       TEXT,
			job_json        TEXT,
			seen_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
	`)
	return err
}

func (db *DB) SaveJob(p SaveJobParams) error {
	_, err := db.conn.Exec(`
		INSERT INTO jobs (
			id, title, company, location, remote, url, status,
			fit_score, recommendation, matched_skills, skill_gaps,
			red_flags_hit, reasoning, job_json, seen_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			status         = excluded.status,
			fit_score      = excluded.fit_score,
			recommendation = excluded.recommendation,
			matched_skills = excluded.matched_skills,
			skill_gaps     = excluded.skill_gaps,
			red_flags_hit  = excluded.red_flags_hit,
			reasoning      = excluded.reasoning,
			job_json       = excluded.job_json,
			updated_at     = CURRENT_TIMESTAMP
	`,
		p.ID, p.Title, p.Company, p.Location, p.Remote, p.URL, p.Status,
		p.FitScore, p.Recommendation, p.MatchedSkills, p.SkillGaps,
		p.RedFlagsHit, p.Reasoning, p.JobJSON, p.SeenAt,
	)
	return err
}

func (db *DB) ListJobs(status string) ([]JobRow, error) {
	query := `SELECT id, title, company, location, url, status,
	                 fit_score, recommendation, reasoning, seen_at, updated_at
	          FROM jobs`
	var args []any
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY seen_at DESC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var jobs []JobRow
	for rows.Next() {
		var j JobRow
		if err := rows.Scan(
			&j.ID, &j.Title, &j.Company, &j.Location, &j.URL,
			&j.Status, &j.FitScore, &j.Recommendation, &j.Reasoning,
			&j.SeenAt, &j.UpdatedAt,
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func (db *DB) GetStatusSummary() (*StatusSummary, error) {
	rows, err := db.conn.Query(`SELECT status, COUNT(*) FROM jobs GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	s := &StatusSummary{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		switch status {
		case "to_evaluate":
			s.ToEvaluate = count
		case "to_apply":
			s.ToApply = count
		case "applied":
			s.Applied = count
		case "rejected":
			s.Rejected = count
		case "archived":
			s.Archived = count
		}
		s.Total += count
	}
	return s, rows.Err()
}
