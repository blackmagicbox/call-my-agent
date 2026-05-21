# Call My Agent

An MCP (Model Context Protocol) server that helps you job-hunt with Claude. It exposes your candidate profile, job evaluation tools, and cover letter generation as MCP tools — so Claude can assist you directly from the browser extension or Claude Desktop.

## Features

- **Candidate Profile** — Serve your profile (skills, experience, education, preferences) to Claude so it has full context when helping you.
- **Job Fitness Evaluation** *(planned)* — Score job listings against your defined criteria and preferences.
- **Cover Letter Generation** *(planned)* — Generate tailored cover letters matched to specific job postings.
- **Application Tracking** *(planned)* — Keep track of where you've applied and the status of each application.

## Prerequisites

- Go 1.23+

## Getting Started

### 1. Build

```bash
go build -o call-my-agent ./cmd/server
```

### 2. Configure your profile

Copy the example profile and fill in your details:

```bash
cp data/profile.example.json data/profile.json
```

Edit `data/profile.json` with your information — name, skills, experience, education, job preferences, etc.

### 3. Register with Claude Desktop

Add the server to your Claude Desktop configuration (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "call-my-agent": {
      "command": "/absolute/path/to/call-my-agent"
    }
  }
}
```

The server communicates over stdio using the MCP protocol.

## Development

```bash
go run ./cmd/server          # run directly
go test ./...                # run all tests
```

## Project Structure

```
cmd/server/       Application entrypoint
internal/
  profile/        Candidate profile data structures
  job/            Job posting data structures
  tools/          MCP tool handlers (stubs)
  db/             Persistence layer (stub)
data/
  profile.json          Your candidate profile
  profile.example.json  Template profile
```

## License

This project is licensed under the [GNU Affero General Public License v3.0](LICENSE).
