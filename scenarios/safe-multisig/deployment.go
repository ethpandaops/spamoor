package safemultisig

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/safe-multisig/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// DeploymentInfo holds the addresses and bound instances of the shared Safe
// infrastructure deployed (or resolved) for a scenario run.
type DeploymentInfo struct {
	SingletonAddr common.Address
	Singleton     *contract.Safe
	FactoryAddr   common.Address
	Factory       *contract.SafeProxyFactory
	GasBurnerAddr common.Address
	GasBurner     *contract.GasBurner

	// proxyCreationCode is SafeProxyFactory.proxyCreationCode(), read once and
	// reused to compute counterfactual proxy addresses for restart-safety.
	proxyCreationCode []byte
}

// safeSpec is the deterministic description of one multisig: which child wallets
// own it, the signing threshold, and the CREATE2 salt nonce used to deploy it.
// Everything is derived deterministically from (executor wallet, slot index) so a
// restarted scenario reproduces the same safes instead of creating duplicates.
type safeSpec struct {
	executor  *spamoor.Wallet
	owners    []*signer // sorted ascending by address
	threshold int
	saltNonce *big.Int
	addr      common.Address // counterfactual proxy address
}

// DeployContracts deploys the shared Safe stack from the "deployer" wallet: the
// Safe singleton (master copy), the SafeProxyFactory, and the GasBurner call
// target. It uses the CREATE2 deployment factory (like the other multi-contract
// scenarios) so a restarted scenario re-uses the already-deployed addresses.
func (s *Scenario) DeployContracts(ctx context.Context) (*DeploymentInfo, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.deployClientGroup()),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
	if deployerWallet == nil {
		return nil, scenario.ErrNoWallet
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := []*types.Transaction{}
	info := &DeploymentInfo{}

	// All three are no-constructor-arg contracts deployed via the CREATE2 factory
	// with a zero seed, so their addresses are stable across scenarios.
	deploy := func(metadata *bind.MetaData) (common.Address, error) {
		initCode := common.FromHex(metadata.Bin)
		seed := [32]byte{}
		addr, tx, err := s.walletPool.GetDeploymentFactory().GetContractDeployment(ctx, initCode, seed, client, deployerWallet, feeCap, tipCap, false)
		if err != nil {
			return common.Address{}, err
		}
		if tx != nil {
			deploymentTxs = append(deploymentTxs, tx)
		}
		return addr, nil
	}

	info.SingletonAddr, err = deploy(contract.SafeMetaData)
	if err != nil {
		return nil, fmt.Errorf("could not deploy Safe singleton: %w", err)
	}
	info.Singleton, err = contract.NewSafe(info.SingletonAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create Safe instance: %w", err)
	}

	info.FactoryAddr, err = deploy(contract.SafeProxyFactoryMetaData)
	if err != nil {
		return nil, fmt.Errorf("could not deploy SafeProxyFactory: %w", err)
	}
	info.Factory, err = contract.NewSafeProxyFactory(info.FactoryAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create SafeProxyFactory instance: %w", err)
	}

	info.GasBurnerAddr, err = deploy(contract.GasBurnerMetaData)
	if err != nil {
		return nil, fmt.Errorf("could not deploy GasBurner: %w", err)
	}
	info.GasBurner, err = contract.NewGasBurner(info.GasBurnerAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create GasBurner instance: %w", err)
	}

	if len(deploymentTxs) > 0 {
		_, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, deployerWallet, deploymentTxs, &spamoor.BatchOptions{
			SendTransactionOptions: spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: s.options.ClientGroup,
			},
			MaxRetries:   3,
			PendingLimit: 10,
		})
		if err != nil {
			return nil, fmt.Errorf("could not send deployment txs: %w", err)
		}
		s.logger.Infof("deployed safe infrastructure (%d txs)", len(deploymentTxs))
	}

	// Cache the proxy creation code used for counterfactual address computation.
	info.proxyCreationCode, err = info.Factory.ProxyCreationCode(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("could not read proxy creation code: %w", err)
	}

	return info, nil
}

