package spamoor

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/txtypes"
)

// sharedHTTPClient is reused across all RPC clients to enable TCP connection
// pooling. Go's default MaxIdleConnsPerHost is 2, which forces a new TCP
// connection for nearly every request under any concurrency. This transport
// raises the per-host pool to 100, letting high-throughput callers (like
// transaction spammers) reuse connections instead of churning them.
var sharedHTTPClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

// ClientType represents the type of Ethereum client
type ClientType string

const (
	ClientTypeClient  ClientType = "client"
	ClientTypeBuilder ClientType = "builder"
)

// String returns the string representation of the client type
func (ct ClientType) String() string {
	return string(ct)
}

// IsValid checks if the client type is valid
func (ct ClientType) IsValid() bool {
	switch ct {
	case ClientTypeClient, ClientTypeBuilder:
		return true
	default:
		return false
	}
}

// Client represents an Ethereum RPC client with additional functionality for transaction management,
// gas estimation caching, and block height tracking. Engine calls go through the raw
// JSON-RPC client so that spamoor's own encoding and decoding applies.
type Client struct {
	Timeout        time.Duration
	rpchost        string
	rpcClient      *rpc.Client
	externalClient *ExternalClientOptions
	logger         *logrus.Entry

	ethClient     *ethclient.Client
	ethClientOnce sync.Once

	clientGroups     []string
	clientType       ClientType
	clientTypeFromDB ClientType // Override from database
	enabled          bool
	nameOverride     string

	gasSuggestionMutex sync.Mutex
	lastGasSuggestion  time.Time
	lastGasCap         *big.Int
	lastTipCap         *big.Int

	blockHeight      uint64
	blockHeightTime  time.Time
	blockHeightMutex sync.Mutex

	clientVersion     string
	clientVersionTime time.Time

	// Request statistics
	statsMutex       sync.Mutex
	totalRequests    uint64
	totalTxRequests  uint64
	totalRpcFailures uint64
}

type ClientOptions struct {
	RpcHost        string
	ExternalClient *ExternalClientOptions
}

type ExternalClientOptions struct {
	GetBlockHeight func(ctx context.Context) (uint64, error)
}

