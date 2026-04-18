package app

import "generator/internal/domain"

// pathRespectSkillGenerator はソースのディレクトリ階層を維持してスキルを生成する。
// Codex のようにネストされたスキル構造をサポートするツール向け。
type pathRespectSkillGenerator struct{}

func NewPathRespectSkillGenerator() SkillGenerator {
	return &pathRespectSkillGenerator{}
}

func (g *pathRespectSkillGenerator) Generate(fsys FileSystem, inRoot, outRoot string) error {
	return walkSkills(fsys, inRoot, outRoot, func(outRoot, relPath string, _ []byte) (string, error) {
		return domain.SkillOutputPath(outRoot, relPath), nil
	})
}

// flatSkillGenerator はフロントマターの name を使い skills/<name>/SKILL.md のフラット構造で生成する。
// Claude のようにスキルをフラットに配置する必要があるツール向け。
type flatSkillGenerator struct{}

func NewFlatSkillGenerator() SkillGenerator {
	return &flatSkillGenerator{}
}

func (g *flatSkillGenerator) Generate(fsys FileSystem, inRoot, outRoot string) error {
	return walkSkills(fsys, inRoot, outRoot, func(outRoot, relPath string, content []byte) (string, error) {
		name, err := domain.ParseName(content)
		if err != nil {
			return "", err
		}
		return domain.FlatSkillOutputPath(outRoot, name), nil
	})
}
