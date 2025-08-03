package shared

import (
	"github.com/grafana/grafana-foundation-sdk/go/loki"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func PrometheusQuery(expression string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(expression).
		Range().
		Exemplar(false).
		Format(prometheus.PromQueryFormatTimeSeries).
		EditorMode(prometheus.QueryEditorModeCode).
		LegendFormat("__auto")
}

func LokiQuery(expression string) *loki.DataqueryBuilder {
	return loki.NewDataqueryBuilder().
		Expr(expression).
		EditorMode(loki.QueryEditorModeCode).
		QueryType("range")
}
