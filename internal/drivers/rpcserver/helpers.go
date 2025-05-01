package rpcserver

import (
	"strings"

	"connectrpc.com/connect"
)

// ...
func SplitProcedure(spec connect.Spec) (service, method string) {
	name := strings.TrimLeft(spec.Procedure, "/")

	service, method, ok := strings.Cut(name, "/")
	if !ok {
		return "", service
	}

	return service, method
}

// ...
func ProtocolToSystem(peer connect.Peer) string {
	switch peer.Protocol {
	case "grpcweb":
		return "grpc_web"

	case "connect":
		return "connect_rpc"

	default:
		return peer.Protocol
	}
}
