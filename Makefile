#!make
.PHONY: migrate server
include .env

export PLANETSCALE_DB
export PLANETSCALE_DB_USERNAME
export PLANETSCALE_DB_PASSWORD
export PLANETSCALE_DB_HOST
export TWITCH_CLIENT_ID
export TWITCH_CLIENT_SECRET
export TWITCH_USERNAME
export TWITCH_OAUTH
export SECRET
export NEXT_PUBLIC_BASE_URL
export VERCEL_ENV

migrate:
	go run cmd/migrate/main.go

server:
	go run server/main.go

refresh:
	go run cmd/refresh/main.go

eventsub:
	go run cmd/eventsub/main.go

docker: 
	docker build . -t gempbot

run_docker:
	docker run --env-file=.env -p 3010:3010 gempbot

tunnel:
	lt --port 3000 --subdomain gempir
