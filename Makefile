
build-generator:
	cd scripts/generator && go build -o ../../bin/generator cmd/main.go

gen-docs:
	$(MAKE) build-generator
	./bin/generator -input=./docs/ja -output=./out/.codex/
	./bin/generator -input=./docs/ja -output=./out/.claude/ -flat-skills

tree-out:
	tree ./out/.codex
	tree ./out/.claude

deploy-codex-docs:
	mkdir -p ~/.codex ~/.codex/skills
	cp -R ./out/.codex/skills/. ~/.codex/skills/
	cp -R ./out/.codex/agents ~/.codex/.
	cp ./out/.codex/Agents.md ~/.codex/.

deploy-claude-docs:
	mkdir -p ~/.claude ~/.claude/skills
	cp -R ./out/.claude/skills/. ~/.claude/skills/
	cp -R ./out/.claude/agents ~/.claude/.
	cp ./claude/CLAUDE.md ~/.claude/.
	cp ./claude/settings.json ~/.claude/.

apply: gen-docs deploy-codex-docs deploy-claude-docs
