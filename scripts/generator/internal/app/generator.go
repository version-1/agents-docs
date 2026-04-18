package app

import (
	"bytes"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"

	"generator/internal/domain"
)

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadDir(name string) ([]os.DirEntry, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error
	WalkDir(root string, fn iofs.WalkDirFunc) error
}

type SkillGenerator interface {
	Generate(fsys FileSystem, inRoot, outRoot string) error
}

type Generator struct {
	fs     FileSystem
	skills SkillGenerator
}

func NewGenerator(fs FileSystem, skills SkillGenerator) *Generator {
	return &Generator{fs: fs, skills: skills}
}

func (g *Generator) Run(inRoot, outRoot string) error {
	inExists, err := dirExists(g.fs, inRoot)
	if err != nil {
		return fmt.Errorf("check input dir: %w", err)
	}
	if !inExists {
		return fmt.Errorf("input dir does not exist: %s", inRoot)
	}

	outExists, err := dirExists(g.fs, outRoot)
	if err != nil {
		return fmt.Errorf("check output dir: %w", err)
	}
	if outExists {
		if err := removeDirContents(g.fs, outRoot); err != nil {
			return fmt.Errorf("clear output dir: %w", err)
		}
	}

	if err := g.fs.MkdirAll(outRoot, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := copyDocs(g.fs, inRoot, outRoot); err != nil {
		return fmt.Errorf("copy docs: %w", err)
	}

	if err := g.skills.Generate(g.fs, inRoot, outRoot); err != nil {
		return fmt.Errorf("generate skills: %w", err)
	}

	return nil
}

func removeDirContents(fsys FileSystem, dir string) error {
	entries, err := fsys.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if err := fsys.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func dirExists(fsys FileSystem, path string) (bool, error) {
	info, err := fsys.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return false, fmt.Errorf("%s is not a directory", path)
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func copyDocs(fsys FileSystem, inRoot, outRoot string) error {
	agentMdPath := filepath.Join(inRoot, "Agents.md")
	if err := copyFile(fsys, agentMdPath, filepath.Join(outRoot, "Agents.md"), 0o644); err != nil {
		return err
	}

	agentsDir := "agents"
	if err := copyDir(fsys, filepath.Join(inRoot, agentsDir), filepath.Join(outRoot, agentsDir)); err != nil {
		return err
	}

	return nil
}

func copyDir(fsys FileSystem, src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	info, err := fsys.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	return fsys.WalkDir(src, func(path string, d iofs.DirEntry, err error) error {
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
			return fsys.MkdirAll(target, info.Mode())

		case info.Mode().IsRegular():
			return copyFile(fsys, path, target, info.Mode())

		case info.Mode()&os.ModeSymlink != 0:
			return nil

		default:
			return nil
		}
	})
}

func copyFile(fsys FileSystem, src, dst string, mod os.FileMode) error {
	b, err := fsys.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %q: %w", src, err)
	}

	if err := fsys.WriteFile(dst, b, mod); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}
	return nil
}

func ensureTrailingNewline(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	if len(b) == 0 || b[len(b)-1] == '\n' {
		return []byte(b)
	}
	return append(b, '\n')
}

// walkSkills は skills ディレクトリを走査し、各スキルファイルに対して outputPathFn で
// 出力パスを解決して書き出す共通ロジック。
// outputPathFn は outRoot、relPath、ファイル内容を受け取り、出力パスを返す。
func walkSkills(fsys FileSystem, inRoot, outRoot string, outputPathFn func(outRoot, relPath string, content []byte) (string, error)) error {
	var found int

	walkFn := func(path string, d iofs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				if path == inRoot {
					return nil
				}
				return iofs.SkipDir
			}

			if !strings.Contains(path, "skills") && path != inRoot {
				return iofs.SkipDir
			}
			return nil
		}

		if !domain.IsMarkdown(d.Name()) {
			return nil
		}

		relPath, err := filepath.Rel(inRoot, path)
		if err != nil {
			return fmt.Errorf("rel %q from %q: %w", path, inRoot, err)
		}

		if domain.IsExcluded(relPath) {
			return nil
		}

		b, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %q: %w", path, err)
		}

		found++
		outPath, err := outputPathFn(outRoot, relPath, b)
		if err != nil {
			return fmt.Errorf("resolve output path for %q: %w", relPath, err)
		}
		outDir := filepath.Dir(outPath)

		if err := fsys.MkdirAll(outDir, 0o755); err != nil {
			return fmt.Errorf("mkdir %q: %w", outDir, err)
		}

		b = ensureTrailingNewline(b)

		if err := fsys.WriteFile(outPath, b, 0o644); err != nil {
			return fmt.Errorf("write %q: %w", outPath, err)
		}

		return nil
	}

	if err := fsys.WalkDir(inRoot, walkFn); err != nil {
		return err
	}

	if found == 0 {
		fmt.Println("no markdown files found")
		return nil
	}

	fmt.Printf("done: %d file(s)\n", found)
	return nil
}
