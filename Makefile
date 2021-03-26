BIN_DIR=bin

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

code-check:
	make lint
	make tests

code-clean:
	make imports
	make format

lint:
	golint -set_exit_status ./clients/...
	golint -set_exit_status ./config/...
	golint -set_exit_status ./database/...
	golint -set_exit_status ./dto/...
	golint -set_exit_status ./log/...
	golint -set_exit_status ./services/...

imports:
	goimports -d -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:
	go fmt $(shell go list ./... | grep -v /vendor/)

tests:
	go test ./...

create-if-not-exists-env:
	if [ ! -f .env ]; then cp .env.example .env; fi


build:
	make vendor
	make create-if-not-exists-env
	go build -o ./bin/toggl-to-jira-osx ./main.go

build-cp:
	# For CROSS platform build there must be Docker installed.
	docker pull karalabe/xgo-latest
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(BIN_DIR)/ --out toggl-to-jira ./

.PHONY: vendor