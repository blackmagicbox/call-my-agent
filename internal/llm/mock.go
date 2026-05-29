package llm

import "context"

// MockProvider is a test double for Provider. Set Response/Err before each call.
type MockProvider struct {
	Response string
	Err      error
	// LastSystem and LastUser capture the most recent call for assertion in tests.
	LastSystem string
	LastUser   string
}

func (m *MockProvider) Complete(_ context.Context, system, user string) (string, error) {
	m.LastSystem = system
	m.LastUser = user
	return m.Response, m.Err
}
