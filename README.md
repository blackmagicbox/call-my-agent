# Call My Agent — MCP Server

An MCP (Model Context Protocol) server designed to work with the Claude chrome extension, helping users streamline their job search workflow.

## Overview

Call My Agent provides tools that enable Claude to:

- **Job Fitness Assessment** — Evaluate job listings against user-defined search criteria and preferences, providing a fitness score and summary.
- **Cover Letter Generation** — Generate tailored cover letters matched to specific job postings, incorporating the user's experience and the role's requirements.
- **Search Criteria Management** — Store and manage the user's job search preferences, skills, and career goals to drive consistent evaluations.

## Prerequisites

- Go 1.26+

## Installation

```bash
go install github.com/blackmagicbox/call-my-agent@latest
```

## Usage

```bash
call-my-agent
```

Configure the Claude chrome extension to connect to this MCP server to enable job search assistance directly in your browser.

## Configuration

_Coming soon_ — details on how to define search criteria, experience profiles, and server options.

## License

_TBD_

