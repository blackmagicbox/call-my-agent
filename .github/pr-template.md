## MCP Server Skeleton & Candidate Profile Tool

This PR bootstraps the MCP (Model Context Protocol) server for "Call My Agent". Its main contribution is to:

- Implement the initial Go MCP server that communicates via stdio (no HTTP)
- Load a candidate profile from `data/profile.json` at startup
- Register the `get_candidate_profile` MCP tool, which exposes the loaded profile data (skills, experience, education, preferences, languages) to Claude or other compatible clients
- Add CEFR-based language proficiency support to the profile schema
- Revise documentation (README.md and CLAUDE.md) to clarify architecture, setup, and project guidelines

### **Why it matters?**

This PR defines the foundation for future MCP tools (job fit scoring, cover letter generation, job tracking) and ensures contributors have a clear starting point with standards for project structure, extensibility, and data-driven configuration. All project logic now cleanly loads from maintainable JSON rather than Go constants and is future-proofed for localization and additional features.

### **Getting started:**
- Fill out your `data/profile.json` before running the server.
- See README.md and CLAUDE.md for setup details, design principles, and technical roadmap.

Please refer to the tool registry and code organization guidelines before adding functionality. Your contributions will help expand the agent's capabilities for job seekers!
