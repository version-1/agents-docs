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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{DryRun: true, NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, "dest.txt")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote destination file: %v", err)
	}
	if !strings.Contains(out.String(), "DRY-RUN") || !strings.Contains(out.String(), "copied") {
		t.Fatalf("expected dry-run copy output, got:\n%s", out.String())
	}
}

func TestRunnerResolvesSourceFromCurrentDirectory(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "src.txt"), []byte("src"), 0644); err != nil {
		t.Fatal(err)
	}

	config := filepath.Join(configDir, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "src.txt", "destination": "../dest.txt"}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(root, "dest.txt"), "src")
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(root, "dest", "keep.txt"), "keep")
	assertFileContent(t, filepath.Join(root, "dest", "nested", "keep.md"), "md")
	assertNotExist(t, filepath.Join(root, "dest", "ignore.tmp"))
	assertNotExist(t, filepath.Join(root, "dest", "nested", "ignore.log"))
	assertNotExist(t, filepath.Join(root, "dest", "cache", "data.txt"))
	if !strings.Contains(out.String(), "skipped") {
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertNotExist(t, filepath.Join(root, "dest.txt"))
	if !strings.Contains(out.String(), "skipped") {
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dstDir, "current.txt"), "current")
	assertNotExist(t, filepath.Join(dstDir, "stale.txt"))
	if !strings.Contains(out.String(), "replace:") {
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dstDir, "current.txt"), "current")
	assertFileContent(t, filepath.Join(dstDir, "stale.txt"), "stale")
}

func TestRunnerFlattenCopiesSkillDirsToDestinationRoot(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src")
	for _, dir := range []string{
		filepath.Join(srcDir, "internal", "role-planner", "assets"),
		filepath.Join(srcDir, "external", "empirical-prompt-tuning"),
		filepath.Join(srcDir, "internal", "not-a-skill"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	for path, content := range map[string]string{
		filepath.Join(srcDir, "internal", "role-planner", "SKILL.md"):              "planner",
		filepath.Join(srcDir, "internal", "role-planner", "assets", "prompt.md"):   "prompt",
		filepath.Join(srcDir, "external", "empirical-prompt-tuning", "SKILL.md"):   "external",
		filepath.Join(srcDir, "internal", "not-a-skill", "README.md"):              "readme",
		filepath.Join(srcDir, "internal", "role-planner", "assets", "ignored.tmp"): "tmp",
		filepath.Join(srcDir, "external", "empirical-prompt-tuning", ".DS_Store"):  "store",
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
      "flatten": true,
      "exclude": ["*.tmp", "**/.DS_Store"]
    }
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(root, "dest", "role-planner", "SKILL.md"), "planner")
	assertFileContent(t, filepath.Join(root, "dest", "role-planner", "assets", "prompt.md"), "prompt")
	assertFileContent(t, filepath.Join(root, "dest", "empirical-prompt-tuning", "SKILL.md"), "external")
	assertNotExist(t, filepath.Join(root, "dest", "internal", "role-planner", "SKILL.md"))
	assertNotExist(t, filepath.Join(root, "dest", "not-a-skill", "README.md"))
	assertNotExist(t, filepath.Join(root, "dest", "role-planner", "assets", "ignored.tmp"))
	assertNotExist(t, filepath.Join(root, "dest", "empirical-prompt-tuning", ".DS_Store"))
	if !strings.Contains(out.String(), "flattened-skill-dirs") {
		t.Fatalf("expected flatten output, got:\n%s", out.String())
	}
}

func TestRunnerFlattenRejectsDuplicateSkillDirNames(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{
		filepath.Join(root, "src", "internal", "same"),
		filepath.Join(root, "src", "external", "same"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(filepath.Base(filepath.Dir(dir))), 0644); err != nil {
			t.Fatal(err)
		}
	}

	config := filepath.Join(root, "deploy.json")
	writeConfig(t, config, `{
  "items": [
    {"source": "src", "destination": "dest", "flatten": true}
  ]
}`)

	var out bytes.Buffer
	runner := NewRunner(&out)
	err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) })
	if err == nil {
		t.Fatal("expected duplicate target error")
	}
	if !strings.Contains(err.Error(), "flatten target conflict") {
		t.Fatalf("expected conflict error, got %v", err)
	}
	assertNotExist(t, filepath.Join(root, "dest"))
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{DryRun: true, NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dstDir, "stale.txt"), "stale")
	if !strings.Contains(out.String(), "replace:") {
		t.Fatalf("expected remove output, got:\n%s", out.String())
	}
}

func TestRunnerBacksUpExistingDestinationFile(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	dst := filepath.Join(root, "dest.txt")
	if err := os.WriteFile(src, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("old"), 0644); err != nil {
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, dst, "new")
	backupPath := backupPathFromOutput(t, out.String())
	assertFileContent(t, backupPath, "old")
	if !strings.Contains(backupPath, filepath.Join(".deploy-backups")) {
		t.Fatalf("expected backup path under .deploy-backups, got %s", backupPath)
	}
}

func TestRunnerBacksUpExistingDestinationDirectoryBeforeReplace(t *testing.T) {
	root := t.TempDir()
	srcDir := filepath.Join(root, "src")
	dstDir := filepath.Join(root, "dest")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dstDir, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "current.txt"), []byte("current"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dstDir, "nested", "old.txt"), []byte("old"), 0644); err != nil {
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dstDir, "current.txt"), "current")
	assertNotExist(t, filepath.Join(dstDir, "nested", "old.txt"))
	backupPath := backupPathFromOutput(t, out.String())
	assertFileContent(t, filepath.Join(backupPath, "nested", "old.txt"), "old")
}

func TestRunnerDryRunBackupDoesNotWriteBackup(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	dst := filepath.Join(root, "dest.txt")
	if err := os.WriteFile(src, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("old"), 0644); err != nil {
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
	if err := runFromDir(t, root, func() error { return runner.Run(config, Options{DryRun: true, NoColor: true}) }); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, dst, "old")
	backupPath := backupPathFromOutput(t, out.String())
	assertNotExist(t, backupPath)
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

func backupPathFromOutput(t *testing.T, output string) string {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "backup: ") {
			return strings.TrimPrefix(line, "backup: ")
		}
	}
	t.Fatalf("backup output not found:\n%s", output)
	return ""
}

func runFromDir(t *testing.T, dir string, fn func() error) error {
	t.Helper()

	current, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(current); err != nil {
			t.Fatal(err)
		}
	}()

	return fn()
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
