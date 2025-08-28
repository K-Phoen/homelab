package keepalived

import (
	"fmt"

	"github.com/K-Phoen/homelab/grafana/dashboards/shared"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/statushistory"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

type Options struct {
	Title         string
	ScriptName    string
	VirtualRouter int
}

func keepalivedStatus() *stat.PanelBuilder {
	return shared.StatPanel("Keepalived Status").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`keepalived_up{instance=~"$instance"}`).Instant().LegendFormat("{{ instance }}")).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Values(false).
			Calcs([]string{"last"}),
		).
		ColorMode(common.BigValueColorModeBackground).
		Unit(units.NoUnit).
		Mappings([]dashboard.ValueMapping{
			{
				ValueMap: &dashboard.ValueMap{
					Type: "value",
					Options: map[string]dashboard.ValueMappingResult{
						"0": {Index: cog.ToPtr[int32](0), Text: cog.ToPtr("DOWN")},
						"1": {Index: cog.ToPtr[int32](1), Text: cog.ToPtr("UP")},
					},
				},
			},
		}).
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{
					{Color: "red"},
					{Color: "green", Value: cog.ToPtr[float64](1)},
				}),
		)
}

func exporterStatus() *stat.PanelBuilder {
	return shared.StatPanel("Keepalived-exporter Status").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`keepalived_exporter_build_info{instance=~"$instance"}`).Instant().LegendFormat("{{ instance }}")).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Values(false).
			Calcs([]string{"last"}),
		).
		ColorMode(common.BigValueColorModeBackground).
		Unit(units.NoUnit).
		Mappings([]dashboard.ValueMapping{
			{
				ValueMap: &dashboard.ValueMap{
					Type: "value",
					Options: map[string]dashboard.ValueMappingResult{
						"0": {Index: cog.ToPtr[int32](0), Text: cog.ToPtr("DOWN")},
						"1": {Index: cog.ToPtr[int32](1), Text: cog.ToPtr("UP")},
					},
				},
			},
		}).
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{
					{Color: "red"},
					{Color: "green", Value: cog.ToPtr[float64](1)},
				}),
		)
}

func scriptStatus(opts Options) *stat.PanelBuilder {
	return shared.StatPanel("Script Status").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(fmt.Sprintf(`keepalived_script_status{instance=~"$instance", name="%s"}`, opts.ScriptName)).Instant().LegendFormat("{{ instance }}")).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Values(false).
			Calcs([]string{"last"}),
		).
		ColorMode(common.BigValueColorModeBackground).
		Unit(units.NoUnit).
		Mappings([]dashboard.ValueMapping{
			{
				ValueMap: &dashboard.ValueMap{
					Type: "value",
					Options: map[string]dashboard.ValueMappingResult{
						"0": {Index: cog.ToPtr[int32](0), Text: cog.ToPtr("KO")},
						"1": {Index: cog.ToPtr[int32](1), Text: cog.ToPtr("OK")},
					},
				},
			},
		}).
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{
					{Color: "red"},
					{Color: "green", Value: cog.ToPtr[float64](1)},
				}),
		)
}

func nodeStatus() *statushistory.PanelBuilder {
	return statushistory.NewPanelBuilder().
		Title("Node Status").
		Transparent(true).
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`keepalived_vrrp_state{instance=~"$instance"}`).LegendFormat("{{ instance }}")).
		ColWidth(0.5).
		RowHeight(0.5).
		Legend(
			common.NewVizLegendOptionsBuilder().
				ShowLegend(true).
				DisplayMode(common.LegendDisplayModeList).
				Placement(common.LegendPlacementBottom),
		).
		Mappings([]dashboard.ValueMapping{
			{
				ValueMap: &dashboard.ValueMap{
					Type: "value",
					Options: map[string]dashboard.ValueMappingResult{
						"0": {Index: cog.ToPtr[int32](0), Color: cog.ToPtr("red"), Text: cog.ToPtr("Unknown")},
						"1": {Index: cog.ToPtr[int32](1), Color: cog.ToPtr("orange"), Text: cog.ToPtr("Backup")},
						"2": {Index: cog.ToPtr[int32](2), Color: cog.ToPtr("green"), Text: cog.ToPtr("Master")},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder())
}

func masterChangesOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Master Changes").
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_become_master_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("{{ instance }} become").
				RefId("A"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_release_master_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("{{ instance }} release").
				RefId("B"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Short).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{Id: "custom.transform", Value: "negative-Y"},
		})
}

func advertisementsOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Advertisements").
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_advertisements_sent_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("{{ instance }} sent").
				RefId("A"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_advertisements_received_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("{{ instance }} received").
				RefId("B"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Short).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{Id: "custom.transform", Value: "negative-Y"},
		}).
		Legend(
			common.NewVizLegendOptionsBuilder().
				ShowLegend(true).
				DisplayMode(common.LegendDisplayModeTable).
				Placement(common.LegendPlacementBottom).
				Calcs([]string{"min", "max", "mean", "last"}).
				SortDesc(true).
				SortBy("max"),
		)
}

func priorityZeroOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Priority Zero").
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_priority_zero_sent_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("{{ instance }} sent").
				RefId("A"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_priority_zero_received_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("{{ instance }} received").
				RefId("B"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Short).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{Id: "custom.transform", Value: "negative-Y"},
		}).
		Legend(
			common.NewVizLegendOptionsBuilder().
				ShowLegend(true).
				DisplayMode(common.LegendDisplayModeTable).
				Placement(common.LegendPlacementBottom).
				Calcs([]string{"min", "max", "mean", "last"}).
				SortDesc(true).
				SortBy("max"),
		)
}

func gratuitousARPDelayOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Gratuitous ARP Delay").
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_gratuitous_arp_delay_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("{{ instance }}"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Milliseconds).
		Legend(
			common.NewVizLegendOptionsBuilder().
				ShowLegend(true).
				DisplayMode(common.LegendDisplayModeTable).
				Placement(common.LegendPlacementBottom).
				Calcs([]string{"min", "max", "mean", "last"}).
				SortDesc(true).
				SortBy("max"),
		)
}

func errorsOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Errors").
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_advertisements_interval_errors_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("Advertisement interval errors - {{ instance }}"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_packet_length_errors_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("Packet length errors - {{ instance }}"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_ip_ttl_errors_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("IP TTL errors - {{ instance }}"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_invalid_type_received_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("Invalid type received - {{ instance }}"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_address_list_errors_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("Address list errors - {{ instance }}"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_authentication_invalid_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("Authentication invalid - {{ instance }}"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_authentication_failure_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("Authentication failure - {{ instance }}"),
		).
		WithTarget(
			shared.PrometheusQuery(fmt.Sprintf(`rate(keepalived_authentication_mismatch_total{instance=~"$instance", vrid="%d"}[$__rate_interval])`, opts.VirtualRouter)).
				LegendFormat("Authentication mismatch - {{ instance }}"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Short).
		Legend(
			common.NewVizLegendOptionsBuilder().
				ShowLegend(true).
				DisplayMode(common.LegendDisplayModeTable).
				Placement(common.LegendPlacementRight).
				Calcs([]string{"min", "max", "mean", "last"}).
				SortDesc(true).
				SortBy("max"),
		)
}

func OverviewDashboard(opts Options) *dashboard.DashboardBuilder {
	return dashboard.NewDashboardBuilder(fmt.Sprintf("Keepalived Overview - %s", opts.Title)).
		Tags([]string{"generated"}).
		Readonly().
		Refresh("30s").
		Time("now-3h", "now").
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Annotations(shared.DefaultAnnotations()).
		WithVariable(
			dashboard.NewQueryVariableBuilder("instance").
				Label("Instance").
				Datasource(shared.DefaultPrometheusDatasource()).
				Refresh(dashboard.VariableRefreshOnTimeRangeChanged).
				IncludeAll(true).
				Query(dashboard.StringOrMap{
					Map: map[string]any{
						"query": "label_values(keepalived_up, instance)",
						"refId": "StandardVariableQuery",
					},
				}),
		).
		WithPanel(keepalivedStatus().Span(8).Height(5)).
		WithPanel(exporterStatus().Span(8).Height(5)).
		WithPanel(scriptStatus(opts).Span(8).Height(5)).
		WithPanel(nodeStatus().Span(24).Height(7)).
		WithPanel(masterChangesOverTime(opts).Height(8)).
		WithPanel(advertisementsOverTime(opts).Height(8)).
		WithPanel(priorityZeroOverTime(opts).Height(8)).
		WithPanel(gratuitousARPDelayOverTime(opts).Height(8)).
		WithPanel(errorsOverTime(opts).Span(24).Height(8))
}
