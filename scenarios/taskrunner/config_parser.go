package taskrunner

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// TasksConfig represents the complete task configuration
type TasksConfig struct {
	InitTasks      []Task `yaml:"init" json:"init"`
	ExecutionTasks []Task `yaml:"execution" json:"execution"`
}

// RawTasksConfig represents the raw configuration before parsing tasks
type RawTasksConfig struct {
	InitTasks      []*TaskConfig `yaml:"init" json:"init"`
	ExecutionTasks []*TaskConfig `yaml:"execution" json:"execution"`
}

// ParseTasksConfig parses task configuration from YAML or JSON data
func ParseTasksConfig(data []byte) (*TasksConfig, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty configuration data")
	}

	// Parse raw configuration first
	rawConfig, err := parseRawConfig(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Convert to task instances
	config := &TasksConfig{}

	// Parse init tasks
	for i, taskConfig := range rawConfig.InitTasks {
		task, err := CreateTask(taskConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create init task %d (%s): %w", i, taskConfig.Type, err)
		}

		if err := task.Validate(); err != nil {
			return nil, fmt.Errorf("init task %d (%s) validation failed: %w", i, taskConfig.Type, err)
		}

		config.InitTasks = append(config.InitTasks, task)
	}

	// Parse execution tasks
	for i, taskConfig := range rawConfig.ExecutionTasks {
		task, err := CreateTask(taskConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create execution task %d (%s): %w", i, taskConfig.Type, err)
		}

		if err := task.Validate(); err != nil {
			return nil, fmt.Errorf("execution task %d (%s) validation failed: %w", i, taskConfig.Type, err)
		}

		config.ExecutionTasks = append(config.ExecutionTasks, task)
	}

	// Validate configuration requirements
	if err := validateTasksConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// parseRawConfig parses the raw configuration data (YAML or JSON)
func parseRawConfig(data []byte) (*RawTasksConfig, error) {
	var config RawTasksConfig

	// Detect format and parse accordingly
	if isJSON(data) {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse as JSON: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse as YAML: %w", err)
		}
	}

	return &config, nil
}

// isJSON checks if the data appears to be JSON format
func isJSON(data []byte) bool {
	trimmed := strings.TrimSpace(string(data))
	return strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
}

// validateTasksConfig validates the overall task configuration
func validateTasksConfig(config *TasksConfig) error {
	// Must have at least one execution task
	if len(config.ExecutionTasks) == 0 {
		return fmt.Errorf("configuration must have at least one execution task")
	}

	// Validate task names are unique within each phase
	if err := validateTaskNames(config.InitTasks, "init"); err != nil {
		return err
	}

	if err := validateTaskNames(config.ExecutionTasks, "execution"); err != nil {
		return err
	}

	// Validate contract references - ensure referenced contracts exist
	if err := validateContractReferences(config); err != nil {
		return err
	}

	return nil
}

// validateTaskNames ensures task names are unique within a phase
func validateTaskNames(tasks []Task, phase string) error {
	names := make(map[string]int)

	for i, task := range tasks {
		name := task.GetName()
		if name == "" {
			continue // Unnamed tasks are allowed
		}

		if existing, exists := names[name]; exists {
			return fmt.Errorf("duplicate task name '%s' in %s phase (tasks %d and %d)",
				name, phase, existing, i)
		}

		names[name] = i
	}

	return nil
}

// validateContractReferences validates that all contract references can potentially be resolved
func validateContractReferences(config *TasksConfig) error {
	// Collect all named tasks from init phase
	initContracts := make(map[string]bool)
	for _, task := range config.InitTasks {
		if name := task.GetName(); name != "" && task.GetType() == "deploy" {
			initContracts[name] = true
		}
	}

	// Validate execution task references
	execContracts := make(map[string]bool)
	for i, task := range config.ExecutionTasks {
		// Check references in this task
		refs := extractContractReferences(task)

		for _, ref := range refs {
			// Check if reference exists in init contracts or previous execution contracts
			if !initContracts[ref] && !execContracts[ref] {
				return fmt.Errorf("execution task %d references unknown contract '%s'", i, ref)
			}
		}

		// Add this task's deployed contracts to execution contracts
		if name := task.GetName(); name != "" && task.GetType() == "deploy" {
			execContracts[name] = true
		}
	}

	return nil
}

// extractContractReferences extracts contract references from a task
func extractContractReferences(task Task) []string {
	var refs []string

	switch t := task.(type) {
	case *CallTask:
		if ref := extractContractRef(t.Target); ref != "" {
			refs = append(refs, ref)
		}

		// Check arguments for contract references
		for _, arg := range t.CallArgs {
			if argStr, ok := arg.(string); ok {
				if ref := extractContractRef(argStr); ref != "" {
					refs = append(refs, ref)
				}
			}
		}
	}

	return refs
}

// extractContractRef extracts contract name from a reference string like {contract:name} or {contract:name:nonce}
func extractContractRef(str string) string {
	if strings.HasPrefix(str, "{contract:") && strings.HasSuffix(str, "}") {
		content := str[10 : len(str)-1]
		// Handle both {contract:name} and {contract:name:nonce} formats
		// Split by colon and take the first part as the contract name
		parts := strings.Split(content, ":")
		if len(parts) >= 1 {
			return parts[0]
		}
	}
	return ""
}

// ValidateConfigString validates a configuration string without full parsing
func ValidateConfigString(configStr string) error {
	if strings.TrimSpace(configStr) == "" {
		return fmt.Errorf("empty configuration string")
	}

	// Try to parse as basic structure
	_, err := parseRawConfig([]byte(configStr))
	return err
}
