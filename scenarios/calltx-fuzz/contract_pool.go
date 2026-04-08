package calltxfuzz

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// ContractPool manages a deterministic pool of deployed fuzz contracts.
// Contracts are generated from a seed so the pool is reproducible.
// It operates as a ring buffer: new deployments rotate out the oldest contract.
//
// Deployment schedule (by effectiveTxID):
//   - 0 to poolSize-1: initial pool fill (deploy index = effectiveTxID)
//   - poolSize onwards: deploy every deployInterval txs
//     (deploy index = poolSize + (effectiveTxID - poolSize) / deployInterval)
//
// This deterministic schedule allows reconstructing the pool state
// at any txIdOffset without replaying the full history.
type ContractPool struct {
	mu             sync.RWMutex
	contracts      []common.Address // Deployed contract addresses
	bytecodes      [][]byte         // Corresponding runtime bytecodes
	poolSize       uint64
	deployInterval uint64 // Post-init deploy interval (0 = no post-init deploys)
	seed           string
	maxCodeSize    uint64
	minCodeSize    uint64
	gasLimit       uint64
	logger         *logrus.Entry
}

// NewContractPool creates a new deterministic contract pool.
// deployRatio determines how often post-init deployments occur (0 = never).
// For example, deployRatio=0.1 means every 10th tx after the initial fill
// deploys a new contract that rotates out the oldest.
func NewContractPool(
	poolSize uint64,
	seed string,
	maxCodeSize, minCodeSize, gasLimit uint64,
	deployRatio float64,
	logger *logrus.Entry,
) *ContractPool {
	var deployInterval uint64
	if deployRatio > 0 && deployRatio <= 1 {
		deployInterval = uint64(math.Round(1.0 / deployRatio))
		if deployInterval == 0 {
			deployInterval = 1
		}
	}

	return &ContractPool{
		contracts:      make([]common.Address, 0, poolSize),
		bytecodes:      make([][]byte, 0, poolSize),
		poolSize:       poolSize,
		deployInterval: deployInterval,
		seed:           seed,
		maxCodeSize:    maxCodeSize,
		minCodeSize:    minCodeSize,
		gasLimit:       gasLimit,
		logger:         logger,
	}
}

// GetAddresses returns a snapshot of all contract addresses currently in the pool.
func (p *ContractPool) GetAddresses() []common.Address {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]common.Address, len(p.contracts))
	copy(result, p.contracts)

	return result
}

// GetRandomContract returns a random contract address from the pool.
func (p *ContractPool) GetRandomContract(rng interface{ Intn(int) int }) common.Address {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.contracts) == 0 {
		return common.Address{}
	}

	return p.contracts[rng.Intn(len(p.contracts))]
}

// Size returns the current number of contracts in the pool.
func (p *ContractPool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.contracts)
}

// GetDeployInterval returns the post-init deploy interval.
func (p *ContractPool) GetDeployInterval() uint64 {
	return p.deployInterval
}

// ShouldDeploy returns true if the given effectiveTxID should be a deployment
// in the main transaction loop. Initial pool fill (txIDs 0 to poolSize-1) is
// handled by InitPool, so this only returns true for post-init rotations.
func (p *ContractPool) ShouldDeploy(effectiveTxID uint64) bool {
	if effectiveTxID < p.poolSize {
		return false // Handled by InitPool
	}
	if p.deployInterval == 0 {
		return false
	}

	return (effectiveTxID-p.poolSize)%p.deployInterval == 0
}

// DeployIndexForTx returns the contract deploy index for a given effectiveTxID.
// Only meaningful when ShouldDeploy returns true or for initial pool fill.
func (p *ContractPool) DeployIndexForTx(effectiveTxID uint64) uint64 {
	if effectiveTxID < p.poolSize {
		return effectiveTxID
	}
	if p.deployInterval == 0 {
		return p.poolSize - 1
	}

	return p.poolSize + (effectiveTxID-p.poolSize)/p.deployInterval
}

