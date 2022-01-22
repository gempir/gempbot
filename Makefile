#!make
.PHONY: migrate server
include .env

export VERCEL
export VERCEL_ENV
export VERCEL_URL
export VERCEL_GIT_PROVIDER
export VERCEL_GIT_REPO_SLUG
export VERCEL_GIT_REPO_OWNER
export VERCEL_GIT_REPO_ID
export VERCEL_GIT_COMMIT_REF
export VERCEL_GIT_COMMIT_SHA
export VERCEL_GIT_COMMIT_MESSAGE
export VERCEL_GIT_COMMIT_AUTHOR_LOGIN
export VERCEL_GIT_COMMIT_AUTHOR_NAME
export NEXT_PUBLIC_API_BASE_URL
export SEVEN_TV_TOKEN
export BTTV_TOKEN
export WEBHOOK_BASE_URL
export NEXT_PUBLIC_LOGTAIL_SOURCE_TOKEN
export LOGTAIL_SOURCE_TOKEN
export TWITCH_OAUTH
export TWITCH_USERNAME
export NEXT_PUBLIC_BASE_URL
export PLANETSCALE_DB
export PLANETSCALE_DB_USERNAME
export PLANETSCALE_DB_PASSWORD
export PLANETSCALE_DB_HOST
export NEXT_PUBLIC_TWITCH_CLIENT_ID
export SECRET
export TWITCH_CLIENT_SECRET
export TWITCH_CLIENT_ID

build_server:
	go build

migrate:
	go run cmd/migrate/main.go

test:
	go test ./internal/...

docker: 
	docker build . -t gempbot

run_docker:
	docker run --env-file=.env -p 3010:3010 gempbot

tunnel:
	lt --port 3010 --subdomain gempir
