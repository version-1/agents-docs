package domain

import "testing"

func TestParseName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		content   string
		expected  string
		expectErr bool
	}{
		{
			name:     "valid frontmatter",
			content:  "---\nname: language-go\ndescription: Go skill\n---\n# content",
			expected: "language-go",
		},
		{
			name:     "name with spaces around value",
			content:  "---\nname:   code-review  \n---\nbody",
			expected: "code-review",
		},
		{
			name:      "no frontmatter",
			content:   "# Just a heading\nsome content",
			expectErr: true,
		},
		{
			name:      "frontmatter without name",
			content:   "---\ndescription: something\n---\nbody",
			expectErr: true,
		},
		{
			name:      "empty name",
			content:   "---\nname:\n---\nbody",
			expectErr: true,
		},
		{
			name:      "unclosed frontmatter",
			content:   "---\nname: test\nno closing",
			expectErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseName([]byte(tc.content))
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("ParseName() = %q, want %q", got, tc.expected)
			}
		})
	}
}
