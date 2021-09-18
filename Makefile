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

migrate:
	go run cmd/migrate/main.go

bot:
	go run cmd/bot/main.go

