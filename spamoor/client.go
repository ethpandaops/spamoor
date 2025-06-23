package spamoor

import (
	"context"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

// Client represents an Ethereum RPC client with additional functionality for transaction management,
// gas estimation caching, and block height tracking. It wraps the standard go-ethereum ethclient
// with enhanced features for spam testing and transaction automation.
type Client struct {
	Timeout   time.Duration
	rpchost   string
	client    *ethclient.Client
	rpcClient *rpc.Client
	logger    *logrus.Entry

	clientGroups []string
	enabled      bool
	nameOverride string

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

// NewClient creates a new Client instance with the specified RPC host URL.
// The rpchost parameter supports special prefixes:
//   - headers(key:value|key2:value2) - sets custom HTTP headers
//   - group(name) - assigns the client to a named group (can be used multiple times)
//   - group(name1,name2,name3) - assigns the client to multiple groups (comma-separated)
//   - name(custom_name) - sets a custom display name override
//
// Example: "headers(Authorization:Bearer token|User-Agent:MyApp)group(mainnet)group(primary)name(My Custom Node)http://localhost:8545"
// Example: "group(mainnet,primary,backup)name(MainNet Primary)http://localhost:8545"
func NewClient(rpchost string) (*Client, error) {
	headers := map[string]string{}
	clientGroups := []string{"default"}
	nameOverride := ""

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

			// Parse comma-separated groups
			for _, group := range strings.Split(groupStr, ",") {
				group = strings.TrimSpace(group)
				if group != "" {
					// Check if group already exists to avoid duplicates
					exists := false
					for _, existing := range clientGroups {
						if existing == group {
							exists = true
							break
						}
					}
					if !exists {
						clientGroups = append(clientGroups, group)
					}
				}
			}
		} else if strings.HasPrefix(rpchost, "name(") {
			nameEnd := strings.Index(rpchost, ")")
			nameOverride = rpchost[5:nameEnd]
			rpchost = rpchost[nameEnd+1:]
		} else {
			break
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
		client:       ethclient.NewClient(rpcClient),
		rpcClient:    rpcClient,
		rpchost:      rpchost,
		logger:       logrus.WithField("rpc", rpchost),
		clientGroups: clientGroups,
		enabled:      true,
		nameOverride: nameOverride,
	}, nil
}

// GetName returns a shortened name for the client derived from the RPC host URL,
// removing common suffixes like ".ethpandaops.io".
// If a name override is set, it returns the override instead.
func (client *Client) GetName() string {
	if client.nameOverride != "" {
		return client.nameOverride
	}
	url, _ := url.Parse(client.rpchost)
	name := strings.TrimSuffix(url.Host, ".ethpandaops.io")
	return name
}

// GetClientGroup returns the first client group name assigned during initialization.
// Defaults to "default" if no group was specified. For multiple groups, use GetClientGroups().
func (client *Client) GetClientGroup() string {
	if len(client.clientGroups) > 0 {
		return client.clientGroups[0]
	}
	return "default"
}

// GetClientGroups returns all client group names assigned to this client.
func (client *Client) GetClientGroups() []string {
	if len(client.clientGroups) == 0 {
		return []string{"default"}
	}
	return client.clientGroups
}

// HasGroup checks if the client belongs to the specified group.
func (client *Client) HasGroup(group string) bool {
	if group == "" {
		group = "default"
	}
	for _, clientGroup := range client.clientGroups {
		if clientGroup == group {
			return true
		}
	}
	return false
}

// GetEthClient returns the underlying go-ethereum ethclient.Client instance.
func (client *Client) GetEthClient() *ethclient.Client {
	return client.client
}

// GetRPCHost returns the original RPC host URL used to create this client.
func (client *Client) GetRPCHost() string {
	return client.rpchost
}

// UpdateWallet refreshes the wallet's chain ID, nonce, and balance by querying the blockchain.
// If the wallet doesn't have a chain ID set, it will be fetched and assigned.
func (client *Client) UpdateWallet(ctx context.Context, wallet *Wallet) error {
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

// SetClientGroups sets multiple client group names for the client, replacing all existing groups.
func (client *Client) SetClientGroups(groups []string) {
	// Remove duplicates and empty strings
	uniqueGroups := []string{}
	for _, group := range groups {
		group = strings.TrimSpace(group)
		if group == "" {
			continue
		}
		exists := false
		for _, existing := range uniqueGroups {
			if existing == group {
				exists = true
				break
			}
		}
		if !exists {
			uniqueGroups = append(uniqueGroups, group)
		}
	}

	client.clientGroups = uniqueGroups
}

// IsEnabled returns whether the client is enabled for selection.
func (client *Client) IsEnabled() bool {
	return client.enabled
}

// SetEnabled sets the enabled state of the client.
// Disabled clients will not be considered for selection in the client pool.
func (client *Client) SetEnabled(enabled bool) {
	client.enabled = enabled
}

// GetNameOverride returns the name override for the client.
func (client *Client) GetNameOverride() string {
	return client.nameOverride
}

// SetNameOverride sets a custom name override for the client.
// If set, this name will be used instead of the auto-generated name from the RPC host.
func (client *Client) SetNameOverride(name string) {
	client.nameOverride = name
}

func (client *Client) getContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if client.Timeout > 0 {
		return context.WithTimeout(ctx, client.Timeout)
	}
	return context.WithCancel(ctx)
}

// GetChainId returns the chain ID of the connected Ethereum network.
func (client *Client) GetChainId(ctx context.Context) (*big.Int, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.ChainID(ctx)
}

// GetNonceAt returns the nonce for the given address at the specified block number.
// If blockNumber is nil, returns the nonce at the latest block.
func (client *Client) GetNonceAt(ctx context.Context, wallet common.Address, blockNumber *big.Int) (uint64, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.NonceAt(ctx, wallet, blockNumber)
}

// GetPendingNonceAt returns the pending nonce for the given address,
// including transactions in the mempool.
func (client *Client) GetPendingNonceAt(ctx context.Context, wallet common.Address) (uint64, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.PendingNonceAt(ctx, wallet)
}

// GetBalanceAt returns the balance of the given address at the latest block.
func (client *Client) GetBalanceAt(ctx context.Context, wallet common.Address) (*big.Int, error) {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	return client.client.BalanceAt(ctx, wallet, nil)
}

// GetSuggestedFee returns suggested gas price and tip cap for transactions.
// Results are cached for 12 seconds to reduce RPC calls.
// Returns (gasCap, tipCap, error).
func (client *Client) GetSuggestedFee(ctx context.Context) (*big.Int, *big.Int, error) {
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
func (client *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	client.logger.Tracef("submitted transaction %v", tx.Hash().String())

	return client.client.SendTransaction(ctx, tx)
}

// SendRawTransaction submits a raw transaction bytes to the network using eth_sendRawTransaction RPC call.
func (client *Client) SendRawTransaction(ctx context.Context, tx []byte) error {
	return client.client.Client().CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(tx))
}

// GetTransactionReceipt retrieves the receipt for a given transaction hash.
// Logs the request at trace level.
func (client *Client) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	client.logger.Tracef("get receipt: 0x%x", txHash.Bytes())

	return client.client.TransactionReceipt(ctx, txHash)
}

// GetBlockHeight returns the current block number.
// Results are cached for 12 seconds to reduce RPC calls.
func (client *Client) GetBlockHeight(ctx context.Context) (uint64, error) {
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
func (client *Client) GetLastBlockHeight() (uint64, time.Time) {
	return client.blockHeight, client.blockHeightTime
}

// GetClientVersion returns the client version string from the web3_clientVersion RPC call.
// Results are cached for 30 minutes to reduce RPC calls.
func (client *Client) GetClientVersion(ctx context.Context) (string, error) {
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
