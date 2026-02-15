package scenario

import (
	"fmt"
	"slices"
	"sync"
)

// ScenarioPlugin represents the minimal interface needed to track plugin state.
// This interface is implemented by plugin.LoadedPlugin to avoid circular dependencies.
type ScenarioPlugin interface {
	GetName() string
	GetDescription() string
	AddRunning()
	RemoveRunning()
}

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
	Descriptor *Descriptor
	Source     ScenarioSource
	Plugin     ScenarioPlugin // nil for native scenarios
}

// Registry manages scenario registration with protection for native scenarios.
type Registry struct {
	mu          sync.RWMutex
	entries     map[string]*ScenarioEntry // keyed by scenario name
	nativeNames map[string]bool           // tracks native scenario names (cannot be overridden)
}

// NewRegistry creates a new Registry with the given native scenario descriptors.
func NewRegistry(nativeDescriptors []*Descriptor) *Registry {
	r := &Registry{
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
func (r *Registry) Register(entry *ScenarioEntry) (*ScenarioEntry, error) {
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
func (r *Registry) Get(name string) *ScenarioEntry {
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
func (r *Registry) GetDescriptor(name string) *Descriptor {
	entry := r.Get(name)
	if entry == nil {
		return nil
	}

	return entry.Descriptor
}

// GetAll returns all registered scenario descriptors.
func (r *Registry) GetAll() []*Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	descriptors := make([]*Descriptor, 0, len(r.entries))
	for _, entry := range r.entries {
		descriptors = append(descriptors, entry.Descriptor)
	}

	return descriptors
}

// GetNames returns all registered scenario names.
func (r *Registry) GetNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.entries))
	for name := range r.entries {
		names = append(names, name)
	}

	return names
}

// IsNative returns true if the scenario name is a native (built-in) scenario.
func (r *Registry) IsNative(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.nativeNames[name]
}

// Remove removes a scenario from the registry and returns the old entry.
// Returns an error if attempting to remove a native scenario.
func (r *Registry) Remove(name string) (*ScenarioEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.nativeNames[name] {
		return nil, fmt.Errorf("cannot remove native scenario: %s", name)
	}

	old := r.entries[name]
	delete(r.entries, name)

	return old, nil
}

// GetPluginScenarios returns all scenario entries that are from plugins.
func (r *Registry) GetPluginScenarios() []*ScenarioEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make([]*ScenarioEntry, 0, 8)
	for _, entry := range r.entries {
		if entry.Source == ScenarioSourcePlugin {
			entries = append(entries, entry)
		}
	}

	return entries
}
