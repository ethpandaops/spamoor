package scenarios

import (
	"sync"

	"github.com/ethpandaops/spamoor/plugin"
	"github.com/ethpandaops/spamoor/scenario"

	blobaverage "github.com/ethpandaops/spamoor/scenarios/blob-average"
	blobcombined "github.com/ethpandaops/spamoor/scenarios/blob-combined"
	blobconflicting "github.com/ethpandaops/spamoor/scenarios/blob-conflicting"
	blobreplacements "github.com/ethpandaops/spamoor/scenarios/blob-replacements"
	"github.com/ethpandaops/spamoor/scenarios/blobs"
	"github.com/ethpandaops/spamoor/scenarios/calltx"
	deploydestruct "github.com/ethpandaops/spamoor/scenarios/deploy-destruct"
	"github.com/ethpandaops/spamoor/scenarios/deploytx"
	"github.com/ethpandaops/spamoor/scenarios/eoatx"
	"github.com/ethpandaops/spamoor/scenarios/erc1155tx"
	"github.com/ethpandaops/spamoor/scenarios/erc20tx"
	"github.com/ethpandaops/spamoor/scenarios/erc721tx"
	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
	"github.com/ethpandaops/spamoor/scenarios/factorydeploytx"
	"github.com/ethpandaops/spamoor/scenarios/gasburnertx"
	"github.com/ethpandaops/spamoor/scenarios/geastx"
	replayeest "github.com/ethpandaops/spamoor/scenarios/replay-eest"
	"github.com/ethpandaops/spamoor/scenarios/setcodetx"
	erc20bloater "github.com/ethpandaops/spamoor/scenarios/statebloat/erc20_bloater"
	"github.com/ethpandaops/spamoor/scenarios/storagespam"
	"github.com/ethpandaops/spamoor/scenarios/taskrunner"
	uniswapswaps "github.com/ethpandaops/spamoor/scenarios/uniswap-swaps"
	"github.com/ethpandaops/spamoor/scenarios/wallets"
	"github.com/ethpandaops/spamoor/scenarios/xentoken"
)

var (
	pluginRegistry   *plugin.PluginRegistry
	scenarioRegistry *plugin.ScenarioRegistry
	initOnce         sync.Once
)

// nativeScenarios contains all built-in scenario descriptors for the spamoor tool.
// This registry includes scenarios for testing various Ethereum transaction types and patterns.
// Each descriptor defines the configuration, constructor, and metadata for a specific test scenario.
var nativeScenarios = []*scenario.Descriptor{
	&blobaverage.ScenarioDescriptor,
	&blobcombined.ScenarioDescriptor,
	&blobconflicting.ScenarioDescriptor,
	&blobs.ScenarioDescriptor,
	&blobreplacements.ScenarioDescriptor,
	&calltx.ScenarioDescriptor,
	&deploydestruct.ScenarioDescriptor,
	&deploytx.ScenarioDescriptor,
	&eoatx.ScenarioDescriptor,
	&erc20bloater.ScenarioDescriptor,
	&erc20tx.ScenarioDescriptor,
	&erc721tx.ScenarioDescriptor,
	&erc1155tx.ScenarioDescriptor,
	&evmfuzz.ScenarioDescriptor,
	&factorydeploytx.ScenarioDescriptor,
	&gasburnertx.ScenarioDescriptor,
	&geastx.ScenarioDescriptor,
	&replayeest.ScenarioDescriptor,
	&setcodetx.ScenarioDescriptor,
	&storagespam.ScenarioDescriptor,
	&taskrunner.ScenarioDescriptor,
	&uniswapswaps.ScenarioDescriptor,
	&wallets.ScenarioDescriptor,
	&xentoken.ScenarioDescriptor,
}

// InitRegistries initializes the plugin and scenario registries.
// This function is idempotent and safe to call multiple times.
func InitRegistries() (*plugin.PluginRegistry, *plugin.ScenarioRegistry) {
	initOnce.Do(func() {
		pluginRegistry = plugin.NewPluginRegistry()
		scenarioRegistry = plugin.NewScenarioRegistry(nativeScenarios)
	})

	return pluginRegistry, scenarioRegistry
}

// GetPluginRegistry returns the global plugin registry.
// Panics if InitRegistries() has not been called.
func GetPluginRegistry() *plugin.PluginRegistry {
	if pluginRegistry == nil {
		InitRegistries()
	}

	return pluginRegistry
}

// GetScenarioRegistry returns the global scenario registry.
// Panics if InitRegistries() has not been called.
func GetScenarioRegistry() *plugin.ScenarioRegistry {
	if scenarioRegistry == nil {
		InitRegistries()
	}

	return scenarioRegistry
}

// GetScenario finds and returns a scenario descriptor by name.
// It performs a lookup through the scenario registry and returns
// the matching descriptor, or nil if no scenario with the given name exists.
// This function is thread-safe.
func GetScenario(name string) *scenario.Descriptor {
	return GetScenarioRegistry().GetDescriptor(name)
}

// GetScenarioEntry finds and returns a scenario entry by name.
// The entry includes the descriptor and metadata about its source (native or plugin).
// Returns nil if no scenario with the given name exists.
// This function is thread-safe.
func GetScenarioEntry(name string) *plugin.ScenarioEntry {
	return GetScenarioRegistry().Get(name)
}

// GetScenarioNames returns a slice containing the names of all registered scenarios.
// This is useful for CLI help text, validation, and displaying available options
// to users. This function is thread-safe.
func GetScenarioNames() []string {
	return GetScenarioRegistry().GetNames()
}
