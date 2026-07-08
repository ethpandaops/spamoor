package aavelending

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/aave-lending/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// Library link placeholders for the unlinked @aave/core-v3@1.19.3 creation
// bytecode (see contract/compile.sh). solc emits `__$<keccak34>$__` markers for
// external libraries; these strings are stable for the pinned package version
// and are resolved to deployed library addresses at deploy time. Pool links all
// seven logic libraries, PoolConfigurator links ConfiguratorLogic, and
// FlashLoanLogic links BorrowLogic.
const (
	phBorrowLogic       = "__$c3724b8d563dc83a94e797176cddecb3b9$__"
	phBridgeLogic       = "__$b06080f092f400a43662c3f835a4d9baa8$__"
	phEModeLogic        = "__$e4b9550ff526a295e1233dea02821b9004$__"
	phFlashLoanLogic    = "__$d5ddd09ae98762b8929dd85e54b218e259$__"
	phLiquidationLogic  = "__$f598c634f2d943205ac23f707b80075cbb$__"
	phPoolLogic         = "__$563c746fa3df0f1858d85f6ef4258864be$__"
	phSupplyLogic       = "__$db79717e66442ee197e8271d032a066e34$__"
	phConfiguratorLogic = "__$3ddc574512022f331a6a4c7e4bbb5c67b6$__"
)

// CREATE2 salts: the deployment factory derives each address from the deployer,
// the salt and the init code (creation bytecode + constructor args). Contracts
// with distinct init code therefore get distinct addresses at salt 0, so a salt
// is only needed to disambiguate deployments of identical init code. Here that
// is solely the two mock price aggregators (same bytecode and price), which are
// salted by their reserve index; everything else deploys at salt 0.

// Fixed gas limits for deterministic deployment/setup transactions. Using fixed
// limits avoids eth_estimateGas round trips and, more importantly, lets us batch
// setup calls whose targets are deployed earlier in the same batch (estimation
// would revert against a not-yet-mined contract).
const (
	gasSetter      = 250000  // provider/ACL setters, addPoolAdmin
	gasProxyDeploy = 5000000 // setPoolImpl/setPoolConfiguratorImpl deploy + initialize a proxy
	gasMint        = 400000  // MintableToken.mint: writes fresh balance + totalSupply slots (~233k under the Amsterdam state-creation fee schedule)
	gasApprove     = 250000  // MintableToken.approve: fresh allowance slot (~128k under Amsterdam)
)

// Market risk parameters (basis points). USD-priced single-tier market.
const (
	reserveLTV                = 8000      // 80%
	reserveLiquidationThresh  = 8500      // 85%
	reserveLiquidationBonus   = 10500     // 105% (must be > 10000)
	reserveFactor             = 1000      // 10%
	oracleBaseCurrencyUnit    = 100000000 // 1e8, USD with 8 decimals
	oraclePriceAnswer         = 100000000 // 1e8 == $1.00 per token
	reserveUnderlyingDecimals = 18
	variableInterestRateMode  = 2 // DataTypes.InterestRateMode.VARIABLE
)

var marketID = "Spamoor Aave Market"

// ray is 1e27, the fixed-point unit used by Aave's interest rate math.
var ray = new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil)

// maxUint256 is used for unlimited allowances and for repay/withdraw "all".
var maxUint256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

// rayFrac returns num/den expressed in ray (1e27) fixed point.
func rayFrac(num, den int64) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(ray, big.NewInt(num)), big.NewInt(den))
}

// TokenInfo describes a single deployed mock reserve asset, its price source and
// the reserve tokens (aToken / variable debt token) created for it by the pool.
type TokenInfo struct {
	Addr        common.Address
	AggAddr     common.Address
	ATokenAddr  common.Address
	VarDebtAddr common.Address
	Name        string
	Symbol      string

	Token   *contract.MintableToken
	Agg     *contract.MockAggregator
	AToken  *contract.AToken
	VarDebt *contract.VariableDebtToken
}

// DeploymentInfo holds every address (and the bound instances needed by the hot
// path) of the Aave V3 market deployed by the scenario.
type DeploymentInfo struct {
	ProviderAddr         common.Address
	ACLManagerAddr       common.Address
	PoolImplAddr         common.Address
	ConfiguratorImplAddr common.Address
	DataProviderAddr     common.Address
	OracleAddr           common.Address
	RateStrategyAddr     common.Address
	ATokenImplAddr       common.Address
	VarDebtImplAddr      common.Address
	StableDebtImplAddr   common.Address

	// proxies resolved from the addresses provider after wiring
	PoolAddr         common.Address
	ConfiguratorAddr common.Address

	// bound instances for building hot-path transactions
	Pool *contract.Pool

	Tokens []TokenInfo

	libs map[string]common.Address
}

