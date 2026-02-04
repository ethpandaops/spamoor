package scenarios

import (
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/loader"
	"github.com/sirupsen/logrus"

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
	// scenarioMu protects ScenarioDescriptors during hot-reload operations
	scenarioMu sync.RWMutex

	// externalScenarios tracks dynamically loaded scenarios for reload support
	externalScenarios []*scenario.Descriptor

	// externalDir is the last directory used for loading external scenarios
	externalDir string
)

// init automatically loads external scenarios from scenarios/external/ at startup
func init() {
	loadExternalScenariosOnStartup()
}

// loadExternalScenariosOnStartup attempts to load external scenarios from the default directory.
// This runs during package initialization and logs warnings if the directory exists but cannot be read.
func loadExternalScenariosOnStartup() {
	// Try to find external scenarios directory
	dir := findExternalDir()
	if dir == "" {
		return
	}

	logger := logrus.StandardLogger()
	l := loader.NewScenarioLoader(logger)

	// Load from each subdirectory (eoatx/, erc20_bloater/, etc.)
	entries, err := os.ReadDir(dir)
	if err != nil {
		logger.Warnf("failed to read external scenarios directory %s: %v", dir, err)
		return
	}

	var newScenarios []*scenario.Descriptor
	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			descriptors := l.LoadFromDir(subdir)
			newScenarios = append(newScenarios, descriptors...)
		}
	}

	if len(newScenarios) > 0 {
		scenarioMu.Lock()
		externalScenarios = newScenarios
		externalDir = dir
		ScenarioDescriptors = append(ScenarioDescriptors, newScenarios...)
		scenarioMu.Unlock()
		logger.Infof("auto-loaded %d external scenario(s) from %s", len(newScenarios), dir)
	}
}

// findExternalDir attempts to locate the external scenarios directory.
func findExternalDir() string {
	// Check relative to working directory
	candidates := []string{
		"scenarios/external",
		"./scenarios/external",
	}

	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}

	return ""
}

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

// ReloadExternalScenarios reloads all dynamic scenarios from the specified directory.
// This function is safe to call while spamoor is running (thread-safe).
// It removes previously loaded external scenarios and loads new ones from the directory.
// Returns the number of scenarios loaded and any error encountered.
func ReloadExternalScenarios(dir string, logger logrus.FieldLogger) (int, error) {
	l := loader.NewScenarioLoader(logger)

	// Load from each subdirectory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	var newScenarios []*scenario.Descriptor
	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			descriptors := l.LoadFromDir(subdir)
			newScenarios = append(newScenarios, descriptors...)
		}
	}

	// Lock and update the scenario list
	scenarioMu.Lock()
	defer scenarioMu.Unlock()

	// Build a set of external scenario names for quick lookup
	externalNames := make(map[string]bool)
	for _, ext := range externalScenarios {
		externalNames[ext.Name] = true
	}

	// Filter out old external scenarios from the main list
	filtered := make([]*scenario.Descriptor, 0, len(ScenarioDescriptors))
	for _, desc := range ScenarioDescriptors {
		if !externalNames[desc.Name] {
			filtered = append(filtered, desc)
		}
	}

	// Add new external scenarios
	ScenarioDescriptors = append(filtered, newScenarios...)
	externalScenarios = newScenarios
	externalDir = dir

	logger.Infof("reloaded %d external scenario(s) from %s", len(newScenarios), dir)
	return len(newScenarios), nil
}

// GetExternalDir returns the directory used for external scenarios.
func GetExternalDir() string {
	scenarioMu.RLock()
	defer scenarioMu.RUnlock()
	return externalDir
}

// LoadDynamicScenarios loads scenarios from Go source files in the specified directory
// and adds them to the ScenarioDescriptors list. This enables runtime scenario loading
// without recompilation. This function is thread-safe.
func LoadDynamicScenarios(dir string, logger logrus.FieldLogger) {
	l := loader.NewScenarioLoader(logger)
	dynamicScenarios := l.LoadFromDir(dir)

	scenarioMu.Lock()
	defer scenarioMu.Unlock()
	ScenarioDescriptors = append(ScenarioDescriptors, dynamicScenarios...)
}

// LoadDynamicScenarioFromFile loads a single scenario from a Go source file
// and adds it to the ScenarioDescriptors list. This function is thread-safe.
func LoadDynamicScenarioFromFile(path string, logger logrus.FieldLogger) error {
	l := loader.NewScenarioLoader(logger)
	desc, err := l.LoadFromFile(path)
	if err != nil {
		return err
	}

	scenarioMu.Lock()
	defer scenarioMu.Unlock()
	ScenarioDescriptors = append(ScenarioDescriptors, desc)
	return nil
}
