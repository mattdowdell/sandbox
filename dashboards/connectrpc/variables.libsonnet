local g = import '../g.libsonnet';
local var = g.dashboard.variable;

{
  serviceName:
    var.query.new('service_name')
    + var.query.generalOptions.withLabel('Service Name')
    + var.query.generalOptions.withDescription('Filter panels by service name.')
    + var.query.refresh.onTime()
    + var.query.selectionOptions.withIncludeAll(true)
    + var.query.withSort(1) // alphabetical asc
    + var.query.queryTypes.withLabelValues(
      'service.name',
      'rpc.server.duration_bucket{rpc.system="connect_rpc"}',
    )
    + var.query.generalOptions.withCurrent('All', '$__all'),

  rpcService:
    var.query.new('rpc_service')
    + var.query.generalOptions.withLabel('RPC Service')
    + var.query.generalOptions.withDescription('Filter panels by RPC service.')
    + var.query.refresh.onTime()
    + var.query.selectionOptions.withIncludeAll(true)
    + var.query.withSort(1) // alphabetical asc
    + var.query.generalOptions.withCurrent('All', '$__all')
    + var.query.queryTypes.withLabelValues(
      'rpc.service',
      'rpc.server.duration_bucket{rpc.system="connect_rpc",service.name=~"$%s"}'
      % [self.serviceName.name],
    ),

  rpcMethod:
    var.query.new('rpc_method')
    + var.query.generalOptions.withLabel('RPC Method')
    + var.query.generalOptions.withDescription('Filter panels by RPC method.')
    + var.query.refresh.onTime()
    + var.query.selectionOptions.withIncludeAll(true)
    + var.query.withSort(1) // alphabetical asc
    + var.query.generalOptions.withCurrent('All', '$__all')
    + var.query.queryTypes.withLabelValues(
      'rpc.method',
      'rpc.server.duration_bucket{rpc.system="connect_rpc",service.name=~"$%s",rpc.service=~"$%s"}'
      % [self.serviceName.name, self.rpcService.name],
    ),
}

/*
"templating": {
    "list": [
      {
        "current": {
          "text": "All",
          "value": "$__all"
        },
        //"definition": "label_values(rpc.server.duration_bucket{rpc.system=\"connect_rpc\"},service.name)",
        //"description": "Filter panels by service name.",
        //"includeAll": true,
        //"label": "Service Name",
        //"name": "service_name",
        "options": [],
        "query": {
          "qryType": 1,
          "query": "label_values(rpc.server.duration_bucket{rpc.system=\"connect_rpc\"},service.name)",
          "refId": "VariableQueryEditor-VariableQuery"
        },
        //"refresh": 1,
        "regex": "",
        //"sort": 1,
        //"type": "query"
      },
      {
        "current": {
          "text": "All",
          "value": "$__all"
        },
        "definition": "label_values(rpc.server.duration_bucket{rpc.system=\"connect_rpc\",service.name=~\"$service_name\"},rpc.service)",
        "description": "Filter panels by RPC service.",
        "includeAll": true,
        "label": "RPC Service",
        "name": "rpc_service",
        "options": [],
        "query": {
          "qryType": 1,
          "query": "label_values(rpc.server.duration_bucket{rpc.system=\"connect_rpc\",service.name=~\"$service_name\"},rpc.service)",
          "refId": "VariableQueryEditor-VariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "sort": 1,
        "type": "query"
      },
      {
        "current": {
          "text": "All",
          "value": "$__all"
        },
        "definition": "label_values(rpc.server.duration_bucket{rpc.system=\"connect_rpc\",service.name=~\"$service_name\", rpc.service=~\"$rpc_service\"},rpc.method)",
        "description": "Filter panels by RPC method.",
        "includeAll": true,
        "label": "RPC Method",
        "name": "rpc_method",
        "options": [],
        "query": {
          "qryType": 1,
          "query": "label_values(rpc.server.duration_bucket{rpc.system=\"connect_rpc\",service.name=~\"$service_name\", rpc.service=~\"$rpc_service\"},rpc.method)",
          "refId": "VariableQueryEditor-VariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "sort": 1,
        "type": "query"
      }
    ]
*/
