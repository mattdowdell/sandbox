package connectrpc

import (
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"

	"github.com/mattdowdell/sandbox/pkg/dashboards/shared"
)

// ...
func ServiceNameVariable() *dashboard.QueryVariableBuilder {
	return shared.NewQueryLabelValuesVariable(
		"service_name",                   /*name*/
		"Service Name",                   /*label*/
		"Filter panels by service name.", /*description*/
		"rpc.server.duration_bucket",     /*metrics*/
		"service.name",                   /*queryLabel*/
		shared.LabelEqual("rpc.system", "connect_rpc"),
	)
}

// ...
func RPCServiceVariable() *dashboard.QueryVariableBuilder {
	return shared.NewQueryLabelValuesVariable(
		"rpc_service",                   /*name*/
		"RPC Service",                   /*label*/
		"Filter panels by RPC service.", /*description*/
		"rpc.server.duration_bucket",    /*metrics*/
		"rpc.service",                   /*queryLabel*/
		shared.LabelEqual("rpc.system", "connect_rpc"),
		shared.LabelMatch("service.name", "$service_name"),
	)
}

// ...
func RPCMethodVariable() *dashboard.QueryVariableBuilder {
	return shared.NewQueryLabelValuesVariable(
		"rpc_method",                   /*name*/
		"RPC Method",                   /*label*/
		"Filter panels by RPC method.", /*description*/
		"rpc.server.duration_bucket",   /*metrics*/
		"rpc.method",                   /*queryLabel*/
		shared.LabelEqual("rpc.system", "connect_rpc"),
		shared.LabelMatch("service.name", "$service_name"),
		shared.LabelMatch("rpc.service", "$rpc_service"),
	)
}
