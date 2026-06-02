package spamoor

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txbuilder"
)

// CanonicalCreate2FactoryAddress is the address of the Arachnid
// "deterministic-deployment-proxy" CREATE2 factory. It is present on most
// public Ethereum networks at this well-known address and is the first place
// the DeploymentFactory looks when initializing.
//
// See https://github.com/Arachnid/deterministic-deployment-proxy.
var canonicalCreate2FactoryAddress = common.HexToAddress("0x4e59b44847b379578588920cA78FbF26c0B4956C")

// create2FactoryInitcode is the init code of the Arachnid deterministic
// deployment proxy. When deployed it yields a minimal CREATE2 factory whose
// calldata convention is: salt (32 bytes) followed by the contract init code.
// The factory CREATE2-deploys the init code with the given salt and returns the
// 20-byte address of the deployed contract.
var create2FactoryCode = common.FromHex("0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf3")
var create2FactoryInitcode = append(common.FromHex("0x604580600e600039806000f350fe"), create2FactoryCode...)

type DeploymentFactory struct {
	txpool        *TxPool
	rootWallet    *Wallet
	initMutex     sync.Mutex
	isInitialized bool
	isDeploying   bool
	factoryAddr   common.Address

	contractMutex       sync.Mutex
	contractDeployments map[common.Address]*ContractDeployment
}

type ContractDeployment struct {
	mu  *sync.Mutex
	tx  *types.Transaction
	err error
}

func newDeploymentFactory(txpool *TxPool, rootWallet *Wallet) *DeploymentFactory {
	return &DeploymentFactory{
		txpool:              txpool,
		rootWallet:          rootWallet,
		contractDeployments: make(map[common.Address]*ContractDeployment),
	}
}

func (f *DeploymentFactory) lazyInit(ctx context.Context) error {
	if f.isInitialized {
		return nil
	}

	f.initMutex.Lock()
	defer f.initMutex.Unlock()

	if f.isInitialized {
		return nil
	}

	client := f.txpool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientByIndex, 0), WithoutBuilder())
	if client == nil {
		return fmt.Errorf("no client available")
	}

	// 1. canonical proxy address
	if code, err := client.GetCodeAt(ctx, canonicalCreate2FactoryAddress); err == nil && len(code) > 0 {
		if bytes.Equal(code, create2FactoryCode) {
			f.txpool.options.Logger.Infof("using canonical CREATE2 factory at %s", canonicalCreate2FactoryAddress.Hex())
			f.factoryAddr = canonicalCreate2FactoryAddress
			f.isInitialized = true
			return nil
		} else {
			f.txpool.options.Logger.Warnf("canonical CREATE2 factory at %s is not the expected code", canonicalCreate2FactoryAddress.Hex())
		}
	}

	// 2. previously self-deployed factory at the root wallet's nonce 0
	candidate := crypto.CreateAddress(f.rootWallet.GetAddress(), 0)
	if code, err := client.GetCodeAt(ctx, candidate); err == nil {
		if bytes.Equal(code, create2FactoryCode) {
			f.txpool.options.Logger.Infof("using self-deployed CREATE2 factory at %s", candidate.Hex())
			f.factoryAddr = candidate
			f.isInitialized = true
			return nil
		} else {
			f.txpool.options.Logger.Warnf("previously deployed CREATE2 factory at %s is not the expected code", candidate.Hex())
		}
	}

	// 3. deploy the factory ourselves
	addr, err := f.deployFactory(ctx, client)
	if err != nil {
		return err
	}

	f.factoryAddr = addr
	f.isInitialized = true
	return nil
}

