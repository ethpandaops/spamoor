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
	"github.com/ethpandaops/spamoor/scenarios/setcodetx"
	contractdeploy "github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy"
	eoadelegation "github.com/ethpandaops/spamoor/scenarios/statebloat/eoa_delegation"

	//erc20maxtransfers "github.com/ethpandaops/spamoor/scenarios/statebloat/erc20_max_transfers"
	//extcodesizeoverload "github.com/ethpandaops/spamoor/scenarios/statebloat/extcodesize-overload"
	randsstorebloater "github.com/ethpandaops/spamoor/scenarios/statebloat/rand_sstore_bloater"
	uniswapswaps "github.com/ethpandaops/spamoor/scenarios/uniswap-swaps"
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
	&setcodetx.ScenarioDescriptor,
	&randsstorebloater.ScenarioDescriptor,
	&uniswapswaps.ScenarioDescriptor,
	&wallets.ScenarioDescriptor,
	&contractdeploy.ScenarioDescriptor,
	&eoadelegation.ScenarioDescriptor,
	//&erc20maxtransfers.ScenarioDescriptor,
	//&extcodesizeoverload.ScenarioDescriptor,
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
