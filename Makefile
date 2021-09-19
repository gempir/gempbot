#!make
.PHONY: migrate bot
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

bot:
	go run cmd/bot/main.go

refresh:
	go run cmd/refresh/main.go

docker: 
	sudo docker build . -t rg.fr-par.scw.cloud/funcscwgempbotp9rlmser/bot