// countDeploysBefore returns the total number of deployments with
// effectiveTxID strictly less than the given value.
func (p *ContractPool) countDeploysBefore(effectiveTxID uint64) uint64 {
	if effectiveTxID == 0 {
		return 0
	}
	if effectiveTxID <= p.poolSize {
		return effectiveTxID
	}

	initialDeploys := p.poolSize
	if p.deployInterval == 0 {
		return initialDeploys
	}

	// Post-init deploys at: poolSize, poolSize+N, poolSize+2N, ...
	// Count those with txID < effectiveTxID.
	postInitDeploys := (effectiveTxID - p.poolSize - 1) / p.deployInterval
	postInitDeploys++ // +1 for the deploy at poolSize itself

	return initialDeploys + postInitDeploys
}

// generateRuntimeBytecode generates runtime bytecode for a deploy index.
// peerAddresses are addresses of other pool contracts that the generated
// bytecode can cross-call. When nil, templates fall back to ADDRESS/CALLER.
func (p *ContractPool) generateRuntimeBytecode(
	deployIdx uint64,
	peerAddresses []common.Address,
) []byte {
	codeSize := p.minCodeSize
	if p.maxCodeSize > p.minCodeSize {
		rng := newQuickRNG(deployIdx, p.seed)
		codeSize = p.minCodeSize + uint64(rng.Intn(int(p.maxCodeSize-p.minCodeSize)))
	}

	gen := NewCallFuzzGenerator(
		deployIdx, p.seed, int(codeSize), p.gasLimit, peerAddresses,
	)

	return gen.Generate()
}

// wrapInInitCode wraps runtime bytecode in init code that deploys it.
// Init code: PUSH2 <size> DUP1 PUSH2 <offset> PUSH0 CODECOPY PUSH0 RETURN <runtime>
func wrapInInitCode(runtime []byte) []byte {
	runtimeLen := len(runtime)

	// Init code: PUSH2 <size> DUP1 PUSH2 <offset> PUSH0 CODECOPY PUSH0 RETURN
	// That's: 3 + 1 + 3 + 1 + 1 + 1 + 1 = 11 bytes of init code
	const initCodeLen = 11

	initCode := make([]byte, 0, initCodeLen+runtimeLen)

	// PUSH2 <runtimeLen>
	initCode = append(initCode, 0x61, byte(runtimeLen>>8), byte(runtimeLen))
	// DUP1
	initCode = append(initCode, 0x80)
	// PUSH2 <initCodeLen> (offset where runtime begins)
	initCode = append(initCode, 0x61, byte(initCodeLen>>8), byte(initCodeLen))
	// PUSH0 (destination offset in memory)
	initCode = append(initCode, 0x5f)
	// CODECOPY
	initCode = append(initCode, 0x39)
	// PUSH0 (memory offset for RETURN)
	initCode = append(initCode, 0x5f)
	// RETURN
	initCode = append(initCode, 0xf3)

	// Append runtime bytecode
	initCode = append(initCode, runtime...)

	return initCode
}

