package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDirCopiesNestedFilesAndEmptyDirs(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	if err := os.MkdirAll(filepath.Join(src, "nested", "empty"), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "nested", "file.txt"), []byte("content"), 0640); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dst, "existing"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := CopyDir(src, dst); err != nil {
		t.Fatal(err)
	}

	assertContent(t, filepath.Join(dst, "nested", "file.txt"), "content")
	assertMode(t, filepath.Join(dst, "nested", "file.txt"), 0640)
	assertDir(t, filepath.Join(dst, "nested", "empty"))
	assertDir(t, filepath.Join(dst, "existing"))
}

func assertContent(t *testing.T, path, want string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != want {
		t.Fatalf("content mismatch: got %q want %q", string(b), want)
	}
}

func assertMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Fatalf("mode mismatch: got %v want %v", got, want)
	}
}

func assertDir(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatalf("expected directory: %s", path)
	}
}
