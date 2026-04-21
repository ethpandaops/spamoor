package plugin

import (
	"github.com/ethpandaops/spamoor/plugins/_example-plugin/example1"
	"github.com/ethpandaops/spamoor/scenario"
)

// PluginDescriptor defines the plugin metadata and scenarios.
var PluginDescriptor = scenario.PluginDescriptor{
	Name:        "example-plugin",
	Description: "Example plugin with sample scenarios",
	Categories: []*scenario.Category{
		{
			Name:        "Simple",
			Description: "Simple example scenarios",
			Descriptors: []*scenario.Descriptor{
				&example1.ScenarioDescriptor,
			},
		},
	},
}
