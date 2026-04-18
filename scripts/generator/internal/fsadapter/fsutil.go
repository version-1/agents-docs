package fsadapter

import (
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
)

func DirExists(fsys FileSystem, path string) (bool, error) {
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

func RemoveDirContents(fsys FileSystem, dir string) error {
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

func CopyFile(fsys FileSystem, src, dst string, mod os.FileMode) error {
	b, err := fsys.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %q: %w", src, err)
	}

	if err := fsys.WriteFile(dst, b, mod); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}
	return nil
}

func CopyDir(fsys FileSystem, src, dst string) error {
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
			return CopyFile(fsys, path, target, info.Mode())

		case info.Mode()&os.ModeSymlink != 0:
			return nil

		default:
			return nil
		}
	})
}
