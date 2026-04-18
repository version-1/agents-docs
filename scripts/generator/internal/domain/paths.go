package domain

import (
	"path/filepath"
	"strings"
)

var excludedRelPaths = map[string]struct{}{
	filepath.Join("skills", "skills.md"): {},
	"Agents.md":                         {},
}

func IsMarkdown(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".md")
}

func IsExcluded(relPath string) bool {
	_, ok := excludedRelPaths[filepath.Clean(relPath)]
	return ok
}

func SkillOutputPath(outRoot, relPath string) string {
	relNoExt := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	return filepath.Join(outRoot, relNoExt, "SKILL.md")
}

// FlatSkillOutputPath はフロントマターの name を使い、
// skills/<name>/SKILL.md のフラットな構造で出力パスを返す。
func FlatSkillOutputPath(outRoot, name string) string {
	return filepath.Join(outRoot, "skills", name, "SKILL.md")
}
