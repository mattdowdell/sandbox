package rpcserver

import (
	"strings"
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
func ProtocolToSystem(protocol string) string {
	switch protocol {
	case "grpcweb":
		return "grpc_web"

	case "connect":
		return "connect_rpc"

	default:
		return protocol
	}
}
