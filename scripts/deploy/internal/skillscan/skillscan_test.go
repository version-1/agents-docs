package skillscan

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type exactMatcher map[string]bool

func (m exactMatcher) Match(rel string) bool {
	return m[rel]
}

func TestWalkSkillDirsFindsRegularSkillFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "internal", "coding", "SKILL.md"))
	writeFile(t, filepath.Join(root, "internal", "coding", "assets", "prompt.md"))
	writeFile(t, filepath.Join(root, "external", "grill-me", "SKILL.md"))
	writeFile(t, filepath.Join(root, "README.md"))

	var got []Dir
	err := WalkSkillDirs(Options{Root: root}, func(dir Dir) error {
		got = append(got, dir)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	wantNames := []string{"grill-me", "coding"}
	var gotNames []string
	for _, dir := range got {
		gotNames = append(gotNames, dir.Name)
		if filepath.Base(dir.SkillFile) != "SKILL.md" {
			t.Fatalf("unexpected skill file path: %s", dir.SkillFile)
		}
	}
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("skill names mismatch: got %v want %v", gotNames, wantNames)
	}
}

func TestWalkSkillDirsSkipsExcludedEntries(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "ignored", "SKILL.md"))
	writeFile(t, filepath.Join(root, "ignored", "nested", "SKILL.md"))
	writeFile(t, filepath.Join(root, "hidden", "SKILL.md"))
	writeFile(t, filepath.Join(root, "visible", "SKILL.md"))

	var skipped []string
	var got []string
	err := WalkSkillDirs(Options{
		Root:    root,
		Matcher: exactMatcher{"ignored": true, "hidden/SKILL.md": true},
		OnSkip: func(path string, reason SkipReason) {
			skipped = append(skipped, string(reason)+":"+filepath.ToSlash(mustRel(t, root, path)))
		},
	}, func(dir Dir) error {
		got = append(got, dir.Name)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, []string{"visible"}) {
		t.Fatalf("visited skills mismatch: got %v", got)
	}
	wantSkipped := []string{"excluded:hidden/SKILL.md", "excluded:ignored"}
	if !reflect.DeepEqual(skipped, wantSkipped) {
		t.Fatalf("skipped mismatch: got %v want %v", skipped, wantSkipped)
	}
}

func TestWalkSkillDirsSkipsNonRegularSkillFile(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "not-regular", "SKILL.md"), 0755); err != nil {
		t.Fatal(err)
	}

	var skipped []SkipReason
	err := WalkSkillDirs(Options{
		Root:   root,
		OnSkip: func(_ string, reason SkipReason) { skipped = append(skipped, reason) },
	}, func(Dir) error {
		t.Fatal("non-regular SKILL.md should not be visited")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(skipped, []SkipReason{SkipNonRegularSkill}) {
		t.Fatalf("skipped mismatch: got %v", skipped)
	}
}

func TestWalkSkillDirsReturnsVisitorError(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "coding", "SKILL.md"))
	wantErr := errors.New("stop")

	err := WalkSkillDirs(Options{Root: root}, func(Dir) error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error mismatch: got %v want %v", err, wantErr)
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
}

func mustRel(t *testing.T, root, path string) string {
	t.Helper()
	rel, err := filepath.Rel(root, path)
	if err != nil {
		t.Fatal(err)
	}
	return rel
}
