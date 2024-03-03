#!make
.PHONY: *
include .env.development
include .env

export $(shell sed 's/=.*//' .env.development)
export $(shell sed 's/=.*//' .env)

build_server:
	go run main.go

ysweet:
	cd web && yarn ysweet-dev

test:
	go test ./internal/...

staticcheck:
	staticcheck ./...

web: 
	cd web && yarn dev

deploy:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o gempbot main.go
	ssh -o StrictHostKeyChecking=no -p 32022 -i ansible/.ssh_key ubuntu@o1.gempir.com "sudo systemctl stop gempbot"
	rsync -avz -e "ssh -o StrictHostKeyChecking=no -p 32022 -i ansible/.ssh_key" gempbot ubuntu@o1.gempir.com:/home/gempbot/
	ssh -o StrictHostKeyChecking=no -p 32022 -i ansible/.ssh_key ubuntu@o1.gempir.com "sudo chown gempbot:gempbot /home/gempbot/gempbot"
	ssh -o StrictHostKeyChecking=no -p 32022 -i ansible/.ssh_key ubuntu@o1.gempir.com "sudo systemctl restart gempbot-migrate && sudo systemctl start gempbot"

ansible:
	cd ansible && ansible-vault decrypt ssh_key.vault --output=.ssh_key
	chmod 600 ansible/.ssh_key

provision:
	cd ansible && OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES ansible-playbook -i hosts playbook.yml --private-key=.ssh_key

migrate:
	go run main.go migrate

docker: 
	docker build . -t gempbot

run_docker:
	docker run --env-file=.env -p 3010:3010 gempbot

tunnel:
	npx localtunnel --port 3010 --subdomain gempbot

ollama:
	ssh -N -R 11434:127.0.0.1:11434 o1