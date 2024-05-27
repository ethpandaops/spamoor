package txbuilder

import (
	"context"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Timeout time.Duration
	rpchost string
	client  *ethclient.Client
	logger  *logrus.Entry

	gasSuggestionMutex sync.Mutex
	lastGasSuggestion  time.Time
	lastGasCap         *big.Int
	lastTipCap         *big.Int

	blockHeight      uint64
	blockHeightTime  time.Time
	blockHeightMutex sync.Mutex

	awaitNonceMutex     sync.Mutex
	awaitNonceWalletMap map[common.Address]*clientNonceAwait

	awaitNextBlockHeight uint64
	awaitNextBlockMutex  sync.Mutex
	awaitNextBlockWaiter *sync.RWMutex
}

type clientNonceAwait struct {
	mutex        sync.Mutex
	running      bool
	awaitNonces  map[uint64]*sync.RWMutex
	blockHeight  uint64
	errorResult  error
	currentNonce uint64
}

func NewClient(rpchost string) (*Client, error) {
	headers := map[string]string{}

	if strings.HasPrefix(rpchost, "headers(") {

		headersEnd := strings.Index(rpchost, ")")
		headersStr := rpchost[8:headersEnd]
		rpchost = rpchost[headersEnd+1:]

		for _, headerStr := range strings.Split(headersStr, "|") {
			headerParts := strings.Split(headerStr, ":")
			headers[strings.Trim(headerParts[0], " ")] = strings.Trim(headerParts[1], " ")
		}
	}

	ctx := context.Background()
	rpcClient, err := rpc.DialContext(ctx, rpchost)
	if err != nil {
		return nil, err
	}

	for hKey, hVal := range headers {
		rpcClient.SetHeader(hKey, hVal)
	}

	return &Client{
		client:              ethclient.NewClient(rpcClient),
		rpchost:             rpchost,
		logger:              logrus.WithField("client", rpchost),
		awaitNonceWalletMap: make(map[common.Address]*clientNonceAwait),
	}, nil
}

func (client *Client) GetName() string {
	url, _ := url.Parse(client.rpchost)
	name := strings.TrimSuffix(url.Host, ".ethpandaops.io")
	return name
}

func (client *Client) GetEthClient() *ethclient.Client {
	return client.client
}

func (client *Client) GetRPCHost() string {
	return client.rpchost
}

func (client *Client) UpdateWallet(wallet *Wallet) error {
	if wallet.GetChainId() == nil {
		chainId, err := client.GetChainId()
		if err != nil {
			return err
		}
		wallet.SetChainId(chainId)
	}

	nonce, err := client.GetNonceAt(wallet.GetAddress())
	if err != nil {
		return err
	}
	wallet.SetNonce(nonce)

	balance, err := client.GetBalanceAt(wallet.GetAddress())
	if err != nil {
		return err
	}
	wallet.SetBalance(balance)

	return nil
}

func (client *Client) getContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if client.Timeout > 0 {
		return context.WithTimeout(ctx, client.Timeout)
	}
	return context.WithCancel(ctx)
}

func (client *Client) GetChainId() (*big.Int, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.ChainID(ctx)
}

func (client *Client) GetNonceAt(wallet common.Address) (uint64, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.NonceAt(ctx, wallet, nil)
}

func (client *Client) GetPendingNonceAt(wallet common.Address) (uint64, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.PendingNonceAt(ctx, wallet)
}

func (client *Client) GetBalanceAt(wallet common.Address) (*big.Int, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.BalanceAt(ctx, wallet, nil)
}

func (client *Client) GetSuggestedFee() (*big.Int, *big.Int, error) {
	client.gasSuggestionMutex.Lock()
	defer client.gasSuggestionMutex.Unlock()

	if time.Since(client.lastGasSuggestion) < 12*time.Second {
		return client.lastGasCap, client.lastTipCap, nil
	}

	ctx, cancel := client.getContext()
	defer cancel()

	gasCap, err := client.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}
	tipCap, err := client.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, err
	}

	client.lastGasSuggestion = time.Now()
	client.lastGasCap = gasCap
	client.lastTipCap = tipCap
	return gasCap, tipCap, nil
}

func (client *Client) SendTransaction(tx *types.Transaction) error {
	client.logger.Tracef("submitted transaction %v", tx.Hash().String())

	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.SendTransaction(ctx, tx)
}

