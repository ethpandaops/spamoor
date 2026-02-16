package spamoor

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ClientSelectionMode defines how clients are selected from the pool.
type ClientSelectionMode uint8

var (
	// SelectClientByIndex selects a client by index (modulo pool size).
	SelectClientByIndex ClientSelectionMode = 0
	// SelectClientRandom selects a random client from the pool.
	SelectClientRandom ClientSelectionMode = 1
	// SelectClientRoundRobin selects clients in round-robin fashion.
	SelectClientRoundRobin ClientSelectionMode = 2
)

// ClientSelectionOption represents an option that can be passed to GetClient.
type ClientSelectionOption interface {
	apply(*clientSelectionOptions)
}

// clientSelectionOptions holds the parsed options for GetClient.
type clientSelectionOptions struct {
	selectionMode ClientSelectionMode
	index         int
	group         string
	excludeTypes  map[ClientType]bool
}

// Default client options - round-robin selection, default group, no exclusions.
func defaultClientSelectionOptions() *clientSelectionOptions {
	return &clientSelectionOptions{
		selectionMode: SelectClientRoundRobin,
		index:         0,
		group:         "",
		excludeTypes:  make(map[ClientType]bool),
	}
}

// selectionModeOption implements ClientOption for selection mode.
type selectionModeOption struct {
	mode  ClientSelectionMode
	index int
}

func (o selectionModeOption) apply(opts *clientSelectionOptions) {
	opts.selectionMode = o.mode
	opts.index = o.index
}

// WithClientSelectionMode sets the client selection mode.
// For SelectClientByIndex mode, provide the index as the second parameter.
func WithClientSelectionMode(mode ClientSelectionMode, args ...int) ClientSelectionOption {
	index := 0
	if len(args) > 0 {
		index = args[0]
	}
	return selectionModeOption{mode: mode, index: index}
}

// groupOption implements ClientOption for group filtering.
type groupOption struct {
	group string
}

func (o groupOption) apply(opts *clientSelectionOptions) {
	opts.group = o.group
}

// WithClientGroup sets the client group filter.
// Use "" for default group, "*" for any group, or specify a group name.
func WithClientGroup(group string) ClientSelectionOption {
	return groupOption{group: group}
}

// excludeTypeOption implements ClientOption for excluding client types.
type excludeTypeOption struct {
	clientType ClientType
}

func (o excludeTypeOption) apply(opts *clientSelectionOptions) {
	opts.excludeTypes[o.clientType] = true
}

// WithoutBuilder excludes builder clients from selection.
func WithoutBuilder() ClientSelectionOption {
	return excludeTypeOption{clientType: ClientTypeBuilder}
}

// WithoutClientType excludes clients of the specified type from selection.
func WithoutClientType(clientType ClientType) ClientSelectionOption {
	return excludeTypeOption{clientType: clientType}
}

// ClientPool manages a pool of Ethereum RPC clients with health monitoring and selection strategies.
// It automatically monitors client health by checking block heights and maintains a list of "good" clients
// that are within 2 blocks of the highest observed block height.
type ClientPool struct {
	ctx            context.Context
	logger         logrus.FieldLogger
	allClients     []*Client
	goodClients    []*Client
	chainId        *big.Int
	selectionMutex sync.Mutex
	rrClientIdx    int
}

// NewClientPool creates a new ClientPool with the specified RPC hosts and logger.
// The pool must be initialized with InitClients() and prepared with PrepareClients() before use.
func NewClientPool(ctx context.Context, logger logrus.FieldLogger) *ClientPool {
	return &ClientPool{
		ctx:    ctx,
		logger: logger,
	}
}

