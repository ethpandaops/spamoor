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

type testFixture struct {
	t          *testing.T
	cmd        *exec.Cmd
	client     *ethclient.Client
	ctx        context.Context
	logger     *logrus.Entry
	clientPool *spamoor.ClientPool
	rootWallet *txbuilder.Wallet
	txpool     *txbuilder.TxPool
	walletPool *spamoor.WalletPool
	scenario   *Scenario
}

func setupTestFixture(t *testing.T) *testFixture {
	// Start Anvil server with legacy transactions
	cmd := exec.Command("anvil", "--hardfork", "pectra")
	err := cmd.Start()
	require.NoError(t, err)

	// Wait for Anvil to start
	time.Sleep(2 * time.Second)

	// Connect to Anvil
	client, err := ethclient.Dial("http://localhost:8545")
	require.NoError(t, err)

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
	scenario := newScenario(logger).(*Scenario)

	return &testFixture{
		t:          t,
		cmd:        cmd,
		client:     client,
		ctx:        ctx,
		logger:     logger,
		clientPool: clientPool,
		rootWallet: rootWallet,
		txpool:     txpool,
		walletPool: walletPool,
		scenario:   scenario,
	}
}

func (f *testFixture) teardown() {
	if f.cmd != nil && f.cmd.Process != nil {
		f.cmd.Process.Kill()
	}
	if f.client != nil {
		f.client.Close()
	}
}

func (f *testFixture) verifyContractDeployment(runCtx context.Context) {
	// Verify contract deployment
	block, err := f.client.BlockByNumber(runCtx, nil)
	require.NoError(f.t, err)
	assert.Greater(f.t, len(block.Transactions()), 0)

	// Check the first transaction
	tx, _, err := f.client.TransactionByHash(runCtx, block.Transactions()[0].Hash())
	require.NoError(f.t, err)
	receipt, err := f.client.TransactionReceipt(runCtx, tx.Hash())
	require.NoError(f.t, err)
	assert.Equal(f.t, uint64(1), receipt.Status)

	// Verify contract was created
	code, err := f.client.CodeAt(runCtx, receipt.ContractAddress, nil)
	require.NoError(f.t, err)
	assert.Greater(f.t, len(code), 0)
}

func TestContractDeployWithContractsPerBlock(t *testing.T) {
	fixture := setupTestFixture(t)
	defer fixture.teardown()

	// Initialize with contracts_per_block configuration
	config := `
max_pending: 0
max_wallets: 1
rebroadcast: 2
base_fee: 20
tip_fee: 2
gas_per_block: 0
client_group: default
contracts_per_block: 6
max_transactions: 1
`
	require.NoError(t, fixture.scenario.Init(fixture.walletPool, config))

	// Run scenario
	runCtx, cancel := context.WithTimeout(fixture.ctx, 30*time.Second)
	defer cancel()

	err := fixture.scenario.Run(runCtx)
	require.NoError(t, err)

	fixture.verifyContractDeployment(runCtx)
}

func TestContractDeployWithGasPerBlock(t *testing.T) {
	fixture := setupTestFixture(t)
	defer fixture.teardown()

	// Initialize with gas_per_block configuration
	config := `
max_pending: 0
max_wallets: 1
rebroadcast: 2
base_fee: 20
tip_fee: 2
gas_per_block: 30000000
client_group: default
contracts_per_block: 0
max_transactions: 1
`
	require.NoError(t, fixture.scenario.Init(fixture.walletPool, config))

	// Run scenario
	runCtx, cancel := context.WithTimeout(fixture.ctx, 30*time.Second)
	defer cancel()

	err := fixture.scenario.Run(runCtx)
	require.NoError(t, err)

	fixture.verifyContractDeployment(runCtx)
}

func TestContractDeployWithInvalidConfig(t *testing.T) {
	fixture := setupTestFixture(t)
	defer fixture.teardown()

	// Initialize with invalid configuration
	config := `
max_pending: 0
max_wallets: 1
rebroadcast: 1
base_fee: 20
tip_fee: 2
gas_per_block: 0
client_group: default
contracts_per_block: 0
max_transactions: 1
`
	// We expect Init to fail with this configuration
	err := fixture.scenario.Init(fixture.walletPool, config)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "neither gas per block limit nor contracts per block set")
}
