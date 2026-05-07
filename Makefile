.PHONY: help setup setup-hooks bootstrap-go ensure-go build check-docker ci-lint lint lint-go lint-yaml lint-markdown fmt fmt-go fmt-check test readiness-check pre-pr gen-manifest-doc

# Optional local toolchain: official tarball under .gitignored .tools/go (see bootstrap-go).
# Prefer it when present so Make works without a global install; otherwise use `go` on PATH.
TOOLS_GO := $(CURDIR)/.tools/go/bin/go
GO_CMD = $(if $(wildcard $(TOOLS_GO)),$(TOOLS_GO),go)
GOFMT_CMD = $(if $(wildcard $(TOOLS_GO)),$(dir $(TOOLS_GO))gofmt,gofmt)
GO_BOOTSTRAP_VERSION ?= 1.22.12

GO_FILES := $(shell git ls-files '*.go')
BINARY := bin/hado
GO_BIN = $(shell $(GO_CMD) env GOPATH 2>/dev/null)/bin
GOBCE = $(GO_BIN)/gobce
GOBCE_PACKAGE ?= github.com/keyskey/gobce/cmd/gobce@latest
COVERPROFILE ?= coverage.out
READINESS_COVERAGE ?= hado-coverage.json
READINESS_MANIFEST ?= hado.yaml
READINESS_STANDARD ?= standards/cli-service.yaml
MARKDOWNLINT_CLI2_IMAGE ?= davidanson/markdownlint-cli2:v0.22.1

help:
	@echo "Available targets:"
	@echo "  make bootstrap-go # Install Go $(GO_BOOTSTRAP_VERSION) into .tools/go (network; not global)"
	@echo "  make ensure-go    # Use .tools/go if complete, else PATH go, else bootstrap"
	@echo "  make setup        # ensure-go + go install gobce"
	@echo "  make build        # Build hado CLI binary"
	@echo "  make lint         # Run YAML, Markdown, and Go lint checks"
	@echo "  make fmt          # Format Go source files"
	@echo "  make fmt-check    # Used in make ci-lint (does not run go test)"
	@echo "  make test         # Run Go tests"
	@echo "  make gen-manifest-doc # Regenerate docs/hado.manifest.reference.yaml (commented reference manifest)"
	@echo "  make readiness-check # Generate HADO coverage evidence and run charge/fire"
	@echo "  make ci-lint      # fmt-check + lint (GitHub Lint job + pre-push; no go test)"
	@echo "  make setup-hooks  # pre-push runs: make ci-lint"
	@echo "  make pre-pr       # ci-lint + test (local only; go test also runs in CI via readiness-check)"

ensure-go:
	@if [ -x "$(TOOLS_GO)" ] && [ -f "$(CURDIR)/.tools/go/src/bytes/buffer.go" ]; then exit 0; fi
	@if command -v go >/dev/null 2>&1 && [ ! -f "$(TOOLS_GO)" ]; then exit 0; fi
	@$(MAKE) bootstrap-go

bootstrap-go:
	@set -e; \
	root="$(CURDIR)"; \
	if [ -x "$$root/.tools/go/bin/go" ] && "$$root/.tools/go/bin/go" version 2>/dev/null | grep -q "go$(GO_BOOTSTRAP_VERSION)" \
		&& [ -f "$$root/.tools/go/src/bytes/buffer.go" ]; then \
		echo "bootstrap-go: .tools/go already has go$(GO_BOOTSTRAP_VERSION)"; \
		exit 0; \
	fi; \
	ver="$(GO_BOOTSTRAP_VERSION)"; \
	case $$(uname -s) in \
		Darwin) os=darwin ;; \
		Linux) os=linux ;; \
		*) echo "bootstrap-go: unsupported OS ($$(uname -s))"; exit 1 ;; \
	esac; \
	case $$(uname -m) in \
		arm64|aarch64) arch=arm64 ;; \
		x86_64) arch=amd64 ;; \
		*) echo "bootstrap-go: unsupported CPU ($$(uname -m))"; exit 1 ;; \
	esac; \
	pl="$$os-$$arch"; \
	url="https://go.dev/dl/go$$ver.$$pl.tar.gz"; \
	echo "bootstrap-go: fetching $$url"; \
	rm -rf "$$root/.tools/go"; \
	mkdir -p "$$root/.tools"; \
	curl -fSL "$$url" | tar -xzf - -C "$$root/.tools"; \
	test -x "$$root/.tools/go/bin/go"; \
	test -f "$$root/.tools/go/src/bytes/buffer.go"; \
	echo "bootstrap-go: installed at $(TOOLS_GO)"

setup: ensure-go
	@echo "Go toolchain: $(GO_CMD)"
	$(GO_CMD) install "$(GOBCE_PACKAGE)"
	@if [ -n "$${GITHUB_PATH:-}" ]; then \
		echo "$(GO_BIN)" >> "$$GITHUB_PATH"; \
	fi
	@echo "gobce is installed at $(GOBCE)."
	@echo "Tip: run 'make setup-hooks' to block git push when fmt-check or lint would fail CI."

setup-hooks:
	@git config core.hooksPath .githooks
	@echo "Git hooks path set to .githooks (pre-push runs: make ci-lint)."

build: ensure-go
	@mkdir -p bin
	$(GO_CMD) build -o "$(BINARY)" ./cmd/hado

gen-manifest-doc: ensure-go
	@mkdir -p bin
	$(GO_CMD) run ./cmd/hado manifest doc --out docs/hado.manifest.reference.yaml
	@echo "Wrote docs/hado.manifest.reference.yaml"

check-docker:
	@command -v docker >/dev/null 2>&1 || { echo "docker is required for YAML/Markdown lint."; exit 1; }

lint: lint-yaml lint-markdown lint-go

# Same as .github/workflows/lint.yml and .githooks/pre-push. No go test — Test workflow
# runs `go test` once inside `make readiness-check`.
ci-lint: fmt-check lint

lint-yaml: check-docker
	docker run --rm -v "$$(pwd):/work" -w /work python:3.12-alpine sh -lc "pip install --no-cache-dir yamllint >/dev/null && yamllint ."

lint-markdown: check-docker
	docker run --rm -v "$$(pwd):/work" -w /work "$(MARKDOWNLINT_CLI2_IMAGE)" "**/*.md" "#node_modules" "#.git" "#.tools"

lint-go: ensure-go
	$(GO_CMD) vet ./...

fmt: fmt-go

fmt-go: ensure-go
	@if [ -n "$(GO_FILES)" ]; then \
		set -e; $(GOFMT_CMD) -w $(GO_FILES); \
	fi

fmt-check: ensure-go
	@if [ -z "$(GO_FILES)" ]; then \
		exit 0; \
	fi; \
	set -e; unformatted="$$($(GOFMT_CMD) -l $(GO_FILES))"; \
	if [ -n "$$unformatted" ]; then \
		echo "These Go files need gofmt:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

test: ensure-go
	$(GO_CMD) test ./...

readiness-check: ensure-go
	@command -v "$(GOBCE)" >/dev/null 2>&1 || { echo "gobce is required. Run: make setup"; exit 1; }
	$(GO_CMD) test ./... -coverprofile="$(COVERPROFILE)"
	"$(GOBCE)" analyze --coverprofile "$(COVERPROFILE)" --format json --output "$(READINESS_COVERAGE)"
	$(GO_CMD) run ./cmd/hado fire --standard "$(READINESS_STANDARD)" --manifest "$(READINESS_MANIFEST)"

pre-pr: ci-lint test
	@echo "pre-pr: OK."