// linkBytecode resolves the library placeholders in bin to deployed library
// addresses. It returns an error if a placeholder is missing or if any
// unresolved placeholder remains afterwards.
func linkBytecode(bin string, links map[string]common.Address) (string, error) {
	for placeholder, addr := range links {
		if !strings.Contains(bin, placeholder) {
			return "", fmt.Errorf("library placeholder %s not found in bytecode", placeholder)
		}
		// a placeholder occupies exactly 40 hex chars (20 bytes), same as an address
		bin = strings.ReplaceAll(bin, placeholder, hex.EncodeToString(addr.Bytes()))
	}
	if i := strings.Index(bin, "__$"); i >= 0 {
		end := i + 40
		if end > len(bin) {
			end = len(bin)
		}
		return "", fmt.Errorf("unresolved library placeholder near %q", bin[i:end])
	}
	return bin, nil
}

// buildDeployTx builds a deterministic CREATE2 deployment transaction for the
// given contract via the wallet pool's deployment factory. links may be nil for
// contracts without external libraries. The returned transaction is nil when the
// contract is already deployed (the address is still returned), which keeps
// re-runs idempotent.
//
// global controls the CREATE2 seed: global contracts (the stateless logic
// libraries and the permissionless-mint mock tokens) deploy at a
// deployer-independent address so they are deployed once and shared across all
// deployers; non-global contracts mix in the deployer address so every deployer
// key gets its own isolated Aave market state (Aave is very state heavy, so each
// run owns a distinct market rather than contending on shared reserve state).
func (s *Scenario) buildDeployTx(ctx context.Context, client *spamoor.Client, deployerWallet *spamoor.Wallet, feeCap, tipCap *big.Int, metadata *bind.MetaData, links map[string]common.Address, global bool, salt uint32, params ...interface{}) (common.Address, *types.Transaction, error) {
	parsed, err := metadata.GetAbi()
	if err != nil {
		return common.Address{}, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, fmt.Errorf("GetAbi returned nil")
	}

	bin := metadata.Bin
	if len(links) > 0 {
		bin, err = linkBytecode(bin, links)
		if err != nil {
			return common.Address{}, nil, err
		}
	}

	initCode := common.FromHex(bin)
	packed, err := parsed.Pack("", params...)
	if err != nil {
		return common.Address{}, nil, err
	}
	initCode = append(initCode, packed...)

	seed := [32]byte{}
	if !global {
		copy(seed[:], deployerWallet.GetAddress().Bytes())
	}
	binary.BigEndian.PutUint32(seed[28:], salt)

	return s.walletPool.GetDeploymentFactory().GetContractDeployment(ctx, initCode, seed, client, deployerWallet, feeCap, tipCap, false)
}

// DeployAaveMarket deploys the full Aave V3 market and seeds borrowable
// liquidity. It is split into sequential phases because several steps depend on
// state produced by earlier ones (library addresses linked into the Pool, the
// pool/configurator proxies created by the addresses provider, the reserves
// created by initReserves). Each phase is mined before the next is built.
func (s *Scenario) DeployAaveMarket(ctx context.Context) (*DeploymentInfo, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.DeployClientGroup),
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

	info := &DeploymentInfo{libs: make(map[string]common.Address, 8)}

	if err := s.deployImplementations(ctx, client, deployerWallet, feeCap, tipCap, info); err != nil {
		return nil, err
	}
	if err := s.wireMarket(ctx, client, deployerWallet, feeCap, tipCap, info); err != nil {
		return nil, err
	}
	if err := s.deployTokenImpls(ctx, client, deployerWallet, feeCap, tipCap, info); err != nil {
		return nil, err
	}
	if err := s.bindContracts(info); err != nil {
		return nil, err
	}
	if err := s.initReserves(ctx, client, deployerWallet, feeCap, tipCap, info); err != nil {
		return nil, err
	}
	if err := s.configureReserves(ctx, client, deployerWallet, feeCap, tipCap, info); err != nil {
		return nil, err
	}
	if err := s.bindReserveTokens(ctx, info); err != nil {
		return nil, err
	}
	if err := s.seedLiquidity(ctx, client, deployerWallet, feeCap, tipCap, info); err != nil {
		return nil, err
	}

	return info, nil
}

