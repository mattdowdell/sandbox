# Development

## Layout

_TODO: describe the repository layout and a quick clean architecture overview_

## Development Environment

_TODO: describe the development environment and how to manage it._

## Go

_TODO: Discuss the use of Go, package selection philosophy, and useful just recipes._

## CI

_TODO: discuss the purpose of CI and what it includes_

## API

_TODO: Discuss where the RPC API is defined, how to use it. Also cover validation._

## Database

_TODO: discuss exec'ing into the db, seeding it, etc._
_TODO: discuss primary key selection + other standards._
_TODO: move standards discussion to a separate doc._

## Observability

_TODO: merge the intro, dashboards and explore sections with the (outdated) observability doc._
_TODO: let this be a guide for using observability, make the observability doc a guide for adding it._

The development environment includes an observability stack, including Grafana (dashboard),
VictoriaMetrics (metrics) and Jaeger (tracing). Grafana can be accessed at [`localhost:3000`].
VictoriaMetrics and Jaeger can be accessed at [`localhost:8428`] and [`localhost:16686`]
respectively, although there is seldom a need to access them directly.

[`localhost:3000`]: http://localhost:3000
[`localhost:8428`]: http://localhost:8428
[`localhost:16686`]: http://localhost:16686

### Dashboards

A few dashboards are pre-installed.

- ConnectRPC: Monitor the ConnectRPC requests to the RPC server.
- gRPC: Monitor the gRPC requests to the RPC server.
- Database: Monitor the Database queries made by the RPC server.
- Go Runtime: Monitor the Go runtime for the RPC server.

Most of the dashboards focus on the 4 Golden Signals and aim to provide a quick view of each
components health for a given time range.

New dashboards can be manually created/updated, exported and saved in the [`config/dashboards/`]
directory. These dashboards will automatically be imported into Grafana when it starts.

[`config/dashboards/`]: ../config/dashboards/

### Explore

The explore view allows metrics beyond those in dashboard to be viewed. It also allows inspection of
traces.

<!--

TODO: make observability doc into a standard for adding o11y, let this be a usage guide.

### Logs

_Discuss guidelines for logs._

### Metrics

_Discuss guidelines for metrics._

### Traces

_Discuss guidelines for traces._

-->
