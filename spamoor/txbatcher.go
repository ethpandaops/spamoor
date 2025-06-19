package spamoor

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/txbuilder"
	geas "github.com/fjl/geas/asm"
	"github.com/holiman/uint256"
)

// Assembly code for the batcher contract initialization
const batcherGeasInit = `
;; Init code
push @.start
codesize
sub
dup1
push @.start
push0
codecopy
push0
return

.start:
`

// Assembly code for the main batcher contract logic that forwards funds to multiple addresses
const batcherGeasCode = `
push 0                ;; [offset]
jump @loop

loop:
	DUP1              ;; [offset, offset]
	CALLDATASIZE      ;; [calldatasize, offset, offset]
	GT                ;; [calldatasize > offset, offset]
	ISZERO            ;; [calldatasize <= offset, offset]
	JUMPI @exit       ;; [offset]

	;; forward funds
	PUSH 0            ;; [0, offset]
	PUSH 0            ;; [0, 0, offset]
	PUSH 0            ;; [0, 0, 0, offset]
	PUSH 0            ;; [0, 0, 0, 0, offset]

	;; get calldata
	DUP5              ;; [offset, 0, 0, 0, 0, offset]
	CALLDATALOAD      ;; [calldata, 0, 0, 0, 0, offset]

	;; get amount
	PUSH 160          ;; [160, calldata, 0, 0, 0, 0, offset]
	DUP2              ;; [calldata, 160, calldata, 0, 0, 0, 0, offset]
	DUP2              ;; [160, calldata, 160, calldata, 0, 0, 0, 0, offset]
	SHL               ;; [calldata<<160, 160, calldata, 0, 0, 0, 0, offset]
	SWAP1             ;; [160, calldata<<160, calldata, 0, 0, 0, 0, offset]
	SHR               ;; [amount, calldata, 0, 0, 0, 0, offset]
	
	;; get address
	SWAP1             ;; [calldata, amount, 0, 0, 0, 0, offset]
	PUSH 96           ;; [96, calldata, amount, 0, 0, 0, 0, offset]
	SHR               ;; [address, amount, 0, 0, 0, 0, offset]
	
	;; forward funds
	PUSH 30000        ;; [30000, address, amount, 0, 0, 0, 0, offset]
	CALL              ;; [success, offset]
	POP               ;; [offset]

	;; increase offset
	PUSH 32           ;; [32, offset]
	ADD               ;; [offset+32]

	jump @loop

exit:
	SELFBALANCE       ;; [selfbalance]
	DUP1              ;; [selfbalance, selfbalance]
	ISZERO            ;; [selfbalance == 0, selfbalance]
	JUMPI @exit2      ;; [selfbalance]

	;; return leftover funds
	PUSH 0            ;; [0, selfbalance]
	PUSH 0            ;; [0, 0, selfbalance]
	PUSH 0            ;; [0, 0, 0, selfbalance]
	PUSH 0            ;; [0, 0, 0, 0, selfbalance]
	SWAP4             ;; [selfbalance, 0, 0, 0, 0]
	CALLER            ;; [caller, selfbalance, 0, 0, 0, 0]
	GAS               ;; [gas, caller, selfbalance, 0, 0, 0, 0]
	CALL              ;; [success]

exit2:
	STOP
`

const (
	// BatcherTxLimit is the maximum number of transactions that can be batched in a single call.
	BatcherTxLimit = 50
	// BatcherBaseGas is the base gas cost for executing a batcher transaction.
	BatcherBaseGas = 50000
	// BatcherGasPerTx is the additional gas cost per transaction in the batch.
	BatcherGasPerTx = 35000
)

// TxBatcher manages the deployment and operation of a smart contract that batches
// multiple fund transfers into a single transaction. It compiles and deploys
// assembly code that efficiently forwards funds to multiple recipients.
type TxBatcher struct {
	txpool     *TxPool
	isDeployed bool
	deployMtx  sync.Mutex
	address    common.Address
}

// NewTxBatcher creates a new TxBatcher instance with the specified transaction pool.
// The batcher must be deployed with Deploy() before it can be used.
func NewTxBatcher(txpool *TxPool) *TxBatcher {
	return &TxBatcher{
		txpool: txpool,
	}
}

// Deploy compiles and deploys the batcher smart contract to the blockchain.
// It uses assembly code to create an efficient contract that can forward funds
// to multiple addresses in a single transaction. The deployment is protected
// by a mutex to ensure it only happens once.
//
// Parameters:
//   - ctx: context for the deployment transaction
//   - wallet: wallet to deploy the contract from
//   - client: optional client to use (if nil, uses pool's default client)
func (b *TxBatcher) Deploy(ctx context.Context, wallet *Wallet, client *Client) error {
	b.deployMtx.Lock()
	defer b.deployMtx.Unlock()

	if b.isDeployed {
		return nil
	}

	b.isDeployed = true

	compiler := geas.NewCompiler(nil)

	initcode := compiler.CompileString(batcherGeasInit)
	if initcode == nil {
		return fmt.Errorf("failed to compile initcode")
	}

	batcherGeasCode := compiler.CompileString(batcherGeasCode)
	if batcherGeasCode == nil {
		return fmt.Errorf("failed to compile batcher geas code")
	}

	deployData := append(initcode, batcherGeasCode...)

	if client == nil {
		client = b.txpool.options.ClientPool.GetClient(SelectClientByIndex, 0, "")
		if client == nil {
			return fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(ctx)
	if err != nil {
		return err
	}
	if feeCap.Cmp(big.NewInt(400000000000)) < 0 {
		feeCap = big.NewInt(400000000000)
	}
	if tipCap.Cmp(big.NewInt(200000000000)) < 0 {
		tipCap = big.NewInt(200000000000)
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       300000,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      deployData,
	})
	if err != nil {
		return err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return err
	}

	err = b.txpool.submitter.Send(ctx, wallet, tx, &SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return err
	}

	b.address = crypto.CreateAddress(wallet.GetAddress(), tx.Nonce())

	return nil
}

// GetAddress returns the deployed contract address.
// Returns zero address if the contract hasn't been deployed yet.
func (b *TxBatcher) GetAddress() common.Address {
	return b.address
}

// GetRequestCalldata encodes funding requests into calldata format expected by the batcher contract.
// Each request is encoded as 32 bytes: 20 bytes for the address and 12 bytes for the amount.
// Returns the encoded calldata that can be used in a transaction to the batcher contract.
func (b *TxBatcher) GetRequestCalldata(requests []*FundingRequest) ([]byte, error) {
	calldata := make([]byte, len(requests)*32)
	for i, request := range requests {
		copy(calldata[i*32:], request.Wallet.GetAddress().Bytes())
		amountBytes := request.Amount.Bytes32()
		copy(calldata[i*32+20:], amountBytes[20:])
	}
	return calldata, nil
}
