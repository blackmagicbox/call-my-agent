package llm

import (
	"errors"
	"fmt"
	"os"
)

func FromConfig() (Provider, error) {
	provider := os.Getenv("LLM_PROVIDER")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")

	if provider == "" {
		return nil, errors.New("provider not set")
	}

	switch provider {
	case "gemini":
		if apiKey == "" {
			return nil, errors.New("missing api key for the selected provider")
		}
		return NewGeminiProvider(apiKey, model), nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}
