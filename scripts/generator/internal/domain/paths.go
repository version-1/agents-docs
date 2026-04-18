package domain

import (
	"path/filepath"
	"strings"
)

func SkillOutputPath(outRoot, relPath string) string {
	relNoExt := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	return filepath.Join(outRoot, relNoExt, "SKILL.md")
}

// FlatSkillOutputPath はフロントマターの name を使い、
// skills/<name>/SKILL.md のフラットな構造で出力パスを返す。
func FlatSkillOutputPath(outRoot, name string) string {
	return filepath.Join(outRoot, "skills", name, "SKILL.md")
}
