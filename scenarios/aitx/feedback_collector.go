package aitx

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type PayloadFeedback struct {
	PayloadID           string `json:"payload_id"`
	PayloadIndex        int    `json:"payload_index"`
	Description         string `json:"description"`
	CompilationStatus   string `json:"compilation_status"` // "success", "failed"
	CompilationError    string `json:"compilation_error,omitempty"`
	ExecutionStatus     string `json:"execution_status"` // "success", "failed", "not_executed"
	ExecutionError      string `json:"execution_error,omitempty"`
	GasUsed             uint64 `json:"gas_used"`
	ExecutionCount      int    `json:"execution_count"`
	LastExecutionResult string `json:"last_execution_result,omitempty"`
}

type FeedbackCollector struct {
	payloadFeedbacks []PayloadFeedback // Ordered list of payload feedback
	payloadLookup    map[string]int    // Map payload ID to index for quick lookup
	mutex            sync.RWMutex
	maxResults       uint64
	currentBatchID   int // ID of current AI response batch
	logger           logrus.FieldLogger
}

func NewFeedbackCollector(maxResults uint64, logger logrus.FieldLogger) *FeedbackCollector {
	return &FeedbackCollector{
		payloadFeedbacks: make([]PayloadFeedback, 0),
		payloadLookup:    make(map[string]int),
		maxResults:       maxResults,
		currentBatchID:   0,
		logger:           logger.WithField("component", "feedback_collector"),
	}
}

// RegisterPayload adds a new payload in order for tracking
func (fc *FeedbackCollector) RegisterPayload(payloadID, description string, batchID int) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	// Only register payloads for the current batch
	if batchID != fc.currentBatchID {
		fc.logger.Debugf("ignoring payload registration from old batch %d (current: %d): %s",
			batchID, fc.currentBatchID, payloadID)
		return
	}

	// Check if already registered
	if _, exists := fc.payloadLookup[payloadID]; exists {
		return
	}

	// Add new payload feedback entry
	index := len(fc.payloadFeedbacks)
	fc.payloadFeedbacks = append(fc.payloadFeedbacks, PayloadFeedback{
		PayloadID:         payloadID,
		PayloadIndex:      index,
		Description:       description,
		CompilationStatus: "pending",
		ExecutionStatus:   "not_executed",
		GasUsed:           0,
		ExecutionCount:    0,
	})
	fc.payloadLookup[payloadID] = index

	fc.logger.Debugf("registered payload %d for batch %d: %s (%s)", index, batchID, payloadID, description)
}

