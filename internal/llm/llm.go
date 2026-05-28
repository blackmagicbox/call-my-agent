package llm

import "context"

// Provider abstracts an LLM backend so callers can swap between
// Gemini, Claude, or any other model without changing business logic.
type Provider interface {
	// Complete sends a single-turn request with a system prompt and a user
	// message, returning the model's text response.
	Complete(ctx context.Context, system string, user string) (string, error)
}
