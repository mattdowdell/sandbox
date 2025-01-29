# Observability

Metrics, traces and logs can be used to observe behaviour and diagnose issues
within the application. Metrics and traces are created using
[`opentelemetry-go`], whilst logs are produced using Go's [`log/slog`] package
and formatted as JSON.

[`opentelemetry-go`]: https://github.com/open-telemetry/opentelemetry-go
[`log/slog`]: https://pkg.go.dev/log/slog

## Visualisations

Metrics and traces can both be viewed using [Grafana]. During local development,
this is accessible at [`localhost:3000`].

For generally browsing metrics and traces, the "Explore" interface can be selected
in the left sidebar. This is useful for ad-hoc metrics queries and for browsing
traces. Metrics can be viewed using the VictoriaMetrics datasource, whilst
traces are available using the Jaeger datasource.

[`localhost:3000`]: http://localhost:3000/

### Dashboards

Dashboards are also available for viewing multiple metrics at once, and can be
found under the "Dashboards" interface in the left sidebar.

The "Go Runtime" dashboard provides a view of the behaviour of the Go runtime.
It tracks the number of goroutines vs the processor limit, i.e. the number of
parallel items that can be worked on. The number of goroutines will typically be
higher at any given time, whilst the processor limit is usually static.
Additionally, the dashboard tracks memory usage and allocations.

The "RPC" dashboard provides a view into the error rate, latency and traffic for
the application. The dropdowns can be used to drill down into the data for a
specific RPC service or method. This dashboard is aimed at measuring compliance
with SLOs which are stated at the top of the dashboard.

## Queries

Whilst VictoriaMetrics is very similar to Prometheus for querying metrics, there
are some differences. The notable differences are:

- Labels and metric names can contain `.`. For example, a label in Prometheus
  might be `rpc_method`, whilst it is `rpc.method` in VictoriaMetrics. This
  difference means that public dashboards will likely not work without some
  modifications.
- _TODO: more?_

See [MetricsQL] for further documentation on the query language.

[MetricsQL]: https://docs.victoriametrics.com/metricsql/

## Exporters

Metrics are exported to a Victoria Metrics instance accessible at
[`localhost:8428`]. Traces are exported to a Jaeger instance accessible at
[`localhost:16686`] In the majority of cases, these can be ignored in favour of
Grafana. However, they can be used to verify issues exporting metrics or
traces.

[`localhost:8428`]: http://localhost:8428
[`localhost:16686`]: http://localhost:16686/
