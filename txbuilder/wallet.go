package txbuilder

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

type Wallet struct {
	nonceMutex   sync.Mutex
	balanceMutex sync.RWMutex
	privkey      *ecdsa.PrivateKey
	address      common.Address
	chainid      *big.Int
	nonce        uint64
	balance      *big.Int
}

func NewWallet(privkey string) (*Wallet, error) {
	wallet := &Wallet{}
	err := wallet.loadPrivateKey(privkey)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wallet *Wallet) loadPrivateKey(privkey string) error {
	var privateKey *ecdsa.PrivateKey
	if privkey == "" {
		var err error
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			return err
		}
	} else {
		var err error
		privateKey, err = crypto.HexToECDSA(privkey)
		if err != nil {
			return err
		}
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	wallet.privkey = privateKey
	wallet.address = crypto.PubkeyToAddress(*publicKeyECDSA)
	return nil
}

func (wallet *Wallet) GetAddress() common.Address {
	return wallet.address
}

func (wallet *Wallet) GetChainId() *big.Int {
	return wallet.chainid
}

func (wallet *Wallet) GetNonce() uint64 {
	return wallet.nonce
}

func (wallet *Wallet) GetBalance() *big.Int {
	wallet.balanceMutex.RLock()
	defer wallet.balanceMutex.RUnlock()
	return wallet.balance
}

func (wallet *Wallet) SetChainId(chainid *big.Int) {
	wallet.chainid = chainid
}

func (wallet *Wallet) SetNonce(nonce uint64) {
	wallet.nonce = nonce
}

func (wallet *Wallet) SetBalance(balance *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = balance
}

func (wallet *Wallet) SubBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Sub(wallet.balance, amount)
}

func (wallet *Wallet) AddBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Add(wallet.balance, amount)
}

func (wallet *Wallet) BuildDynamicFeeTx(txData *types.DynamicFeeTx) (*types.Transaction, error) {
	wallet.nonceMutex.Lock()
	txData.ChainID = wallet.chainid
	txData.Nonce = wallet.nonce
	wallet.nonce++
	wallet.nonceMutex.Unlock()
	return wallet.signTx(txData)
}

func (wallet *Wallet) BuildBlobTx(txData *types.BlobTx) (*types.Transaction, error) {
	wallet.nonceMutex.Lock()
	txData.ChainID = uint256.MustFromBig(wallet.chainid)
	txData.Nonce = wallet.nonce
	wallet.nonce++
	wallet.nonceMutex.Unlock()
	return wallet.signTx(txData)
}

func (wallet *Wallet) BuildBoundTx(txData *TxMetadata, buildFn func(transactOpts *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	transactor, err := bind.NewKeyedTransactorWithChainID(wallet.privkey, wallet.chainid)
	if err != nil {
		return nil, err
	}

	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	transactor.Context = context.Background()
	transactor.From = wallet.address
	transactor.Nonce = big.NewInt(0).SetUint64(wallet.nonce)
	wallet.nonce++

	transactor.GasTipCap = txData.GasTipCap.ToBig()
	transactor.GasFeeCap = txData.GasFeeCap.ToBig()
	transactor.GasLimit = txData.Gas
	transactor.Value = txData.Value.ToBig()
	transactor.NoSend = true

	tx, err := buildFn(transactor)
	if err != nil {
		wallet.nonce--
		return nil, err
	}

	return tx, nil
}

func (wallet *Wallet) ReplaceDynamicFeeTx(txData *types.DynamicFeeTx, nonce uint64) (*types.Transaction, error) {
	txData.ChainID = wallet.chainid
	txData.Nonce = nonce
	return wallet.signTx(txData)
}

func (wallet *Wallet) ReplaceBlobTx(txData *types.BlobTx, nonce uint64) (*types.Transaction, error) {
	txData.ChainID = uint256.MustFromBig(wallet.chainid)
	txData.Nonce = nonce
	return wallet.signTx(txData)
}

func (wallet *Wallet) signTx(txData types.TxData) (*types.Transaction, error) {
	tx := types.NewTx(txData)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(wallet.chainid), wallet.privkey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