func (f *DeploymentFactory) deployFactory(ctx context.Context, client *Client) (common.Address, error) {
	feeCap, tipCap, err := client.GetSuggestedFee(ctx)
	if err != nil {
		return common.Address{}, err
	}
	if feeCap.Cmp(big.NewInt(400000000000)) < 0 {
		feeCap = big.NewInt(400000000000)
	}
	if tipCap.Cmp(big.NewInt(200000000000)) < 0 {
		tipCap = big.NewInt(200000000000)
	}

	// Estimate the deploy gas so we stay correct under EIP-8037 where
	// account creation + per-byte code deposit dominate the cost. Falls back
	// to the formula upper bound if the RPC path fails.
	deployGas, estErr := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  f.rootWallet.GetAddress(),
		To:    nil,
		Value: new(big.Int),
		Data:  create2FactoryInitcode,
	})
	if estErr != nil || deployGas == 0 {
		deployGas = fallbackDeployGas(len(create2FactoryInitcode), f.txpool.IsAmsterdam(), f.txpool.GetCostPerStateByte())
	} else if f.txpool.IsAmsterdam() {
		deployGas = deployGas * 12 / 10
	} else {
		deployGas = deployGas * 11 / 10
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       deployGas,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      create2FactoryInitcode,
	})
	if err != nil {
		return common.Address{}, err
	}

	tx, err := f.rootWallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return common.Address{}, err
	}

	factoryAddr := crypto.CreateAddress(f.rootWallet.GetAddress(), tx.Nonce())
	f.isDeploying = true

	err = f.txpool.SendTransaction(ctx, f.rootWallet, tx, &SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			f.isDeploying = false
			if receipt.Status == types.ReceiptStatusSuccessful {
				f.txpool.options.Logger.Infof("deployed CREATE2 factory at %s (block #%v)", factoryAddr.Hex(), receipt.BlockNumber.Uint64())
			} else {
				f.txpool.options.Logger.Errorf("failed to deploy CREATE2 factory (block #%v)", receipt.BlockNumber.Uint64())
			}
		},
	})
	if err != nil {
		return common.Address{}, err
	}

	f.txpool.options.Logger.Infof("deploying CREATE2 factory at %s", factoryAddr.Hex())
	return factoryAddr, nil
}

func (f *DeploymentFactory) GetContractDeployment(ctx context.Context, initcode []byte, salt [32]byte, client *Client, deployer *Wallet, feeCap *big.Int, tipCap *big.Int, submit bool) (common.Address, *types.Transaction, error) {
	err := f.lazyInit(ctx)
	if err != nil {
		return common.Address{}, nil, err
	}

	contractInitHash := crypto.Keccak256(initcode)
	contractAddr := crypto.CreateAddress2(f.factoryAddr, salt, contractInitHash)

	if client == nil {
		client = f.txpool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientByIndex, 0), WithoutBuilder())
		if client == nil {
			return common.Address{}, nil, fmt.Errorf("no client available")
		}
	}

	if code, err := client.GetCodeAt(ctx, contractAddr); err == nil && len(code) > 0 {
		return contractAddr, nil, nil
	}

	if feeCap == nil || tipCap == nil {
		clientFeeCap, clientTipCap, err := client.GetSuggestedFee(ctx)
		if err != nil {
			return common.Address{}, nil, err
		}
		if clientFeeCap.Cmp(big.NewInt(400000000000)) < 0 {
			clientFeeCap = big.NewInt(400000000000)
		}
		if clientTipCap.Cmp(big.NewInt(200000000000)) < 0 {
			clientTipCap = big.NewInt(200000000000)
		}

		if feeCap == nil {
			feeCap = clientFeeCap
		}
		if tipCap == nil {
			tipCap = clientTipCap
		}
	}

	// Estimate the deploy gas so we stay correct under EIP-8037 where
	// account creation + per-byte code deposit dominate the cost. Falls back
	// to the formula upper bound if the RPC path fails.

	deployData := append(salt[:], initcode...)
	deployGas, estErr := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  deployer.GetAddress(),
		To:    &f.factoryAddr,
		Value: new(big.Int),
		Data:  deployData,
	})
	if estErr != nil || deployGas == 0 || f.isDeploying {
		deployGas = fallbackDeployGas(len(deployData), f.txpool.IsAmsterdam(), f.txpool.GetCostPerStateByte()) + 50_000
	} else if f.txpool.IsAmsterdam() {
		deployGas = deployGas * 12 / 10
	} else {
		deployGas = deployGas * 11 / 10
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       deployGas,
		To:        &f.factoryAddr,
		Value:     uint256.NewInt(0),
		Data:      deployData,
	})
	if err != nil {
		return common.Address{}, nil, err
	}

	tx, err := deployer.BuildDynamicFeeTx(txData)
	if err != nil {
		return common.Address{}, nil, err
	}

	if submit {
		err = f.txpool.SendTransaction(ctx, deployer, tx, &SendTransactionOptions{
			Client:      client,
			Rebroadcast: true,
		})
		if err != nil {
			return common.Address{}, nil, err
		}
	}

	return contractAddr, tx, nil
}
