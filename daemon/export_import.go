package daemon

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/scenarios"
)

// ExportSpammerConfig is an alias for scenario.SpammerConfig to maintain API compatibility
type ExportSpammerConfig = configs.SpammerConfig

// ImportItem is an alias for scenario.ConfigImportItem to maintain API compatibility
type ImportItem = configs.ConfigImportItem

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
func (d *Daemon) ImportSpammers(input string, userEmail string) (*ImportResult, error) {
	// Resolve all includes and get the final spammer configs
	importConfigs, err := configs.ResolveConfigImports(input, "", make(map[string]bool))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve import configs: %w", err)
	}

	// Validate the resolved data
	validation, err := d.validateImportConfigs(importConfigs)
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

	imported := 0
	var importedSpammers []ImportedSpammerInfo
	var importErrors []string
	var importWarnings []string

	for _, importConfig := range importConfigs {
		// Skip invalid scenarios (already validated)
		scenarioDescriptor := scenarios.GetScenario(importConfig.Scenario)
		if scenarioDescriptor == nil {
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
		configYAML, err := configs.MergeScenarioConfiguration(scenarioDescriptor, importConfig.Config)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to merge config for spammer '%s': %v", finalName, err)
			importErrors = append(importErrors, errMsg)
			continue
		}

		// Create the spammer (never start immediately for safety, pass isImport=true)
		spammer, err := d.NewSpammer(
			importConfig.Scenario,
			configYAML,
			finalName,
			importConfig.Description,
			false,
			userEmail,
			true, // isImport=true
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
			Start:       importConfig.Start,
		})
	}

	result := &ImportResult{
		ImportedCount: imported,
		Validation:    validation,
		Imported:      importedSpammers,
		Errors:        importErrors,
		Warnings:      importWarnings,
		Message:       fmt.Sprintf("Successfully imported %d out of %d spammers", imported, validation.TotalSpammers),
	}

	// Audit log the import
	if d.auditLogger != nil && userEmail != "" {
		skippedCount := validation.TotalSpammers - validation.ValidSpammers
		err := d.auditLogger.LogSpammersImport(userEmail, imported, skippedCount, input)
		if err != nil {
			d.logger.Errorf("Failed to create audit log for spammers import: %v", err)
		}
	}

	return result, nil
}

// validateImportConfigs validates spammer configurations
func (d *Daemon) validateImportConfigs(importConfigs []ExportSpammerConfig) (*ImportValidationResult, error) {
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
		scenarioDescriptor := scenarios.GetScenario(importConfig.Scenario)
		if scenarioDescriptor == nil {
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
	Start       *bool  `json:"start,omitempty"`
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
