DIR = $(shell basename $(CURDIR))

PROJECT = go-events-enricher
export APP_NAME := go-events-enricher
export VERSION := $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))
export BUILD_TIME:= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

LDFLAGS = -ldflags "-s -w -X github.com/Alveona/go-events-enricher/version.VERSION=${VERSION} -X github.com/Alveona/go-events-enricher/version.BUILD_TIME=${BUILD_TIME}"

SWAGGER_VERSION=$(shell swagger version)
define REQUIRED_SWAGGER_VERSION
version: v0.24.0 commit: 094421c7a2e1a7982cb28fdc0fc2bda610bbcc56
endef

checkswagger:
    ifneq (${SWAGGER_VERSION}, ${REQUIRED_SWAGGER_VERSION})
	    @echo "Required swagger version: ${REQUIRED_SWAGGER_VERSION}"
	    @echo "Installed swagger version: ${SWAGGER_VERSION}"
	    @echo "Download correct version there: https://github.com/go-swagger/go-swagger/releases/tag/v0.24.0"
	    @exit 1
    endif


init: pre-install generate

pre-install:
	@echo "Installing swagger..."
	@go get -u github.com/go-swagger/go-swagger/cmd/swagger@v0.24.0
	@echo "Installing golangci-lint..."
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint


docker-server-run: #? run application using docker
	@echo "Starting server in docker container..."
	@docker compose -f docker-compose.yml up

docker-server-down: #? stop application in docker
	@echo "Stopping server in docker container..."
	@docker compose -f docker-compose.yml down

mocks:
	@echo "Regenerate mocks..."
	@go generate ./...

dep: checkpath
	go mod tidy
	go mod vendor

run: #? run application
	@echo "Running app..."
	go run $(LDFLAGS) main.go

test:
	@echo "Testing Go packages..."
	@go test ./app/... -cover -count=1

lint: #? run golangci-lint
	@echo "Running golangci-lint..."
	@golangci-lint run

generate: clean generate-server


generate-server: #? generate server
	@echo "Generating server..."
	@swagger generate server -f ./swagger-doc/swagger.yml -t ./app/generated --exclude-main

clean: #? remove generated files
	@echo "Removing generated files..."
	@rm -rf ./app/generated/clients
	@rm -rf ./app/generated/models
	@rm -rf ./app/generated/restapi/operations
	@rm -rf ./app/generated/restapi/doc.go
	@rm -rf ./app/generated/restapi/embedded_spec.go
	@rm -rf ./app/generated/restapi/server.go