// setupSafes deterministically derives the per-wallet multisigs, creates the ones
// that don't exist yet (restart-safe via CREATE2 address probing), validates the
// off-chain SafeTx hashing against the on-chain Safe once, and populates the
// per-executor safe pools used by the transaction loop.
func (s *Scenario) setupSafes(ctx context.Context) error {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.deployClientGroup()),
	)
	if client == nil {
		return scenario.ErrNoClients
	}
	ethClient := client.GetEthClient()

	specs, err := s.buildSafeSpecs(ctx)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("no safes to set up (wallet count is zero)")
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return fmt.Errorf("could not get tx fee: %w", err)
	}

	walletTxs := make(map[*spamoor.Wallet][]*types.Transaction)
	createdCount := 0
	for _, spec := range specs {
		code, err := ethClient.CodeAt(ctx, spec.addr, nil)
		if err != nil {
			return fmt.Errorf("could not check safe code at %s: %w", spec.addr.Hex(), err)
		}
		if len(code) > 0 {
			// Already deployed (restart) - reuse without a creation tx.
			continue
		}

		initializer, err := s.encodeSetup(spec)
		if err != nil {
			return err
		}

		createTx, err := spec.executor.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return s.deploymentInfo.Factory.CreateProxyWithNonce(transactOpts, s.deploymentInfo.SingletonAddr, initializer, spec.saltNonce)
		})
		if err != nil {
			return fmt.Errorf("could not build safe creation tx: %w", err)
		}
		walletTxs[spec.executor] = append(walletTxs[spec.executor], createTx)
		createdCount++
	}

	if createdCount > 0 {
		s.logger.Infof("creating %d safe multisigs across %d wallets...", createdCount, len(walletTxs))
		_, err := s.walletPool.GetTxPool().SendMultiTransactionBatch(ctx, walletTxs, &spamoor.BatchOptions{
			SendTransactionOptions: spamoor.SendTransactionOptions{
				ClientGroup: s.options.ClientGroup,
			},
			MaxRetries:   3,
			PendingLimit: 50,
			LogFn: func(confirmedCount int, totalCount int) {
				s.logger.Infof("creating safes... (%v/%v)", confirmedCount, totalCount)
			},
			LogInterval: 20,
		})
		if err != nil {
			return fmt.Errorf("could not send safe creation txs: %w", err)
		}
	} else {
		s.logger.Infof("all %d safe multisigs already deployed, reusing", len(specs))
	}

	// Verify every safe now has code, read each safe's current nonce, and build
	// the per-executor pools.
	chainID := s.walletPool.GetRootWallet().GetWallet().GetChainId()
	verified := false
	for _, spec := range specs {
		code, err := ethClient.CodeAt(ctx, spec.addr, nil)
		if err != nil {
			return fmt.Errorf("could not verify safe code at %s: %w", spec.addr.Hex(), err)
		}
		if len(code) == 0 {
			return fmt.Errorf("safe %s has no code after creation", spec.addr.Hex())
		}

		safeInstance, err := contract.NewSafe(spec.addr, ethClient)
		if err != nil {
			return fmt.Errorf("could not create safe instance %s: %w", spec.addr.Hex(), err)
		}
		nonce, err := safeInstance.Nonce(&bind.CallOpts{Context: ctx})
		if err != nil {
			return fmt.Errorf("could not read safe nonce %s: %w", spec.addr.Hex(), err)
		}

		domainSep, err := computeDomainSeparator(chainID, spec.addr)
		if err != nil {
			return err
		}

		// One-time validation of the whole off-chain hashing path against the
		// on-chain Safe, covering the domain separator and SafeTx struct hashing.
		if !verified {
			if err := s.verifyHashing(ctx, safeInstance, spec.addr, domainSep); err != nil {
				return err
			}
			verified = true
		}

		entry := &safeEntry{
			addr:      spec.addr,
			instance:  safeInstance,
			owners:    spec.owners,
			threshold: spec.threshold,
			nonce:     nonce.Uint64(),
			domainSep: domainSep,
			balance:   big.NewInt(0),
		}
		s.poolForWallet(spec.executor.GetAddress()).safes = append(s.poolForWallet(spec.executor.GetAddress()).safes, entry)
	}

	return nil
}

