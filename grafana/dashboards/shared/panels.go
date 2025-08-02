package shared

import (
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
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
		GradientMode(common.GraphGradientModeOpacity).
		ColorScheme(dashboard.NewFieldColorBuilder().Mode(dashboard.FieldColorModeIdPaletteClassic)).
		Legend(common.NewVizLegendOptionsBuilder().
			DisplayMode(common.LegendDisplayModeList).
			Placement(common.LegendPlacementBottom).
			ShowLegend(true),
		)
}
