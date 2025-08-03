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

func DefaultLokiDatasource() dashboard.DataSourceRef {
	return dashboard.DataSourceRef{
		Type: cog.ToPtr("loki"),
		Uid:  cog.ToPtr("grafanacloud-logs"),
	}
}

func AllVariableOption() dashboard.VariableOption {
	return dashboard.VariableOption{
		Text:  dashboard.StringOrArrayOfString{String: cog.ToPtr("All")},
		Value: dashboard.StringOrArrayOfString{ArrayOfString: []string{"$__all"}},
	}
}
