// +build tools

package main

import (
	_ "github.com/daixiang0/gci"
	_ "mvdan.cc/gofumpt"
	_ "github.com/vektra/mockery/v2"
	_ "connectrpc.com/connect/cmd/protoc-gen-connect-go"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "github.com/google/wire/cmd/wire"
)
