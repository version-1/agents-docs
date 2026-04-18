package app

import (
	"fmt"
	"path/filepath"

	"generator/internal/domain"
	"generator/internal/fsadapter"
)

type Generator struct {
	fs     fsadapter.FileSystem
	skills domain.SkillGenerator
}

func NewGenerator(fs fsadapter.FileSystem, skills domain.SkillGenerator) *Generator {
	return &Generator{fs: fs, skills: skills}
}

func (g *Generator) Run(inRoot, outRoot string) error {
	inExists, err := fsadapter.DirExists(g.fs, inRoot)
	if err != nil {
		return fmt.Errorf("check input dir: %w", err)
	}
	if !inExists {
		return fmt.Errorf("input dir does not exist: %s", inRoot)
	}

	outExists, err := fsadapter.DirExists(g.fs, outRoot)
	if err != nil {
		return fmt.Errorf("check output dir: %w", err)
	}
	if outExists {
		if err := fsadapter.RemoveDirContents(g.fs, outRoot); err != nil {
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

func copyDocs(fsys fsadapter.FileSystem, inRoot, outRoot string) error {
	agentMdPath := filepath.Join(inRoot, "Agents.md")
	if err := fsadapter.CopyFile(fsys, agentMdPath, filepath.Join(outRoot, "Agents.md"), 0o644); err != nil {
		return err
	}

	agentsDir := "agents"
	if err := fsadapter.CopyDir(fsys, filepath.Join(inRoot, agentsDir), filepath.Join(outRoot, agentsDir)); err != nil {
		return err
	}

	return nil
}
