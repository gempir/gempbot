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

ansible:
	cd ansible && ansible-vault decrypt ssh_key.vault --output=.ssh_key
	python3 -m pip install jmespath

provision:
	cd ansible && OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES ansible-playbook -i hosts playbook.yml --private-key=.ssh_key

migrate:
	go run main.go migrate

proxy:
	flyctl proxy 5000:5432 -a gempbot-db

docker: 
	docker build . -t gempbot

run_docker:
	docker run --env-file=.env -p 3010:3010 gempbot

tunnel:
	bore -id gempbot -s bore.services -p 2200 -ls localhost -lp 3010
