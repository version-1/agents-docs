package skillscan

import (
	"fmt"
	"os"
	"path/filepath"
)

type Matcher interface {
	Match(rel string) bool
}

type Dir struct {
	Name      string
	Path      string
	SkillFile string
}

type SkipReason string

const (
	SkipExcluded        SkipReason = "excluded"
	SkipNonRegularSkill SkipReason = "non-regular-skill"
)

type Options struct {
	Root    string
	Matcher Matcher
	// OnSkip is called for excluded directories, excluded SKILL.md files, and non-regular SKILL.md files.
	// The path is absolute or relative according to Root and the walked directory tree.
	OnSkip func(path string, reason SkipReason)
}

// WalkSkillDirs visits directories whose direct child SKILL.md is a regular file.
// Matcher is evaluated against root-relative directory paths and SKILL.md paths.
// Excluded directories are not descended into.
func WalkSkillDirs(opts Options, visit func(Dir) error) error {
	return opts.walkDir(opts.Root, visit)
}

func (o Options) walkDir(dir string, visit func(Dir) error) error {
	rel, err := filepath.Rel(o.Root, dir)
	if err != nil {
		return err
	}
	if o.matches(rel) {
		o.skip(dir, SkipExcluded)
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %q: %w", dir, err)
	}

	if err := o.visitSkillDir(dir, entries, visit); err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if err := o.walkDir(filepath.Join(dir, entry.Name()), visit); err != nil {
			return err
		}
	}
	return nil
}

func (o Options) visitSkillDir(dir string, entries []os.DirEntry, visit func(Dir) error) error {
	for _, entry := range entries {
		if entry.Name() != "SKILL.md" {
			continue
		}
		skillFile := filepath.Join(dir, entry.Name())
		skillRel, err := filepath.Rel(o.Root, skillFile)
		if err != nil {
			return err
		}
		if o.matches(skillRel) {
			o.skip(skillFile, SkipExcluded)
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("stat skill file %q: %w", skillFile, err)
		}
		if !info.Mode().IsRegular() {
			o.skip(skillFile, SkipNonRegularSkill)
			return nil
		}
		if err := visit(Dir{Name: filepath.Base(dir), Path: dir, SkillFile: skillFile}); err != nil {
			return fmt.Errorf("visit skill dir %q: %w", dir, err)
		}
		return nil
	}
	return nil
}

func (o Options) matches(rel string) bool {
	return o.Matcher != nil && o.Matcher.Match(rel)
}

func (o Options) skip(path string, reason SkipReason) {
	if o.OnSkip != nil {
		o.OnSkip(path, reason)
	}
}
