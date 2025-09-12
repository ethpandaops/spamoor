package daemon

import (
	"fmt"
	"strings"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/jmoiron/sqlx"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/sirupsen/logrus"
)

// AuditLogger handles creating audit log entries with proper diff generation
type AuditLogger struct {
	daemon      *Daemon
	logger      *logrus.Entry
	userHeader  string
	defaultUser string
}

// NewAuditLogger creates a new audit logger instance
func NewAuditLogger(daemon *Daemon, userHeader string, defaultUser string) *AuditLogger {
	return &AuditLogger{
		daemon:      daemon,
		logger:      daemon.logger.WithField("module", "audit"),
		userHeader:  userHeader,
		defaultUser: defaultUser,
	}
}

// GetUserFromRequest extracts the user email from the request header
func (al *AuditLogger) GetUserFromRequest(headers map[string][]string) string {
	if headers == nil {
		return al.defaultUser
	}

	// Try exact header match first
	if values, ok := headers[al.userHeader]; ok && len(values) > 0 {
		return values[0]
	}

	// Try case-insensitive match
	for key, values := range headers {
		if strings.EqualFold(key, al.userHeader) && len(values) > 0 {
			return values[0]
		}
	}

	return al.defaultUser
}

// LogSpammerCreate logs a spammer creation action
func (al *AuditLogger) LogSpammerCreate(tx *sqlx.Tx, userEmail string, spammerID int64, name string, scenario string, config string) error {
	// For creation, the diff is just the new config
	diff := al.generateConfigDiff("", config, "spammer configuration")

	metadata := db.AuditMetadata{
		"scenario": scenario,
	}

	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(db.AuditActionSpammerCreate),
		EntityType: string(db.AuditEntitySpammer),
		EntityID:   fmt.Sprintf("%d", spammerID),
		EntityName: name,
		Diff:       diff,
		Metadata:   db.MarshalAuditMetadata(metadata),
	}

	return al.daemon.db.InsertAuditLog(tx, log)
}

// LogSpammerUpdate logs a spammer update action
func (al *AuditLogger) LogSpammerUpdate(tx *sqlx.Tx, userEmail string, spammerID int64, oldName, newName, oldDesc, newDesc, oldConfig, newConfig string) error {
	var diffs []string
	var changedFields []string

	// Track what changed for metadata
	if oldName != newName {
		diffs = append(diffs, al.generateFieldDiff("name", oldName, newName))
		changedFields = append(changedFields, "name")
	}

	if oldDesc != newDesc {
		diffs = append(diffs, al.generateFieldDiff("description", oldDesc, newDesc))
		changedFields = append(changedFields, "description")
	}

	if oldConfig != newConfig {
		diffs = append(diffs, al.generateConfigDiff(oldConfig, newConfig, "config"))
		changedFields = append(changedFields, "config")
	}

	diff := strings.Join(diffs, "\n")

	// If only name or description changed (single simple property), add metadata
	metadata := db.AuditMetadata{}
	if len(changedFields) == 1 && (changedFields[0] == "name" || changedFields[0] == "description") {
		metadata["single_property_change"] = true
		metadata["property_name"] = changedFields[0]
		if changedFields[0] == "name" {
			metadata["old_value"] = oldName
			metadata["new_value"] = newName
		} else {
			metadata["old_value"] = oldDesc
			metadata["new_value"] = newDesc
		}
	}

	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(db.AuditActionSpammerUpdate),
		EntityType: string(db.AuditEntitySpammer),
		EntityID:   fmt.Sprintf("%d", spammerID),
		EntityName: newName,
		Diff:       diff,
		Metadata:   db.MarshalAuditMetadata(metadata),
	}

	return al.daemon.db.InsertAuditLog(tx, log)
}

// LogSpammerDelete logs a spammer deletion action
func (al *AuditLogger) LogSpammerDelete(tx *sqlx.Tx, userEmail string, spammerID int64, name string) error {
	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(db.AuditActionSpammerDelete),
		EntityType: string(db.AuditEntitySpammer),
		EntityID:   fmt.Sprintf("%d", spammerID),
		EntityName: name,
		Diff:       "",
		Metadata:   "{}",
	}

	return al.daemon.db.InsertAuditLog(tx, log)
}

// LogSpammerAction logs simple spammer actions (start, pause, reclaim)
func (al *AuditLogger) LogSpammerAction(userEmail string, action db.AuditActionType, spammerID int64, name string, metadata db.AuditMetadata) error {
	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(action),
		EntityType: string(db.AuditEntitySpammer),
		EntityID:   fmt.Sprintf("%d", spammerID),
		EntityName: name,
		Diff:       "",
		Metadata:   db.MarshalAuditMetadata(metadata),
	}

	return al.daemon.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return al.daemon.db.InsertAuditLog(tx, log)
	})
}

