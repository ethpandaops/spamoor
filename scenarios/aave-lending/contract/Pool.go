// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DataTypesEModeCategory is an auto generated low-level Go binding around an user-defined struct.
type DataTypesEModeCategory struct {
	Ltv                  uint16
	LiquidationThreshold uint16
	LiquidationBonus     uint16
	PriceSource          common.Address
	Label                string
}

// DataTypesReserveConfigurationMap is an auto generated low-level Go binding around an user-defined struct.
type DataTypesReserveConfigurationMap struct {
	Data *big.Int
}

// DataTypesReserveData is an auto generated low-level Go binding around an user-defined struct.
type DataTypesReserveData struct {
	Configuration               DataTypesReserveConfigurationMap
	LiquidityIndex              *big.Int
	CurrentLiquidityRate        *big.Int
	VariableBorrowIndex         *big.Int
	CurrentVariableBorrowRate   *big.Int
	CurrentStableBorrowRate     *big.Int
	LastUpdateTimestamp         *big.Int
	Id                          uint16
	ATokenAddress               common.Address
	StableDebtTokenAddress      common.Address
	VariableDebtTokenAddress    common.Address
	InterestRateStrategyAddress common.Address
	AccruedToTreasury           *big.Int
	Unbacked                    *big.Int
	IsolationModeTotalDebt      *big.Int
}

// DataTypesUserConfigurationMap is an auto generated low-level Go binding around an user-defined struct.
type DataTypesUserConfigurationMap struct {
	Data *big.Int
}

