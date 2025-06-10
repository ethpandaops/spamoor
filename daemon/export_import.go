package daemon

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"gopkg.in/yaml.v3"
)

// ExportSpammerConfig represents a spammer configuration for export/import.
// This uses the same format as StartupSpammerConfig to maintain compatibility.
type ExportSpammerConfig struct {
	Scenario    string                 `yaml:"scenario"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Config      map[string]interface{} `yaml:"config"`
}

// ExportSpammers exports the specified spammer IDs to YAML format.
// If no IDs are provided, exports all spammers.
// Returns the YAML string representation of the spammers.
func (d *Daemon) ExportSpammers(spammerIDs ...int64) (string, error) {
	var spammersToExport []*Spammer

	if len(spammerIDs) == 0 {
		// Export all spammers
		spammersToExport = d.GetAllSpammers()
	} else {
		// Export specific spammers
		for _, id := range spammerIDs {
			spammer := d.GetSpammer(id)
			if spammer == nil {
				return "", fmt.Errorf("spammer with ID %d not found", id)
			}
			spammersToExport = append(spammersToExport, spammer)
		}
	}

	if len(spammersToExport) == 0 {
		return "", fmt.Errorf("no spammers to export")
	}

	var exportConfigs []ExportSpammerConfig
	for _, spammer := range spammersToExport {
		// Parse the existing config to extract only the custom parts
		var spammerConfig map[string]interface{}
		if err := yaml.Unmarshal([]byte(spammer.GetConfig()), &spammerConfig); err != nil {
			return "", fmt.Errorf("failed to parse config for spammer %d: %w", spammer.GetID(), err)
		}

		exportConfig := ExportSpammerConfig{
			Scenario:    spammer.GetScenario(),
			Name:        spammer.GetName(),
			Description: spammer.GetDescription(),
			Config:      spammerConfig,
		}

		exportConfigs = append(exportConfigs, exportConfig)
	}

	yamlData, err := yaml.Marshal(exportConfigs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal export data: %w", err)
	}

	return string(yamlData), nil
}

// ImportSpammers imports spammers from YAML data, file path, or URL.
// Handles deduplication by checking name combinations and validates before importing.
// Returns validation results and the number of spammers imported.
func (d *Daemon) ImportSpammers(input string) (*ImportResult, error) {
	var yamlData string
	var err error

	// Check if input is a URL
	if isURL(input) {
		yamlData, err = d.downloadFromURL(input)
		if err != nil {
			return nil, fmt.Errorf("failed to download from URL: %w", err)
		}
	} else if isFilePath(input) {
		// Try to read as file path
		yamlData, err = d.readFromFile(input)
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
	} else {
		// Treat as raw YAML data
		yamlData = input
	}

	// Validate the YAML data first
	validation, err := d.validateImportData(yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to validate import data: %w", err)
	}

	if validation.ValidSpammers == 0 {
		return &ImportResult{
			ImportedCount: 0,
			Validation:    validation,
			Message:       "No valid spammers found to import",
		}, nil
	}

	// Parse the validated data
	var importConfigs []ExportSpammerConfig
	if err := yaml.Unmarshal([]byte(yamlData), &importConfigs); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	imported := 0
	var importedSpammers []ImportedSpammerInfo
	var importErrors []string
	var importWarnings []string

	for _, importConfig := range importConfigs {
		// Skip invalid scenarios (already validated)
		scenario := scenarios.GetScenario(importConfig.Scenario)
		if scenario == nil {
			importErrors = append(importErrors, fmt.Sprintf("Scenario '%s' not found", importConfig.Scenario))
			continue
		}

		// Check if spammer with same scenario + name already exists - skip if it does
		if d.spammerExistsByScenarioAndName(importConfig.Scenario, importConfig.Name) {
			importWarnings = append(importWarnings, fmt.Sprintf("Skipped spammer '%s' (%s) - already exists", importConfig.Name, importConfig.Scenario))
			continue
		}

		finalName := importConfig.Name

		// Merge default configuration with imported config
		configYAML, err := d.mergeConfiguration(scenario, importConfig.Config)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to merge config for spammer '%s': %v", finalName, err)
			importErrors = append(importErrors, errMsg)
			continue
		}

		// Create the spammer (never start immediately for safety)
		spammer, err := d.NewSpammer(
			importConfig.Scenario,
			configYAML,
			finalName,
			importConfig.Description,
			false,
		)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to create spammer '%s': %v", finalName, err)
			importErrors = append(importErrors, errMsg)
			continue
		}

		imported++
		importedSpammers = append(importedSpammers, ImportedSpammerInfo{
			ID:          spammer.GetID(),
			Name:        finalName,
			Scenario:    importConfig.Scenario,
			Description: importConfig.Description,
		})
	}

	return &ImportResult{
		ImportedCount: imported,
		Validation:    validation,
		Imported:      importedSpammers,
		Errors:        importErrors,
		Warnings:      importWarnings,
		Message:       fmt.Sprintf("Successfully imported %d out of %d spammers", imported, validation.TotalSpammers),
	}, nil
}

// isURL checks if the input string is a valid URL
func isURL(input string) bool {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

// isFilePath checks if the input string looks like a file path
func isFilePath(input string) bool {
	// Check if it contains YAML content markers (likely raw YAML)
	if len(input) > 0 && (input[0] == '-' || input[0] == '[' || input[0] == '{') {
		return false
	}

	// Check if file exists
	if _, err := os.Stat(input); err == nil {
		return true
	}

	// Check if it looks like a path (contains / or \)
	return len(input) > 0 && (input[0] == '/' || input[0] == '.' || input[0] == '~' ||
		(len(input) > 1 && input[1] == ':')) // Windows drive letter
}

// readFromFile reads YAML data from a local file
func (d *Daemon) readFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(data), nil
}

// downloadFromURL downloads YAML data from a remote URL
func (d *Daemon) downloadFromURL(urlStr string) (string, error) {
	// Validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	// Download the YAML data
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to download from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	yamlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(yamlData), nil
}

// spammerExistsByScenarioAndName checks if a spammer with the given scenario and name already exists
func (d *Daemon) spammerExistsByScenarioAndName(scenario, name string) bool {
	existingSpammers := d.GetAllSpammers()
	for _, existing := range existingSpammers {
		if existing.GetScenario() == scenario && existing.GetName() == name {
			return true
		}
	}
	return false
}

// mergeConfiguration merges scenario defaults with imported configuration
func (d *Daemon) mergeConfiguration(scenario *scenariotypes.ScenarioDescriptor, importedConfig map[string]interface{}) (string, error) {
	// Get default configurations
	defaultYaml, err := yaml.Marshal(scenario.DefaultOptions)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default config: %w", err)
	}

	defaultWalletConfig := spamoor.GetDefaultWalletConfig(scenario.Name)
	defaultWalletConfigYaml, err := yaml.Marshal(defaultWalletConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default wallet config: %w", err)
	}

	// Merge configurations
	mergedConfig := map[string]interface{}{}

	if err := yaml.Unmarshal(defaultWalletConfigYaml, &mergedConfig); err != nil {
		return "", fmt.Errorf("failed to unmarshal default wallet config: %w", err)
	}

	if err := yaml.Unmarshal(defaultYaml, &mergedConfig); err != nil {
		return "", fmt.Errorf("failed to unmarshal default config: %w", err)
	}

	// Apply imported config overrides
	for k, v := range importedConfig {
		mergedConfig[k] = v
	}

	configYAML, err := yaml.Marshal(mergedConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged config: %w", err)
	}

	return string(configYAML), nil
}

// validateImportData validates the YAML import data without actually importing.
// Returns information about what would be imported and any potential issues.
func (d *Daemon) validateImportData(yamlData string) (*ImportValidationResult, error) {
	var importConfigs []ExportSpammerConfig
	if err := yaml.Unmarshal([]byte(yamlData), &importConfigs); err != nil {
		return nil, fmt.Errorf("failed to parse import data: %w", err)
	}

	result := &ImportValidationResult{
		TotalSpammers:    len(importConfigs),
		ValidSpammers:    0,
		Duplicates:       []string{},
		InvalidScenarios: []string{},
		Spammers:         []SpammerValidationInfo{},
	}

	existingSpammers := d.GetAllSpammers()
	existingNames := make(map[string]bool)
	for _, existing := range existingSpammers {
		existingNames[existing.GetName()] = true
	}

	for _, importConfig := range importConfigs {
		info := SpammerValidationInfo{
			Name:        importConfig.Name,
			Scenario:    importConfig.Scenario,
			Description: importConfig.Description,
			Valid:       true,
			Issues:      []string{},
		}

		// Check scenario validity
		scenario := scenarios.GetScenario(importConfig.Scenario)
		if scenario == nil {
			info.Valid = false
			info.Issues = append(info.Issues, "Unknown scenario")
			result.InvalidScenarios = append(result.InvalidScenarios, importConfig.Scenario)
		}

		// Check name duplicates
		if existingNames[importConfig.Name] {
			info.Issues = append(info.Issues, "Name already exists, will be renamed")
			result.Duplicates = append(result.Duplicates, importConfig.Name)
		}

		if info.Valid {
			result.ValidSpammers++
		}

		result.Spammers = append(result.Spammers, info)
	}

	return result, nil
}

// ImportResult contains the results of an import operation
type ImportResult struct {
	ImportedCount int                     `json:"imported_count"`
	Validation    *ImportValidationResult `json:"validation"`
	Imported      []ImportedSpammerInfo   `json:"imported"`
	Errors        []string                `json:"errors"`
	Warnings      []string                `json:"warnings"`
	Message       string                  `json:"message"`
}

type ImportedSpammerInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Scenario    string `json:"scenario"`
	Description string `json:"description"`
}

// ImportValidationResult contains validation results for import data
type ImportValidationResult struct {
	TotalSpammers    int                     `json:"total_spammers"`
	ValidSpammers    int                     `json:"valid_spammers"`
	Duplicates       []string                `json:"duplicates"`
	InvalidScenarios []string                `json:"invalid_scenarios"`
	Spammers         []SpammerValidationInfo `json:"spammers"`
}

// SpammerValidationInfo contains validation info for a single spammer
type SpammerValidationInfo struct {
	Name        string   `json:"name"`
	Scenario    string   `json:"scenario"`
	Description string   `json:"description"`
	Valid       bool     `json:"valid"`
	Issues      []string `json:"issues"`
}
