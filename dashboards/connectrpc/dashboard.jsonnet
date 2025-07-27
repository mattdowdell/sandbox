local g = import '../g.libsonnet';

local row = g.panel.row;
local variables = import './variables.libsonnet';
local panels = import './panels.libsonnet';

g.dashboard.new('ConnectRPC (generated)')
+ g.dashboard.withTimezone('browser')
+ g.dashboard.time.withFrom('now-1h')
+ g.dashboard.time.withTo('now')
+ g.dashboard.withRefresh('1m')
+ g.dashboard.withVariables([
  variables.serviceName,
  variables.rpcService,
  variables.rpcMethod,
])
+ g.dashboard.withPanels([
  row.new('Latency')
  + row.withCollapsed(false),
  panels.stat.small(
    'P99 Average (Success)',
    'The P99 latency for success responses across the current time range.',
  ),
  panels.stat.small(
    'P99 Average (Error)',
    'The P99 latency for error responses across the current time range.',
  ),
  panels.stat.small(
    'P90 Average (Success)',
    'The P90 latency for success responses across the current time range.',
  ),
  panels.stat.small(
    'P90 Average (Error)',
    'The P90 latency for error responses across the current time range.',
  ),
  panels.stat.small(
    'P50 Average (Success)',
    'The P50 latency for success responses across the current time range.',
  ),
  panels.stat.small(
    'P50 Average (Error)',
    'The P50 latency for error responses across the current time range.',
  ),
])
