# https://just.systems/man/en/

is_docker := path_exists("/.dockerenv")
db_host := if is_docker == "true" { "postgres" } else { "localhost" }
db_port := "5432"
db_user := "postgres"
db_pass := "secret"
export BUILDKIT_PROGRESS := "plain"
now := shell("date +%s")

[private]
default:
    @just --list

# Start the development environment.
dev-up:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        up \
        --detach \
        --build \
        --wait

# Exec into the development environment.
dev-exec:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        exec dev \
        bash -l

# List containers in the development environment.
dev-ps:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        ps \
        --all

# Stop the development environment.
dev-down:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        down \
        --volumes \
        --remove-orphans

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

# Install go tools.
install-tools:
    cat tools.go | grep _ | awk -F'"' '{print $2}' | xargs -tI % go install %
    @# TODO: provide via tools.go; currently depends on outdated protovalidate-go
    go install github.com/bufbuild/buf/cmd/buf@latest

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
gen-go: gen-go-jet gen-go-mockery gen-go-wire

# Run the Go jet generator
gen-go-jet:
    jet \
        -source=postgres \
        -host={{ db_host }} \
        -port={{ db_port }} \
        -user={{ db_user }} \
        -password={{ db_pass }} \
        -dbname=postgres \
        -path ./internal/adapters/datastore/models/

# Run the Go mockery generator.
gen-go-mockery:
    rm -rf mocks/
    mockery

# Run the Go wire generator.
gen-go-wire:
    wire gen ./cmd/...

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
unit timeout="30s":
    go test -timeout={{ timeout }} -count=1 -cover -coverprofile=cover.out ./internal/... ./pkg/...
    @echo "Total coverage: `go tool cover -func=cover.out | tail -n 1 | awk '{print $3}'`"
    go tool cover -html cover.out -o cover.html

# Scan the repository for issues.
scan: scan-gitleaks scan-trivy scan-zizmor

# Scan the repository for secrets with Gitleaks.
scan-gitleaks:
    gitleaks dir

# Scan the repository for issues using Trivy.
scan-trivy:
    trivy fs .

# Scan actions and workflows using Zizmor.
scan-zizmor:
    zizmor --persona pedantic .github/actions/*/*.yml .github/workflows/*.yml

# Build all binaries.
build:
    CGO_ENABLED=0 go build -trimpath -ldflags="-buildid= -s -w" -o ./dist/ ./cmd/...;

# Exec into the database.
db-exec:
    PGPASSWORD={{ db_pass }} psql \
        --host {{ db_host }} \
        --username {{ db_user }}

# Insert sample data into the database.
db-seed:
    PGPASSWORD={{ db_pass }} psql \
        --host {{ db_host }} \
        --username {{ db_user }} \
        --echo-all \
        --file ./tools/seed.sql

# Build all containers.
container-build: container-build-rpc

# Build the example-rpc container.
container-build-rpc: (_container-build "example-rpc")

[private]
_container-build service:
    SOURCE_DATE_EPOCH=0 docker buildx build \
        --pull \
        --no-cache \
        --target runtime \
        --build-arg SERVICE={{ service }} \
        --build-arg SOURCE_DATE_EPOCH=0 \
        --tag {{ service }}:{{ now }} \
        --tag {{ service }}:local \
        .

# Scan all containers.
container-scan: container-scan-rpc

# Scan the example-rpc container
container-scan-rpc: (_container-scan "example-rpc")

[private]
_container-scan service:
    trivy image \
        --config trivy.yaml \
        --docker-host unix://{{ env('HOME') }}/.colima/default/docker.sock \
        {{ service }}:local