// LogClientUpdate logs a client configuration update with specific action types
func (al *AuditLogger) LogClientUpdate(tx *sqlx.Tx, userEmail string, rpcURL string, clientName string, changes map[string]interface{}) error {
	// If single property change, use specific action type and metadata
	if len(changes) == 1 {
		for field, change := range changes {
			if changeMap, ok := change.(map[string]interface{}); ok {
				oldVal := fmt.Sprintf("%v", changeMap["old"])
				newVal := fmt.Sprintf("%v", changeMap["new"])

				// Determine specific action type
				var actionType db.AuditActionType
				switch field {
				case "name":
					actionType = db.AuditActionClientNameUpdate
				case "client_type":
					actionType = db.AuditActionClientTypeUpdate
				case "tags":
					actionType = db.AuditActionClientGroupUpdate
				case "enabled":
					actionType = db.AuditActionClientToggle
				default:
					actionType = db.AuditActionClientUpdate
				}

				// Create metadata for simple property change display
				metadata := db.AuditMetadata{
					"single_property_change": true,
					"property_name":          field,
					"old_value":              oldVal,
					"new_value":              newVal,
				}

				// For enabled field, also store the boolean value for display logic
				if field == "enabled" {
					metadata["enabled"] = changeMap["new"]
				}

				// For name changes, use the old name as entity name so it's clear which client was changed
				entityName := clientName
				if field == "name" {
					entityName = oldVal
				}

				log := &db.AuditLog{
					UserEmail:  userEmail,
					ActionType: string(actionType),
					EntityType: string(db.AuditEntityClient),
					EntityID:   rpcURL,
					EntityName: entityName,
					Diff:       al.generateFieldDiff(field, oldVal, newVal),
					Metadata:   db.MarshalAuditMetadata(metadata),
				}

				return al.daemon.db.InsertAuditLog(tx, log)
			}
		}
	}

	// Multiple property changes - use general update action
	var diffs []string

	// Generate diffs for each changed field
	for field, change := range changes {
		if changeMap, ok := change.(map[string]interface{}); ok {
			oldVal := fmt.Sprintf("%v", changeMap["old"])
			newVal := fmt.Sprintf("%v", changeMap["new"])
			diffs = append(diffs, al.generateFieldDiff(field, oldVal, newVal))
		}
	}

	diff := strings.Join(diffs, "\n")

	metadata := db.AuditMetadata{
		"changes": changes,
	}

	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(db.AuditActionClientUpdate),
		EntityType: string(db.AuditEntityClient),
		EntityID:   rpcURL,
		EntityName: clientName,
		Diff:       diff,
		Metadata:   db.MarshalAuditMetadata(metadata),
	}

	return al.daemon.db.InsertAuditLog(tx, log)
}

// LogSpammersImport logs a spammers import action
func (al *AuditLogger) LogSpammersImport(userEmail string, importedCount int, skippedCount int, source string) error {
	metadata := db.AuditMetadata{
		"imported_count": importedCount,
		"skipped_count":  skippedCount,
		"source":         source,
	}

	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(db.AuditActionSpammersImport),
		EntityType: string(db.AuditEntitySystem),
		EntityID:   "import",
		EntityName: fmt.Sprintf("Import %d spammers", importedCount),
		Diff:       "",
		Metadata:   db.MarshalAuditMetadata(metadata),
	}

	return al.daemon.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return al.daemon.db.InsertAuditLog(tx, log)
	})
}

// LogSpammersExport logs a spammers export action
func (al *AuditLogger) LogSpammersExport(userEmail string, exportedIDs []int64) error {
	metadata := db.AuditMetadata{
		"exported_ids":   exportedIDs,
		"exported_count": len(exportedIDs),
	}

	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(db.AuditActionSpammersExport),
		EntityType: string(db.AuditEntitySystem),
		EntityID:   "export",
		EntityName: fmt.Sprintf("Export %d spammers", len(exportedIDs)),
		Diff:       "",
		Metadata:   db.MarshalAuditMetadata(metadata),
	}

	return al.daemon.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return al.daemon.db.InsertAuditLog(tx, log)
	})
}

// LogRootWalletTransaction logs a transaction sent from the root wallet
func (al *AuditLogger) LogRootWalletTransaction(userEmail string, toAddress string, valueWei string, txHash string, gasLimit uint64, maxFee string, maxTip string, data string) error {
	metadata := db.AuditMetadata{
		"to_address": toAddress,
		"value_wei":  valueWei,
		"tx_hash":    txHash,
		"gas_limit":  gasLimit,
		"max_fee":    maxFee,
		"max_tip":    maxTip,
		"data":       data,
	}

	log := &db.AuditLog{
		UserEmail:  userEmail,
		ActionType: string(db.AuditActionRootWalletSend),
		EntityType: string(db.AuditEntityRootWallet),
		EntityID:   "root",
		EntityName: "Root Wallet Transaction",
		Diff:       fmt.Sprintf("Sent %s wei to %s\nTransaction Hash: %s\nGas Limit: %d\nMax Fee: %s\nMax Tip: %s", valueWei, toAddress, txHash, gasLimit, maxFee, maxTip),
		Metadata:   db.MarshalAuditMetadata(metadata),
	}

	return al.daemon.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return al.daemon.db.InsertAuditLog(tx, log)
	})
}

// generateConfigDiff generates a unified diff for configuration changes
func (al *AuditLogger) generateConfigDiff(oldConfig, newConfig, filename string) string {
	if oldConfig == "" && newConfig == "" {
		return ""
	}

	// Split configs into lines, preserving newlines for proper diff generation
	oldLines := strings.SplitAfter(oldConfig, "\n")
	newLines := strings.SplitAfter(newConfig, "\n")

	// Generate unified diff
	diff := difflib.UnifiedDiff{
		A:        oldLines,
		B:        newLines,
		FromFile: filename + ".old",
		ToFile:   filename + ".new",
		Context:  0,
	}

	diffStr, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		al.logger.Errorf("Failed to generate diff: %v", err)
		return fmt.Sprintf("Error generating diff: %v", err)
	}

	return diffStr
}

// generateFieldDiff generates a simple field diff
func (al *AuditLogger) generateFieldDiff(fieldName, oldValue, newValue string) string {
	return fmt.Sprintf("--- %s\n+++ %s\n@@ -1 +1 @@\n-%s\n+%s", fieldName, fieldName, oldValue, newValue)
}
