package ensnames

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/ens-names/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// usdOraclePrice is the static ETH/USD price fed to the DummyOracle ($1600,
// 8 decimals).
var usdOraclePrice = big.NewInt(160000000000)

// rentPrices are the StablePriceOracle prices in attoUSD per second by label
// length (1, 2, 3, 4, 5+ letters) - the values of the official ENS
// deployment: $640/yr for 3 letters, $160/yr for 4, $5/yr for 5+.
var rentPrices = []*big.Int{
	big.NewInt(0),
	big.NewInt(0),
	big.NewInt(20294266869609),
	big.NewInt(5073566717402),
	big.NewInt(158548959919),
}

// wrapperMetadataURL is the (dummy) NameWrapper token metadata endpoint.
const wrapperMetadataURL = "https://ens.spamoor.local/metadata/{id}"

// DeploymentInfo holds the addresses and bound instances of the ENS stack
// deployed (or resolved) for a scenario run.
type DeploymentInfo struct {
	ExecutorAddr common.Address
	Executor     *contract.EnsExecutor

	RegistryAddr common.Address
	Registry     *contract.ENSRegistry

	BaseAddr common.Address
	Base     *contract.BaseRegistrarImplementation

	DummyOracleAddr common.Address
	PriceOracleAddr common.Address

	ReverseAddr common.Address
	Reverse     *contract.ReverseRegistrar

	DefaultReverseAddr common.Address
	DefaultReverse     *contract.DefaultReverseRegistrar

	MetadataAddr common.Address

	WrapperAddr common.Address
	Wrapper     *contract.NameWrapper

	ControllerAddr common.Address
	Controller     *contract.ETHRegistrarController

	ResolverAddr common.Address
	Resolver     *contract.PublicResolver

	SpamControllerAddr common.Address
	SpamController     *contract.SpamRegistrarController
}

// stackDeployer bundles the shared state of one DeployContracts run so the
// phase helpers don't need to thread half a dozen parameters around.
type stackDeployer struct {
	s          *Scenario
	client     *spamoor.Client
	wallet     *spamoor.Wallet // "deployer" well-known wallet
	feeCap     *big.Int
	tipCap     *big.Int
	info       *DeploymentInfo
	operatorOK bool
}

// deploymentSalt derives the CREATE2 salt of the EnsExecutor from the
// configured deployment seed. An empty seed yields the zero salt (one shared,
// stable ENS stack per root key); any other seed yields an independent stack.
func deploymentSalt(seed string) [32]byte {
	salt := [32]byte{}
	if seed != "" {
		copy(salt[:], crypto.Keccak256([]byte(seed)))
	}
	return salt
}

// buildInitcode returns the creation bytecode of a contract with its ABI-packed
// constructor arguments appended.
func buildInitcode(metadata *bind.MetaData, params ...any) ([]byte, error) {
	parsed, err := metadata.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("could not parse abi: %w", err)
	}

	packed, err := parsed.Pack("", params...)
	if err != nil {
		return nil, fmt.Errorf("could not pack constructor args: %w", err)
	}

	return append(common.FromHex(metadata.Bin), packed...), nil
}

// packCall ABI-encodes a method call for use as executor.execute() calldata.
func packCall(metadata *bind.MetaData, method string, params ...any) ([]byte, error) {
	parsed, err := metadata.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("could not parse abi: %w", err)
	}

	data, err := parsed.Pack(method, params...)
	if err != nil {
		return nil, fmt.Errorf("could not pack %s call: %w", method, err)
	}

	return data, nil
}

// executorCreate2Addr computes the address of a contract deployed by the
// executor via CREATE2 with the zero salt (addresses are distinguished by
// initcode, the executor address itself carries the deployment seed).
func executorCreate2Addr(executor common.Address, initcode []byte) common.Address {
	return crypto.CreateAddress2(executor, [32]byte{}, crypto.Keccak256(initcode))
}