// buildSafeSpecs derives the deterministic multisig specs created upfront for
// every child wallet. In recreate mode only one safe per wallet is created
// upfront (the rest grow lazily in the transaction loop); otherwise the full
// per-wallet count is created.
func (s *Scenario) buildSafeSpecs(ctx context.Context) ([]safeSpec, error) {
	walletCount := s.walletPool.GetConfiguredWalletCount()
	if walletCount == 0 {
		return nil, nil
	}

	perWallet := s.options.SafesPerWallet
	if s.options.RecreateRate > 0 {
		perWallet = 1
	}

	specs := make([]safeSpec, 0, walletCount*perWallet)
	for w := uint64(0); w < walletCount; w++ {
		executor := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(w))
		if executor == nil {
			continue
		}
		for j := uint64(0); j < perWallet; j++ {
			spec, err := s.buildSafeSpec(ctx, executor, w, j, walletCount)
			if err != nil {
				return nil, err
			}
			specs = append(specs, spec)
		}
	}
	return specs, nil
}

// buildSafeSpec derives one multisig spec deterministically from the executor
// wallet address and the per-wallet slot index. The owner count, owner set, and
// threshold are all seed-derived so they are reproduced exactly across restarts.
func (s *Scenario) buildSafeSpec(ctx context.Context, executor *spamoor.Wallet, executorIdx, slot, walletCount uint64) (safeSpec, error) {
	seedBuf := make([]byte, 0, 32)
	seedBuf = append(seedBuf, executor.GetAddress().Bytes()...)
	slotBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(slotBytes, slot)
	seedBuf = append(seedBuf, slotBytes...)
	seed := crypto.Keccak256(seedBuf)

	// Salt nonce unique per (executor, slot) so addresses never collide even when
	// two safes happen to share an owner set.
	saltNonce := new(big.Int).SetUint64(executorIdx*s.options.SafesPerWallet + slot)

	return s.buildSafeSpecFromSeed(executor, seed, saltNonce, walletCount)
}

// buildSafeSpecFromSeed derives a multisig spec (owner count, owner set,
// threshold, and counterfactual proxy address) from a seed and salt. A
// deterministic seed reproduces the same safe across restarts (initial setup); a
// random seed produces a fresh pseudo-random safe (recreation).
func (s *Scenario) buildSafeSpecFromSeed(executor *spamoor.Wallet, seed []byte, saltNonce *big.Int, walletCount uint64) (safeSpec, error) {
	seedInt := new(big.Int).SetBytes(seed)

	// Owner count varies within [minOwners, maxOwners], clamped to the available
	// wallet pool size.
	maxOwners := s.options.MaxOwners
	if maxOwners > walletCount {
		maxOwners = walletCount
	}
	minOwners := s.options.MinOwners
	if minOwners > maxOwners {
		minOwners = maxOwners
	}
	span := maxOwners - minOwners + 1
	ownerCount := minOwners + new(big.Int).Mod(seedInt, big.NewInt(int64(span))).Uint64()
	if ownerCount == 0 {
		ownerCount = 1
	}

	owners, err := s.selectOwners(seed, ownerCount, walletCount)
	if err != nil {
		return safeSpec{}, err
	}
	sortSignersByAddress(owners)

	// Threshold: 0 means n-of-n; otherwise clamp the configured value to the
	// owner count.
	threshold := int(ownerCount)
	if s.options.Threshold > 0 && uint64(threshold) > s.options.Threshold {
		threshold = int(s.options.Threshold)
	}
	if threshold < 1 {
		threshold = 1
	}

	addr, err := s.computeProxyAddress(owners, threshold, saltNonce)
	if err != nil {
		return safeSpec{}, err
	}

	return safeSpec{
		executor:  executor,
		owners:    owners,
		threshold: threshold,
		saltNonce: saltNonce,
		addr:      addr,
	}, nil
}

// selectOwners picks ownerCount distinct child wallets deterministically from the
// seed. It walks a seed-derived stride over the wallet index space, deduping, so
// the same owner set is reproduced across restarts.
func (s *Scenario) selectOwners(seed []byte, ownerCount, walletCount uint64) ([]*signer, error) {
	if ownerCount > walletCount {
		ownerCount = walletCount
	}

	owners := make([]*signer, 0, ownerCount)
	used := make(map[uint64]bool, ownerCount)
	cursor := new(big.Int).SetBytes(seed).Uint64() % walletCount
	for uint64(len(owners)) < ownerCount {
		if !used[cursor] {
			used[cursor] = true
			wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(cursor))
			if wallet == nil {
				return nil, fmt.Errorf("owner wallet %d not available", cursor)
			}
			key := wallet.GetPrivateKey()
			if key == nil {
				return nil, fmt.Errorf("owner wallet %s has no private key", wallet.GetAddress().Hex())
			}
			owners = append(owners, &signer{addr: wallet.GetAddress(), key: key})
		}
		cursor = (cursor + 1) % walletCount
	}
	return owners, nil
}

