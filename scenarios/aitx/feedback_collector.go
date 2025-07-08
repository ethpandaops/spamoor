package aitx

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type PayloadStats struct {
	Description     string
	SuccessCount    int
	FailureCount    int
	TotalGasUsed    uint64
	SuccessfulCalls []TransactionResult // Successful transactions for this payload
	FailedCalls     []TransactionResult // Failed transactions for this payload
}

type FeedbackCollector struct {
	payloadStats   map[string]*PayloadStats // Per-payload statistics
	mutex          sync.RWMutex
	maxResults     uint64
	currentBatchID int // ID of current AI response batch
	logger         logrus.FieldLogger
}

func NewFeedbackCollector(maxResults uint64, logger logrus.FieldLogger) *FeedbackCollector {
	return &FeedbackCollector{
		payloadStats:   make(map[string]*PayloadStats),
		maxResults:     maxResults,
		currentBatchID: 0,
		logger:         logger.WithField("component", "feedback_collector"),
	}
}

func (fc *FeedbackCollector) RecordResult(result TransactionResult, batchID int) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	// Only record results for the current batch
	if batchID != fc.currentBatchID {
		fc.logger.Debugf("ignoring result from old batch %d (current: %d): %s",
			batchID, fc.currentBatchID, result.PayloadDescription)
		return
	}

	// Get or create payload stats
	stats, exists := fc.payloadStats[result.PayloadDescription]
	if !exists {
		stats = &PayloadStats{
			Description:     result.PayloadDescription,
			SuccessfulCalls: make([]TransactionResult, 0),
			FailedCalls:     make([]TransactionResult, 0),
		}
		fc.payloadStats[result.PayloadDescription] = stats
	}

	// Record result
	if result.Status == "success" {
		stats.SuccessCount++
		stats.TotalGasUsed += result.GasUsed
		stats.SuccessfulCalls = append(stats.SuccessfulCalls, result)

		// Keep only recent successful calls (limit to maxResults/2 per payload)
		maxPerPayload := int(fc.maxResults / 2)
		if len(stats.SuccessfulCalls) > maxPerPayload {
			stats.SuccessfulCalls = stats.SuccessfulCalls[len(stats.SuccessfulCalls)-maxPerPayload:]
		}
	} else {
		stats.FailureCount++
		stats.FailedCalls = append(stats.FailedCalls, result)

		// Keep only recent failed calls (limit to maxResults/2 per payload)
		maxPerPayload := int(fc.maxResults / 2)
		if len(stats.FailedCalls) > maxPerPayload {
			stats.FailedCalls = stats.FailedCalls[len(stats.FailedCalls)-maxPerPayload:]
		}
	}

	fc.logger.Debugf("recorded transaction result for batch %d: %s (%s) - %s (success: %d, failed: %d)",
		batchID, result.PayloadDescription, result.PayloadType, result.Status,
		stats.SuccessCount, stats.FailureCount)
}

func (fc *FeedbackCollector) RecordFailure(payload *PayloadInstance, status, errorMsg string, batchID int) {
	result := TransactionResult{
		PayloadType:        payload.Type,
		PayloadDescription: payload.Description,
		Status:             status,
		GasUsed:            0,
		BlockExecTime:      "N/A",
		ErrorMessage:       errorMsg,
	}
	fc.RecordResult(result, batchID)
}

