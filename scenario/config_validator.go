package scenario

import (
	"fmt"
	"reflect"
	"strings"

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
	Required     bool
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
		if fieldInfo, exists := cv.validFields[fieldName]; exists {
			// Validate field type and value
			if validationErr := cv.validateFieldValue(fieldName, fieldValue, fieldInfo); validationErr != nil {
				result.Valid = false
				result.Errors = append(result.Errors, validationErr.Error())
			}
		} else {
			// Field doesn't exist in valid fields - this is now an error that crashes
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

	// Check all required fields are present
	for fieldName, fieldInfo := range cv.validFields {
		if fieldInfo.Required {
			if _, exists := config[fieldName]; !exists {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("Required field '%s' is missing", fieldName))
			}
		}
	}

	return result
}

// validateFieldValue validates a specific field value against its expected type
func (cv *ConfigValidator) validateFieldValue(fieldName string, value interface{}, fieldInfo FieldInfo) error {
	// Type validation would go here - for now just basic checks
	switch fieldInfo.Type.Kind() {
	case reflect.Uint64:
		if _, ok := value.(int); !ok {
			if _, ok := value.(uint64); !ok {
				return fmt.Errorf("field '%s' must be a positive integer", fieldName)
			}
		}
	case reflect.String:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field '%s' must be a string", fieldName)
		}
	case reflect.Bool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field '%s' must be a boolean", fieldName)
		}
	}

	return nil
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
	if descriptor == nil {
		return make(map[string]FieldInfo)
	}

	return ExtractFieldsFromStruct(descriptor.DefaultOptions)
}

// ExtractFieldsFromStruct uses reflection to extract field information from a struct
func ExtractFieldsFromStruct(structValue interface{}) map[string]FieldInfo {
	fields := make(map[string]FieldInfo)

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
			Required:     false, // Generally config fields are optional with defaults
			DefaultValue: fieldValue.Interface(),
			Description:  field.Name, // Could be enhanced with struct tags for descriptions
		}
	}

	// Add common wallet configuration fields that are valid for all scenarios
	walletFields := getCommonWalletFields()
	for key, value := range walletFields {
		fields[key] = value
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

// getCommonWalletFields returns wallet configuration fields that are valid for all scenarios
func getCommonWalletFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"seed":            {Type: reflect.TypeOf(""), Required: false, Description: "Seed for deterministic wallet generation"},
		"refill_amount":   {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Amount of ETH to fund/refill each child wallet with"},
		"refill_balance":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Min amount of ETH each child wallet should hold before refilling"},
		"refill_interval": {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Interval for child wallet balance check and refilling if needed (in sec)"},
		"wallet_count":    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of child wallets to create"},
	}
}
