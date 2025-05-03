package rpcserver

import (
	"strings"

	"connectrpc.com/connect"
)

// SplitProcedure splits the RPC procedure into RPC service and RPC method components.
func SplitProcedure(procedure string) (service, method string) {
	name := strings.TrimLeft(procedure, "/")

	service, method, ok := strings.Cut(name, "/")
	if !ok {
		return "", service
	}

	return service, method
}

// ProtocolToSystem converts the request protocol to the format used by the OpenTelemetry semantic
// convention for RPC systems.
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
