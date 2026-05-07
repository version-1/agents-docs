package fileops

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Options struct {
	DryRun bool
}

type Report struct {
	CopiedFiles int
	CreatedDirs int
	Skipped     int
}

type Matcher interface {
	Match(rel string) bool
}

type TreeOptions struct {
	CopyRoot        string
	ExcludeRoot     string
	DestinationRoot string
	Matcher         Matcher
	Report          *Report
	Options         Options
}

func EnsureDir(path string, mode os.FileMode, opts Options) error {
	if opts.DryRun {
		return nil
	}
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdir %q: %w", path, err)
	}
	return nil
}

func CopyFile(src, dst string, mode os.FileMode, opts Options) error {
	if opts.DryRun {
		return nil
	}
	return CopyFileWithParents(src, dst, mode)
}

func CopyFileWithParents(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir %q: %w", filepath.Dir(dst), err)
	}
	return CopyFileWithoutMkdir(src, dst, mode)
}

func CopyFileWithoutMkdir(src, dst string, mode os.FileMode) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("read %q: %w", src, err)
	}
	defer func() {
		if closeErr := in.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("read %q: %w", src, closeErr)
		}
	}()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}
	defer func() {
		if closeErr := out.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("write %q: %w", dst, closeErr)
		}
	}()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}
	return nil
}

func CopyDir(src, dst string) error {
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
			// WalkDir visits parent directories before their files, so the parent exists here.
			return CopyFileWithoutMkdir(path, target, info.Mode())
		default:
			return nil
		}
	})
}

func CopyTree(opts TreeOptions) error {
	return filepath.WalkDir(opts.CopyRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		skip, err := opts.shouldSkip(path, d.IsDir())
		if skip || err != nil {
			return err
		}

		return copyTreeEntry(path, d, opts)
	})
}

func (o TreeOptions) shouldSkip(path string, isDir bool) (bool, error) {
	excludeRel, err := filepath.Rel(o.ExcludeRoot, path)
	if err != nil {
		return false, err
	}
	if o.Matcher == nil || !o.Matcher.Match(excludeRel) {
		return false, nil
	}
	o.Report.Skipped++
	if isDir {
		return true, filepath.SkipDir
	}
	return true, nil
}

func copyTreeEntry(path string, d fs.DirEntry, opts TreeOptions) error {
	destinationRel, err := opts.destinationRel(path)
	if err != nil {
		return err
	}
	target := filepath.Join(opts.DestinationRoot, destinationRel)

	if opts.Options.DryRun {
		countDryRunEntry(d, opts.Report)
		return nil
	}

	info, err := d.Info()
	if err != nil {
		return err
	}

	switch {
	case d.IsDir():
		if err := EnsureDir(target, info.Mode(), opts.Options); err != nil {
			return err
		}
		opts.Report.CreatedDirs++
	case info.Mode().IsRegular():
		if err := CopyFileWithoutMkdir(path, target, info.Mode()); err != nil {
			return err
		}
		opts.Report.CopiedFiles++
	default:
		opts.Report.Skipped++
	}
	return nil
}

func countDryRunEntry(d fs.DirEntry, report *Report) {
	switch {
	case d.IsDir():
		report.CreatedDirs++
	case d.Type().IsRegular():
		report.CopiedFiles++
	default:
		report.Skipped++
	}
}

func (o TreeOptions) destinationRel(path string) (string, error) {
	if o.CopyRoot == o.ExcludeRoot {
		return filepath.Rel(o.ExcludeRoot, path)
	}
	return filepath.Rel(o.CopyRoot, path)
}
