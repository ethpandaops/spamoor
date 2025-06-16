package validators

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// TransactionValidator validates transaction ordering and nonce consistency
type TransactionValidator struct {
	// Transaction storage by address
	transactionsByAddress map[common.Address][]*types.Transaction
	transactionOrder      []TransactionRecord
	
	// Validation state
	nonceGaps           int
	orderingViolations  int
	validationErrors    []string
	
	// Runtime state
	started             bool
	stopped             bool
	
	// Synchronization
	mutex               sync.RWMutex
	
	// Configuration
	logger              *logrus.Entry
}

// TransactionRecord represents a captured transaction with metadata
type TransactionRecord struct {
	Transaction   *types.Transaction
	SenderAddress common.Address
	CaptureTime   time.Time
	Nonce         uint64
	Hash          common.Hash
}

// NewTransactionValidator creates a new transaction validator
func NewTransactionValidator(logger *logrus.Entry) *TransactionValidator {
	return &TransactionValidator{
		transactionsByAddress: make(map[common.Address][]*types.Transaction),
		transactionOrder:      make([]TransactionRecord, 0),
		validationErrors:      make([]string, 0),
		logger:               logger.WithField("component", "transaction-validator"),
	}
}

// Start begins transaction validation
func (v *TransactionValidator) Start() {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	
	if v.started {
		return
	}
	
	v.started = true
	v.stopped = false
	v.logger.Info("Transaction validator started")
}

// Stop ends transaction validation and performs final validation
func (v *TransactionValidator) Stop() {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	
	if v.stopped {
		return
	}
	
	v.stopped = true
	
	// Perform final validation
	v.performFinalValidation()
	
	v.logger.WithFields(logrus.Fields{
		"total_transactions":   len(v.transactionOrder),
		"nonce_gaps":          v.nonceGaps,
		"ordering_violations": v.orderingViolations,
		"validation_errors":   len(v.validationErrors),
	}).Info("Transaction validator stopped")
}

// CaptureTransaction captures a transaction for validation
func (v *TransactionValidator) CaptureTransaction(tx *types.Transaction, senderAddress common.Address) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	
	if v.stopped {
		return
	}
	
	// Create transaction record
	record := TransactionRecord{
		Transaction:   tx,
		SenderAddress: senderAddress,
		CaptureTime:   time.Now(),
		Nonce:         tx.Nonce(),
		Hash:          tx.Hash(),
	}
	
	// Add to global order
	v.transactionOrder = append(v.transactionOrder, record)
	
	// Add to address-specific list
	if _, exists := v.transactionsByAddress[senderAddress]; !exists {
		v.transactionsByAddress[senderAddress] = make([]*types.Transaction, 0)
	}
	v.transactionsByAddress[senderAddress] = append(v.transactionsByAddress[senderAddress], tx)
	
	// Perform real-time validation
	v.validateTransactionNonce(senderAddress, tx)
	
	v.logger.WithFields(logrus.Fields{
		"tx_hash":        tx.Hash().Hex(),
		"sender":         senderAddress.Hex(),
		"nonce":          tx.Nonce(),
		"total_captured": len(v.transactionOrder),
	}).Debug("Transaction captured for validation")
}

// validateTransactionNonce validates nonce consistency for a specific address
func (v *TransactionValidator) validateTransactionNonce(address common.Address, tx *types.Transaction) {
	txs := v.transactionsByAddress[address]
	if len(txs) <= 1 {
		return // First transaction for this address, nothing to validate
	}
	
	// Check if this transaction's nonce follows the expected sequence
	previousTx := txs[len(txs)-2] // Second to last transaction
	expectedNonce := previousTx.Nonce() + 1
	actualNonce := tx.Nonce()
	
	if actualNonce != expectedNonce {
		v.nonceGaps++
		errorMsg := fmt.Sprintf("Nonce gap detected for %s: expected %d, got %d", 
			address.Hex(), expectedNonce, actualNonce)
		v.validationErrors = append(v.validationErrors, errorMsg)
		
		v.logger.WithFields(logrus.Fields{
			"address":        address.Hex(),
			"expected_nonce": expectedNonce,
			"actual_nonce":   actualNonce,
			"tx_hash":        tx.Hash().Hex(),
		}).Warn("Nonce gap detected")
	}
}

// performFinalValidation performs comprehensive validation after all transactions are captured
func (v *TransactionValidator) performFinalValidation() {
	v.logger.Info("Performing final transaction validation")
	
	// Validate nonce sequences for each address
	for address, txs := range v.transactionsByAddress {
		v.validateNonceSequence(address, txs)
	}
	
	// Validate global transaction ordering
	v.validateGlobalOrdering()
}

