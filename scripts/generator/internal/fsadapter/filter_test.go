package fsadapter

import (
	"path/filepath"
	"testing"
)

func TestIsMarkdown(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "lowercase", input: "readme.md", expected: true},
		{name: "uppercase", input: "README.MD", expected: true},
		{name: "other", input: "readme.txt", expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := IsMarkdown(tc.input); got != tc.expected {
				t.Fatalf("IsMarkdown(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestIsExcluded(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "agents-md", input: "Agents.md", expected: true},
		{name: "skills-md", input: filepath.Join("skills", "skills.md"), expected: true},
		{name: "other", input: filepath.Join("skills", "roles", "scout.md"), expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := IsExcluded(tc.input); got != tc.expected {
				t.Fatalf("IsExcluded(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}
