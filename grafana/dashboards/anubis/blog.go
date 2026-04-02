package anubis

import (
	"fmt"

	"github.com/K-Phoen/homelab/grafana/dashboards/shared"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/piechart"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func requestRatesTimeseries(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Request rates").
		WithTarget(shared.PrometheusQuery("sum(irate(anubis_policy_results{instance=~\"$instance\",container=\"anubis\",job=\"%s\"}[$__rate_interval])) by (action, rule)", opts.Integration).
			LegendFormat("{{action}} - {{rule}}"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.RequestsPerSecond).
		Min(0).
		GradientMode(common.GraphGradientModeNone).
		ThresholdsStyle(common.NewGraphThresholdsStyleConfigBuilder().Mode(common.GraphThresholdsStyleModeOff)).
		FillOpacity(0).
		ShowPoints(common.VisibilityModeAuto).
		PointSize(5)
}

func verdictsPieChart(opts Options) *piechart.PanelBuilder {
	return shared.PieChartPanel("Verdicts").
		WithTarget(shared.PrometheusQuery("sum by(action) (increase(anubis_policy_results{instance=~\"$instance\",container=\"anubis\", job=\"%s\"}[$__range]))", opts.Integration).
			LegendFormat("{{action}}"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.NoUnit).
		Min(0).
		ColorScheme(dashboard.NewFieldColorBuilder().Mode("palette-classic")).
		OverrideByName("WEIGH", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "yellow", "mode": "fixed"}},
		}).
		OverrideByName("CHALLENGE", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "light-orange", "mode": "fixed"}},
		}).
		OverrideByName("DENY", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "red", "mode": "fixed"}}},
		).
		OverrideByName("ALLOW", []dashboard.DynamicConfigValue{
			{Id: "color", Value: map[string]any{"fixedColor": "green", "mode": "fixed"}}},
		)
}

func challengeResultsPieChart(opts Options) *piechart.PanelBuilder {
	return shared.PieChartPanel("Challenge results").
		WithTarget(shared.PrometheusQuery("increase(anubis_challenges_validated{instance=~\"$instance\",container=\"anubis\",job=\"%s\"}[$__range])", opts.Integration).
			LegendFormat("Validated"),
		).
		WithTarget(shared.PrometheusQuery("increase(anubis_failed_validations{instance=~\"$instance\",container=\"anubis\",job=\"%s\"}[$__range])", opts.Integration).
			LegendFormat("Failed"),
		).
		WithTarget(shared.PrometheusQuery("sum by(job, instance) (increase(anubis_challenges_issued{job=\"%[1]s\",container=\"anubis\",instance=~\"$instance\"}[$__range])) -  increase(anubis_challenges_validated{job=\"%[1]s\",instance=~\"$instance\"}[$__range]) - increase(anubis_failed_validations{job=\"%[1]s\",instance=~\"$instance\"}[$__range])", opts.Integration).
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

func logsViewer(opts Options) *logs.PanelBuilder {
	return shared.LogPanel("Allowed requests").
		WithTarget(shared.LokiQuery(fmt.Sprintf("{cluster=~\"homelab\", namespace=\"%s\", container=\"%s\"}", opts.Namespace, opts.Container))).
		Datasource(shared.DefaultLokiDatasource())
}

func logsVolumeTimeseries(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Logs volume").
		WithTarget(shared.LokiQuery(fmt.Sprintf("sum (count_over_time({cluster=~\"homelab\", namespace=\"%s\", container=\"%s\"} [$__auto]))", opts.Namespace, opts.Container)).
			LegendFormat("Logs"),
		).
		Datasource(shared.DefaultLokiDatasource()).
		DrawStyle(common.GraphDrawStyleBars).
		GradientMode(common.GraphGradientModeNone).
		ShowPoints(common.VisibilityModeNever).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false))
}

type Options struct {
	Namespace   string
	Container   string
	Website     string
	Integration string
}

func Dashboard(opts Options) *dashboard.DashboardBuilder {
	return dashboard.NewDashboardBuilder(fmt.Sprintf("Anubis – %s", opts.Website)).
		Tags([]string{"generated"}).
		Readonly().
		Refresh("30s").
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Annotations(shared.DefaultAnnotations()).
		WithVariable(dashboard.NewQueryVariableBuilder("instance").
			Label("Instance").
			Query(dashboard.StringOrMap{
				Map: map[string]any{
					"qryType": 1,
					"query":   fmt.Sprintf("label_values(anubis_policy_results{job=~\"%s\"},instance)", opts.Integration), "refId": "PrometheusVariableQueryEditor-VariableQuery",
				},
			}).
			Datasource(shared.DefaultPrometheusDatasource()).
			Current(shared.AllVariableOption()).
			Multi(true).
			Refresh(dashboard.VariableRefreshOnDashboardLoad).
			Sort(dashboard.VariableSortNaturalAsc).
			AllValue(".+").
			IncludeAll(true),
		).
		WithPanel(requestRatesTimeseries(opts).Span(12).Height(8)).
		WithPanel(verdictsPieChart(opts).Span(6).Height(8)).
		WithPanel(challengeResultsPieChart(opts).Span(6).Height(8)).
		WithPanel(logsVolumeTimeseries(opts).Span(24).Height(6)).
		WithPanel(logsViewer(opts).Span(24).Height(14))
}