// validateNonceSequence validates the complete nonce sequence for an address
func (v *TransactionValidator) validateNonceSequence(address common.Address, txs []*types.Transaction) {
	if len(txs) <= 1 {
		return
	}
	
	for i := 1; i < len(txs); i++ {
		expectedNonce := txs[i-1].Nonce() + 1
		actualNonce := txs[i].Nonce()
		
		if actualNonce != expectedNonce {
			// This might have been caught in real-time validation, but double-check
			found := false
			for _, errMsg := range v.validationErrors {
				if errMsg == fmt.Sprintf("Nonce gap detected for %s: expected %d, got %d", 
					address.Hex(), expectedNonce, actualNonce) {
					found = true
					break
				}
			}
			
			if !found {
				v.nonceGaps++
				errorMsg := fmt.Sprintf("Final validation - Nonce gap for %s: expected %d, got %d", 
					address.Hex(), expectedNonce, actualNonce)
				v.validationErrors = append(v.validationErrors, errorMsg)
			}
		}
	}
	
	v.logger.WithFields(logrus.Fields{
		"address":           address.Hex(),
		"transaction_count": len(txs),
		"first_nonce":       txs[0].Nonce(),
		"last_nonce":        txs[len(txs)-1].Nonce(),
	}).Debug("Validated nonce sequence for address")
}

// validateGlobalOrdering validates global transaction ordering consistency
func (v *TransactionValidator) validateGlobalOrdering() {
	if len(v.transactionOrder) <= 1 {
		return
	}
	
	// Group transactions by address to check within-address ordering
	addressGroups := make(map[common.Address][]int)
	for i, record := range v.transactionOrder {
		address := record.SenderAddress
		if _, exists := addressGroups[address]; !exists {
			addressGroups[address] = make([]int, 0)
		}
		addressGroups[address] = append(addressGroups[address], i)
	}
	
	// Check ordering within each address group
	for address, indices := range addressGroups {
		for i := 1; i < len(indices); i++ {
			prevIdx := indices[i-1]
			currIdx := indices[i]
			
			prevRecord := v.transactionOrder[prevIdx]
			currRecord := v.transactionOrder[currIdx]
			
			// Check nonce ordering
			if currRecord.Nonce <= prevRecord.Nonce {
				v.orderingViolations++
				errorMsg := fmt.Sprintf("Ordering violation for %s: nonce %d came before nonce %d in global order", 
					address.Hex(), currRecord.Nonce, prevRecord.Nonce)
				v.validationErrors = append(v.validationErrors, errorMsg)
				
				v.logger.WithFields(logrus.Fields{
					"address":      address.Hex(),
					"prev_nonce":   prevRecord.Nonce,
					"curr_nonce":   currRecord.Nonce,
					"prev_time":    prevRecord.CaptureTime,
					"curr_time":    currRecord.CaptureTime,
				}).Warn("Transaction ordering violation detected")
			}
			
			// Check time ordering (transactions should be captured in time order)
			if currRecord.CaptureTime.Before(prevRecord.CaptureTime) {
				v.orderingViolations++
				errorMsg := fmt.Sprintf("Time ordering violation for %s: later transaction captured before earlier one", 
					address.Hex())
				v.validationErrors = append(v.validationErrors, errorMsg)
			}
		}
	}
	
	v.logger.WithFields(logrus.Fields{
		"total_addresses":     len(addressGroups),
		"total_transactions":  len(v.transactionOrder),
		"ordering_violations": v.orderingViolations,
	}).Info("Global ordering validation completed")
}

// GetValidationResults returns the current validation results
func (v *TransactionValidator) GetValidationResults() ValidationResults {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	
	return ValidationResults{
		TotalTransactions:   len(v.transactionOrder),
		NonceGaps:          v.nonceGaps,
		OrderingViolations: v.orderingViolations,
		ValidationErrors:   make([]string, len(v.validationErrors)),
		AddressCount:       len(v.transactionsByAddress),
	}
}

// ValidationResults contains the results of transaction validation
type ValidationResults struct {
	TotalTransactions   int      `json:"total_transactions"`
	NonceGaps          int      `json:"nonce_gaps"`
	OrderingViolations int      `json:"ordering_violations"`
	ValidationErrors   []string `json:"validation_errors"`
	AddressCount       int      `json:"address_count"`
}

// GetTransactionsByAddress returns all transactions for a specific address
func (v *TransactionValidator) GetTransactionsByAddress(address common.Address) []*types.Transaction {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	
	if txs, exists := v.transactionsByAddress[address]; exists {
		// Return a copy to prevent external modification
		result := make([]*types.Transaction, len(txs))
		copy(result, txs)
		return result
	}
	
	return nil
}

// GetAllTransactions returns all captured transactions in order
func (v *TransactionValidator) GetAllTransactions() []TransactionRecord {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	result := make([]TransactionRecord, len(v.transactionOrder))
	copy(result, v.transactionOrder)
	return result
}

// GetAddresses returns all addresses that have sent transactions
func (v *TransactionValidator) GetAddresses() []common.Address {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	
	addresses := make([]common.Address, 0, len(v.transactionsByAddress))
	for address := range v.transactionsByAddress {
		addresses = append(addresses, address)
	}
	
	return addresses
}

// Reset resets the validator state (useful for multiple test runs)
func (v *TransactionValidator) Reset() {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	
	v.transactionsByAddress = make(map[common.Address][]*types.Transaction)
	v.transactionOrder = make([]TransactionRecord, 0)
	v.nonceGaps = 0
	v.orderingViolations = 0
	v.validationErrors = make([]string, 0)
	v.started = false
	v.stopped = false
	
	v.logger.Info("Transaction validator reset")
}