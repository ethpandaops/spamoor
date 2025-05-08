package contractdeploy

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

func TestContractDeploy(t *testing.T) {
	// Start Anvil server with legacy transactions
	cmd := exec.Command("anvil", "--hardfork", "pectra")
	err := cmd.Start()
	require.NoError(t, err)
	defer cmd.Process.Kill()

	// Wait for Anvil to start
	time.Sleep(2 * time.Second)

	// Connect to Anvil
	client, err := ethclient.Dial("http://localhost:8545")
	require.NoError(t, err)
	defer client.Close()

	// Create context
	ctx := context.Background()
	logger := logrus.New().WithField("test", "contract-deploy")

	// Create client pool
	clientPool := spamoor.NewClientPool(ctx, []string{"http://localhost:8545"}, logger)
	err = clientPool.PrepareClients()
	require.NoError(t, err)

	// Create root wallet using a pre-funded Anvil account
	rootWallet, err := spamoor.InitRootWallet(ctx, "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", clientPool.GetClient(spamoor.SelectClientByIndex, 0, ""), logger)
	require.NoError(t, err)

	// Create tx pool
	txpool := txbuilder.NewTxPool(&txbuilder.TxPoolOptions{
		Context: ctx,
		GetClientFn: func(index int, random bool) *txbuilder.Client {
			mode := spamoor.SelectClientByIndex
			if random {
				mode = spamoor.SelectClientRandom
			}
			return clientPool.GetClient(mode, index, "")
		},
		GetClientCountFn: func() int {
			return len(clientPool.GetAllClients())
		},
	})

	// Create wallet pool
	walletPool := spamoor.NewWalletPool(ctx, logger, rootWallet, clientPool, txpool)
	walletPool.SetWalletCount(1)
	walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(1)))  // 1 ETH
	walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(1))) // 1 ETH
	walletPool.SetRefillInterval(60)

	// Create scenario instance
	scenario := newScenario(logger)

	// Configure scenario with 6 contracts per tx which is the maximum allowed by the latest fork.
	config := `
max_pending: 0
max_wallets: 1
rebroadcast: 30
base_fee: 20
tip_fee: 2
gas_per_block: 30000000
client_group: default
contracts_per_tx: 6
`
	err = scenario.Init(walletPool, config)
	require.NoError(t, err)

	// Prepare wallets
	err = walletPool.PrepareWallets(true)
	require.NoError(t, err)

	// Create context with timeout
	runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Run scenario
	err = scenario.Run(runCtx)
	require.NoError(t, err)

	// Verify contract deployment
	block, err := client.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	assert.Greater(t, block.Transactions().Len(), 0)

	// Get transaction receipt
	tx := block.Transactions()[0]
	receipt, err := client.TransactionReceipt(ctx, tx.Hash())
	require.NoError(t, err)
	assert.Equal(t, uint64(1), receipt.Status)

	// Verify contract was created
	code, err := client.CodeAt(ctx, receipt.ContractAddress, nil)
	require.NoError(t, err)
	assert.Greater(t, len(code), 0)
}