// InitClients initializes the clients in the pool.
func (pool *ClientPool) InitClients(clientOptions []*ClientOptions) error {
	pool.allClients = make([]*Client, 0)

	var chainId *big.Int
	for _, clientOption := range clientOptions {
		client, err := NewClient(clientOption)
		if err != nil {
			pool.logger.Errorf("failed creating client for '%v': %v", clientOption.RpcHost, err.Error())
			continue
		}

		pool.allClients = append(pool.allClients, client)

		if chainId == nil {
			client.Timeout = 10 * time.Second
			cliChainId, err := client.GetChainId(pool.ctx)
			if err != nil {
				pool.logger.Errorf("failed getting chainid from '%v': %v", client.GetRPCHost(), err.Error())
				continue
			}
			chainId = cliChainId
		}
		client.Timeout = 30 * time.Second
	}

	if len(pool.allClients) == 0 {
		return fmt.Errorf("no rpc hosts provided")
	}

	if chainId == nil {
		return fmt.Errorf("no useable clients")
	}

	pool.chainId = chainId

	pool.logger.Infof("initialized client pool with %v clients (chain id: %v)", len(pool.allClients), pool.chainId)

	return nil
}

// PrepareClients initializes all clients in the pool and starts health monitoring.
// It creates Client instances for each RPC host, determines the chain ID,
// and begins periodic health checks. Returns an error if no usable clients are found.
func (pool *ClientPool) PrepareClients() error {
	err := pool.watchClientStatus()
	if err != nil {
		return err
	}
	// watch client status
	go pool.watchClientStatusLoop()

	return nil
}

// watchClientStatusLoop continuously monitors client health in the background.
// It periodically calls watchClientStatus() to check all clients and update the good clients list.
// Runs every 2 minutes normally, but reduces to 10 seconds on errors. Exits when context is cancelled.
func (pool *ClientPool) watchClientStatusLoop() {
	sleepTime := 2 * time.Minute
	for {
		select {
		case <-pool.ctx.Done():
			return
		case <-time.After(sleepTime):
		}

		err := pool.watchClientStatus()
		if err != nil {
			pool.logger.Warnf("could not check client status: %v", err)
			sleepTime = 10 * time.Second
		} else {
			sleepTime = 2 * time.Minute
		}
	}
}

// watchClientStatus checks the health of all clients by querying their current block height.
// It runs concurrent health checks and updates the goodClients list with clients that are
// within 2 blocks of the highest observed block height. Builder clients are assumed to always
// be online and skip the eth_blockNumber check. Logs the results of the health check.
func (pool *ClientPool) watchClientStatus() error {
	wg := &sync.WaitGroup{}
	mtx := sync.Mutex{}
	clientHeads := map[int]uint64{}
	highestHead := uint64(0)

	// First, check all non-builder clients to determine the highest block height
	for idx, client := range pool.allClients {
		if client.IsBuilder() {
			continue // Skip builders in the initial pass
		}

		wg.Add(1)
		go func(idx int, client *Client) {
			defer wg.Done()

			var blockHeight uint64
			var err error
			hasBlockHeight := false

			if client.externalClient != nil {
				blockHeight, err = client.externalClient.GetBlockHeight(pool.ctx)
				if err != nil {
					pool.logger.Warnf("external client check failed: %v", err)
				} else {
					hasBlockHeight = true
				}
			}

			if !hasBlockHeight {
				blockHeight, err = client.GetBlockHeight(pool.ctx)
			}
			if err != nil {
				pool.logger.Warnf("client check failed: %v", err)
			} else {
				mtx.Lock()
				clientHeads[idx] = blockHeight
				if blockHeight > highestHead {
					highestHead = blockHeight
				}
				mtx.Unlock()
			}
		}(idx, client)
	}
	wg.Wait()

	// Now mark all builder clients as healthy with the highest block height
	for idx, client := range pool.allClients {
		if client.IsBuilder() {
			clientHeads[idx] = highestHead // Assume builders are at the latest block
		}
	}

	goodClients := make([]*Client, 0)
	goodHead := highestHead
	if goodHead > 2 {
		goodHead -= 2
	}
	for idx, client := range pool.allClients {
		if client.IsBuilder() || clientHeads[idx] >= goodHead {
			goodClients = append(goodClients, client)
		}
	}
	pool.goodClients = goodClients
	pool.logger.Infof("client check completed (%v good clients, %v bad clients)", len(goodClients), len(pool.allClients)-len(goodClients))

	return nil
}

