package scenario

import (
	"context"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Category describes the category of a scenario.
type Category struct {
	Name        string
	Description string
	Descriptors []*Descriptor
	Children    []*Category
}

// Descriptor describes a scenario.
type Descriptor struct {
	Name           string
	Aliases        []string
	Description    string
	DefaultOptions any
	NewScenario    func(logger logrus.FieldLogger) Scenario
}

// PluginDescriptor describes a plugin.
type PluginDescriptor struct {
	Name        string
	Description string
	Categories  []*Category
}

// GetAllScenarios returns all scenario descriptors from all categories (flattened).
func (p *PluginDescriptor) GetAllScenarios() []*Descriptor {
	var scenarios []*Descriptor
	var collect func(cat *Category)
	collect = func(cat *Category) {
		scenarios = append(scenarios, cat.Descriptors...)
		for _, child := range cat.Children {
			collect(child)
		}
	}
	for _, cat := range p.Categories {
		collect(cat)
	}
	return scenarios
}

// Options contains the options for the scenario initialization.
type Options struct {
	WalletPool *spamoor.WalletPool
	Config     string
	GlobalCfg  map[string]any
	PluginPath string // Path to plugin resources (empty for native scenarios)
}

type Scenario interface {
	// Flags registers the scenario's flags with the given flag set.
	Flags(flags *pflag.FlagSet) error
	// Init initializes the scenario with the given options.
	Init(options *Options) error
	// Run runs the scenario.
	Run(ctx context.Context) error
}
