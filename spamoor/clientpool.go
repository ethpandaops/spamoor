package spamoor

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/sirupsen/logrus"
)

type ClientSelectionMode uint8

var (
	SelectClientByIndex    ClientSelectionMode = 0
	SelectClientRandom     ClientSelectionMode = 1
	SelectClientRoundRobin ClientSelectionMode = 2
)

type ClientPool struct {
	ctx            context.Context
	rpcHosts       []string
	logger         logrus.FieldLogger
	allClients     []*txbuilder.Client
	goodClients    []*txbuilder.Client
	chainId        *big.Int
	selectionMutex sync.Mutex
	rrClientIdx    int
}

func NewClientPool(ctx context.Context, rpcHosts []string, logger logrus.FieldLogger) *ClientPool {
	return &ClientPool{
		ctx:      ctx,
		rpcHosts: rpcHosts,
		logger:   logger,
	}
}

func (pool *ClientPool) PrepareClients() error {
	pool.allClients = make([]*txbuilder.Client, 0)

	var chainId *big.Int
	for _, rpcHost := range pool.rpcHosts {
		client, err := txbuilder.NewClient(rpcHost)
		if err != nil {
			pool.logger.Errorf("failed creating client for '%v': %v", client.GetRPCHost(), err.Error())
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

	err := pool.watchClientStatus()
	if err != nil {
		return err
	}
	// watch client status
	go pool.watchClientStatusLoop()

	return nil
}

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

func (pool *ClientPool) watchClientStatus() error {
	wg := &sync.WaitGroup{}
	mtx := sync.Mutex{}
	clientHeads := map[int]uint64{}
	highestHead := uint64(0)

	for idx, client := range pool.allClients {
		wg.Add(1)
		go func(idx int, client *txbuilder.Client) {
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

	goodClients := make([]*txbuilder.Client, 0)
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

func (pool *ClientPool) GetClient(mode ClientSelectionMode, input int, group string) *txbuilder.Client {
	pool.selectionMutex.Lock()
	defer pool.selectionMutex.Unlock()

	if len(pool.goodClients) == 0 {
		return nil
	}

	clientCandidates := make([]*txbuilder.Client, 0)

	if group == "" {
		for _, client := range pool.goodClients {
			if client.GetClientGroup() == "default" {
				clientCandidates = append(clientCandidates, client)
			}
		}
	} else if group == "*" {
		clientCandidates = pool.goodClients
	}

	if len(clientCandidates) == 0 {
		for _, client := range pool.goodClients {
			if group == "" || client.GetClientGroup() == group {
				clientCandidates = append(clientCandidates, client)
			}
		}
	}

	if len(clientCandidates) == 0 {
		return nil
	}

	switch mode {
	case SelectClientByIndex:
		input = input % len(clientCandidates)
	case SelectClientRandom:
		input = rand.Intn(len(clientCandidates))
	case SelectClientRoundRobin:
		input = pool.rrClientIdx
		pool.rrClientIdx++
		if pool.rrClientIdx >= len(clientCandidates) {
			pool.rrClientIdx = 0
		}
	}
	return clientCandidates[input]
}

func (pool *ClientPool) GetAllClients() []*txbuilder.Client {
	clients := make([]*txbuilder.Client, len(pool.allClients))
	copy(clients, pool.allClients)
	return clients
}

func (pool *ClientPool) GetAllGoodClients() []*txbuilder.Client {
	clients := make([]*txbuilder.Client, len(pool.goodClients))
	copy(clients, pool.goodClients)
	return clients
}
