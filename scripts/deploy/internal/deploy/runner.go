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
	DryRun  bool
	NoColor bool
}

type Runner struct {
	out io.Writer
}

type itemReport struct {
	copiedFiles int
	createdDirs int
	skipped     int
}

type flattenedSkillDir struct {
	source string
	target string
}

const (
	colorReset   = "\033[0m"
	colorFaint   = "\033[2m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

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

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current directory: %w", err)
	}

	configDir := filepath.Dir(absConfigPath)
	backupRoot := filepath.Join(configDir, ".deploy-backups", time.Now().Format("20060102-150405"))
	for i, item := range cfg.Items {
		src, err := resolveSourcePath(cwd, item.Source)
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

		if err := r.deployItem(i, src, dst, matcher, item.Replace, item.Flatten, backupRoot, opts); err != nil {
			return err
		}
	}
	return nil
}

func (r Runner) deployItem(index int, src, dst string, matcher excludeMatcher, replace, flatten bool, backupRoot string, opts Options) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source for items[%d] %q: %w", index, src, err)
	}

	report := &itemReport{}
	if info.IsDir() {
		if flatten {
			return r.deployFlattenedSkillDirs(index, src, dst, matcher, replace, backupRoot, report, opts)
		}
		return r.deployDir(index, src, dst, matcher, replace, backupRoot, report, opts)
	}
	if flatten {
		return fmt.Errorf("items[%d]: flatten requires a directory source", index)
	}
	if info.Mode().IsRegular() {
		if matcher.Match(filepath.Base(src)) {
			r.printItemHeader(index, "file", src, dst, opts)
			report.skipped++
			r.printSummary(report, opts)
			return nil
		}
		return r.deployFile(index, src, dst, info.Mode(), replace, backupRoot, report, opts)
	}
	return fmt.Errorf("unsupported source for items[%d] %q: only regular files and directories are supported", index, src)
}

func (r Runner) deployFlattenedSkillDirs(index int, src, dst string, matcher excludeMatcher, replace bool, backupRoot string, report *itemReport, opts Options) error {
	r.printItemHeader(index, "flattened-skill-dirs", src, dst, opts)
	skillDirs, err := findFlattenedSkillDirs(src, dst, matcher, report)
	if err != nil {
		return err
	}

	if err := r.backupDestination(dst, backupRoot, opts); err != nil {
		return err
	}
	if err := r.replaceDestination(dst, replace, opts); err != nil {
		return err
	}

	for _, skillDir := range skillDirs {
		if err := r.copySkillDir(skillDir.source, src, skillDir.target, matcher, report, opts); err != nil {
			return err
		}
	}
	r.printSummary(report, opts)
	return nil
}

func findFlattenedSkillDirs(src, dst string, matcher excludeMatcher, report *itemReport) ([]flattenedSkillDir, error) {
	targets := map[string]string{}
	var skillDirs []flattenedSkillDir
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if matcher.Match(rel) {
			report.skipped++
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			return nil
		}

		skillFile := filepath.Join(path, "SKILL.md")
		skillRel, err := filepath.Rel(src, skillFile)
		if err != nil {
			return err
		}
		if matcher.Match(skillRel) {
			report.skipped++
			return nil
		}
		info, err := os.Stat(skillFile)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return fmt.Errorf("stat skill file %q: %w", skillFile, err)
		}
		if !info.Mode().IsRegular() {
			report.skipped++
			return nil
		}

		target := filepath.Join(dst, filepath.Base(path))
		if existing, ok := targets[target]; ok {
			return fmt.Errorf("flatten target conflict %q from %q and %q", target, existing, path)
		}
		targets[target] = path
		skillDirs = append(skillDirs, flattenedSkillDir{source: path, target: target})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return skillDirs, nil
}

