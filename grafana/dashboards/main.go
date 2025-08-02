package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/K-Phoen/homelab/grafana/dashboards/homelab"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/resource"
)

const defaultOutputDir = "./dashboards"

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
			fn:     homelab.LivingRoomDashboard,
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
