
build-deploy:
	cd scripts/deploy && go build -o ../../bin/deploy cmd/main.go

deploy-dry-run:
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json -external-skills=./external-skills.json -dry-run

deploy:
	$(MAKE) build-deploy
	./bin/deploy -config=./deploy.json -external-skills=./external-skills.json

apply: deploy
