package tools

import "strings"

// StripJSONFences strips ```json ... ``` wrapping that LLMs sometimes
// add despite being instructed not to.
func StripJSONFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
