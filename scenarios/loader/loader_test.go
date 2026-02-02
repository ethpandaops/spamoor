package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoadScenarioFromFile(t *testing.T) {
	// Find the example scenario file
	// Go up from scenarios/loader to repo root, then into examples
	examplePath := filepath.Join("..", "..", "examples", "dynamic-scenarios", "simple_transfer.go")

	// Check if file exists
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skipf("Example scenario file not found at %s - skipping test", examplePath)
	}

	// Create loader
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	loader := NewScenarioLoader(logger)

	// Load the scenario
	desc, err := loader.LoadFromFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to load scenario: %v", err)
	}

	// Verify descriptor
	if desc == nil {
		t.Fatal("Descriptor is nil")
	}

	if desc.Name != "dyn-transfer" {
		t.Errorf("Expected name 'dyn-transfer', got '%s'", desc.Name)
	}

	if desc.Description == "" {
		t.Error("Description should not be empty")
	}

	if desc.NewScenario == nil {
		t.Error("NewScenario function should not be nil")
	}

	// Try to create a scenario instance
	scenarioInstance := desc.NewScenario(logger)
	if scenarioInstance == nil {
		t.Error("NewScenario returned nil")
	}

	t.Logf("Successfully loaded scenario: %s - %s", desc.Name, desc.Description)
}

func TestLoadFromDir(t *testing.T) {
	// Find the examples directory
	exampleDir := filepath.Join("..", "..", "examples", "dynamic-scenarios")

	// Check if directory exists
	if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
		t.Skipf("Example directory not found at %s - skipping test", exampleDir)
	}

	// Create loader
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	loader := NewScenarioLoader(logger)

	// Load scenarios from directory
	descriptors := loader.LoadFromDir(exampleDir)

	if len(descriptors) == 0 {
		t.Error("Expected at least one scenario to be loaded")
	}

	for _, desc := range descriptors {
		t.Logf("Loaded: %s - %s", desc.Name, desc.Description)
	}
}

func TestNewInterpreter(t *testing.T) {
	logger := logrus.New()
	loader := NewScenarioLoader(logger)

	// Just verify we can create an interpreter without panicking
	interp := loader.newInterpreter()
	if interp == nil {
		t.Error("Interpreter should not be nil")
	}
}
