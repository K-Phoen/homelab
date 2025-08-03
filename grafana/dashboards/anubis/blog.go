package anubis

import (
	"github.com/K-Phoen/homelab/grafana/dashboards/shared"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/piechart"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func requestRatesTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Request rates").
		WithTarget(shared.PrometheusQuery("sum(rate(anubis_policy_results{instance=~\"$instance\",container=\"anubis\",job=\"integrations/anubis-blog\"}[$__rate_interval])) by (action)").
			LegendFormat("{{action}}"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Transformations([]dashboard.DataTransformerConfig{
			{Id: "renameByRegex", Options: map[string]any{"regex": "ALLOW", "renamePattern": "Allowed"}},
			{Id: "renameByRegex", Options: map[string]any{"renamePattern": "Challenged", "regex": "CHALLENGE"}},
			{Id: "renameByRegex", Options: map[string]any{"renamePattern": "Denied", "regex": "DENY"}}},
		).
		Unit(units.RequestsPerSecond).
		Min(0).
		GradientMode(common.GraphGradientModeNone).
		ThresholdsStyle(common.NewGraphThresholdsStyleConfigBuilder().Mode(common.GraphThresholdsStyleModeOff)).
		FillOpacity(0).
		ShowPoints(common.VisibilityModeAuto).
		PointSize(5)
}

func verdictsPieChart() *piechart.PanelBuilder {
	return shared.PieChartPanel("Verdicts").
		WithTarget(shared.PrometheusQuery("sum by(action) (increase(anubis_policy_results{instance=~\"$instance\",container=\"anubis\", job=\"integrations/anubis-blog\"}[$__range]))").
			LegendFormat("{{action}}"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Transformations([]dashboard.DataTransformerConfig{
			{Id: "renameByRegex", Options: map[string]any{"regex": "CHALLENGE", "renamePattern": "Challenge"}},
			{Id: "renameByRegex", Options: map[string]any{"regex": "ALLOW", "renamePattern": "Allow"}},
			{Id: "renameByRegex", Options: map[string]any{"regex": "DENY", "renamePattern": "Deny"}}},
		).
		Unit(units.NoUnit).
		Min(0).
		ColorScheme(dashboard.NewFieldColorBuilder().Mode("palette-classic")).
		OverrideByName("Challenge", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "light-orange", "mode": "fixed"}},
		}).
		OverrideByName("Deny", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "red", "mode": "fixed"}}},
		).
		OverrideByName("Allow", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "green", "mode": "fixed"}}},
		)
}

func challengeResultsPieChart() *piechart.PanelBuilder {
	return shared.PieChartPanel("Challenge results").
		WithTarget(shared.PrometheusQuery("increase(anubis_challenges_validated{instance=~\"$instance\",container=\"anubis\",job=\"integrations/anubis-blog\"}[$__range])").
			LegendFormat("Validated"),
		).
		WithTarget(shared.PrometheusQuery("increase(anubis_failed_validations{instance=~\"$instance\",container=\"anubis\",job=\"integrations/anubis-blog\"}[$__range])").
			LegendFormat("Failed"),
		).
		WithTarget(shared.PrometheusQuery("sum by(job, instance) (increase(anubis_challenges_issued{job=\"integrations/anubis-blog\",container=\"anubis\",instance=~\"$instance\"}[$__range])) -  increase(anubis_challenges_validated{job=\"integrations/anubis-blog\",instance=~\"$instance\"}[$__range]) - increase(anubis_failed_validations{job=\"integrations/anubis-blog\",instance=~\"$instance\"}[$__range])").
			LegendFormat("No response"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.NoUnit).
		OverrideByName("No response", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "light-orange", "mode": "fixed"}}},
		).
		OverrideByName("Validated", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"mode": "fixed", "fixedColor": "green"}}},
		).
		OverrideByName("Failed", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "red", "mode": "fixed"}}},
		)
}

func allowedRequestsLogs() *logs.PanelBuilder {
	return shared.LogPanel("Allowed requests").
		WithTarget(shared.LokiQuery("{cluster=~\"homelab\", namespace=\"blog\", container=\"blog\"}")).
		Datasource(shared.DefaultLokiDatasource())
}

func logsVolumeTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Logs volume").
		WithTarget(shared.LokiQuery("sum (count_over_time({cluster=~\"homelab\", namespace=\"blog\", container=\"blog\"} [$__auto]))").
			LegendFormat("Logs"),
		).
		Datasource(shared.DefaultLokiDatasource()).
		DrawStyle(common.GraphDrawStyleBars).
		GradientMode(common.GraphGradientModeNone).
		ShowPoints(common.VisibilityModeNever).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false))
}

func BlogDashboard() *dashboard.DashboardBuilder {
	return dashboard.NewDashboardBuilder("Anubis – blog.kevingomez.fr").
		Readonly().
		Refresh("30s").
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Annotations(shared.DefaultAnnotations()).
		WithVariable(dashboard.NewQueryVariableBuilder("instance").
			Label("Instance").
			Query(dashboard.StringOrMap{
				Map: map[string]any{
					"qryType": 1,
					"query":   "label_values(anubis_challenges_issued{job=~\"integrations/anubis-blog\"},instance)", "refId": "PrometheusVariableQueryEditor-VariableQuery",
				},
			}).
			Datasource(shared.DefaultPrometheusDatasource()).
			Current(shared.AllVariableOption()).
			Multi(true).
			Refresh(dashboard.VariableRefreshOnDashboardLoad).
			Sort(dashboard.VariableSortNaturalAsc).
			IncludeAll(true),
		).
		WithPanel(requestRatesTimeseries().Span(12).Height(8)).
		WithPanel(verdictsPieChart().Span(6).Height(8)).
		WithPanel(challengeResultsPieChart().Span(6).Height(8)).
		WithPanel(allowedRequestsLogs().Span(12).Height(8)).
		WithPanel(logsVolumeTimeseries().Span(12).Height(8))
}
