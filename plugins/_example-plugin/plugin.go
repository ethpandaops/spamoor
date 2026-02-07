package plugin

import (
	"github.com/ethpandaops/spamoor/plugins/_example-plugin/example1"
	"github.com/ethpandaops/spamoor/scenario"
)

// PluginDescriptor defines the plugin metadata and scenarios.
var PluginDescriptor = scenario.PluginDescriptor{
	Name:        "test-plugin",
	Description: "Test plugin with sample scenarios",
	Scenarios: []*scenario.Descriptor{
		&example1.ScenarioDescriptor,
	},
}
