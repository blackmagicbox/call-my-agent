// Package testutil provides shared helpers for test code across the project.
// Import it in any _test.go file that needs fixture loading or common unmarshalling.
package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// LoadFixture reads a JSON fixture file from the top-level testdata/ directory.
// The path is resolved relative to this file's location so it works regardless
// of which package calls it.
func LoadFixture(t *testing.T, name string) []byte {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("testutil: runtime.Caller failed — cannot resolve fixture path")
	}
	root := filepath.Join(filepath.Dir(filename), "..", "..", "testdata")
	data, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatalf("testutil: could not load fixture %q: %v", name, err)
	}
	return data
}

// MustUnmarshal unmarshals JSON bytes into T or immediately fails the test.
func MustUnmarshal[T any](t *testing.T, data []byte) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("testutil: unmarshal failed: %v", err)
	}
	return v
}
