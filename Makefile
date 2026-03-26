.DEFAULT_GOAL := help

run: ## Run the game
	go run .
.PHONY: run

test: ## Run tests
	go test ./...
.PHONY: test

build: ## Make a macOS executable binary
	go build -o bin/mac/pacvim .
.PHONY: build

build-win: ## Make a Windows executable binary
	GOOS=windows GOARCH=amd64 go build -o bin/win/pacvim.exe .
.PHONY: build-win

wasm: ## Build WebAssembly (output: docs/wasm/)
	mkdir -p docs/wasm
	GOOS=js GOARCH=wasm go build -o docs/wasm/pacvim.wasm .
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" docs/wasm/
.PHONY: wasm

clean: ## Remove binary files
	rm -f bin/mac/pacvim bin/win/pacvim.exe docs/wasm/pacvim.wasm
.PHONY: clean

help:
	@echo "Usage:\n    make \033[36m<command>\033[0m\n\nCommands:" >&2
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "%4s\033[36m%-10s\033[0m\n%8s%s\n", "", $$1, "", $$2}' >&2
.PHONY: help