// DeployContracts deploys (or resolves) the full ENS stack. The EnsExecutor is
// deployed through the shared CREATE2 factory with the root wallet as admin;
// every ENS contract is then CREATE2-deployed by the executor itself so the
// executor - not the factory - owns the registry root node and all Ownable
// contracts (see EnsExecutor.sol). All addresses are deterministic per
// (root key, deployment seed) and every step is checked against chain state
// first, so restarts and concurrently running scenario instances converge on
// the same stack without redeploying or re-wiring anything.
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

	d := &stackDeployer{
		s:      s,
		client: client,
		wallet: deployerWallet,
		feeCap: feeCap,
		tipCap: tipCap,
		info:   &DeploymentInfo{},
	}

	if err := d.deployExecutor(ctx); err != nil {
		return nil, err
	}
	if err := d.computeAddresses(); err != nil {
		return nil, err
	}
	if err := d.deployCoreContracts(ctx); err != nil {
		return nil, err
	}
	if err := d.wireRootNodes(ctx); err != nil {
		return nil, err
	}
	if err := d.deployDependentContracts(ctx); err != nil {
		return nil, err
	}
	if err := d.wireGrants(ctx); err != nil {
		return nil, err
	}
	if err := d.bindContracts(); err != nil {
		return nil, err
	}

	// Log all core addresses so external tooling (resolvers, indexers,
	// dashboards) can be pointed at the stack.
	s.logger.Infof("ens stack ready:")
	s.logger.Infof("  ENSRegistry:             %s", d.info.RegistryAddr.Hex())
	s.logger.Infof("  BaseRegistrar (.eth):    %s", d.info.BaseAddr.Hex())
	s.logger.Infof("  ETHRegistrarController:  %s", d.info.ControllerAddr.Hex())
	s.logger.Infof("  PublicResolver:          %s", d.info.ResolverAddr.Hex())
	s.logger.Infof("  ReverseRegistrar:        %s", d.info.ReverseAddr.Hex())
	s.logger.Infof("  DefaultReverseRegistrar: %s", d.info.DefaultReverseAddr.Hex())
	s.logger.Infof("  NameWrapper:             %s", d.info.WrapperAddr.Hex())
	s.logger.Infof("  StablePriceOracle:       %s (usd feed: %s)", d.info.PriceOracleAddr.Hex(), d.info.DummyOracleAddr.Hex())
	s.logger.Infof("  SpamRegistrarController: %s", d.info.SpamControllerAddr.Hex())
	s.logger.Infof("  EnsExecutor (owner):     %s", d.info.ExecutorAddr.Hex())

	return d.info, nil
}

// deployExecutor deploys the EnsExecutor through the shared CREATE2 factory
// with the root wallet address as admin constructor arg (the root wallet is
// the trust anchor that later authorizes per-scenario operator wallets).
func (d *stackDeployer) deployExecutor(ctx context.Context) error {
	rootAddr := d.s.walletPool.GetRootWallet().GetWallet().GetAddress()

	initcode, err := buildInitcode(contract.EnsExecutorMetaData, rootAddr)
	if err != nil {
		return fmt.Errorf("could not build EnsExecutor initcode: %w", err)
	}

	salt := deploymentSalt(d.s.options.DeploymentSeed)
	addr, tx, err := d.s.walletPool.GetDeploymentFactory().GetContractDeployment(ctx, initcode, salt, d.client, d.wallet, d.feeCap, d.tipCap, false)
	if err != nil {
		return fmt.Errorf("could not deploy EnsExecutor: %w", err)
	}

	if tx != nil {
		if err := d.sendBatch(ctx, []*types.Transaction{tx}, "deploying ens executor"); err != nil {
			return err
		}
	}

	d.info.ExecutorAddr = addr
	d.info.Executor, err = contract.NewEnsExecutor(addr, d.client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind EnsExecutor: %w", err)
	}

	return nil
}

