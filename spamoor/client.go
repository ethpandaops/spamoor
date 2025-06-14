package spamoor

import (
	"context"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoortypes"
)

// rpcClient represents an Ethereum RPC client with additional functionality for transaction management,
// gas estimation caching, and block height tracking. It wraps the standard go-ethereum ethclient
// with enhanced features for spam testing and transaction automation.
type rpcClient struct {
	rpchost   string
	client    *ethclient.Client
	rpcClient *rpc.Client
	logger    *logrus.Entry

	timeout     time.Duration
	clientGroup string
	enabled     bool

	gasSuggestionMutex sync.Mutex
	lastGasSuggestion  time.Time
	lastGasCap         *big.Int
	lastTipCap         *big.Int

	blockHeight      uint64
	blockHeightTime  time.Time
	blockHeightMutex sync.Mutex

	clientVersion     string
	clientVersionTime time.Time
}

// NewClient creates a new rpcClient instance with the specified RPC host URL and returns it as Client interface.
// The rpchost parameter supports special prefixes:
//   - headers(key:value|key2:value2) - sets custom HTTP headers
//   - group(name) - assigns the client to a named group
//
// Example: "headers(Authorization:Bearer token|User-Agent:MyApp)group(mainnet)http://localhost:8545"
func NewClient(rpchost string) (spamoortypes.Client, error) {
	headers := map[string]string{}
	clientGroup := "default"

	for {
		if strings.HasPrefix(rpchost, "headers(") {
			headersEnd := strings.Index(rpchost, ")")
			headersStr := rpchost[8:headersEnd]
			rpchost = rpchost[headersEnd+1:]

			for _, headerStr := range strings.Split(headersStr, "|") {
				headerParts := strings.Split(headerStr, ":")
				headers[strings.Trim(headerParts[0], " ")] = strings.Trim(headerParts[1], " ")
			}
		} else if strings.HasPrefix(rpchost, "group(") {
			groupEnd := strings.Index(rpchost, ")")
			groupStr := rpchost[6:groupEnd]
			rpchost = rpchost[groupEnd+1:]
			clientGroup = groupStr
		} else {
			break
		}
	}

	ctx := context.Background()
	client, err := rpc.DialContext(ctx, rpchost)
	if err != nil {
		return nil, err
	}

	for hKey, hVal := range headers {
		client.SetHeader(hKey, hVal)
	}

	return &rpcClient{
		client:      ethclient.NewClient(client),
		rpcClient:   client,
		rpchost:     rpchost,
		logger:      logrus.WithField("rpc", rpchost),
		clientGroup: clientGroup,
		enabled:     true,
	}, nil
}

// GetName returns a shortened name for the client derived from the RPC host URL,
// removing common suffixes like ".ethpandaops.io".
func (client *rpcClient) GetName() string {
	url, _ := url.Parse(client.rpchost)
	name := strings.TrimSuffix(url.Host, ".ethpandaops.io")
	return name
}

// GetClientGroup returns the client group name assigned during initialization.
// Defaults to "default" if no group was specified.
func (client *rpcClient) GetClientGroup() string {
	return client.clientGroup
}

// GetEthClient returns the underlying go-ethereum ethclient.Client instance.
func (client *rpcClient) GetEthClient() bind.ContractBackend {
	return client.client
}

// GetRPCHost returns the original RPC host URL used to create this client.
func (client *rpcClient) GetRPCHost() string {
	return client.rpchost
}

// GetTimeout returns the timeout for the client.
func (client *rpcClient) GetTimeout() time.Duration {
	return client.timeout
}

// SetTimeout sets the timeout for the client.
func (client *rpcClient) SetTimeout(timeout time.Duration) {
	client.timeout = timeout
}

// UpdateWallet refreshes the wallet's chain ID, nonce, and balance by querying the blockchain.
// If the wallet doesn't have a chain ID set, it will be fetched and assigned.
func (client *rpcClient) UpdateWallet(ctx context.Context, wallet spamoortypes.Wallet) error {
	if wallet.GetChainId() == nil {
		chainId, err := client.GetChainId(ctx)
		if err != nil {
			return err
		}
		wallet.SetChainId(chainId)
	}

	nonce, err := client.GetNonceAt(ctx, wallet.GetAddress(), nil)
	if err != nil {
		return err
	}
	wallet.SetNonce(nonce)

	balance, err := client.GetBalanceAt(ctx, wallet.GetAddress())
	if err != nil {
		return err
	}
	wallet.SetBalance(balance)

	return nil
}

// SetClientGroup sets the client group name for the client.
// This is used to group clients together and target them with specific scenarios.
func (client *rpcClient) SetClientGroup(group string) {
	client.clientGroup = group
}

