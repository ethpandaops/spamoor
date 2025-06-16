package testingutils

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// TransactionValidator interface for capturing transactions
type TransactionValidator interface {
	CaptureTransaction(tx *types.Transaction, senderAddress common.Address)
}

// ValidatingMockClient extends MockClient with transaction capture capabilities
type ValidatingMockClient struct {
	*MockClient
	validator TransactionValidator
	logger    *logrus.Entry
}

// NewValidatingMockClient creates a new validating mock client
func NewValidatingMockClient(validator TransactionValidator, logger *logrus.Entry) *ValidatingMockClient {
	return &ValidatingMockClient{
		MockClient: NewMockClient(),
		validator:  validator,
		logger:     logger.WithField("component", "validating-mock-client"),
	}
}

// NewValidatingMockClientFromExisting creates a validating mock client from an existing mock client
func NewValidatingMockClientFromExisting(mockClient *MockClient, validator TransactionValidator, logger *logrus.Entry) *ValidatingMockClient {
	return &ValidatingMockClient{
		MockClient: mockClient,
		validator:  validator,
		logger:     logger.WithField("component", "validating-mock-client"),
	}
}

// SendTransaction overrides the MockClient SendTransaction to capture transactions
func (c *ValidatingMockClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	// Capture transaction for validation if validator is available
	if c.validator != nil {
		// We need to determine the sender address from the transaction
		senderAddress, err := c.getSenderAddress(tx)
		if err != nil {
			c.logger.WithError(err).WithField("tx_hash", tx.Hash().Hex()).Warn("Failed to determine sender address for transaction validation")
		} else {
			c.validator.CaptureTransaction(tx, senderAddress)
			c.logger.WithFields(logrus.Fields{
				"tx_hash": tx.Hash().Hex(),
				"sender":  senderAddress.Hex(),
				"nonce":   tx.Nonce(),
			}).Debug("Transaction captured for validation")
		}
	}
	
	// Continue with normal mock behavior
	return c.MockClient.SendTransaction(ctx, tx)
}

// getSenderAddress extracts the sender address from a signed transaction
func (c *ValidatingMockClient) getSenderAddress(tx *types.Transaction) (common.Address, error) {
	// For signed transactions, we can recover the sender
	chainId, err := c.MockClient.GetChainId(context.Background())
	if err != nil {
		return common.Address{}, err
	}
	
	signer := types.LatestSignerForChainID(chainId)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

// SetValidator sets or updates the transaction validator
func (c *ValidatingMockClient) SetValidator(validator TransactionValidator) {
	c.validator = validator
}

// GetValidator returns the current transaction validator
func (c *ValidatingMockClient) GetValidator() TransactionValidator {
	return c.validator
}