// encodeSetup ABI-encodes the Safe.setup(...) initializer used by
// createProxyWithNonce: the owners and threshold with no module/fallback/refund
// configuration (the zero-refund execution path).
func (s *Scenario) encodeSetup(spec safeSpec) ([]byte, error) {
	safeABI, err := contract.SafeMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("could not load Safe abi: %w", err)
	}
	ownerAddrs := make([]common.Address, len(spec.owners))
	for i, o := range spec.owners {
		ownerAddrs[i] = o.addr
	}
	initializer, err := safeABI.Pack(
		"setup",
		ownerAddrs,
		big.NewInt(int64(spec.threshold)),
		common.Address{}, // to
		[]byte{},         // data
		common.Address{}, // fallbackHandler
		common.Address{}, // paymentToken
		big.NewInt(0),    // payment
		common.Address{}, // paymentReceiver
	)
	if err != nil {
		return nil, fmt.Errorf("could not pack setup: %w", err)
	}
	return initializer, nil
}

// computeProxyAddress reproduces SafeProxyFactory.createProxyWithNonce's CREATE2
// address derivation so the scenario can detect already-deployed safes:
//
//	salt           = keccak256(keccak256(initializer) || uint256(saltNonce))
//	deploymentData = proxyCreationCode || uint256(uint160(singleton))
//	address        = CREATE2(factory, salt, keccak256(deploymentData))
func (s *Scenario) computeProxyAddress(owners []*signer, threshold int, saltNonce *big.Int) (common.Address, error) {
	initializer, err := s.encodeSetup(safeSpec{owners: owners, threshold: threshold})
	if err != nil {
		return common.Address{}, err
	}

	saltPreimage := make([]byte, 0, 64)
	saltPreimage = append(saltPreimage, crypto.Keccak256(initializer)...)
	saltPreimage = append(saltPreimage, common.LeftPadBytes(saltNonce.Bytes(), 32)...)
	salt := crypto.Keccak256Hash(saltPreimage)

	deploymentData := make([]byte, 0, len(s.deploymentInfo.proxyCreationCode)+32)
	deploymentData = append(deploymentData, s.deploymentInfo.proxyCreationCode...)
	deploymentData = append(deploymentData, common.LeftPadBytes(s.deploymentInfo.SingletonAddr.Bytes(), 32)...)

	return crypto.CreateAddress2(s.deploymentInfo.FactoryAddr, salt, crypto.Keccak256(deploymentData)), nil
}

// verifyHashing validates the off-chain EIP-712 implementation against the
// on-chain Safe exactly once at startup: the domain separator and a sample
// SafeTx hash must match what the deployed contract computes. A mismatch means
// the signing path would produce invalid signatures, so it fails fast.
func (s *Scenario) verifyHashing(ctx context.Context, safeInstance *contract.Safe, addr common.Address, goDomainSep [32]byte) error {
	onchainDomainSep, err := safeInstance.DomainSeparator(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("could not read domain separator: %w", err)
	}
	if onchainDomainSep != goDomainSep {
		return fmt.Errorf("domain separator mismatch for safe %s: on-chain %x vs computed %x", addr.Hex(), onchainDomainSep, goDomainSep)
	}

	sample := &safeTxParams{
		To:        s.deploymentInfo.GasBurnerAddr,
		Value:     big.NewInt(0),
		Data:      []byte{0x01, 0x02, 0x03},
		Operation: 0,
		Nonce:     big.NewInt(0),
	}
	goHash, err := computeSafeTxHash(goDomainSep, sample)
	if err != nil {
		return err
	}
	onchainHash, err := safeInstance.GetTransactionHash(&bind.CallOpts{Context: ctx},
		sample.To, sample.Value, sample.Data, sample.Operation,
		big.NewInt(0), big.NewInt(0), big.NewInt(0),
		common.Address{}, common.Address{}, sample.Nonce,
	)
	if err != nil {
		return fmt.Errorf("could not read on-chain transaction hash: %w", err)
	}
	if onchainHash != goHash {
		return fmt.Errorf("SafeTx hash mismatch for safe %s: on-chain %x vs computed %x", addr.Hex(), onchainHash, goHash)
	}

	s.logger.Debugf("validated off-chain SafeTx hashing against on-chain safe %s", addr.Hex())
	return nil
}
