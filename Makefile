.PHONY: help setup build check-docker lint lint-go lint-yaml lint-markdown fmt fmt-go fmt-check test ci

GO_FILES := $(shell git ls-files '*.go')
BINARY := bin/hado

help:
	@echo "Available targets:"
	@echo "  make setup        # Verify local Go toolchain"
	@echo "  make build        # Build hado CLI binary"
	@echo "  make lint         # Run YAML, Markdown, and Go lint checks"
	@echo "  make fmt          # Format Go source files"
	@echo "  make fmt-check    # Check Go formatting (CI equivalent)"
	@echo "  make test         # Run Go tests"
	@echo "  make ci           # Run fmt-check, lint, and test"

setup:
	@command -v go >/dev/null 2>&1 || { echo "go is required."; exit 1; }
	@echo "Go toolchain is available."

build:
	@mkdir -p bin
	go build -o "$(BINARY)" ./cmd/hado

check-docker:
	@command -v docker >/dev/null 2>&1 || { echo "docker is required for YAML/Markdown lint."; exit 1; }

lint: lint-yaml lint-markdown lint-go

lint-yaml: check-docker
	docker run --rm -v "$$(pwd):/work" -w /work python:3.12-alpine sh -lc "pip install --no-cache-dir yamllint >/dev/null && yamllint ."

lint-markdown: check-docker
	docker run --rm -v "$$(pwd):/work" -w /work node:22-alpine sh -lc "npx --yes markdownlint-cli2"

lint-go:
	go vet ./...

fmt: fmt-go

fmt-go:
	@if [ -n "$(GO_FILES)" ]; then \
		gofmt -w $(GO_FILES); \
	fi

fmt-check:
	@if [ -z "$(GO_FILES)" ]; then \
		exit 0; \
	fi; \
	unformatted="$$(gofmt -l $(GO_FILES))"; \
	if [ -n "$$unformatted" ]; then \
		echo "These Go files need gofmt:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

test:
	go test ./...

ci: fmt-check lint test
