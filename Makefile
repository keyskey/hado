.PHONY: help setup build check-docker lint lint-go lint-yaml lint-markdown fmt fmt-go fmt-check test readiness-check ci

GO_FILES := $(shell git ls-files '*.go')
BINARY := bin/hado
GO_BIN := $(shell go env GOPATH 2>/dev/null)/bin
COVERPROFILE ?= coverage.out
GOBCE ?= $(GO_BIN)/gobce
GOBCE_PACKAGE ?= github.com/keyskey/gobce/cmd/gobce@latest
READINESS_COVERAGE ?= hado-coverage.json
READINESS_MANIFEST ?= hado.yaml
READINESS_STANDARD ?= standards/cli-service.yaml
# Pinned image avoids npx + rolling npm deps on node:22-alpine (CI flakiness).
MARKDOWNLINT_CLI2_IMAGE ?= davidanson/markdownlint-cli2:v0.22.1

help:
	@echo "Available targets:"
	@echo "  make setup        # Verify Go and install development tools"
	@echo "  make build        # Build hado CLI binary"
	@echo "  make lint         # Run YAML, Markdown, and Go lint checks"
	@echo "  make fmt          # Format Go source files"
	@echo "  make fmt-check    # Check Go formatting (CI equivalent)"
	@echo "  make test         # Run Go tests"
	@echo "  make readiness-check # Generate HADO coverage evidence and evaluate readiness"
	@echo "  make ci           # Run fmt-check, lint, and test"

setup:
	@command -v go >/dev/null 2>&1 || { echo "go is required."; exit 1; }
	@echo "Go toolchain is available."
	go install "$(GOBCE_PACKAGE)"
	@if [ -n "$${GITHUB_PATH:-}" ]; then \
		echo "$(GO_BIN)" >> "$$GITHUB_PATH"; \
	fi
	@echo "gobce is installed at $(GOBCE)."

build:
	@mkdir -p bin
	go build -o "$(BINARY)" ./cmd/hado

check-docker:
	@command -v docker >/dev/null 2>&1 || { echo "docker is required for YAML/Markdown lint."; exit 1; }

lint: lint-yaml lint-markdown lint-go

lint-yaml: check-docker
	docker run --rm -v "$$(pwd):/work" -w /work python:3.12-alpine sh -lc "pip install --no-cache-dir yamllint >/dev/null && yamllint ."

lint-markdown: check-docker
	docker run --rm -v "$$(pwd):/work" -w /work "$(MARKDOWNLINT_CLI2_IMAGE)" "**/*.md" "#node_modules" "#.git" "#.tools"

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

readiness-check:
	@command -v "$(GOBCE)" >/dev/null 2>&1 || { echo "gobce is required. Run: make setup"; exit 1; }
	go test ./... -coverprofile="$(COVERPROFILE)"
	"$(GOBCE)" analyze --coverprofile "$(COVERPROFILE)" --format json --output "$(READINESS_COVERAGE)"
	go run ./cmd/hado evaluate --standard "$(READINESS_STANDARD)" --manifest "$(READINESS_MANIFEST)"

ci: fmt-check lint test