// deployImplementations (phase 1) deploys the logic libraries, the addresses
// provider, the Pool/PoolConfigurator implementations (linked), the data
// provider, the interest rate strategy, the mock reserve tokens, their price
// aggregators and the oracle. Every constructor argument here is a CREATE2
// address known ahead of mining, so it all fits in a single batch.
func (s *Scenario) deployImplementations(ctx context.Context, client *spamoor.Client, deployerWallet *spamoor.Wallet, feeCap, tipCap *big.Int, info *DeploymentInfo) error {
	deployerAddr := deployerWallet.GetAddress()
	var txs []*types.Transaction

	// global marks stateless, shareable contracts (the logic libraries and the
	// permissionless mock tokens) so they deploy once at a deployer-independent
	// address and are reused by every deployer. The market contracts are
	// deployer-specific so each deployer key owns an isolated Aave state.
	const (
		global = true
		perKey = false
	)
	deploy := func(what string, metadata *bind.MetaData, links map[string]common.Address, global bool, salt uint32, params ...interface{}) (common.Address, error) {
		addr, tx, err := s.buildDeployTx(ctx, client, deployerWallet, feeCap, tipCap, metadata, links, global, salt, params...)
		if err != nil {
			return common.Address{}, fmt.Errorf("could not deploy %s: %w", what, err)
		}
		if tx != nil {
			txs = append(txs, tx)
		}
		return addr, nil
	}

	var err error
	if info.ProviderAddr, err = deploy("PoolAddressesProvider", contract.PoolAddressesProviderMetaData, nil, perKey, 0, marketID, deployerAddr); err != nil {
		return err
	}

	// logic libraries (global: stateless, shared across deployers)
	if info.libs[phSupplyLogic], err = deploy("SupplyLogic", contract.SupplyLogicMetaData, nil, global, 0); err != nil {
		return err
	}
	if info.libs[phBorrowLogic], err = deploy("BorrowLogic", contract.BorrowLogicMetaData, nil, global, 0); err != nil {
		return err
	}
	if info.libs[phLiquidationLogic], err = deploy("LiquidationLogic", contract.LiquidationLogicMetaData, nil, global, 0); err != nil {
		return err
	}
	if info.libs[phEModeLogic], err = deploy("EModeLogic", contract.EModeLogicMetaData, nil, global, 0); err != nil {
		return err
	}
	if info.libs[phBridgeLogic], err = deploy("BridgeLogic", contract.BridgeLogicMetaData, nil, global, 0); err != nil {
		return err
	}
	if info.libs[phPoolLogic], err = deploy("PoolLogic", contract.PoolLogicMetaData, nil, global, 0); err != nil {
		return err
	}
	// ConfiguratorLogic is linked into PoolConfigurator only (not the Pool), so
	// it is kept out of info.libs, which carries exactly the seven Pool libs.
	cfgLogicAddr, err := deploy("ConfiguratorLogic", contract.ConfiguratorLogicMetaData, nil, global, 0)
	if err != nil {
		return err
	}
	// FlashLoanLogic links BorrowLogic, so it must be linked after BorrowLogic's
	// address is known. Both are global, so the linked address (and thus
	// FlashLoanLogic's own address) is identical across deployers.
	if info.libs[phFlashLoanLogic], err = deploy("FlashLoanLogic", contract.FlashLoanLogicMetaData, map[string]common.Address{
		phBorrowLogic: info.libs[phBorrowLogic],
	}, global, 0); err != nil {
		return err
	}

	// Pool implementation links all seven logic libraries. It is deployer-specific
	// (its constructor binds the per-key addresses provider).
	if info.PoolImplAddr, err = deploy("Pool", contract.PoolMetaData, info.libs, perKey, 0, info.ProviderAddr); err != nil {
		return err
	}
	// PoolConfigurator implementation links ConfiguratorLogic only.
	if info.ConfiguratorImplAddr, err = deploy("PoolConfigurator", contract.PoolConfiguratorMetaData, map[string]common.Address{
		phConfiguratorLogic: cfgLogicAddr,
	}, perKey, 0); err != nil {
		return err
	}

	if info.DataProviderAddr, err = deploy("AaveProtocolDataProvider", contract.AaveProtocolDataProviderMetaData, nil, perKey, 0, info.ProviderAddr); err != nil {
		return err
	}

	// Shared interest rate strategy (constructor model in core-v3 1.19.3).
	if info.RateStrategyAddr, err = deploy("DefaultReserveInterestRateStrategy", contract.DefaultReserveInterestRateStrategyMetaData, nil, perKey, 0,
		info.ProviderAddr,
		rayFrac(80, 100), // optimalUsageRatio 0.80
		big.NewInt(0),    // baseVariableBorrowRate
		rayFrac(4, 100),  // variableRateSlope1 0.04
		rayFrac(75, 100), // variableRateSlope2 0.75
		rayFrac(2, 100),  // stableRateSlope1 0.02
		rayFrac(75, 100), // stableRateSlope2 0.75
		rayFrac(1, 100),  // baseStableRateOffset 0.01
		rayFrac(8, 100),  // stableRateExcessOffset 0.08
		rayFrac(20, 100), // optimalStableToTotalDebtRatio 0.20
	); err != nil {
		return err
	}

	// mock reserve tokens + price aggregators
	info.Tokens = make([]TokenInfo, 2)
	for i := 0; i < 2; i++ {
		name := fmt.Sprintf("Aave Mock Token %d", i)
		symbol := fmt.Sprintf("AMT%d", i)
		tokenAddr, err := deploy(fmt.Sprintf("token %d", i), contract.MintableTokenMetaData, nil, global, 0, name, symbol)
		if err != nil {
			return err
		}
		aggAddr, err := deploy(fmt.Sprintf("aggregator %d", i), contract.MockAggregatorMetaData, nil, perKey, uint32(i), big.NewInt(oraclePriceAnswer))
		if err != nil {
			return err
		}
		info.Tokens[i] = TokenInfo{Addr: tokenAddr, AggAddr: aggAddr, Name: name, Symbol: symbol}
	}

	// oracle (references the token + aggregator addresses)
	assets := []common.Address{info.Tokens[0].Addr, info.Tokens[1].Addr}
	sources := []common.Address{info.Tokens[0].AggAddr, info.Tokens[1].AggAddr}
	if info.OracleAddr, err = deploy("AaveOracle", contract.AaveOracleMetaData, nil, perKey, 0,
		info.ProviderAddr, assets, sources, common.Address{}, common.Address{}, big.NewInt(oracleBaseCurrencyUnit),
	); err != nil {
		return err
	}

	return s.sendBatch(ctx, deployerWallet, client, txs, "deploying aave implementations")
}

