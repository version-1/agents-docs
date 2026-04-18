package domain

import "testing"

func TestParseFrontmatter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		content     string
		wantName    string
		wantDesc    string
		expectErr   bool
	}{
		{
			name:     "name and description",
			content:  "---\nname: language-go\ndescription: Go skill\n---\n# content",
			wantName: "language-go",
			wantDesc: "Go skill",
		},
		{
			name:     "name only",
			content:  "---\nname: code-review\n---\nbody",
			wantName: "code-review",
			wantDesc: "",
		},
		{
			name:     "spaces around values",
			content:  "---\nname:   code-review  \ndescription:  some desc  \n---\nbody",
			wantName: "code-review",
			wantDesc: "some desc",
		},
		{
			name:      "no frontmatter",
			content:   "# Just a heading\nsome content",
			expectErr: true,
		},
		{
			name:      "missing name",
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
			got, err := ParseFrontmatter([]byte(tc.content))
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tc.wantName {
				t.Fatalf("Name = %q, want %q", got.Name, tc.wantName)
			}
			if got.Description != tc.wantDesc {
				t.Fatalf("Description = %q, want %q", got.Description, tc.wantDesc)
			}
		})
	}
}
