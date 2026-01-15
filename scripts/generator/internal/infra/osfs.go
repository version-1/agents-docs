package infra

import (
	"io/fs"
	"os"
	"path/filepath"
)

type OSFS struct{}

func (OSFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSFS) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (OSFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (OSFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (OSFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OSFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (OSFS) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}
