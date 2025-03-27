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
	wg := &sync.WaitGroup{}
	mtx := sync.Mutex{}

	var chainId *big.Int
	for _, rpcHost := range pool.rpcHosts {
		wg.Add(1)

		go func(rpcHost string) {
			defer wg.Done()

			client, err := txbuilder.NewClient(rpcHost)
			if err != nil {
				pool.logger.Errorf("failed creating client for '%v': %v", client.GetRPCHost(), err.Error())
				return
			}
			client.Timeout = 10 * time.Second
			cliChainId, err := client.GetChainId(pool.ctx)
			if err != nil {
				pool.logger.Errorf("failed getting chainid from '%v': %v", client.GetRPCHost(), err.Error())
				return
			}
			if chainId == nil {
				chainId = cliChainId
			} else if cliChainId.Cmp(chainId) != 0 {
				pool.logger.Errorf("chainid missmatch from %v (chain ids: %v, %v)", client.GetRPCHost(), cliChainId, chainId)
				return
			}
			client.Timeout = 30 * time.Second
			mtx.Lock()
			pool.allClients = append(pool.allClients, client)
			mtx.Unlock()
		}(rpcHost)
	}

	wg.Wait()
	pool.chainId = chainId
	if len(pool.allClients) == 0 {
		return fmt.Errorf("no useable clients")
	}

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

func (pool *ClientPool) GetClient(mode ClientSelectionMode, input int) *txbuilder.Client {
	pool.selectionMutex.Lock()
	defer pool.selectionMutex.Unlock()

	if len(pool.goodClients) == 0 {
		return nil
	}

	switch mode {
	case SelectClientByIndex:
		input = input % len(pool.goodClients)
	case SelectClientRandom:
		input = rand.Intn(len(pool.goodClients))
	case SelectClientRoundRobin:
		input = pool.rrClientIdx
		pool.rrClientIdx++
		if pool.rrClientIdx >= len(pool.goodClients) {
			pool.rrClientIdx = 0
		}
	}
	return pool.goodClients[input]
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
