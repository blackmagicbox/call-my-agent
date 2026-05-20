package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blackmagicbox/call-my-agent/internal/db"
	"github.com/blackmagicbox/call-my-agent/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	dbPath := dbFilePath()

	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	s := server.NewMCPServer("call-my-agent", "0.1.0",
		server.WithToolCapabilities(false),
	)

	tools.RegisterTracking(s, database)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func dbFilePath() string {
	if p := os.Getenv("CALL_MY_AGENT_DB"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "call-my-agent.db"
	}
	dir := filepath.Join(home, ".call-my-agent")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "call-my-agent.db"
	}
	return filepath.Join(dir, "jobs.db")
}
