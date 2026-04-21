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
	"github.com/ethpandaops/spamoor/scenarios/eip2780tx"
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
	"github.com/ethpandaops/spamoor/scenarios/storagerefundtx"
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
			&eip2780tx.ScenarioDescriptor,
			&eoatx.ScenarioDescriptor,
			&erc20tx.ScenarioDescriptor,
			&erc721tx.ScenarioDescriptor,
			&erc1155tx.ScenarioDescriptor,
			&evmfuzz.ScenarioDescriptor,
			&gasburnertx.ScenarioDescriptor,
			&setcodetx.ScenarioDescriptor,
			&storagerefundtx.ScenarioDescriptor,
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
	scenarioRegistry *scenario.Registry
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
func InitRegistries() (*plugin.PluginRegistry, *scenario.Registry) {
	initOnce.Do(func() {
		pluginRegistry = plugin.NewPluginRegistry()
		scenarioRegistry = scenario.NewRegistry(nativeScenarios)
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
func GetScenarioRegistry() *scenario.Registry {
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
func GetScenarioEntry(name string) *scenario.ScenarioEntry {
	return GetScenarioRegistry().Get(name)
}

// GetScenarioNames returns a slice containing the names of all registered scenarios.
// This is useful for CLI help text, validation, and displaying available options
// to users. This function is thread-safe.
func GetScenarioNames() []string {
	return GetScenarioRegistry().GetNames()
}

// GetScenarioCategories returns a slice containing all scenario categories including
// both native and plugin scenarios. Plugin categories are merged with native categories:
// - Plugin scenarios are added to existing categories if the category name matches
// - New categories from plugins are appended to the category list
// - Category descriptions can be updated by plugins (later plugin registrations override earlier ones)
// This is useful for CLI help text, validation, and displaying available options to users.
func GetScenarioCategories() []*scenario.Category {
	// Deep copy native categories so we can safely add plugin scenarios
	result := deepCopyCategories(nativeScenarioCategories)

	// Build a map for quick category lookup by name (including nested categories)
	categoryMap := make(map[string]*scenario.Category)
	var buildCategoryMap func(cats []*scenario.Category)
	buildCategoryMap = func(cats []*scenario.Category) {
		for _, cat := range cats {
			categoryMap[cat.Name] = cat
			buildCategoryMap(cat.Children)
		}
	}
	buildCategoryMap(result)

	// Get all loaded plugins and merge their categories
	plugins := GetPluginRegistry().GetAll()

	// Sort plugins by name for consistent ordering
	slices.SortFunc(plugins, func(a, b *plugin.LoadedPlugin) int {
		if a.Descriptor.Name < b.Descriptor.Name {
			return -1
		}
		if a.Descriptor.Name > b.Descriptor.Name {
			return 1
		}
		return 0
	})

	for _, loadedPlugin := range plugins {
		if loadedPlugin.Descriptor == nil {
			continue
		}

		// Merge each category from the plugin
		for _, pluginCat := range loadedPlugin.Descriptor.Categories {
			result = mergeCategory(result, categoryMap, pluginCat)
		}
	}

	return result
}

// deepCopyCategories creates a deep copy of categories.
func deepCopyCategories(cats []*scenario.Category) []*scenario.Category {
	if cats == nil {
		return nil
	}

	result := make([]*scenario.Category, len(cats))
	for i, cat := range cats {
		result[i] = &scenario.Category{
			Name:        cat.Name,
			Description: cat.Description,
			Descriptors: append([]*scenario.Descriptor(nil), cat.Descriptors...),
			Children:    deepCopyCategories(cat.Children),
		}
	}

	return result
}

// mergeCategory merges a plugin category into the result categories.
// If a category with the same name exists, scenarios are added to it.
// Otherwise, a new category is created.
func mergeCategory(result []*scenario.Category, categoryMap map[string]*scenario.Category, pluginCat *scenario.Category) []*scenario.Category {
	if existingCat, ok := categoryMap[pluginCat.Name]; ok {
		// Category exists - add scenarios and update description if provided
		if pluginCat.Description != "" {
			existingCat.Description = pluginCat.Description
		}

		existingCat.Descriptors = append(existingCat.Descriptors, pluginCat.Descriptors...)

		// Recursively merge children
		for _, child := range pluginCat.Children {
			existingCat.Children = mergeCategory(existingCat.Children, categoryMap, child)
		}
	} else {
		// Category doesn't exist - create a deep copy and add it
		newCat := &scenario.Category{
			Name:        pluginCat.Name,
			Description: pluginCat.Description,
			Descriptors: append([]*scenario.Descriptor(nil), pluginCat.Descriptors...),
			Children:    deepCopyCategories(pluginCat.Children),
		}
		result = append(result, newCat)
		categoryMap[newCat.Name] = newCat

		// Also add children to the map
		var addChildrenToMap func(cats []*scenario.Category)
		addChildrenToMap = func(cats []*scenario.Category) {
			for _, cat := range cats {
				categoryMap[cat.Name] = cat
				addChildrenToMap(cat.Children)
			}
		}
		addChildrenToMap(newCat.Children)
	}

	return result
}
