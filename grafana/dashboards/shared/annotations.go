package shared

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func DefaultAnnotations() []cog.Builder[dashboard.AnnotationQuery] {
	return []cog.Builder[dashboard.AnnotationQuery]{
		dashboard.NewAnnotationQueryBuilder().
			Name("Annotations & Alerts").
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr("grafana"), Uid: cog.ToPtr("-- Grafana --")}).
			Hide(true).
			IconColor("rgba(0, 211, 255, 1)").
			Type("dashboard").
			BuiltIn(1),
	}
}
