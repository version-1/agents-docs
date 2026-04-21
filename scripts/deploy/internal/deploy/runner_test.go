package deploy

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunnerCopiesFilesAndDirectoryContents(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src", "dir")
	if err := os.MkdirAll(filepath.Join(srcDir, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "nested", "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	srcFile := filepath.Join(root, "src", "single.txt")
	if err := os.WriteFile(srcFile, []byte("single"), 0600); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "src/dir", "destination": "dest/dir"},
    {"source": "src/single.txt", "destination": "dest/single.txt"}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runner.Run(config, Options{}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(root, "dest", "dir", "nested", "a.txt"), "a")
	assertFileContent(t, filepath.Join(root, "dest", "single.txt"), "single")
	if !strings.Contains(out.String(), "DEPLOY") {
		t.Fatalf("expected deploy output, got:\n%s", out.String())
	}
}

func TestRunnerDryRunDoesNotWriteFiles(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	if err := os.WriteFile(src, []byte("src"), 0644); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "src.txt", "destination": "dest.txt"}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runner.Run(config, Options{DryRun: true}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, "dest.txt")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote destination file: %v", err)
	}
	if !strings.Contains(out.String(), "DRY-RUN") || !strings.Contains(out.String(), "COPY") {
		t.Fatalf("expected dry-run copy output, got:\n%s", out.String())
	}
}

func TestLoadConfigRequiresItems(t *testing.T) {
	root := t.TempDir()
	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{"items":[]}`)

	_, err := LoadConfig(config)
	if err == nil {
		t.Fatal("expected error")
	}
}

func writeConfig(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != want {
		t.Fatalf("content mismatch for %s: got %q want %q", path, string(b), want)
	}
}