// NewClient creates a new Client instance with the specified RPC host URL.
// The rpchost parameter supports special prefixes:
//   - headers(key:value|key2:value2) - sets custom HTTP headers
//   - group(name) - assigns the client to a named group (can be used multiple times)
//   - group(name1,name2,name3) - assigns the client to multiple groups (comma-separated)
//   - group(-name) - removes the client from a group. Every client starts in
//     "default", which is the group used by selections that name no group
//     (wallet funding, deployments, scenarios without --client-group), so
//     "group(mygroup,-default)" is how a client is reserved for scenarios
//     that ask for "mygroup" and reached by nothing else.
//   - name(custom_name) - sets a custom display name override
//   - type(client_type) - sets the client type (e.g., "builder" for builder clients)
//
// Example: "headers(Authorization:Bearer token|User-Agent:MyApp)group(mainnet)group(primary)name(My Custom Node)http://localhost:8545"
// Example: "group(mainnet,primary,backup)name(MainNet Primary)http://localhost:8545"
// Example: "type(builder)group(builders)name(Builder Node)http://localhost:8545"
// Example: "group(private,-default)name(Private Intake)http://localhost:8080/rpc"
func NewClient(options *ClientOptions) (*Client, error) {
	headers := map[string]string{}
	clientGroups := []string{"default"}
	nameOverride := ""
	clientType := ClientTypeClient
	rpchost := options.RpcHost

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

			// Parse comma-separated groups. A "-" prefix removes a group
			// instead of adding it, which is the only way to take a client
			// out of the implicit "default" group.
			for _, group := range strings.Split(groupStr, ",") {
				group = strings.TrimSpace(group)
				if group == "" {
					continue
				}

				if remove, found := strings.CutPrefix(group, "-"); found {
					if remove == "" {
						return nil, fmt.Errorf("invalid client group \"-\": expected a group name after the '-'")
					}

					clientGroups = slices.DeleteFunc(clientGroups, func(existing string) bool {
						return existing == remove
					})

					continue
				}

				if !slices.Contains(clientGroups, group) {
					clientGroups = append(clientGroups, group)
				}
			}
		} else if strings.HasPrefix(rpchost, "name(") {
			nameEnd := strings.Index(rpchost, ")")
			nameOverride = rpchost[5:nameEnd]
			rpchost = rpchost[nameEnd+1:]
		} else if strings.HasPrefix(rpchost, "type(") {
			typeEnd := strings.Index(rpchost, ")")
			typeStr := rpchost[5:typeEnd]
			rpchost = rpchost[typeEnd+1:]

			parsedType := ClientType(typeStr)
			if !parsedType.IsValid() {
				return nil, fmt.Errorf("invalid client type '%s', supported types: client, builder", typeStr)
			}
			clientType = parsedType
		} else {
			break
		}
	}

	ctx := context.Background()

	dialOpts := []rpc.ClientOption{rpc.WithWebsocketMessageSizeLimit(0)}
	if strings.HasPrefix(rpchost, "http://") || strings.HasPrefix(rpchost, "https://") {
		dialOpts = append(dialOpts, rpc.WithHTTPClient(sharedHTTPClient))
	}

	rpcClient, err := rpc.DialOptions(ctx, rpchost, dialOpts...)
	if err != nil {
		return nil, err
	}

	for hKey, hVal := range headers {
		rpcClient.SetHeader(hKey, hVal)
	}

	// An empty group set makes the client unselectable by every code path
	// while GetClientGroups still reports "default", so refuse it rather than
	// silently shipping a client nothing can ever reach.
	if len(clientGroups) == 0 {
		return nil, fmt.Errorf("client %s has no client groups left: the group(...) prefix removed every group", rpchost)
	}

	return &Client{
		rpcClient:    rpcClient,
		rpchost:      rpchost,
		logger:       logrus.WithField("rpc", rpchost),
		clientGroups: clientGroups,
		clientType:   clientType,
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

// GetClientType returns the client type.
// If a database override is set, it returns the override; otherwise returns the parsed type.
func (client *Client) GetClientType() ClientType {
	if client.clientTypeFromDB != "" {
		return client.clientTypeFromDB
	}
	return client.clientType
}

// IsBuilder returns true if the client is a builder type.
// Uses the database override if available, otherwise the parsed type.
func (client *Client) IsBuilder() bool {
	return client.GetClientType() == ClientTypeBuilder
}

// GetEthClient returns a go-ethereum ethclient.Client sharing this client's RPC
// connection, for use as the bind.ContractBackend of abigen contract bindings.
// Constructed on first use.
func (client *Client) GetEthClient() *ethclient.Client {
	client.ethClientOnce.Do(func() {
		client.ethClient = ethclient.NewClient(client.rpcClient)
	})

	return client.ethClient
}

// GetRPCHost returns the original RPC host URL used to create this client.
func (client *Client) GetRPCHost() string {
	return client.rpchost
}

// UpdateWallet refreshes the wallet's chain ID, nonce, and balance by querying the blockchain.
// If the wallet doesn't have a chain ID set, it will be fetched and assigned.
func (client *Client) UpdateWallet(ctx context.Context, wallet *Wallet) error {
	return wallet.UpdateWallet(ctx, client, true)
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

// SetClientTypeOverride sets a client type override from the database.
// If set, this type will be used instead of the type parsed from the RPC URL.
func (client *Client) SetClientTypeOverride(clientType string) {
	if clientType == "" {
		client.clientTypeFromDB = ""
		return
	}

	parsedType := ClientType(clientType)
	if parsedType.IsValid() {
		client.clientTypeFromDB = parsedType
	}
}

func (client *Client) getContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if client.Timeout > 0 {
		return context.WithTimeout(ctx, client.Timeout)
	}
	return context.WithCancel(ctx)
}

// callRPC performs a JSON-RPC call with the client's timeout and records request
// statistics. isTxRequest marks submission calls, which are counted separately.
func (client *Client) callRPC(ctx context.Context, isTxRequest bool, result any, method string, args ...any) error {
	ctx, cancel := client.getContext(ctx)
	defer cancel()

	err := client.rpcClient.CallContext(ctx, result, method, args...)
	client.incrementRequestStats(isTxRequest, err != nil)

	return err
}

// GetChainId returns the chain ID of the connected Ethereum network.
func (client *Client) GetChainId(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	if err := client.callRPC(ctx, false, &result, "eth_chainId"); err != nil {
		return nil, err
	}

	return (*big.Int)(&result), nil
}

// GetNonceAt returns the nonce for the given address at the specified block number.
// If blockNumber is nil, returns the nonce at the latest block.
func (client *Client) GetNonceAt(ctx context.Context, wallet common.Address, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	if err := client.callRPC(ctx, false, &result, "eth_getTransactionCount", wallet, toBlockNumberArg(blockNumber)); err != nil {
		return 0, err
	}

	return uint64(result), nil
}

// GetPendingNonceAt returns the pending nonce for the given address,
// including transactions in the mempool.
func (client *Client) GetPendingNonceAt(ctx context.Context, wallet common.Address) (uint64, error) {
	var result hexutil.Uint64
	if err := client.callRPC(ctx, false, &result, "eth_getTransactionCount", wallet, "pending"); err != nil {
		return 0, err
	}

	return uint64(result), nil
}

// GetBalanceAt returns the balance of the given address at the latest block.
func (client *Client) GetBalanceAt(ctx context.Context, wallet common.Address) (*big.Int, error) {
	var result hexutil.Big
	if err := client.callRPC(ctx, false, &result, "eth_getBalance", wallet, "latest"); err != nil {
		return nil, err
	}

	return (*big.Int)(&result), nil
}

// GetCodeAt returns the contract code deployed at the given address at the latest
// block. An empty result means no contract is deployed there. Used by deployment
// helpers to detect whether a (CREATE2) contract already exists before deploying.
func (client *Client) GetCodeAt(ctx context.Context, address common.Address) ([]byte, error) {
	var result hexutil.Bytes
	if err := client.callRPC(ctx, false, &result, "eth_getCode", address, "latest"); err != nil {
		return nil, err
	}

	return result, nil
}

// EstimateGas runs eth_estimateGas against the chain for the provided call
// message and returns the suggested gas limit. Used by deploy paths that need
// an accurate pre-send estimate under EIP-8037 where static guesses tend to
// under-allocate.
func (client *Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	var result hexutil.Uint64
	if err := client.callRPC(ctx, false, &result, "eth_estimateGas", toCallArg(call)); err != nil {
		return 0, err
	}

	return uint64(result), nil
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

	var gasCapResult hexutil.Big
	if err := client.callRPC(ctx, false, &gasCapResult, "eth_gasPrice"); err != nil {
		return nil, nil, err
	}

	gasCap := (*big.Int)(&gasCapResult)

	var tipCapResult hexutil.Big

	tipCap := big.NewInt(2000000000)
	if err := client.callRPC(ctx, false, &tipCapResult, "eth_maxPriorityFeePerGas"); err == nil {
		tipCap = (*big.Int)(&tipCapResult)
	}

	client.lastGasSuggestion = time.Now()
	client.lastGasCap = gasCap
	client.lastTipCap = tipCap

	return gasCap, tipCap, nil
}

// SendTransaction submits a transaction to the network in its network encoding,
// which for blob transactions includes the sidecar.
func (client *Client) SendTransaction(ctx context.Context, tx *txtypes.Transaction) error {
	client.logger.Tracef("submitted transaction %v", tx.Hash().String())

	encoded, err := tx.MarshalNetwork()
	if err != nil {
		return fmt.Errorf("failed encoding transaction: %w", err)
	}

	return client.SendRawTransaction(ctx, encoded)
}

// SendRawTransaction submits raw transaction bytes to the network using the
// eth_sendRawTransaction RPC call.
func (client *Client) SendRawTransaction(ctx context.Context, tx []byte) error {
	return client.callRPC(ctx, true, nil, "eth_sendRawTransaction", hexutil.Encode(tx))
}

// GetTransactionReceipt retrieves the receipt for a given transaction hash.
// Returns a nil receipt without error when the transaction is not yet mined.
func (client *Client) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*txtypes.Receipt, error) {
	client.logger.Tracef("get receipt: 0x%x", txHash.Bytes())

	var receipt *txtypes.Receipt
	if err := client.callRPC(ctx, false, &receipt, "eth_getTransactionReceipt", txHash); err != nil {
		return nil, err
	}

	if receipt == nil {
		return nil, ethereum.NotFound
	}

	return receipt, nil
}

