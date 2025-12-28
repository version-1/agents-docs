
build-generator:
	cd scripts/generator && go build -o ../../bin/generator cmd/main.go

gen-docs:
	$(MAKE) build-generator
	./bin/generator -input=./docs/ja -output=./out/.codex/