// wireMarket (phase 2) deploys the ACLManager and wires the addresses provider:
// it sets the ACL admin/manager, grants the deployer pool-admin, registers the
// oracle and data provider, and installs the Pool/PoolConfigurator
// implementations (which makes the provider deploy their proxies). The calls use
// fixed gas and run in nonce order within one batch so each executes against the
// state produced by the previous one. Skipped entirely on a re-run where the
// configurator proxy already exists.
func (s *Scenario) wireMarket(ctx context.Context, client *spamoor.Client, deployerWallet *spamoor.Wallet, feeCap, tipCap *big.Int, info *DeploymentInfo) error {
	providerCaller, err := contract.NewPoolAddressesProvider(info.ProviderAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind addresses provider: %w", err)
	}
	if cfg, err := providerCaller.GetPoolConfigurator(&bind.CallOpts{Context: ctx}); err == nil && cfg != (common.Address{}) {
		s.logger.Infof("aave market already wired, skipping")
		return nil
	}

	deployerAddr := deployerWallet.GetAddress()
	provider, err := contract.NewPoolAddressesProvider(info.ProviderAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind addresses provider: %w", err)
	}

	var txs []*types.Transaction
	boundCall := func(what string, gas uint64, build func(*bind.TransactOpts) (*types.Transaction, error)) error {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       gas,
			Value:     uint256.NewInt(0),
		}, build)
		if err != nil {
			return fmt.Errorf("could not build %s tx: %w", what, err)
		}
		txs = append(txs, tx)
		return nil
	}

	// ACL admin must be set before the ACLManager constructor reads it.
	if err := boundCall("setACLAdmin", gasSetter, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return provider.SetACLAdmin(opts, deployerAddr)
	}); err != nil {
		return err
	}

	// Deploy ACLManager (constructor reads provider.getACLAdmin(), set above and
	// executed earlier in this batch).
	aclAddr, aclTx, err := s.buildDeployTx(ctx, client, deployerWallet, feeCap, tipCap, contract.ACLManagerMetaData, nil, false, 0, info.ProviderAddr)
	if err != nil {
		return fmt.Errorf("could not deploy ACLManager: %w", err)
	}
	if aclTx != nil {
		txs = append(txs, aclTx)
	}
	info.ACLManagerAddr = aclAddr

	acl, err := contract.NewACLManager(aclAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind ACLManager: %w", err)
	}

	if err := boundCall("setACLManager", gasSetter, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return provider.SetACLManager(opts, aclAddr)
	}); err != nil {
		return err
	}
	if err := boundCall("addPoolAdmin", gasSetter, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return acl.AddPoolAdmin(opts, deployerAddr)
	}); err != nil {
		return err
	}
	if err := boundCall("setPriceOracle", gasSetter, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return provider.SetPriceOracle(opts, info.OracleAddr)
	}); err != nil {
		return err
	}
	if err := boundCall("setPoolDataProvider", gasSetter, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return provider.SetPoolDataProvider(opts, info.DataProviderAddr)
	}); err != nil {
		return err
	}
	if err := boundCall("setPoolImpl", gasProxyDeploy, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return provider.SetPoolImpl(opts, info.PoolImplAddr)
	}); err != nil {
		return err
	}
	if err := boundCall("setPoolConfiguratorImpl", gasProxyDeploy, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return provider.SetPoolConfiguratorImpl(opts, info.ConfiguratorImplAddr)
	}); err != nil {
		return err
	}

	return s.sendBatch(ctx, deployerWallet, client, txs, "wiring aave market")
}