func (fc *FeedbackCollector) GenerateFeedback() *TransactionFeedback {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	if len(fc.payloadStats) == 0 {
		return nil
	}

	var totalTxs, successfulTxs, failedTxs uint64
	var allGasValues []uint64
	var allResults []TransactionResult

	// Aggregate stats from all payloads
	for _, stats := range fc.payloadStats {
		totalTxs += uint64(stats.SuccessCount + stats.FailureCount)
		successfulTxs += uint64(stats.SuccessCount)
		failedTxs += uint64(stats.FailureCount)

		// Collect gas values from successful calls
		for _, result := range stats.SuccessfulCalls {
			if result.GasUsed > 0 {
				allGasValues = append(allGasValues, result.GasUsed)
			}
			allResults = append(allResults, result)
		}

		// Include failed calls in results
		allResults = append(allResults, stats.FailedCalls...)
	}

	var avgGas, medianGas uint64
	if len(allGasValues) > 0 {
		sort.Slice(allGasValues, func(i, j int) bool { return allGasValues[i] < allGasValues[j] })

		var total uint64
		for _, gas := range allGasValues {
			total += gas
		}
		avgGas = total / uint64(len(allGasValues))
		medianGas = allGasValues[len(allGasValues)/2]
	}

	summary := fc.generateDetailedSummary()

	return &TransactionFeedback{
		TotalTransactions:    totalTxs,
		SuccessfulTxs:        successfulTxs,
		FailedTxs:            failedTxs,
		AverageGasUsed:       avgGas,
		MedianGasUsed:        medianGas,
		AverageBlockExecTime: "N/A",
		RecentResults:        allResults,
		Summary:              summary,
	}
}

func (fc *FeedbackCollector) generateDetailedSummary() string {
	if len(fc.payloadStats) == 0 {
		return "No transaction results for current batch yet."
	}

	var summaryParts []string

	for description, stats := range fc.payloadStats {
		total := stats.SuccessCount + stats.FailureCount
		if total == 0 {
			continue
		}

		successRate := float64(stats.SuccessCount) / float64(total) * 100

		// Calculate average gas for this payload
		var avgGas uint64
		if stats.SuccessCount > 0 {
			avgGas = stats.TotalGasUsed / uint64(stats.SuccessCount)
		}

		// Get recent logs from successful calls
		var recentLogs []string
		for _, result := range stats.SuccessfulCalls {
			if len(result.LogData) > 0 {
				recentLogs = append(recentLogs, result.LogData[0]) // Take first log
				if len(recentLogs) >= 2 {                          // Limit to 2 logs per payload
					break
				}
			}
		}

		// Get recent error messages from failed calls
		var recentErrors []string
		for _, result := range stats.FailedCalls {
			if result.ErrorMessage != "" {
				recentErrors = append(recentErrors, result.ErrorMessage)
				if len(recentErrors) >= 2 { // Limit to 2 errors per payload
					break
				}
			}
		}

		payloadSummary := fmt.Sprintf("'%s': %.1f%% success (%d/%d), avg_gas: %d",
			description, successRate, stats.SuccessCount, total, avgGas)

		if len(recentLogs) > 0 {
			payloadSummary += fmt.Sprintf(", recent_logs: [%s]", strings.Join(recentLogs, ", "))
		}

		if len(recentErrors) > 0 {
			payloadSummary += fmt.Sprintf(", recent_errors: [%s]", strings.Join(recentErrors, ", "))
		}

		summaryParts = append(summaryParts, payloadSummary)
	}

	return fmt.Sprintf("Detailed payload analysis: %s", strings.Join(summaryParts, " | "))
}

// StartNewBatch resets the feedback collector for a new AI response batch
func (fc *FeedbackCollector) StartNewBatch() {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	fc.currentBatchID++
	fc.payloadStats = make(map[string]*PayloadStats)
	fc.logger.Infof("started new feedback batch %d", fc.currentBatchID)
}

func (fc *FeedbackCollector) GetCurrentBatchID() int {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()
	return fc.currentBatchID
}

func (fc *FeedbackCollector) GetCurrentBatchStats() (uint64, uint64, uint64) {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	var totalTxs, successfulTxs, failedTxs uint64
	for _, stats := range fc.payloadStats {
		totalTxs += uint64(stats.SuccessCount + stats.FailureCount)
		successfulTxs += uint64(stats.SuccessCount)
		failedTxs += uint64(stats.FailureCount)
	}

	return totalTxs, successfulTxs, failedTxs
}

// GetFailedPayloads returns payloads that have failures for immediate feedback
func (fc *FeedbackCollector) GetFailedPayloads() map[string]*PayloadStats {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	failedPayloads := make(map[string]*PayloadStats)
	for description, stats := range fc.payloadStats {
		if stats.FailureCount > 0 {
			failedPayloads[description] = stats
		}
	}

	return failedPayloads
}