// computeAddresses derives the deterministic addresses of all ENS contracts
// from the executor address and each contract's initcode. Constructor args
// only reference other computed addresses, so everything is known upfront.
func (d *stackDeployer) computeAddresses() error {
	info := d.info

	registryInit, err := buildInitcode(contract.ENSRegistryMetaData)
	if err != nil {
		return err
	}
	info.RegistryAddr = executorCreate2Addr(info.ExecutorAddr, registryInit)

	baseInit, err := buildInitcode(contract.BaseRegistrarImplementationMetaData, info.RegistryAddr, ethNode)
	if err != nil {
		return err
	}
	info.BaseAddr = executorCreate2Addr(info.ExecutorAddr, baseInit)

	dummyOracleInit, err := buildInitcode(contract.DummyOracleMetaData, usdOraclePrice)
	if err != nil {
		return err
	}
	info.DummyOracleAddr = executorCreate2Addr(info.ExecutorAddr, dummyOracleInit)

	priceOracleInit, err := buildInitcode(contract.StablePriceOracleMetaData, info.DummyOracleAddr, rentPrices)
	if err != nil {
		return err
	}
	info.PriceOracleAddr = executorCreate2Addr(info.ExecutorAddr, priceOracleInit)

	reverseInit, err := buildInitcode(contract.ReverseRegistrarMetaData, info.RegistryAddr)
	if err != nil {
		return err
	}
	info.ReverseAddr = executorCreate2Addr(info.ExecutorAddr, reverseInit)

	defaultReverseInit, err := buildInitcode(contract.DefaultReverseRegistrarMetaData)
	if err != nil {
		return err
	}
	info.DefaultReverseAddr = executorCreate2Addr(info.ExecutorAddr, defaultReverseInit)

	metadataInit, err := buildInitcode(contract.StaticMetadataServiceMetaData, wrapperMetadataURL)
	if err != nil {
		return err
	}
	info.MetadataAddr = executorCreate2Addr(info.ExecutorAddr, metadataInit)

	wrapperInit, err := buildInitcode(contract.NameWrapperMetaData, info.RegistryAddr, info.BaseAddr, info.MetadataAddr)
	if err != nil {
		return err
	}
	info.WrapperAddr = executorCreate2Addr(info.ExecutorAddr, wrapperInit)

	controllerInit, err := buildInitcode(contract.ETHRegistrarControllerMetaData,
		info.BaseAddr, info.PriceOracleAddr,
		new(big.Int).SetUint64(d.s.options.MinCommitmentAge), new(big.Int).SetUint64(d.s.options.MaxCommitmentAge),
		info.ReverseAddr, info.DefaultReverseAddr, info.RegistryAddr)
	if err != nil {
		return err
	}
	info.ControllerAddr = executorCreate2Addr(info.ExecutorAddr, controllerInit)

	resolverInit, err := buildInitcode(contract.PublicResolverMetaData,
		info.RegistryAddr, info.WrapperAddr, info.ControllerAddr, info.ReverseAddr)
	if err != nil {
		return err
	}
	info.ResolverAddr = executorCreate2Addr(info.ExecutorAddr, resolverInit)

	spamControllerInit, err := buildInitcode(contract.SpamRegistrarControllerMetaData,
		info.BaseAddr, info.RegistryAddr, info.ReverseAddr, info.ResolverAddr)
	if err != nil {
		return err
	}
	info.SpamControllerAddr = executorCreate2Addr(info.ExecutorAddr, spamControllerInit)

	return nil
}

// ensureOperator makes sure the scenario's deployer wallet is authorized on
// the executor, sending the single root-wallet setOperator tx if it is not.
// This is the only transaction the root wallet ever sends for this scenario;
// it is skipped entirely when the stack is already deployed and wired.
func (d *stackDeployer) ensureOperator(ctx context.Context) error {
	if d.operatorOK {
		return nil
	}

	deployerAddr := d.wallet.GetAddress()
	isOperator, err := d.info.Executor.Operators(&bind.CallOpts{Context: ctx}, deployerAddr)
	if err != nil {
		return fmt.Errorf("could not check executor operator: %w", err)
	}
	if isOperator {
		d.operatorOK = true
		return nil
	}

	rootWallet := d.s.walletPool.GetRootWallet()
	err = rootWallet.WithWalletLock(ctx, 1, nil, d.s.walletPool.GetClientPool(), func(reason string) {
		d.s.logger.Infof("root wallet is locked, %s", reason)
	}, func() error {
		tx, err := rootWallet.GetWallet().BuildBoundTxWithEstimate(ctx, d.client, d.s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(d.feeCap),
			GasTipCap: uint256.MustFromBig(d.tipCap),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return d.info.Executor.SetOperator(transactOpts, deployerAddr, true)
		})
		if err != nil {
			return fmt.Errorf("could not build setOperator tx: %w", err)
		}

		return d.sendBatchFrom(ctx, rootWallet.GetWallet(), []*types.Transaction{tx}, "authorizing ens deployer")
	})
	if err != nil {
		return err
	}

	isOperator, err = d.info.Executor.Operators(&bind.CallOpts{Context: ctx}, deployerAddr)
	if err != nil {
		return fmt.Errorf("could not verify executor operator: %w", err)
	}
	if !isOperator {
		return fmt.Errorf("deployer %s is not an executor operator after setOperator", deployerAddr.Hex())
	}

	d.operatorOK = true
	return nil
}

