all: 
	echo "Yo!"
	
deploy-staging:
	gcloud app deploy app-staging.yaml --project msubot-staging

deploy-staging-cron:
	gcloud app deploy cron.yaml app-staging.yaml --project msubot-staging

deploy-prod:
	gcloud app deploy cron.yaml app-prod.yaml --project msu-bot


generate:
	oapi-codegen -generate="types" -package="api" -o="server/gen/api/types.go" ./api/MSUBot-Appengine-1.0.0.yaml
	oapi-codegen -generate="server" -package="api" -o="server/gen/api/service.go" ./api/MSUBot-Appengine-1.0.0.yaml