.PHONY: build test local-run local-up dev-deploy prd-deploy

build:
	sam build --use-container

test:
	cd ./line-bot/ && \
	go test -v ./...

local-run: build
	sam local invoke \
	--parameter-overrides 'ENV=local AWSREGION=ap-northeast-1' 

local-up: build
	sam local start-api \
	--parameter-overrides 'ENV=local AWSREGION=ap-northeast-1' 

dev-deploy: build
	sam deploy \
	--config-env dev \
	--parameter-overrides 'ENV=dev AWSREGION=ap-northeast-1' \
	--force-upload \
	--no-fail-on-empty-changeset

prd-deploy: build
	sam deploy \
	--config-env prd \
	--parameter-overrides 'ENV=prd AWSREGION=ap-northeast-1' \
	--force-upload \
	--no-fail-on-empty-changeset