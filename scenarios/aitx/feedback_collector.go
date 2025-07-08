package aitx

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type FeedbackCollector struct {
	results           []TransactionResult
	mutex             sync.RWMutex
	maxResults        uint64
	totalTransactions uint64
	successfulTxs     uint64
	failedTxs         uint64
	logger            logrus.FieldLogger
}

func NewFeedbackCollector(maxResults uint64, logger logrus.FieldLogger) *FeedbackCollector {
	return &FeedbackCollector{
		results:    make([]TransactionResult, 0, maxResults),
		maxResults: maxResults,
		logger:     logger.WithField("component", "feedback_collector"),
	}
}

func (fc *FeedbackCollector) RecordResult(result TransactionResult) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	fc.results = append(fc.results, result)
	if uint64(len(fc.results)) > fc.maxResults {
		fc.results = fc.results[1:]
	}

	fc.totalTransactions++
	if result.Status == "success" {
		fc.successfulTxs++
	} else {
		fc.failedTxs++
	}

	fc.logger.Debugf("recorded transaction result: %s (%s) - %s",
		result.PayloadDescription, result.PayloadType, result.Status)
}

func (fc *FeedbackCollector) RecordFailure(payload *PayloadInstance, status, errorMsg string) {
	result := TransactionResult{
		PayloadType:        payload.Type,
		PayloadDescription: payload.Description,
		Status:             status,
		GasUsed:            0,
		BlockExecTime:      "N/A",
		ErrorMessage:       errorMsg,
	}
	fc.RecordResult(result)
}

func (fc *FeedbackCollector) GenerateFeedback() *TransactionFeedback {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	if len(fc.results) == 0 {
		return nil
	}

	gasValues := make([]uint64, 0, len(fc.results))
	for _, result := range fc.results {
		if result.Status == "success" && result.GasUsed > 0 {
			gasValues = append(gasValues, result.GasUsed)
		}
	}

	var avgGas, medianGas uint64
	if len(gasValues) > 0 {
		sort.Slice(gasValues, func(i, j int) bool { return gasValues[i] < gasValues[j] })

		var total uint64
		for _, gas := range gasValues {
			total += gas
		}
		avgGas = total / uint64(len(gasValues))
		medianGas = gasValues[len(gasValues)/2]
	}

	summary := fc.generateSummary()

	return &TransactionFeedback{
		TotalTransactions:    fc.totalTransactions,
		SuccessfulTxs:        fc.successfulTxs,
		FailedTxs:            fc.failedTxs,
		AverageGasUsed:       avgGas,
		MedianGasUsed:        medianGas,
		AverageBlockExecTime: "N/A",
		RecentResults:        fc.getRecentResults(10),
		Summary:              summary,
	}
}

func (fc *FeedbackCollector) generateSummary() string {
	if len(fc.results) == 0 {
		return "No transaction results yet."
	}

	typeSuccess := make(map[string]int)
	typeTotal := make(map[string]int)

	for _, result := range fc.results {
		typeTotal[result.PayloadType]++
		if result.Status == "success" {
			typeSuccess[result.PayloadType]++
		}
	}

	var summaryParts []string
	for payloadType, total := range typeTotal {
		success := typeSuccess[payloadType]
		successRate := float64(success) / float64(total) * 100
		summaryParts = append(summaryParts,
			fmt.Sprintf("%s: %.1f%% success (%d/%d)", payloadType, successRate, success, total))
	}

	return fmt.Sprintf("Pattern analysis: %s", strings.Join(summaryParts, ", "))
}

func (fc *FeedbackCollector) getRecentResults(count int) []TransactionResult {
	if len(fc.results) <= count {
		return fc.results
	}
	return fc.results[len(fc.results)-count:]
}

func (fc *FeedbackCollector) GetStats() (uint64, uint64, uint64) {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()
	return fc.totalTransactions, fc.successfulTxs, fc.failedTxs
}
