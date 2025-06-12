# https://just.systems/man/en/

import '.tools.just'

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
[group('development environment')]
dev-up:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        up \
        --detach \
        --build \
        --wait

# Exec into the development environment.
[group('development environment')]
dev-exec:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        exec dev \
        bash -l

# List containers in the development environment.
[group('development environment')]
dev-ps:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        ps \
        --all

# Stop the development environment.
[group('development environment')]
dev-down:
    docker compose \
        --file compose.yaml \
        --file compose-dev.yaml \
        down \
        --volumes \
        --remove-orphans

# Restart the development environment.
[group('development environment')]
dev-restart: dev-down dev-up

# Run all automated code modifications.
checks: tidy vendor gen fmt

# Tidy all dependencies.
[group('dependencies')]
tidy: tidy-buf tidy-go

# Tidy Protobuf dependencies.
[group('dependencies')]
tidy-buf: install-buf
    {{ buf }} dep prune
    {{ buf }} dep update

# Tidy Go dependencies
[group('dependencies')]
tidy-go:
    go mod tidy

# Vendor all dependencies.
[group('dependencies')]
vendor: vendor-go

# Vendor Go dependencies.
[group('dependencies')]
vendor-go:
    go mod vendor

# Run all formatters.
[group('formatters')]
fmt: fmt-buf fmt-go fmt-just fmt-yaml

# Run the Protobuf formatter.
[group('formatters')]
fmt-buf: install-buf
    {{ buf }} format --config buf.yaml --write

# Run the Go formatter.
[group('formatters')]
fmt-go: install-gofumpt install-gci
    {{ gofumpt }} -l -w .
    {{ gci }} write \
        --skip-vendor \
        --skip-generated \
        -s standard \
        -s default \
        -s localmodule \
        .

# Run the Justfile formatter.
[group('formatters')]
fmt-just:
    just --unstable --fmt

# Run the YAML formatter
[group('formatters')]
fmt-yaml: install-yamlfmt
    {{ yamlfmt }} .

# Run all code generators.
[group('generators')]
gen: gen-buf gen-go gen-just

# Run the Protobuf generator.
[group('generators')]
gen-buf: install-buf install-protoc-gen-connect-go install-protoc-gen-go
    {{ buf }} generate --clean --config buf.yaml

# Run the Go generators.
[group('generators')]
gen-go: gen-go-jet gen-go-mockery gen-go-wire

# Run the Go jet generator
[group('generators')]
gen-go-jet: install-jet
    {{ jet }} \
        -source=postgres \
        -host={{ db_host }} \
        -port={{ db_port }} \
        -user={{ db_user }} \
        -password={{ db_pass }} \
        -dbname=postgres \
        -path ./internal/adapters/datastore/models/

# Run the Go mockery generator.
[group('generators')]
gen-go-mockery: install-mockery
    rm -rf mocks/
    {{ mockery }}

# Run the Go wire generator.
[group('generators')]
gen-go-wire: install-wire
    {{ wire }} gen ./cmd/...

# Run the justfile generator.
[group('generators')]
gen-just:
    ./tools/regen-tools.sh > .tools.just

# Check for uncommitted changes.
[private]
dirty:
    git diff --exit-code

# Run all linters.
[group('linters')]
lint: lint-buf lint-go

# Run the Protobuf linter.
[group('linters')]
lint-buf: install-buf
    {{ buf }} lint --config buf.yaml

# Run the Go linter.
[group('linters')]
lint-go:
    golangci-lint run

# Run all linter fixers.
[group('linters')]
lint-fix:

# Run the Go linter fixer.
[group('linters')]
lint-fix-go:
    golangci-lint run --fix

# Run the Go unit tests.
[group('tests')]
unit timeout="30s":
    go test -timeout={{ timeout }} -count=1 -cover -coverprofile=cover.out ./internal/... ./pkg/...
    @go run ./tools/filter-coverage/ -output=unit.out cover.out
    @echo "Total coverage: `go tool cover -func=unit.out | tail -n 1 | awk '{print $3}'`"
    go tool cover -html unit.out -o unit.html

# Run the functional tests.
[group('tests')]
functional:
    go test ./tests/functional/ \
        --godog.strict \
        --test.v \
        --test.count=1

# Summarise functional test coverage.
[group('tests')]
functional-cover:
    go tool covdata percent -i=.covdata
    @go tool covdata textfmt -i=.covdata -o functional.out
    @go run ./tools/filter-coverage/ -output=functional2.out functional.out
    @echo "Total coverage: `go tool cover -func=functional2.out | tail -n 1 | awk '{print $3}'`"
    go tool cover -html functional2.out -o functional.html

# Delete functional test coverage artifacts.
[group('tests')]
functional-cover-clean:
    rm -f .covdata/cov*

# Scan the repository for issues.
[group('scanners')]
scan: scan-gitleaks scan-trivy scan-zizmor

# Scan the repository for secrets with Gitleaks.
[group('scanners')]
scan-gitleaks:
    gitleaks dir --verbose

# Scan the repository for issues using Trivy.
[group('scanners')]
scan-trivy:
    trivy fs .

# Scan actions and workflows using Zizmor.
[group('scanners')]
scan-zizmor:
    zizmor --persona pedantic .

[private]
scan-zizmor-ci:
    docker buildx build --tag zizmor:local ./.github/actions/zizmor/
    docker run \
        --rm \
        --workdir /github/workspace \
        -e "INPUT_INPUTS=." \
        -v ".":"/github/workspace" \
        zizmor:local \
        "--persona=pedantic"

# Build all binaries.
[group('builds')]
build:
    CGO_ENABLED=0 go build -trimpath -ldflags="-buildid= -s -w" -o ./dist/ ./cmd/...;

# Exec into the database.
[group('development environment')]
db-exec:
    PGPASSWORD={{ db_pass }} psql \
        --host {{ db_host }} \
        --username {{ db_user }}

# Insert sample data into the database.
[group('development environment')]
db-seed:
    PGPASSWORD={{ db_pass }} psql \
        --host {{ db_host }} \
        --username {{ db_user }} \
        --echo-all \
        --file ./tools/seed.sql

# Build all containers.
[group('builds')]
container-build go_build_args="": (container-build-rpc go_build_args)

# Build the example-rpc container.
[group('builds')]
container-build-rpc go_build_args="": (_container-build "example-rpc" go_build_args)

[private]
_container-build service go_build_args="":
    SOURCE_DATE_EPOCH=0 docker buildx build \
        --pull \
        --target runtime \
        --build-arg SERVICE={{ service }} \
        --build-arg SOURCE_DATE_EPOCH=0 \
        --build-arg GO_BUILD_ARGS="{{ go_build_args }}" \
        --tag {{ service }}:{{ now }} \
        --tag {{ service }}:local \
        .

# Scan all containers.
[group('scanners')]
container-scan: container-scan-rpc

# Scan the example-rpc container
[group('scanners')]
container-scan-rpc: (_container-scan "example-rpc")

[private]
_container-scan service:
    trivy image \
        --config trivy.yaml \
        --docker-host unix://{{ env('HOME') }}/.colima/default/docker.sock \
        {{ service }}:local
