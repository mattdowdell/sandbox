# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.9. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for buf variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(BUF)
#	@echo "Running buf"
#	@$(BUF) <flags/args..>
#
BUF := $(GOBIN)/buf-v1.55.1
$(BUF): $(BINGO_DIR)/buf.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/buf-v1.55.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=buf.mod -o=$(GOBIN)/buf-v1.55.1 "github.com/bufbuild/buf/cmd/buf"

GODOG := $(GOBIN)/godog-v0.15.1
$(GODOG): $(BINGO_DIR)/godog.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/godog-v0.15.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=godog.mod -o=$(GOBIN)/godog-v0.15.1 "github.com/cucumber/godog/cmd/godog"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v2.3.0
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v2.3.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v2.3.0 "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"

JB := $(GOBIN)/jb-v0.6.0
$(JB): $(BINGO_DIR)/jb.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/jb-v0.6.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=jb.mod -o=$(GOBIN)/jb-v0.6.0 "github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb"

JSONNET := $(GOBIN)/jsonnet-v0.21.0
$(JSONNET): $(BINGO_DIR)/jsonnet.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/jsonnet-v0.21.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=jsonnet.mod -o=$(GOBIN)/jsonnet-v0.21.0 "github.com/google/go-jsonnet/cmd/jsonnet"

MOCKERY := $(GOBIN)/mockery-v3.5.1
$(MOCKERY): $(BINGO_DIR)/mockery.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/mockery-v3.5.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=mockery.mod -o=$(GOBIN)/mockery-v3.5.1 "github.com/vektra/mockery/v3"

PROTOC_GEN_CONNECT_GO := $(GOBIN)/protoc-gen-connect-go-v1.18.1
$(PROTOC_GEN_CONNECT_GO): $(BINGO_DIR)/protoc-gen-connect-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-connect-go-v1.18.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-connect-go.mod -o=$(GOBIN)/protoc-gen-connect-go-v1.18.1 "connectrpc.com/connect/cmd/protoc-gen-connect-go"

PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go-v1.36.6
$(PROTOC_GEN_GO): $(BINGO_DIR)/protoc-gen-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-v1.36.6"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-go.mod -o=$(GOBIN)/protoc-gen-go-v1.36.6 "google.golang.org/protobuf/cmd/protoc-gen-go"

WIRE := $(GOBIN)/wire-v0.6.0
$(WIRE): $(BINGO_DIR)/wire.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/wire-v0.6.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=wire.mod -o=$(GOBIN)/wire-v0.6.0 "github.com/google/wire/cmd/wire"

YAMLFMT := $(GOBIN)/yamlfmt-v0.17.2
$(YAMLFMT): $(BINGO_DIR)/yamlfmt.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/yamlfmt-v0.17.2"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=yamlfmt.mod -o=$(GOBIN)/yamlfmt-v0.17.2 "github.com/google/yamlfmt/cmd/yamlfmt"

