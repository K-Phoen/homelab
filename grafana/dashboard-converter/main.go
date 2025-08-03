package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func main() {
	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	inputJSON, err := os.ReadFile("een0nilkaqpz4a.json")
	if err != nil {
		panic(err)
	}

	manifest := struct {
		Spec dashboard.Dashboard `json:"spec"`
	}{}
	if err = json.Unmarshal(inputJSON, &manifest); err != nil {
		panic(err)
	}

	converted := dashboard.DashboardConverter(manifest.Spec)
	fmt.Println(converted)
}
