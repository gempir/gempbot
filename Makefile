#!make
.PHONY: *
include .env
include .env.development

export $(shell sed 's/=.*//' .env)
export $(shell sed 's/=.*//' .env.development)

build_server:
	go run main.go

web: 
	yarn dev

migrate:
	go run main.go migrate

proxy:
	flyctl proxy 5433 -a gempbot-db

test:
	go test ./internal/...

docker: 
	docker build . -t gempbot

run_docker:
	docker run --env-file=.env -p 3010:3010 gempbot

tunnel:
	lt --print-requests --port 3010 --subdomain gempir
