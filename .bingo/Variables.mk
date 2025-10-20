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
BUF := $(GOBIN)/buf-v1.58.0
$(BUF): $(BINGO_DIR)/buf.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/buf-v1.58.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=buf.mod -o=$(GOBIN)/buf-v1.58.0 "github.com/bufbuild/buf/cmd/buf"

CTLPTL := $(GOBIN)/ctlptl-v0.8.43
$(CTLPTL): $(BINGO_DIR)/ctlptl.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/ctlptl-v0.8.43"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=ctlptl.mod -o=$(GOBIN)/ctlptl-v0.8.43 "github.com/tilt-dev/ctlptl/cmd/ctlptl"

GODOG := $(GOBIN)/godog-v0.15.1
$(GODOG): $(BINGO_DIR)/godog.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/godog-v0.15.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=godog.mod -o=$(GOBIN)/godog-v0.15.1 "github.com/cucumber/godog/cmd/godog"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v2.5.0
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v2.5.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v2.5.0 "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"

KIND := $(GOBIN)/kind-v0.30.0
$(KIND): $(BINGO_DIR)/kind.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/kind-v0.30.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=kind.mod -o=$(GOBIN)/kind-v0.30.0 "sigs.k8s.io/kind"

KUBE_SCORE := $(GOBIN)/kube-score-v1.20.0
$(KUBE_SCORE): $(BINGO_DIR)/kube-score.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/kube-score-v1.20.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=kube-score.mod -o=$(GOBIN)/kube-score-v1.20.0 "github.com/zegl/kube-score/cmd/kube-score"

KUBECONFORM := $(GOBIN)/kubeconform-v0.7.0
$(KUBECONFORM): $(BINGO_DIR)/kubeconform.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/kubeconform-v0.7.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=kubeconform.mod -o=$(GOBIN)/kubeconform-v0.7.0 "github.com/yannh/kubeconform/cmd/kubeconform"

MOCKERY := $(GOBIN)/mockery-v3.5.5
$(MOCKERY): $(BINGO_DIR)/mockery.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/mockery-v3.5.5"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=mockery.mod -o=$(GOBIN)/mockery-v3.5.5 "github.com/vektra/mockery/v3"

PROTOC_GEN_CONNECT_GO := $(GOBIN)/protoc-gen-connect-go-v1.19.1
$(PROTOC_GEN_CONNECT_GO): $(BINGO_DIR)/protoc-gen-connect-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-connect-go-v1.19.1"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-connect-go.mod -o=$(GOBIN)/protoc-gen-connect-go-v1.19.1 "connectrpc.com/connect/cmd/protoc-gen-connect-go"

PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go-v1.36.10
$(PROTOC_GEN_GO): $(BINGO_DIR)/protoc-gen-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-v1.36.10"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=protoc-gen-go.mod -o=$(GOBIN)/protoc-gen-go-v1.36.10 "google.golang.org/protobuf/cmd/protoc-gen-go"

YAMLFMT := $(GOBIN)/yamlfmt-v0.19.0
$(YAMLFMT): $(BINGO_DIR)/yamlfmt.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/yamlfmt-v0.19.0"
	@cd $(BINGO_DIR) && GOWORK=off $(GO) build -mod=mod -modfile=yamlfmt.mod -o=$(GOBIN)/yamlfmt-v0.19.0 "github.com/google/yamlfmt/cmd/yamlfmt"

