package scenarios

import (
	"fmt"
	"slices"
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

// nativeScenarioCategories contains the tree of native scenario categories and descriptors.
// The order of the categories and descriptors is important for the CLI help text,
// validation, and displaying available options to users.
var nativeScenarioCategories = []*scenario.Category{
	{
		Name:        "Simple",
		Description: "Simple scenarios",
		Descriptors: []*scenario.Descriptor{
			&blobaverage.ScenarioDescriptor,
			&blobcombined.ScenarioDescriptor,
			&blobconflicting.ScenarioDescriptor,
			&blobs.ScenarioDescriptor,
			&blobreplacements.ScenarioDescriptor,
			&deploydestruct.ScenarioDescriptor,
			&eoatx.ScenarioDescriptor,
			&erc20tx.ScenarioDescriptor,
			&erc721tx.ScenarioDescriptor,
			&erc1155tx.ScenarioDescriptor,
			&evmfuzz.ScenarioDescriptor,
			&gasburnertx.ScenarioDescriptor,
			&setcodetx.ScenarioDescriptor,
			&uniswapswaps.ScenarioDescriptor,
			&xentoken.ScenarioDescriptor,
		},
	},
	{
		Name:        "Complex",
		Description: "Complex scenarios (require additional configuration)",
		Descriptors: []*scenario.Descriptor{
			&calltx.ScenarioDescriptor,
			&deploytx.ScenarioDescriptor,
			&factorydeploytx.ScenarioDescriptor,
			&geastx.ScenarioDescriptor,
			&replayeest.ScenarioDescriptor,
			&storagespam.ScenarioDescriptor,
			&taskrunner.ScenarioDescriptor,
		},
	},
	{
		Name:        "Bloatnet",
		Description: "Scenarios specifically designed for state bloating",
		Descriptors: []*scenario.Descriptor{
			&erc20bloater.ScenarioDescriptor,
		},
	},
	{
		Name:        "Utility",
		Description: "Utility scenarios",
		Descriptors: []*scenario.Descriptor{
			&wallets.ScenarioDescriptor,
		},
	},
}

var (
	pluginRegistry   *plugin.PluginRegistry
	scenarioRegistry *plugin.ScenarioRegistry
	initOnce         sync.Once
)

// nativeScenarios contains all built-in scenario descriptors flattened from categories.
var nativeScenarios []*scenario.Descriptor

func init() {
	// Flatten all native scenario descriptors from categories
	nameMap := make(map[string]*scenario.Descriptor, 32)

	var addCategory func(category *scenario.Category)
	addCategory = func(category *scenario.Category) {
		for _, descriptor := range category.Descriptors {
			if _, ok := nameMap[descriptor.Name]; ok {
				panic(fmt.Sprintf("scenario descriptor %s already registered", descriptor.Name))
			}
			nameMap[descriptor.Name] = descriptor
		}
		nativeScenarios = append(nativeScenarios, category.Descriptors...)
		for _, child := range category.Children {
			addCategory(child)
		}
	}
	for _, category := range nativeScenarioCategories {
		addCategory(category)
	}
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

// GetScenarioCategories returns a slice containing all scenario categories including
// both native and plugin scenarios. Plugin scenarios are grouped under a "Plugins"
// category with sub-categories for each plugin.
// This is useful for CLI help text, validation, and displaying available options to users.
func GetScenarioCategories() []*scenario.Category {
	// Start with a copy of native categories
	result := make([]*scenario.Category, len(nativeScenarioCategories))
	copy(result, nativeScenarioCategories)

	// Get plugin scenarios and group them by plugin name
	registry := GetScenarioRegistry()
	pluginScenarios := registry.GetPluginScenarios()

	if len(pluginScenarios) > 0 {
		// Group scenarios by plugin name
		pluginGroups := make(map[string][]*scenario.Descriptor, 8)
		pluginDescriptions := make(map[string]string, 8)

		for _, entry := range pluginScenarios {
			if entry.Plugin != nil && entry.Plugin.Descriptor != nil {
				pluginName := entry.Plugin.Descriptor.Name
				pluginGroups[pluginName] = append(pluginGroups[pluginName], entry.Descriptor)
				if entry.Plugin.Descriptor.Description != "" {
					pluginDescriptions[pluginName] = entry.Plugin.Descriptor.Description
				}
			}
		}

		// Create sub-categories for each plugin
		pluginChildren := make([]*scenario.Category, 0, len(pluginGroups))
		for pluginName, descriptors := range pluginGroups {
			desc := pluginDescriptions[pluginName]
			if desc == "" {
				desc = fmt.Sprintf("Scenarios from %s plugin", pluginName)
			}
			pluginChildren = append(pluginChildren, &scenario.Category{
				Name:        pluginName,
				Description: desc,
				Descriptors: descriptors,
			})
		}

		// Add plugin category if there are any plugin scenarios
		if len(pluginChildren) > 0 {
			// Sort plugin categories by name for consistent ordering
			slices.SortFunc(pluginChildren, func(a, b *scenario.Category) int {
				if a.Name < b.Name {
					return -1
				}
				if a.Name > b.Name {
					return 1
				}
				return 0
			})

			pluginsCategory := &scenario.Category{
				Name:        "Plugins",
				Description: "Scenarios loaded from plugins",
				Children:    pluginChildren,
			}
			result = append(result, pluginsCategory)
		}
	}

	return result
}
