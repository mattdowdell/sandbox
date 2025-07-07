package connectrpc

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

const (
	// TODO: figure out how to vary this based on the SLO
	aboveLatencyThresholdStatQuery = `(
	sum(increase(
		rpc.server.duration_bucket{
			rpc.system="connect_rpc",
			service.name=~"$service_name",
			rpc.service=~"$rpc_service",
			rpc.method=~"$rpc_method",
			le="250"
		}[$__rate_interval]
	)) - sum(increase(
		rpc.server.duration_bucket{
			rpc.system="connect_rpc",
			service.name=~"$service_name",
			rpc.service=~"$rpc_service",
			rpc.method=~"$rpc_method",
			le="100"
		}[$__rate_interval]
	))
) / sum(increase(
	rpc.server.duration_count{
		rpc.system="connect_rpc",
		service.name=~"$service_name",
		rpc.service=~"$rpc_service",
		rpc.method=~"$rpc_method"
	}[$__rate_interval]
))`

	aboveLatencyThresholdTimeseriesQuery = `(
	sum by(service.name,rpc.service,rpc.method) (increase(
		rpc.server.duration_bucket{
			rpc.system="connect_rpc",
			service.name=~"$service_name",
			rpc.service=~"$rpc_service",
			rpc.method=~"$rpc_method",
			le="250"
		}[$__rate_interval]
	)) - sum by(service.name,rpc.service,rpc.method) (increase(
		rpc.server.duration_bucket{
			rpc.system="connect_rpc",
			service.name=~"$service_name",
			rpc.service=~"$rpc_service",
			rpc.method=~"$rpc_method",
			le="100"
		}[$__rate_interval]
	))
) / sum by(service.name,rpc.service,rpc.method) (increase(
	rpc.server.duration_count{
		rpc.system="connect_rpc",
		service.name=~"$service_name",
		rpc.service=~"$rpc_service",
		rpc.method=~"$rpc_method"
	}[$__rate_interval]
))`

	// TODO: make error codes a constant somewhere?
	errorResponsesStatQuery = `sum(increase(
	rpc.server.duration_count{
		rpc.system="connect_rpc",
		service.name=~"$service_name",
		rpc.service=~"$rpc_service",
		rpc.method=~"$rpc_method",
		rpc.connect_rpc.error_code=~"unknown|deadline_exceeded|unimplemented|internal|unavailable|data_loss"
	}[$__rate_interval]
)) / sum(increase(
	rpc.server.duration_count{
		rpc.system="connect_rpc",
		service.name=~"$service_name",
		rpc.service=~"$rpc_service",
		rpc.method=~"$rpc_method"
	}[$__rate_interval]
))`

	errorResponsesTimeseriesQuery = `sum by(service.name,rpc.service,rpc.method) (increase(
	rpc.server.duration_count{
		rpc.system="connect_rpc",
		service.name=~"$service_name",
		rpc.service=~"$rpc_service",
		rpc.method=~"$rpc_method",
		rpc.connect_rpc.error_code=~"unknown|deadline_exceeded|unimplemented|internal|unavailable|data_loss"
	}[$__rate_interval]
)) / sum by(service.name,rpc.service,rpc.method) (increase(
	rpc.server.duration_count{
		rpc.system="connect_rpc",
		service.name=~"$service_name",
		rpc.service=~"$rpc_service",
		rpc.method=~"$rpc_method"
	}[$__rate_interval]
))`
)

// ...
func AboveLatencyThresholdStatPanel() *stat.PanelBuilder {
	return newStatPanel(
		"Above Latency Threshold",
		"The percentage of requests that exceeded the latency threshold across the current time range.",
		aboveLatencyThresholdStatQuery,
	)
}

// ...
func ErrorResponsesStatPanel() *stat.PanelBuilder {
	return newStatPanel(
		"Error Responses",
		"The percentage of requests that resulted in error responses across the current time range.",
		errorResponsesStatQuery,
	)
}

func AboveLatencyThresholdTimeseriesPanel() *timeseries.PanelBuilder {
	return newTimeseriesPanel(
		"Above Latency Threshold",
		"The percentage of requests that exceeded the latency threshold.",
		aboveLatencyThresholdTimeseriesQuery,
	)
}

func ErrorResponsesTimeseriesPanel() *timeseries.PanelBuilder {
	return newTimeseriesPanel(
		"Error Responses",
		"The percentage of requests that resulted in error responses.",
		errorResponsesTimeseriesQuery,
	)
}

func newStatPanel(title, description, query string) *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title(title).
		Description(description).
		Height(8).
		Span(4).
		Unit(units.PercentUnit).
		Min(0).
		Max(1).
		Decimals(2). // TODO: make this longer depending on the SLO
		NoValue("N/A").
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().Steps([]dashboard.Threshold{
				{
					Color: "green",
				},
				{
					Value: cog.ToPtr(0.01), // TODO: vary this based on the SLO
					Color: "red",
				},
			}),
		).
		ReduceOptions(common.NewReduceDataOptionsBuilder().Calcs([]string{"mean"})).
		WithTarget(
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("Requests"),
		)
}

func newTimeseriesPanel(title, description, query string) *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title(title).
		Description(description).
		AxisSoftMax(0.1).
		Height(8).
		Span(8).
		Unit(units.PercentUnit).
		Min(0).
		Max(1).
		WithTarget(
			prometheus.NewDataqueryBuilder().
				Expr(query).
				LegendFormat("{{service.name}} {{rpc.service}}/{{rpc.method}}"),
		)
}
