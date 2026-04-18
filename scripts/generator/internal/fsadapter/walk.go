package fsadapter

import (
	"bytes"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"strings"
)

// OutputPathFn は出力パスを解決する関数。
// outRoot、relPath、ファイル内容を受け取り、出力先パスを返す。
type OutputPathFn func(outRoot, relPath string, content []byte) (string, error)

// WalkResources は inRoot 配下の skills ディレクトリを走査し、
// 各マークダウンファイルに対して outputPathFn で出力パスを解決して書き出す。
func WalkResources(fsys FileSystem, inRoot, outRoot string, outputPathFn OutputPathFn) error {
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

		if !IsMarkdown(d.Name()) {
			return nil
		}

		relPath, err := filepath.Rel(inRoot, path)
		if err != nil {
			return fmt.Errorf("rel %q from %q: %w", path, inRoot, err)
		}

		if IsExcluded(relPath) {
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

func ensureTrailingNewline(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	if len(b) == 0 || b[len(b)-1] == '\n' {
		return []byte(b)
	}
	return append(b, '\n')
}
