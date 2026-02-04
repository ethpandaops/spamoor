package plugin

import "github.com/ethpandaops/spamoor/scenario"

type Descriptor struct {
	Name        string
	Description string
	Scenarios   []*scenario.Descriptor
}
