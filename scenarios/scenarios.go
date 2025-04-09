package scenarios

import (
	"github.com/ethpandaops/spamoor/scenariotypes"

	blobcombined "github.com/ethpandaops/spamoor/scenarios/blob-combined"
	blobconflicting "github.com/ethpandaops/spamoor/scenarios/blob-conflicting"
	blobreplacements "github.com/ethpandaops/spamoor/scenarios/blob-replacements"
	"github.com/ethpandaops/spamoor/scenarios/blobs"
	deploydestruct "github.com/ethpandaops/spamoor/scenarios/deploy-destruct"
	"github.com/ethpandaops/spamoor/scenarios/deploytx"
	"github.com/ethpandaops/spamoor/scenarios/eoatx"
	"github.com/ethpandaops/spamoor/scenarios/erctx"
	"github.com/ethpandaops/spamoor/scenarios/gasburnertx"
	mevswaps "github.com/ethpandaops/spamoor/scenarios/mev-swaps"
	"github.com/ethpandaops/spamoor/scenarios/setcodetx"
	"github.com/ethpandaops/spamoor/scenarios/wallets"
)

var ScenarioDescriptors = []*scenariotypes.ScenarioDescriptor{
	&blobcombined.ScenarioDescriptor,
	&blobconflicting.ScenarioDescriptor,
	&blobs.ScenarioDescriptor,
	&blobreplacements.ScenarioDescriptor,
	&deploydestruct.ScenarioDescriptor,
	&deploytx.ScenarioDescriptor,
	&eoatx.ScenarioDescriptor,
	&erctx.ScenarioDescriptor,
	&gasburnertx.ScenarioDescriptor,
	&mevswaps.ScenarioDescriptor,
	&setcodetx.ScenarioDescriptor,
	&wallets.ScenarioDescriptor,
}

func GetScenario(name string) *scenariotypes.ScenarioDescriptor {
	for _, scenario := range ScenarioDescriptors {
		if scenario.Name == name {
			return scenario
		}
	}

	return nil
}

func GetScenarioNames() []string {
	names := make([]string, len(ScenarioDescriptors))
	for i, scenario := range ScenarioDescriptors {
		names[i] = scenario.Name
	}
	return names
}
