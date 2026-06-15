
LOCAL_CONFIG_FLAG := $(if $(wildcard config.local.json),-local-config=./config.local.json)

build-deploy:
	cd scripts/deploy && go build -o ../../bin/deploy cmd/main.go

deploy-dry-run:
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json -external-skills=./external-skills.json $(LOCAL_CONFIG_FLAG) -dry-run

deploy:
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json -external-skills=./external-skills.json $(LOCAL_CONFIG_FLAG)

apply: deploy
