// Package plugin provides dynamic plugin loading using Yaegi.
package plugin

import (
	"fmt"
	"slices"
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
type LoadedPlugin struct {
	Descriptor *scenario.PluginDescriptor
	Metadata   *PluginMetadata  // Metadata from plugin.yaml
	TempDir    string           // Base temp directory
	PluginPath string           // Path to plugin source (for scenario Options.PluginPath)
	SourceType PluginSourceType // How the plugin was loaded

	mu           sync.RWMutex
	scenarios    map[string]bool // scenario names from this plugin currently in ScenarioRegistry
	runningCount atomic.Int32    // number of running spammers using scenarios from this plugin
	cleanedUp    bool
}

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
		scenarios:  make(map[string]bool, len(descriptor.Scenarios)),
	}
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

// AddRunning increments the running spammer count.
func (p *LoadedPlugin) AddRunning() {
	p.runningCount.Add(1)
}

// RemoveRunning decrements the running spammer count.
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

// ================== SCENARIO REGISTRY ==================

// ScenarioSource indicates whether a scenario is native or from a plugin.
type ScenarioSource int

const (
	// ScenarioSourceNative indicates a built-in scenario.
	ScenarioSourceNative ScenarioSource = iota
	// ScenarioSourcePlugin indicates a scenario loaded from a plugin.
	ScenarioSourcePlugin
)

// ScenarioEntry wraps a scenario descriptor with metadata about its source.
type ScenarioEntry struct {
	Descriptor *scenario.Descriptor
	Source     ScenarioSource
	Plugin     *LoadedPlugin // nil for native scenarios
}

// ScenarioRegistry manages scenario registration with protection for native scenarios.
type ScenarioRegistry struct {
	mu          sync.RWMutex
	entries     map[string]*ScenarioEntry // keyed by scenario name
	nativeNames map[string]bool           // tracks native scenario names (cannot be overridden)
}

// NewScenarioRegistry creates a new ScenarioRegistry with the given native scenario descriptors.
func NewScenarioRegistry(nativeDescriptors []*scenario.Descriptor) *ScenarioRegistry {
	r := &ScenarioRegistry{
		entries:     make(map[string]*ScenarioEntry, len(nativeDescriptors)+16),
		nativeNames: make(map[string]bool, len(nativeDescriptors)),
	}

	// Register native scenarios
	for _, desc := range nativeDescriptors {
		entry := &ScenarioEntry{
			Descriptor: desc,
			Source:     ScenarioSourceNative,
			Plugin:     nil,
		}
		r.entries[desc.Name] = entry
		r.nativeNames[desc.Name] = true

		// Also mark aliases as native
		for _, alias := range desc.Aliases {
			r.nativeNames[alias] = true
		}
	}

	return r
}

// Register adds or replaces a scenario in the registry.
// Returns an error if attempting to override a native scenario.
// Returns the old entry if one was replaced.
func (r *ScenarioRegistry) Register(entry *ScenarioEntry) (*ScenarioEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := entry.Descriptor.Name

	// Check if trying to override a native scenario
	if r.nativeNames[name] {
		return nil, fmt.Errorf("cannot override native scenario: %s", name)
	}

	// Also check aliases
	for _, alias := range entry.Descriptor.Aliases {
		if r.nativeNames[alias] {
			return nil, fmt.Errorf("cannot override native scenario alias: %s", alias)
		}
	}

	old := r.entries[name]
	r.entries[name] = entry

	return old, nil
}

// Get retrieves a scenario entry by name or alias.
func (r *ScenarioRegistry) Get(name string) *ScenarioEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Direct name lookup
	if entry, ok := r.entries[name]; ok {
		return entry
	}

	// Search by alias
	for _, entry := range r.entries {
		if slices.Contains(entry.Descriptor.Aliases, name) {
			return entry
		}
	}

	return nil
}

// GetDescriptor retrieves just the scenario descriptor by name or alias.
func (r *ScenarioRegistry) GetDescriptor(name string) *scenario.Descriptor {
	entry := r.Get(name)
	if entry == nil {
		return nil
	}

	return entry.Descriptor
}

// GetAll returns all registered scenario descriptors.
func (r *ScenarioRegistry) GetAll() []*scenario.Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	descriptors := make([]*scenario.Descriptor, 0, len(r.entries))
	for _, entry := range r.entries {
		descriptors = append(descriptors, entry.Descriptor)
	}

	return descriptors
}

// GetNames returns all registered scenario names.
func (r *ScenarioRegistry) GetNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.entries))
	for name := range r.entries {
		names = append(names, name)
	}

	return names
}

// IsNative returns true if the scenario name is a native (built-in) scenario.
func (r *ScenarioRegistry) IsNative(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.nativeNames[name]
}

// Remove removes a scenario from the registry and returns the old entry.
// Returns an error if attempting to remove a native scenario.
func (r *ScenarioRegistry) Remove(name string) (*ScenarioEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.nativeNames[name] {
		return nil, fmt.Errorf("cannot remove native scenario: %s", name)
	}

	old := r.entries[name]
	delete(r.entries, name)

	return old, nil
}
