package deploy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"deploy/internal/config"
	"deploy/internal/external"
	"deploy/internal/fileops"
	"deploy/internal/matcher"
	"deploy/internal/pathutil"
	"deploy/internal/skillscan"
)

type Options struct {
	DryRun             bool
	NoColor            bool
	ExternalSkillsPath string
}

type Runner struct {
	out     io.Writer
	fetcher external.Fetcher
}

type itemReport = fileops.Report

type flattenedSkillDir struct {
	source string
	target string
}

type fetchedExternalSkill struct {
	skill  external.Skill
	source string
}

type headerField struct {
	label string
	value string
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
	return Runner{out: out, fetcher: external.GitFetcher{}}
}

func newRunnerWithFetcher(out io.Writer, fetcher external.Fetcher) Runner {
	return Runner{out: out, fetcher: fetcher}
}

func (r Runner) Run(configPath string, opts Options) error {
	absConfigPath, err := pathutil.ResolveConfigPath(configPath)
	if err != nil {
		return err
	}

	cfg, err := config.Load(absConfigPath)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current directory: %w", err)
	}

	configDir := filepath.Dir(absConfigPath)
	backupRoot := filepath.Join(configDir, ".deploy-backups", time.Now().Format("20060102-150405"))

	externalSkills, externalConfigDir, cleanupExternalSkills, err := r.prepareOptionalExternalSkills(opts.ExternalSkillsPath, cfg, cwd)
	if err != nil {
		return err
	}
	if cleanupExternalSkills != nil {
		defer cleanupExternalSkills()
	}

	if err := r.deployConfiguredItems(cfg, cwd, configDir, backupRoot, opts); err != nil {
		return err
	}
	return r.deployExternalSkills(externalSkills, externalConfigDir, len(cfg.Items), backupRoot, opts)
}

func (r Runner) prepareOptionalExternalSkills(path string, cfg config.Config, cwd string) ([]fetchedExternalSkill, string, func(), error) {
	var cleanup func()
	var skills []fetchedExternalSkill
	var configDir string
	var err error
	if path == "" {
		return nil, "", nil, nil
	}
	skills, configDir, cleanup, err = r.prepareExternalSkills(path, cfg, cwd)
	if err != nil {
		return nil, "", nil, err
	}
	return skills, configDir, cleanup, nil
}

func (r Runner) deployConfiguredItems(cfg config.Config, cwd, configDir, backupRoot string, opts Options) error {
	for i, item := range cfg.Items {
		src, dst, excludeMatcher, err := resolveDeployItem(i, item, cwd, configDir)
		if err != nil {
			return err
		}
		if err := r.deployItem(i, src, dst, excludeMatcher, item.Replace, item.Flatten, backupRoot, opts); err != nil {
			return err
		}
	}
	return nil
}

func resolveDeployItem(index int, item config.Item, cwd, configDir string) (string, string, matcher.Matcher, error) {
	src, err := pathutil.ResolveSourcePath(cwd, item.Source)
	if err != nil {
		return "", "", matcher.Matcher{}, fmt.Errorf("resolve source for items[%d]: %w", index, err)
	}
	dst, err := pathutil.ResolveItemPath(configDir, item.Destination)
	if err != nil {
		return "", "", matcher.Matcher{}, fmt.Errorf("resolve destination for items[%d]: %w", index, err)
	}

	excludeMatcher, err := matcher.New(item.Exclude)
	if err != nil {
		return "", "", matcher.Matcher{}, fmt.Errorf("items[%d]: %w", index, err)
	}
	return src, dst, excludeMatcher, nil
}

func (r Runner) deployExternalSkills(skills []fetchedExternalSkill, externalConfigDir string, startIndex int, backupRoot string, opts Options) error {
	if len(skills) == 0 {
		return nil
	}
	index := startIndex
	for i, fetched := range skills {
		for j, destination := range fetched.skill.Destination {
			dst, err := pathutil.ResolveItemPath(externalConfigDir, destination)
			if err != nil {
				return fmt.Errorf("resolve destination for externalSkills[%d].destination[%d]: %w", i, j, err)
			}
			report := &itemReport{}
			if err := r.deployExternalSkill(index, fetched.skill, fetched.source, dst, backupRoot, report, opts); err != nil {
				return err
			}
			index++
		}
	}
	return nil
}