// GetBlockByNumber retrieves a full block including its transaction objects.
func (client *Client) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*txtypes.Block, error) {
	var block *txtypes.Block
	if err := client.callRPC(ctx, false, &block, "eth_getBlockByNumber", rpc.BlockNumber(blockNumber), true); err != nil {
		return nil, err
	}

	if block == nil {
		return nil, ethereum.NotFound
	}

	return block, nil
}

// GetHeaderByNumber retrieves a block header. A nil blockNumber requests the latest.
func (client *Client) GetHeaderByNumber(ctx context.Context, blockNumber *big.Int) (*txtypes.Header, error) {
	var header *txtypes.Header
	if err := client.callRPC(ctx, false, &header, "eth_getBlockByNumber", toBlockNumberArg(blockNumber), false); err != nil {
		return nil, err
	}

	if header == nil {
		return nil, ethereum.NotFound
	}

	return header, nil
}

// GetBlockReceipts retrieves all transaction receipts of a block.
func (client *Client) GetBlockReceipts(ctx context.Context, blockHash common.Hash) ([]*txtypes.Receipt, error) {
	var receipts []*txtypes.Receipt
	if err := client.callRPC(ctx, false, &receipts, "eth_getBlockReceipts", rpc.BlockNumberOrHash{BlockHash: &blockHash}); err != nil {
		return nil, err
	}

	if receipts == nil {
		return nil, ethereum.NotFound
	}

	return receipts, nil
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

	var result hexutil.Uint64
	if err := client.callRPC(ctx, false, &result, "eth_blockNumber"); err != nil {
		return client.blockHeight, err
	}

	blockHeight := uint64(result)
	if blockHeight > client.blockHeight {
		client.blockHeight = blockHeight
		client.blockHeightTime = time.Now()
	}

	return client.blockHeight, nil
}

