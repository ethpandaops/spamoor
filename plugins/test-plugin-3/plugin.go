package plugin

import (
	"github.com/ethpandaops/spamoor/plugins/test-plugin-3/erc20tx"
	"github.com/ethpandaops/spamoor/scenario"
)

// PluginDescriptor defines the plugin metadata and scenarios.
// Note: This uses an anonymous struct compatible with plugin.Descriptor
// to avoid import cycles with the plugin package.
var PluginDescriptor = struct {
	Name        string
	Description string
	Scenarios   []*scenario.Descriptor
}{
	Name:        "test-plugin",
	Description: "Test plugin with sample scenarios",
	Scenarios: []*scenario.Descriptor{
		&erc20tx.ScenarioDescriptor,
	},
}
