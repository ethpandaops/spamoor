#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd $SCRIPT_DIR

# Mock underlying assets and price source (own minimal contracts).
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 200" MintableToken
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 200" MockAggregator

# Aave V3 core contracts (canonical precompiled artifacts from @aave/core-v3).
# Using the published creation bytecode keeps the deployed contracts identical to
# mainnet Aave V3. The package is pinned to 1.19.3 so the library link placeholders
# (hardcoded in deployment.go) stay stable. Pool, PoolConfigurator and
# FlashLoanLogic ship with unlinked __$...$__ library placeholders that
# deployment.go resolves at deploy time.
AAVE="https://unpkg.com/@aave/core-v3@1.19.3/artifacts/contracts"

# Logic libraries (linked into Pool / PoolConfigurator / FlashLoanLogic).
gen_from_artifact "$AAVE/protocol/libraries/logic/SupplyLogic.sol/SupplyLogic.json" SupplyLogic
gen_from_artifact "$AAVE/protocol/libraries/logic/BorrowLogic.sol/BorrowLogic.json" BorrowLogic
gen_from_artifact "$AAVE/protocol/libraries/logic/LiquidationLogic.sol/LiquidationLogic.json" LiquidationLogic
gen_from_artifact "$AAVE/protocol/libraries/logic/EModeLogic.sol/EModeLogic.json" EModeLogic
gen_from_artifact "$AAVE/protocol/libraries/logic/BridgeLogic.sol/BridgeLogic.json" BridgeLogic
gen_from_artifact "$AAVE/protocol/libraries/logic/FlashLoanLogic.sol/FlashLoanLogic.json" FlashLoanLogic
gen_from_artifact "$AAVE/protocol/libraries/logic/PoolLogic.sol/PoolLogic.json" PoolLogic
gen_from_artifact "$AAVE/protocol/libraries/logic/ConfiguratorLogic.sol/ConfiguratorLogic.json" ConfiguratorLogic

# Configuration + core protocol.
gen_from_artifact "$AAVE/protocol/configuration/PoolAddressesProvider.sol/PoolAddressesProvider.json" PoolAddressesProvider
gen_from_artifact "$AAVE/protocol/configuration/ACLManager.sol/ACLManager.json" ACLManager
gen_from_artifact "$AAVE/protocol/pool/Pool.sol/Pool.json" Pool
gen_from_artifact "$AAVE/protocol/pool/PoolConfigurator.sol/PoolConfigurator.json" PoolConfigurator
gen_from_artifact "$AAVE/protocol/pool/DefaultReserveInterestRateStrategy.sol/DefaultReserveInterestRateStrategy.json" DefaultReserveInterestRateStrategy

# Misc helpers + oracle.
gen_from_artifact "$AAVE/misc/AaveProtocolDataProvider.sol/AaveProtocolDataProvider.json" AaveProtocolDataProvider
gen_from_artifact "$AAVE/misc/AaveOracle.sol/AaveOracle.json" AaveOracle

# Tokenization (impls reused across reserves; proxies deployed by initReserves).
gen_from_artifact "$AAVE/protocol/tokenization/AToken.sol/AToken.json" AToken
gen_from_artifact "$AAVE/protocol/tokenization/VariableDebtToken.sol/VariableDebtToken.json" VariableDebtToken
gen_from_artifact "$AAVE/protocol/tokenization/StableDebtToken.sol/StableDebtToken.json" StableDebtToken
