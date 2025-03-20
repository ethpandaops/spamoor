package tester

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethpandaops/spamoor/txbuilder"
)

func (tester *Tester) PrepareClients() error {
	tester.allClients = make([]*txbuilder.Client, 0)
	wg := &sync.WaitGroup{}
	mtx := sync.Mutex{}

	var chainId *big.Int
	for _, rpcHost := range tester.config.RpcHosts {
		wg.Add(1)

		go func(rpcHost string) {
			defer wg.Done()

			client, err := txbuilder.NewClient(rpcHost)
			if err != nil {
				tester.logger.Errorf("failed creating client for '%v': %v", client.GetRPCHost(), err.Error())
				return
			}
			client.Timeout = 10 * time.Second
			cliChainId, err := client.GetChainId(tester.ctx)
			if err != nil {
				tester.logger.Errorf("failed getting chainid from '%v': %v", client.GetRPCHost(), err.Error())
				return
			}
			if chainId == nil {
				chainId = cliChainId
			} else if cliChainId.Cmp(chainId) != 0 {
				tester.logger.Errorf("chainid missmatch from %v (chain ids: %v, %v)", client.GetRPCHost(), cliChainId, chainId)
				return
			}
			client.Timeout = 30 * time.Second
			mtx.Lock()
			tester.allClients = append(tester.allClients, client)
			mtx.Unlock()
		}(rpcHost)
	}

	wg.Wait()
	tester.chainId = chainId
	if len(tester.allClients) == 0 {
		return fmt.Errorf("no useable clients")
	}
	return nil
}

func (tester *Tester) watchClientStatus() error {
	wg := &sync.WaitGroup{}
	mtx := sync.Mutex{}
	clientHeads := map[int]uint64{}
	highestHead := uint64(0)

	for idx, client := range tester.allClients {
		wg.Add(1)
		go func(idx int, client *txbuilder.Client) {
			defer wg.Done()

			blockHeight, err := client.GetBlockHeight(tester.ctx)
			if err != nil {
				tester.logger.Warnf("client check failed: %v", err)
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
	for idx, client := range tester.allClients {
		if clientHeads[idx] >= goodHead {
			goodClients = append(goodClients, client)
		}
	}
	tester.goodClients = goodClients
	tester.logger.Warnf("client check completed (%v good clients, %v bad clients)", len(goodClients), len(tester.allClients)-len(goodClients))

	return nil
}