// deployTx builds an executor.deploy() tx for one contract if it does not
// exist yet, ensuring operator rights first.
func (d *stackDeployer) deployTx(ctx context.Context, name string, addr common.Address, initcode []byte) (*types.Transaction, error) {
	if code, err := d.client.GetCodeAt(ctx, addr); err == nil && len(code) > 0 {
		d.s.logger.Debugf("%s already deployed at %s, reusing", name, addr.Hex())
		return nil, nil
	}

	if err := d.ensureOperator(ctx); err != nil {
		return nil, err
	}

	tx, err := d.wallet.BuildBoundTxWithEstimate(ctx, d.client, d.s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(d.feeCap),
		GasTipCap: uint256.MustFromBig(d.tipCap),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return d.info.Executor.Deploy(transactOpts, [32]byte{}, initcode)
	})
	if err != nil {
		return nil, fmt.Errorf("could not build %s deploy tx: %w", name, err)
	}

	return tx, nil
}

// execTx builds an executor.execute() tx performing an owner-gated call,
// ensuring operator rights first.
func (d *stackDeployer) execTx(ctx context.Context, what string, target common.Address, data []byte) (*types.Transaction, error) {
	if err := d.ensureOperator(ctx); err != nil {
		return nil, err
	}

	tx, err := d.wallet.BuildBoundTxWithEstimate(ctx, d.client, d.s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(d.feeCap),
		GasTipCap: uint256.MustFromBig(d.tipCap),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return d.info.Executor.Execute(transactOpts, target, big.NewInt(0), data)
	})
	if err != nil {
		return nil, fmt.Errorf("could not build %s tx: %w", what, err)
	}

	return tx, nil
}

// deployContractSet deploys a list of contracts via the executor in one mined
// batch, skipping the ones that already exist.
func (d *stackDeployer) deployContractSet(ctx context.Context, action string, contracts []executorDeployment) error {
	txs := make([]*types.Transaction, 0, len(contracts))
	for _, c := range contracts {
		initcode, err := buildInitcode(c.metadata, c.params...)
		if err != nil {
			return fmt.Errorf("could not build %s initcode: %w", c.name, err)
		}

		tx, err := d.deployTx(ctx, c.name, c.addr, initcode)
		if err != nil {
			return err
		}
		if tx != nil {
			txs = append(txs, tx)
		}
	}

	return d.sendBatch(ctx, txs, action)
}

// executorDeployment describes one contract deployed via executor.deploy().
type executorDeployment struct {
	name     string
	addr     common.Address
	metadata *bind.MetaData
	params   []any
}

// deployCoreContracts deploys the base contracts in two mined batches:
// registry, base registrar, oracles and metadata service first, then the
// ReverseRegistrar - its constructor reads registry.owner(addr.reverse), so
// its gas estimation needs the registry to already exist on-chain.
func (d *stackDeployer) deployCoreContracts(ctx context.Context) error {
	info := d.info

	err := d.deployContractSet(ctx, "deploying ens core contracts", []executorDeployment{
		{"ENSRegistry", info.RegistryAddr, contract.ENSRegistryMetaData, nil},
		{"BaseRegistrarImplementation", info.BaseAddr, contract.BaseRegistrarImplementationMetaData, []any{info.RegistryAddr, ethNode}},
		{"DummyOracle", info.DummyOracleAddr, contract.DummyOracleMetaData, []any{usdOraclePrice}},
		{"StablePriceOracle", info.PriceOracleAddr, contract.StablePriceOracleMetaData, []any{info.DummyOracleAddr, rentPrices}},
		{"DefaultReverseRegistrar", info.DefaultReverseAddr, contract.DefaultReverseRegistrarMetaData, nil},
		{"StaticMetadataService", info.MetadataAddr, contract.StaticMetadataServiceMetaData, []any{wrapperMetadataURL}},
	})
	if err != nil {
		return err
	}

	return d.deployContractSet(ctx, "deploying ens reverse registrar", []executorDeployment{
		{"ReverseRegistrar", info.ReverseAddr, contract.ReverseRegistrarMetaData, []any{info.RegistryAddr}},
	})
}

