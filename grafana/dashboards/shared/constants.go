package shared

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func DefaultPrometheusDatasource() dashboard.DataSourceRef {
	return dashboard.DataSourceRef{
		Type: cog.ToPtr("prometheus"),
		Uid:  cog.ToPtr("grafanacloud-prom"),
	}
}
