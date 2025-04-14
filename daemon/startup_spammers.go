package daemon

import (
	"fmt"
	"os"

	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// StartupSpammerConfig represents a single spammer configuration in the startup config file
type StartupSpammerConfig struct {
	Scenario string                 `yaml:"scenario"`
	Name     string                 `yaml:"name"`
	Config   map[string]interface{} `yaml:"config"`
}

// LoadStartupSpammers loads the startup spammers configuration from a file
func (d *Daemon) LoadStartupSpammers(configFile string, logger logrus.FieldLogger) ([]StartupSpammerConfig, error) {
	if configFile == "" {
		return nil, nil
	}

	logger.Infof("loading startup spammers from %s", configFile)

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read startup config file: %w", err)
	}

	var spammers []StartupSpammerConfig
	err = yaml.Unmarshal(data, &spammers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse startup config file: %w", err)
	}

	return spammers, nil
}

// AddStartupSpammers adds the startup spammers to the daemon
func (d *Daemon) AddStartupSpammers(spammers []StartupSpammerConfig) error {
	for _, spammerConfig := range spammers {
		scenario := scenarios.GetScenario(spammerConfig.Scenario)
		if scenario == nil {
			return fmt.Errorf("scenario not found: %s", spammerConfig.Scenario)
		}

		defaultYaml, err := yaml.Marshal(scenario.DefaultOptions)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}

		defaultWalletConfig := spamoor.GetDefaultWalletConfig(scenario.Name)
		defaultWalletConfigYaml, err := yaml.Marshal(defaultWalletConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal default wallet config: %w", err)
		}

		// merge default config with spammer config
		mergedConfig := map[string]interface{}{}

		err = yaml.Unmarshal(defaultWalletConfigYaml, &mergedConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal default config: %w", err)
		}

		err = yaml.Unmarshal(defaultYaml, &mergedConfig)
		if err != nil {
			return fmt.Errorf("failed to unmarshal default config: %w", err)
		}

		for k, v := range spammerConfig.Config {
			mergedConfig[k] = v
		}

		// Convert config to YAML string
		configYAML, err := yaml.Marshal(mergedConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal spammer config: %w", err)
		}

		// Create a name if not provided
		name := spammerConfig.Name
		if name == "" {
			name = fmt.Sprintf("Startup %s", spammerConfig.Scenario)
		}

		// Create the spammer
		_, err = d.NewSpammer(
			spammerConfig.Scenario,
			string(configYAML),
			name,
			"Created from startup configuration",
			true,
		)
		if err != nil {
			return fmt.Errorf("failed to create startup spammer: %w", err)
		}
	}

	return nil
}
