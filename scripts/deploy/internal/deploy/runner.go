package deploy

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Options struct {
	DryRun bool
}

type Runner struct {
	out io.Writer
}

func NewRunner(out io.Writer) Runner {
	return Runner{out: out}
}

func (r Runner) Run(configPath string, opts Options) error {
	absConfigPath, err := resolveConfigPath(configPath)
	if err != nil {
		return err
	}

	cfg, err := LoadConfig(absConfigPath)
	if err != nil {
		return err
	}

	configDir := filepath.Dir(absConfigPath)
	backupRoot := filepath.Join(configDir, ".deploy-backups", time.Now().Format("20060102-150405"))
	for i, item := range cfg.Items {
		src, err := resolveItemPath(configDir, item.Source)
		if err != nil {
			return fmt.Errorf("resolve source for items[%d]: %w", i, err)
		}
		dst, err := resolveItemPath(configDir, item.Destination)
		if err != nil {
			return fmt.Errorf("resolve destination for items[%d]: %w", i, err)
		}

		matcher, err := newExcludeMatcher(item.Exclude)
		if err != nil {
			return fmt.Errorf("items[%d]: %w", i, err)
		}

		if err := r.deployItem(i, src, dst, matcher, item.Replace, backupRoot, opts); err != nil {
			return err
		}
	}
	return nil
}

func (r Runner) deployItem(index int, src, dst string, matcher excludeMatcher, replace bool, backupRoot string, opts Options) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source for items[%d] %q: %w", index, src, err)
	}

	if info.IsDir() {
		return r.deployDir(index, src, dst, matcher, replace, backupRoot, opts)
	}
	if info.Mode().IsRegular() {
		if matcher.Match(filepath.Base(src)) {
			fmt.Fprintf(r.out, "SKIP     %s\n", src)
			return nil
		}
		return r.deployFile(index, src, dst, info.Mode(), replace, backupRoot, opts)
	}
	return fmt.Errorf("unsupported source for items[%d] %q: only regular files and directories are supported", index, src)
}

func (r Runner) deployDir(index int, src, dst string, matcher excludeMatcher, replace bool, backupRoot string, opts Options) error {
	if opts.DryRun {
		fmt.Fprintf(r.out, "DRY-RUN item[%d] dir  %s -> %s\n", index, src, dst)
	} else {
		fmt.Fprintf(r.out, "DEPLOY  item[%d] dir  %s -> %s\n", index, src, dst)
	}
	if err := r.backupDestination(dst, backupRoot, opts); err != nil {
		return err
	}
	if err := r.replaceDestination(dst, replace, opts); err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if matcher.Match(rel) {
			fmt.Fprintf(r.out, "SKIP     %s\n", path)
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		switch {
		case d.IsDir():
			return r.ensureDir(target, info.Mode(), opts)
		case info.Mode().IsRegular():
			return r.copyFile(path, target, info.Mode(), opts)
		default:
			fmt.Fprintf(r.out, "SKIP     %s\n", path)
			return nil
		}
	})
}

func (r Runner) deployFile(index int, src, dst string, mode os.FileMode, replace bool, backupRoot string, opts Options) error {
	if opts.DryRun {
		fmt.Fprintf(r.out, "DRY-RUN item[%d] file %s -> %s\n", index, src, dst)
	} else {
		fmt.Fprintf(r.out, "DEPLOY  item[%d] file %s -> %s\n", index, src, dst)
	}
	if err := r.backupDestination(dst, backupRoot, opts); err != nil {
		return err
	}
	if err := r.replaceDestination(dst, replace, opts); err != nil {
		return err
	}
	return r.copyFile(src, dst, mode, opts)
}

func (r Runner) ensureDir(path string, mode os.FileMode, opts Options) error {
	if opts.DryRun {
		fmt.Fprintf(r.out, "MKDIR    %s\n", path)
		return nil
	}
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdir %q: %w", path, err)
	}
	fmt.Fprintf(r.out, "MKDIR    %s\n", path)
	return nil
}

func (r Runner) copyFile(src, dst string, mode os.FileMode, opts Options) error {
	if opts.DryRun {
		fmt.Fprintf(r.out, "COPY     %s -> %s\n", src, dst)
		return nil
	}

	b, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %q: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir %q: %w", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, b, mode); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}

	fmt.Fprintf(r.out, "COPY     %s -> %s\n", src, dst)
	return nil
}

func (r Runner) backupDestination(dst, backupRoot string, opts Options) error {
	info, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat backup source %q: %w", dst, err)
	}

	backupPath := backupPathFor(backupRoot, dst)
	if opts.DryRun {
		fmt.Fprintf(r.out, "BACKUP   %s -> %s\n", dst, backupPath)
		return nil
	}

	if info.IsDir() {
		if err := copyDir(dst, backupPath); err != nil {
			return fmt.Errorf("backup dir %q to %q: %w", dst, backupPath, err)
		}
	} else if info.Mode().IsRegular() {
		if err := copyFile(dst, backupPath, info.Mode()); err != nil {
			return fmt.Errorf("backup file %q to %q: %w", dst, backupPath, err)
		}
	} else {
		fmt.Fprintf(r.out, "SKIP     %s\n", dst)
		return nil
	}

	fmt.Fprintf(r.out, "BACKUP   %s -> %s\n", dst, backupPath)
	return nil
}

func (r Runner) replaceDestination(dst string, replace bool, opts Options) error {
	if !replace {
		return nil
	}
	if opts.DryRun {
		fmt.Fprintf(r.out, "REMOVE   %s\n", dst)
		return nil
	}
	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("remove %q: %w", dst, err)
	}
	fmt.Fprintf(r.out, "REMOVE   %s\n", dst)
	return nil
}

func backupPathFor(backupRoot, dst string) string {
	volume := filepath.VolumeName(dst)
	withoutVolume := strings.TrimPrefix(dst, volume)
	withoutVolume = strings.TrimPrefix(filepath.Clean(withoutVolume), string(filepath.Separator))
	if volume != "" {
		volume = strings.TrimSuffix(volume, ":")
		return filepath.Join(backupRoot, volume, withoutVolume)
	}
	return filepath.Join(backupRoot, withoutVolume)
}

func copyFile(src, dst string, mode os.FileMode) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %q: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir %q: %w", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, b, mode); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}
	return nil
}

func copyDir(src, dst string) error {
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
		default:
			return nil
		}
	})
}
