package metallb

import (
	"github.com/K-Phoen/homelab/grafana/dashboards/shared"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func addressUtilizationTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Address utilization").
		WithTarget(shared.PrometheusQuery(`metallb_allocator_addresses_in_use_total * 100 / metallb_allocator_addresses_total`).LegendFormat("{{ pool }}")).
		Datasource(shared.DefaultPrometheusDatasource()).
		Min(0).
		Unit(units.Percent)
}

func staleConfigStat() *stat.PanelBuilder {
	return shared.StatPanel("Stale config").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`sum(metallb_k8s_client_config_stale_bool)`).Instant()).
		Unit(units.NoUnit).
		ShowPercentChange(true)
}

func layer2Timeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Layer 2").
		WithTarget(shared.PrometheusQuery(`sum(rate(metallb_layer2_gratuitous_sent[$__rate_interval])) by (instance, ip)`).LegendFormat("gratuitous_sent instance={{instance}} ip={{ip}}")).
		WithTarget(shared.PrometheusQuery(`sum(rate(metallb_layer2_requests_received[$__rate_interval])) by (instance, ip)`).LegendFormat("requests_received instance={{instance}} ip={{ip}}")).
		WithTarget(shared.PrometheusQuery(`sum(rate(metallb_layer2_responses_sent[$__rate_interval])) by (instance, ip)`).LegendFormat("responses_sent instance={{instance}} ip={{ip}}")).
		Datasource(shared.DefaultPrometheusDatasource()).
		Min(0).
		Unit(units.NoUnit)
}

func clientUpdatesTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Layer 2").
		WithTarget(shared.PrometheusQuery(`rate(metallb_k8s_client_update_errors_total[$__rate_interval])`).LegendFormat("Errors - {{ container }}")).
		WithTarget(shared.PrometheusQuery(`rate(metallb_k8s_client_updates_total[$__rate_interval])`).LegendFormat("Updates - {{ container }}")).
		Datasource(shared.DefaultPrometheusDatasource()).
		Min(0).
		Unit(units.NoUnit)
}

func OverviewDashboard() *dashboard.DashboardBuilder {
	return dashboard.NewDashboardBuilder("MetalLb Overview").
		Tags([]string{"generated"}).
		Readonly().
		Refresh("30s").
		Time("now-1h", "now").
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Annotations(shared.DefaultAnnotations()).
		WithPanel(addressUtilizationTimeseries().Span(21).Height(6)).
		WithPanel(staleConfigStat().Span(3).Height(6)).
		WithPanel(layer2Timeseries().Span(24).Height(6)).
		WithPanel(clientUpdatesTimeseries().Span(24).Height(6))
}
