package spamoor

import (
	"context"
	"errors"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoortypes"
)

// ClientPool manages a pool of Ethereum RPC clients with health monitoring and selection strategies.
// It automatically monitors client health by checking block heights and maintains a list of "good" clients
// that are within 2 blocks of the highest observed block height.
type ClientPool struct {
	ctx            context.Context
	clientFactory  func(string) (spamoortypes.Client, error)
	rpcHosts       []string
	logger         logrus.FieldLogger
	allClients     []spamoortypes.Client
	goodClients    []spamoortypes.Client
	chainId        *big.Int
	selectionMutex sync.Mutex
	rrClientIdx    int
}

// NewClientPool creates a new ClientPool with the specified RPC hosts and logger.
// The pool must be initialized with PrepareClients() before use.
func NewClientPool(ctx context.Context, rpcHosts []string, logger logrus.FieldLogger) *ClientPool {
	return &ClientPool{
		ctx:           ctx,
		rpcHosts:      rpcHosts,
		logger:        logger,
		allClients:    make([]spamoortypes.Client, 0),
		goodClients:   make([]spamoortypes.Client, 0),
		clientFactory: NewClient,
	}
}

// PrepareClients initializes all clients in the pool and starts health monitoring.
// It creates Client instances for each RPC host, determines the chain ID,
// and begins periodic health checks. Returns an error if no usable clients are found.
func (pool *ClientPool) PrepareClients() error {
	if len(pool.rpcHosts) == 0 {
		return errors.New("no rpc hosts provided")
	}

	for _, rpchost := range pool.rpcHosts {
		client, err := pool.clientFactory(rpchost)
		if err != nil {
			pool.logger.Warnf("failed to create client for '%v': %v", rpchost, err)
			continue
		}
		pool.allClients = append(pool.allClients, client)
	}

	if len(pool.allClients) == 0 {
		return errors.New("no useable clients")
	}

	// Get chain ID from first client
	chainId, err := pool.allClients[0].GetChainId(pool.ctx)
	if err != nil {
		return err
	}
	pool.chainId = chainId

	pool.logger.Infof("initialized client pool with %v clients (chain id: %v)", len(pool.allClients), pool.chainId)

	// Start monitoring client status
	pool.watchClientStatus()
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
// within 2 blocks of the highest observed block height. Logs the results of the health check.
func (pool *ClientPool) watchClientStatus() error {
	wg := &sync.WaitGroup{}
	mtx := sync.Mutex{}
	clientHeads := map[int]uint64{}
	highestHead := uint64(0)

	for idx, client := range pool.allClients {
		wg.Add(1)
		go func(idx int, client spamoortypes.Client) {
			defer wg.Done()

			blockHeight, err := client.GetBlockHeight(pool.ctx)
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

	goodClients := make([]spamoortypes.Client, 0)
	goodHead := highestHead
	if goodHead > 2 {
		goodHead -= 2
	}
	for idx, client := range pool.allClients {
		if clientHeads[idx] >= goodHead {
			goodClients = append(goodClients, client)
		}
	}
	pool.goodClients = goodClients
	pool.logger.Infof("client check completed (%v good clients, %v bad clients)", len(goodClients), len(pool.allClients)-len(goodClients))

	return nil
}

// GetClient returns a client from the pool based on the specified selection mode.
// Parameters:
//   - mode: how to select the client (by index, random, or round-robin)
//   - input: used as index when mode is SelectClientByIndex
//   - group: client group filter ("" for default, "*" for any, or specific group name)
//
// Returns nil if no suitable clients are available.
func (pool *ClientPool) GetClient(mode spamoortypes.ClientSelectionMode, input int, group string) spamoortypes.Client {
	pool.selectionMutex.Lock()
	defer pool.selectionMutex.Unlock()

	if len(pool.goodClients) == 0 {
		return nil
	}

	clientCandidates := make([]spamoortypes.Client, 0)

	if group == "" {
		for _, client := range pool.goodClients {
			if client.IsEnabled() && client.GetClientGroup() == "default" {
				clientCandidates = append(clientCandidates, client)
			}
		}
	} else if group == "*" {
		for _, client := range pool.goodClients {
			if client.IsEnabled() {
				clientCandidates = append(clientCandidates, client)
			}
		}
	}

	if len(clientCandidates) == 0 {
		for _, client := range pool.goodClients {
			if client.IsEnabled() && (group == "" || client.GetClientGroup() == group) {
				clientCandidates = append(clientCandidates, client)
			}
		}
	}

	if len(clientCandidates) == 0 {
		return nil
	}

	switch mode {
	case spamoortypes.SelectClientByIndex:
		input = input % len(clientCandidates)
	case spamoortypes.SelectClientRandom:
		input = rand.Intn(len(clientCandidates))
	case spamoortypes.SelectClientRoundRobin:
		input = pool.rrClientIdx
		pool.rrClientIdx++
		if pool.rrClientIdx >= len(clientCandidates) {
			pool.rrClientIdx = 0
		}
	}
	return clientCandidates[input]
}

// GetAllClients returns a copy of all clients in the pool, regardless of their health status.
func (pool *ClientPool) GetAllClients() []spamoortypes.Client {
	clients := make([]spamoortypes.Client, len(pool.allClients))
	copy(clients, pool.allClients)
	return clients
}

// GetAllGoodClients returns a copy of all clients currently considered healthy
// (within 2 blocks of the highest observed block height).
func (pool *ClientPool) GetAllGoodClients() []spamoortypes.Client {
	clients := make([]spamoortypes.Client, len(pool.goodClients))
	copy(clients, pool.goodClients)
	return clients
}
