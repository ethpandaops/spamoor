package daemon

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// ImportSpammersOnStartup imports spammers from a file or URL on startup.
// This method reuses the import functionality but starts the spammers after importing
// them, which is the key difference from regular imports.
// Only imports on first startup for safety.
func (d *Daemon) ImportSpammersOnStartup(source string, logger logrus.FieldLogger) error {
	logger.Infof("importing startup spammers from %s", source)

	// Import the spammers using the unified import functionality
	result, err := d.ImportSpammers(source)
	if err != nil {
		return fmt.Errorf("failed to import startup spammers: %w", err)
	}

	if len(result.Errors) > 0 {
		logger.Warnf("encountered %d errors during import:", len(result.Errors))
		for _, importError := range result.Errors {
			logger.Warnf("  - %s", importError)
		}
	}

	if len(result.Warnings) > 0 {
		logger.Infof("%d warnings during import:", len(result.Warnings))
		for _, warning := range result.Warnings {
			logger.Infof("  - %s", warning)
		}
	}

	logger.Infof("successfully imported %d spammers", len(result.Imported))

	// Start all imported spammers (this is the key difference from regular import)
	startedCount := 0
	if len(result.Imported) > 0 {
		logger.Info("starting imported startup spammers:")
		for _, importedInfo := range result.Imported {
			spammer := d.GetSpammer(importedInfo.ID)
			if spammer != nil {
				err := spammer.Start()
				if err != nil {
					logger.Errorf("failed to start imported spammer %s: %v", importedInfo.Name, err)
				} else {
					logger.Infof("  - Started ID %d: %s (%s)", importedInfo.ID, importedInfo.Name, importedInfo.Scenario)
					startedCount++
				}
			}
		}
	}

	logger.Infof("successfully started %d out of %d imported spammers", startedCount, len(result.Imported))
	return nil
}
