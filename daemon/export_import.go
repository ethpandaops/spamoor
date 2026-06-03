package daemon

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/scenario"
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

	exportSet := make(map[int64]*Spammer, len(spammersToExport))
	for _, s := range spammersToExport {
		exportSet[s.GetID()] = s
	}

	// Pull in the parent group of any selected member so the import can re-link it even
	// when only the member was selected.
	for _, s := range spammersToExport {
		if s.GetGroupID() != 0 {
			if _, ok := exportSet[s.GetGroupID()]; !ok {
				if parent := d.GetSpammer(s.GetGroupID()); parent != nil {
					exportSet[parent.GetID()] = parent
				}
			}
		}
	}

	// Pull in all members of any group in the set (selected directly, or added above as a
	// parent) so exporting a group always carries its full membership.
	groupsInSet := make([]*Spammer, 0, len(exportSet))
	for _, s := range exportSet {
		if s.IsGroup() {
			groupsInSet = append(groupsInSet, s)
		}
	}
	for _, group := range groupsInSet {
		for _, member := range d.getGroupMembersFromMap(group.GetID()) {
			if _, ok := exportSet[member.GetID()]; !ok {
				exportSet[member.GetID()] = member
			}
		}
	}

	// Emit group rows first so the file is human-readable top-down; the importer is
	// order-independent regardless.
	ordered := make([]*Spammer, 0, len(exportSet))
	for _, s := range exportSet {
		ordered = append(ordered, s)
	}
	sort.SliceStable(ordered, func(a, b int) bool {
		if ordered[a].IsGroup() != ordered[b].IsGroup() {
			return ordered[a].IsGroup()
		}
		return ordered[a].GetID() < ordered[b].GetID()
	})

	var exportConfigs []ExportSpammerConfig
	for _, spammer := range ordered {
		exportConfig, err := d.spammerToExportConfig(spammer)
		if err != nil {
			return "", err
		}
		exportConfigs = append(exportConfigs, exportConfig)
	}

	yamlData, err := yaml.Marshal(exportConfigs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal export data: %w", err)
	}

	return string(yamlData), nil
}

// spammerToExportConfig converts a spammer (standalone, group, or member) into its
// export representation, including group overlay/totals or member weight metadata.
func (d *Daemon) spammerToExportConfig(spammer *Spammer) (ExportSpammerConfig, error) {
	// Parse into a yaml.Node so the config's key ordering and comments are preserved.
	configNode, err := configs.ParseConfigNode(spammer.GetConfig())
	if err != nil {
		return ExportSpammerConfig{}, fmt.Errorf("failed to parse config for spammer %d: %w", spammer.GetID(), err)
	}

	exportConfig := ExportSpammerConfig{
		Scenario:    spammer.GetScenario(),
		Name:        spammer.GetName(),
		Description: spammer.GetDescription(),
	}
	if configNode != nil {
		exportConfig.Config = *configNode
	}

	switch {
	case spammer.IsGroup():
		gc, err := configs.ParseGroupConfig(spammer.GetGroupConfig())
		if err != nil {
			return ExportSpammerConfig{}, err
		}
		exportConfig.GroupConfig = map[string]interface{}{
			"throughput_mode":   gc.ThroughputMode,
			"total_throughput":  gc.TotalThroughput,
			"total_count":       gc.TotalCount,
			"total_max_pending": gc.TotalMaxPending,
		}
	case spammer.GetGroupID() != 0:
		if parent := d.GetSpammer(spammer.GetGroupID()); parent != nil {
			exportConfig.Group = parent.GetName()
		}
		mc, err := configs.ParseMemberConfig(spammer.GetGroupConfig())
		if err == nil {
			exportConfig.GroupConfig = map[string]interface{}{
				"weight":     mc.Weight,
				"enabled":    mc.Enabled,
				"sort_order": mc.SortOrder,
			}
		}
	}

	return exportConfig, nil
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

	// Map group name -> group id, seeded with existing groups so members can link to
	// groups that already exist or are created during this import.
	groupIDByName := make(map[string]int64)
	for _, existing := range d.GetAllSpammers() {
		if existing.IsGroup() {
			groupIDByName[existing.GetName()] = existing.GetID()
		}
	}

	// Pass 1: create group rows first so members can be linked afterwards.
	for _, importConfig := range importConfigs {
		if importConfig.Scenario != scenario.GroupScenarioName {
			continue
		}

		if _, exists := groupIDByName[importConfig.Name]; exists {
			importWarnings = append(importWarnings, fmt.Sprintf("Skipped group '%s' - already exists", importConfig.Name))
			continue
		}

		// The group's config is a sparse overlay; marshal it verbatim (no defaults merge),
		// preserving its key order and comments.
		var overlayYAML string
		if importConfig.Config.Kind != 0 {
			data, err := yaml.Marshal(&importConfig.Config)
			if err != nil {
				importErrors = append(importErrors, fmt.Sprintf("Failed to marshal overlay for group '%s': %v", importConfig.Name, err))
				continue
			}
			overlayYAML = string(data)
		}

		groupCfg := configs.GroupConfigFromMap(importConfig.GroupConfig)
		group, err := d.NewGroup(importConfig.Name, importConfig.Description, overlayYAML, groupCfg, userEmail)
		if err != nil {
			importErrors = append(importErrors, fmt.Sprintf("Failed to create group '%s': %v", importConfig.Name, err))
			continue
		}

		groupIDByName[importConfig.Name] = group.GetID()
		imported++
		importedSpammers = append(importedSpammers, ImportedSpammerInfo{
			ID:          group.GetID(),
			Name:        importConfig.Name,
			Scenario:    scenario.GroupScenarioName,
			Description: importConfig.Description,
		})
	}

	// Pass 2: create standalone spammers and group members, then link members.
	for _, importConfig := range importConfigs {
		if importConfig.Scenario == scenario.GroupScenarioName {
			continue
		}

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
		configYAML, err := configs.MergeScenarioConfiguration(scenarioDescriptor, &importConfig.Config)
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

		// Link to its group if the member references one that exists.
		if importConfig.Group != "" {
			if groupID, ok := groupIDByName[importConfig.Group]; ok {
				memberCfg := configs.MemberConfigFromMap(importConfig.GroupConfig)
				if err := d.AddSpammerToGroup(spammer.GetID(), groupID, memberCfg, userEmail); err != nil {
					importWarnings = append(importWarnings, fmt.Sprintf("Imported '%s' but could not add to group '%s': %v", finalName, importConfig.Group, err))
				}
			} else {
				importWarnings = append(importWarnings, fmt.Sprintf("Imported '%s' but group '%s' was not found", finalName, importConfig.Group))
			}
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

		// Check scenario validity (group rows use the reserved sentinel and are valid)
		if importConfig.Scenario != scenario.GroupScenarioName {
			scenarioDescriptor := scenarios.GetScenario(importConfig.Scenario)
			if scenarioDescriptor == nil {
				info.Valid = false
				info.Issues = append(info.Issues, "Unknown scenario")
				result.InvalidScenarios = append(result.InvalidScenarios, importConfig.Scenario)
			}
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
