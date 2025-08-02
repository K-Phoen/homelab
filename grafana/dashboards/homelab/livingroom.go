package homelab

import (
	"github.com/K-Phoen/homelab/grafana/dashboards/shared"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"github.com/grafana/grafana-foundation-sdk/go/units"
)

func livingRoomTempStat() *stat.PanelBuilder {
	return shared.StatPanel("Temperature").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery("homelab_sensors_temp")).
		Unit(units.Celsius).
		Decimals(1)
}

func livingRoomTempTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(
			shared.PrometheusQuery("homelab_sensors_temp").
				LegendFormat("Temperature"),
		).
		Unit(units.Celsius).
		Decimals(1).
		MaxDataPoints(11000).
		SpanNulls(common.BoolOrFloat64{Bool: cog.ToPtr(true)}).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false))
}

func livingRoomHumidityStat() *stat.PanelBuilder {
	return shared.StatPanel("Humidity").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery("homelab_sensors_humidity")).
		Unit(units.Humidity).
		Decimals(1)
}

func livingRoomHumidityTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(
			shared.PrometheusQuery("homelab_sensors_humidity").
				LegendFormat("Humidity"),
		).
		Unit(units.Celsius).
		Decimals(1).
		MaxDataPoints(11000).
		SpanNulls(common.BoolOrFloat64{Bool: cog.ToPtr(true)}).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false))
}

func livingRoomPressureStat() *stat.PanelBuilder {
	return shared.StatPanel("Atmospheric pressure").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(shared.PrometheusQuery("homelab_sensors_pressure / 100")).
		Unit(units.Hectopascals).
		Decimals(1)
}

func livingRoomPressureTimeseries() *timeseries.PanelBuilder {
	return shared.TimeseriesPanel("").
		Datasource(shared.DefaultPrometheusDatasource()).
		WithTarget(
			shared.PrometheusQuery("homelab_sensors_humidity").
				LegendFormat("Atmospheric pressure"),
		).
		Unit(units.Celsius).
		Decimals(1).
		MaxDataPoints(11000).
		SpanNulls(common.BoolOrFloat64{Bool: cog.ToPtr(true)}).
		Legend(common.NewVizLegendOptionsBuilder().ShowLegend(false))
}

func LivingRoomDashboard() *dashboard.DashboardBuilder {
	return dashboard.NewDashboardBuilder("Living room").
		Title("Living room").
		Readonly().
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Annotations(shared.DefaultAnnotations()).
		WithPanel(livingRoomTempStat().Span(6).Height(8)).
		WithPanel(livingRoomTempTimeseries().Span(18).Height(8)).
		WithPanel(livingRoomHumidityStat().Span(6).Height(8)).
		WithPanel(livingRoomHumidityTimeseries().Span(18).Height(8)).
		WithPanel(livingRoomPressureStat().Span(6).Height(8)).
		WithPanel(livingRoomPressureTimeseries().Span(18).Height(8))
}
