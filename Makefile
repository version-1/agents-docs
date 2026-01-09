
build-generator:
	cd scripts/generator && go build -o ../../bin/generator cmd/main.go

gen-docs:
	$(MAKE) build-generator
	./bin/generator -input=./docs/ja -output=./out/.codex/

deploy-codex-docs:
	cp -pr ./out/.codex/skills/* ~/.codex/skills/
	cp -pr ./out/.codex/agents ~/.codex/.
	cp -p ./out/.codex/Agents.md ~/.codex/.
