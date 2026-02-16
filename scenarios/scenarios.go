package scenarios

import (
	"fmt"
	"slices"

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
	extcodesizesetup "github.com/ethpandaops/spamoor/scenarios/statebloat/extcodesize_setup"
	storagetriebrancher "github.com/ethpandaops/spamoor/scenarios/statebloat/storage_trie_brancher"
	"github.com/ethpandaops/spamoor/scenarios/storagespam"
	"github.com/ethpandaops/spamoor/scenarios/taskrunner"
	uniswapswaps "github.com/ethpandaops/spamoor/scenarios/uniswap-swaps"
	"github.com/ethpandaops/spamoor/scenarios/wallets"
	"github.com/ethpandaops/spamoor/scenarios/xentoken"
)

// scenarioDescriptorTree contains the tree of scenario categories and descriptors.
// The order of the categories and descriptors is important for the CLI help text,
// validation, and displaying available options to users.
var scenarioDescriptorTree = []*scenario.Category{
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
			&extcodesizesetup.ScenarioDescriptor,
			&storagetriebrancher.ScenarioDescriptor,
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

// scenarioDescriptors contains all available scenario descriptors for the spamoor tool.
var scenarioDescriptors []*scenario.Descriptor

func init() {
	nameMap := make(map[string]*scenario.Descriptor)

	var addCategory func(category *scenario.Category)
	addCategory = func(category *scenario.Category) {
		for _, descriptor := range category.Descriptors {
			if _, ok := nameMap[descriptor.Name]; ok {
				panic(fmt.Sprintf("scenario descriptor %s already registered", descriptor.Name))
			}
			nameMap[descriptor.Name] = descriptor
		}
		scenarioDescriptors = append(scenarioDescriptors, category.Descriptors...)
		for _, child := range category.Children {
			addCategory(child)
		}
	}
	for _, category := range scenarioDescriptorTree {
		addCategory(category)
	}
}

// GetScenario finds and returns a scenario descriptor by name.
// It performs a linear search through all registered scenarios and returns
// the matching descriptor, or nil if no scenario with the given name exists.
func GetScenario(name string) *scenario.Descriptor {
	for _, scenario := range scenarioDescriptors {
		if scenario.Name == name {
			return scenario
		}
		if len(scenario.Aliases) > 0 && slices.Contains(scenario.Aliases, name) {
			return scenario
		}
	}

	return nil
}

// GetScenarioNames returns a slice containing the names of all registered scenarios.
// This is useful for CLI help text, validation, and displaying available options
// to users. The order matches the order in ScenarioDescriptors.
func GetScenarioNames() []string {
	names := make([]string, len(scenarioDescriptors))
	for i, scenario := range scenarioDescriptors {
		names[i] = scenario.Name
	}
	return names
}

// GetScenarioCategories returns a slice containing the categories of all registered scenarios.
// This is useful for CLI help text, validation, and displaying available options
// to users. The order matches the order in scenarioDescriptorTree.
func GetScenarioCategories() []*scenario.Category {
	return scenarioDescriptorTree
}
