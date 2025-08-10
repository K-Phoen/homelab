package gitea

import (
	"github.com/K-Phoen/homelab/grafana/dashboards/shared"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func generalStat() *stat.PanelBuilder {
	return shared.StatPanel("").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery("gitea_organizations").LegendFormat("Organizations")).
		WithTarget(shared.PrometheusQuery("gitea_teams").LegendFormat("Teams")).
		WithTarget(shared.PrometheusQuery("gitea_users").LegendFormat("Users")).
		WithTarget(shared.PrometheusQuery("gitea_repositories").LegendFormat("Repositories")).
		WithTarget(shared.PrometheusQuery("gitea_milestones").LegendFormat("Milestones")).
		WithTarget(shared.PrometheusQuery("gitea_releases").LegendFormat("Releases")).
		WithTarget(shared.PrometheusQuery("gitea_issues_open").LegendFormat("Issues opened")).
		WithTarget(shared.PrometheusQuery("gitea_issues_closed").LegendFormat("Issues closed")).
		WithTarget(shared.PrometheusQuery("gitea_webhooks").LegendFormat("Webhooks")).
		Unit(units.NoUnit)
}

func versionStat() *stat.PanelBuilder {
	return shared.StatPanel("Version").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`gitea_build_info{job="integrations/gitea"}`).Format(prometheus.PromQueryFormatTable)).
		ReduceOptions(
			common.NewReduceDataOptionsBuilder().
				Calcs([]string{"lastNotNull"}).
				Fields("version"),
		)
}

func uptimeStat() *stat.PanelBuilder {
	return shared.StatPanel("Uptime").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`time()-process_start_time_seconds{job="integrations/gitea"}`)).
		Unit(units.Seconds).
		GraphMode(common.BigValueGraphModeArea)
}

func memoryUsageTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Memory usage").
		WithTarget(shared.PrometheusQuery(`process_resident_memory_bytes{job="integrations/gitea"}`)).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.BytesSI).
		Min(0).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false))
}

func cpuUsageTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("CPU usage").
		WithTarget(shared.PrometheusQuery(`rate(process_cpu_seconds_total{job="integrations/gitea"}[$__rate_interval])*100`)).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Percent).
		Min(0).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false))
}

func fileDescriptorsUsageTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("File descriptors usage").
		WithTarget(shared.PrometheusQuery(`process_open_fds{job="integrations/gitea"}`).RefId("A").LegendFormat("Open")).
		WithTarget(shared.PrometheusQuery(`process_max_fds{job="integrations/gitea"}`).RefId("B").LegendFormat("Maximum")).
		Datasource(shared.DefaultPrometheusDatasource()).
		Min(0).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false)).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "custom.lineStyle",
				Value: map[string]any{"dash": []int{10, 10}},
			},
			{
				Id:    "color",
				Value: map[string]any{"fixedColor": "red", "mode": "fixed"},
			},
			{
				Id:    "custom.fillOpacity",
				Value: 0,
			},
		})
}

func OverviewDashboard() *dashboard.DashboardBuilder {
	return dashboard.NewDashboardBuilder("Gitea Overview").
		Readonly().
		Refresh("30s").
		Time("now-1h", "now").
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Annotations(shared.DefaultAnnotations()).
		WithRow(dashboard.NewRowBuilder("General")).
		WithPanel(versionStat().Span(4).Height(4)).
		WithPanel(generalStat().Span(14).Height(4)).
		WithPanel(uptimeStat().Span(6).Height(4)).
		WithRow(dashboard.NewRowBuilder("System")).
		WithPanel(memoryUsageTimeseries().Span(8).Height(6)).
		WithPanel(cpuUsageTimeseries().Span(8).Height(6)).
		WithPanel(fileDescriptorsUsageTimeseries().Span(8).Height(6))
}
