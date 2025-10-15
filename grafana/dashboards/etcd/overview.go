package etcd

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
	Title   string
	Cluster string
}

func etcdStatus(opts Options) *stat.PanelBuilder {
	return shared.StatPanel("etcd Status").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`etcd_server_has_leader{cluster="%s"}`, opts.Cluster).Instant().LegendFormat("{{ instance }}")).
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

func rpcRateOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("RPC rate").
		WithTarget(
			shared.PrometheusQuery(`sum(rate(grpc_server_started_total{job=~".*etcd.*", cluster="%s", grpc_type="unary"}[$__rate_interval]))`, opts.Cluster).
				LegendFormat("RPC rate"),
		).
		WithTarget(
			shared.PrometheusQuery(`sum(rate(grpc_server_handled_total{job=~".*etcd.*", cluster="%s", grpc_type="unary", grpc_code=~"Unknown|FailedPrecondition|ResourceExhausted|Internal|Unavailable|DataLoss|DeadlineExceeded"}[$__rate_interval]))`, opts.Cluster).
				LegendFormat("RPC failed rate"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.OpsPerSecond)
}

func activeStreamsOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Active streams").
		WithTarget(
			shared.PrometheusQuery(`sum(grpc_server_started_total{job=~".*etcd.*",cluster="%[1]s",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"}) - sum(grpc_server_handled_total{cluster="%[1]s",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"})`, opts.Cluster).
				LegendFormat("Watch streams"),
		).
		WithTarget(
			shared.PrometheusQuery(`sum(grpc_server_started_total{job=~".*etcd.*",cluster="%[1]s",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"}) - sum(grpc_server_handled_total{cluster="%[1]s",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"})`, opts.Cluster).
				LegendFormat("Lease streams"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Short)
}

func nodeStatus(opts Options) *statushistory.PanelBuilder {
	return statushistory.NewPanelBuilder().
		Title("Node Status").
		Transparent(true).
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery(`etcd_server_is_leader{job=~".*etcd.*", cluster="%s"}`, opts.Cluster).LegendFormat("{{ instance }}")).
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
						"0": {Index: cog.ToPtr[int32](0), Color: cog.ToPtr("orange"), Text: cog.ToPtr("Learner")},
						"1": {Index: cog.ToPtr[int32](1), Color: cog.ToPtr("green"), Text: cog.ToPtr("Leader")},
					},
				},
			},
		}).
		Thresholds(dashboard.NewThresholdsConfigBuilder().Mode(dashboard.ThresholdsModeAbsolute))
}

func dbSizeOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("DB size").
		WithTarget(
			shared.PrometheusQuery(`etcd_mvcc_db_total_size_in_bytes{job=~".*etcd.*", cluster="%s"}`, opts.Cluster).
				LegendFormat("{{instance}} DB size"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.BytesIEC)
}

func dbSyncDurationOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Disk sync duration").
		WithTarget(
			shared.PrometheusQuery(`histogram_quantile(0.99, sum(rate(etcd_disk_wal_fsync_duration_seconds_bucket{job=~".*etcd.*", cluster="%s"}[$__rate_interval])) by (instance, le))`, opts.Cluster).
				LegendFormat("{{instance}} WAL fsync"),
		).
		WithTarget(
			shared.PrometheusQuery(`histogram_quantile(0.99, sum(rate(etcd_disk_backend_commit_duration_seconds_bucket{job=~".*etcd.*", cluster="%s"}[$__rate_interval])) by (instance, le))`, opts.Cluster).
				LegendFormat("{{instance}} DB fsync"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Seconds)
}

func memoryOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Memory").
		WithTarget(
			shared.PrometheusQuery(`process_resident_memory_bytes{job=~".*etcd.*", cluster="%s"}`, opts.Cluster).
				LegendFormat("{{instance}} resident memory"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.BytesIEC)
}

func clientTrafficOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Client traffic").
		WithTarget(
			shared.PrometheusQuery(`rate(etcd_network_client_grpc_received_bytes_total{job=~".*etcd.*", cluster="%s"}[$__rate_interval])`, opts.Cluster).
				LegendFormat("{{instance}} traffic in").
				RefId("A"),
		).
		WithTarget(
			shared.PrometheusQuery(`rate(etcd_network_client_grpc_sent_bytes_total{job=~".*etcd.*", cluster="%s"}[$__rate_interval])`, opts.Cluster).
				LegendFormat("{{instance}} traffic out").
				RefId("B"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.BytesPerSecondSI).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "custom.transform",
				Value: "negative-Y",
			},
		})
}

func peerTrafficOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Peer traffic").
		WithTarget(
			shared.PrometheusQuery(`sum(rate(etcd_network_peer_received_bytes_total{job=~".*etcd.*", cluster="%s"}[$__rate_interval])) by (instance)`, opts.Cluster).
				LegendFormat("{{instance}} traffic in").
				RefId("A"),
		).
		WithTarget(
			shared.PrometheusQuery(`rate(etcd_network_peer_sent_bytes_total{job=~".*etcd.*", cluster="%s"}[$__rate_interval])`, opts.Cluster).
				LegendFormat("{{instance}} traffic out").
				RefId("B"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.BytesPerSecondSI).
		OverrideByQuery("B", []dashboard.DynamicConfigValue{
			{
				Id:    "custom.transform",
				Value: "negative-Y",
			},
		})
}

func raftProposalsOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Raft proposals").
		WithTarget(
			shared.PrometheusQuery(`changes(etcd_server_leader_changes_seen_total{job=~".*etcd.*", cluster="%s"}[1d])`, opts.Cluster).
				LegendFormat("{{instance}} total leader elections per day"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Short)
}

func electionsPerDayOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Total leader elections per day").
		WithTarget(
			shared.PrometheusQuery(`changes(etcd_server_leader_changes_seen_total{job=~".*etcd.*", cluster="%s"}[1d])`, opts.Cluster).
				LegendFormat("{{instance}} total leader elections per day"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Short)
}

func peerRoundtripTimeOverTime(opts Options) *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("Peer round trip time").
		WithTarget(
			shared.PrometheusQuery(`histogram_quantile(0.99, sum by (instance, le) (rate(etcd_network_peer_round_trip_time_seconds_bucket{job=~".*etcd.*", cluster="%s"}[$__rate_interval])))`, opts.Cluster).
				LegendFormat("{{instance}} peer round trip time"),
		).
		Datasource(shared.DefaultPrometheusDatasource()).
		Unit(units.Seconds)
}

func OverviewDashboard(opts Options) *dashboard.DashboardBuilder {
	return dashboard.NewDashboardBuilder(fmt.Sprintf("etcd Overview - %s", opts.Title)).
		Tags([]string{"generated", "etcd"}).
		Readonly().
		Refresh("30s").
		Time("now-1h", "now").
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Annotations(shared.DefaultAnnotations()).
		WithPanel(etcdStatus(opts).Span(8).Height(6)).
		WithPanel(rpcRateOverTime(opts).Span(8).Height(6)).
		WithPanel(activeStreamsOverTime(opts).Span(8).Height(6)).
		WithPanel(nodeStatus(opts).Span(24).Height(7)).
		WithRow(dashboard.NewRowBuilder("Database")).
		WithPanel(dbSizeOverTime(opts).Span(8).Height(7)).
		WithPanel(dbSyncDurationOverTime(opts).Span(8).Height(7)).
		WithPanel(memoryOverTime(opts).Span(8).Height(7)).
		WithRow(dashboard.NewRowBuilder("Networking")).
		WithPanel(clientTrafficOverTime(opts).Span(12).Height(7)).
		WithPanel(peerTrafficOverTime(opts).Span(12).Height(7)).
		WithPanel(raftProposalsOverTime(opts).Span(8).Height(7)).
		WithPanel(electionsPerDayOverTime(opts).Span(8).Height(7)).
		WithPanel(peerRoundtripTimeOverTime(opts).Span(8).Height(7))
}
