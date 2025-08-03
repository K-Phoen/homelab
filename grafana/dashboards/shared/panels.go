package shared

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/piechart"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

func StatPanel(title string) *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title(title).
		Transparent(true).
		GraphMode(common.BigValueGraphModeNone).
		ColorMode(common.BigValueColorModeValue).
		JustifyMode(common.BigValueJustifyModeAuto).
		TextMode(common.BigValueTextModeAuto).
		Orientation(common.VizOrientationAuto).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Values(false).
			Calcs([]string{"lastNotNull"}),
		).
		PercentChangeColorMode(common.PercentChangeColorModeSameAsValue).
		ColorScheme(
			dashboard.NewFieldColorBuilder().
				Mode(dashboard.FieldColorModeIdFixed).
				FixedColor("green"),
		)
}

func TimeseriesPanel(title string) *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		Title(title).
		Transparent(true).
		FillOpacity(10).
		LineWidth(1).
		LineInterpolation(common.LineInterpolationSmooth).
		GradientMode(common.GraphGradientModeOpacity).
		ColorScheme(dashboard.NewFieldColorBuilder().Mode(dashboard.FieldColorModeIdPaletteClassic)).
		Legend(common.NewVizLegendOptionsBuilder().
			DisplayMode(common.LegendDisplayModeList).
			Placement(common.LegendPlacementBottom).
			ShowLegend(true),
		).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone).
			HideZeros(false),
		)
}

func PieChartPanel(title string) *piechart.PanelBuilder {
	return piechart.NewPanelBuilder().
		Title(title).
		Transparent(true).
		ColorScheme(dashboard.NewFieldColorBuilder().Mode(dashboard.FieldColorModeIdPaletteClassic)).
		PieType(piechart.PieChartTypePie).
		ReduceOptions(common.NewReduceDataOptionsBuilder().
			Values(false).
			Calcs([]string{"lastNotNull"}),
		).
		Legend(piechart.NewPieChartLegendOptionsBuilder().
			DisplayMode(common.LegendDisplayModeList).
			Placement(common.LegendPlacementBottom).
			ShowLegend(true),
		).
		Tooltip(common.NewVizTooltipOptionsBuilder().
			Mode(common.TooltipDisplayModeMulti).
			Sort(common.SortOrderNone).
			HideZeros(false),
		)
}

func LogPanel(title string) *logs.PanelBuilder {
	return logs.NewPanelBuilder().
		Title(title).
		Transparent(true).
		ShowTime(true).
		SortOrder(common.LogsSortOrderDescending).
		EnableLogDetails(true)
}