// IsEnabled returns whether the client is enabled for selection.
func (client *rpcClient) IsEnabled() bool {
	return client.enabled
}

// SetEnabled sets the enabled state of the client.
// Disabled clients will not be considered for selection in the client pool.
func (client *rpcClient) SetEnabled(enabled bool) {
	client.enabled = enabled
}

func (client *rpcClient) getContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if client.timeout > 0 {
		return context.WithTimeout(ctx, client.timeout)
	}
	return context.WithCancel(ctx)
}

// GetChainId returns the chain ID of the connected Ethereum network.
func (client *rpcClient) GetChainId(ctx context.Context) (*big.Int, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.ChainID(ctx)
}

// GetNonceAt returns the nonce for the given address at the specified block number.
// If blockNumber is nil, returns the nonce at the latest block.
func (client *rpcClient) GetNonceAt(ctx context.Context, wallet common.Address, blockNumber *big.Int) (uint64, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.NonceAt(ctx, wallet, blockNumber)
}

// GetPendingNonceAt returns the pending nonce for the given address,
// including transactions in the mempool.
func (client *rpcClient) GetPendingNonceAt(ctx context.Context, wallet common.Address) (uint64, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.PendingNonceAt(ctx, wallet)
}

// GetBalanceAt returns the balance of the given address at the latest block.
func (client *rpcClient) GetBalanceAt(ctx context.Context, wallet common.Address) (*big.Int, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.BalanceAt(ctx, wallet, nil)
}

// GetSuggestedFee returns suggested gas price and tip cap for transactions.
// Results are cached for 12 seconds to reduce RPC calls.
// Returns (gasCap, tipCap, error).
func (client *rpcClient) GetSuggestedFee(ctx context.Context) (*big.Int, *big.Int, error) {
	client.gasSuggestionMutex.Lock()
	defer client.gasSuggestionMutex.Unlock()

	if time.Since(client.lastGasSuggestion) < 12*time.Second {
		return client.lastGasCap, client.lastTipCap, nil
	}

	ctx, cancel := client.getContext(ctx)
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

// SendTransaction submits a transaction to the network using the provided context.
// Logs the transaction hash at trace level.
func (client *rpcClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	client.logger.Tracef("submitted transaction %v", tx.Hash().String())

	return client.client.SendTransaction(ctx, tx)
}

// SendRawTransaction submits a raw transaction bytes to the network using eth_sendRawTransaction RPC call.
func (client *rpcClient) SendRawTransaction(ctx context.Context, tx []byte) error {
	return client.client.Client().CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(tx))
}

// GetTransactionReceipt retrieves the receipt for a given transaction hash.
// Logs the request at trace level.
func (client *rpcClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	client.logger.Tracef("get receipt: 0x%x", txHash.Bytes())

	return client.client.TransactionReceipt(ctx, txHash)
}

// GetBlockHeight returns the current block number.
// Results are cached for 12 seconds to reduce RPC calls.
func (client *rpcClient) GetBlockHeight(ctx context.Context) (uint64, error) {
	client.blockHeightMutex.Lock()
	defer client.blockHeightMutex.Unlock()

	if time.Since(client.blockHeightTime) < 12*time.Second {
		return client.blockHeight, nil
	}

	client.logger.Tracef("get block number")

	ctx, cancel := client.getContext(ctx)
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

// GetLastBlockHeight returns the last cached block height and the time it was retrieved.
func (client *rpcClient) GetLastBlockHeight() (uint64, time.Time) {
	return client.blockHeight, client.blockHeightTime
}

// GetClientVersion returns the client version string from the web3_clientVersion RPC call.
// Results are cached for 30 minutes to reduce RPC calls.
func (client *rpcClient) GetClientVersion(ctx context.Context) (string, error) {
	if time.Since(client.clientVersionTime) < 30*time.Minute {
		return client.clientVersion, nil
	}

	var result string
	err := client.rpcClient.CallContext(ctx, &result, "web3_clientVersion")
	if err != nil {
		return client.clientVersion, err
	}

	client.clientVersion = result
	client.clientVersionTime = time.Now()

	return result, nil
}

// GetBlock retrieves a block by its number.
func (client *rpcClient) GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

// GetBlockReceipts retrieves all transaction receipts for a block.
func (client *rpcClient) GetBlockReceipts(ctx context.Context, blockNumber uint64) ([]*types.Receipt, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	blockNum := rpc.BlockNumber(blockNumber)

	return client.client.BlockReceipts(ctx, rpc.BlockNumberOrHash{
		BlockNumber: &blockNum,
	})
}

// GetLatestGasLimit returns the latest gas limit from the latest block.
func (client *rpcClient) GetLatestGasLimit(ctx context.Context) (uint64, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	header, err := client.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return header.GasLimit, nil
}
