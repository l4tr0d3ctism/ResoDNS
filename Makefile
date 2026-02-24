PKG := resodns
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')
REVISION := $(shell git rev-parse --short HEAD)
VERSION ?= 1.0.0
RELEASE_NAME := resodns-$(VERSION)

.SILENT: ;
.PHONY: all build release

all: build

lint: ## Lint the files
	golint -set_exit_status $(PKG_LIST)
	staticcheck ./...

test: ## Run unit tests
	go fmt $(PKG_LIST)
	go vet $(PKG_LIST)
	go test -race -timeout 30s -cover -count 1 $(PKG_LIST)

msan: ## Run memory sanitizer
	go test -msan $(PKG_LIST)

build: ## Build the binary file
	go build -trimpath -ldflags="-s -w" -o resodns .

cover: ## Code coverage
	go test -coverprofile=cover.out $(PKG_LIST)

release: ## Package resodns + massdns binaries into dist/ (no build; binaries must exist)
	@mkdir -p dist/$(RELEASE_NAME)
	@cp massdns/bin/massdns dist/$(RELEASE_NAME)/
	@if [ -f resodns ]; then cp resodns dist/$(RELEASE_NAME)/; fi
	@tar czvf dist/$(RELEASE_NAME).tar.gz -C dist $(RELEASE_NAME)
	@rm -rf dist/$(RELEASE_NAME)
	@echo "Release: dist/$(RELEASE_NAME).tar.gz"

clean: ## Remove previous build
	rm -f cover.out
	go clean

help: ## Display this help screen
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
