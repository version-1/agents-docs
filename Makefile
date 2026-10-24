
LOCAL_CONFIG_FLAG := $(if $(wildcard config.local.json),-local-config=./config.local.json)
TAG_FLAG := $(if $(TAG),-tag='$(TAG)')
ifeq ($(TAG),)
PLUGIN_DEPLOY_TARGET := build-plugins
PLUGIN_DRY_RUN_TARGET := build-plugins
else ifeq ($(TAG),skills)
PLUGIN_DEPLOY_TARGET := build-plugins
PLUGIN_DRY_RUN_TARGET := build-plugins
endif

build-plugin-builder:
	cd scripts/deploy && go build -o ../../bin/plugins cmd/plugins/main.go

build-plugins: build-plugin-builder
	./bin/plugins -groups=./skill-groups.json -external-skills=./external-skills.json -output=./out/marketplace

plugins-dry-run:
	$(MAKE) build-plugin-builder
	./bin/plugins -groups=./skill-groups.json -external-skills=./external-skills.json -output=./out/marketplace -dry-run

deploy-plugins-dry-run: build-plugins
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json -tag='skills' -dry-run

deploy-plugins:
	$(MAKE) build-plugins
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json -tag='skills'

build-deploy:
	cd scripts/deploy && go build -o ../../bin/deploy cmd/main.go

deploy-dry-run:
	$(if $(PLUGIN_DRY_RUN_TARGET),$(MAKE) $(PLUGIN_DRY_RUN_TARGET))
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json $(LOCAL_CONFIG_FLAG) $(TAG_FLAG) -dry-run

deploy:
	$(if $(PLUGIN_DEPLOY_TARGET),$(MAKE) $(PLUGIN_DEPLOY_TARGET))
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json $(LOCAL_CONFIG_FLAG) $(TAG_FLAG)

apply: deploy
