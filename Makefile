all: 
	echo "Yo!"
	
deploy-staging:
	gcloud app deploy app-staging.yaml --project msubot-staging

deploy-staging-cron:
	gcloud app deploy cron.yaml app-staging.yaml --project msubot-staging

deploy-prod:
	gcloud app deploy cron.yaml app-prod.yaml --project msu-bot


generate:
	go generate ./...