func (r Runner) prepareExternalSkills(path string, cfg config.Config, cwd string) ([]fetchedExternalSkill, string, func(), error) {
	absExternalPath, err := pathutil.ResolveConfigPath(path)
	if err != nil {
		return nil, "", nil, err
	}
	skills, err := external.Load(absExternalPath)
	if err != nil {
		return nil, "", nil, err
	}
	if err := external.ValidateConflicts(skills, cfg, cwd); err != nil {
		return nil, "", nil, err
	}

	workDir, err := os.MkdirTemp("", "deploy-external-skills-*")
	if err != nil {
		return nil, "", nil, fmt.Errorf("create external skills temp dir: %w", err)
	}
	cleanupOnError := true
	defer func() {
		if cleanupOnError {
			_ = os.RemoveAll(workDir)
		}
	}()

	fetched := make([]fetchedExternalSkill, 0, len(skills))
	for _, skill := range skills {
		src, err := r.fetcher.Fetch(skill, workDir)
		if err != nil {
			return nil, "", nil, fmt.Errorf("fetch external skill %q from %s: %w", skill.Name, skill.URL, err)
		}
		if err := external.ValidateFetched(skill, src); err != nil {
			return nil, "", nil, err
		}
		fetched = append(fetched, fetchedExternalSkill{skill: skill, source: src})
	}

	cleanupOnError = false
	return fetched, filepath.Dir(absExternalPath), func() { _ = os.RemoveAll(workDir) }, nil
}

func (r Runner) deployExternalSkill(index int, skill external.Skill, src, dst, backupRoot string, report *itemReport, opts Options) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat external skill %q source %q: %w", skill.Name, src, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("external skill %q source %q must be a directory", skill.Name, src)
	}

	r.printExternalSkillHeader(index, skill, src, dst, opts)
	if err := r.backupDestination(dst, backupRoot, opts); err != nil {
		return err
	}
	if err := r.replaceDestination(dst, true, opts); err != nil {
		return err
	}
	if err := r.copySkillDir(src, src, dst, matcher.Matcher{}, report, opts); err != nil {
		return err
	}
	r.printSummary(report, opts)
	return nil
}

func (r Runner) deployItem(index int, src, dst string, excludeMatcher matcher.Matcher, replace, flatten bool, backupRoot string, opts Options) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source for items[%d] %q: %w", index, src, err)
	}

	report := &itemReport{}
	if info.IsDir() {
		if flatten {
			return r.deployFlattenedSkillDirs(index, src, dst, excludeMatcher, replace, backupRoot, report, opts)
		}
		return r.deployDir(index, src, dst, excludeMatcher, replace, backupRoot, report, opts)
	}
	if flatten {
		return fmt.Errorf("items[%d]: flatten requires a directory source", index)
	}
	if info.Mode().IsRegular() {
		if excludeMatcher.Match(filepath.Base(src)) {
			r.printItemHeader(index, "file", src, dst, opts)
			report.Skipped++
			r.printSummary(report, opts)
			return nil
		}
		return r.deployFile(index, src, dst, info.Mode(), replace, backupRoot, report, opts)
	}
	return fmt.Errorf("unsupported source for items[%d] %q: only regular files and directories are supported", index, src)
}

func (r Runner) deployFlattenedSkillDirs(index int, src, dst string, excludeMatcher matcher.Matcher, replace bool, backupRoot string, report *itemReport, opts Options) error {
	r.printItemHeader(index, "flattened-skill-dirs", src, dst, opts)
	skillDirs, err := findFlattenedSkillDirs(src, dst, excludeMatcher, report)
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
		if err := r.copySkillDir(skillDir.source, src, skillDir.target, excludeMatcher, report, opts); err != nil {
			return err
		}
	}
	r.printSummary(report, opts)
	return nil
}