// PoolMetaData contains all meta data concerning the Pool contract.
var PoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"provider\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"backer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"BackUnbacked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumDataTypes.InterestRateMode\",\"name\":\"interestRateMode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"Borrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumDataTypes.InterestRateMode\",\"name\":\"interestRateMode\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"FlashLoan\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"}],\"name\":\"IsolationModeTotalDebtUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"collateralAsset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"debtAsset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"debtToCover\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidatedCollateralAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"receiveAToken\",\"type\":\"bool\"}],\"name\":\"LiquidationCall\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"MintUnbacked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountMinted\",\"type\":\"uint256\"}],\"name\":\"MintedToTreasury\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"RebalanceStableBorrowRate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"repayer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"useATokens\",\"type\":\"bool\"}],\"name\":\"Repay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidityRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stableBorrowRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"variableBorrowRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidityIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"variableBorrowIndex\",\"type\":\"uint256\"}],\"name\":\"ReserveDataUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"ReserveUsedAsCollateralDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"ReserveUsedAsCollateralEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"Supply\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumDataTypes.InterestRateMode\",\"name\":\"interestRateMode\",\"type\":\"uint8\"}],\"name\":\"SwapBorrowRateMode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"categoryId\",\"type\":\"uint8\"}],\"name\":\"UserEModeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADDRESSES_PROVIDER\",\"outputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BRIDGE_PROTOCOL_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FLASHLOAN_PREMIUM_TOTAL\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FLASHLOAN_PREMIUM_TO_PROTOCOL\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_NUMBER_RESERVES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_STABLE_RATE_BORROW_SIZE_PERCENT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POOL_REVISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"backUnbacked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"}],\"name\":\"borrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"id\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint16\",\"name\":\"ltv\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"liquidationThreshold\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"liquidationBonus\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceSource\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"internalType\":\"structDataTypes.EModeCategory\",\"name\":\"category\",\"type\":\"tuple\"}],\"name\":\"configureEModeCategory\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"dropReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balanceFromBefore\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balanceToBefore\",\"type\":\"uint256\"}],\"name\":\"finalizeTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiverAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"interestRateModes\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"flashLoan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiverAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"},{\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"flashLoanSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getConfiguration\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.ReserveConfigurationMap\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"id\",\"type\":\"uint8\"}],\"name\":\"getEModeCategoryData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint16\",\"name\":\"ltv\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"liquidationThreshold\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"liquidationBonus\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"priceSource\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"internalType\":\"structDataTypes.EModeCategory\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"id\",\"type\":\"uint16\"}],\"name\":\"getReserveAddressById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getReserveData\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.ReserveConfigurationMap\",\"name\":\"configuration\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"liquidityIndex\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"currentLiquidityRate\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"variableBorrowIndex\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"currentVariableBorrowRate\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"currentStableBorrowRate\",\"type\":\"uint128\"},{\"internalType\":\"uint40\",\"name\":\"lastUpdateTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"id\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"aTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stableDebtTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"variableDebtTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategyAddress\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"accruedToTreasury\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"unbacked\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"isolationModeTotalDebt\",\"type\":\"uint128\"}],\"internalType\":\"structDataTypes.ReserveData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getReserveNormalizedIncome\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getReserveNormalizedVariableDebt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReservesList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserAccountData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalCollateralBase\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebtBase\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"availableBorrowsBase\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentLiquidationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ltv\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"healthFactor\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserConfiguration\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.UserConfigurationMap\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserEMode\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"aTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stableDebtAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"variableDebtAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategyAddress\",\"type\":\"address\"}],\"name\":\"initReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"provider\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collateralAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"debtAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"debtToCover\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"receiveAToken\",\"type\":\"bool\"}],\"name\":\"liquidationCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"}],\"name\":\"mintToTreasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"mintUnbacked\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"rebalanceStableBorrowRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"}],\"name\":\"repay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"}],\"name\":\"repayWithATokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"permitV\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"permitR\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"permitS\",\"type\":\"bytes32\"}],\"name\":\"repayWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"rescueTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"resetIsolationModeTotalDebt\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.ReserveConfigurationMap\",\"name\":\"configuration\",\"type\":\"tuple\"}],\"name\":\"setConfiguration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rateStrategyAddress\",\"type\":\"address\"}],\"name\":\"setReserveInterestRateStrategyAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"categoryId\",\"type\":\"uint8\"}],\"name\":\"setUserEMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useAsCollateral\",\"type\":\"bool\"}],\"name\":\"setUserUseReserveAsCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"supply\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"permitV\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"permitR\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"permitS\",\"type\":\"bytes32\"}],\"name\":\"supplyWithPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"}],\"name\":\"swapBorrowRateMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"protocolFee\",\"type\":\"uint256\"}],\"name\":\"updateBridgeProtocolFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"flashLoanPremiumTotal\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"flashLoanPremiumToProtocol\",\"type\":\"uint128\"}],\"name\":\"updateFlashloanPremiums\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052600080553480156200001557600080fd5b50604051620055a3380380620055a383398101604081905262000038916200004a565b6001600160a01b03166080526200007c565b6000602082840312156200005d57600080fd5b81516001600160a01b03811681146200007557600080fd5b9392505050565b6080516154a9620000fa6000396000818161035601528181610a4e01528181610b40015281816110430152818161166601528181611a0c0152818161211d015281816121ee015281816124410152818161273c0152818161299b015281816130120152818161360f015281816137b6015261394301526154a96000f3fe608060405234801561001057600080fd5b50600436106103095760003560e01c80637a708e921161019c578063d15e0053116100ee578063e82fec2f11610097578063ee3e210b11610071578063ee3e210b1461096d578063f51e435b14610980578063f8119d511461099357600080fd5b8063e82fec2f14610922578063e8eda9df146106da578063eddf1b791461093457600080fd5b8063d5ed3933116100c8578063d5ed3933146108e9578063d65dc7a1146108fc578063e43e88a11461090f57600080fd5b8063d15e0053146108ae578063d1946dbc146108c1578063d579ea7d146108d657600080fd5b8063bcb6e52211610150578063c4d66de81161012a578063c4d66de814610875578063cd11238214610888578063cea9d26f1461089b57600080fd5b8063bcb6e522146107d3578063bf92857c146107e6578063c44b11f71461082657600080fd5b80639cd19996116101815780639cd199961461079a578063a415bcad146107ad578063ab9c4b5d146107c057600080fd5b80637a708e921461077457806394ba89a21461078757600080fd5b8063386497fd11610260578063617ba0371161020957806369a933a5116101e357806369a933a5146107135780636a99c036146107265780636c6f6ae11461075457600080fd5b8063617ba037146106da57806363c9b860146106ed57806369328dec1461070057600080fd5b8063527517971161023a578063527517971461067a578063573ade81146106b45780635a3b74b9146106c757600080fd5b8063386497fd146105f657806342b0b77c146106095780634417a5831461061c57600080fd5b80631d2118f9116102c25780632dad97d41161029c5780632dad97d4146104025780633036b4391461041557806335ea6a751461042857600080fd5b80631d2118f9146103d4578063272d9072146103e757806328530a47146103ef57600080fd5b806302c205f0116102f357806302c205f01461033e5780630542975c14610351578063074b2e431461039d57600080fd5b8062a718a91461030e5780630148170e14610323575b600080fd5b61032161031c366004613e1b565b6109a2565b005b61032b600181565b6040519081526020015b60405180910390f35b61032161034c366004613ea6565b610c1d565b6103787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610335565b603a546fffffffffffffffffffffffffffffffff165b6040516fffffffffffffffffffffffffffffffff9091168152602001610335565b6103216103e2366004613f25565b610dcd565b60395461032b565b6103216103fd366004613f5e565b610fbb565b61032b610410366004613f79565b61119a565b610321610423366004613fae565b6112de565b6105e9610436366004613fc7565b604080516102008101825260006101e08201818152825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101829052610100810182905261012081018290526101408101829052610160810182905261018081018290526101a081018290526101c08101919091525073ffffffffffffffffffffffffffffffffffffffff90811660009081526034602090815260409182902082516102008101845281546101e08201908152815260018201546fffffffffffffffffffffffffffffffff80821694830194909452700100000000000000000000000000000000908190048416948201949094526002820154808416606083015284900483166080820152600382015480841660a083015284810464ffffffffff1660c08301527501000000000000000000000000000000000000000000900461ffff1660e0820152600482015485166101008201526005820154851661012082015260068201548516610140820152600782015490941661016085015260088101548083166101808601529290920481166101a0840152600990910154166101c082015290565b6040516103359190613fe4565b61032b610604366004613fc7565b6112eb565b6103216106173660046141aa565b61131f565b61066b61062a366004613fc7565b604080516020808201835260009182905273ffffffffffffffffffffffffffffffffffffffff93909316815260358352819020815192830190915254815290565b60405190518152602001610335565b61037861068836600461422c565b61ffff1660009081526036602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b61032b6106c2366004614247565b611499565b6103216106d5366004614291565b6115f2565b6103216106e83660046142bf565b6117c7565b6103216106fb366004613fc7565b6118ca565b61032b61070e366004614310565b611946565b6103216107213660046142bf565b611b65565b603a5470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff166103b3565b610767610762366004613f5e565b611c12565b60405161033591906143bd565b610321610782366004614420565b611d4c565b610321610795366004614483565b611ed8565b6103216107a83660046144f4565b611f59565b6103216107bb366004614536565b611fae565b6103216107ce366004614575565b612294565b6103216107e136600461468f565b61264d565b6107f96107f4366004613fc7565b612684565b604080519687526020870195909552938501929092526060840152608083015260a082015260c001610335565b61066b610834366004613fc7565b604080516020808201835260009182905273ffffffffffffffffffffffffffffffffffffffff93909316815260348352819020815192830190915254815290565b610321610883366004613fc7565b6128b3565b610321610896366004613f25565b612ab7565b6103216108a93660046146c2565b612b40565b61032b6108bc366004613fc7565b612bed565b6108c9612c1b565b6040516103359190614703565b6103216108e4366004614804565b612d57565b6103216108f736600461493c565b612ec3565b61032b61090a366004613f79565b61314a565b61032161091d366004613fc7565b6131ea565b603b5467ffffffffffffffff1661032b565b61032b610942366004613fc7565b73ffffffffffffffffffffffffffffffffffffffff1660009081526038602052604090205460ff1690565b61032b61097b3660046149a1565b61325f565b61032161098e3660046149e7565b61343a565b60405160808152602001610335565b73__$f598c634f2d943205ac23f707b80075cbb$__6383c1087d6034603660356037604051806101200160405280603b60089054906101000a900461ffff1661ffff1681526020018981526020018c73ffffffffffffffffffffffffffffffffffffffff1681526020018b73ffffffffffffffffffffffffffffffffffffffff1681526020018a73ffffffffffffffffffffffffffffffffffffffff16815260200188151581526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fca513a86040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ab7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610adb9190614a46565b73ffffffffffffffffffffffffffffffffffffffff90811682528b81166000908152603860209081526040918290205460ff168185015281517f5eb88d3d000000000000000000000000000000000000000000000000000000008152825192909401937f000000000000000000000000000000000000000000000000000000000000000090931692635eb88d3d92600480830193928290030181865afa158015610b89573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bad9190614a46565b73ffffffffffffffffffffffffffffffffffffffff168152506040518663ffffffff1660e01b8152600401610be6959493929190614a63565b60006040518083038186803b158015610bfe57600080fd5b505af4158015610c12573d6000803e3d6000fd5b505050505050505050565b6040517fd505accf000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018890526064810185905260ff8416608482015260a4810183905260c4810182905273ffffffffffffffffffffffffffffffffffffffff89169063d505accf9060e401600060405180830381600087803b158015610caf57600080fd5b505af1158015610cc3573d6000803e3d6000fd5b5050505073ffffffffffffffffffffffffffffffffffffffff86811660008181526035602090815260409182902082516080810184528d861681529182018c815282840194855261ffff8b81166060850190815294517f1913f16100000000000000000000000000000000000000000000000000000000815260346004820152603660248201526044810193909352925186166064830152516084820152925190931660a48301525190911660c482015273__$db79717e66442ee197e8271d032a066e34$__90631913f1619060e40160006040518083038186803b158015610dab57600080fd5b505af4158015610dbf573d6000803e3d6000fd5b505050505050505050505050565b610dd56135f6565b60408051808201909152600281527f3737000000000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff8316610e60576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b60405180910390fd5b5073ffffffffffffffffffffffffffffffffffffffff82166000908152603460205260409020600301547501000000000000000000000000000000000000000000900461ffff16151580610ef657506000805260366020527f4cb2b152c1b54ce671907a93c300fd5aa72383a9d4ec19a81e3333632ae92e005473ffffffffffffffffffffffffffffffffffffffff8381169116145b6040518060400160405280600281526020017f383200000000000000000000000000000000000000000000000000000000000081525090610f64576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b5073ffffffffffffffffffffffffffffffffffffffff918216600090815260346020526040902060070180547fffffffffffffffffffffffff00000000000000000000000000000000000000001691909216179055565b73__$e4b9550ff526a295e1233dea02821b9004$__635d5dc3136034603660376038603560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518060600160405280603b60089054906101000a900461ffff1661ffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fca513a86040518163ffffffff1660e01b8152600401602060405180830381865afa1580156110ac573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110d09190614a46565b73ffffffffffffffffffffffffffffffffffffffff1681526020018960ff168152506040518763ffffffff1660e01b81526004016111679695949392919095865260208087019590955260408087019490945260608601929092526080850152805160a08501529182015173ffffffffffffffffffffffffffffffffffffffff1660c0840152015160ff1660e08201526101000190565b60006040518083038186803b15801561117f57600080fd5b505af4158015611193573d6000803e3d6000fd5b5050505050565b600073__$c3724b8d563dc83a94e797176cddecb3b9$__6340e95de660346036603560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518060a001604052808a73ffffffffffffffffffffffffffffffffffffffff16815260200189815260200188600281111561123857611238614b57565b600281111561124957611249614b57565b81523360208201526001604091820152517fffffffff0000000000000000000000000000000000000000000000000000000060e087901b1681526112939493929190600401614bc1565b602060405180830381865af41580156112b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112d49190614c34565b90505b9392505050565b6112e66135f6565b603955565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260346020526040812061131990613724565b92915050565b60006040518060e001604052808873ffffffffffffffffffffffffffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff16815260200186815260200185858080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525093855250505061ffff8516602080840191909152603a546fffffffffffffffffffffffffffffffff70010000000000000000000000000000000082048116604080870191909152911660609094019390935273ffffffffffffffffffffffffffffffffffffffff8a1682526034905281902090517fa1fe0e8d00000000000000000000000000000000000000000000000000000000815291925073__$d5ddd09ae98762b8929dd85e54b218e259$__9163a1fe0e8d91611460918590600401614c4d565b60006040518083038186803b15801561147857600080fd5b505af415801561148c573d6000803e3d6000fd5b5050505050505050505050565b600073__$c3724b8d563dc83a94e797176cddecb3b9$__6340e95de660346036603560008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518060a001604052808b73ffffffffffffffffffffffffffffffffffffffff1681526020018a815260200189600281111561153757611537614b57565b600281111561154857611548614b57565b815273ffffffffffffffffffffffffffffffffffffffff891660208201526000604091820152517fffffffff0000000000000000000000000000000000000000000000000000000060e087901b1681526115a89493929190600401614bc1565b602060405180830381865af41580156115c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115e99190614c34565b95945050505050565b73__$db79717e66442ee197e8271d032a066e34$__63bf697a26603460366037603560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208787603b60089054906101000a900461ffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fca513a86040518163ffffffff1660e01b8152600401602060405180830381865afa1580156116cf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116f39190614a46565b336000908152603860205260409081902054905160e08b901b7fffffffff00000000000000000000000000000000000000000000000000000000168152600481019990995260248901979097526044880195909552606487019390935273ffffffffffffffffffffffffffffffffffffffff9182166084870152151560a486015261ffff90911660c48501521660e483015260ff16610104820152610124015b60006040518083038186803b1580156117ab57600080fd5b505af41580156117bf573d6000803e3d6000fd5b505050505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000818152603560209081526040918290208251608081018452898616815291820188815282840194855261ffff8781166060850190815294517f1913f16100000000000000000000000000000000000000000000000000000000815260346004820152603660248201526044810193909352925186166064830152516084820152925190931660a48301525190911660c482015273__$db79717e66442ee197e8271d032a066e34$__90631913f1619060e4015b60006040518083038186803b1580156118ac57600080fd5b505af41580156118c0573d6000803e3d6000fd5b5050505050505050565b6118d26135f6565b6040517f9cf57023000000000000000000000000000000000000000000000000000000008152603460048201526036602482015273ffffffffffffffffffffffffffffffffffffffff8216604482015273__$563c746fa3df0f1858d85f6ef4258864be$__90639cf5702390606401611167565b600073__$db79717e66442ee197e8271d032a066e34$__63186dea44603460366037603560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518060c001604052808b73ffffffffffffffffffffffffffffffffffffffff1681526020018a81526020018973ffffffffffffffffffffffffffffffffffffffff168152602001603b60089054906101000a900461ffff1661ffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fca513a86040518163ffffffff1660e01b8152600401602060405180830381865afa158015611a75573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a999190614a46565b73ffffffffffffffffffffffffffffffffffffffff9081168252336000908152603860209081526040918290205460ff90811694820194909452815160e08b901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260048101999099526024890197909752604488019590955260648701939093528151831660848701529381015160a486015291820151811660c4850152606082015160e485015260808201511661010484015260a001511661012482015261014401611293565b611b6d6137b4565b73ffffffffffffffffffffffffffffffffffffffff8281166000818152603560205260409081902090517f0413c86f0000000000000000000000000000000000000000000000000000000081526034600482015260366024820152604481019190915291861660648301526084820185905260a482015261ffff821660c482015273__$b06080f092f400a43662c3f835a4d9baa8$__90630413c86f9060e401611894565b6040805160a081018252600080825260208201819052918101829052606080820192909252608081019190915260ff8216600090815260376020908152604091829020825160a081018452815461ffff8082168352620100008204811694830194909452640100000000810490931693810193909352660100000000000090910473ffffffffffffffffffffffffffffffffffffffff166060830152600181018054608084019190611cc390614cd8565b80601f0160208091040260200160405190810160405280929190818152602001828054611cef90614cd8565b8015611d3c5780601f10611d1157610100808354040283529160200191611d3c565b820191906000526020600020905b815481529060010190602001808311611d1f57829003601f168201915b5050505050815250509050919050565b611d546135f6565b73__$563c746fa3df0f1858d85f6ef4258864be$__6369fc1bdf603460366040518060e001604052808a73ffffffffffffffffffffffffffffffffffffffff1681526020018973ffffffffffffffffffffffffffffffffffffffff1681526020018873ffffffffffffffffffffffffffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff168152602001603b60089054906101000a900461ffff1661ffff168152602001611e2b608090565b61ffff168152506040518463ffffffff1660e01b8152600401611e5093929190614d26565b602060405180830381865af4158015611e6d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e919190614db6565b1561119357603b805468010000000000000000900461ffff16906008611eb683614e02565b91906101000a81548161ffff021916908361ffff160217905550505050505050565b73ffffffffffffffffffffffffffffffffffffffff82166000908152603460209081526040808320338452603590925290912073__$c3724b8d563dc83a94e797176cddecb3b9$__9163eac4d7039185856002811115611f3a57611f3a614b57565b6040518563ffffffff1660e01b81526004016117939493929190614e24565b6040517f48c2ca8c00000000000000000000000000000000000000000000000000000000815273__$563c746fa3df0f1858d85f6ef4258864be$__906348c2ca8c906117939060349086908690600401614e5b565b73__$c3724b8d563dc83a94e797176cddecb3b9$__631e6473f9603460366037603560008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518061018001604052808c73ffffffffffffffffffffffffffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018873ffffffffffffffffffffffffffffffffffffffff1681526020018b81526020018a600281111561208557612085614b57565b600281111561209657612096614b57565b815261ffff808b166020808401919091526001604080850191909152603b5467ffffffffffffffff81166060860152680100000000000000009004909216608084015281517ffca513a8000000000000000000000000000000000000000000000000000000008152915160a09093019273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169263fca513a89260048083019391928290030181865afa158015612165573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121899190614a46565b73ffffffffffffffffffffffffffffffffffffffff90811682528981166000908152603860209081526040918290205460ff168185015281517f5eb88d3d000000000000000000000000000000000000000000000000000000008152825192909401937f000000000000000000000000000000000000000000000000000000000000000090931692635eb88d3d92600480830193928290030181865afa158015612237573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061225b9190614a46565b73ffffffffffffffffffffffffffffffffffffffff168152506040518663ffffffff1660e01b8152600401610be6959493929190614ec0565b6000604051806101c001604052808d73ffffffffffffffffffffffffffffffffffffffff1681526020018c8c808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505050908252506040805160208c810282810182019093528c82529283019290918d918d9182918501908490808284376000920191909152505050908252506040805160208a810282810182019093528a82529283019290918b918b91829185019084908082843760009201919091525050509082525073ffffffffffffffffffffffffffffffffffffffff871660208083019190915260408051601f88018390048302810183018252878152920191908790879081908401838280828437600092018290525093855250505061ffff808616602080850191909152603a546fffffffffffffffffffffffffffffffff7001000000000000000000000000000000008204811660408088019190915291166060860152603b5467ffffffffffffffff8116608087015268010000000000000000900490921660a085015273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811660c08601819052908b16845260388252928290205460ff1660e085015281517f707cd71600000000000000000000000000000000000000000000000000000000815291516101009094019363707cd7169260048082019392918290030181865afa1580156124d4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124f89190614a46565b6040517ffa50f29700000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff919091169063fa50f29790602401602060405180830381865afa158015612564573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125889190614db6565b1515905273ffffffffffffffffffffffffffffffffffffffff86166000908152603560205260409081902090517f2e7263ea00000000000000000000000000000000000000000000000000000000815291925073__$d5ddd09ae98762b8929dd85e54b218e259$__91632e7263ea9161260f91603491603691603791908890600401615069565b60006040518083038186803b15801561262757600080fd5b505af415801561263b573d6000803e3d6000fd5b50505050505050505050505050505050565b6126556135f6565b6fffffffffffffffffffffffffffffffff90811670010000000000000000000000000000000002911617603a55565b6040805173ffffffffffffffffffffffffffffffffffffffff83811660008181526035602090815285822060c0860187525460a086019081528552603b5468010000000000000000900461ffff16818601528486019290925284517ffca513a8000000000000000000000000000000000000000000000000000000008152945190948594859485948594859473__$563c746fa3df0f1858d85f6ef4258864be$__946326ec273f9460349460369460379460608501937f0000000000000000000000000000000000000000000000000000000000000000169263fca513a8926004808401938290030181865afa158015612782573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127a69190614a46565b73ffffffffffffffffffffffffffffffffffffffff90811682528e81166000908152603860209081526040918290205460ff90811694820194909452815160e08a901b7fffffffff00000000000000000000000000000000000000000000000000000000168152600481019890985260248801969096526044870194909452825151606487015293820151608486015291810151831660a4850152606081015190921660c48401526080909101511660e48201526101040160c060405180830381865af415801561287b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061289f919061520f565b949c939b5091995097509550909350915050565b6001805460ff16806128c45750303b155b806128d0575060005481115b61295c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f436f6e747261637420696e7374616e63652068617320616c726561647920626560448201527f656e20696e697469616c697a65640000000000000000000000000000000000006064820152608401610e57565b60015460ff1615801561299957600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168117905560008290555b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16146040518060400160405280600281526020017f313200000000000000000000000000000000000000000000000000000000000081525090612a56576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b50603b80547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000166109c41790558015612ab257600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690555b505050565b73ffffffffffffffffffffffffffffffffffffffff8281166000818152603460205260409081902090517f6973f74400000000000000000000000000000000000000000000000000000000815260048101919091526024810191909152908216604482015273__$c3724b8d563dc83a94e797176cddecb3b9$__90636973f74490606401611793565b612b48613941565b6040517f87b322b200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8085166004830152831660248201526044810182905273__$563c746fa3df0f1858d85f6ef4258864be$__906387b322b29060640160006040518083038186803b158015612bd057600080fd5b505af4158015612be4573d6000803e3d6000fd5b50505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260346020526040812061131990613ace565b603b5460609068010000000000000000900461ffff166000808267ffffffffffffffff811115612c4d57612c4d61475d565b604051908082528060200260200182016040528015612c76578160200160208202803683370190505b50905060005b83811015612d4d5760008181526036602052604090205473ffffffffffffffffffffffffffffffffffffffff1615612d2d5760008181526036602052604090205473ffffffffffffffffffffffffffffffffffffffff1682612cde8584615259565b81518110612cee57612cee615270565b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050612d3b565b82612d378161529f565b9350505b80612d458161529f565b915050612c7c565b5091038152919050565b612d5f6135f6565b60408051808201909152600281527f3136000000000000000000000000000000000000000000000000000000000000602082015260ff8316612dce576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b5060ff8216600090815260376020908152604091829020835181548386015194860151606087015173ffffffffffffffffffffffffffffffffffffffff166601000000000000027fffffffffffff0000000000000000000000000000000000000000ffffffffffff61ffff92831664010000000002167fffffffffffff00000000000000000000000000000000000000000000ffffffff97831662010000027fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000909416929094169190911791909117949094161792909217825560808301518051849392611193926001850192910190613d42565b73ffffffffffffffffffffffffffffffffffffffff868116600090815260346020908152604091829020600401548251808401909352600283527f3131000000000000000000000000000000000000000000000000000000000000918301919091529091163314612f61576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b5073__$db79717e66442ee197e8271d032a066e34$__638a5dadd160346036603760356040518061012001604052808d73ffffffffffffffffffffffffffffffffffffffff1681526020018c73ffffffffffffffffffffffffffffffffffffffff1681526020018b73ffffffffffffffffffffffffffffffffffffffff1681526020018a8152602001898152602001888152602001603b60089054906101000a900461ffff1661ffff1681526020017f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fca513a86040518163ffffffff1660e01b8152600401602060405180830381865afa15801561307b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061309f9190614a46565b73ffffffffffffffffffffffffffffffffffffffff90811682528d166000908152603860209081526040918290205460ff16920191909152517fffffffff0000000000000000000000000000000000000000000000000000000060e088901b1681526131129594939291906004016152d8565b60006040518083038186803b15801561312a57600080fd5b505af415801561313e573d6000803e3d6000fd5b50505050505050505050565b60006131546137b4565b73ffffffffffffffffffffffffffffffffffffffff84166000818152603460205260409081902060395491517f8e743248000000000000000000000000000000000000000000000000000000008152600481019190915260248101929092526044820185905260648201849052608482015273__$b06080f092f400a43662c3f835a4d9baa8$__90638e7432489060a401611293565b6131f26135f6565b6040517f1e3b41450000000000000000000000000000000000000000000000000000000081526034600482015273ffffffffffffffffffffffffffffffffffffffff8216602482015273__$563c746fa3df0f1858d85f6ef4258864be$__90631e3b414590604401611167565b6040517fd505accf000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018890526064810185905260ff8416608482015260a4810183905260c4810182905260009073ffffffffffffffffffffffffffffffffffffffff8a169063d505accf9060e401600060405180830381600087803b1580156132f457600080fd5b505af1158015613308573d6000803e3d6000fd5b5050505060006040518060a001604052808b73ffffffffffffffffffffffffffffffffffffffff1681526020018a815260200189600281111561334d5761334d614b57565b600281111561335e5761335e614b57565b815273ffffffffffffffffffffffffffffffffffffffff89166020808301829052600060409384018190529182526035905281902090517f40e95de600000000000000000000000000000000000000000000000000000000815291925073__$c3724b8d563dc83a94e797176cddecb3b9$__916340e95de6916133eb916034916036918790600401614bc1565b602060405180830381865af4158015613408573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061342c9190614c34565b9a9950505050505050505050565b6134426135f6565b60408051808201909152600281527f3737000000000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff83166134c4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b5073ffffffffffffffffffffffffffffffffffffffff82166000908152603460205260409020600301547501000000000000000000000000000000000000000000900461ffff1615158061355a57506000805260366020527f4cb2b152c1b54ce671907a93c300fd5aa72383a9d4ec19a81e3333632ae92e005473ffffffffffffffffffffffffffffffffffffffff8381169116145b6040518060400160405280600281526020017f3832000000000000000000000000000000000000000000000000000000000000815250906135c8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b5073ffffffffffffffffffffffffffffffffffffffff91909116600090815260346020526040902090359055565b3373ffffffffffffffffffffffffffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663631adfca6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613678573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061369c9190614a46565b73ffffffffffffffffffffffffffffffffffffffff16146040518060400160405280600281526020017f313000000000000000000000000000000000000000000000000000000000000081525090613721576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b50565b6003810154600090700100000000000000000000000000000000900464ffffffffff164281141561376a575050600201546fffffffffffffffffffffffffffffffff1690565b60028301546112d7906fffffffffffffffffffffffffffffffff808216916137a8917001000000000000000000000000000000009091041684613b52565b90613b5f565b50919050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663707cd7166040518163ffffffff1660e01b8152600401602060405180830381865afa15801561381f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138439190614a46565b6040517f726600ce00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff919091169063726600ce90602401602060405180830381865afa1580156138af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138d39190614db6565b6040518060400160405280600181526020017f360000000000000000000000000000000000000000000000000000000000000081525090613721576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663707cd7166040518163ffffffff1660e01b8152600401602060405180830381865afa1580156139ac573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139d09190614a46565b6040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9190911690637be53ca190602401602060405180830381865afa158015613a3c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a609190614db6565b6040518060400160405280600181526020017f310000000000000000000000000000000000000000000000000000000000000081525090613721576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e579190614b44565b6003810154600090700100000000000000000000000000000000900464ffffffffff1642811415613b14575050600101546fffffffffffffffffffffffffffffffff1690565b60018301546112d7906fffffffffffffffffffffffffffffffff808216916137a8917001000000000000000000000000000000009091041684613bb6565b60006112d7838342613bfb565b600081157ffffffffffffffffffffffffffffffffffffffffffe6268e1b017bfe18bffffff83900484111517613b9457600080fd5b506b033b2e3c9fd0803ce800000091026b019d971e4fe8401e74000000010490565b600080613bca64ffffffffff841642615259565b613bd490856153b4565b6301e1338090049050613bf3816b033b2e3c9fd0803ce8000000615420565b949350505050565b600080613c0f64ffffffffff851684615259565b905080613c2b576b033b2e3c9fd0803ce80000009150506112d7565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81016000808060028511613c61576000613c66565b600285035b925066038882915c4000613c7a8a80613b5f565b81613c8757613c876153f1565b0491506301e13380613c99838b613b5f565b81613ca657613ca66153f1565b049050600082613cb686886153b4565b613cc091906153b4565b60029004905060008285613cd4888a6153b4565b613cde91906153b4565b613ce891906153b4565b60069004905080826301e13380613cff8a8f6153b4565b613d099190615438565b613d1f906b033b2e3c9fd0803ce8000000615420565b613d299190615420565b613d339190615420565b9b9a5050505050505050505050565b828054613d4e90614cd8565b90600052602060002090601f016020900481019282613d705760008555613db6565b82601f10613d8957805160ff1916838001178555613db6565b82800160010185558215613db6579182015b82811115613db6578251825591602001919060010190613d9b565b50613dc2929150613dc6565b5090565b5b80821115613dc25760008155600101613dc7565b73ffffffffffffffffffffffffffffffffffffffff8116811461372157600080fd5b8035613e0881613ddb565b919050565b801515811461372157600080fd5b600080600080600060a08688031215613e3357600080fd5b8535613e3e81613ddb565b94506020860135613e4e81613ddb565b93506040860135613e5e81613ddb565b9250606086013591506080860135613e7581613e0d565b809150509295509295909350565b803561ffff81168114613e0857600080fd5b803560ff81168114613e0857600080fd5b600080600080600080600080610100898b031215613ec357600080fd5b8835613ece81613ddb565b9750602089013596506040890135613ee581613ddb565b9550613ef360608a01613e83565b945060808901359350613f0860a08a01613e95565b925060c0890135915060e089013590509295985092959890939650565b60008060408385031215613f3857600080fd5b8235613f4381613ddb565b91506020830135613f5381613ddb565b809150509250929050565b600060208284031215613f7057600080fd5b6112d782613e95565b600080600060608486031215613f8e57600080fd5b8335613f9981613ddb565b95602085013595506040909401359392505050565b600060208284031215613fc057600080fd5b5035919050565b600060208284031215613fd957600080fd5b81356112d781613ddb565b81515181526101e08101602083015161401160208401826fffffffffffffffffffffffffffffffff169052565b50604083015161403560408401826fffffffffffffffffffffffffffffffff169052565b50606083015161405960608401826fffffffffffffffffffffffffffffffff169052565b50608083015161407d60808401826fffffffffffffffffffffffffffffffff169052565b5060a08301516140a160a08401826fffffffffffffffffffffffffffffffff169052565b5060c08301516140ba60c084018264ffffffffff169052565b5060e08301516140d060e084018261ffff169052565b506101008381015173ffffffffffffffffffffffffffffffffffffffff9081169184019190915261012080850151821690840152610140808501518216908401526101608085015190911690830152610180808401516fffffffffffffffffffffffffffffffff908116918401919091526101a0808501518216908401526101c09384015116929091019190915290565b60008083601f84011261417357600080fd5b50813567ffffffffffffffff81111561418b57600080fd5b6020830191508360208285010111156141a357600080fd5b9250929050565b60008060008060008060a087890312156141c357600080fd5b86356141ce81613ddb565b955060208701356141de81613ddb565b945060408701359350606087013567ffffffffffffffff81111561420157600080fd5b61420d89828a01614161565b9094509250614220905060808801613e83565b90509295509295509295565b60006020828403121561423e57600080fd5b6112d782613e83565b6000806000806080858703121561425d57600080fd5b843561426881613ddb565b93506020850135925060408501359150606085013561428681613ddb565b939692955090935050565b600080604083850312156142a457600080fd5b82356142af81613ddb565b91506020830135613f5381613e0d565b600080600080608085870312156142d557600080fd5b84356142e081613ddb565b93506020850135925060408501356142f781613ddb565b915061430560608601613e83565b905092959194509250565b60008060006060848603121561432557600080fd5b833561433081613ddb565b925060208401359150604084013561434781613ddb565b809150509250925092565b6000815180845260005b818110156143785760208185018101518683018201520161435c565b8181111561438a576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061ffff8084511660208401528060208501511660408401528060408501511660608401525073ffffffffffffffffffffffffffffffffffffffff6060840151166080830152608083015160a080840152613bf360c0840182614352565b600080600080600060a0868803121561443857600080fd5b853561444381613ddb565b9450602086013561445381613ddb565b9350604086013561446381613ddb565b9250606086013561447381613ddb565b91506080860135613e7581613ddb565b6000806040838503121561449657600080fd5b82356144a181613ddb565b946020939093013593505050565b60008083601f8401126144c157600080fd5b50813567ffffffffffffffff8111156144d957600080fd5b6020830191508360208260051b85010111156141a357600080fd5b6000806020838503121561450757600080fd5b823567ffffffffffffffff81111561451e57600080fd5b61452a858286016144af565b90969095509350505050565b600080600080600060a0868803121561454e57600080fd5b853561455981613ddb565b9450602086013593506040860135925061447360608701613e83565b600080600080600080600080600080600060e08c8e03121561459657600080fd5b61459f8c613dfd565b9a5067ffffffffffffffff8060208e013511156145bb57600080fd5b6145cb8e60208f01358f016144af565b909b50995060408d01358110156145e157600080fd5b6145f18e60408f01358f016144af565b909950975060608d013581101561460757600080fd5b6146178e60608f01358f016144af565b909750955061462860808e01613dfd565b94508060a08e0135111561463b57600080fd5b5061464c8d60a08e01358e01614161565b909350915061465d60c08d01613e83565b90509295989b509295989b9093969950565b80356fffffffffffffffffffffffffffffffff81168114613e0857600080fd5b600080604083850312156146a257600080fd5b6146ab8361466f565b91506146b96020840161466f565b90509250929050565b6000806000606084860312156146d757600080fd5b83356146e281613ddb565b925060208401356146f281613ddb565b929592945050506040919091013590565b6020808252825182820181905260009190848201906040850190845b8181101561475157835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161471f565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff811182821017156147af576147af61475d565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156147fc576147fc61475d565b604052919050565b6000806040838503121561481757600080fd5b61482083613e95565b915060208084013567ffffffffffffffff8082111561483e57600080fd5b9085019060a0828803121561485257600080fd5b61485a61478c565b61486383613e83565b8152614870848401613e83565b8482015261488060408401613e83565b6040820152606083013561489381613ddb565b60608201526080830135828111156148aa57600080fd5b80840193505087601f8401126148bf57600080fd5b8235828111156148d1576148d161475d565b614901857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016147b5565b9250808352888582860101111561491757600080fd5b8085850186850137600085828501015250816080820152809450505050509250929050565b60008060008060008060c0878903121561495557600080fd5b863561496081613ddb565b9550602087013561497081613ddb565b9450604087013561498081613ddb565b959894975094956060810135955060808101359460a0909101359350915050565b600080600080600080600080610100898b0312156149be57600080fd5b88356149c981613ddb565b975060208901359650604089013595506060890135613ef381613ddb565b60008082840360408112156149fb57600080fd5b8335614a0681613ddb565b925060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082011215614a3857600080fd5b506020830190509250929050565b600060208284031215614a5857600080fd5b81516112d781613ddb565b60006101a08201905086825285602083015284604083015283606083015282516080830152602083015160a0830152604083015173ffffffffffffffffffffffffffffffffffffffff80821660c08501528060608601511660e085015250506080830151610100614aeb8185018373ffffffffffffffffffffffffffffffffffffffff169052565b60a0850151151561012085015260c085015173ffffffffffffffffffffffffffffffffffffffff90811661014086015260e086015160ff166101608601529085015190811661018085015290505b509695505050505050565b6020815260006112d76020830184614352565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60038110614bbd577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b60006101008201905085825284602083015283604083015273ffffffffffffffffffffffffffffffffffffffff808451166060840152602084015160808401526040840151614c1360a0850182614b86565b5060608401511660c0830152608090920151151560e0909101529392505050565b600060208284031215614c4657600080fd5b5051919050565b82815260406020820152600073ffffffffffffffffffffffffffffffffffffffff8084511660408401528060208501511660608401525060408301516080830152606083015160e060a0840152614ca8610120840182614352565b905061ffff60808501511660c084015260a084015160e084015260c0840151610100840152809150509392505050565b600181811c90821680614cec57607f821691505b602082108114156137ae577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006101208201905084825283602083015273ffffffffffffffffffffffffffffffffffffffff8084511660408401528060208501511660608401528060408501511660808401528060608501511660a08401528060808501511660c08401525060a0830151614d9c60e084018261ffff169052565b5060c083015161ffff811661010084015250949350505050565b600060208284031215614dc857600080fd5b81516112d781613e0d565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061ffff80831681811415614e1a57614e1a614dd3565b6001019392505050565b8481526020810184905273ffffffffffffffffffffffffffffffffffffffff83166040820152608081016115e96060830184614b86565b83815260406020808301829052908201839052600090849060608401835b86811015614eb4578335614e8c81613ddb565b73ffffffffffffffffffffffffffffffffffffffff1682529282019290820190600101614e79565b50979650505050505050565b858152602081018590526040810184905260608101839052815173ffffffffffffffffffffffffffffffffffffffff1660808201526102008101602083015173ffffffffffffffffffffffffffffffffffffffff811660a084015250604083015173ffffffffffffffffffffffffffffffffffffffff811660c084015250606083015160e08301526080830151610100614f5c81850183614b86565b60a08501519150610120614f758186018461ffff169052565b60c08601519250610140614f8c8187018515159052565b60e087015161016087810191909152928701516101808701529086015173ffffffffffffffffffffffffffffffffffffffff9081166101a08701529086015160ff166101c0860152908501519081166101e08501529050614b39565b600081518084526020808501945080840160005b8381101561502e57815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101614ffc565b509495945050505050565b600081518084526020808501945080840160005b8381101561502e5781518752958201959082019060010161504d565b85815284602082015283604082015282606082015260a060808201526150a860a08201835173ffffffffffffffffffffffffffffffffffffffff169052565b600060208301516101c08060c08501526150c6610260850183614fe8565b915060408501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60808685030160e08701526151028483615039565b9350606087015191506101008187860301818801526151218584615039565b94506080880151925061012061514e8189018573ffffffffffffffffffffffffffffffffffffffff169052565b60a089015193506101408389880301818a015261516b8786614352565b965060c08a015194506101609350615188848a018661ffff169052565b60e08a0151945061018085818b0152838b015195506101a0935085848b0152828b0151878b0152818b01516101e08b0152848b015196506151e26102008b018873ffffffffffffffffffffffffffffffffffffffff169052565b8a015160ff81166102208b015295506151f9915050565b8701518015156102408801529250614eb4915050565b60008060008060008060c0878903121561522857600080fd5b865195506020870151945060408701519350606087015192506080870151915060a087015190509295509295509295565b60008282101561526b5761526b614dd3565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156152d1576152d1614dd3565b5060010190565b60006101a08201905086825285602083015284604083015283606083015273ffffffffffffffffffffffffffffffffffffffff8084511660808401528060208501511660a084015250604083015161534860c084018273ffffffffffffffffffffffffffffffffffffffff169052565b50606083015160e08301526080830151610100818185015260a085015161012085015260c085015161014085015260e085015191506153a061016085018373ffffffffffffffffffffffffffffffffffffffff169052565b84015160ff81166101808501529050614b39565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156153ec576153ec614dd3565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000821982111561543357615433614dd3565b500190565b60008261546e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea264697066735822122080aadf89d3e671d640d0df3da4ce810529e5e059dd15f51852bbce03f60f26c164736f6c634300080a0033",
}

// PoolABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolMetaData.ABI instead.
var PoolABI = PoolMetaData.ABI

// PoolBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PoolMetaData.Bin instead.
var PoolBin = PoolMetaData.Bin

// DeployPool deploys a new Ethereum contract, binding an instance of Pool to it.
func DeployPool(auth *bind.TransactOpts, backend bind.ContractBackend, provider common.Address) (common.Address, *types.Transaction, *Pool, error) {
	parsed, err := PoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PoolBin), backend, provider)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Pool{PoolCaller: PoolCaller{contract: contract}, PoolTransactor: PoolTransactor{contract: contract}, PoolFilterer: PoolFilterer{contract: contract}}, nil
}

// Pool is an auto generated Go binding around an Ethereum contract.
type Pool struct {
	PoolCaller     // Read-only binding to the contract
	PoolTransactor // Write-only binding to the contract
	PoolFilterer   // Log filterer for contract events
}

// PoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolSession struct {
	Contract     *Pool             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolCallerSession struct {
	Contract *PoolCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolTransactorSession struct {
	Contract     *PoolTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolRaw struct {
	Contract *Pool // Generic contract binding to access the raw methods on
}

// PoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolCallerRaw struct {
	Contract *PoolCaller // Generic read-only contract binding to access the raw methods on
}

// PoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolTransactorRaw struct {
	Contract *PoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPool creates a new instance of Pool, bound to a specific deployed contract.
func NewPool(address common.Address, backend bind.ContractBackend) (*Pool, error) {
	contract, err := bindPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pool{PoolCaller: PoolCaller{contract: contract}, PoolTransactor: PoolTransactor{contract: contract}, PoolFilterer: PoolFilterer{contract: contract}}, nil
}

// NewPoolCaller creates a new read-only instance of Pool, bound to a specific deployed contract.
func NewPoolCaller(address common.Address, caller bind.ContractCaller) (*PoolCaller, error) {
	contract, err := bindPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolCaller{contract: contract}, nil
}

// NewPoolTransactor creates a new write-only instance of Pool, bound to a specific deployed contract.
func NewPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolTransactor, error) {
	contract, err := bindPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolTransactor{contract: contract}, nil
}

// NewPoolFilterer creates a new log filterer instance of Pool, bound to a specific deployed contract.
func NewPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolFilterer, error) {
	contract, err := bindPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolFilterer{contract: contract}, nil
}

// bindPool binds a generic wrapper to an already deployed contract.
func bindPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pool *PoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pool.Contract.PoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pool *PoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pool.Contract.PoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pool *PoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pool.Contract.PoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pool *PoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pool *PoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pool *PoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pool.Contract.contract.Transact(opts, method, params...)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_Pool *PoolCaller) ADDRESSESPROVIDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "ADDRESSES_PROVIDER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_Pool *PoolSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _Pool.Contract.ADDRESSESPROVIDER(&_Pool.CallOpts)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_Pool *PoolCallerSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _Pool.Contract.ADDRESSESPROVIDER(&_Pool.CallOpts)
}

// BRIDGEPROTOCOLFEE is a free data retrieval call binding the contract method 0x272d9072.
//
// Solidity: function BRIDGE_PROTOCOL_FEE() view returns(uint256)
func (_Pool *PoolCaller) BRIDGEPROTOCOLFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "BRIDGE_PROTOCOL_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BRIDGEPROTOCOLFEE is a free data retrieval call binding the contract method 0x272d9072.
//
// Solidity: function BRIDGE_PROTOCOL_FEE() view returns(uint256)
func (_Pool *PoolSession) BRIDGEPROTOCOLFEE() (*big.Int, error) {
	return _Pool.Contract.BRIDGEPROTOCOLFEE(&_Pool.CallOpts)
}

// BRIDGEPROTOCOLFEE is a free data retrieval call binding the contract method 0x272d9072.
//
// Solidity: function BRIDGE_PROTOCOL_FEE() view returns(uint256)
func (_Pool *PoolCallerSession) BRIDGEPROTOCOLFEE() (*big.Int, error) {
	return _Pool.Contract.BRIDGEPROTOCOLFEE(&_Pool.CallOpts)
}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint128)
func (_Pool *PoolCaller) FLASHLOANPREMIUMTOTAL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "FLASHLOAN_PREMIUM_TOTAL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint128)
func (_Pool *PoolSession) FLASHLOANPREMIUMTOTAL() (*big.Int, error) {
	return _Pool.Contract.FLASHLOANPREMIUMTOTAL(&_Pool.CallOpts)
}

// FLASHLOANPREMIUMTOTAL is a free data retrieval call binding the contract method 0x074b2e43.
//
// Solidity: function FLASHLOAN_PREMIUM_TOTAL() view returns(uint128)
func (_Pool *PoolCallerSession) FLASHLOANPREMIUMTOTAL() (*big.Int, error) {
	return _Pool.Contract.FLASHLOANPREMIUMTOTAL(&_Pool.CallOpts)
}

// FLASHLOANPREMIUMTOPROTOCOL is a free data retrieval call binding the contract method 0x6a99c036.
//
// Solidity: function FLASHLOAN_PREMIUM_TO_PROTOCOL() view returns(uint128)
func (_Pool *PoolCaller) FLASHLOANPREMIUMTOPROTOCOL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "FLASHLOAN_PREMIUM_TO_PROTOCOL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FLASHLOANPREMIUMTOPROTOCOL is a free data retrieval call binding the contract method 0x6a99c036.
//
// Solidity: function FLASHLOAN_PREMIUM_TO_PROTOCOL() view returns(uint128)
func (_Pool *PoolSession) FLASHLOANPREMIUMTOPROTOCOL() (*big.Int, error) {
	return _Pool.Contract.FLASHLOANPREMIUMTOPROTOCOL(&_Pool.CallOpts)
}

// FLASHLOANPREMIUMTOPROTOCOL is a free data retrieval call binding the contract method 0x6a99c036.
//
// Solidity: function FLASHLOAN_PREMIUM_TO_PROTOCOL() view returns(uint128)
func (_Pool *PoolCallerSession) FLASHLOANPREMIUMTOPROTOCOL() (*big.Int, error) {
	return _Pool.Contract.FLASHLOANPREMIUMTOPROTOCOL(&_Pool.CallOpts)
}

// MAXNUMBERRESERVES is a free data retrieval call binding the contract method 0xf8119d51.
//
// Solidity: function MAX_NUMBER_RESERVES() view returns(uint16)
func (_Pool *PoolCaller) MAXNUMBERRESERVES(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "MAX_NUMBER_RESERVES")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MAXNUMBERRESERVES is a free data retrieval call binding the contract method 0xf8119d51.
//
// Solidity: function MAX_NUMBER_RESERVES() view returns(uint16)
func (_Pool *PoolSession) MAXNUMBERRESERVES() (uint16, error) {
	return _Pool.Contract.MAXNUMBERRESERVES(&_Pool.CallOpts)
}

// MAXNUMBERRESERVES is a free data retrieval call binding the contract method 0xf8119d51.
//
// Solidity: function MAX_NUMBER_RESERVES() view returns(uint16)
func (_Pool *PoolCallerSession) MAXNUMBERRESERVES() (uint16, error) {
	return _Pool.Contract.MAXNUMBERRESERVES(&_Pool.CallOpts)
}

// MAXSTABLERATEBORROWSIZEPERCENT is a free data retrieval call binding the contract method 0xe82fec2f.
//
// Solidity: function MAX_STABLE_RATE_BORROW_SIZE_PERCENT() view returns(uint256)
func (_Pool *PoolCaller) MAXSTABLERATEBORROWSIZEPERCENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "MAX_STABLE_RATE_BORROW_SIZE_PERCENT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSTABLERATEBORROWSIZEPERCENT is a free data retrieval call binding the contract method 0xe82fec2f.
//
// Solidity: function MAX_STABLE_RATE_BORROW_SIZE_PERCENT() view returns(uint256)
func (_Pool *PoolSession) MAXSTABLERATEBORROWSIZEPERCENT() (*big.Int, error) {
	return _Pool.Contract.MAXSTABLERATEBORROWSIZEPERCENT(&_Pool.CallOpts)
}

