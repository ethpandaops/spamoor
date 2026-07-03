package scenario

import (
	"slices"
	"testing"
)

type stubPlugin struct {
	name string
}

func (p *stubPlugin) GetName() string        { return p.name }
func (p *stubPlugin) GetDescription() string { return "stub plugin" }
func (p *stubPlugin) AddRunning()            {}
func (p *stubPlugin) RemoveRunning()         {}

func newTestRegistry() *Registry {
	return NewRegistry([]*Descriptor{
		{Name: "eoatx", Aliases: []string{"transfer-tx"}, Description: "native eoa scenario"},
		{Name: "erctx", Description: "native erc scenario"},
	})
}

func pluginEntry(name string, aliases ...string) *ScenarioEntry {
	return &ScenarioEntry{
		Descriptor: &Descriptor{Name: name, Aliases: aliases, Description: "plugin scenario"},
		Source:     ScenarioSourcePlugin,
		Plugin:     &stubPlugin{name: "test-plugin"},
	}
}

func TestRegistryGet(t *testing.T) {
	r := newTestRegistry()

	if entry := r.Get("eoatx"); entry == nil || entry.Descriptor.Name != "eoatx" {
		t.Fatalf("expected to find scenario by name, got %v", entry)
	}

	if entry := r.Get("transfer-tx"); entry == nil || entry.Descriptor.Name != "eoatx" {
		t.Fatalf("expected to find scenario by alias, got %v", entry)
	}

	if entry := r.Get("unknown"); entry != nil {
		t.Fatalf("expected nil for unknown scenario, got %v", entry)
	}
}

func TestRegistryRegister(t *testing.T) {
	r := newTestRegistry()

	// New plugin scenario registers without replacing anything
	old, err := r.Register(pluginEntry("plugin-scenario"))
	if err != nil {
		t.Fatalf("unexpected error registering plugin scenario: %v", err)
	}
	if old != nil {
		t.Fatalf("expected no replaced entry, got %v", old)
	}

	// Re-registering the same name replaces the entry and returns the old one
	replacement := pluginEntry("plugin-scenario")
	old, err = r.Register(replacement)
	if err != nil {
		t.Fatalf("unexpected error replacing plugin scenario: %v", err)
	}
	if old == nil || old.Descriptor.Name != "plugin-scenario" {
		t.Fatalf("expected replaced entry to be returned, got %v", old)
	}
	if got := r.Get("plugin-scenario"); got != replacement {
		t.Fatalf("expected registry to hold the replacement entry")
	}

	// Native scenarios cannot be overridden by name
	if _, err := r.Register(pluginEntry("eoatx")); err == nil {
		t.Fatal("expected error when overriding a native scenario")
	}

	// Native scenario aliases cannot be shadowed either
	if _, err := r.Register(pluginEntry("other-name", "transfer-tx")); err == nil {
		t.Fatal("expected error when a plugin alias collides with a native alias")
	}
}

func TestRegistryGetDescriptor(t *testing.T) {
	r := newTestRegistry()

	if desc := r.GetDescriptor("erctx"); desc == nil || desc.Name != "erctx" {
		t.Fatalf("expected descriptor for registered scenario, got %v", desc)
	}

	if desc := r.GetDescriptor("unknown"); desc != nil {
		t.Fatalf("expected nil descriptor for unknown scenario, got %v", desc)
	}
}

func TestRegistryGetAllAndNames(t *testing.T) {
	r := newTestRegistry()

	if _, err := r.Register(pluginEntry("plugin-scenario")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	descriptors := r.GetAll()
	if len(descriptors) != 3 {
		t.Fatalf("expected 3 descriptors, got %d", len(descriptors))
	}

	names := r.GetNames()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	for _, expected := range []string{"eoatx", "erctx", "plugin-scenario"} {
		if !slices.Contains(names, expected) {
			t.Fatalf("expected names to contain %q, got %v", expected, names)
		}
	}
}

func TestRegistryIsNative(t *testing.T) {
	r := newTestRegistry()

	if _, err := r.Register(pluginEntry("plugin-scenario")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !r.IsNative("eoatx") {
		t.Fatal("expected native scenario name to be native")
	}
	if !r.IsNative("transfer-tx") {
		t.Fatal("expected native scenario alias to be native")
	}
	if r.IsNative("plugin-scenario") {
		t.Fatal("expected plugin scenario to not be native")
	}
	if r.IsNative("unknown") {
		t.Fatal("expected unknown scenario to not be native")
	}
}

func TestRegistryRemove(t *testing.T) {
	r := newTestRegistry()

	if _, err := r.Remove("eoatx"); err == nil {
		t.Fatal("expected error when removing a native scenario")
	}

	entry := pluginEntry("plugin-scenario")
	if _, err := r.Register(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	old, err := r.Remove("plugin-scenario")
	if err != nil {
		t.Fatalf("unexpected error removing plugin scenario: %v", err)
	}
	if old != entry {
		t.Fatalf("expected removed entry to be returned, got %v", old)
	}
	if got := r.Get("plugin-scenario"); got != nil {
		t.Fatalf("expected scenario to be gone after removal, got %v", got)
	}

	// Removing an unknown scenario is a no-op
	old, err = r.Remove("unknown")
	if err != nil {
		t.Fatalf("unexpected error removing unknown scenario: %v", err)
	}
	if old != nil {
		t.Fatalf("expected nil entry for unknown scenario, got %v", old)
	}
}

func TestRegistryGetPluginScenarios(t *testing.T) {
	r := newTestRegistry()

	if entries := r.GetPluginScenarios(); len(entries) != 0 {
		t.Fatalf("expected no plugin scenarios initially, got %d", len(entries))
	}

	if _, err := r.Register(pluginEntry("plugin-scenario")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := r.GetPluginScenarios()
	if len(entries) != 1 || entries[0].Descriptor.Name != "plugin-scenario" {
		t.Fatalf("expected exactly the plugin scenario, got %v", entries)
	}
}