func findFlattenedSkillDirs(src, dst string, excludeMatcher matcher.Matcher, report *itemReport) ([]flattenedSkillDir, error) {
	targets := map[string]string{}
	var skillDirs []flattenedSkillDir
	err := skillscan.WalkSkillDirs(skillscan.Options{
		Root:    src,
		Matcher: excludeMatcher,
		OnSkip:  func(string, skillscan.SkipReason) { report.Skipped++ },
	}, func(dir skillscan.Dir) error {
		target := filepath.Join(dst, dir.Name)
		if existing, ok := targets[target]; ok {
			return fmt.Errorf("flatten target conflict %q from %q and %q", target, existing, dir.Path)
		}
		targets[target] = dir.Path
		skillDirs = append(skillDirs, flattenedSkillDir{source: dir.Path, target: target})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return skillDirs, nil
}

func (r Runner) copySkillDir(skillDir, excludeRoot, destinationRoot string, excludeMatcher matcher.Matcher, report *itemReport, opts Options) error {
	return fileops.CopyTree(fileops.TreeOptions{
		CopyRoot:        skillDir,
		ExcludeRoot:     excludeRoot,
		DestinationRoot: destinationRoot,
		Matcher:         excludeMatcher,
		Report:          report,
		Options:         fileOptions(opts),
	})
}

func (r Runner) deployDir(index int, src, dst string, excludeMatcher matcher.Matcher, replace bool, backupRoot string, report *itemReport, opts Options) error {
	r.printItemHeader(index, "dir", src, dst, opts)
	if err := r.backupDestination(dst, backupRoot, opts); err != nil {
		return err
	}
	if err := r.replaceDestination(dst, replace, opts); err != nil {
		return err
	}

	if err := fileops.CopyTree(fileops.TreeOptions{
		CopyRoot:        src,
		ExcludeRoot:     src,
		DestinationRoot: dst,
		Matcher:         excludeMatcher,
		Report:          report,
		Options:         fileOptions(opts),
	}); err != nil {
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
	if err := fileops.CopyFile(src, dst, mode, fileOptions(opts)); err != nil {
		return err
	}
	report.CopiedFiles++
	r.printSummary(report, opts)
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
		if err := fileops.CopyDir(dst, backupPath); err != nil {
			return fmt.Errorf("backup dir %q to %q: %w", dst, backupPath, err)
		}
	} else if info.Mode().IsRegular() {
		if err := fileops.CopyFileWithParents(dst, backupPath, info.Mode()); err != nil {
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
	r.printHeader(index, kind, []headerField{
		{label: "source", value: src},
		{label: "destination", value: dst},
	}, opts)
}

func (r Runner) printExternalSkillHeader(index int, skill external.Skill, src, dst string, opts Options) {
	r.printHeader(index, "external-skill", []headerField{
		{label: "name", value: skill.Name},
		{label: "url", value: skill.URL},
		{label: "source", value: src},
		{label: "destination", value: dst},
	}, opts)
}

func (r Runner) printHeader(index int, kind string, fields []headerField, opts Options) {
	mode := "DEPLOY"
	modeColor := colorGreen
	if opts.DryRun {
		mode = "DRY-RUN"
		modeColor = colorYellow
	}
	fmt.Fprintf(r.out, "\n%s item[%d] %s\n", colorize("["+mode+"]", modeColor, opts), index, colorize(kind, colorCyan, opts))
	labelWidth := maxHeaderLabelWidth(fields)
	for _, field := range fields {
		label := field.label + ":"
		padding := strings.Repeat(" ", labelWidth-len(field.label)+1)
		fmt.Fprintf(r.out, "  %s%s%s\n", colorize(label, colorFaint, opts), padding, field.value)
	}
}

func maxHeaderLabelWidth(fields []headerField) int {
	width := 0
	for _, field := range fields {
		if len(field.label) > width {
			width = len(field.label)
		}
	}
	return width
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
		colorize(fmt.Sprintf("%d copied", report.CopiedFiles), colorGreen, opts),
		colorize(fmt.Sprintf("%d dirs", report.CreatedDirs), colorCyan, opts),
		colorize(fmt.Sprintf("%d skipped", report.Skipped), colorYellow, opts),
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

func fileOptions(opts Options) fileops.Options {
	return fileops.Options{DryRun: opts.DryRun}
}
