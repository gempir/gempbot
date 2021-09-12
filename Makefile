.PHONY: frontend functions

build: frontend functions

frontend:
	yarn build

functions:
	GOBIN=${PWD}/netlify-functions go install ./functions/...