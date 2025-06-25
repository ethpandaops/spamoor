package scenario

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ConfigValidator provides validation for scenario configurations
type ConfigValidator struct {
	scenarioName string
	validFields  map[string]FieldInfo
	logger       logrus.FieldLogger
}

// FieldInfo contains metadata about a valid configuration field
type FieldInfo struct {
	Type         reflect.Type
	DefaultValue interface{}
	Description  string
}

// ValidationResult contains the results of configuration validation
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// NewConfigValidator creates a new configuration validator for a specific scenario
func NewConfigValidator(scenarioName string, validFields map[string]FieldInfo, logger logrus.FieldLogger) *ConfigValidator {
	return &ConfigValidator{
		scenarioName: scenarioName,
		validFields:  validFields,
		logger:       logger,
	}
}

// ValidateConfig validates a YAML configuration against the scenario's valid fields
func (cv *ConfigValidator) ValidateConfig(configYAML string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	if configYAML == "" {
		return result // Empty config is valid
	}

	// Parse the YAML configuration
	var config map[string]interface{}
	if err := yaml.Unmarshal([]byte(configYAML), &config); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid YAML configuration: %v", err))
		return result
	}

	// Collect all invalid fields
	var invalidFields []string
	var suggestions []string

	// Check each field in the provided configuration
	for fieldName, fieldValue := range config {
		if _, exists := cv.validFields[fieldName]; !exists {
			// Field doesn't exist in valid fields
			invalidFields = append(invalidFields, fmt.Sprintf("'%s' (value: %v)", fieldName, fieldValue))

			// Check for common typos (underscore vs dash)
			if suggestion := cv.suggestCorrectField(fieldName); suggestion != "" {
				suggestions = append(suggestions, fmt.Sprintf("'%s' -> '%s'", fieldName, suggestion))
			}
		}
	}

	// Convert invalid fields to errors instead of warnings
	if len(invalidFields) > 0 {
		result.Valid = false
		errorMessage := fmt.Sprintf("invalid configuration fields detected for scenario '%s': %s",
			cv.scenarioName, strings.Join(invalidFields, ", "))

		if len(suggestions) > 0 {
			errorMessage += fmt.Sprintf(" - possible corrections: %s", strings.Join(suggestions, ", "))
		}

		result.Errors = append(result.Errors, errorMessage)
	}

	return result
}

// suggestCorrectField suggests a correct field name based on common typos
func (cv *ConfigValidator) suggestCorrectField(invalidField string) string {
	// Common transformations to check
	transformations := []func(string) string{
		func(s string) string { return strings.ReplaceAll(s, "-", "_") }, // dash to underscore
		func(s string) string { return strings.ReplaceAll(s, "_", "-") }, // underscore to dash
		strings.ToLower, // case insensitive
	}

	for _, transform := range transformations {
		transformed := transform(invalidField)
		if _, exists := cv.validFields[transformed]; exists {
			return transformed
		}
	}

	// If no direct transformation works, find the closest match
	// This is a simple approach - could be enhanced with edit distance algorithms
	invalidLower := strings.ToLower(invalidField)
	for validField := range cv.validFields {
		validLower := strings.ToLower(validField)
		if strings.Contains(validLower, invalidLower) || strings.Contains(invalidLower, validLower) {
			return validField
		}
	}

	return ""
}

// GetScenarioValidFields returns the valid fields for a scenario by extracting them from the scenario descriptor
func GetScenarioValidFields(descriptor *Descriptor) map[string]FieldInfo {
	var fields map[string]FieldInfo

	fields = extractFieldsFromStruct(descriptor.DefaultOptions, fields)
	fields = extractFieldsFromStruct(&spamoor.WalletPoolConfig{}, fields)

	return fields
}

// extractFieldsFromStruct uses reflection to extract field information from a struct
func extractFieldsFromStruct(structValue interface{}, fields map[string]FieldInfo) map[string]FieldInfo {
	if fields == nil {
		fields = make(map[string]FieldInfo)
	}

	val := reflect.ValueOf(structValue)
	typ := reflect.TypeOf(structValue)

	// Handle pointers
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get the yaml tag name
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" || yamlTag == "-" {
			continue
		}

		// Parse yaml tag (could be "name,omitempty" etc)
		yamlName := strings.Split(yamlTag, ",")[0]
		if yamlName == "" {
			yamlName = strings.ToLower(field.Name)
		}

		fields[yamlName] = FieldInfo{
			Type:         field.Type,
			DefaultValue: fieldValue.Interface(),
			Description:  field.Tag.Get("usage"),
		}
	}

	return fields
}

// ValidateScenarioConfig validates configuration for a scenario using reflection-based field extraction
func ValidateScenarioConfig(descriptor *Descriptor, configYAML string, logger logrus.FieldLogger) error {
	if descriptor == nil {
		return fmt.Errorf("scenario descriptor is nil")
	}

	validFields := GetScenarioValidFields(descriptor)
	validator := NewConfigValidator(descriptor.Name, validFields, logger)

	validationResult := validator.ValidateConfig(configYAML)
	if !validationResult.Valid {
		for _, err := range validationResult.Errors {
			logger.Errorf("Configuration validation error: %s", err)
		}
		return fmt.Errorf("configuration validation failed")
	}

	return nil
}

// ParseAndValidateConfig is a generalized helper that validates and parses YAML config for any scenario
func ParseAndValidateConfig(descriptor *Descriptor, configYAML string, target interface{}, logger logrus.FieldLogger) error {
	// Validate the configuration first
	if err := ValidateScenarioConfig(descriptor, configYAML, logger); err != nil {
		return err
	}

	// Parse the YAML into the target struct
	if configYAML != "" {
		if err := yaml.Unmarshal([]byte(configYAML), target); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	return nil
}