// deployTokenImpls (phase 3) resolves the pool/configurator proxy addresses
// (created by phase 2) and deploys the aToken / debt token implementations,
// whose constructors call pool.ADDRESSES_PROVIDER() on the now-live pool proxy.
func (s *Scenario) deployTokenImpls(ctx context.Context, client *spamoor.Client, deployerWallet *spamoor.Wallet, feeCap, tipCap *big.Int, info *DeploymentInfo) error {
	provider, err := contract.NewPoolAddressesProvider(info.ProviderAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind addresses provider: %w", err)
	}
	callOpts := &bind.CallOpts{Context: ctx}
	if info.PoolAddr, err = provider.GetPool(callOpts); err != nil {
		return fmt.Errorf("could not read pool proxy: %w", err)
	}
	if info.ConfiguratorAddr, err = provider.GetPoolConfigurator(callOpts); err != nil {
		return fmt.Errorf("could not read configurator proxy: %w", err)
	}
	if info.PoolAddr == (common.Address{}) || info.ConfiguratorAddr == (common.Address{}) {
		return fmt.Errorf("pool or configurator proxy not deployed")
	}

	var txs []*types.Transaction
	// token implementations are deployer-specific (their constructor binds the
	// per-key pool proxy).
	deploy := func(what string, metadata *bind.MetaData, salt uint32) (common.Address, error) {
		addr, tx, err := s.buildDeployTx(ctx, client, deployerWallet, feeCap, tipCap, metadata, nil, false, salt, info.PoolAddr)
		if err != nil {
			return common.Address{}, fmt.Errorf("could not deploy %s: %w", what, err)
		}
		if tx != nil {
			txs = append(txs, tx)
		}
		return addr, nil
	}

	if info.ATokenImplAddr, err = deploy("AToken", contract.ATokenMetaData, 0); err != nil {
		return err
	}
	if info.VarDebtImplAddr, err = deploy("VariableDebtToken", contract.VariableDebtTokenMetaData, 0); err != nil {
		return err
	}
	if info.StableDebtImplAddr, err = deploy("StableDebtToken", contract.StableDebtTokenMetaData, 0); err != nil {
		return err
	}

	return s.sendBatch(ctx, deployerWallet, client, txs, "deploying token implementations")
}

// bindContracts binds the instances used on the hot path to a non-builder client
// for static calls and transaction building.
func (s *Scenario) bindContracts(info *DeploymentInfo) error {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithoutBuilder(),
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	var err error
	if info.Pool, err = contract.NewPool(info.PoolAddr, client.GetEthClient()); err != nil {
		return fmt.Errorf("could not bind pool: %w", err)
	}
	for i := range info.Tokens {
		if info.Tokens[i].Token, err = contract.NewMintableToken(info.Tokens[i].Addr, client.GetEthClient()); err != nil {
			return fmt.Errorf("could not bind token %d: %w", i, err)
		}
		if info.Tokens[i].Agg, err = contract.NewMockAggregator(info.Tokens[i].AggAddr, client.GetEthClient()); err != nil {
			return fmt.Errorf("could not bind aggregator %d: %w", i, err)
		}
	}
	return nil
}

