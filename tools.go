//go:build tools
// +build tools

package main

import (
	_ "connectrpc.com/connect/cmd/protoc-gen-connect-go"
	_ "github.com/daixiang0/gci"
	_ "github.com/go-jet/jet/v2/cmd/jet"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/vektra/mockery/v2"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "mvdan.cc/gofumpt"
)
