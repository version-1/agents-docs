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

func TestRunnerExcludesFilesByGlob(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src")
	for _, dir := range []string{
		filepath.Join(srcDir, "nested"),
		filepath.Join(srcDir, "cache"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	for path, content := range map[string]string{
		filepath.Join(srcDir, "keep.txt"):             "keep",
		filepath.Join(srcDir, "ignore.tmp"):           "tmp",
		filepath.Join(srcDir, "nested", "keep.md"):    "md",
		filepath.Join(srcDir, "nested", "ignore.log"): "log",
		filepath.Join(srcDir, "cache", "data.txt"):    "cache",
	} {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {
      "source": "src",
      "destination": "dest",
      "exclude": ["*.tmp", "**/*.log", "cache/**"]
    }
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runner.Run(config, Options{}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(root, "dest", "keep.txt"), "keep")
	assertFileContent(t, filepath.Join(root, "dest", "nested", "keep.md"), "md")
	assertNotExist(t, filepath.Join(root, "dest", "ignore.tmp"))
	assertNotExist(t, filepath.Join(root, "dest", "nested", "ignore.log"))
	assertNotExist(t, filepath.Join(root, "dest", "cache", "data.txt"))
	if !strings.Contains(out.String(), "SKIP") {
		t.Fatalf("expected skip output, got:\n%s", out.String())
	}
}

func TestRunnerExcludesSingleFileByGlob(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "secret.txt")
	if err := os.WriteFile(src, []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "secret.txt", "destination": "dest.txt", "exclude": ["secret.*"]}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runner.Run(config, Options{}); err != nil {
		t.Fatal(err)
	}

	assertNotExist(t, filepath.Join(root, "dest.txt"))
	if !strings.Contains(out.String(), "SKIP") {
		t.Fatalf("expected skip output, got:\n%s", out.String())
	}
}

func TestRunnerReplaceRemovesDestinationBeforeCopyingDirectory(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src")
	dstDir := filepath.Join(root, "dest")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "current.txt"), []byte("current"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dstDir, "stale.txt"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "src", "destination": "dest", "replace": true}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runner.Run(config, Options{}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dstDir, "current.txt"), "current")
	assertNotExist(t, filepath.Join(dstDir, "stale.txt"))
	if !strings.Contains(out.String(), "REMOVE") {
		t.Fatalf("expected remove output, got:\n%s", out.String())
	}
}

func TestRunnerKeepsDestinationExtrasWhenReplaceIsFalse(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src")
	dstDir := filepath.Join(root, "dest")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "current.txt"), []byte("current"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dstDir, "stale.txt"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "src", "destination": "dest"}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runner.Run(config, Options{}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dstDir, "current.txt"), "current")
	assertFileContent(t, filepath.Join(dstDir, "stale.txt"), "stale")
}

func TestRunnerDryRunReplaceDoesNotRemoveDestination(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src")
	dstDir := filepath.Join(root, "dest")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "current.txt"), []byte("current"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dstDir, "stale.txt"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "src", "destination": "dest", "replace": true}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runner.Run(config, Options{DryRun: true}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dstDir, "stale.txt"), "stale")
	if !strings.Contains(out.String(), "REMOVE") {
		t.Fatalf("expected remove output, got:\n%s", out.String())
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

func assertNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s not to exist: %v", path, err)
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
