BIN_DIR=bin

go-vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

code-check:
	make lint
	make tests

code-clean:
	make imports
	make format

lint:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./src/...

imports:
	goimports -d -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:
	go fmt $(shell go list ./... | grep -v /vendor/)

tests:
	go test ./...

create-if-not-exists-env:
	if [ ! -f .env ]; then cp .env.example .env; fi

build-binary:
	make go-vendor
	make create-if-not-exists-env
	go build -o ./bin/toggl-to-jira ./main.go

build:
	docker build -t toggl-to-jira .

check-security:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

generate-mocks:
	(which mockery || go install github.com/vektra/mockery/v2@latest)
	mockery --all --dir=src --keeptree
	make go-vendor

release:
	make build-cp
	zip -r release.zip bin/*