func (r Runner) copySkillDir(skillDir, sourceRoot, targetRoot string, matcher excludeMatcher, report *itemReport, opts Options) error {
	return filepath.WalkDir(skillDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		sourceRel, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		if matcher.Match(sourceRel) {
			report.skipped++
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		skillRel, err := filepath.Rel(skillDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(targetRoot, skillRel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		switch {
		case d.IsDir():
			if err := r.ensureDir(target, info.Mode(), opts); err != nil {
				return err
			}
			report.createdDirs++
			return nil
		case info.Mode().IsRegular():
			if err := r.copyFile(path, target, info.Mode(), opts); err != nil {
				return err
			}
			report.copiedFiles++
			return nil
		default:
			report.skipped++
			return nil
		}
	})
}

func (r Runner) deployDir(index int, src, dst string, matcher excludeMatcher, replace bool, backupRoot string, report *itemReport, opts Options) error {
	r.printItemHeader(index, "dir", src, dst, opts)
	if err := r.backupDestination(dst, backupRoot, opts); err != nil {
		return err
	}
	if err := r.replaceDestination(dst, replace, opts); err != nil {
		return err
	}

	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if matcher.Match(rel) {
			report.skipped++
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
			if err := r.ensureDir(target, info.Mode(), opts); err != nil {
				return err
			}
			report.createdDirs++
			return nil
		case info.Mode().IsRegular():
			if err := r.copyFile(path, target, info.Mode(), opts); err != nil {
				return err
			}
			report.copiedFiles++
			return nil
		default:
			report.skipped++
			return nil
		}
	})
	if err != nil {
		return err
	}
	r.printSummary(report, opts)
	return nil
}

func (r Runner) deployFile(index int, src, dst string, mode os.FileMode, replace bool, backupRoot string, report *itemReport, opts Options) error {
	r.printItemHeader(index, "file", src, dst, opts)
	if err := r.backupDestination(dst, backupRoot, opts); err != nil {
		return err
	}
	if err := r.replaceDestination(dst, replace, opts); err != nil {
		return err
	}
	if err := r.copyFile(src, dst, mode, opts); err != nil {
		return err
	}
	report.copiedFiles++
	r.printSummary(report, opts)
	return nil
}

func (r Runner) ensureDir(path string, mode os.FileMode, opts Options) error {
	if opts.DryRun {
		return nil
	}
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdir %q: %w", path, err)
	}
	return nil
}

func (r Runner) copyFile(src, dst string, mode os.FileMode, opts Options) error {
	if opts.DryRun {
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
		r.printBackup(backupPath, opts)
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
		return nil
	}

	r.printBackup(backupPath, opts)
	return nil
}

func (r Runner) replaceDestination(dst string, replace bool, opts Options) error {
	if !replace {
		return nil
	}
	if opts.DryRun {
		r.printReplace("remove existing destination", opts)
		return nil
	}
	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("remove %q: %w", dst, err)
	}
	r.printReplace("removed existing destination", opts)
	return nil
}

func (r Runner) printItemHeader(index int, kind, src, dst string, opts Options) {
	mode := "DEPLOY"
	modeColor := colorGreen
	if opts.DryRun {
		mode = "DRY-RUN"
		modeColor = colorYellow
	}
	fmt.Fprintf(r.out, "\n%s item[%d] %s\n", colorize("["+mode+"]", modeColor, opts), index, colorize(kind, colorCyan, opts))
	fmt.Fprintf(r.out, "  %s      %s\n", colorize("source:", colorFaint, opts), src)
	fmt.Fprintf(r.out, "  %s %s\n", colorize("destination:", colorFaint, opts), dst)
}

func (r Runner) printBackup(path string, opts Options) {
	fmt.Fprintf(r.out, "  %s %s\n", colorize("backup:", colorBlue, opts), path)
}

func (r Runner) printReplace(message string, opts Options) {
	fmt.Fprintf(r.out, "  %s %s\n", colorize("replace:", colorMagenta, opts), message)
}

func (r Runner) printSummary(report *itemReport, opts Options) {
	fmt.Fprintf(
		r.out,
		"  %s %s, %s, %s\n",
		colorize("summary:", colorFaint, opts),
		colorize(fmt.Sprintf("%d copied", report.copiedFiles), colorGreen, opts),
		colorize(fmt.Sprintf("%d dirs", report.createdDirs), colorCyan, opts),
		colorize(fmt.Sprintf("%d skipped", report.skipped), colorYellow, opts),
	)
}

func colorize(s, color string, opts Options) string {
	if opts.NoColor {
		return s
	}
	return color + s + colorReset
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
