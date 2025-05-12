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

	// Create client pool
	ctx := context.Background()
	logger := logrus.New().WithField("test", "contract-deploy")
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

	// Initialize wallet pool
	err = walletPool.PrepareWallets(true)
	require.NoError(t, err)

	// Create scenario instance
	scenario := newScenario(logger)

	// Test case 1: Using contracts_per_block
	t.Run("contracts_per_block", func(t *testing.T) {
		// Initialize with contracts_per_block configuration
		config := `
max_pending: 0
max_wallets: 1
rebroadcast: 30
base_fee: 20
tip_fee: 2
gas_per_block: 0
client_group: default
contracts_per_block: 6
`
		require.NoError(t, scenario.Init(walletPool, config))

		// Run scenario
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := scenario.Run(runCtx)
		require.NoError(t, err)

		// Verify contract deployment
		block, err := client.BlockByNumber(runCtx, nil)
		require.NoError(t, err)
		assert.Greater(t, len(block.Transactions()), 0)

		// Check the first transaction
		tx, _, err := client.TransactionByHash(runCtx, block.Transactions()[0].Hash())
		require.NoError(t, err)
		receipt, err := client.TransactionReceipt(runCtx, tx.Hash())
		require.NoError(t, err)
		assert.Equal(t, uint64(1), receipt.Status)

		// Verify contract was created
		code, err := client.CodeAt(runCtx, receipt.ContractAddress, nil)
		require.NoError(t, err)
		assert.Greater(t, len(code), 0)
	})

	// Test case 2: Using gas_per_block
	t.Run("gas_per_block", func(t *testing.T) {
		// Initialize with gas_per_block configuration
		config := `
max_pending: 0
max_wallets: 1
rebroadcast: 30
base_fee: 20
tip_fee: 2
gas_per_block: 30000000
client_group: default
contracts_per_block: 0
`
		require.NoError(t, scenario.Init(walletPool, config))

		// Run scenario
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := scenario.Run(runCtx)
		require.NoError(t, err)

		// Verify contract deployment
		block, err := client.BlockByNumber(runCtx, nil)
		require.NoError(t, err)
		assert.Greater(t, len(block.Transactions()), 0)

		// Check the first transaction
		tx, _, err := client.TransactionByHash(runCtx, block.Transactions()[0].Hash())
		require.NoError(t, err)
		receipt, err := client.TransactionReceipt(runCtx, tx.Hash())
		require.NoError(t, err)
		assert.Equal(t, uint64(1), receipt.Status)

		// Verify contract was created
		code, err := client.CodeAt(runCtx, receipt.ContractAddress, nil)
		require.NoError(t, err)
		assert.Greater(t, len(code), 0)
	})

	// Test case 3: Invalid configuration
	t.Run("invalid_config", func(t *testing.T) {
		// Initialize with invalid configuration
		config := `
max_pending: 0
max_wallets: 1
rebroadcast: 30
base_fee: 20
tip_fee: 2
gas_per_block: 0
client_group: default
contracts_per_block: 0
`
		require.NoError(t, scenario.Init(walletPool, config))

		// Run scenario
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := scenario.Run(runCtx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "neither gas per block limit nor contracts per block set")
	})
}