// bindReserveTokens resolves and binds the aToken / variable debt token created
// for each reserve by initReserves. It must run after the reserves exist; the
// action engine reads these balances to gate supply/borrow/repay/withdraw and to
// find liquidatable positions.
func (s *Scenario) bindReserveTokens(ctx context.Context, info *DeploymentInfo) error {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithoutBuilder(),
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	callOpts := &bind.CallOpts{Context: ctx}
	for i := range info.Tokens {
		data, err := info.Pool.GetReserveData(callOpts, info.Tokens[i].Addr)
		if err != nil {
			return fmt.Errorf("could not read reserve data for token %d: %w", i, err)
		}
		if data.ATokenAddress == (common.Address{}) || data.VariableDebtTokenAddress == (common.Address{}) {
			return fmt.Errorf("reserve tokens for token %d not initialized", i)
		}
		info.Tokens[i].ATokenAddr = data.ATokenAddress
		info.Tokens[i].VarDebtAddr = data.VariableDebtTokenAddress
		if info.Tokens[i].AToken, err = contract.NewAToken(data.ATokenAddress, client.GetEthClient()); err != nil {
			return fmt.Errorf("could not bind aToken %d: %w", i, err)
		}
		if info.Tokens[i].VarDebt, err = contract.NewVariableDebtToken(data.VariableDebtTokenAddress, client.GetEthClient()); err != nil {
			return fmt.Errorf("could not bind variable debt token %d: %w", i, err)
		}
	}
	return nil
}

// initReserves (phase 4) registers both reserves on the pool via the
// configurator. This deploys and initializes the aToken/debt token proxies for
// each reserve. Skipped per-market when the first reserve already exists.
func (s *Scenario) initReserves(ctx context.Context, client *spamoor.Client, deployerWallet *spamoor.Wallet, feeCap, tipCap *big.Int, info *DeploymentInfo) error {
	if data, err := info.Pool.GetReserveData(&bind.CallOpts{Context: ctx}, info.Tokens[0].Addr); err == nil && data.ATokenAddress != (common.Address{}) {
		s.logger.Infof("aave reserves already initialized, skipping")
		return nil
	}

	deployerAddr := deployerWallet.GetAddress()
	configurator, err := contract.NewPoolConfigurator(info.ConfiguratorAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind configurator: %w", err)
	}

	inputs := make([]contract.ConfiguratorInputTypesInitReserveInput, 0, len(info.Tokens))
	for _, t := range info.Tokens {
		inputs = append(inputs, contract.ConfiguratorInputTypesInitReserveInput{
			ATokenImpl:                  info.ATokenImplAddr,
			StableDebtTokenImpl:         info.StableDebtImplAddr,
			VariableDebtTokenImpl:       info.VarDebtImplAddr,
			UnderlyingAssetDecimals:     reserveUnderlyingDecimals,
			InterestRateStrategyAddress: info.RateStrategyAddr,
			UnderlyingAsset:             t.Addr,
			Treasury:                    deployerAddr,
			IncentivesController:        common.Address{},
			ATokenName:                  fmt.Sprintf("Aave %s", t.Symbol),
			ATokenSymbol:                fmt.Sprintf("a%s", t.Symbol),
			VariableDebtTokenName:       fmt.Sprintf("Aave Variable Debt %s", t.Symbol),
			VariableDebtTokenSymbol:     fmt.Sprintf("variableDebt%s", t.Symbol),
			StableDebtTokenName:         fmt.Sprintf("Aave Stable Debt %s", t.Symbol),
			StableDebtTokenSymbol:       fmt.Sprintf("stableDebt%s", t.Symbol),
			Params:                      []byte{},
		})
	}

	// All dependencies are mined, so gas estimation works here.
	tx, err := deployerWallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return configurator.InitReserves(opts, inputs)
	})
	if err != nil {
		return fmt.Errorf("could not build initReserves tx: %w", err)
	}

	return s.sendBatch(ctx, deployerWallet, client, []*types.Transaction{tx}, "initializing reserves")
}

