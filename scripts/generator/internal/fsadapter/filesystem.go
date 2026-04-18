package fsadapter

import (
	iofs "io/fs"
	"os"
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