// MAXSTABLERATEBORROWSIZEPERCENT is a free data retrieval call binding the contract method 0xe82fec2f.
//
// Solidity: function MAX_STABLE_RATE_BORROW_SIZE_PERCENT() view returns(uint256)
func (_Pool *PoolCallerSession) MAXSTABLERATEBORROWSIZEPERCENT() (*big.Int, error) {
	return _Pool.Contract.MAXSTABLERATEBORROWSIZEPERCENT(&_Pool.CallOpts)
}

// POOLREVISION is a free data retrieval call binding the contract method 0x0148170e.
//
// Solidity: function POOL_REVISION() view returns(uint256)
func (_Pool *PoolCaller) POOLREVISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "POOL_REVISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// POOLREVISION is a free data retrieval call binding the contract method 0x0148170e.
//
// Solidity: function POOL_REVISION() view returns(uint256)
func (_Pool *PoolSession) POOLREVISION() (*big.Int, error) {
	return _Pool.Contract.POOLREVISION(&_Pool.CallOpts)
}

// POOLREVISION is a free data retrieval call binding the contract method 0x0148170e.
//
// Solidity: function POOL_REVISION() view returns(uint256)
func (_Pool *PoolCallerSession) POOLREVISION() (*big.Int, error) {
	return _Pool.Contract.POOLREVISION(&_Pool.CallOpts)
}

// GetConfiguration is a free data retrieval call binding the contract method 0xc44b11f7.
//
// Solidity: function getConfiguration(address asset) view returns((uint256))
func (_Pool *PoolCaller) GetConfiguration(opts *bind.CallOpts, asset common.Address) (DataTypesReserveConfigurationMap, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getConfiguration", asset)

	if err != nil {
		return *new(DataTypesReserveConfigurationMap), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesReserveConfigurationMap)).(*DataTypesReserveConfigurationMap)

	return out0, err

}

// GetConfiguration is a free data retrieval call binding the contract method 0xc44b11f7.
//
// Solidity: function getConfiguration(address asset) view returns((uint256))
func (_Pool *PoolSession) GetConfiguration(asset common.Address) (DataTypesReserveConfigurationMap, error) {
	return _Pool.Contract.GetConfiguration(&_Pool.CallOpts, asset)
}

// GetConfiguration is a free data retrieval call binding the contract method 0xc44b11f7.
//
// Solidity: function getConfiguration(address asset) view returns((uint256))
func (_Pool *PoolCallerSession) GetConfiguration(asset common.Address) (DataTypesReserveConfigurationMap, error) {
	return _Pool.Contract.GetConfiguration(&_Pool.CallOpts, asset)
}

// GetEModeCategoryData is a free data retrieval call binding the contract method 0x6c6f6ae1.
//
// Solidity: function getEModeCategoryData(uint8 id) view returns((uint16,uint16,uint16,address,string))
func (_Pool *PoolCaller) GetEModeCategoryData(opts *bind.CallOpts, id uint8) (DataTypesEModeCategory, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getEModeCategoryData", id)

	if err != nil {
		return *new(DataTypesEModeCategory), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesEModeCategory)).(*DataTypesEModeCategory)

	return out0, err

}

// GetEModeCategoryData is a free data retrieval call binding the contract method 0x6c6f6ae1.
//
// Solidity: function getEModeCategoryData(uint8 id) view returns((uint16,uint16,uint16,address,string))
func (_Pool *PoolSession) GetEModeCategoryData(id uint8) (DataTypesEModeCategory, error) {
	return _Pool.Contract.GetEModeCategoryData(&_Pool.CallOpts, id)
}

// GetEModeCategoryData is a free data retrieval call binding the contract method 0x6c6f6ae1.
//
// Solidity: function getEModeCategoryData(uint8 id) view returns((uint16,uint16,uint16,address,string))
func (_Pool *PoolCallerSession) GetEModeCategoryData(id uint8) (DataTypesEModeCategory, error) {
	return _Pool.Contract.GetEModeCategoryData(&_Pool.CallOpts, id)
}

// GetReserveAddressById is a free data retrieval call binding the contract method 0x52751797.
//
// Solidity: function getReserveAddressById(uint16 id) view returns(address)
func (_Pool *PoolCaller) GetReserveAddressById(opts *bind.CallOpts, id uint16) (common.Address, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getReserveAddressById", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetReserveAddressById is a free data retrieval call binding the contract method 0x52751797.
//
// Solidity: function getReserveAddressById(uint16 id) view returns(address)
func (_Pool *PoolSession) GetReserveAddressById(id uint16) (common.Address, error) {
	return _Pool.Contract.GetReserveAddressById(&_Pool.CallOpts, id)
}

// GetReserveAddressById is a free data retrieval call binding the contract method 0x52751797.
//
// Solidity: function getReserveAddressById(uint16 id) view returns(address)
func (_Pool *PoolCallerSession) GetReserveAddressById(id uint16) (common.Address, error) {
	return _Pool.Contract.GetReserveAddressById(&_Pool.CallOpts, id)
}

// GetReserveData is a free data retrieval call binding the contract method 0x35ea6a75.
//
// Solidity: function getReserveData(address asset) view returns(((uint256),uint128,uint128,uint128,uint128,uint128,uint40,uint16,address,address,address,address,uint128,uint128,uint128))
func (_Pool *PoolCaller) GetReserveData(opts *bind.CallOpts, asset common.Address) (DataTypesReserveData, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getReserveData", asset)

	if err != nil {
		return *new(DataTypesReserveData), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesReserveData)).(*DataTypesReserveData)

	return out0, err

}

// GetReserveData is a free data retrieval call binding the contract method 0x35ea6a75.
//
// Solidity: function getReserveData(address asset) view returns(((uint256),uint128,uint128,uint128,uint128,uint128,uint40,uint16,address,address,address,address,uint128,uint128,uint128))
func (_Pool *PoolSession) GetReserveData(asset common.Address) (DataTypesReserveData, error) {
	return _Pool.Contract.GetReserveData(&_Pool.CallOpts, asset)
}

// GetReserveData is a free data retrieval call binding the contract method 0x35ea6a75.
//
// Solidity: function getReserveData(address asset) view returns(((uint256),uint128,uint128,uint128,uint128,uint128,uint40,uint16,address,address,address,address,uint128,uint128,uint128))
func (_Pool *PoolCallerSession) GetReserveData(asset common.Address) (DataTypesReserveData, error) {
	return _Pool.Contract.GetReserveData(&_Pool.CallOpts, asset)
}

// GetReserveNormalizedIncome is a free data retrieval call binding the contract method 0xd15e0053.
//
// Solidity: function getReserveNormalizedIncome(address asset) view returns(uint256)
func (_Pool *PoolCaller) GetReserveNormalizedIncome(opts *bind.CallOpts, asset common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getReserveNormalizedIncome", asset)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetReserveNormalizedIncome is a free data retrieval call binding the contract method 0xd15e0053.
//
// Solidity: function getReserveNormalizedIncome(address asset) view returns(uint256)
func (_Pool *PoolSession) GetReserveNormalizedIncome(asset common.Address) (*big.Int, error) {
	return _Pool.Contract.GetReserveNormalizedIncome(&_Pool.CallOpts, asset)
}

// GetReserveNormalizedIncome is a free data retrieval call binding the contract method 0xd15e0053.
//
// Solidity: function getReserveNormalizedIncome(address asset) view returns(uint256)
func (_Pool *PoolCallerSession) GetReserveNormalizedIncome(asset common.Address) (*big.Int, error) {
	return _Pool.Contract.GetReserveNormalizedIncome(&_Pool.CallOpts, asset)
}

// GetReserveNormalizedVariableDebt is a free data retrieval call binding the contract method 0x386497fd.
//
// Solidity: function getReserveNormalizedVariableDebt(address asset) view returns(uint256)
func (_Pool *PoolCaller) GetReserveNormalizedVariableDebt(opts *bind.CallOpts, asset common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getReserveNormalizedVariableDebt", asset)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetReserveNormalizedVariableDebt is a free data retrieval call binding the contract method 0x386497fd.
//
// Solidity: function getReserveNormalizedVariableDebt(address asset) view returns(uint256)
func (_Pool *PoolSession) GetReserveNormalizedVariableDebt(asset common.Address) (*big.Int, error) {
	return _Pool.Contract.GetReserveNormalizedVariableDebt(&_Pool.CallOpts, asset)
}

// GetReserveNormalizedVariableDebt is a free data retrieval call binding the contract method 0x386497fd.
//
// Solidity: function getReserveNormalizedVariableDebt(address asset) view returns(uint256)
func (_Pool *PoolCallerSession) GetReserveNormalizedVariableDebt(asset common.Address) (*big.Int, error) {
	return _Pool.Contract.GetReserveNormalizedVariableDebt(&_Pool.CallOpts, asset)
}

// GetReservesList is a free data retrieval call binding the contract method 0xd1946dbc.
//
// Solidity: function getReservesList() view returns(address[])
func (_Pool *PoolCaller) GetReservesList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getReservesList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetReservesList is a free data retrieval call binding the contract method 0xd1946dbc.
//
// Solidity: function getReservesList() view returns(address[])
func (_Pool *PoolSession) GetReservesList() ([]common.Address, error) {
	return _Pool.Contract.GetReservesList(&_Pool.CallOpts)
}

// GetReservesList is a free data retrieval call binding the contract method 0xd1946dbc.
//
// Solidity: function getReservesList() view returns(address[])
func (_Pool *PoolCallerSession) GetReservesList() ([]common.Address, error) {
	return _Pool.Contract.GetReservesList(&_Pool.CallOpts)
}

// GetUserAccountData is a free data retrieval call binding the contract method 0xbf92857c.
//
// Solidity: function getUserAccountData(address user) view returns(uint256 totalCollateralBase, uint256 totalDebtBase, uint256 availableBorrowsBase, uint256 currentLiquidationThreshold, uint256 ltv, uint256 healthFactor)
func (_Pool *PoolCaller) GetUserAccountData(opts *bind.CallOpts, user common.Address) (struct {
	TotalCollateralBase         *big.Int
	TotalDebtBase               *big.Int
	AvailableBorrowsBase        *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getUserAccountData", user)

	outstruct := new(struct {
		TotalCollateralBase         *big.Int
		TotalDebtBase               *big.Int
		AvailableBorrowsBase        *big.Int
		CurrentLiquidationThreshold *big.Int
		Ltv                         *big.Int
		HealthFactor                *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalCollateralBase = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalDebtBase = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.AvailableBorrowsBase = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.CurrentLiquidationThreshold = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Ltv = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.HealthFactor = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUserAccountData is a free data retrieval call binding the contract method 0xbf92857c.
//
// Solidity: function getUserAccountData(address user) view returns(uint256 totalCollateralBase, uint256 totalDebtBase, uint256 availableBorrowsBase, uint256 currentLiquidationThreshold, uint256 ltv, uint256 healthFactor)
func (_Pool *PoolSession) GetUserAccountData(user common.Address) (struct {
	TotalCollateralBase         *big.Int
	TotalDebtBase               *big.Int
	AvailableBorrowsBase        *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}, error) {
	return _Pool.Contract.GetUserAccountData(&_Pool.CallOpts, user)
}

// GetUserAccountData is a free data retrieval call binding the contract method 0xbf92857c.
//
// Solidity: function getUserAccountData(address user) view returns(uint256 totalCollateralBase, uint256 totalDebtBase, uint256 availableBorrowsBase, uint256 currentLiquidationThreshold, uint256 ltv, uint256 healthFactor)
func (_Pool *PoolCallerSession) GetUserAccountData(user common.Address) (struct {
	TotalCollateralBase         *big.Int
	TotalDebtBase               *big.Int
	AvailableBorrowsBase        *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}, error) {
	return _Pool.Contract.GetUserAccountData(&_Pool.CallOpts, user)
}

// GetUserConfiguration is a free data retrieval call binding the contract method 0x4417a583.
//
// Solidity: function getUserConfiguration(address user) view returns((uint256))
func (_Pool *PoolCaller) GetUserConfiguration(opts *bind.CallOpts, user common.Address) (DataTypesUserConfigurationMap, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getUserConfiguration", user)

	if err != nil {
		return *new(DataTypesUserConfigurationMap), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesUserConfigurationMap)).(*DataTypesUserConfigurationMap)

	return out0, err

}

// GetUserConfiguration is a free data retrieval call binding the contract method 0x4417a583.
//
// Solidity: function getUserConfiguration(address user) view returns((uint256))
func (_Pool *PoolSession) GetUserConfiguration(user common.Address) (DataTypesUserConfigurationMap, error) {
	return _Pool.Contract.GetUserConfiguration(&_Pool.CallOpts, user)
}

// GetUserConfiguration is a free data retrieval call binding the contract method 0x4417a583.
//
// Solidity: function getUserConfiguration(address user) view returns((uint256))
func (_Pool *PoolCallerSession) GetUserConfiguration(user common.Address) (DataTypesUserConfigurationMap, error) {
	return _Pool.Contract.GetUserConfiguration(&_Pool.CallOpts, user)
}

// GetUserEMode is a free data retrieval call binding the contract method 0xeddf1b79.
//
// Solidity: function getUserEMode(address user) view returns(uint256)
func (_Pool *PoolCaller) GetUserEMode(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getUserEMode", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserEMode is a free data retrieval call binding the contract method 0xeddf1b79.
//
// Solidity: function getUserEMode(address user) view returns(uint256)
func (_Pool *PoolSession) GetUserEMode(user common.Address) (*big.Int, error) {
	return _Pool.Contract.GetUserEMode(&_Pool.CallOpts, user)
}

// GetUserEMode is a free data retrieval call binding the contract method 0xeddf1b79.
//
// Solidity: function getUserEMode(address user) view returns(uint256)
func (_Pool *PoolCallerSession) GetUserEMode(user common.Address) (*big.Int, error) {
	return _Pool.Contract.GetUserEMode(&_Pool.CallOpts, user)
}

// BackUnbacked is a paid mutator transaction binding the contract method 0xd65dc7a1.
//
// Solidity: function backUnbacked(address asset, uint256 amount, uint256 fee) returns(uint256)
func (_Pool *PoolTransactor) BackUnbacked(opts *bind.TransactOpts, asset common.Address, amount *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "backUnbacked", asset, amount, fee)
}

// BackUnbacked is a paid mutator transaction binding the contract method 0xd65dc7a1.
//
// Solidity: function backUnbacked(address asset, uint256 amount, uint256 fee) returns(uint256)
func (_Pool *PoolSession) BackUnbacked(asset common.Address, amount *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.BackUnbacked(&_Pool.TransactOpts, asset, amount, fee)
}

// BackUnbacked is a paid mutator transaction binding the contract method 0xd65dc7a1.
//
// Solidity: function backUnbacked(address asset, uint256 amount, uint256 fee) returns(uint256)
func (_Pool *PoolTransactorSession) BackUnbacked(asset common.Address, amount *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.BackUnbacked(&_Pool.TransactOpts, asset, amount, fee)
}

// Borrow is a paid mutator transaction binding the contract method 0xa415bcad.
//
// Solidity: function borrow(address asset, uint256 amount, uint256 interestRateMode, uint16 referralCode, address onBehalfOf) returns()
func (_Pool *PoolTransactor) Borrow(opts *bind.TransactOpts, asset common.Address, amount *big.Int, interestRateMode *big.Int, referralCode uint16, onBehalfOf common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "borrow", asset, amount, interestRateMode, referralCode, onBehalfOf)
}

// Borrow is a paid mutator transaction binding the contract method 0xa415bcad.
//
// Solidity: function borrow(address asset, uint256 amount, uint256 interestRateMode, uint16 referralCode, address onBehalfOf) returns()
func (_Pool *PoolSession) Borrow(asset common.Address, amount *big.Int, interestRateMode *big.Int, referralCode uint16, onBehalfOf common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Borrow(&_Pool.TransactOpts, asset, amount, interestRateMode, referralCode, onBehalfOf)
}

// Borrow is a paid mutator transaction binding the contract method 0xa415bcad.
//
// Solidity: function borrow(address asset, uint256 amount, uint256 interestRateMode, uint16 referralCode, address onBehalfOf) returns()
func (_Pool *PoolTransactorSession) Borrow(asset common.Address, amount *big.Int, interestRateMode *big.Int, referralCode uint16, onBehalfOf common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Borrow(&_Pool.TransactOpts, asset, amount, interestRateMode, referralCode, onBehalfOf)
}

// ConfigureEModeCategory is a paid mutator transaction binding the contract method 0xd579ea7d.
//
// Solidity: function configureEModeCategory(uint8 id, (uint16,uint16,uint16,address,string) category) returns()
func (_Pool *PoolTransactor) ConfigureEModeCategory(opts *bind.TransactOpts, id uint8, category DataTypesEModeCategory) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "configureEModeCategory", id, category)
}

// ConfigureEModeCategory is a paid mutator transaction binding the contract method 0xd579ea7d.
//
// Solidity: function configureEModeCategory(uint8 id, (uint16,uint16,uint16,address,string) category) returns()
func (_Pool *PoolSession) ConfigureEModeCategory(id uint8, category DataTypesEModeCategory) (*types.Transaction, error) {
	return _Pool.Contract.ConfigureEModeCategory(&_Pool.TransactOpts, id, category)
}

// ConfigureEModeCategory is a paid mutator transaction binding the contract method 0xd579ea7d.
//
// Solidity: function configureEModeCategory(uint8 id, (uint16,uint16,uint16,address,string) category) returns()
func (_Pool *PoolTransactorSession) ConfigureEModeCategory(id uint8, category DataTypesEModeCategory) (*types.Transaction, error) {
	return _Pool.Contract.ConfigureEModeCategory(&_Pool.TransactOpts, id, category)
}

// Deposit is a paid mutator transaction binding the contract method 0xe8eda9df.
//
// Solidity: function deposit(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolTransactor) Deposit(opts *bind.TransactOpts, asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "deposit", asset, amount, onBehalfOf, referralCode)
}

// Deposit is a paid mutator transaction binding the contract method 0xe8eda9df.
//
// Solidity: function deposit(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolSession) Deposit(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.Deposit(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode)
}

