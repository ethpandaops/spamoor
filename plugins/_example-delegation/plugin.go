package plugin

import (
	"github.com/ethpandaops/spamoor/plugins/_example-delegation/delegatedcounter"
	"github.com/ethpandaops/spamoor/scenario"
)

// PluginDescriptor defines the plugin metadata and scenarios.
var PluginDescriptor = scenario.PluginDescriptor{
	Name:        "example-delegation",
	Description: "Example plugin: EIP-7702 delegation workload that needs receipt-dependent setup",
	Categories: []*scenario.Category{
		{
			Name:        "Examples",
			Description: "Example plugin scenarios",
			Descriptors: []*scenario.Descriptor{
				&delegatedcounter.ScenarioDescriptor,
			},
		},
	},
}
