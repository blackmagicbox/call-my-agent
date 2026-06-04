package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

type SaveJobParams struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Company        string    `json:"company"`
	Location       string    `json:"location"`
	Remote         string    `json:"remote"`
	URL            string    `json:"url"`
	Status         string    `json:"status"`
	FitScore       int       `json:"fit_score"`
	Recommendation string    `json:"recommendation"`
	MatchedSkills  string    `json:"matched_skills"`
	SkillGaps      string    `json:"skill_gaps"`
	RedFlags       string    `json:"red_flags"`
	Reasoning      string    `json:"reasoning"`
	JobJSON        string    `json:"job_json"`
	SeenAt         time.Time `json:"seen_at"`
}

type JobRow struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Company        string    `json:"company"`
	Location       string    `json:"location"`
	URL            string    `json:"url"`
	Status         string    `json:"status"`
	FitScore       int       `json:"fit_score"`
	Recommendation string    `json:"recommendation"`
	Reasoning      string    `json:"reasoning"`
	SeenAt         time.Time `json:"seen_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (db *DB) ListJobs(status string) ([]JobRow, error) {
	sql
}

func (db *DB) SaveJob(params SaveJobParams) error {
	sqlStmt := `INSERT INTO jobs(
    	id, title, company, location, remote, url, status,
        fit_score, recommendation, matched_skills,
        skill_gaps, red_flags_hit, reasoning,
        job_json, seen_at, updated_at) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT DO UPDATE SET
            status = excluded.status,
            fit_score = excluded.fit_score,
			recommendation = excluded.recommendation,
			matched_skills = excluded.matched_skills,
			skill_gaps = excluded.skill_gaps,
			red_flags_hit = excluded.red_flags_hit,
			reasoning = excluded.reasoning,
			job_json = excluded.job_json,
			updated_at = CURRENT_TIMESTAMP();`

	_, err := db.conn.Exec(sqlStmt, params)
	return err
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS jobs (
 		id TEXT PRIMARY KEY,
 		title TEXT NOT NULL,
 		company TEXT NOT NULL,
 		location TEXT,
 		remote TEXT,
 		url TEXT,
 		status TEXT NOT NULL DEFAULT 'to_evaluate',
 		fit_score INTEGER DEFAULT 0,
 		recommendation TEXT,
 		matched_skills TEXT,
 		skill_gaps TEXT,
 		red_flags_hit TEXT,
 		reasoning TEXT,
 		job_json TEXT,
 		seen_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);
	CREATE INDEX IF NOT EXISTS jobs_idx_status ON jobs(status);`

	_, err := db.conn.Exec(sqlStmt)

	return err
}
