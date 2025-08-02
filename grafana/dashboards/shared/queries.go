package shared

import "github.com/grafana/grafana-foundation-sdk/go/prometheus"

func PrometheusQuery(expression string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(expression).
		Range().
		Format(prometheus.PromQueryFormatTimeSeries).
		EditorMode(prometheus.QueryEditorModeCode).
		LegendFormat("__auto")
}