// Deposit is a paid mutator transaction binding the contract method 0xe8eda9df.
//
// Solidity: function deposit(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolTransactorSession) Deposit(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.Deposit(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode)
}

// DropReserve is a paid mutator transaction binding the contract method 0x63c9b860.
//
// Solidity: function dropReserve(address asset) returns()
func (_Pool *PoolTransactor) DropReserve(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "dropReserve", asset)
}

// DropReserve is a paid mutator transaction binding the contract method 0x63c9b860.
//
// Solidity: function dropReserve(address asset) returns()
func (_Pool *PoolSession) DropReserve(asset common.Address) (*types.Transaction, error) {
	return _Pool.Contract.DropReserve(&_Pool.TransactOpts, asset)
}

// DropReserve is a paid mutator transaction binding the contract method 0x63c9b860.
//
// Solidity: function dropReserve(address asset) returns()
func (_Pool *PoolTransactorSession) DropReserve(asset common.Address) (*types.Transaction, error) {
	return _Pool.Contract.DropReserve(&_Pool.TransactOpts, asset)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xd5ed3933.
//
// Solidity: function finalizeTransfer(address asset, address from, address to, uint256 amount, uint256 balanceFromBefore, uint256 balanceToBefore) returns()
func (_Pool *PoolTransactor) FinalizeTransfer(opts *bind.TransactOpts, asset common.Address, from common.Address, to common.Address, amount *big.Int, balanceFromBefore *big.Int, balanceToBefore *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "finalizeTransfer", asset, from, to, amount, balanceFromBefore, balanceToBefore)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xd5ed3933.
//
// Solidity: function finalizeTransfer(address asset, address from, address to, uint256 amount, uint256 balanceFromBefore, uint256 balanceToBefore) returns()
func (_Pool *PoolSession) FinalizeTransfer(asset common.Address, from common.Address, to common.Address, amount *big.Int, balanceFromBefore *big.Int, balanceToBefore *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.FinalizeTransfer(&_Pool.TransactOpts, asset, from, to, amount, balanceFromBefore, balanceToBefore)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xd5ed3933.
//
// Solidity: function finalizeTransfer(address asset, address from, address to, uint256 amount, uint256 balanceFromBefore, uint256 balanceToBefore) returns()
func (_Pool *PoolTransactorSession) FinalizeTransfer(asset common.Address, from common.Address, to common.Address, amount *big.Int, balanceFromBefore *big.Int, balanceToBefore *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.FinalizeTransfer(&_Pool.TransactOpts, asset, from, to, amount, balanceFromBefore, balanceToBefore)
}

// FlashLoan is a paid mutator transaction binding the contract method 0xab9c4b5d.
//
// Solidity: function flashLoan(address receiverAddress, address[] assets, uint256[] amounts, uint256[] interestRateModes, address onBehalfOf, bytes params, uint16 referralCode) returns()
func (_Pool *PoolTransactor) FlashLoan(opts *bind.TransactOpts, receiverAddress common.Address, assets []common.Address, amounts []*big.Int, interestRateModes []*big.Int, onBehalfOf common.Address, params []byte, referralCode uint16) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "flashLoan", receiverAddress, assets, amounts, interestRateModes, onBehalfOf, params, referralCode)
}

// FlashLoan is a paid mutator transaction binding the contract method 0xab9c4b5d.
//
// Solidity: function flashLoan(address receiverAddress, address[] assets, uint256[] amounts, uint256[] interestRateModes, address onBehalfOf, bytes params, uint16 referralCode) returns()
func (_Pool *PoolSession) FlashLoan(receiverAddress common.Address, assets []common.Address, amounts []*big.Int, interestRateModes []*big.Int, onBehalfOf common.Address, params []byte, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.FlashLoan(&_Pool.TransactOpts, receiverAddress, assets, amounts, interestRateModes, onBehalfOf, params, referralCode)
}

// FlashLoan is a paid mutator transaction binding the contract method 0xab9c4b5d.
//
// Solidity: function flashLoan(address receiverAddress, address[] assets, uint256[] amounts, uint256[] interestRateModes, address onBehalfOf, bytes params, uint16 referralCode) returns()
func (_Pool *PoolTransactorSession) FlashLoan(receiverAddress common.Address, assets []common.Address, amounts []*big.Int, interestRateModes []*big.Int, onBehalfOf common.Address, params []byte, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.FlashLoan(&_Pool.TransactOpts, receiverAddress, assets, amounts, interestRateModes, onBehalfOf, params, referralCode)
}

// FlashLoanSimple is a paid mutator transaction binding the contract method 0x42b0b77c.
//
// Solidity: function flashLoanSimple(address receiverAddress, address asset, uint256 amount, bytes params, uint16 referralCode) returns()
func (_Pool *PoolTransactor) FlashLoanSimple(opts *bind.TransactOpts, receiverAddress common.Address, asset common.Address, amount *big.Int, params []byte, referralCode uint16) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "flashLoanSimple", receiverAddress, asset, amount, params, referralCode)
}

// FlashLoanSimple is a paid mutator transaction binding the contract method 0x42b0b77c.
//
// Solidity: function flashLoanSimple(address receiverAddress, address asset, uint256 amount, bytes params, uint16 referralCode) returns()
func (_Pool *PoolSession) FlashLoanSimple(receiverAddress common.Address, asset common.Address, amount *big.Int, params []byte, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.FlashLoanSimple(&_Pool.TransactOpts, receiverAddress, asset, amount, params, referralCode)
}

// FlashLoanSimple is a paid mutator transaction binding the contract method 0x42b0b77c.
//
// Solidity: function flashLoanSimple(address receiverAddress, address asset, uint256 amount, bytes params, uint16 referralCode) returns()
func (_Pool *PoolTransactorSession) FlashLoanSimple(receiverAddress common.Address, asset common.Address, amount *big.Int, params []byte, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.FlashLoanSimple(&_Pool.TransactOpts, receiverAddress, asset, amount, params, referralCode)
}

// InitReserve is a paid mutator transaction binding the contract method 0x7a708e92.
//
// Solidity: function initReserve(address asset, address aTokenAddress, address stableDebtAddress, address variableDebtAddress, address interestRateStrategyAddress) returns()
func (_Pool *PoolTransactor) InitReserve(opts *bind.TransactOpts, asset common.Address, aTokenAddress common.Address, stableDebtAddress common.Address, variableDebtAddress common.Address, interestRateStrategyAddress common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "initReserve", asset, aTokenAddress, stableDebtAddress, variableDebtAddress, interestRateStrategyAddress)
}

// InitReserve is a paid mutator transaction binding the contract method 0x7a708e92.
//
// Solidity: function initReserve(address asset, address aTokenAddress, address stableDebtAddress, address variableDebtAddress, address interestRateStrategyAddress) returns()
func (_Pool *PoolSession) InitReserve(asset common.Address, aTokenAddress common.Address, stableDebtAddress common.Address, variableDebtAddress common.Address, interestRateStrategyAddress common.Address) (*types.Transaction, error) {
	return _Pool.Contract.InitReserve(&_Pool.TransactOpts, asset, aTokenAddress, stableDebtAddress, variableDebtAddress, interestRateStrategyAddress)
}

// InitReserve is a paid mutator transaction binding the contract method 0x7a708e92.
//
// Solidity: function initReserve(address asset, address aTokenAddress, address stableDebtAddress, address variableDebtAddress, address interestRateStrategyAddress) returns()
func (_Pool *PoolTransactorSession) InitReserve(asset common.Address, aTokenAddress common.Address, stableDebtAddress common.Address, variableDebtAddress common.Address, interestRateStrategyAddress common.Address) (*types.Transaction, error) {
	return _Pool.Contract.InitReserve(&_Pool.TransactOpts, asset, aTokenAddress, stableDebtAddress, variableDebtAddress, interestRateStrategyAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address provider) returns()
func (_Pool *PoolTransactor) Initialize(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "initialize", provider)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address provider) returns()
func (_Pool *PoolSession) Initialize(provider common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Initialize(&_Pool.TransactOpts, provider)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address provider) returns()
func (_Pool *PoolTransactorSession) Initialize(provider common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Initialize(&_Pool.TransactOpts, provider)
}

// LiquidationCall is a paid mutator transaction binding the contract method 0x00a718a9.
//
// Solidity: function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover, bool receiveAToken) returns()
func (_Pool *PoolTransactor) LiquidationCall(opts *bind.TransactOpts, collateralAsset common.Address, debtAsset common.Address, user common.Address, debtToCover *big.Int, receiveAToken bool) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "liquidationCall", collateralAsset, debtAsset, user, debtToCover, receiveAToken)
}

// LiquidationCall is a paid mutator transaction binding the contract method 0x00a718a9.
//
// Solidity: function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover, bool receiveAToken) returns()
func (_Pool *PoolSession) LiquidationCall(collateralAsset common.Address, debtAsset common.Address, user common.Address, debtToCover *big.Int, receiveAToken bool) (*types.Transaction, error) {
	return _Pool.Contract.LiquidationCall(&_Pool.TransactOpts, collateralAsset, debtAsset, user, debtToCover, receiveAToken)
}

// LiquidationCall is a paid mutator transaction binding the contract method 0x00a718a9.
//
// Solidity: function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover, bool receiveAToken) returns()
func (_Pool *PoolTransactorSession) LiquidationCall(collateralAsset common.Address, debtAsset common.Address, user common.Address, debtToCover *big.Int, receiveAToken bool) (*types.Transaction, error) {
	return _Pool.Contract.LiquidationCall(&_Pool.TransactOpts, collateralAsset, debtAsset, user, debtToCover, receiveAToken)
}

// MintToTreasury is a paid mutator transaction binding the contract method 0x9cd19996.
//
// Solidity: function mintToTreasury(address[] assets) returns()
func (_Pool *PoolTransactor) MintToTreasury(opts *bind.TransactOpts, assets []common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "mintToTreasury", assets)
}

// MintToTreasury is a paid mutator transaction binding the contract method 0x9cd19996.
//
// Solidity: function mintToTreasury(address[] assets) returns()
func (_Pool *PoolSession) MintToTreasury(assets []common.Address) (*types.Transaction, error) {
	return _Pool.Contract.MintToTreasury(&_Pool.TransactOpts, assets)
}

// MintToTreasury is a paid mutator transaction binding the contract method 0x9cd19996.
//
// Solidity: function mintToTreasury(address[] assets) returns()
func (_Pool *PoolTransactorSession) MintToTreasury(assets []common.Address) (*types.Transaction, error) {
	return _Pool.Contract.MintToTreasury(&_Pool.TransactOpts, assets)
}

// MintUnbacked is a paid mutator transaction binding the contract method 0x69a933a5.
//
// Solidity: function mintUnbacked(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolTransactor) MintUnbacked(opts *bind.TransactOpts, asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "mintUnbacked", asset, amount, onBehalfOf, referralCode)
}

// MintUnbacked is a paid mutator transaction binding the contract method 0x69a933a5.
//
// Solidity: function mintUnbacked(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolSession) MintUnbacked(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.MintUnbacked(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode)
}

// MintUnbacked is a paid mutator transaction binding the contract method 0x69a933a5.
//
// Solidity: function mintUnbacked(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolTransactorSession) MintUnbacked(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.MintUnbacked(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode)
}

// RebalanceStableBorrowRate is a paid mutator transaction binding the contract method 0xcd112382.
//
// Solidity: function rebalanceStableBorrowRate(address asset, address user) returns()
func (_Pool *PoolTransactor) RebalanceStableBorrowRate(opts *bind.TransactOpts, asset common.Address, user common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "rebalanceStableBorrowRate", asset, user)
}

// RebalanceStableBorrowRate is a paid mutator transaction binding the contract method 0xcd112382.
//
// Solidity: function rebalanceStableBorrowRate(address asset, address user) returns()
func (_Pool *PoolSession) RebalanceStableBorrowRate(asset common.Address, user common.Address) (*types.Transaction, error) {
	return _Pool.Contract.RebalanceStableBorrowRate(&_Pool.TransactOpts, asset, user)
}

// RebalanceStableBorrowRate is a paid mutator transaction binding the contract method 0xcd112382.
//
// Solidity: function rebalanceStableBorrowRate(address asset, address user) returns()
func (_Pool *PoolTransactorSession) RebalanceStableBorrowRate(asset common.Address, user common.Address) (*types.Transaction, error) {
	return _Pool.Contract.RebalanceStableBorrowRate(&_Pool.TransactOpts, asset, user)
}

// Repay is a paid mutator transaction binding the contract method 0x573ade81.
//
// Solidity: function repay(address asset, uint256 amount, uint256 interestRateMode, address onBehalfOf) returns(uint256)
func (_Pool *PoolTransactor) Repay(opts *bind.TransactOpts, asset common.Address, amount *big.Int, interestRateMode *big.Int, onBehalfOf common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "repay", asset, amount, interestRateMode, onBehalfOf)
}

// Repay is a paid mutator transaction binding the contract method 0x573ade81.
//
// Solidity: function repay(address asset, uint256 amount, uint256 interestRateMode, address onBehalfOf) returns(uint256)
func (_Pool *PoolSession) Repay(asset common.Address, amount *big.Int, interestRateMode *big.Int, onBehalfOf common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Repay(&_Pool.TransactOpts, asset, amount, interestRateMode, onBehalfOf)
}

// Repay is a paid mutator transaction binding the contract method 0x573ade81.
//
// Solidity: function repay(address asset, uint256 amount, uint256 interestRateMode, address onBehalfOf) returns(uint256)
func (_Pool *PoolTransactorSession) Repay(asset common.Address, amount *big.Int, interestRateMode *big.Int, onBehalfOf common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Repay(&_Pool.TransactOpts, asset, amount, interestRateMode, onBehalfOf)
}

// RepayWithATokens is a paid mutator transaction binding the contract method 0x2dad97d4.
//
// Solidity: function repayWithATokens(address asset, uint256 amount, uint256 interestRateMode) returns(uint256)
func (_Pool *PoolTransactor) RepayWithATokens(opts *bind.TransactOpts, asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "repayWithATokens", asset, amount, interestRateMode)
}

// RepayWithATokens is a paid mutator transaction binding the contract method 0x2dad97d4.
//
// Solidity: function repayWithATokens(address asset, uint256 amount, uint256 interestRateMode) returns(uint256)
func (_Pool *PoolSession) RepayWithATokens(asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.RepayWithATokens(&_Pool.TransactOpts, asset, amount, interestRateMode)
}

// RepayWithATokens is a paid mutator transaction binding the contract method 0x2dad97d4.
//
// Solidity: function repayWithATokens(address asset, uint256 amount, uint256 interestRateMode) returns(uint256)
func (_Pool *PoolTransactorSession) RepayWithATokens(asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.RepayWithATokens(&_Pool.TransactOpts, asset, amount, interestRateMode)
}

// RepayWithPermit is a paid mutator transaction binding the contract method 0xee3e210b.
//
// Solidity: function repayWithPermit(address asset, uint256 amount, uint256 interestRateMode, address onBehalfOf, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS) returns(uint256)
func (_Pool *PoolTransactor) RepayWithPermit(opts *bind.TransactOpts, asset common.Address, amount *big.Int, interestRateMode *big.Int, onBehalfOf common.Address, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "repayWithPermit", asset, amount, interestRateMode, onBehalfOf, deadline, permitV, permitR, permitS)
}

// RepayWithPermit is a paid mutator transaction binding the contract method 0xee3e210b.
//
// Solidity: function repayWithPermit(address asset, uint256 amount, uint256 interestRateMode, address onBehalfOf, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS) returns(uint256)
func (_Pool *PoolSession) RepayWithPermit(asset common.Address, amount *big.Int, interestRateMode *big.Int, onBehalfOf common.Address, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte) (*types.Transaction, error) {
	return _Pool.Contract.RepayWithPermit(&_Pool.TransactOpts, asset, amount, interestRateMode, onBehalfOf, deadline, permitV, permitR, permitS)
}

// RepayWithPermit is a paid mutator transaction binding the contract method 0xee3e210b.
//
// Solidity: function repayWithPermit(address asset, uint256 amount, uint256 interestRateMode, address onBehalfOf, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS) returns(uint256)
func (_Pool *PoolTransactorSession) RepayWithPermit(asset common.Address, amount *big.Int, interestRateMode *big.Int, onBehalfOf common.Address, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte) (*types.Transaction, error) {
	return _Pool.Contract.RepayWithPermit(&_Pool.TransactOpts, asset, amount, interestRateMode, onBehalfOf, deadline, permitV, permitR, permitS)
}