// wireRootNodes hands the top-level registry nodes to their owners: eth to the
// base registrar, reverse to the executor and addr.reverse to the reverse
// registrar. The addr.reverse step runs in a second mined batch because its
// gas estimation requires the executor to already own the reverse node.
func (d *stackDeployer) wireRootNodes(ctx context.Context) error {
	info := d.info

	registry, err := contract.NewENSRegistry(info.RegistryAddr, d.client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind ENSRegistry: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx}

	txs := make([]*types.Transaction, 0, 2)

	if owner, err := registry.Owner(callOpts, ethNode); err != nil {
		return fmt.Errorf("could not read eth node owner: %w", err)
	} else if owner != info.BaseAddr {
		data, err := packCall(contract.ENSRegistryMetaData, "setSubnodeOwner", rootNode, labelhash("eth"), info.BaseAddr)
		if err != nil {
			return err
		}
		tx, err := d.execTx(ctx, "setSubnodeOwner(eth)", info.RegistryAddr, data)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	if owner, err := registry.Owner(callOpts, reverseNode); err != nil {
		return fmt.Errorf("could not read reverse node owner: %w", err)
	} else if owner != info.ExecutorAddr {
		data, err := packCall(contract.ENSRegistryMetaData, "setSubnodeOwner", rootNode, labelhash("reverse"), info.ExecutorAddr)
		if err != nil {
			return err
		}
		tx, err := d.execTx(ctx, "setSubnodeOwner(reverse)", info.RegistryAddr, data)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	if err := d.sendBatch(ctx, txs, "wiring ens root nodes"); err != nil {
		return err
	}

	if owner, err := registry.Owner(callOpts, addrReverseNode); err != nil {
		return fmt.Errorf("could not read addr.reverse node owner: %w", err)
	} else if owner != info.ReverseAddr {
		data, err := packCall(contract.ENSRegistryMetaData, "setSubnodeOwner", reverseNode, labelhash("addr"), info.ReverseAddr)
		if err != nil {
			return err
		}
		tx, err := d.execTx(ctx, "setSubnodeOwner(addr.reverse)", info.RegistryAddr, data)
		if err != nil {
			return err
		}
		if err := d.sendBatch(ctx, []*types.Transaction{tx}, "wiring addr.reverse node"); err != nil {
			return err
		}
	}

	return nil
}

// deployDependentContracts deploys the contracts whose constructors need the
// wired registry state (NameWrapper and PublicResolver claim their reverse
// record from the reverse registrar at construction), plus the registrar
// controllers.
func (d *stackDeployer) deployDependentContracts(ctx context.Context) error {
	info := d.info

	return d.deployContractSet(ctx, "deploying ens registrar contracts", []executorDeployment{
		{"NameWrapper", info.WrapperAddr, contract.NameWrapperMetaData, []any{info.RegistryAddr, info.BaseAddr, info.MetadataAddr}},
		{"ETHRegistrarController", info.ControllerAddr, contract.ETHRegistrarControllerMetaData, []any{
			info.BaseAddr, info.PriceOracleAddr,
			new(big.Int).SetUint64(d.s.options.MinCommitmentAge), new(big.Int).SetUint64(d.s.options.MaxCommitmentAge),
			info.ReverseAddr, info.DefaultReverseAddr, info.RegistryAddr}},
		{"PublicResolver", info.ResolverAddr, contract.PublicResolverMetaData, []any{
			info.RegistryAddr, info.WrapperAddr, info.ControllerAddr, info.ReverseAddr}},
		{"SpamRegistrarController", info.SpamControllerAddr, contract.SpamRegistrarControllerMetaData, []any{
			info.BaseAddr, info.RegistryAddr, info.ReverseAddr, info.ResolverAddr}},
	})
}

// wireGrants applies the controller grants and the default reverse resolver:
// the commit-reveal controller, the NameWrapper and the SpamRegistrarController
// on the base registrar, plus the reverse registrar controller rights that let
// the controllers set reverse records for third-party addresses.
func (d *stackDeployer) wireGrants(ctx context.Context) error {
	info := d.info
	callOpts := &bind.CallOpts{Context: ctx}

	base, err := contract.NewBaseRegistrarImplementation(info.BaseAddr, d.client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind BaseRegistrarImplementation: %w", err)
	}
	reverse, err := contract.NewReverseRegistrar(info.ReverseAddr, d.client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind ReverseRegistrar: %w", err)
	}
	defaultReverse, err := contract.NewDefaultReverseRegistrar(info.DefaultReverseAddr, d.client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind DefaultReverseRegistrar: %w", err)
	}

	txs := make([]*types.Transaction, 0, 7)

	for _, grant := range []struct {
		what       string
		controller common.Address
	}{
		{"addController(eth controller)", info.ControllerAddr},
		{"addController(name wrapper)", info.WrapperAddr},
		{"addController(spam controller)", info.SpamControllerAddr},
	} {
		enabled, err := base.Controllers(callOpts, grant.controller)
		if err != nil {
			return fmt.Errorf("could not read base controller state: %w", err)
		}
		if enabled {
			continue
		}

		data, err := packCall(contract.BaseRegistrarImplementationMetaData, "addController", grant.controller)
		if err != nil {
			return err
		}
		tx, err := d.execTx(ctx, grant.what, info.BaseAddr, data)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	for _, grant := range []struct {
		what       string
		controller common.Address
	}{
		{"reverse setController(eth controller)", info.ControllerAddr},
		{"reverse setController(spam controller)", info.SpamControllerAddr},
	} {
		enabled, err := reverse.Controllers(callOpts, grant.controller)
		if err != nil {
			return fmt.Errorf("could not read reverse controller state: %w", err)
		}
		if enabled {
			continue
		}

		data, err := packCall(contract.ReverseRegistrarMetaData, "setController", grant.controller, true)
		if err != nil {
			return err
		}
		tx, err := d.execTx(ctx, grant.what, info.ReverseAddr, data)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	if enabled, err := defaultReverse.Controllers(callOpts, info.ControllerAddr); err != nil {
		return fmt.Errorf("could not read default reverse controller state: %w", err)
	} else if !enabled {
		data, err := packCall(contract.DefaultReverseRegistrarMetaData, "setController", info.ControllerAddr, true)
		if err != nil {
			return err
		}
		tx, err := d.execTx(ctx, "default reverse setController(eth controller)", info.DefaultReverseAddr, data)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	if resolver, err := reverse.DefaultResolver(callOpts); err != nil {
		return fmt.Errorf("could not read reverse default resolver: %w", err)
	} else if resolver != info.ResolverAddr {
		data, err := packCall(contract.ReverseRegistrarMetaData, "setDefaultResolver", info.ResolverAddr)
		if err != nil {
			return err
		}
		tx, err := d.execTx(ctx, "reverse setDefaultResolver", info.ReverseAddr, data)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	return d.sendBatch(ctx, txs, "wiring ens grants")
}

// bindContracts creates the bound contract instances for the scenario run.
func (d *stackDeployer) bindContracts() error {
	info := d.info
	ethClient := d.client.GetEthClient()

	var err error
	if info.Registry, err = contract.NewENSRegistry(info.RegistryAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind ENSRegistry: %w", err)
	}
	if info.Base, err = contract.NewBaseRegistrarImplementation(info.BaseAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind BaseRegistrarImplementation: %w", err)
	}
	if info.Reverse, err = contract.NewReverseRegistrar(info.ReverseAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind ReverseRegistrar: %w", err)
	}
	if info.DefaultReverse, err = contract.NewDefaultReverseRegistrar(info.DefaultReverseAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind DefaultReverseRegistrar: %w", err)
	}
	if info.Wrapper, err = contract.NewNameWrapper(info.WrapperAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind NameWrapper: %w", err)
	}
	if info.Controller, err = contract.NewETHRegistrarController(info.ControllerAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind ETHRegistrarController: %w", err)
	}
	if info.Resolver, err = contract.NewPublicResolver(info.ResolverAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind PublicResolver: %w", err)
	}
	if info.SpamController, err = contract.NewSpamRegistrarController(info.SpamControllerAddr, ethClient); err != nil {
		return fmt.Errorf("could not bind SpamRegistrarController: %w", err)
	}

	return nil
}

// sendBatch submits a batch of deployment txs from the deployer wallet and
// awaits confirmation.
func (d *stackDeployer) sendBatch(ctx context.Context, txs []*types.Transaction, action string) error {
	return d.sendBatchFrom(ctx, d.wallet, txs, action)
}

// sendBatchFrom submits a batch of txs from the given wallet and awaits
// confirmation.
func (d *stackDeployer) sendBatchFrom(ctx context.Context, wallet *spamoor.Wallet, txs []*types.Transaction, action string) error {
	if len(txs) == 0 {
		return nil
	}

	_, err := d.s.walletPool.GetTxPool().SendTransactionBatch(ctx, wallet, txs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      d.client,
			ClientGroup: d.s.deployClientGroup(),
		},
		MaxRetries:   3,
		PendingLimit: 10,
		LogFn: func(confirmedCount int, totalCount int) {
			d.s.logger.Infof("%s... (%v/%v)", action, confirmedCount, totalCount)
		},
		LogInterval: 10,
	})
	if err != nil {
		return fmt.Errorf("could not send %s txs: %w", action, err)
	}

	d.s.logger.Infof("%s complete (%d txs)", action, len(txs))
	return nil
}
