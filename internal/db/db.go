package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
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

func (db *DB) Close() error {
	return db.conn.Close()
}