// RecordCompilationResult records the compilation status for a payload
func (fc *FeedbackCollector) RecordCompilationResult(payloadID string, success bool, errorMsg string, batchID int) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	if batchID != fc.currentBatchID {
		return
	}

	index, exists := fc.payloadLookup[payloadID]
	if !exists {
		fc.logger.Warnf("compilation result for unknown payload: %s", payloadID)
		return
	}

	if success {
		fc.payloadFeedbacks[index].CompilationStatus = "success"
		fc.payloadFeedbacks[index].CompilationError = ""
	} else {
		fc.payloadFeedbacks[index].CompilationStatus = "failed"
		fc.payloadFeedbacks[index].CompilationError = errorMsg
	}

	fc.logger.Debugf("recorded compilation result for payload %d: %s - %s",
		index, payloadID, fc.payloadFeedbacks[index].CompilationStatus)
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

	index, exists := fc.payloadLookup[result.PayloadDescription]
	if !exists {
		fc.logger.Warnf("execution result for unknown payload: %s", result.PayloadDescription)
		return
	}

	// Update execution feedback
	fc.payloadFeedbacks[index].ExecutionCount++
	fc.payloadFeedbacks[index].GasUsed = result.GasUsed

	if result.Status == "success" {
		fc.payloadFeedbacks[index].ExecutionStatus = "success"
		fc.payloadFeedbacks[index].ExecutionError = ""
		fc.payloadFeedbacks[index].LastExecutionResult = fmt.Sprintf("gas:%d", result.GasUsed)
	} else {
		fc.payloadFeedbacks[index].ExecutionStatus = "failed"
		fc.payloadFeedbacks[index].ExecutionError = result.ErrorMessage
		fc.payloadFeedbacks[index].LastExecutionResult = fmt.Sprintf("error:%s", result.ErrorMessage)
	}

	fc.logger.Debugf("recorded execution result for payload %d: %s - %s (count: %d)",
		index, result.PayloadDescription, fc.payloadFeedbacks[index].ExecutionStatus,
		fc.payloadFeedbacks[index].ExecutionCount)
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

	if len(fc.payloadFeedbacks) == 0 {
		return nil
	}

	var totalTxs, successfulTxs, failedTxs uint64
	var allGasValues []uint64
	var allResults []TransactionResult

	// Process payloads in order
	for _, feedback := range fc.payloadFeedbacks {
		if feedback.ExecutionCount > 0 {
			totalTxs += uint64(feedback.ExecutionCount)

			if feedback.ExecutionStatus == "success" {
				successfulTxs += uint64(feedback.ExecutionCount)
				if feedback.GasUsed > 0 {
					allGasValues = append(allGasValues, feedback.GasUsed)
				}
			} else {
				failedTxs += uint64(feedback.ExecutionCount)
			}

			// Create a result entry for this payload
			result := TransactionResult{
				PayloadType:        "geas",
				PayloadDescription: feedback.Description,
				Status:             feedback.ExecutionStatus,
				GasUsed:            feedback.GasUsed,
				BlockExecTime:      "N/A",
				ErrorMessage:       feedback.ExecutionError,
			}
			allResults = append(allResults, result)
		}
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
	if len(fc.payloadFeedbacks) == 0 {
		return "No payloads received for current batch yet."
	}

	var summaryParts []string

	// Generate feedback for each payload in order
	for _, feedback := range fc.payloadFeedbacks {
		var status string
		if feedback.CompilationStatus == "failed" {
			status = fmt.Sprintf("compilation_failed: %s", feedback.CompilationError)
		} else if feedback.ExecutionStatus == "not_executed" {
			status = "not_executed"
		} else if feedback.ExecutionStatus == "failed" {
			status = fmt.Sprintf("execution_failed: %s", feedback.ExecutionError)
		} else {
			status = fmt.Sprintf("success: gas=%d", feedback.GasUsed)
		}

		payloadSummary := fmt.Sprintf("%s ('%s'): %s",
			feedback.PayloadID, feedback.Description, status)
		summaryParts = append(summaryParts, payloadSummary)
	}

	return fmt.Sprintf("ORDERED PAYLOAD FEEDBACK: %s", strings.Join(summaryParts, " | "))
}

// StartNewBatch resets the feedback collector for a new AI response batch
func (fc *FeedbackCollector) StartNewBatch() {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	fc.currentBatchID++
	fc.payloadFeedbacks = make([]PayloadFeedback, 0)
	fc.payloadLookup = make(map[string]int)
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
	for _, feedback := range fc.payloadFeedbacks {
		if feedback.ExecutionCount > 0 {
			totalTxs += uint64(feedback.ExecutionCount)
			if feedback.ExecutionStatus == "success" {
				successfulTxs += uint64(feedback.ExecutionCount)
			} else {
				failedTxs += uint64(feedback.ExecutionCount)
			}
		}
	}

	return totalTxs, successfulTxs, failedTxs
}

// GetFailedPayloads returns payloads that have failures for immediate feedback
func (fc *FeedbackCollector) GetFailedPayloads() []PayloadFeedback {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	var failedPayloads []PayloadFeedback
	for _, feedback := range fc.payloadFeedbacks {
		if feedback.CompilationStatus == "failed" || feedback.ExecutionStatus == "failed" {
			failedPayloads = append(failedPayloads, feedback)
		}
	}

	return failedPayloads
}
