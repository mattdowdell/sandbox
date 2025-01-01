# https://just.systems/man/en/

[private]
default:
	@just --list

# Tidy all dependencies.
tidy: tidy-buf tidy-go

# Tidy Protobuf dependencies.
tidy-buf:
	buf dep update

# Tidy Go dependencies
tidy-go:
	go mod tidy

# Vendor Go dependencies.
vendor:
	go mod vendor

# Run all formatters.
fmt: fmt-buf fmt-go

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

# Run all formatter checks.
fmt-check: fmt-check-buf fmt-check-go

# Run the Protobuf formatter check.
fmt-check-buf:
	buf format --config buf.yaml --diff --exit-code

# Run the Go formatter check (unimplemented).
fmt-check-go:
	@echo '{{ style("error") }}unimplemented{{ NORMAL }}' && exit 1

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

# Run all code generators.
gen: gen-buf gen-go

# Run the Protobuf generator.
gen-buf:
	buf generate --clean --config buf.yaml

# Run the Go generators.
gen-go: gen-go-wire

# Run the Go wire generator.
gen-go-wire:
	wire gen ./cmd/...

# Run the Go mockery generator.
gen-go-mockery:
	mockery

# Run the Go unit tests.
unit:
	go test -count=1 -cover ./...

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
