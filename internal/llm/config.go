package llm

import (
	"fmt"
	"os"
)

func FromConfig() (Provider, error) {
	provider := os.Getenv("LLM_PROVIDER")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")

	switch provider {
	case "gemini":
		if model == "" {
			model = "gemini-3.5-flash"
		}
		return NewGemini(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unknown LLM_PROVIDER: %s", provider)
	}
}