// toBlockNumberArg converts a block number to its JSON-RPC representation, mapping
// nil to "latest".
func toBlockNumberArg(blockNumber *big.Int) any {
	if blockNumber == nil {
		return "latest"
	}

	return hexutil.EncodeBig(blockNumber)
}

// toCallArg converts a call message to the JSON-RPC argument object.
func toCallArg(msg ethereum.CallMsg) any {
	arg := map[string]any{
		"from": msg.From,
		"to":   msg.To,
	}

	if len(msg.Data) > 0 {
		arg["input"] = hexutil.Bytes(msg.Data)
	}

	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}

	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}

	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}

	if msg.GasFeeCap != nil {
		arg["maxFeePerGas"] = (*hexutil.Big)(msg.GasFeeCap)
	}

	if msg.GasTipCap != nil {
		arg["maxPriorityFeePerGas"] = (*hexutil.Big)(msg.GasTipCap)
	}

	if msg.AccessList != nil {
		arg["accessList"] = msg.AccessList
	}

	if msg.BlobGasFeeCap != nil {
		arg["maxFeePerBlobGas"] = (*hexutil.Big)(msg.BlobGasFeeCap)
	}

	if msg.BlobHashes != nil {
		arg["blobVersionedHashes"] = msg.BlobHashes
	}

	return arg
}

// GetBlockNumber returns the current block number, bypassing the cache used by
// GetBlockHeight. Head tracking needs an unsmoothed value.
func (client *Client) GetBlockNumber(ctx context.Context) (uint64, error) {
	var result hexutil.Uint64
	if err := client.callRPC(ctx, false, &result, "eth_blockNumber"); err != nil {
		return 0, err
	}

	return uint64(result), nil
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
	client.incrementRequestStats(false, err != nil)
	if err != nil {
		return client.clientVersion, err
	}

	client.clientVersion = result
	client.clientVersionTime = time.Now()

	return result, nil
}

// GetRequestStats returns the current request statistics for this client.
// Returns total requests, transaction requests, and RPC failures.
func (client *Client) GetRequestStats() (total, txRequests, rpcFailures uint64) {
	client.statsMutex.Lock()
	defer client.statsMutex.Unlock()
	return client.totalRequests, client.totalTxRequests, client.totalRpcFailures
}

// incrementRequestStats updates the request statistics based on the request type and outcome.
func (client *Client) incrementRequestStats(isTxRequest bool, failed bool) {
	client.statsMutex.Lock()
	defer client.statsMutex.Unlock()

	client.totalRequests++
	if isTxRequest {
		client.totalTxRequests++
	}
	if failed {
		client.totalRpcFailures++
	}
}