// configureReserves (phase 5) enables each reserve as collateral, enables
// variable borrowing, sets the reserve factor and activates the reserve.
func (s *Scenario) configureReserves(ctx context.Context, client *spamoor.Client, deployerWallet *spamoor.Wallet, feeCap, tipCap *big.Int, info *DeploymentInfo) error {
	configurator, err := contract.NewPoolConfigurator(info.ConfiguratorAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not bind configurator: %w", err)
	}

	var txs []*types.Transaction
	boundCall := func(what string, build func(*bind.TransactOpts) (*types.Transaction, error)) error {
		tx, err := deployerWallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, build)
		if err != nil {
			return fmt.Errorf("could not build %s tx: %w", what, err)
		}
		txs = append(txs, tx)
		return nil
	}

	for _, t := range info.Tokens {
		asset := t.Addr
		if err := boundCall("configureReserveAsCollateral", func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return configurator.ConfigureReserveAsCollateral(opts, asset, big.NewInt(reserveLTV), big.NewInt(reserveLiquidationThresh), big.NewInt(reserveLiquidationBonus))
		}); err != nil {
			return err
		}
		if err := boundCall("setReserveBorrowing", func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return configurator.SetReserveBorrowing(opts, asset, true)
		}); err != nil {
			return err
		}
		if err := boundCall("setReserveFactor", func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return configurator.SetReserveFactor(opts, asset, big.NewInt(reserveFactor))
		}); err != nil {
			return err
		}
		if err := boundCall("setReserveActive", func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return configurator.SetReserveActive(opts, asset, true)
		}); err != nil {
			return err
		}
	}

	return s.sendBatch(ctx, deployerWallet, client, txs, "configuring reserves")
}

// seedLiquidity (phase 6) has the deployer mint and supply a large amount of
// every reserve so child wallets can borrow immediately at the cold start. Fixed
// gas is used (and mint/approve/supply share one batch) because supply cannot be
// gas-estimated before the approve in front of it is mined.
func (s *Scenario) seedLiquidity(ctx context.Context, client *spamoor.Client, deployerWallet *spamoor.Wallet, feeCap, tipCap *big.Int, info *DeploymentInfo) error {
	deployerAddr := deployerWallet.GetAddress()

	// resume guard: if the first reserve already holds seeded liquidity, skip
	if data, err := info.Pool.GetReserveData(&bind.CallOpts{Context: ctx}, info.Tokens[0].Addr); err == nil && data.ATokenAddress != (common.Address{}) {
		if bal, err := info.Tokens[0].Token.BalanceOf(&bind.CallOpts{Context: ctx}, data.ATokenAddress); err == nil && bal.Sign() > 0 {
			s.logger.Infof("aave reserves already seeded, skipping")
			return nil
		}
	}

	var txs []*types.Transaction
	boundCall := func(what string, gas uint64, build func(*bind.TransactOpts) (*types.Transaction, error)) error {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       gas,
			Value:     uint256.NewInt(0),
		}, build)
		if err != nil {
			return fmt.Errorf("could not build %s tx: %w", what, err)
		}
		txs = append(txs, tx)
		return nil
	}

	for i := range info.Tokens {
		token := info.Tokens[i].Token
		asset := info.Tokens[i].Addr
		if err := boundCall("seed mint", gasMint, func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return token.Mint(opts, deployerAddr, s.seedAmount)
		}); err != nil {
			return err
		}
		if err := boundCall("seed approve", gasApprove, func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return token.Approve(opts, info.PoolAddr, maxUint256)
		}); err != nil {
			return err
		}
		if err := boundCall("seed supply", s.actionGasLimit(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return info.Pool.Supply(opts, asset, s.seedAmount, deployerAddr, 0)
		}); err != nil {
			return err
		}
	}

	return s.sendBatch(ctx, deployerWallet, client, txs, "seeding reserve liquidity")
}

