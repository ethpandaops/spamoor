package scenarios

import (
	"slices"
	"sync"

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
	// scenarioMu protects ScenarioDescriptors during registration operations
	scenarioMu sync.RWMutex
)

// ScenarioDescriptors contains all available scenario descriptors for the spamoor tool.
// This registry includes scenarios for testing various Ethereum transaction types and patterns.
// Each descriptor defines the configuration, constructor, and metadata for a specific test scenario.
var ScenarioDescriptors = []*scenario.Descriptor{
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

// GetScenario finds and returns a scenario descriptor by name.
// It performs a linear search through all registered scenarios and returns
// the matching descriptor, or nil if no scenario with the given name exists.
// This function is thread-safe.
func GetScenario(name string) *scenario.Descriptor {
	scenarioMu.RLock()
	defer scenarioMu.RUnlock()

	for _, s := range ScenarioDescriptors {
		if s.Name == name {
			return s
		}
		if len(s.Aliases) > 0 && slices.Contains(s.Aliases, name) {
			return s
		}
	}

	return nil
}

// GetScenarioNames returns a slice containing the names of all registered scenarios.
// This is useful for CLI help text, validation, and displaying available options
// to users. The order matches the order in ScenarioDescriptors.
// This function is thread-safe.
func GetScenarioNames() []string {
	scenarioMu.RLock()
	defer scenarioMu.RUnlock()

	names := make([]string, len(ScenarioDescriptors))
	for i, s := range ScenarioDescriptors {
		names[i] = s.Name
	}
	return names
}

// RegisterScenario adds a scenario descriptor to the global registry.
// This function is thread-safe and can be used to register scenarios
// loaded from plugins or other external sources.
func RegisterScenario(desc *scenario.Descriptor) {
	scenarioMu.Lock()
	defer scenarioMu.Unlock()
	ScenarioDescriptors = append(ScenarioDescriptors, desc)
}
