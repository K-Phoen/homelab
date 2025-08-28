package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/K-Phoen/homelab/grafana/dashboards/anubis"
	"github.com/K-Phoen/homelab/grafana/dashboards/gitea"
	"github.com/K-Phoen/homelab/grafana/dashboards/keepalived"
	"github.com/K-Phoen/homelab/grafana/dashboards/metallb"
	"github.com/K-Phoen/homelab/grafana/dashboards/rooms"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/resource"
)

const defaultOutputDir = "./"

func main() {
	outputDir := defaultOutputDir
	if len(os.Args) == 2 {
		outputDir = os.Args[1]
	}

	factories := []struct {
		name   string
		folder string
		fn     func() *dashboard.DashboardBuilder
	}{
		{
			// "Living Room"
			name:   "de81hbkf5okqoe",
			folder: "cdl4fwl71924gc",
			fn:     rooms.LivingRoomDashboard,
		},
		{
			// "Anubis – blog.kevingomez.fr"
			name:   "een0nilkaqpz4a",
			folder: "cdl4fwl71924gc",
			fn:     anubis.BlogDashboard,
		},
		{
			// "Gitea Overview"
			name:   "gitea-overview",
			folder: "cdl4fwl71924gc",
			fn:     gitea.OverviewDashboard,
		},
		{
			// "MetalLb Overview"
			name:   "metallb-overview",
			folder: "cdl4fwl71924gc",
			fn:     metallb.OverviewDashboard,
		},
		{
			// "Keepalived Overview - blocky"
			name:   "keepalived-overview-blocky",
			folder: "cdl4fwl71924gc",
			fn: func() *dashboard.DashboardBuilder {
				return keepalived.OverviewDashboard(keepalived.Options{
					Title:         "blocky",
					ScriptName:    "blocky_responding",
					VirtualRouter: 42,
				})
			},
		},
	}

	for _, factory := range factories {
		spec, err := factory.fn().Build()
		if err != nil {
			panic(fmt.Errorf("could not build dashboard \"%s\": %w", factory.name, err))
		}

		manifest := resource.Manifest{
			ApiVersion: "dashboard.grafana.app/v1beta1",
			Kind:       "Dashboard",
			Metadata: resource.Metadata{
				Name:      factory.name,
				Namespace: cog.ToPtr("stacks-660328"),
				Annotations: map[string]string{
					"grafana.app/folder": factory.folder,
				},
			},
			Spec: spec,
		}

		payload, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			panic(fmt.Errorf("could not marshal dashboard \"%s\": %w", factory.name, err))
		}

		if err := os.WriteFile(path.Join(outputDir, factory.name+".json"), payload, 0644); err != nil {
			panic(fmt.Errorf("could not write dashboard \"%s\": %w", factory.name, err))
		}
	}
}
