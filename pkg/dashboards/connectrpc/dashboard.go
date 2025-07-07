package connectrpc

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func NewDashboard() (dashboard.Dashboard, error) {
	return dashboard.NewDashboardBuilder("ConnectRPC (generated)").
		Uid("my-connect-rpc"). // inherit the uid from the output file?
		Refresh("1m").         // TODO: make configurable
		Time("now-1h", "now"). // TODO: make configurable
		Timezone(common.TimeZoneBrowser).
		Variables([]cog.Builder[dashboard.VariableModel]{
			ServiceNameVariable(),
			RPCServiceVariable(),
			RPCMethodVariable(),
		}).
		WithRow(dashboard.NewRowBuilder("Errors")).
		WithPanel(AboveLatencyThresholdStatPanel()).
		WithPanel(ErrorResponsesStatPanel()).
		WithPanel(AboveLatencyThresholdTimeseriesPanel()).
		WithPanel(ErrorResponsesTimeseriesPanel()).
		WithRow(dashboard.NewRowBuilder("Latency")).
		WithRow(dashboard.NewRowBuilder("Traffic")).
		Build()
}
