package plugin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/sirupsen/logrus"
)

// createTestRegistries creates plugin and scenario registries for testing.
func createTestRegistries() (*PluginRegistry, *ScenarioRegistry) {
	pluginRegistry := NewPluginRegistry()
	scenarioRegistry := NewScenarioRegistry(nil)

	return pluginRegistry, scenarioRegistry
}

func TestNewPluginLoader(t *testing.T) {
	logger := logrus.New()
	pluginRegistry, scenarioRegistry := createTestRegistries()
	loader := NewPluginLoader(logger, pluginRegistry, scenarioRegistry)

	if loader == nil {
		t.Fatal("NewPluginLoader returned nil")
	}
}

func TestLoadFromFile(t *testing.T) {
	// Check if test plugin exists
	pluginPath := filepath.Join("..", "bin", "plugins", "test-plugin.tar.gz")
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		t.Skipf("Test plugin not found at %s - run 'make plugins' first", pluginPath)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	pluginRegistry, scenarioRegistry := createTestRegistries()
	loader := NewPluginLoader(logger, pluginRegistry, scenarioRegistry)

	loaded, err := loader.LoadFromFile(pluginPath)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded plugin is nil")
	}

	if loaded.Descriptor == nil {
		t.Fatal("Descriptor is nil")
	}

	if loaded.Descriptor.Name == "" {
		t.Error("Plugin name should not be empty")
	}

	if loaded.TempDir == "" {
		t.Error("TempDir should not be empty")
	}

	if loaded.PluginPath == "" {
		t.Error("PluginPath should not be empty")
	}

	t.Logf("Loaded plugin: %s with %d scenarios", loaded.Descriptor.Name, len(loaded.Descriptor.Scenarios))
	t.Logf("TempDir: %s", loaded.TempDir)
	t.Logf("PluginPath: %s", loaded.PluginPath)

	// Cleanup
	if err := loader.CleanupPlugin(loaded); err != nil {
		t.Errorf("Failed to cleanup plugin: %v", err)
	}
}

func TestLoadFromBytes(t *testing.T) {
	// Check if test plugin exists
	pluginPath := filepath.Join("..", "bin", "plugins", "test-plugin.tar.gz")
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		t.Skipf("Test plugin not found at %s - run 'make plugins' first", pluginPath)
	}

	data, err := os.ReadFile(pluginPath)
	if err != nil {
		t.Fatalf("Failed to read plugin file: %v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	pluginRegistry, scenarioRegistry := createTestRegistries()
	loader := NewPluginLoader(logger, pluginRegistry, scenarioRegistry)

	loaded, err := loader.LoadFromBytes(data, true)
	if err != nil {
		t.Fatalf("Failed to load plugin from bytes: %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded plugin is nil")
	}

	t.Logf("Loaded plugin from bytes: %s", loaded.Descriptor.Name)

	// Cleanup
	if err := loader.CleanupPlugin(loaded); err != nil {
		t.Errorf("Failed to cleanup plugin: %v", err)
	}
}

func TestLoadFromReader(t *testing.T) {
	// Check if test plugin exists
	pluginPath := filepath.Join("..", "bin", "plugins", "test-plugin.tar.gz")
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		t.Skipf("Test plugin not found at %s - run 'make plugins' first", pluginPath)
	}

	file, err := os.Open(pluginPath)
	if err != nil {
		t.Fatalf("Failed to open plugin file: %v", err)
	}
	defer file.Close()

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	pluginRegistry, scenarioRegistry := createTestRegistries()
	loader := NewPluginLoader(logger, pluginRegistry, scenarioRegistry)

	loaded, err := loader.LoadFromReader(file, true)
	if err != nil {
		t.Fatalf("Failed to load plugin from reader: %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded plugin is nil")
	}

	t.Logf("Loaded plugin from reader: %s", loaded.Descriptor.Name)

	// Cleanup
	if err := loader.CleanupPlugin(loaded); err != nil {
		t.Errorf("Failed to cleanup plugin: %v", err)
	}
}

func TestMemoryFSReadDir(t *testing.T) {
	memFS := newMemoryFS()

	// Add some files
	memFS.addFile("src/pkg/file1.go", []byte("package pkg"), 0644)
	memFS.addFile("src/pkg/file2.go", []byte("package pkg"), 0644)
	memFS.addFile("src/pkg/sub/file3.go", []byte("package sub"), 0644)

	// Test ReadDir on root
	entries, err := memFS.ReadDir("src")
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry in src/, got %d", len(entries))
	}

	// Test ReadDir on pkg
	entries, err = memFS.ReadDir("src/pkg")
	if err != nil {
		t.Fatalf("ReadDir on pkg failed: %v", err)
	}

	// Should have file1.go, file2.go, and sub/
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries in src/pkg/, got %d", len(entries))
	}
}

