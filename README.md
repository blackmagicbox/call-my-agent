# Call My Agent

![CI](https://github.com/blackmagicbox/call-my-agent/actions/workflows/ci.yml/badge.svg)

An MCP server that turns Claude into your personal job-hunting assistant. Browse LinkedIn, open a listing, and ask Claude to evaluate it — you get a fit score, a tailored cover letter, and the job saved to a local database, all without leaving the browser.

---

## How it works

Call My Agent runs as a local MCP server. Claude connects to it and gains access to these tools:

| Tool | What it does |
|------|-------------|
| `get_candidate_profile` | Gives Claude full context about you — skills, experience, preferences |
| `evaluate_job` | Scores a job listing against your profile (0–100 fit score + reasoning) |
| `generate_cover_letter` | Writes a tailored cover letter for a specific listing |
| `save_job` | Saves an evaluated job to a local SQLite database |
| `list_jobs` | Lists your saved jobs, optionally filtered by status |
| `get_application_status` | Summary of your pipeline (to apply, applied, rejected, etc.) |

You browse normally. Claude reads the page. You ask. It answers.

---

## Prerequisites

- **Go 1.23+** — [install](https://go.dev/dl/)
- **Claude Desktop** or **Claude in Chrome** — [download](https://claude.ai/download)
- A **Gemini API key** — [get a free key at Google AI Studio](https://aistudio.google.com/apikey)

---

## Setup

### 1. Clone and build

```bash
git clone https://github.com/blackmagicbox/call-my-agent.git
cd call-my-agent
go build -o call-my-agent ./cmd/server
```

### 2. Set up your profile

Your profile is what Claude uses to evaluate every job. Copy the example and fill in your details:

```bash
cp data/profile.example.json data/profile.json
```

Open `data/profile.json` in any editor. The key sections to fill in:

```jsonc
{
  "name": "Your Name",
  "title": "Senior SRE",
  "skills": ["Go", "Kubernetes", "Terraform", "GCP"],
  "experience": [ ... ],          // your work history with bullet points
  "preferences": {
    "target_roles":   ["Senior SRE", "Staff SRE"],
    "target_levels":  ["L5", "L6"],
    "locations":      ["Berlin", "Remote"],
    "remote_preference": "preferred",  // "required" | "preferred" | "not_interested"
    "red_flags":      ["no Kubernetes", "PHP"],
    "min_salary":     120000,
    "salary_currency": "EUR"
  }
}
```

The `red_flags` list is especially important — any job that matches one will be flagged in the evaluation.

> **Privacy:** `data/profile.json` is gitignored and never committed.

### 3. Connect to Claude Desktop

Find your Claude Desktop config file:

| OS | Path |
|----|------|
| macOS | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| Windows | `%APPDATA%\Claude\claude_desktop_config.json` |
| Linux | `~/.config/Claude/claude_desktop_config.json` |

Add the server entry:

```json
{
  "mcpServers": {
    "call-my-agent": {
      "command": "/absolute/path/to/call-my-agent",
      "env": {
        "LLM_PROVIDER": "gemini",
        "LLM_API_KEY": "your-key-here",
        "LLM_MODEL": "gemini-2.5-flash"
      }
    }
  }
}
```

Restart Claude Desktop. You should see the tools appear in the toolbar.

---

## Using it

Open a LinkedIn job listing in your browser, then ask Claude in natural language:

**Evaluate a job:**
> *"Evaluate this job for me."*

Claude will return:
- **Fit score** (0–100)
- **Recommendation** — `apply`, `consider`, or `skip`
- **Matched skills** and **skill gaps**
- **Red flags** triggered from your preferences
- **Reasoning** for the score

**Generate a cover letter:**
> *"Write me a cover letter for this one."*

**Save a job:**
> *"Save this job."*

**Review your pipeline:**
> *"Show me everything I've saved so far."*
> *"Show me all jobs I've marked to apply for."*
> *"What does my application pipeline look like?"*

---

## Job statuses

Saved jobs move through a simple pipeline:

```
to_evaluate → to_apply → applied → rejected / archived
```

You can filter by any status when listing jobs.

---

## Environment variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LLM_PROVIDER` | Yes | — | LLM backend. Currently only `gemini` is supported. |
| `LLM_API_KEY` | Yes | — | Your API key for the configured provider. |
| `LLM_MODEL` | No | `gemini-2.5-flash` | Model name to use. |

---

## Privacy

- Your profile and jobs database are local files — nothing is stored remotely.
- Job descriptions and your profile are sent to the LLM API when evaluating. Use a key with appropriate quota limits if privacy is a concern.
- The database lives at `data/jobs.db` (gitignored).

---

## Contributing / Development

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

[GNU Affero General Public License v3.0](LICENSE)
