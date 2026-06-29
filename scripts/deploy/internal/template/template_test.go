package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpand(t *testing.T) {
	vars := Vars{
		"codex": map[string]any{
			"model": "o3-pro",
			"sandbox": map[string]any{
				"writable_roots": []any{"/tmp", "/home"},
			},
			"debug":  true,
			"port":   float64(8080),
			"ratio":  float64(3.14),
			"empty":  []any{},
			"nested": map[string]any{"key": "val"},
			"table": map[string]any{
				"~/code/shared-lib": true,
				"~/code/app":        true,
			},
		},
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple variable",
			input: `model = "{{codex.model}}"`,
			want:  `model = "o3-pro"`,
		},
		{
			name:  "array variable",
			input: `writable_roots = {{codex.sandbox.writable_roots}}`,
			want:  `writable_roots = ["/tmp", "/home"]`,
		},
		{
			name:  "empty array",
			input: `roots = {{codex.empty}}`,
			want:  `roots = []`,
		},
		{
			name:  "bool variable",
			input: `debug = {{codex.debug}}`,
			want:  `debug = true`,
		},
		{
			name:  "integer variable",
			input: `port = {{codex.port}}`,
			want:  `port = 8080`,
		},
		{
			name:  "float variable",
			input: `ratio = {{codex.ratio}}`,
			want:  `ratio = 3.14`,
		},
		{
			name:  "multiple variables",
			input: `model = "{{codex.model}}" roots = {{codex.sandbox.writable_roots}}`,
			want:  `model = "o3-pro" roots = ["/tmp", "/home"]`,
		},
		{
			name:  "no variables",
			input: `model = "gpt-5.5"`,
			want:  `model = "gpt-5.5"`,
		},
		{
			name:  "env var syntax preserved",
			input: `key = "${ENV_VAR}"`,
			want:  `key = "${ENV_VAR}"`,
		},
		{
			name:    "undefined variable",
			input:   `model = "{{codex.unknown}}"`,
			wantErr: true,
		},
		{
			name:    "intermediate not object",
			input:   `val = "{{codex.model.sub}}"`,
			wantErr: true,
		},
		{
			name:  "object variable as table entries",
			input: `[permissions.project-edit.workspace_roots]` + "\n" + `{{codex.table}}`,
			want: `[permissions.project-edit.workspace_roots]
"~/code/app" = true
"~/code/shared-lib" = true`,
		},
		{
			name:  "single object variable as table entry",
			input: `{{codex.nested}}`,
			want:  `"key" = "val"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.input, vars)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLoadVars(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.local.json")

	content := `{"codex": {"model": "o3-pro"}}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	vars, err := LoadVars(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := Expand("{{codex.model}}", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "o3-pro" {
		t.Errorf("got %q, want %q", got, "o3-pro")
	}
}

func TestLoadVars_TrailingComma(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.local.json")

	content := `{
  "codex": {
    "sandbox": {
      "writable_roots": [
        "/tmp",
        "/home",
      ],
    },
  },
}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	vars, err := LoadVars(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := Expand("roots = {{codex.sandbox.writable_roots}}", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `roots = ["/tmp", "/home"]`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLoadVars_NotFound(t *testing.T) {
	_, err := LoadVars("/nonexistent/config.local.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestExpandFile(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "config.toml")
	content := `model = "{{codex.model}}"` + "\n"
	if err := os.WriteFile(tmplPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	vars := Vars{"codex": map[string]any{"model": "o3-pro"}}
	got, err := ExpandFile(tmplPath, vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "model = \"o3-pro\"\n"
	if string(got) != want {
		t.Errorf("got %q, want %q", string(got), want)
	}
}
