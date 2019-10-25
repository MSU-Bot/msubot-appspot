.PHONY: help

help:  ## Prints this message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

format: ## Format the go files
	gofmt -w ./..

build: ## Builds the binary
	go build .

run-binary: build ## Runs the binary
	./msubot-appspot

run-fast: ## Runs the server withough building
	go run .
	
deploy-staging: ## Deploys to the staging enviroment
	gcloud app deploy app-staging.yaml --project msubot-staging

deploy-staging-cron: ## Deploys cron to staging
	gcloud app deploy cron.yaml app-staging.yaml --project msubot-staging

deploy-prod: ## deploy everything to Production
	gcloud app deploy cron.yaml app-prod.yaml --project msu-bot


setup: ## Sets up tooling and other dependencies
	echo "Sooon!"