// FundAndApproveWallets mints the configured starting balance of every reserve
// token to each child wallet and sets an unlimited allowance to the pool so the
// supply/repay calls can pull funds. It mirrors the curve-swaps funding pattern:
// per-wallet RPC fan-out is bounded by the number of healthy clients and the
// transactions use fixed gas.
func (s *Scenario) FundAndApproveWallets(ctx context.Context, info *DeploymentInfo) error {
	s.logger.Infof("funding and approving wallets...")

	wallets := s.walletPool.GetAllWallets()
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return fmt.Errorf("could not get tx fee: %w", err)
	}

	concurrency := len(s.walletPool.GetClientPool().GetAllGoodClients())
	concurrency = min(max(concurrency, 1), 50)

	var (
		setupTxs     []*types.Transaction
		setupWallets []*spamoor.Wallet
		mu           sync.Mutex
		wg           sync.WaitGroup
	)
	sem := make(chan struct{}, concurrency)

	buildTx := func(wallet *spamoor.Wallet, gas uint64, build func(*bind.TransactOpts) (*types.Transaction, error)) {
		tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       gas,
			Value:     uint256.NewInt(0),
		}, build)
		if err != nil {
			s.logger.Errorf("could not build setup tx for %v: %v", wallet.GetAddress(), err)
			return
		}
		mu.Lock()
		setupTxs = append(setupTxs, tx)
		setupWallets = append(setupWallets, wallet)
		mu.Unlock()
	}

	for idx, wallet := range wallets {
		if ctx.Err() != nil {
			break
		}
		wg.Add(1)
		go func(idx int, wallet *spamoor.Wallet) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return
			}

			rclient := s.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, idx),
				spamoor.WithClientGroup(s.options.ClientGroup),
				spamoor.WithoutBuilder(),
			)
			if rclient == nil {
				rclient = client
			}
			callOpts := &bind.CallOpts{Context: ctx}

			for i := range info.Tokens {
				token, err := contract.NewMintableToken(info.Tokens[i].Addr, rclient.GetEthClient())
				if err != nil {
					s.logger.Errorf("could not bind token %v: %v", info.Tokens[i].Addr, err)
					continue
				}

				if balance, err := token.BalanceOf(callOpts, wallet.GetAddress()); err != nil {
					s.logger.Errorf("could not read token balance for %v: %v", wallet.GetAddress(), err)
				} else if balance.Cmp(s.walletFunding) < 0 {
					buildTx(wallet, gasMint, func(opts *bind.TransactOpts) (*types.Transaction, error) {
						return token.Mint(opts, wallet.GetAddress(), s.walletFunding)
					})
				}

				allowance, err := token.Allowance(callOpts, wallet.GetAddress(), info.PoolAddr)
				if err != nil {
					s.logger.Errorf("could not check allowance for %v: %v", wallet.GetAddress(), err)
					continue
				}
				if allowance.Cmp(maxUint256) >= 0 {
					continue
				}
				buildTx(wallet, gasApprove, func(opts *bind.TransactOpts) (*types.Transaction, error) {
					return token.Approve(opts, info.PoolAddr, maxUint256)
				})
			}
		}(idx, wallet)
	}
	wg.Wait()

	if ctx.Err() != nil {
		return ctx.Err()
	}
	if len(setupTxs) == 0 {
		s.logger.Infof("no funding/approval transactions needed")
		return nil
	}

	s.logger.Infof("sending %d funding/approval transactions...", len(setupTxs))
	for i, tx := range setupTxs {
		txClient := s.walletPool.GetClient(
			spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, i),
			spamoor.WithClientGroup(s.options.ClientGroup),
		)
		if txClient == nil {
			txClient = client
		}

		wg.Add(1)
		go func(tx *types.Transaction, client *spamoor.Client, wallet *spamoor.Wallet) {
			s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: s.options.ClientGroup,
				Rebroadcast: true,
				OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					if err != nil {
						s.logger.Errorf("funding/approval tx failed: %v", err)
					}
					wg.Done()
				},
			})
		}(tx, txClient, setupWallets[i])
	}
	wg.Wait()
	s.logger.Infof("all funding/approval transactions sent")

	return nil
}

// sendBatch submits a batch of transactions from a single wallet and waits for
// them to confirm, logging progress. It is a no-op for an empty batch.
func (s *Scenario) sendBatch(ctx context.Context, wallet *spamoor.Wallet, client *spamoor.Client, txs []*types.Transaction, action string) error {
	if len(txs) == 0 {
		return nil
	}

	_, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, wallet, txs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: s.options.DeployClientGroup,
		},
		MaxRetries:   3,
		PendingLimit: 10,
		LogFn: func(confirmedCount int, totalCount int) {
			s.logger.Infof("%s... (%v/%v)", action, confirmedCount, totalCount)
		},
		LogInterval: 10,
	})
	if err != nil {
		return fmt.Errorf("could not %s: %w", action, err)
	}
	s.logger.Infof("%s complete. (%v/%v)", action, len(txs), len(txs))
	return nil
}
