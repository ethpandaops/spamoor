package scenarios

import (
	"github.com/ethpandaops/spamoor/scenarios/gasburnertx"
	"github.com/ethpandaops/spamoor/scenariotypes"

	blobcombined "github.com/ethpandaops/spamoor/scenarios/blob-combined"
	blobconflicting "github.com/ethpandaops/spamoor/scenarios/blob-conflicting"
	blobreplacements "github.com/ethpandaops/spamoor/scenarios/blob-replacements"
	"github.com/ethpandaops/spamoor/scenarios/blobs"
	"github.com/ethpandaops/spamoor/scenarios/deploytx"
	"github.com/ethpandaops/spamoor/scenarios/eoatx"
	"github.com/ethpandaops/spamoor/scenarios/erctx"
	"github.com/ethpandaops/spamoor/scenarios/wallets"
)

var Scenarios map[string]func() scenariotypes.Scenario = map[string]func() scenariotypes.Scenario{
	"blob-combined":     blobcombined.NewScenario,
	"blob-conflicting":  blobconflicting.NewScenario,
	"blobs":             blobs.NewScenario,
	"blob-replacements": blobreplacements.NewScenario,

	"eoatx":       eoatx.NewScenario,
	"erctx":       erctx.NewScenario,
	"deploytx":    deploytx.NewScenario,
	"gasburnertx": gasburnertx.NewScenario,

	"wallets": wallets.NewScenario,
}
