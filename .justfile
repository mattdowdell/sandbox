# https://just.systems/man/en/

[private]
default:
    @just --list

# Start the development environment.
dev-up:
    docker compose up --detach --build

# Exec into the development environment.
dev-exec:
    docker compose exec dev bash -l

# Stop the development environment.
dev-down:
    docker compose down -v

# Restart the development environment.
dev-restart: dev-down dev-up

# Run all automated code modifications.
checks: tidy vendor gen fmt

# Tidy all dependencies.
tidy: tidy-buf tidy-go

# Tidy Protobuf dependencies.
tidy-buf:
    buf dep prune
    buf dep update

# Tidy Go dependencies
tidy-go:
    go mod tidy

# Vendor all dependencies.
vendor: vendor-go

# Vendor Go dependencies.
vendor-go:
    go mod vendor

# Run all formatters.
fmt: fmt-buf fmt-go fmt-just

# Run the Protobuf formatter.
fmt-buf:
    buf format --config buf.yaml --write

# Run the Go formatter.
fmt-go:
    gofumpt -l -w .
    gci write \
        --skip-vendor \
        --skip-generated \
        -s standard \
        -s default \
        -s localmodule \
        .

# Run the Justfile formatter.
fmt-just:
    just --unstable --fmt

# Run all code generators.
gen: gen-buf gen-go

# Run the Protobuf generator.
gen-buf:
    buf generate --clean --config buf.yaml

# Run the Go generators.
gen-go: gen-go-wire gen-go-mockery

# Run the Go wire generator.
gen-go-wire:
    wire gen ./cmd/...

# Run the Go mockery generator.
gen-go-mockery:
    rm -rf mocks/
    mockery

# Check for uncommitted changes.
[private]
dirty:
    git diff --exit-code

# Run all linters.
lint: lint-buf lint-go

# Run the Protobuf linter.
lint-buf:
    buf lint --config buf.yaml

# Run the Go linter.
lint-go:
    golangci-lint run

# Run all linter fixers.
lint-fix:

# Run the Go linter fixer.
lint-fix-go:
    golangci-lint run --fix

# Run the Go unit tests.
unit:
    go test -count=1 -cover ./...

# Scan the repository for issues.
scan: scan-trivy

# Scan the repository for issues using Trivy.
scan-trivy:
    trivy fs --config trivy.yaml .

# Build all containers.
container-build: container-build-rpc

# Build the example-rpc container.
container-build-rpc: (_container-build "example-rpc")

[private]
_container-build service:
    podman build \
        --target runtime \
        --build-arg SERVICE={{ service }} \
        --tag {{ service }}:local \
        .