func (client *Client) SubmitTransaction(txBytes []byte) *common.Hash {
	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(txBytes)
	if err != nil {
		client.logger.Errorf("Error while decoding transaction: %v (%v)", err, len(txBytes))
		return nil
	}

	ctx, cancel := client.getContext()
	defer cancel()

	err = client.client.SendTransaction(ctx, tx)
	if err != nil {
		client.logger.Errorf("Error while sending transaction: %v", err)
		return nil
	}

	txHash := tx.Hash()
	client.logger.Tracef("submitted transaction %v", txHash.String())

	return &txHash
}

func (client *Client) GetTransactionReceipt(txHash []byte) *types.Receipt {
	hash := common.Hash{}
	hash.SetBytes(txHash)

	client.logger.Tracef("get receipt: 0x%x", txHash)

	ctx, cancel := client.getContext()
	defer cancel()

	receipt, err := client.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil
	}

	return receipt
}

func (client *Client) GetBlockHeight() (uint64, error) {
	client.blockHeightMutex.Lock()
	defer client.blockHeightMutex.Unlock()

	if time.Since(client.blockHeightTime) < 12*time.Second {
		return client.blockHeight, nil
	}

	client.logger.Tracef("get block number")

	ctx, cancel := client.getContext()
	defer cancel()

	blockHeight, err := client.client.BlockNumber(ctx)
	if err != nil {
		return blockHeight, err
	}
	if blockHeight > client.blockHeight {
		client.blockHeight = blockHeight
		client.blockHeightTime = time.Now()
	}
	return client.blockHeight, nil
}

func (client *Client) AwaitTransaction(tx *types.Transaction) (*types.Receipt, uint64, error) {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return nil, 0, err
	}

	blockHeight, err := client.AwaitWalletNonce(from, tx.Nonce()+1, 0)
	if err != nil {
		return nil, 0, err
	}

	client.logger.Tracef("get receipt: %v", tx.Hash().String())

	ctx, cancel := client.getContext()
	defer cancel()

	receipt, err := client.client.TransactionReceipt(ctx, tx.Hash())
	if receipt != nil {
		return receipt, blockHeight, nil
	}

	if err != nil && err.Error() != "not found" {
		client.logger.Warnf("receipt error: %v\n", err)
		return nil, blockHeight, err
	}

	return nil, blockHeight, nil
}

func (client *Client) AwaitWalletNonce(wallet common.Address, nonce uint64, blockHeight uint64) (uint64, error) {
	var awaitMutex *sync.RWMutex

	client.awaitNonceMutex.Lock()
	walletNonceAwaiter := client.awaitNonceWalletMap[wallet]
	if walletNonceAwaiter == nil {
		walletNonceAwaiter = &clientNonceAwait{
			awaitNonces: make(map[uint64]*sync.RWMutex),
			blockHeight: blockHeight,
		}
		client.awaitNonceWalletMap[wallet] = walletNonceAwaiter
	}
	client.awaitNonceMutex.Unlock()

	walletNonceAwaiter.mutex.Lock()
	if nonce <= walletNonceAwaiter.currentNonce {
		walletNonceAwaiter.mutex.Unlock()
		return walletNonceAwaiter.blockHeight, nil
	}
	awaitMutex = walletNonceAwaiter.awaitNonces[nonce]
	if awaitMutex == nil {
		awaitMutex = &sync.RWMutex{}
		awaitMutex.Lock()
		walletNonceAwaiter.awaitNonces[nonce] = awaitMutex
	}

	if !walletNonceAwaiter.running {
		walletNonceAwaiter.running = true
		go func() {
			var err error
			retryCount := 0
			for {
				err = client.processWalletNonceAwaiter(wallet, walletNonceAwaiter)
				if err == nil {
					break
				}
				client.logger.Warnf("error while awaiting nonce inclusion: %v", err)
				retryCount++
				if retryCount > 10 {
					break
				}
				time.Sleep(2 * time.Second)
			}
			if err != nil {
				// can't check nonce - client is probably dead or unsynced
				// cancel nonceAwaiter, bubble up error
				walletNonceAwaiter.errorResult = err
				client.disposeWalletNonceAwaiter(wallet, walletNonceAwaiter, true)
				for nonce, mtx := range walletNonceAwaiter.awaitNonces {
					mtx.Unlock()
					delete(walletNonceAwaiter.awaitNonces, nonce)
				}
			}
		}()
	}
	walletNonceAwaiter.mutex.Unlock()

	awaitMutex.RLock()
	defer awaitMutex.RUnlock()
	return walletNonceAwaiter.blockHeight, walletNonceAwaiter.errorResult
}

