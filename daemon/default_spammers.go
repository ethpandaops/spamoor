package daemon

import (
	"embed"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/scenario"
)

// defaultSpammersFS embeds the built-in default spammer definitions. Each file holds a
// list of spammer/group configs in the regular export/import format; files are imported
// in lexical filename order on first launch.
//
//go:embed default-spammers/*.yaml
var defaultSpammersFS embed.FS

// DefaultSpammerConfig is a built-in default spammer definition: a regular import
// config plus a short technical key. The key is a stable machine identifier for the
// --startup-defaults flag; the display name/description may change over time, the key
// must not.
type DefaultSpammerConfig struct {
	Key                   string `yaml:"key"`
	configs.SpammerConfig `yaml:",inline"`
}

// LoadDefaultSpammerConfigs parses the embedded default spammer definitions and returns
// them in import order (lexical filename order, entries in file order).
func LoadDefaultSpammerConfigs() ([]DefaultSpammerConfig, error) {
	entries, err := defaultSpammersFS.ReadDir("default-spammers")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded default spammers: %w", err)
	}

	allConfigs := make([]DefaultSpammerConfig, 0, len(entries)*8)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := defaultSpammersFS.ReadFile("default-spammers/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read default spammer file %s: %w", entry.Name(), err)
		}

		var fileConfigs []DefaultSpammerConfig
		if err := yaml.Unmarshal(data, &fileConfigs); err != nil {
			return nil, fmt.Errorf("failed to parse default spammer file %s: %w", entry.Name(), err)
		}

		allConfigs = append(allConfigs, fileConfigs...)
	}

	return allConfigs, nil
}

// ImportDefaultSpammers inserts the built-in default spammers and spammer groups into
// the database. All defaults are created paused; the ones listed in autoStart (matched
// by their technical key, falling back to their display name) are started afterwards.
// Intended to be called once on the daemon's first launch, after the spammer state has
// been restored.
func (d *Daemon) ImportDefaultSpammers(autoStart []string, logger logrus.FieldLogger) error {
	defaultConfigs, err := LoadDefaultSpammerConfigs()
	if err != nil {
		return err
	}

	importConfigs := make([]configs.SpammerConfig, 0, len(defaultConfigs))
	for _, cfg := range defaultConfigs {
		importConfigs = append(importConfigs, cfg.SpammerConfig)
	}

	result, err := d.importSpammerConfigs(importConfigs, "system", d.computeDefaultSpammerIDs(defaultConfigs))
	if err != nil {
		return fmt.Errorf("failed to import default spammers: %w", err)
	}

	for _, importError := range result.Errors {
		logger.Warnf("default spammer import error: %s", importError)
	}
	for _, warning := range result.Warnings {
		logger.Infof("default spammer import warning: %s", warning)
	}

	logger.Infof("inserted %d default spammers", result.ImportedCount)

	// Resolve the requested technical keys to display names, then start by name.
	if len(autoStart) > 0 {
		d.StartSpammersByName(resolveDefaultStartNames(defaultConfigs, autoStart), logger)
	}

	return nil
}

// computeDefaultSpammerIDs assigns the reserved ids (< 100) to the default spammers in
// definition order: regular entries take the next free sequential id starting at 1,
// group rows snap to the next free multiple of 10 and their members continue from
// there. Ids already taken by existing spammers are skipped; entries that would
// overflow the reserved range get no explicit id and fall back to the runtime
// allocator (ids >= 100).
func (d *Daemon) computeDefaultSpammerIDs(defaultConfigs []DefaultSpammerConfig) map[string]int64 {
	used := make(map[int64]bool, len(defaultConfigs))
	for _, s := range d.GetAllSpammers() {
		used[s.GetID()] = true
	}

	ids := make(map[string]int64, len(defaultConfigs))
	next := int64(1)
	for _, cfg := range defaultConfigs {
		var id int64
		if cfg.Scenario == scenario.GroupScenarioName {
			id = (next + 9) / 10 * 10
			for used[id] {
				id += 10
			}
		} else {
			id = next
			for used[id] {
				id++
			}
		}

		if id >= 100 {
			continue
		}

		used[id] = true
		ids[cfg.Name] = id
		next = id + 1
	}

	return ids
}

// resolveDefaultStartNames maps --startup-defaults values to spammer display names.
// Each value is matched against the defaults' technical keys first; values that match
// no key are kept as-is and treated as display names.
func resolveDefaultStartNames(defaultConfigs []DefaultSpammerConfig, autoStart []string) []string {
	nameByKey := make(map[string]string, len(defaultConfigs))
	for _, cfg := range defaultConfigs {
		if cfg.Key != "" {
			nameByKey[cfg.Key] = cfg.Name
		}
	}

	startNames := make([]string, 0, len(autoStart))
	for _, key := range autoStart {
		if name, ok := nameByKey[key]; ok {
			startNames = append(startNames, name)
		} else {
			startNames = append(startNames, key)
		}
	}

	return startNames
}

// StartSpammersByName starts all spammers or spammer groups whose configured name
// matches one of the given names. Starting a group starts all of its enabled members.
// Unknown names and failed starts are logged but do not abort the remaining names.
func (d *Daemon) StartSpammersByName(names []string, logger logrus.FieldLogger) {
	if len(names) == 0 {
		return
	}

	spammers := d.GetAllSpammers()
	for _, name := range names {
		found := false
		for _, spammer := range spammers {
			if spammer.GetName() != name {
				continue
			}

			found = true
			if err := spammer.Start(); err != nil {
				logger.Errorf("failed to start spammer %d (%s): %v", spammer.GetID(), name, err)
			} else {
				logger.Infof("started spammer %d: %s (%s)", spammer.GetID(), name, spammer.GetScenario())
			}
		}

		if !found {
			logger.Warnf("no spammer or spammer group named %q found to start", name)
		}
	}
}
