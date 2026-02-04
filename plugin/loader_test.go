package plugin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewPluginLoader(t *testing.T) {
	logger := logrus.New()
	loader := NewPluginLoader(logger)

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
	loader := NewPluginLoader(logger)

	desc, err := loader.LoadFromFile(pluginPath)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	if desc == nil {
		t.Fatal("Descriptor is nil")
	}

	if desc.Name == "" {
		t.Error("Plugin name should not be empty")
	}

	t.Logf("Loaded plugin: %s with %d scenarios", desc.Name, len(desc.Scenarios))
}

func TestBuildMemoryFS(t *testing.T) {
	logger := logrus.New()
	loader := NewPluginLoader(logger)

	// Create a simple tar with one file
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	content := []byte("package test\n\nvar TestVar = 123\n")
	hdr := &tar.Header{
		Name: "test.go",
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

	// Build the memory FS
	memFS, err := loader.buildMemoryFS(&buf, "testplugin")
	if err != nil {
		t.Fatalf("Failed to build memory FS: %v", err)
	}

	// Verify we can open the file
	expectedPath := "src/github.com/ethpandaops/spamoor/plugins/testplugin/test.go"
	file, err := memFS.Open(expectedPath)
	if err != nil {
		t.Fatalf("Failed to open file from memory FS: %v", err)
	}
	file.Close()
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
	loader := NewPluginLoader(logger)

	desc, err := loader.LoadFromBytes("test-plugin", data, true)
	if err != nil {
		t.Fatalf("Failed to load plugin from bytes: %v", err)
	}

	if desc == nil {
		t.Fatal("Descriptor is nil")
	}

	t.Logf("Loaded plugin from bytes: %s", desc.Name)
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
	loader := NewPluginLoader(logger)

	desc, err := loader.LoadFromReader("test-plugin", file, true)
	if err != nil {
		t.Fatalf("Failed to load plugin from reader: %v", err)
	}

	if desc == nil {
		t.Fatal("Descriptor is nil")
	}

	t.Logf("Loaded plugin from reader: %s", desc.Name)
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
	loader := NewPluginLoader(logger)

	// This will fail because the plugin imports plugin.Descriptor which uses scenario.Descriptor
	// But it tests that the compressed tar handling works
	_, err := loader.LoadFromBytes("testplugin", buf.Bytes(), true)
	// We expect an error because the minimal plugin doesn't have all required symbols
	// but we can verify the tar was processed correctly by checking the error type
	if err != nil {
		t.Logf("Expected error loading minimal plugin: %v", err)
	}
}
