package scenarios

import (
	"github.com/ethpandaops/spamoor/scenario"

	blobcombined "github.com/ethpandaops/spamoor/scenarios/blob-combined"
	blobconflicting "github.com/ethpandaops/spamoor/scenarios/blob-conflicting"
	blobreplacements "github.com/ethpandaops/spamoor/scenarios/blob-replacements"
	"github.com/ethpandaops/spamoor/scenarios/blobs"
	"github.com/ethpandaops/spamoor/scenarios/calltx"
	deploydestruct "github.com/ethpandaops/spamoor/scenarios/deploy-destruct"
	"github.com/ethpandaops/spamoor/scenarios/deploytx"
	"github.com/ethpandaops/spamoor/scenarios/eoatx"
	"github.com/ethpandaops/spamoor/scenarios/erctx"
	"github.com/ethpandaops/spamoor/scenarios/factorydeploytx"
	"github.com/ethpandaops/spamoor/scenarios/gasburnertx"
	"github.com/ethpandaops/spamoor/scenarios/geastx"
	"github.com/ethpandaops/spamoor/scenarios/setcodetx"
	contractdeploy "github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy"
	"github.com/ethpandaops/spamoor/scenarios/storagespam"
	uniswapswaps "github.com/ethpandaops/spamoor/scenarios/uniswap-swaps"
	"github.com/ethpandaops/spamoor/scenarios/wallets"
	"github.com/ethpandaops/spamoor/scenarios/xentoken"
)

// ScenarioDescriptors contains all available scenario descriptors for the spamoor tool.
// This registry includes scenarios for testing various Ethereum transaction types and patterns.
// Each descriptor defines the configuration, constructor, and metadata for a specific test scenario.
var ScenarioDescriptors = []*scenario.Descriptor{
	&blobcombined.ScenarioDescriptor,
	&blobconflicting.ScenarioDescriptor,
	&blobs.ScenarioDescriptor,
	&blobreplacements.ScenarioDescriptor,
	&calltx.ScenarioDescriptor,
	&deploydestruct.ScenarioDescriptor,
	&deploytx.ScenarioDescriptor,
	&eoatx.ScenarioDescriptor,
	&erctx.ScenarioDescriptor,
	&factorydeploytx.ScenarioDescriptor,
	&gasburnertx.ScenarioDescriptor,
	&geastx.ScenarioDescriptor,
	&setcodetx.ScenarioDescriptor,
	&storagespam.ScenarioDescriptor,
	&uniswapswaps.ScenarioDescriptor,
	&wallets.ScenarioDescriptor,
	&contractdeploy.ScenarioDescriptor,
	&xentoken.ScenarioDescriptor,
}

// GetScenario finds and returns a scenario descriptor by name.
// It performs a linear search through all registered scenarios and returns
// the matching descriptor, or nil if no scenario with the given name exists.
func GetScenario(name string) *scenario.Descriptor {
	for _, scenario := range ScenarioDescriptors {
		if scenario.Name == name {
			return scenario
		}
	}

	return nil
}

// GetScenarioNames returns a slice containing the names of all registered scenarios.
// This is useful for CLI help text, validation, and displaying available options
// to users. The order matches the order in ScenarioDescriptors.
func GetScenarioNames() []string {
	names := make([]string, len(ScenarioDescriptors))
	for i, scenario := range ScenarioDescriptors {
		names[i] = scenario.Name
	}
	return names
}
