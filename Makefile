#!make
.PHONY: *
include .env.development
include .env

export $(shell sed 's/=.*//' .env.development)
export $(shell sed 's/=.*//' .env)

build_server:
	go run main.go

test:
	go test ./internal/...

staticcheck:
	staticcheck ./...

web: 
	yarn dev

migrate:
	go run main.go migrate

proxy:
	flyctl proxy 5000:5432 -a gempbot-db

docker: 
	docker build . -t gempbot

run_docker:
	docker run --env-file=.env -p 3010:3010 gempbot

tunnel:
	lt --print-requests --port 3010 --subdomain gempir
