package shared

import (
	"fmt"
	"strings"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func LabelEqual(name, value string) string {
	return fmt.Sprintf("%s=%q", name, value)
}

func LabelMatch(name, value string) string {
	return fmt.Sprintf("%s=~%q", name, value)
}

// ...
func NewQueryLabelValuesVariable(
	name, label, description, metric, queryLabel string,
	filters ...string,
) *dashboard.QueryVariableBuilder {
	query := fmt.Sprintf("label_values(%s{%s},%s)", metric, strings.Join(filters, ","), queryLabel)
	return NewQueryVariable(name, label, description, query)
}

// ...
func NewQueryVariable(name, label, description, query string) *dashboard.QueryVariableBuilder {
	return dashboard.NewQueryVariableBuilder(name).
		Label(label).
		Description(description).
		Query(dashboard.StringOrMap{
			String: &query,
		}).
		IncludeAll(true).
		Current(AllVariableOption()).
		AllowCustomValue(false).
		Sort(dashboard.VariableSortAlphabeticalAsc)
}

// ...
func AllVariableOption() dashboard.VariableOption {
	return dashboard.VariableOption{
		Text: dashboard.StringOrArrayOfString{
			String: cog.ToPtr("All"),
		},
		Value: dashboard.StringOrArrayOfString{
			String: cog.ToPtr("$__all"),
		},
	}
}