// RescueTokens is a paid mutator transaction binding the contract method 0xcea9d26f.
//
// Solidity: function rescueTokens(address token, address to, uint256 amount) returns()
func (_Pool *PoolTransactor) RescueTokens(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "rescueTokens", token, to, amount)
}

// RescueTokens is a paid mutator transaction binding the contract method 0xcea9d26f.
//
// Solidity: function rescueTokens(address token, address to, uint256 amount) returns()
func (_Pool *PoolSession) RescueTokens(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.RescueTokens(&_Pool.TransactOpts, token, to, amount)
}

// RescueTokens is a paid mutator transaction binding the contract method 0xcea9d26f.
//
// Solidity: function rescueTokens(address token, address to, uint256 amount) returns()
func (_Pool *PoolTransactorSession) RescueTokens(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.RescueTokens(&_Pool.TransactOpts, token, to, amount)
}

// ResetIsolationModeTotalDebt is a paid mutator transaction binding the contract method 0xe43e88a1.
//
// Solidity: function resetIsolationModeTotalDebt(address asset) returns()
func (_Pool *PoolTransactor) ResetIsolationModeTotalDebt(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "resetIsolationModeTotalDebt", asset)
}

// ResetIsolationModeTotalDebt is a paid mutator transaction binding the contract method 0xe43e88a1.
//
// Solidity: function resetIsolationModeTotalDebt(address asset) returns()
func (_Pool *PoolSession) ResetIsolationModeTotalDebt(asset common.Address) (*types.Transaction, error) {
	return _Pool.Contract.ResetIsolationModeTotalDebt(&_Pool.TransactOpts, asset)
}

// ResetIsolationModeTotalDebt is a paid mutator transaction binding the contract method 0xe43e88a1.
//
// Solidity: function resetIsolationModeTotalDebt(address asset) returns()
func (_Pool *PoolTransactorSession) ResetIsolationModeTotalDebt(asset common.Address) (*types.Transaction, error) {
	return _Pool.Contract.ResetIsolationModeTotalDebt(&_Pool.TransactOpts, asset)
}

// SetConfiguration is a paid mutator transaction binding the contract method 0xf51e435b.
//
// Solidity: function setConfiguration(address asset, (uint256) configuration) returns()
func (_Pool *PoolTransactor) SetConfiguration(opts *bind.TransactOpts, asset common.Address, configuration DataTypesReserveConfigurationMap) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "setConfiguration", asset, configuration)
}

// SetConfiguration is a paid mutator transaction binding the contract method 0xf51e435b.
//
// Solidity: function setConfiguration(address asset, (uint256) configuration) returns()
func (_Pool *PoolSession) SetConfiguration(asset common.Address, configuration DataTypesReserveConfigurationMap) (*types.Transaction, error) {
	return _Pool.Contract.SetConfiguration(&_Pool.TransactOpts, asset, configuration)
}

// SetConfiguration is a paid mutator transaction binding the contract method 0xf51e435b.
//
// Solidity: function setConfiguration(address asset, (uint256) configuration) returns()
func (_Pool *PoolTransactorSession) SetConfiguration(asset common.Address, configuration DataTypesReserveConfigurationMap) (*types.Transaction, error) {
	return _Pool.Contract.SetConfiguration(&_Pool.TransactOpts, asset, configuration)
}

// SetReserveInterestRateStrategyAddress is a paid mutator transaction binding the contract method 0x1d2118f9.
//
// Solidity: function setReserveInterestRateStrategyAddress(address asset, address rateStrategyAddress) returns()
func (_Pool *PoolTransactor) SetReserveInterestRateStrategyAddress(opts *bind.TransactOpts, asset common.Address, rateStrategyAddress common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "setReserveInterestRateStrategyAddress", asset, rateStrategyAddress)
}

// SetReserveInterestRateStrategyAddress is a paid mutator transaction binding the contract method 0x1d2118f9.
//
// Solidity: function setReserveInterestRateStrategyAddress(address asset, address rateStrategyAddress) returns()
func (_Pool *PoolSession) SetReserveInterestRateStrategyAddress(asset common.Address, rateStrategyAddress common.Address) (*types.Transaction, error) {
	return _Pool.Contract.SetReserveInterestRateStrategyAddress(&_Pool.TransactOpts, asset, rateStrategyAddress)
}

// SetReserveInterestRateStrategyAddress is a paid mutator transaction binding the contract method 0x1d2118f9.
//
// Solidity: function setReserveInterestRateStrategyAddress(address asset, address rateStrategyAddress) returns()
func (_Pool *PoolTransactorSession) SetReserveInterestRateStrategyAddress(asset common.Address, rateStrategyAddress common.Address) (*types.Transaction, error) {
	return _Pool.Contract.SetReserveInterestRateStrategyAddress(&_Pool.TransactOpts, asset, rateStrategyAddress)
}

// SetUserEMode is a paid mutator transaction binding the contract method 0x28530a47.
//
// Solidity: function setUserEMode(uint8 categoryId) returns()
func (_Pool *PoolTransactor) SetUserEMode(opts *bind.TransactOpts, categoryId uint8) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "setUserEMode", categoryId)
}

// SetUserEMode is a paid mutator transaction binding the contract method 0x28530a47.
//
// Solidity: function setUserEMode(uint8 categoryId) returns()
func (_Pool *PoolSession) SetUserEMode(categoryId uint8) (*types.Transaction, error) {
	return _Pool.Contract.SetUserEMode(&_Pool.TransactOpts, categoryId)
}

// SetUserEMode is a paid mutator transaction binding the contract method 0x28530a47.
//
// Solidity: function setUserEMode(uint8 categoryId) returns()
func (_Pool *PoolTransactorSession) SetUserEMode(categoryId uint8) (*types.Transaction, error) {
	return _Pool.Contract.SetUserEMode(&_Pool.TransactOpts, categoryId)
}

// SetUserUseReserveAsCollateral is a paid mutator transaction binding the contract method 0x5a3b74b9.
//
// Solidity: function setUserUseReserveAsCollateral(address asset, bool useAsCollateral) returns()
func (_Pool *PoolTransactor) SetUserUseReserveAsCollateral(opts *bind.TransactOpts, asset common.Address, useAsCollateral bool) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "setUserUseReserveAsCollateral", asset, useAsCollateral)
}

// SetUserUseReserveAsCollateral is a paid mutator transaction binding the contract method 0x5a3b74b9.
//
// Solidity: function setUserUseReserveAsCollateral(address asset, bool useAsCollateral) returns()
func (_Pool *PoolSession) SetUserUseReserveAsCollateral(asset common.Address, useAsCollateral bool) (*types.Transaction, error) {
	return _Pool.Contract.SetUserUseReserveAsCollateral(&_Pool.TransactOpts, asset, useAsCollateral)
}

// SetUserUseReserveAsCollateral is a paid mutator transaction binding the contract method 0x5a3b74b9.
//
// Solidity: function setUserUseReserveAsCollateral(address asset, bool useAsCollateral) returns()
func (_Pool *PoolTransactorSession) SetUserUseReserveAsCollateral(asset common.Address, useAsCollateral bool) (*types.Transaction, error) {
	return _Pool.Contract.SetUserUseReserveAsCollateral(&_Pool.TransactOpts, asset, useAsCollateral)
}

// Supply is a paid mutator transaction binding the contract method 0x617ba037.
//
// Solidity: function supply(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolTransactor) Supply(opts *bind.TransactOpts, asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "supply", asset, amount, onBehalfOf, referralCode)
}

// Supply is a paid mutator transaction binding the contract method 0x617ba037.
//
// Solidity: function supply(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolSession) Supply(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.Supply(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode)
}

// Supply is a paid mutator transaction binding the contract method 0x617ba037.
//
// Solidity: function supply(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) returns()
func (_Pool *PoolTransactorSession) Supply(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16) (*types.Transaction, error) {
	return _Pool.Contract.Supply(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode)
}

// SupplyWithPermit is a paid mutator transaction binding the contract method 0x02c205f0.
//
// Solidity: function supplyWithPermit(address asset, uint256 amount, address onBehalfOf, uint16 referralCode, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS) returns()
func (_Pool *PoolTransactor) SupplyWithPermit(opts *bind.TransactOpts, asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "supplyWithPermit", asset, amount, onBehalfOf, referralCode, deadline, permitV, permitR, permitS)
}

// SupplyWithPermit is a paid mutator transaction binding the contract method 0x02c205f0.
//
// Solidity: function supplyWithPermit(address asset, uint256 amount, address onBehalfOf, uint16 referralCode, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS) returns()
func (_Pool *PoolSession) SupplyWithPermit(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte) (*types.Transaction, error) {
	return _Pool.Contract.SupplyWithPermit(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode, deadline, permitV, permitR, permitS)
}

// SupplyWithPermit is a paid mutator transaction binding the contract method 0x02c205f0.
//
// Solidity: function supplyWithPermit(address asset, uint256 amount, address onBehalfOf, uint16 referralCode, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS) returns()
func (_Pool *PoolTransactorSession) SupplyWithPermit(asset common.Address, amount *big.Int, onBehalfOf common.Address, referralCode uint16, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte) (*types.Transaction, error) {
	return _Pool.Contract.SupplyWithPermit(&_Pool.TransactOpts, asset, amount, onBehalfOf, referralCode, deadline, permitV, permitR, permitS)
}

// SwapBorrowRateMode is a paid mutator transaction binding the contract method 0x94ba89a2.
//
// Solidity: function swapBorrowRateMode(address asset, uint256 interestRateMode) returns()
func (_Pool *PoolTransactor) SwapBorrowRateMode(opts *bind.TransactOpts, asset common.Address, interestRateMode *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "swapBorrowRateMode", asset, interestRateMode)
}

// SwapBorrowRateMode is a paid mutator transaction binding the contract method 0x94ba89a2.
//
// Solidity: function swapBorrowRateMode(address asset, uint256 interestRateMode) returns()
func (_Pool *PoolSession) SwapBorrowRateMode(asset common.Address, interestRateMode *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.SwapBorrowRateMode(&_Pool.TransactOpts, asset, interestRateMode)
}

// SwapBorrowRateMode is a paid mutator transaction binding the contract method 0x94ba89a2.
//
// Solidity: function swapBorrowRateMode(address asset, uint256 interestRateMode) returns()
func (_Pool *PoolTransactorSession) SwapBorrowRateMode(asset common.Address, interestRateMode *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.SwapBorrowRateMode(&_Pool.TransactOpts, asset, interestRateMode)
}

// UpdateBridgeProtocolFee is a paid mutator transaction binding the contract method 0x3036b439.
//
// Solidity: function updateBridgeProtocolFee(uint256 protocolFee) returns()
func (_Pool *PoolTransactor) UpdateBridgeProtocolFee(opts *bind.TransactOpts, protocolFee *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "updateBridgeProtocolFee", protocolFee)
}

// UpdateBridgeProtocolFee is a paid mutator transaction binding the contract method 0x3036b439.
//
// Solidity: function updateBridgeProtocolFee(uint256 protocolFee) returns()
func (_Pool *PoolSession) UpdateBridgeProtocolFee(protocolFee *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.UpdateBridgeProtocolFee(&_Pool.TransactOpts, protocolFee)
}

// UpdateBridgeProtocolFee is a paid mutator transaction binding the contract method 0x3036b439.
//
// Solidity: function updateBridgeProtocolFee(uint256 protocolFee) returns()
func (_Pool *PoolTransactorSession) UpdateBridgeProtocolFee(protocolFee *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.UpdateBridgeProtocolFee(&_Pool.TransactOpts, protocolFee)
}

// UpdateFlashloanPremiums is a paid mutator transaction binding the contract method 0xbcb6e522.
//
// Solidity: function updateFlashloanPremiums(uint128 flashLoanPremiumTotal, uint128 flashLoanPremiumToProtocol) returns()
func (_Pool *PoolTransactor) UpdateFlashloanPremiums(opts *bind.TransactOpts, flashLoanPremiumTotal *big.Int, flashLoanPremiumToProtocol *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "updateFlashloanPremiums", flashLoanPremiumTotal, flashLoanPremiumToProtocol)
}

// UpdateFlashloanPremiums is a paid mutator transaction binding the contract method 0xbcb6e522.
//
// Solidity: function updateFlashloanPremiums(uint128 flashLoanPremiumTotal, uint128 flashLoanPremiumToProtocol) returns()
func (_Pool *PoolSession) UpdateFlashloanPremiums(flashLoanPremiumTotal *big.Int, flashLoanPremiumToProtocol *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.UpdateFlashloanPremiums(&_Pool.TransactOpts, flashLoanPremiumTotal, flashLoanPremiumToProtocol)
}