// InitPool deploys the initial pool of contracts using multi-wallet batch
// submission. Distributes deployments across child wallets to avoid
// per-wallet pending limits that would throttle a single deployer.
func (p *ContractPool) InitPool(
	ctx context.Context,
	walletPool *spamoor.WalletPool,
	txIdOffset uint64,
	baseFee, tipFee float64,
	baseFeeWei, tipFeeWei string,
	deployGasLimit uint64,
	clientGroup string,
) error {
	// Determine which deploy indices should be in the pool
	totalDeploys := p.countDeploysBefore(txIdOffset)

	// If starting within the initial pool fill phase, deploy the full initial pool
	if txIdOffset <= p.poolSize {
		totalDeploys = p.poolSize
	}

	if totalDeploys == 0 {
		return nil
	}

	startIdx := uint64(0)
	if totalDeploys > p.poolSize {
		startIdx = totalDeploys - p.poolSize
	}

	contractCount := totalDeploys - startIdx

	client := walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(clientGroup),
	)
	if client == nil {
		return fmt.Errorf("no client available for deployment")
	}

	// Collect deploy wallets: spread across child wallets to avoid
	// per-wallet pending limits throttling the deployment.
	const maxDeployWallets = 20
	walletCount := int(contractCount)
	if walletCount > maxDeployWallets {
		walletCount = maxDeployWallets
	}

	wallets := make([]*spamoor.Wallet, walletCount)
	for i := range walletCount {
		w := walletPool.GetWallet(spamoor.SelectWalletByIndex, i)
		if w == nil {
			return fmt.Errorf("deployer wallet %d not found", i)
		}
		if err := w.ResetNoncesIfNeeded(ctx, client); err != nil {
			return fmt.Errorf("deployer wallet %d nonce reset failed: %w", i, err)
		}
		wallets[i] = w
	}

	p.logger.Infof("deploying initial contract pool (%d contracts across %d wallets, deploy indices %d-%d)",
		contractCount, walletCount, startIdx, totalDeploys-1)

	feeCapWei, tipCapWei := spamoor.ResolveFees(baseFee, tipFee, baseFeeWei, tipFeeWei)
	feeCap, tipCap, err := walletPool.GetTxPool().GetSuggestedFees(client, feeCapWei, tipCapWei)
	if err != nil {
		return fmt.Errorf("fee suggestion failed: %w", err)
	}

	// Predict deployment addresses per wallet so bytecodes can cross-reference.
	// Assign contracts round-robin across wallets.
	type deployEntry struct {
		walletIdx int
		nonce     uint64
	}
	deployPlan := make([]deployEntry, contractCount)
	walletNonces := make([]uint64, walletCount)
	for i := range walletCount {
		walletNonces[i] = wallets[i].GetNonce()
	}
	for i := uint64(0); i < contractCount; i++ {
		wIdx := int(i) % walletCount
		deployPlan[i] = deployEntry{walletIdx: wIdx, nonce: walletNonces[wIdx]}
		walletNonces[wIdx]++
	}

	predictedAddrs := make([]common.Address, contractCount)
	for i := uint64(0); i < contractCount; i++ {
		e := deployPlan[i]
		predictedAddrs[i] = crypto.CreateAddress(wallets[e.walletIdx].GetAddress(), e.nonce)
	}

	// Build all deploy transactions, grouped by wallet
	walletTxs := make(map[*spamoor.Wallet][]*types.Transaction, walletCount)
	runtimes := make([][]byte, contractCount)

	// Track which contract index each wallet's tx corresponds to
	type txMapping struct {
		contractIdx int
		runtime     []byte
	}
	walletMappings := make(map[*spamoor.Wallet][]txMapping, walletCount)

	for i := startIdx; i < totalDeploys; i++ {
		localIdx := i - startIdx
		peerAddrs := make([]common.Address, 0, contractCount-1)
		for j := uint64(0); j < contractCount; j++ {
			if j != localIdx {
				peerAddrs = append(peerAddrs, predictedAddrs[j])
			}
		}

		runtime := p.generateRuntimeBytecode(i, peerAddrs)
		initCode := wrapInInitCode(runtime)
		runtimes[localIdx] = runtime

		e := deployPlan[localIdx]
		w := wallets[e.walletIdx]

		txData, buildErr := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       deployGasLimit,
			To:        nil, // Contract creation
			Data:      initCode,
		})
		if buildErr != nil {
			return fmt.Errorf("build deploy tx %d failed: %w", i, buildErr)
		}

		tx, signErr := w.BuildDynamicFeeTx(txData)
		if signErr != nil {
			return fmt.Errorf("sign deploy tx %d failed: %w", i, signErr)
		}

		walletTxs[w] = append(walletTxs[w], tx)
		walletMappings[w] = append(walletMappings[w], txMapping{
			contractIdx: int(localIdx),
			runtime:     runtime,
		})
	}

	// Multi-wallet batch send and wait for all receipts
	allReceipts, err := walletPool.GetTxPool().SendMultiTransactionBatch(ctx, walletTxs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: clientGroup,
			Rebroadcast: true,
		},
		MaxRetries: 3,
		LogFn: func(confirmed, total int) {
			p.logger.Infof("pool deployment progress: %d/%d confirmed", confirmed, total)
		},
		LogInterval: 10,
	})
	if err != nil {
		return fmt.Errorf("batch deploy failed: %w", err)
	}

	// Process receipts and populate the pool in deploy order
	type poolEntry struct {
		addr    common.Address
		runtime []byte
		idx     int
	}
	entries := make([]poolEntry, 0, contractCount)

	for w, receipts := range allReceipts {
		mappings := walletMappings[w]
		for i, receipt := range receipts {
			if i >= len(mappings) {
				break
			}
			m := mappings[i]
			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				entries = append(entries, poolEntry{
					addr:    receipt.ContractAddress,
					runtime: m.runtime,
					idx:     m.contractIdx,
				})
				p.logger.Debugf("deployed contract %d at %s (%d bytes)",
					startIdx+uint64(m.contractIdx), receipt.ContractAddress.Hex(), len(m.runtime))
			} else {
				p.logger.Warnf("deploy tx %d failed (receipt: %v)",
					startIdx+uint64(m.contractIdx), receipt)
			}
		}
	}

	// Sort by original deploy order for deterministic pool state
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].idx < entries[j].idx
	})

	p.mu.Lock()
	for _, e := range entries {
		p.contracts = append(p.contracts, e.addr)
		p.bytecodes = append(p.bytecodes, e.runtime)
	}
	p.mu.Unlock()

	p.logger.Infof("contract pool initialized with %d contracts", len(p.contracts))

	return nil
}

