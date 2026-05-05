# Run `just` to list recipes. See https://just.systems/man/en/

default:
    @just --list

generate:
    go generate ./...

build:
    go build -o bin/server ./cmd/server/...

tidy:
    go mod tidy

all:
    echo "Yo!"

deploy-staging:
    gcloud app deploy app-staging.yaml --project msubot-staging

deploy-staging-cron:
    gcloud app deploy cron.yaml app-staging.yaml --project msubot-staging

deploy-prod:
    gcloud app deploy cron.yaml app-prod.yaml --project msu-bot