func (client *Client) processWalletNonceAwaiter(wallet common.Address, awaiter *clientNonceAwait) error {
	if awaiter.blockHeight == 0 {
		client.logger.Tracef("get block number")

		ctx, cancel := client.getContext()
		defer cancel()

		blockHeight, err := client.client.BlockNumber(ctx)
		if err != nil {
			return err
		}
		awaiter.blockHeight = blockHeight
	}

	fetchLatestNonce := func(blockHeight uint64) error {
		client.logger.Tracef("get nonce for %v at %v", wallet.String(), blockHeight)

		ctx, cancel := client.getContext()
		defer cancel()

		currentNonce, err := client.client.NonceAt(ctx, wallet, big.NewInt(int64(blockHeight)))
		if err != nil {
			return err
		}
		if currentNonce <= awaiter.currentNonce {
			awaiter.blockHeight = blockHeight
			return nil
		}

		awaiter.mutex.Lock()
		defer awaiter.mutex.Unlock()
		awaiter.blockHeight = blockHeight
		awaiter.currentNonce = currentNonce
		for nonce, mtx := range awaiter.awaitNonces {
			if nonce <= currentNonce {
				delete(awaiter.awaitNonces, nonce)
				mtx.Unlock()
			}
		}
		return nil
	}

	for {
		err := fetchLatestNonce(awaiter.blockHeight)
		if err != nil {
			return err
		}

		// break loop if no more nonces to wait for
		awaiter.mutex.Lock()
		awaitNonceCount := len(awaiter.awaitNonces)
		if awaitNonceCount == 0 {
			client.disposeWalletNonceAwaiter(wallet, awaiter, false)
			return nil
		}
		awaiter.mutex.Unlock()

		// await next block
		awaiter.blockHeight, err = client.AwaitNextBlock(awaiter.blockHeight)
		if err != nil {
			return err
		}
	}
}

func (client *Client) disposeWalletNonceAwaiter(wallet common.Address, awaiter *clientNonceAwait, lock bool) {
	if lock {
		awaiter.mutex.Lock()
	}
	awaiter.running = false
	awaiter.mutex.Unlock()
	client.awaitNonceMutex.Lock()
	delete(client.awaitNonceWalletMap, wallet)
	client.awaitNonceMutex.Unlock()
}

func (client *Client) AwaitNextBlock(lastBlockHeight uint64) (uint64, error) {
	if client.awaitNextBlockHeight > lastBlockHeight {
		return client.awaitNextBlockHeight, nil
	}
	client.awaitNextBlockMutex.Lock()
	if client.awaitNextBlockWaiter != nil {
		waitMutex := client.awaitNextBlockWaiter
		client.awaitNextBlockMutex.Unlock()

		waitMutex.RLock()
		defer waitMutex.RUnlock()
		return client.awaitNextBlockHeight, nil
	}
	client.awaitNextBlockWaiter = &sync.RWMutex{}
	client.awaitNextBlockWaiter.Lock()
	defer func() {
		client.awaitNextBlockWaiter.Unlock()
		client.awaitNextBlockMutex.Lock()
		client.awaitNextBlockWaiter = nil
		client.awaitNextBlockMutex.Unlock()
	}()
	client.awaitNextBlockMutex.Unlock()

	if lastBlockHeight == 0 {
		var err error
		lastBlockHeight, err = client.GetBlockHeight()
		if err != nil {
			return 0, err
		}
	}
	client.awaitNextBlockHeight = lastBlockHeight
	for {
		time.Sleep(1 * time.Second)

		client.logger.Tracef("get block number")

		ctx, cancel := client.getContext()
		blockHeight, err := client.client.BlockNumber(ctx)
		cancel()

		if err != nil {
			return lastBlockHeight, err
		}
		client.blockHeight = blockHeight
		client.blockHeightTime = time.Now()
		if blockHeight > lastBlockHeight {
			client.awaitNextBlockHeight = blockHeight
			return blockHeight, nil
		}
	}
}