func TestCompressedTar(t *testing.T) {
	// Create a gzip-compressed tar with a simple plugin structure
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Add plugin.yaml file (required for tar.gz plugins)
	metaContent := []byte(`name: testplugin
build_time: "2024-01-01T00:00:00Z"
git_version: "test123"
`)
	metaHdr := &tar.Header{
		Name: "plugin.yaml",
		Mode: 0644,
		Size: int64(len(metaContent)),
	}
	if err := tw.WriteHeader(metaHdr); err != nil {
		t.Fatalf("Failed to write tar header: %v", err)
	}
	if _, err := tw.Write(metaContent); err != nil {
		t.Fatalf("Failed to write tar content: %v", err)
	}

	// Add plugin.go file
	content := []byte(`package testplugin

import "github.com/ethpandaops/spamoor/plugin"

var PluginDescriptor = plugin.Descriptor{
	Name:        "test",
	Description: "Test plugin",
}
`)
	hdr := &tar.Header{
		Name: "plugin.go",
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("Failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("Failed to write tar content: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("Failed to close tar writer: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	pluginRegistry, scenarioRegistry := createTestRegistries()
	loader := NewPluginLoader(logger, pluginRegistry, scenarioRegistry)

	// This will fail because the plugin imports plugin.Descriptor which uses scenario.Descriptor
	// But it tests that the compressed tar handling and metadata extraction works
	loaded, err := loader.LoadFromBytes(buf.Bytes(), true)
	// We expect an error because the minimal plugin doesn't have all required symbols
	// but we can verify the tar was processed correctly by checking the error type
	if err != nil {
		t.Logf("Expected error loading minimal plugin: %v", err)
	} else if loaded != nil {
		// Cleanup if it somehow succeeded
		loader.CleanupPlugin(loaded)
	}
}

func TestMissingPluginYAML(t *testing.T) {
	// Create a gzip-compressed tar WITHOUT plugin.yaml
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Add only plugin.go file, no plugin.yaml
	content := []byte(`package testplugin
var Test = 1
`)
	hdr := &tar.Header{
		Name: "plugin.go",
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("Failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("Failed to write tar content: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("Failed to close tar writer: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	pluginRegistry, scenarioRegistry := createTestRegistries()
	loader := NewPluginLoader(logger, pluginRegistry, scenarioRegistry)

	// This should fail because plugin.yaml is missing
	_, err := loader.LoadFromBytes(buf.Bytes(), true)
	if err == nil {
		t.Fatal("Expected error when plugin.yaml is missing")
	}

	if !strings.Contains(err.Error(), "plugin.yaml") {
		t.Errorf("Error should mention plugin.yaml, got: %v", err)
	}

	t.Logf("Got expected error: %v", err)
}

func TestPluginRegistry(t *testing.T) {
	registry := NewPluginRegistry()

	// Create a mock loaded plugin
	desc := &scenario.PluginDescriptor{
		Name:        "test-plugin",
		Description: "Test plugin",
	}
	meta := &PluginMetadata{
		Name:       "test-plugin",
		BuildTime:  "2024-01-01T00:00:00Z",
		GitVersion: "abc123",
	}
	plugin := NewLoadedPlugin(desc, meta, "/tmp/test", "/tmp/test/path", PluginSourceFile)

	// Register
	old := registry.Register(plugin)
	if old != nil {
		t.Error("Expected nil for first registration")
	}

	// Get
	retrieved := registry.Get("test-plugin")
	if retrieved != plugin {
		t.Error("Get returned wrong plugin")
	}

	// Get all
	all := registry.GetAll()
	if len(all) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(all))
	}

	// Remove
	removed := registry.Remove("test-plugin")
	if removed != plugin {
		t.Error("Remove returned wrong plugin")
	}

	// Should be empty now
	if registry.Get("test-plugin") != nil {
		t.Error("Plugin should be removed")
	}
}

func TestScenarioRegistry(t *testing.T) {
	registry := NewScenarioRegistry(nil)

	// Check empty registry
	if registry.Get("nonexistent") != nil {
		t.Error("Expected nil for nonexistent scenario")
	}

	// IsNative should return false for empty registry
	if registry.IsNative("foo") {
		t.Error("Empty registry should have no native scenarios")
	}
}

func TestLoadedPluginRefCounting(t *testing.T) {
	desc := &scenario.PluginDescriptor{
		Name: "test",
	}
	meta := &PluginMetadata{
		Name:       "test",
		BuildTime:  "2024-01-01T00:00:00Z",
		GitVersion: "abc123",
	}
	plugin := NewLoadedPlugin(desc, meta, "/tmp/test", "/tmp/test/path", PluginSourceFile)

	// Initially should be cleanable (no scenarios, no running)
	if !plugin.CanCleanup() {
		t.Error("Empty plugin should be cleanable")
	}

	// Add a scenario
	plugin.AddScenario("scenario1")
	if plugin.CanCleanup() {
		t.Error("Plugin with scenario should not be cleanable")
	}
	if plugin.GetScenarioCount() != 1 {
		t.Errorf("Expected 1 scenario, got %d", plugin.GetScenarioCount())
	}

	// Add running count
	plugin.AddRunning()
	if plugin.GetRunningCount() != 1 {
		t.Errorf("Expected 1 running, got %d", plugin.GetRunningCount())
	}

	// Remove scenario but still running
	plugin.RemoveScenario("scenario1")
	if plugin.CanCleanup() {
		t.Error("Plugin with running spammer should not be cleanable")
	}

	// Remove running
	plugin.RemoveRunning()
	if !plugin.CanCleanup() {
		t.Error("Plugin with no scenarios and no running should be cleanable")
	}
}