// GetClient returns a client from the pool based on the specified options.
// By default, it uses round-robin selection with the default group and no type exclusions.
//
// Available options:
//   - WithSelectionMode(mode, index...): Set selection mode (ByIndex, Random, RoundRobin), index optional for ByIndex
//   - WithClientGroup(group): Set group filter ("" for default, "*" for any, or group name)
//   - WithoutBuilder(): Exclude builder clients
//   - WithoutClientType(type): Exclude clients of specified type
//
// Examples:
//   - pool.GetClient() // Round-robin from default group
//   - pool.GetClient(WithSelectionMode(SelectClientRandom)) // Random from default group
//   - pool.GetClient(WithClientGroup("builders")) // Round-robin from builders group
//   - pool.GetClient(WithoutBuilder()) // Round-robin excluding builders
//   - pool.GetClient(WithSelectionMode(SelectClientByIndex, 2)) // Select by index 2
//
// Returns nil if no suitable clients are available.
func (pool *ClientPool) GetClient(options ...ClientSelectionOption) *Client {
	pool.selectionMutex.Lock()
	defer pool.selectionMutex.Unlock()

	if len(pool.goodClients) == 0 {
		return nil
	}

	// Parse options
	opts := defaultClientSelectionOptions()
	for _, option := range options {
		option.apply(opts)
	}

	// Build list of candidate clients
	clientCandidates := make([]*Client, 0)

	for _, client := range pool.goodClients {
		// Check if client is enabled
		if !client.IsEnabled() {
			continue
		}

		// Check group filter
		if opts.group == "" {
			// Empty group means default group
			if !client.HasGroup("default") {
				continue
			}
		} else if opts.group != "*" {
			// Specific group name (wildcard "*" accepts any group)
			if !client.HasGroup(opts.group) {
				continue
			}
		}

		// Check type exclusions
		if opts.excludeTypes[client.GetClientType()] {
			continue
		}

		clientCandidates = append(clientCandidates, client)
	}

	if len(clientCandidates) == 0 {
		return nil
	}

	// Select client based on mode
	var selectedIndex int
	switch opts.selectionMode {
	case SelectClientByIndex:
		selectedIndex = opts.index % len(clientCandidates)
	case SelectClientRandom:
		selectedIndex = rand.Intn(len(clientCandidates))
	case SelectClientRoundRobin:
		selectedIndex = pool.rrClientIdx
		pool.rrClientIdx++
		if pool.rrClientIdx >= len(clientCandidates) {
			pool.rrClientIdx = 0
		}
	default:
		// Fallback to round-robin
		selectedIndex = pool.rrClientIdx
		pool.rrClientIdx++
		if pool.rrClientIdx >= len(clientCandidates) {
			pool.rrClientIdx = 0
		}
	}

	return clientCandidates[selectedIndex]
}

// GetAllClients returns a copy of all clients in the pool, regardless of their health status.
func (pool *ClientPool) GetAllClients() []*Client {
	clients := make([]*Client, len(pool.allClients))
	copy(clients, pool.allClients)
	return clients
}

// GetAllGoodClients returns a copy of all clients currently considered healthy
// (within 2 blocks of the highest observed block height).
func (pool *ClientPool) GetAllGoodClients() []*Client {
	clients := make([]*Client, len(pool.goodClients))
	copy(clients, pool.goodClients)
	return clients
}

// GetChainId returns the chain ID of the clients in the pool.
func (pool *ClientPool) GetChainId() *big.Int {
	return pool.chainId
}
