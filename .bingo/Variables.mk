# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.10. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
BINGO_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Ensure bingo-managed tools are always built for the host platform,
# even when GOOS/GOARCH are set for cross-compilation of other targets.
GOHOSTOS     ?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH   ?= $(shell $(GO) env GOHOSTARCH)
GOHOSTARM    ?= $(shell $(GO) env GOHOSTARM)

# Below generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for actionlint variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(ACTIONLINT)
#	@echo "Running actionlint"
#	@$(ACTIONLINT) <flags/args..>
#
ACTIONLINT := $(GOBIN)/actionlint-v1.7.12
$(ACTIONLINT): $(BINGO_DIR)/actionlint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/actionlint-v1.7.12"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=actionlint.mod -o=$(GOBIN)/actionlint-v1.7.12 "github.com/rhysd/actionlint/cmd/actionlint"

BUF := $(GOBIN)/buf-v1.69.0
$(BUF): $(BINGO_DIR)/buf.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/buf-v1.69.0"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=buf.mod -o=$(GOBIN)/buf-v1.69.0 "github.com/bufbuild/buf/cmd/buf"

CTLPTL := $(GOBIN)/ctlptl-v0.9.2
$(CTLPTL): $(BINGO_DIR)/ctlptl.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/ctlptl-v0.9.2"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=ctlptl.mod -o=$(GOBIN)/ctlptl-v0.9.2 "github.com/tilt-dev/ctlptl/cmd/ctlptl"

GODOG := $(GOBIN)/godog-v0.15.1
$(GODOG): $(BINGO_DIR)/godog.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/godog-v0.15.1"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=godog.mod -o=$(GOBIN)/godog-v0.15.1 "github.com/cucumber/godog/cmd/godog"

GOLANGCI_LINT := $(GOBIN)/golangci-lint-v2.12.2
$(GOLANGCI_LINT): $(BINGO_DIR)/golangci-lint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golangci-lint-v2.12.2"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=golangci-lint.mod -o=$(GOBIN)/golangci-lint-v2.12.2 "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"

KIND := $(GOBIN)/kind-v0.31.0
$(KIND): $(BINGO_DIR)/kind.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/kind-v0.31.0"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=kind.mod -o=$(GOBIN)/kind-v0.31.0 "sigs.k8s.io/kind"

KUBE_SCORE := $(GOBIN)/kube-score-v1.20.0
$(KUBE_SCORE): $(BINGO_DIR)/kube-score.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/kube-score-v1.20.0"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=kube-score.mod -o=$(GOBIN)/kube-score-v1.20.0 "github.com/zegl/kube-score/cmd/kube-score"

KUBECONFORM := $(GOBIN)/kubeconform-v0.7.0
$(KUBECONFORM): $(BINGO_DIR)/kubeconform.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/kubeconform-v0.7.0"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=kubeconform.mod -o=$(GOBIN)/kubeconform-v0.7.0 "github.com/yannh/kubeconform/cmd/kubeconform"

KUSTOMIZE := $(GOBIN)/kustomize-v5.8.1
$(KUSTOMIZE): $(BINGO_DIR)/kustomize.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/kustomize-v5.8.1"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=kustomize.mod -o=$(GOBIN)/kustomize-v5.8.1 "sigs.k8s.io/kustomize/kustomize/v5"

MOCKERY := $(GOBIN)/mockery-v3.7.0
$(MOCKERY): $(BINGO_DIR)/mockery.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/mockery-v3.7.0"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=mockery.mod -o=$(GOBIN)/mockery-v3.7.0 "github.com/vektra/mockery/v3"

PROTOC_GEN_CONNECT_GO := $(GOBIN)/protoc-gen-connect-go-v1.19.1
$(PROTOC_GEN_CONNECT_GO): $(BINGO_DIR)/protoc-gen-connect-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-connect-go-v1.19.1"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=protoc-gen-connect-go.mod -o=$(GOBIN)/protoc-gen-connect-go-v1.19.1 "connectrpc.com/connect/cmd/protoc-gen-connect-go"

PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go-v1.36.11
$(PROTOC_GEN_GO): $(BINGO_DIR)/protoc-gen-go.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/protoc-gen-go-v1.36.11"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=protoc-gen-go.mod -o=$(GOBIN)/protoc-gen-go-v1.36.11 "google.golang.org/protobuf/cmd/protoc-gen-go"

TRIVY := $(GOBIN)/trivy-v0.67.2
$(TRIVY): $(BINGO_DIR)/trivy.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/trivy-v0.67.2"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) GOEXPERIMENT=jsonv2 $(GO) build -mod=mod -modfile=trivy.mod -o=$(GOBIN)/trivy-v0.67.2 "github.com/aquasecurity/trivy/cmd/trivy"

YAMLFMT := $(GOBIN)/yamlfmt-v0.21.0
$(YAMLFMT): $(BINGO_DIR)/yamlfmt.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/yamlfmt-v0.21.0"
	@cd $(BINGO_DIR) && GOWORK=off GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) GOARM=$(GOHOSTARM) $(GO) build -mod=mod -modfile=yamlfmt.mod -o=$(GOBIN)/yamlfmt-v0.21.0 "github.com/google/yamlfmt/cmd/yamlfmt"

