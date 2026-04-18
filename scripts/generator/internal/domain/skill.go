package domain

import "generator/internal/fsadapter"

// SkillGenerator はスキルファイルの生成戦略を表すインターフェース。
type SkillGenerator interface {
	Generate(fsys fsadapter.FileSystem, inRoot, outRoot string) error
}

// pathRespectSkillGenerator はソースのディレクトリ階層を維持してスキルを生成する。
// Codex のようにネストされたスキル構造をサポートするツール向け。
type pathRespectSkillGenerator struct{}

func NewPathRespectSkillGenerator() SkillGenerator {
	return &pathRespectSkillGenerator{}
}

func (g *pathRespectSkillGenerator) Generate(fsys fsadapter.FileSystem, inRoot, outRoot string) error {
	return fsadapter.WalkResources(fsys, inRoot, outRoot, func(outRoot, relPath string, _ []byte) (string, error) {
		return SkillOutputPath(outRoot, relPath), nil
	})
}

// flatSkillGenerator はフロントマターの name を使い skills/<name>/SKILL.md のフラット構造で生成する。
// Claude のようにスキルをフラットに配置する必要があるツール向け。
type flatSkillGenerator struct{}

func NewFlatSkillGenerator() SkillGenerator {
	return &flatSkillGenerator{}
}

func (g *flatSkillGenerator) Generate(fsys fsadapter.FileSystem, inRoot, outRoot string) error {
	return fsadapter.WalkResources(fsys, inRoot, outRoot, func(outRoot, relPath string, content []byte) (string, error) {
		fm, err := ParseFrontmatter(content)
		if err != nil {
			return "", err
		}
		return FlatSkillOutputPath(outRoot, fm.Name), nil
	})
}
