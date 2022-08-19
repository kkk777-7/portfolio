.PHONY: build test local-run local-up dev-deploy

build:
	sam build

test:
	cd ./hello-world/ && \
	go test -v ./...

local-run: build
	sam local invoke \
	--parameter-overrides 'ENV=local'

local-up: build
	sam local start-api \
	--parameter-overrides 'ENV=local'


dev-deploy: build
	sam deploy \
	--parameter-overrides 'ENV=dev' \
	--force-upload \
	--no-fail-on-empty-changeset