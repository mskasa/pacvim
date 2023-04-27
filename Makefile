.DEFAULT_GOAL := help

fmt: ## go fmt
	go fmt
.PHONY: fmt

lint: fmt ## golangci-lint run
	golangci-lint run
.PHONY: lint

deps: lint ## go mod tidy
	go mod tidy
.PHONY: deps

test: deps ## go test -short
	go test -short
.PHONY: test

cover: ## create cover.html
	go test -cover -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

build: test ## Make a macOS executable binary
	go build -o bin/mac/pacvim .
.PHONY: build

build-win: test ## Make a Windows executable binary
	GOOS=windows GOARCH=amd64 go build -o bin/win/pacvim.exe .
.PHONY: build-win

clean: ## Remove binary files
	rm ./bin/mac/pacvim
	rm ./bin/win/pacvim.exe
.PHONY: clean

help:
	@echo "Usage:\n    make \033[36m<command>\033[0m\n\nCommands:" >&2
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "%4s\033[36m%-10s\033[0m\n%8s%s\n", "", $$1, "", $$2}' >&2
.PHONY: help
