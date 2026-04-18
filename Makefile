
build-generator:
	cd scripts/generator && go build -o ../../bin/generator cmd/main.go

gen-docs:
	$(MAKE) build-generator
	./bin/generator -input=./docs/ja -output=./out/.codex/ -mode=codex
	./bin/generator -input=./docs/ja -output=./out/.claude/ -mode=claude

tree-out:
	tree ./out/.codex
	tree ./out/.claude

deploy-claude-specific-docs:
	cp ./claude/CLAUDE.md ~/.claude/.
	cp ./claude/settings.json ~/.claude/.

deploy-codex-specific-docs:
	cp ./out/.codex/Agents.md ~/.codex/.
	cp codex/config.toml ~/.codex/.

deploy-codex-docs:
	mkdir -p ~/.codex ~/.codex/skills
	cp -R ./out/.codex/skills/. ~/.codex/skills/
	cp -R ./out/.codex/agents ~/.codex/.

deploy-claude-docs:
	mkdir -p ~/.claude ~/.claude/skills
	cp -R ./out/.claude/skills/. ~/.claude/skills/
	cp -R ./out/.claude/agents ~/.claude/.

apply: gen-docs deploy-codex-docs deploy-claude-docs deploy-codex-specific-docs deploy-claude-specific-docs
