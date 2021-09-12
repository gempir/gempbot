.PHONY: build

build:
	yarn build
	GOBIN=${PWD}/netlify-functions go install ./functions/...