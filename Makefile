#!make
.PHONY: migrate server
include .env


export TWITCH_OAUTH
export TWITCH_USERNAME
export SECRET
export NEXT_PUBLIC_BASE_URL
export NEXT_PUBLIC_API_BASE_URL
export NEXT_PUBLIC_TWITCH_CLIENT_ID
export TWITCH_CLIENT_SECRET
export DSN

build_server:
	go run main.go

web: 
	yarn dev

migrate:
	go run main.go migrate

proxy:
	flyctl proxy 5432 -a gempbot-db

test:
	go test ./internal/...

docker: 
	docker build . -t gempbot

run_docker:
	docker run --env-file=.env -p 3010:3010 gempbot

tunnel:
	lt --print-requests --port 3010 --subdomain gempir