// DeployNew deploys a new fuzz contract and adds it to the pool,
// rotating out the oldest contract if the pool is full.
// The deploy index is computed from the effectiveTxID.
func (p *ContractPool) DeployNew(
	ctx context.Context,
	walletPool *spamoor.WalletPool,
	effectiveTxID uint64,
	baseFee, tipFee float64,
	baseFeeWei, tipFeeWei string,
	deployGasLimit uint64,
	clientGroup string,
	txIdx uint64,
) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	deployIdx := p.DeployIndexForTx(effectiveTxID)
	peerAddrs := p.GetAddresses() // current pool contracts as peers
	runtime := p.generateRuntimeBytecode(deployIdx, peerAddrs)
	initCode := wrapInInitCode(runtime)

	wallet := walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, int(txIdx))
	if wallet == nil {
		return nil, nil, nil, nil, fmt.Errorf("no wallet available")
	}

	client := walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(clientGroup),
	)
	if client == nil {
		return nil, nil, nil, wallet, fmt.Errorf("no client available")
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, err
	}

	feeCapWei, tipCapWei := spamoor.ResolveFees(baseFee, tipFee, baseFeeWei, tipFeeWei)
	feeCap, tipCap, err := walletPool.GetTxPool().GetSuggestedFees(client, feeCapWei, tipCapWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       deployGasLimit,
		To:        nil,
		Data:      initCode,
	})
	if err != nil {
		return nil, nil, client, wallet, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)

	err = walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: clientGroup,
		Rebroadcast: true,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			if err == nil && receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				p.addContract(receipt.ContractAddress, runtime)
			}
			receiptChan <- receipt
		},
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// addContract adds a contract to the pool, rotating out the oldest if full.
func (p *ContractPool) addContract(addr common.Address, runtime []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if uint64(len(p.contracts)) >= p.poolSize {
		// Ring buffer: rotate out the oldest
		p.contracts = append(p.contracts[1:], addr)
		p.bytecodes = append(p.bytecodes[1:], runtime)
	} else {
		p.contracts = append(p.contracts, addr)
		p.bytecodes = append(p.bytecodes, runtime)
	}
}

// quickRNG is a minimal RNG for non-crypto purposes (pool size selection).
type quickRNG struct {
	state uint64
}

func newQuickRNG(idx uint64, seed string) *quickRNG {
	h := uint64(0x9e3779b97f4a7c15)
	for _, b := range []byte(seed) {
		h ^= uint64(b)
		h *= 0x2545f4914f6cdd1d
	}
	h ^= idx
	h *= 0x2545f4914f6cdd1d
	if h == 0 {
		h = 1
	}

	return &quickRNG{state: h}
}

func (r *quickRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	r.state ^= r.state >> 12
	r.state ^= r.state << 25
	r.state ^= r.state >> 27

	return int((r.state * 0x2545f4914f6cdd1d) % uint64(n))
}
