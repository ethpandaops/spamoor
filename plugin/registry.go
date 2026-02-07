// Package plugin provides dynamic plugin loading using Yaegi.
package plugin

import (
	"sync"
	"sync/atomic"

	"github.com/ethpandaops/spamoor/scenario"
)

// PluginSourceType indicates how a plugin was loaded.
type PluginSourceType int

const (
	// PluginSourceBytes indicates the plugin was loaded from raw bytes (e.g., tar.gz data).
	PluginSourceBytes PluginSourceType = iota
	// PluginSourceFile indicates the plugin was loaded from a tar.gz file path.
	PluginSourceFile
	// PluginSourceURL indicates the plugin was loaded from a remote URL.
	PluginSourceURL
	// PluginSourceLocal indicates the plugin was loaded from a local directory path.
	PluginSourceLocal
)

// String returns a string representation of the PluginSourceType.
func (p PluginSourceType) String() string {
	switch p {
	case PluginSourceBytes:
		return "bytes"
	case PluginSourceFile:
		return "file"
	case PluginSourceURL:
		return "url"
	case PluginSourceLocal:
		return "local"
	default:
		return "unknown"
	}
}

// LoadedPlugin tracks a loaded plugin with its temp directory and reference counting.
// It implements scenario.ScenarioPlugin interface to allow use in the scenario registry.
type LoadedPlugin struct {
	Descriptor *scenario.PluginDescriptor
	Metadata   *PluginMetadata  // Metadata from plugin.yaml
	TempDir    string           // Base temp directory
	PluginPath string           // Path to plugin source (for scenario Options.PluginPath)
	SourceType PluginSourceType // How the plugin was loaded

	mu           sync.RWMutex
	scenarios    map[string]bool // scenario names from this plugin currently in Registry
	runningCount atomic.Int32    // number of running spammers using scenarios from this plugin
	cleanedUp    bool
}

// Ensure LoadedPlugin implements scenario.ScenarioPlugin
var _ scenario.ScenarioPlugin = (*LoadedPlugin)(nil)

// NewLoadedPlugin creates a new LoadedPlugin instance.
func NewLoadedPlugin(
	descriptor *scenario.PluginDescriptor,
	metadata *PluginMetadata,
	tempDir, pluginPath string,
	sourceType PluginSourceType,
) *LoadedPlugin {
	return &LoadedPlugin{
		Descriptor: descriptor,
		Metadata:   metadata,
		TempDir:    tempDir,
		PluginPath: pluginPath,
		SourceType: sourceType,
		scenarios:  make(map[string]bool, len(descriptor.GetAllScenarios())),
	}
}

// GetName returns the plugin name (implements scenario.ScenarioPlugin).
func (p *LoadedPlugin) GetName() string {
	if p.Descriptor == nil {
		return ""
	}
	return p.Descriptor.Name
}

// GetDescription returns the plugin description (implements scenario.ScenarioPlugin).
func (p *LoadedPlugin) GetDescription() string {
	if p.Descriptor == nil {
		return ""
	}
	return p.Descriptor.Description
}

// AddScenario marks a scenario name as registered from this plugin.
func (p *LoadedPlugin) AddScenario(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.scenarios[name] = true
}

// RemoveScenario removes a scenario name from this plugin's tracking.
func (p *LoadedPlugin) RemoveScenario(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.scenarios, name)
}

// GetScenarioCount returns the number of scenarios still registered from this plugin.
func (p *LoadedPlugin) GetScenarioCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.scenarios)
}

// AddRunning increments the running spammer count (implements scenario.ScenarioPlugin).
func (p *LoadedPlugin) AddRunning() {
	p.runningCount.Add(1)
}

// RemoveRunning decrements the running spammer count (implements scenario.ScenarioPlugin).
func (p *LoadedPlugin) RemoveRunning() {
	p.runningCount.Add(-1)
}

// GetRunningCount returns the number of running spammers using this plugin.
func (p *LoadedPlugin) GetRunningCount() int32 {
	return p.runningCount.Load()
}

// CanCleanup returns true if the plugin can be cleaned up:
// all descriptors from this plugin have been replaced AND no running spammer uses them.
func (p *LoadedPlugin) CanCleanup() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.scenarios) == 0 && p.runningCount.Load() == 0
}

// IsCleanedUp returns true if the plugin has already been cleaned up.
func (p *LoadedPlugin) IsCleanedUp() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.cleanedUp
}

// MarkCleanedUp marks the plugin as cleaned up.
func (p *LoadedPlugin) MarkCleanedUp() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cleanedUp = true
}

// ================== PLUGIN REGISTRY ==================

// PluginRegistry tracks all loaded plugins.
type PluginRegistry struct {
	mu      sync.RWMutex
	plugins map[string]*LoadedPlugin // keyed by plugin name
}

// NewPluginRegistry creates a new PluginRegistry.
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]*LoadedPlugin, 8),
	}
}

// Register adds or replaces a plugin in the registry.
// Returns the old plugin if one was replaced.
func (r *PluginRegistry) Register(plugin *LoadedPlugin) *LoadedPlugin {
	r.mu.Lock()
	defer r.mu.Unlock()

	old := r.plugins[plugin.Descriptor.Name]
	r.plugins[plugin.Descriptor.Name] = plugin

	return old
}

// Get retrieves a plugin by name.
func (r *PluginRegistry) Get(name string) *LoadedPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.plugins[name]
}

// Remove removes a plugin from the registry and returns it.
func (r *PluginRegistry) Remove(name string) *LoadedPlugin {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin := r.plugins[name]
	delete(r.plugins, name)

	return plugin
}

// GetAll returns all loaded plugins.
func (r *PluginRegistry) GetAll() []*LoadedPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]*LoadedPlugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}

	return plugins
}
