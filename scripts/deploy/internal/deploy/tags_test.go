package deploy

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunnerTagSelection(t *testing.T) {
	for _, tag := range []string{"", "agents", "skills", "unknown"} {
		t.Run(tag, func(t *testing.T) {
			root := t.TempDir()
			for _, name := range []string{"agents", "skills", "untagged"} {
				if err := os.WriteFile(filepath.Join(root, name), []byte(name), 0600); err != nil {
					t.Fatal(err)
				}
			}
			path := filepath.Join(root, "deploy.json")
			writeConfig(t, path, `{"items":[
				{"source":"agents","destination":"out/agents","tags":["agents","other"]},
				{"source":"skills","destination":"out/skills","tags":["skills"]},
				{"source":"untagged","destination":"out/untagged"}
			]}`)
			opts := Options{Tag: tag}
			if tag == "agents" {
				// A non-skills tag must not even read the external config.
				opts.ExternalSkillsPath = filepath.Join(root, "missing-external.json")
			}
			err := runFromDir(t, root, func() error { return NewRunner(io.Discard).Run(path, opts) })
			if tag == "unknown" {
				if err == nil || !strings.Contains(err.Error(), "no deployment items match") {
					t.Fatalf("expected unknown tag error, got %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"agents", "skills", "untagged"} {
				_, err := os.Stat(filepath.Join(root, "out", name))
				want := tag == "" || tag == name
				if want && err != nil {
					t.Fatalf("expected %s to be copied: %v", name, err)
				}
				if !want && !os.IsNotExist(err) {
					t.Fatalf("unexpected destination %s: %v", name, err)
				}
			}
		})
	}
}
