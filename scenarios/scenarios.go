package scenarios

import (
	"github.com/ethpandaops/goomy-blob/scenariotypes"

	"github.com/ethpandaops/goomy-blob/scenarios/combined"
	"github.com/ethpandaops/goomy-blob/scenarios/conflicting"
	"github.com/ethpandaops/goomy-blob/scenarios/deploytx"
	"github.com/ethpandaops/goomy-blob/scenarios/eoatx"
	"github.com/ethpandaops/goomy-blob/scenarios/erctx"
	"github.com/ethpandaops/goomy-blob/scenarios/normal"
	"github.com/ethpandaops/goomy-blob/scenarios/replacements"
	"github.com/ethpandaops/goomy-blob/scenarios/wallets"
)

var Scenarios map[string]func() scenariotypes.Scenario = map[string]func() scenariotypes.Scenario{
	"combined":     combined.NewScenario,
	"conflicting":  conflicting.NewScenario,
	"normal":       normal.NewScenario,
	"replacements": replacements.NewScenario,

	"eoatx":    eoatx.NewScenario,
	"erctx":    erctx.NewScenario,
	"deploytx": deploytx.NewScenario,

	"wallets": wallets.NewScenario,
}