// UpdateFlashloanPremiums is a paid mutator transaction binding the contract method 0xbcb6e522.
//
// Solidity: function updateFlashloanPremiums(uint128 flashLoanPremiumTotal, uint128 flashLoanPremiumToProtocol) returns()
func (_Pool *PoolTransactorSession) UpdateFlashloanPremiums(flashLoanPremiumTotal *big.Int, flashLoanPremiumToProtocol *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.UpdateFlashloanPremiums(&_Pool.TransactOpts, flashLoanPremiumTotal, flashLoanPremiumToProtocol)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address asset, uint256 amount, address to) returns(uint256)
func (_Pool *PoolTransactor) Withdraw(opts *bind.TransactOpts, asset common.Address, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "withdraw", asset, amount, to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address asset, uint256 amount, address to) returns(uint256)
func (_Pool *PoolSession) Withdraw(asset common.Address, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Withdraw(&_Pool.TransactOpts, asset, amount, to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address asset, uint256 amount, address to) returns(uint256)
func (_Pool *PoolTransactorSession) Withdraw(asset common.Address, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _Pool.Contract.Withdraw(&_Pool.TransactOpts, asset, amount, to)
}

// PoolBackUnbackedIterator is returned from FilterBackUnbacked and is used to iterate over the raw logs and unpacked data for BackUnbacked events raised by the Pool contract.
type PoolBackUnbackedIterator struct {
	Event *PoolBackUnbacked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolBackUnbackedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolBackUnbacked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolBackUnbacked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolBackUnbackedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolBackUnbackedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolBackUnbacked represents a BackUnbacked event raised by the Pool contract.
type PoolBackUnbacked struct {
	Reserve common.Address
	Backer  common.Address
	Amount  *big.Int
	Fee     *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBackUnbacked is a free log retrieval operation binding the contract event 0x281596e92b2d974beb7d4f124df30a0b39067b096893e95011ce4bdad798b759.
//
// Solidity: event BackUnbacked(address indexed reserve, address indexed backer, uint256 amount, uint256 fee)
func (_Pool *PoolFilterer) FilterBackUnbacked(opts *bind.FilterOpts, reserve []common.Address, backer []common.Address) (*PoolBackUnbackedIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var backerRule []interface{}
	for _, backerItem := range backer {
		backerRule = append(backerRule, backerItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "BackUnbacked", reserveRule, backerRule)
	if err != nil {
		return nil, err
	}
	return &PoolBackUnbackedIterator{contract: _Pool.contract, event: "BackUnbacked", logs: logs, sub: sub}, nil
}

// WatchBackUnbacked is a free log subscription operation binding the contract event 0x281596e92b2d974beb7d4f124df30a0b39067b096893e95011ce4bdad798b759.
//
// Solidity: event BackUnbacked(address indexed reserve, address indexed backer, uint256 amount, uint256 fee)
func (_Pool *PoolFilterer) WatchBackUnbacked(opts *bind.WatchOpts, sink chan<- *PoolBackUnbacked, reserve []common.Address, backer []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var backerRule []interface{}
	for _, backerItem := range backer {
		backerRule = append(backerRule, backerItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "BackUnbacked", reserveRule, backerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolBackUnbacked)
				if err := _Pool.contract.UnpackLog(event, "BackUnbacked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBackUnbacked is a log parse operation binding the contract event 0x281596e92b2d974beb7d4f124df30a0b39067b096893e95011ce4bdad798b759.
//
// Solidity: event BackUnbacked(address indexed reserve, address indexed backer, uint256 amount, uint256 fee)
func (_Pool *PoolFilterer) ParseBackUnbacked(log types.Log) (*PoolBackUnbacked, error) {
	event := new(PoolBackUnbacked)
	if err := _Pool.contract.UnpackLog(event, "BackUnbacked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolBorrowIterator is returned from FilterBorrow and is used to iterate over the raw logs and unpacked data for Borrow events raised by the Pool contract.
type PoolBorrowIterator struct {
	Event *PoolBorrow // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolBorrowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolBorrow)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolBorrow)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolBorrowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolBorrowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolBorrow represents a Borrow event raised by the Pool contract.
type PoolBorrow struct {
	Reserve          common.Address
	User             common.Address
	OnBehalfOf       common.Address
	Amount           *big.Int
	InterestRateMode uint8
	BorrowRate       *big.Int
	ReferralCode     uint16
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBorrow is a free log retrieval operation binding the contract event 0xb3d084820fb1a9decffb176436bd02558d15fac9b0ddfed8c465bc7359d7dce0.
//
// Solidity: event Borrow(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint8 interestRateMode, uint256 borrowRate, uint16 indexed referralCode)
func (_Pool *PoolFilterer) FilterBorrow(opts *bind.FilterOpts, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (*PoolBorrowIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "Borrow", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return &PoolBorrowIterator{contract: _Pool.contract, event: "Borrow", logs: logs, sub: sub}, nil
}

// WatchBorrow is a free log subscription operation binding the contract event 0xb3d084820fb1a9decffb176436bd02558d15fac9b0ddfed8c465bc7359d7dce0.
//
// Solidity: event Borrow(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint8 interestRateMode, uint256 borrowRate, uint16 indexed referralCode)
func (_Pool *PoolFilterer) WatchBorrow(opts *bind.WatchOpts, sink chan<- *PoolBorrow, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "Borrow", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolBorrow)
				if err := _Pool.contract.UnpackLog(event, "Borrow", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBorrow is a log parse operation binding the contract event 0xb3d084820fb1a9decffb176436bd02558d15fac9b0ddfed8c465bc7359d7dce0.
//
// Solidity: event Borrow(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint8 interestRateMode, uint256 borrowRate, uint16 indexed referralCode)
func (_Pool *PoolFilterer) ParseBorrow(log types.Log) (*PoolBorrow, error) {
	event := new(PoolBorrow)
	if err := _Pool.contract.UnpackLog(event, "Borrow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolFlashLoanIterator is returned from FilterFlashLoan and is used to iterate over the raw logs and unpacked data for FlashLoan events raised by the Pool contract.
type PoolFlashLoanIterator struct {
	Event *PoolFlashLoan // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolFlashLoanIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolFlashLoan)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolFlashLoan)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolFlashLoanIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolFlashLoanIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolFlashLoan represents a FlashLoan event raised by the Pool contract.
type PoolFlashLoan struct {
	Target           common.Address
	Initiator        common.Address
	Asset            common.Address
	Amount           *big.Int
	InterestRateMode uint8
	Premium          *big.Int
	ReferralCode     uint16
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterFlashLoan is a free log retrieval operation binding the contract event 0xefefaba5e921573100900a3ad9cf29f222d995fb3b6045797eaea7521bd8d6f0.
//
// Solidity: event FlashLoan(address indexed target, address initiator, address indexed asset, uint256 amount, uint8 interestRateMode, uint256 premium, uint16 indexed referralCode)
func (_Pool *PoolFilterer) FilterFlashLoan(opts *bind.FilterOpts, target []common.Address, asset []common.Address, referralCode []uint16) (*PoolFlashLoanIterator, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "FlashLoan", targetRule, assetRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return &PoolFlashLoanIterator{contract: _Pool.contract, event: "FlashLoan", logs: logs, sub: sub}, nil
}

// WatchFlashLoan is a free log subscription operation binding the contract event 0xefefaba5e921573100900a3ad9cf29f222d995fb3b6045797eaea7521bd8d6f0.
//
// Solidity: event FlashLoan(address indexed target, address initiator, address indexed asset, uint256 amount, uint8 interestRateMode, uint256 premium, uint16 indexed referralCode)
func (_Pool *PoolFilterer) WatchFlashLoan(opts *bind.WatchOpts, sink chan<- *PoolFlashLoan, target []common.Address, asset []common.Address, referralCode []uint16) (event.Subscription, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "FlashLoan", targetRule, assetRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolFlashLoan)
				if err := _Pool.contract.UnpackLog(event, "FlashLoan", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFlashLoan is a log parse operation binding the contract event 0xefefaba5e921573100900a3ad9cf29f222d995fb3b6045797eaea7521bd8d6f0.
//
// Solidity: event FlashLoan(address indexed target, address initiator, address indexed asset, uint256 amount, uint8 interestRateMode, uint256 premium, uint16 indexed referralCode)
func (_Pool *PoolFilterer) ParseFlashLoan(log types.Log) (*PoolFlashLoan, error) {
	event := new(PoolFlashLoan)
	if err := _Pool.contract.UnpackLog(event, "FlashLoan", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolIsolationModeTotalDebtUpdatedIterator is returned from FilterIsolationModeTotalDebtUpdated and is used to iterate over the raw logs and unpacked data for IsolationModeTotalDebtUpdated events raised by the Pool contract.
type PoolIsolationModeTotalDebtUpdatedIterator struct {
	Event *PoolIsolationModeTotalDebtUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolIsolationModeTotalDebtUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolIsolationModeTotalDebtUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolIsolationModeTotalDebtUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolIsolationModeTotalDebtUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolIsolationModeTotalDebtUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolIsolationModeTotalDebtUpdated represents a IsolationModeTotalDebtUpdated event raised by the Pool contract.
type PoolIsolationModeTotalDebtUpdated struct {
	Asset     common.Address
	TotalDebt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIsolationModeTotalDebtUpdated is a free log retrieval operation binding the contract event 0xaef84d3b40895fd58c561f3998000f0583abb992a52fbdc99ace8e8de4d676a5.
//
// Solidity: event IsolationModeTotalDebtUpdated(address indexed asset, uint256 totalDebt)
func (_Pool *PoolFilterer) FilterIsolationModeTotalDebtUpdated(opts *bind.FilterOpts, asset []common.Address) (*PoolIsolationModeTotalDebtUpdatedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "IsolationModeTotalDebtUpdated", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolIsolationModeTotalDebtUpdatedIterator{contract: _Pool.contract, event: "IsolationModeTotalDebtUpdated", logs: logs, sub: sub}, nil
}

// WatchIsolationModeTotalDebtUpdated is a free log subscription operation binding the contract event 0xaef84d3b40895fd58c561f3998000f0583abb992a52fbdc99ace8e8de4d676a5.
//
// Solidity: event IsolationModeTotalDebtUpdated(address indexed asset, uint256 totalDebt)
func (_Pool *PoolFilterer) WatchIsolationModeTotalDebtUpdated(opts *bind.WatchOpts, sink chan<- *PoolIsolationModeTotalDebtUpdated, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "IsolationModeTotalDebtUpdated", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolIsolationModeTotalDebtUpdated)
				if err := _Pool.contract.UnpackLog(event, "IsolationModeTotalDebtUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIsolationModeTotalDebtUpdated is a log parse operation binding the contract event 0xaef84d3b40895fd58c561f3998000f0583abb992a52fbdc99ace8e8de4d676a5.
//
// Solidity: event IsolationModeTotalDebtUpdated(address indexed asset, uint256 totalDebt)
func (_Pool *PoolFilterer) ParseIsolationModeTotalDebtUpdated(log types.Log) (*PoolIsolationModeTotalDebtUpdated, error) {
	event := new(PoolIsolationModeTotalDebtUpdated)
	if err := _Pool.contract.UnpackLog(event, "IsolationModeTotalDebtUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolLiquidationCallIterator is returned from FilterLiquidationCall and is used to iterate over the raw logs and unpacked data for LiquidationCall events raised by the Pool contract.
type PoolLiquidationCallIterator struct {
	Event *PoolLiquidationCall // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolLiquidationCallIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolLiquidationCall)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolLiquidationCall)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolLiquidationCallIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolLiquidationCallIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolLiquidationCall represents a LiquidationCall event raised by the Pool contract.
type PoolLiquidationCall struct {
	CollateralAsset            common.Address
	DebtAsset                  common.Address
	User                       common.Address
	DebtToCover                *big.Int
	LiquidatedCollateralAmount *big.Int
	Liquidator                 common.Address
	ReceiveAToken              bool
	Raw                        types.Log // Blockchain specific contextual infos
}

// FilterLiquidationCall is a free log retrieval operation binding the contract event 0xe413a321e8681d831f4dbccbca790d2952b56f977908e45be37335533e005286.
//
// Solidity: event LiquidationCall(address indexed collateralAsset, address indexed debtAsset, address indexed user, uint256 debtToCover, uint256 liquidatedCollateralAmount, address liquidator, bool receiveAToken)
func (_Pool *PoolFilterer) FilterLiquidationCall(opts *bind.FilterOpts, collateralAsset []common.Address, debtAsset []common.Address, user []common.Address) (*PoolLiquidationCallIterator, error) {

	var collateralAssetRule []interface{}
	for _, collateralAssetItem := range collateralAsset {
		collateralAssetRule = append(collateralAssetRule, collateralAssetItem)
	}
	var debtAssetRule []interface{}
	for _, debtAssetItem := range debtAsset {
		debtAssetRule = append(debtAssetRule, debtAssetItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "LiquidationCall", collateralAssetRule, debtAssetRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PoolLiquidationCallIterator{contract: _Pool.contract, event: "LiquidationCall", logs: logs, sub: sub}, nil
}

// WatchLiquidationCall is a free log subscription operation binding the contract event 0xe413a321e8681d831f4dbccbca790d2952b56f977908e45be37335533e005286.
//
// Solidity: event LiquidationCall(address indexed collateralAsset, address indexed debtAsset, address indexed user, uint256 debtToCover, uint256 liquidatedCollateralAmount, address liquidator, bool receiveAToken)
func (_Pool *PoolFilterer) WatchLiquidationCall(opts *bind.WatchOpts, sink chan<- *PoolLiquidationCall, collateralAsset []common.Address, debtAsset []common.Address, user []common.Address) (event.Subscription, error) {

	var collateralAssetRule []interface{}
	for _, collateralAssetItem := range collateralAsset {
		collateralAssetRule = append(collateralAssetRule, collateralAssetItem)
	}
	var debtAssetRule []interface{}
	for _, debtAssetItem := range debtAsset {
		debtAssetRule = append(debtAssetRule, debtAssetItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "LiquidationCall", collateralAssetRule, debtAssetRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolLiquidationCall)
				if err := _Pool.contract.UnpackLog(event, "LiquidationCall", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLiquidationCall is a log parse operation binding the contract event 0xe413a321e8681d831f4dbccbca790d2952b56f977908e45be37335533e005286.
//
// Solidity: event LiquidationCall(address indexed collateralAsset, address indexed debtAsset, address indexed user, uint256 debtToCover, uint256 liquidatedCollateralAmount, address liquidator, bool receiveAToken)
func (_Pool *PoolFilterer) ParseLiquidationCall(log types.Log) (*PoolLiquidationCall, error) {
	event := new(PoolLiquidationCall)
	if err := _Pool.contract.UnpackLog(event, "LiquidationCall", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolMintUnbackedIterator is returned from FilterMintUnbacked and is used to iterate over the raw logs and unpacked data for MintUnbacked events raised by the Pool contract.
type PoolMintUnbackedIterator struct {
	Event *PoolMintUnbacked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolMintUnbackedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolMintUnbacked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolMintUnbacked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolMintUnbackedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolMintUnbackedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolMintUnbacked represents a MintUnbacked event raised by the Pool contract.
type PoolMintUnbacked struct {
	Reserve      common.Address
	User         common.Address
	OnBehalfOf   common.Address
	Amount       *big.Int
	ReferralCode uint16
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMintUnbacked is a free log retrieval operation binding the contract event 0xf25af37b3d3ec226063dc9bdc103ece7eb110a50f340fe854bb7bc1b0676d7d0.
//
// Solidity: event MintUnbacked(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_Pool *PoolFilterer) FilterMintUnbacked(opts *bind.FilterOpts, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (*PoolMintUnbackedIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "MintUnbacked", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return &PoolMintUnbackedIterator{contract: _Pool.contract, event: "MintUnbacked", logs: logs, sub: sub}, nil
}

// WatchMintUnbacked is a free log subscription operation binding the contract event 0xf25af37b3d3ec226063dc9bdc103ece7eb110a50f340fe854bb7bc1b0676d7d0.
//
// Solidity: event MintUnbacked(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_Pool *PoolFilterer) WatchMintUnbacked(opts *bind.WatchOpts, sink chan<- *PoolMintUnbacked, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "MintUnbacked", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolMintUnbacked)
				if err := _Pool.contract.UnpackLog(event, "MintUnbacked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMintUnbacked is a log parse operation binding the contract event 0xf25af37b3d3ec226063dc9bdc103ece7eb110a50f340fe854bb7bc1b0676d7d0.
//
// Solidity: event MintUnbacked(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_Pool *PoolFilterer) ParseMintUnbacked(log types.Log) (*PoolMintUnbacked, error) {
	event := new(PoolMintUnbacked)
	if err := _Pool.contract.UnpackLog(event, "MintUnbacked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolMintedToTreasuryIterator is returned from FilterMintedToTreasury and is used to iterate over the raw logs and unpacked data for MintedToTreasury events raised by the Pool contract.
type PoolMintedToTreasuryIterator struct {
	Event *PoolMintedToTreasury // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolMintedToTreasuryIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolMintedToTreasury)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolMintedToTreasury)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolMintedToTreasuryIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolMintedToTreasuryIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolMintedToTreasury represents a MintedToTreasury event raised by the Pool contract.
type PoolMintedToTreasury struct {
	Reserve      common.Address
	AmountMinted *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMintedToTreasury is a free log retrieval operation binding the contract event 0xbfa21aa5d5f9a1f0120a95e7c0749f389863cbdbfff531aa7339077a5bc919de.
//
// Solidity: event MintedToTreasury(address indexed reserve, uint256 amountMinted)
func (_Pool *PoolFilterer) FilterMintedToTreasury(opts *bind.FilterOpts, reserve []common.Address) (*PoolMintedToTreasuryIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "MintedToTreasury", reserveRule)
	if err != nil {
		return nil, err
	}
	return &PoolMintedToTreasuryIterator{contract: _Pool.contract, event: "MintedToTreasury", logs: logs, sub: sub}, nil
}

// WatchMintedToTreasury is a free log subscription operation binding the contract event 0xbfa21aa5d5f9a1f0120a95e7c0749f389863cbdbfff531aa7339077a5bc919de.
//
// Solidity: event MintedToTreasury(address indexed reserve, uint256 amountMinted)
func (_Pool *PoolFilterer) WatchMintedToTreasury(opts *bind.WatchOpts, sink chan<- *PoolMintedToTreasury, reserve []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "MintedToTreasury", reserveRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolMintedToTreasury)
				if err := _Pool.contract.UnpackLog(event, "MintedToTreasury", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMintedToTreasury is a log parse operation binding the contract event 0xbfa21aa5d5f9a1f0120a95e7c0749f389863cbdbfff531aa7339077a5bc919de.
//
// Solidity: event MintedToTreasury(address indexed reserve, uint256 amountMinted)
func (_Pool *PoolFilterer) ParseMintedToTreasury(log types.Log) (*PoolMintedToTreasury, error) {
	event := new(PoolMintedToTreasury)
	if err := _Pool.contract.UnpackLog(event, "MintedToTreasury", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRebalanceStableBorrowRateIterator is returned from FilterRebalanceStableBorrowRate and is used to iterate over the raw logs and unpacked data for RebalanceStableBorrowRate events raised by the Pool contract.
type PoolRebalanceStableBorrowRateIterator struct {
	Event *PoolRebalanceStableBorrowRate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolRebalanceStableBorrowRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRebalanceStableBorrowRate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolRebalanceStableBorrowRate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolRebalanceStableBorrowRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRebalanceStableBorrowRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRebalanceStableBorrowRate represents a RebalanceStableBorrowRate event raised by the Pool contract.
type PoolRebalanceStableBorrowRate struct {
	Reserve common.Address
	User    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRebalanceStableBorrowRate is a free log retrieval operation binding the contract event 0x9f439ae0c81e41a04d3fdfe07aed54e6a179fb0db15be7702eb66fa8ef6f5300.
//
// Solidity: event RebalanceStableBorrowRate(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) FilterRebalanceStableBorrowRate(opts *bind.FilterOpts, reserve []common.Address, user []common.Address) (*PoolRebalanceStableBorrowRateIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "RebalanceStableBorrowRate", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PoolRebalanceStableBorrowRateIterator{contract: _Pool.contract, event: "RebalanceStableBorrowRate", logs: logs, sub: sub}, nil
}

// WatchRebalanceStableBorrowRate is a free log subscription operation binding the contract event 0x9f439ae0c81e41a04d3fdfe07aed54e6a179fb0db15be7702eb66fa8ef6f5300.
//
// Solidity: event RebalanceStableBorrowRate(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) WatchRebalanceStableBorrowRate(opts *bind.WatchOpts, sink chan<- *PoolRebalanceStableBorrowRate, reserve []common.Address, user []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "RebalanceStableBorrowRate", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRebalanceStableBorrowRate)
				if err := _Pool.contract.UnpackLog(event, "RebalanceStableBorrowRate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRebalanceStableBorrowRate is a log parse operation binding the contract event 0x9f439ae0c81e41a04d3fdfe07aed54e6a179fb0db15be7702eb66fa8ef6f5300.
//
// Solidity: event RebalanceStableBorrowRate(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) ParseRebalanceStableBorrowRate(log types.Log) (*PoolRebalanceStableBorrowRate, error) {
	event := new(PoolRebalanceStableBorrowRate)
	if err := _Pool.contract.UnpackLog(event, "RebalanceStableBorrowRate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolRepayIterator is returned from FilterRepay and is used to iterate over the raw logs and unpacked data for Repay events raised by the Pool contract.
type PoolRepayIterator struct {
	Event *PoolRepay // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolRepayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolRepay)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolRepay)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolRepayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolRepayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolRepay represents a Repay event raised by the Pool contract.
type PoolRepay struct {
	Reserve    common.Address
	User       common.Address
	Repayer    common.Address
	Amount     *big.Int
	UseATokens bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRepay is a free log retrieval operation binding the contract event 0xa534c8dbe71f871f9f3530e97a74601fea17b426cae02e1c5aee42c96c784051.
//
// Solidity: event Repay(address indexed reserve, address indexed user, address indexed repayer, uint256 amount, bool useATokens)
func (_Pool *PoolFilterer) FilterRepay(opts *bind.FilterOpts, reserve []common.Address, user []common.Address, repayer []common.Address) (*PoolRepayIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var repayerRule []interface{}
	for _, repayerItem := range repayer {
		repayerRule = append(repayerRule, repayerItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "Repay", reserveRule, userRule, repayerRule)
	if err != nil {
		return nil, err
	}
	return &PoolRepayIterator{contract: _Pool.contract, event: "Repay", logs: logs, sub: sub}, nil
}

// WatchRepay is a free log subscription operation binding the contract event 0xa534c8dbe71f871f9f3530e97a74601fea17b426cae02e1c5aee42c96c784051.
//
// Solidity: event Repay(address indexed reserve, address indexed user, address indexed repayer, uint256 amount, bool useATokens)
func (_Pool *PoolFilterer) WatchRepay(opts *bind.WatchOpts, sink chan<- *PoolRepay, reserve []common.Address, user []common.Address, repayer []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var repayerRule []interface{}
	for _, repayerItem := range repayer {
		repayerRule = append(repayerRule, repayerItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "Repay", reserveRule, userRule, repayerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolRepay)
				if err := _Pool.contract.UnpackLog(event, "Repay", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRepay is a log parse operation binding the contract event 0xa534c8dbe71f871f9f3530e97a74601fea17b426cae02e1c5aee42c96c784051.
//
// Solidity: event Repay(address indexed reserve, address indexed user, address indexed repayer, uint256 amount, bool useATokens)
func (_Pool *PoolFilterer) ParseRepay(log types.Log) (*PoolRepay, error) {
	event := new(PoolRepay)
	if err := _Pool.contract.UnpackLog(event, "Repay", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolReserveDataUpdatedIterator is returned from FilterReserveDataUpdated and is used to iterate over the raw logs and unpacked data for ReserveDataUpdated events raised by the Pool contract.
type PoolReserveDataUpdatedIterator struct {
	Event *PoolReserveDataUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolReserveDataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolReserveDataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolReserveDataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolReserveDataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolReserveDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolReserveDataUpdated represents a ReserveDataUpdated event raised by the Pool contract.
type PoolReserveDataUpdated struct {
	Reserve             common.Address
	LiquidityRate       *big.Int
	StableBorrowRate    *big.Int
	VariableBorrowRate  *big.Int
	LiquidityIndex      *big.Int
	VariableBorrowIndex *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterReserveDataUpdated is a free log retrieval operation binding the contract event 0x804c9b842b2748a22bb64b345453a3de7ca54a6ca45ce00d415894979e22897a.
//
// Solidity: event ReserveDataUpdated(address indexed reserve, uint256 liquidityRate, uint256 stableBorrowRate, uint256 variableBorrowRate, uint256 liquidityIndex, uint256 variableBorrowIndex)
func (_Pool *PoolFilterer) FilterReserveDataUpdated(opts *bind.FilterOpts, reserve []common.Address) (*PoolReserveDataUpdatedIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "ReserveDataUpdated", reserveRule)
	if err != nil {
		return nil, err
	}
	return &PoolReserveDataUpdatedIterator{contract: _Pool.contract, event: "ReserveDataUpdated", logs: logs, sub: sub}, nil
}

// WatchReserveDataUpdated is a free log subscription operation binding the contract event 0x804c9b842b2748a22bb64b345453a3de7ca54a6ca45ce00d415894979e22897a.
//
// Solidity: event ReserveDataUpdated(address indexed reserve, uint256 liquidityRate, uint256 stableBorrowRate, uint256 variableBorrowRate, uint256 liquidityIndex, uint256 variableBorrowIndex)
func (_Pool *PoolFilterer) WatchReserveDataUpdated(opts *bind.WatchOpts, sink chan<- *PoolReserveDataUpdated, reserve []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "ReserveDataUpdated", reserveRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolReserveDataUpdated)
				if err := _Pool.contract.UnpackLog(event, "ReserveDataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReserveDataUpdated is a log parse operation binding the contract event 0x804c9b842b2748a22bb64b345453a3de7ca54a6ca45ce00d415894979e22897a.
//
// Solidity: event ReserveDataUpdated(address indexed reserve, uint256 liquidityRate, uint256 stableBorrowRate, uint256 variableBorrowRate, uint256 liquidityIndex, uint256 variableBorrowIndex)
func (_Pool *PoolFilterer) ParseReserveDataUpdated(log types.Log) (*PoolReserveDataUpdated, error) {
	event := new(PoolReserveDataUpdated)
	if err := _Pool.contract.UnpackLog(event, "ReserveDataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolReserveUsedAsCollateralDisabledIterator is returned from FilterReserveUsedAsCollateralDisabled and is used to iterate over the raw logs and unpacked data for ReserveUsedAsCollateralDisabled events raised by the Pool contract.
type PoolReserveUsedAsCollateralDisabledIterator struct {
	Event *PoolReserveUsedAsCollateralDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolReserveUsedAsCollateralDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolReserveUsedAsCollateralDisabled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolReserveUsedAsCollateralDisabled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolReserveUsedAsCollateralDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolReserveUsedAsCollateralDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolReserveUsedAsCollateralDisabled represents a ReserveUsedAsCollateralDisabled event raised by the Pool contract.
type PoolReserveUsedAsCollateralDisabled struct {
	Reserve common.Address
	User    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveUsedAsCollateralDisabled is a free log retrieval operation binding the contract event 0x44c58d81365b66dd4b1a7f36c25aa97b8c71c361ee4937adc1a00000227db5dd.
//
// Solidity: event ReserveUsedAsCollateralDisabled(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) FilterReserveUsedAsCollateralDisabled(opts *bind.FilterOpts, reserve []common.Address, user []common.Address) (*PoolReserveUsedAsCollateralDisabledIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "ReserveUsedAsCollateralDisabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PoolReserveUsedAsCollateralDisabledIterator{contract: _Pool.contract, event: "ReserveUsedAsCollateralDisabled", logs: logs, sub: sub}, nil
}

// WatchReserveUsedAsCollateralDisabled is a free log subscription operation binding the contract event 0x44c58d81365b66dd4b1a7f36c25aa97b8c71c361ee4937adc1a00000227db5dd.
//
// Solidity: event ReserveUsedAsCollateralDisabled(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) WatchReserveUsedAsCollateralDisabled(opts *bind.WatchOpts, sink chan<- *PoolReserveUsedAsCollateralDisabled, reserve []common.Address, user []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "ReserveUsedAsCollateralDisabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolReserveUsedAsCollateralDisabled)
				if err := _Pool.contract.UnpackLog(event, "ReserveUsedAsCollateralDisabled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReserveUsedAsCollateralDisabled is a log parse operation binding the contract event 0x44c58d81365b66dd4b1a7f36c25aa97b8c71c361ee4937adc1a00000227db5dd.
//
// Solidity: event ReserveUsedAsCollateralDisabled(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) ParseReserveUsedAsCollateralDisabled(log types.Log) (*PoolReserveUsedAsCollateralDisabled, error) {
	event := new(PoolReserveUsedAsCollateralDisabled)
	if err := _Pool.contract.UnpackLog(event, "ReserveUsedAsCollateralDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolReserveUsedAsCollateralEnabledIterator is returned from FilterReserveUsedAsCollateralEnabled and is used to iterate over the raw logs and unpacked data for ReserveUsedAsCollateralEnabled events raised by the Pool contract.
type PoolReserveUsedAsCollateralEnabledIterator struct {
	Event *PoolReserveUsedAsCollateralEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolReserveUsedAsCollateralEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolReserveUsedAsCollateralEnabled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolReserveUsedAsCollateralEnabled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolReserveUsedAsCollateralEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolReserveUsedAsCollateralEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolReserveUsedAsCollateralEnabled represents a ReserveUsedAsCollateralEnabled event raised by the Pool contract.
type PoolReserveUsedAsCollateralEnabled struct {
	Reserve common.Address
	User    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveUsedAsCollateralEnabled is a free log retrieval operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) FilterReserveUsedAsCollateralEnabled(opts *bind.FilterOpts, reserve []common.Address, user []common.Address) (*PoolReserveUsedAsCollateralEnabledIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "ReserveUsedAsCollateralEnabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PoolReserveUsedAsCollateralEnabledIterator{contract: _Pool.contract, event: "ReserveUsedAsCollateralEnabled", logs: logs, sub: sub}, nil
}

// WatchReserveUsedAsCollateralEnabled is a free log subscription operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) WatchReserveUsedAsCollateralEnabled(opts *bind.WatchOpts, sink chan<- *PoolReserveUsedAsCollateralEnabled, reserve []common.Address, user []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "ReserveUsedAsCollateralEnabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolReserveUsedAsCollateralEnabled)
				if err := _Pool.contract.UnpackLog(event, "ReserveUsedAsCollateralEnabled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReserveUsedAsCollateralEnabled is a log parse operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_Pool *PoolFilterer) ParseReserveUsedAsCollateralEnabled(log types.Log) (*PoolReserveUsedAsCollateralEnabled, error) {
	event := new(PoolReserveUsedAsCollateralEnabled)
	if err := _Pool.contract.UnpackLog(event, "ReserveUsedAsCollateralEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolSupplyIterator is returned from FilterSupply and is used to iterate over the raw logs and unpacked data for Supply events raised by the Pool contract.
type PoolSupplyIterator struct {
	Event *PoolSupply // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolSupplyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolSupply)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolSupply)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolSupplyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolSupplyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolSupply represents a Supply event raised by the Pool contract.
type PoolSupply struct {
	Reserve      common.Address
	User         common.Address
	OnBehalfOf   common.Address
	Amount       *big.Int
	ReferralCode uint16
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSupply is a free log retrieval operation binding the contract event 0x2b627736bca15cd5381dcf80b0bf11fd197d01a037c52b927a881a10fb73ba61.
//
// Solidity: event Supply(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_Pool *PoolFilterer) FilterSupply(opts *bind.FilterOpts, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (*PoolSupplyIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "Supply", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return &PoolSupplyIterator{contract: _Pool.contract, event: "Supply", logs: logs, sub: sub}, nil
}

// WatchSupply is a free log subscription operation binding the contract event 0x2b627736bca15cd5381dcf80b0bf11fd197d01a037c52b927a881a10fb73ba61.
//
// Solidity: event Supply(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_Pool *PoolFilterer) WatchSupply(opts *bind.WatchOpts, sink chan<- *PoolSupply, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "Supply", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolSupply)
				if err := _Pool.contract.UnpackLog(event, "Supply", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSupply is a log parse operation binding the contract event 0x2b627736bca15cd5381dcf80b0bf11fd197d01a037c52b927a881a10fb73ba61.
//
// Solidity: event Supply(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_Pool *PoolFilterer) ParseSupply(log types.Log) (*PoolSupply, error) {
	event := new(PoolSupply)
	if err := _Pool.contract.UnpackLog(event, "Supply", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolSwapBorrowRateModeIterator is returned from FilterSwapBorrowRateMode and is used to iterate over the raw logs and unpacked data for SwapBorrowRateMode events raised by the Pool contract.
type PoolSwapBorrowRateModeIterator struct {
	Event *PoolSwapBorrowRateMode // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolSwapBorrowRateModeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolSwapBorrowRateMode)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolSwapBorrowRateMode)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolSwapBorrowRateModeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolSwapBorrowRateModeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolSwapBorrowRateMode represents a SwapBorrowRateMode event raised by the Pool contract.
type PoolSwapBorrowRateMode struct {
	Reserve          common.Address
	User             common.Address
	InterestRateMode uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSwapBorrowRateMode is a free log retrieval operation binding the contract event 0x7962b394d85a534033ba2efcf43cd36de57b7ebeb3de0ca4428965d9b3ddc481.
//
// Solidity: event SwapBorrowRateMode(address indexed reserve, address indexed user, uint8 interestRateMode)
func (_Pool *PoolFilterer) FilterSwapBorrowRateMode(opts *bind.FilterOpts, reserve []common.Address, user []common.Address) (*PoolSwapBorrowRateModeIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "SwapBorrowRateMode", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PoolSwapBorrowRateModeIterator{contract: _Pool.contract, event: "SwapBorrowRateMode", logs: logs, sub: sub}, nil
}

// WatchSwapBorrowRateMode is a free log subscription operation binding the contract event 0x7962b394d85a534033ba2efcf43cd36de57b7ebeb3de0ca4428965d9b3ddc481.
//
// Solidity: event SwapBorrowRateMode(address indexed reserve, address indexed user, uint8 interestRateMode)
func (_Pool *PoolFilterer) WatchSwapBorrowRateMode(opts *bind.WatchOpts, sink chan<- *PoolSwapBorrowRateMode, reserve []common.Address, user []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "SwapBorrowRateMode", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolSwapBorrowRateMode)
				if err := _Pool.contract.UnpackLog(event, "SwapBorrowRateMode", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSwapBorrowRateMode is a log parse operation binding the contract event 0x7962b394d85a534033ba2efcf43cd36de57b7ebeb3de0ca4428965d9b3ddc481.
//
// Solidity: event SwapBorrowRateMode(address indexed reserve, address indexed user, uint8 interestRateMode)
func (_Pool *PoolFilterer) ParseSwapBorrowRateMode(log types.Log) (*PoolSwapBorrowRateMode, error) {
	event := new(PoolSwapBorrowRateMode)
	if err := _Pool.contract.UnpackLog(event, "SwapBorrowRateMode", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolUserEModeSetIterator is returned from FilterUserEModeSet and is used to iterate over the raw logs and unpacked data for UserEModeSet events raised by the Pool contract.
type PoolUserEModeSetIterator struct {
	Event *PoolUserEModeSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolUserEModeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolUserEModeSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolUserEModeSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolUserEModeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolUserEModeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolUserEModeSet represents a UserEModeSet event raised by the Pool contract.
type PoolUserEModeSet struct {
	User       common.Address
	CategoryId uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUserEModeSet is a free log retrieval operation binding the contract event 0xd728da875fc88944cbf17638bcbe4af0eedaef63becd1d1c57cc097eb4608d84.
//
// Solidity: event UserEModeSet(address indexed user, uint8 categoryId)
func (_Pool *PoolFilterer) FilterUserEModeSet(opts *bind.FilterOpts, user []common.Address) (*PoolUserEModeSetIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "UserEModeSet", userRule)
	if err != nil {
		return nil, err
	}
	return &PoolUserEModeSetIterator{contract: _Pool.contract, event: "UserEModeSet", logs: logs, sub: sub}, nil
}

// WatchUserEModeSet is a free log subscription operation binding the contract event 0xd728da875fc88944cbf17638bcbe4af0eedaef63becd1d1c57cc097eb4608d84.
//
// Solidity: event UserEModeSet(address indexed user, uint8 categoryId)
func (_Pool *PoolFilterer) WatchUserEModeSet(opts *bind.WatchOpts, sink chan<- *PoolUserEModeSet, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "UserEModeSet", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolUserEModeSet)
				if err := _Pool.contract.UnpackLog(event, "UserEModeSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUserEModeSet is a log parse operation binding the contract event 0xd728da875fc88944cbf17638bcbe4af0eedaef63becd1d1c57cc097eb4608d84.
//
// Solidity: event UserEModeSet(address indexed user, uint8 categoryId)
func (_Pool *PoolFilterer) ParseUserEModeSet(log types.Log) (*PoolUserEModeSet, error) {
	event := new(PoolUserEModeSet)
	if err := _Pool.contract.UnpackLog(event, "UserEModeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Pool contract.
type PoolWithdrawIterator struct {
	Event *PoolWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolWithdraw represents a Withdraw event raised by the Pool contract.
type PoolWithdraw struct {
	Reserve common.Address
	User    common.Address
	To      common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x3115d1449a7b732c986cba18244e897a450f61e1bb8d589cd2e69e6c8924f9f7.
//
// Solidity: event Withdraw(address indexed reserve, address indexed user, address indexed to, uint256 amount)
func (_Pool *PoolFilterer) FilterWithdraw(opts *bind.FilterOpts, reserve []common.Address, user []common.Address, to []common.Address) (*PoolWithdrawIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "Withdraw", reserveRule, userRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PoolWithdrawIterator{contract: _Pool.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x3115d1449a7b732c986cba18244e897a450f61e1bb8d589cd2e69e6c8924f9f7.
//
// Solidity: event Withdraw(address indexed reserve, address indexed user, address indexed to, uint256 amount)
func (_Pool *PoolFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *PoolWithdraw, reserve []common.Address, user []common.Address, to []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "Withdraw", reserveRule, userRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolWithdraw)
				if err := _Pool.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdraw is a log parse operation binding the contract event 0x3115d1449a7b732c986cba18244e897a450f61e1bb8d589cd2e69e6c8924f9f7.
//
// Solidity: event Withdraw(address indexed reserve, address indexed user, address indexed to, uint256 amount)
func (_Pool *PoolFilterer) ParseWithdraw(log types.Log) (*PoolWithdraw, error) {
	event := new(PoolWithdraw)
	if err := _Pool.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
