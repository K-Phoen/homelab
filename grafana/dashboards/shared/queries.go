package shared

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/loki"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
)

func PrometheusQuery(expressionFormat string, args ...any) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(fmt.Sprintf(expressionFormat, args...)).
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
