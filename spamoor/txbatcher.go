package spamoor

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/txbuilder"
	geas "github.com/fjl/geas/asm"
	"github.com/holiman/uint256"
)

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
	BatcherTxLimit  = 50
	BatcherBaseGas  = 50000
	BatcherGasPerTx = 35000
)

type TxBatcher struct {
	txpool     *TxPool
	isDeployed bool
	deployMtx  sync.Mutex
	address    common.Address
}

func NewTxBatcher(txpool *TxPool) *TxBatcher {
	return &TxBatcher{
		txpool: txpool,
	}
}

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

	err = b.txpool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     10,
		RebroadcastInterval: 30 * time.Second,
	})
	if err != nil {
		return err
	}

	b.address = crypto.CreateAddress(wallet.GetAddress(), tx.Nonce())

	return nil
}

func (b *TxBatcher) GetAddress() common.Address {
	return b.address
}

func (b *TxBatcher) GetRequestCalldata(requests []*FundingRequest) ([]byte, error) {
	calldata := make([]byte, len(requests)*32)
	for i, request := range requests {
		copy(calldata[i*32:], request.Wallet.GetAddress().Bytes())
		amountBytes := request.Amount.Bytes32()
		copy(calldata[i*32+20:], amountBytes[20:])
	}
	return calldata, nil
}
