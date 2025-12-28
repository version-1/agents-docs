package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		inRoot  = flag.String("input", "", "input root directory to scan (e.g., ./.)")
		outRoot = flag.String("output", "out/.codex/skills", "output directory (e.g., .out/codex/skills)")
	)
	flag.Parse()

	exists, err := dirExists(*outRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error checking output dir:", err)
		os.Exit(1)
	}

	if exists {
		if err := removeDirContents(*outRoot); err != nil {
			fmt.Fprintln(os.Stderr, "error clearing output dir:", err)
			os.Exit(1)
		}
	}

	if err := createDirIfNotExists(*outRoot); err != nil {
		fmt.Fprintln(os.Stderr, "error creating output dir:", err)
		os.Exit(1)
	}

	if err := createDirIfNotExists(*inRoot); err != nil {
		fmt.Fprintln(os.Stderr, "error creating input dir:", err)
		os.Exit(1)
	}

	if err := copyDocs(*inRoot, *outRoot); err != nil {
		fmt.Fprintln(os.Stderr, "error copying docs:", err)
		os.Exit(1)
	}

	if err := generateSkills(*inRoot, *outRoot); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

}

func removeDirContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func createDirIfNotExists(path string) error {
	exists, err := dirExists(path)
	if err != nil {
		return err
	}
	if !exists {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("mkdir %q: %w", path, err)
		}
	}

	return nil
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func copyDocs(inRoot, outRoot string) error {
	// Implementation omitted for brevity; focus is on generateSkills
	agentMdPath := filepath.Join(inRoot, "Agents.md")
	if err := copyFile(agentMdPath, filepath.Join(outRoot, "Agents.md"), 0644); err != nil {
		return err
	}

	agentsDir := "agents"
	if err := copyDir(filepath.Join(inRoot, agentsDir), filepath.Join(outRoot, agentsDir)); err != nil {
		return err
	}

	return nil
}

func copyDir(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		switch {
		case d.IsDir():
			return os.MkdirAll(target, info.Mode())

		case info.Mode().IsRegular():
			return copyFile(path, target, info.Mode())

		case info.Mode()&os.ModeSymlink != 0:
			// 方針: シンボリックリンクは無視
			return nil

		default:
			return nil
		}
	})
}

func copyFile(src, dst string, mod os.FileMode) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %q: %w", src, err)
	}

	if err := os.WriteFile(dst, b, mod); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}
	return nil
}

func generateSkills(inRoot, outRoot string) error {
	// Ensure input exists
	info, err := os.Stat(inRoot)
	if err != nil {
		return fmt.Errorf("stat input root %q: %w", inRoot, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("input root %q is not a directory", inRoot)
	}

	var found int

	walkFn := func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			// skip hidden dirs like .git, .codex, etc.
			if strings.HasPrefix(d.Name(), ".") {
				if path == inRoot {
					return nil // allow root even if it's "."
				}
				return fs.SkipDir
			}

			if !strings.Contains(path, "skills") && !strings.Contains(path, "docs/ja") {
				fmt.Println("skipping dir:", path)
				return fs.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		found++

		fmt.Println("Processing:", path)
		excludePaths := []string{
			"skills/skills.md",
			"Agents.md",
		}
		for _, exclude := range excludePaths {
			excludeFullPath := filepath.Join(inRoot, exclude)
			if path == excludeFullPath {
				return nil
			}
		}

		// Base name without extension, uppercased (SKILL.md style)
		outRootPath := strings.Replace(filepath.Dir(path), "docs/ja", "", 1)
		fmt.Println("Path:", filepath.Dir(path), outRoot, outRootPath)
		outName := filepath.Join(outRootPath, "SKILL.md")
		outPath := filepath.Join(outRoot, outName)

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return fmt.Errorf("mkdir %q: %w", filepath.Dir(outPath), err)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %q: %w", path, err)
		}

		// Optional: ensure trailing newline (handy for markdown)
		b = ensureTrailingNewline(b)

		if err := os.WriteFile(outPath, b, 0o644); err != nil {
			return fmt.Errorf("write %q: %w", outPath, err)
		}

		fmt.Printf("wrote %s\n", outPath)
		return nil
	}

	if err := filepath.WalkDir(inRoot, walkFn); err != nil {
		return err
	}

	if found == 0 {
		fmt.Println("no markdown files found")
		return nil
	}

	fmt.Printf("done: %d file(s)\n", found)
	return nil
}

func ensureTrailingNewline(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n")) // normalize CRLF->LF (optional; remove if you want exact copy)
	if len(b) == 0 || b[len(b)-1] == '\n' {
		return b
	}
	return append(b, '\n')
}
