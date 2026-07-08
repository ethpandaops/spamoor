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

// ConfiguratorInputTypesInitReserveInput is an auto generated low-level Go binding around an user-defined struct.
type ConfiguratorInputTypesInitReserveInput struct {
	ATokenImpl                  common.Address
	StableDebtTokenImpl         common.Address
	VariableDebtTokenImpl       common.Address
	UnderlyingAssetDecimals     uint8
	InterestRateStrategyAddress common.Address
	UnderlyingAsset             common.Address
	Treasury                    common.Address
	IncentivesController        common.Address
	ATokenName                  string
	ATokenSymbol                string
	VariableDebtTokenName       string
	VariableDebtTokenSymbol     string
	StableDebtTokenName         string
	StableDebtTokenSymbol       string
	Params                      []byte
}

// ConfiguratorInputTypesUpdateATokenInput is an auto generated low-level Go binding around an user-defined struct.
type ConfiguratorInputTypesUpdateATokenInput struct {
	Asset                common.Address
	Treasury             common.Address
	IncentivesController common.Address
	Name                 string
	Symbol               string
	Implementation       common.Address
	Params               []byte
}

// ConfiguratorInputTypesUpdateDebtTokenInput is an auto generated low-level Go binding around an user-defined struct.
type ConfiguratorInputTypesUpdateDebtTokenInput struct {
	Asset                common.Address
	IncentivesController common.Address
	Name                 string
	Symbol               string
	Implementation       common.Address
	Params               []byte
}

// PoolConfiguratorMetaData contains all meta data concerning the PoolConfigurator contract.
var PoolConfiguratorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ATokenUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBorrowCap\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBorrowCap\",\"type\":\"uint256\"}],\"name\":\"BorrowCapChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"borrowable\",\"type\":\"bool\"}],\"name\":\"BorrowableInIsolationChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBridgeProtocolFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBridgeProtocolFee\",\"type\":\"uint256\"}],\"name\":\"BridgeProtocolFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ltv\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidationThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidationBonus\",\"type\":\"uint256\"}],\"name\":\"CollateralConfigurationChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldDebtCeiling\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newDebtCeiling\",\"type\":\"uint256\"}],\"name\":\"DebtCeilingChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"oldCategoryId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"newCategoryId\",\"type\":\"uint8\"}],\"name\":\"EModeAssetCategoryChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"categoryId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ltv\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidationThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidationBonus\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"name\":\"EModeCategoryAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"oldFlashloanPremiumToProtocol\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"newFlashloanPremiumToProtocol\",\"type\":\"uint128\"}],\"name\":\"FlashloanPremiumToProtocolUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"oldFlashloanPremiumTotal\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"newFlashloanPremiumTotal\",\"type\":\"uint128\"}],\"name\":\"FlashloanPremiumTotalUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"LiquidationProtocolFeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"name\":\"ReserveActive\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ReserveBorrowing\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"ReserveDropped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldReserveFactor\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newReserveFactor\",\"type\":\"uint256\"}],\"name\":\"ReserveFactorChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ReserveFlashLoaning\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"frozen\",\"type\":\"bool\"}],\"name\":\"ReserveFrozen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stableDebtToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"variableDebtToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"interestRateStrategyAddress\",\"type\":\"address\"}],\"name\":\"ReserveInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldStrategy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newStrategy\",\"type\":\"address\"}],\"name\":\"ReserveInterestRateStrategyChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"ReservePaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ReserveStableRateBorrowing\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"oldState\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newState\",\"type\":\"bool\"}],\"name\":\"SiloedBorrowingChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"StableDebtTokenUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldSupplyCap\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSupplyCap\",\"type\":\"uint256\"}],\"name\":\"SupplyCapChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldUnbackedMintCap\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newUnbackedMintCap\",\"type\":\"uint256\"}],\"name\":\"UnbackedMintCapChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"VariableDebtTokenUpgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CONFIGURATOR_REVISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"ltv\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationBonus\",\"type\":\"uint256\"}],\"name\":\"configureReserveAsCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"dropReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"aTokenImpl\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stableDebtTokenImpl\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"variableDebtTokenImpl\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"underlyingAssetDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategyAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"aTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"aTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"variableDebtTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"variableDebtTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"stableDebtTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"stableDebtTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"internalType\":\"structConfiguratorInputTypes.InitReserveInput[]\",\"name\":\"input\",\"type\":\"tuple[]\"}],\"name\":\"initReserves\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"provider\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"newCategoryId\",\"type\":\"uint8\"}],\"name\":\"setAssetEModeCategory\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newBorrowCap\",\"type\":\"uint256\"}],\"name\":\"setBorrowCap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"borrowable\",\"type\":\"bool\"}],\"name\":\"setBorrowableInIsolation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newDebtCeiling\",\"type\":\"uint256\"}],\"name\":\"setDebtCeiling\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"categoryId\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"ltv\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"liquidationThreshold\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"liquidationBonus\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"name\":\"setEModeCategory\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"setLiquidationProtocolFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"setPoolPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"name\":\"setReserveActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setReserveBorrowing\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newReserveFactor\",\"type\":\"uint256\"}],\"name\":\"setReserveFactor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setReserveFlashLoaning\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"freeze\",\"type\":\"bool\"}],\"name\":\"setReserveFreeze\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newRateStrategyAddress\",\"type\":\"address\"}],\"name\":\"setReserveInterestRateStrategyAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"setReservePause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setReserveStableRateBorrowing\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"newSiloed\",\"type\":\"bool\"}],\"name\":\"setSiloedBorrowing\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newSupplyCap\",\"type\":\"uint256\"}],\"name\":\"setSupplyCap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newUnbackedMintCap\",\"type\":\"uint256\"}],\"name\":\"setUnbackedMintCap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"internalType\":\"structConfiguratorInputTypes.UpdateATokenInput\",\"name\":\"input\",\"type\":\"tuple\"}],\"name\":\"updateAToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newBridgeProtocolFee\",\"type\":\"uint256\"}],\"name\":\"updateBridgeProtocolFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"newFlashloanPremiumToProtocol\",\"type\":\"uint128\"}],\"name\":\"updateFlashloanPremiumToProtocol\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"newFlashloanPremiumTotal\",\"type\":\"uint128\"}],\"name\":\"updateFlashloanPremiumTotal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"internalType\":\"structConfiguratorInputTypes.UpdateDebtTokenInput\",\"name\":\"input\",\"type\":\"tuple\"}],\"name\":\"updateStableDebtToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"internalType\":\"structConfiguratorInputTypes.UpdateDebtTokenInput\",\"name\":\"input\",\"type\":\"tuple\"}],\"name\":\"updateVariableDebtToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000805534801561001457600080fd5b50615ecb80620000256000396000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80637af635a611610104578063aeb4fcc1116100a2578063c4d66de811610071578063c4d66de8146103b8578063d14a0983146103cb578063d4fe3f99146103de578063f213ef0e146103f157600080fd5b8063aeb4fcc11461036c578063b736aaeb1461037f578063bb01c37c14610392578063c19d61e4146103a557600080fd5b80638a751a60116100de5780638a751a601461032057806396e957c414610333578063a7fa83b714610346578063ad4e64321461035957600080fd5b80637af635a6146102e05780637c4e560b146102fa5780638a4936761461030d57600080fd5b806348d9fba91161017157806363c9b8601161014b57806363c9b86014610294578063682cf264146102a75780637626cde3146102ba5780637641f3d9146102cd57600080fd5b806348d9fba91461025b5780634b4e67531461026e578063571f03e51461028157600080fd5b80631df970bd116101ad5780631df970bd1461020f57806326d2cec2146102225780633036b4391461023557806338ae0cc31461024857600080fd5b806302fb45e6146101d4578063145f5892146101e95780631d2118f9146101fc575b600080fd5b6101e76101e2366004614dc3565b610404565b005b6101e76101f7366004614e6d565b6104d5565b6101e761020a366004614e99565b61066f565b6101e761021d366004614ef0565b6107f4565b6101e7610230366004614e6d565b610a90565b6101e7610243366004614f14565b610c8f565b6101e7610256366004614f3b565b610e59565b6101e7610269366004614f3b565b610fe4565b6101e761027c366004614e6d565b611170565b6101e761028f366004614e6d565b61136f565b6101e76102a2366004614f69565b6114ff565b6101e76102b5366004614f3b565b6115d0565b6101e76102c8366004614f86565b6117cf565b6101e76102db366004614fc1565b611874565b6102e8600181565b60405190815260200160405180910390f35b6101e7610308366004614fde565b6119c6565b6101e761031b366004614ef0565b611d50565b6101e761032e366004614f3b565b611fdb565b6101e7610341366004614f3b565b6121de565b6101e7610354366004614f3b565b61235d565b6101e7610367366004614f86565b61250a565b6101e761037a366004614e6d565b61257c565b6101e761038d366004614f3b565b6127a9565b6101e76103a0366004615019565b612936565b6101e76103b3366004615075565b6129a8565b6101e76103c6366004614f69565b61301c565b6101e76103d9366004614e6d565b613235565b6101e76103ec366004615143565b6133c5565b6101e76103ff366004614f3b565b6136a4565b61040c613823565b60355473ffffffffffffffffffffffffffffffffffffffff1660005b828110156104cf5773__$3ddc574512022f331a6a4c7e4bbb5c67b6$__63df59b8b28386868581811061045d5761045d615178565b905060200281019061046f91906151a7565b6040518363ffffffff1660e01b815260040161048c929190615299565b60006040518083038186803b1580156104a457600080fd5b505af41580156104b8573d6000803e3d6000fd5b5050505080806104c790615548565b915050610428565b50505050565b6104dd613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa15801561054e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061057291906156b5565b805190915060b01c640fffffffff1661058b8284613c75565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b1580156105ff57600080fd5b505af1158015610613573d6000803e3d6000fd5b5050604080518481526020810187905273ffffffffffffffffffffffffffffffffffffffff881693507f09808b1fc5abde94edf02fdde393bea0d2e4795999ba31695472848638b5c29f9250015b60405180910390a250505050565b610677613a4e565b6035546040517f35ea6a7500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015260009216906335ea6a75906024016101e060405180830381865afa1580156106e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061070d9190615707565b6101608101516035546040517f1d2118f900000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152868116602483015293945091921690631d2118f990604401600060405180830381600087803b15801561078b57600080fd5b505af115801561079f573d6000803e3d6000fd5b50506040805173ffffffffffffffffffffffffffffffffffffffff85811682528781166020830152881693507fdb8dada53709ce4988154324196790c2e4a60c377e1256790946f83b87db3c33925001610661565b6107fc613d19565b60408051808201909152600281527f313900000000000000000000000000000000000000000000000000000000000060208201526127106fffffffffffffffffffffffffffffffff83161115610888576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b60405180910390fd5b50603554604080517f6a99c036000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff1691636a99c0369160048083019260209291908290030181865afa1580156108f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061091d91906158b3565b603554604080517f074b2e43000000000000000000000000000000000000000000000000000000008152905192935073ffffffffffffffffffffffffffffffffffffffff9091169163bcb6e52291839163074b2e43916004808201926020929091908290030181865afa158015610998573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109bc91906158b3565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b1681526fffffffffffffffffffffffffffffffff91821660048201529085166024820152604401600060405180830381600087803b158015610a2657600080fd5b505af1158015610a3a573d6000803e3d6000fd5b5050604080516fffffffffffffffffffffffffffffffff8086168252861660208201527fe7e0c75e1fc2d0bd83dc85d59f085b3e763107c392fb368e85572b292f1f557693500190505b60405180910390a15050565b610a98613a4e565b60408051808201909152600281527f37300000000000000000000000000000000000000000000000000000000000006020820152612710821115610b09576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b506035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa158015610b7b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b9f91906156b5565b805190915060981c61ffff16610bb58284613eac565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b158015610c2957600080fd5b505af1158015610c3d573d6000803e3d6000fd5b5050604080518481526020810187905273ffffffffffffffffffffffffffffffffffffffff881693507fb5b0a963825337808b6e3154de8e98027595a5cad4219bb3a9bc55b192f4b391925001610661565b610c97613d19565b60408051808201909152600281527f32320000000000000000000000000000000000000000000000000000000000006020820152612710821115610d08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50603554604080517f272d9072000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163272d90729160048083019260209291908290030181865afa158015610d79573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d9d91906158d0565b6035546040517f3036b4390000000000000000000000000000000000000000000000000000000081526004810185905291925073ffffffffffffffffffffffffffffffffffffffff1690633036b43990602401600060405180830381600087803b158015610e0a57600080fd5b505af1158015610e1e573d6000803e3d6000fd5b505060408051848152602081018690527f30b17cb587a89089d003457c432f73e22aeee93de425e92224ba01080260ecd99350019050610a84565b610e61613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa158015610ed2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ef691906156b5565b9050610f028183613f4d565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b158015610f7657600080fd5b505af1158015610f8a573d6000803e3d6000fd5b50506040805173ffffffffffffffffffffffffffffffffffffffff8716815285151560208201527f74adf6aaf58c08bc4f993640385e136522375ea3d1589a10d02adbb906c67d1c935001905060405180910390a1505050565b610fec613f92565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa15801561105d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061108191906156b5565b905061108d81836141b9565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b15801561110157600080fd5b505af1158015611115573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff167fe188d542a5f11925d3a3af33703cdd30a43cb3e8066a3cf68b1b57f61a5a94b583604051611163911515815260200190565b60405180910390a2505050565b611178613a4e565b60408051808201909152600281527f363700000000000000000000000000000000000000000000000000000000000060208201526127108211156111e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b506035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa15801561125b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061127f91906156b5565b805190915060401c61ffff1661129582846141fe565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b15801561130957600080fd5b505af115801561131d573d6000803e3d6000fd5b5050604080518481526020810187905273ffffffffffffffffffffffffffffffffffffffff881693507fb46e2b82b0c2cf3d7d9dece53635e165c53e0eaa7a44f904d61a2b7174826aef925001610661565b611377613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa1580156113e8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061140c91906156b5565b805190915060741c640fffffffff16611425828461429f565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b15801561149957600080fd5b505af11580156114ad573d6000803e3d6000fd5b5050604080518481526020810187905273ffffffffffffffffffffffffffffffffffffffff881693507f0263602682188540a2d633561c0b4453b7d8566285e99f9f6018b8ef2facef49925001610661565b611507613d19565b6035546040517f63c9b86000000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8381166004830152909116906363c9b86090602401600060405180830381600087803b15801561157457600080fd5b505af1158015611588573d6000803e3d6000fd5b505060405173ffffffffffffffffffffffffffffffffffffffff841692507feeec4c06f7adad215cbdb4d2960896c83c26aedce02dde76d36fa28588d62da49150600090a250565b6115d8613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa158015611649573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061166d91906156b5565b9050816116ef57805160408051808201909152600281527f383800000000000000000000000000000000000000000000000000000000000060208201529067080000000000000016156116ed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b505b6116f98183614343565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b15801561176d57600080fd5b505af1158015611781573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff167f2443ba28e8d1d88d531a3d90b981816a4f3b3c7f1fd4085c6029e81d1b7a570d83604051611163911515815260200190565b6117d7613d19565b6035546040517ff5b50e7000000000000000000000000000000000000000000000000000000000815273__$3ddc574512022f331a6a4c7e4bbb5c67b6$__9163f5b50e70916118419173ffffffffffffffffffffffffffffffffffffffff169085906004016158e9565b60006040518083038186803b15801561185957600080fd5b505af415801561186d573d6000803e3d6000fd5b5050505050565b61187c614388565b603554604080517fd1946dbc000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163d1946dbc91600480830192869291908290030181865afa1580156118eb573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261193191908101906159ee565b905060005b81518110156119c157600073ffffffffffffffffffffffffffffffffffffffff1682828151811061196957611969615178565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16146119af576119af8282815181106119a1576119a1615178565b602002602001015184610fe4565b806119b981615548565b915050611936565b505050565b6119ce613a4e565b60408051808201909152600281527f3230000000000000000000000000000000000000000000000000000000000000602082015282841115611a3d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b506035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152600092169063c44b11f790602401602060405180830381865afa158015611aaf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ad391906156b5565b90508215611bcf5760408051808201909152600281527f323000000000000000000000000000000000000000000000000000000000000060208201526127108311611b4b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50612710611b59848461451b565b11156040518060400160405280600281526020017f323000000000000000000000000000000000000000000000000000000000000081525090611bc9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50611c46565b60408051808201909152600281527f323000000000000000000000000000000000000000000000000000000000000060208201528215611c3c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50611c468561455e565b611c50818561470f565b611c5a81846147aa565b611c64818361484b565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b158015611cd857600080fd5b505af1158015611cec573d6000803e3d6000fd5b5050604080518781526020810187905290810185905273ffffffffffffffffffffffffffffffffffffffff881692507f637febbda9275aea2e85c0ff690444c8d87eb2e8339bbede9715abcc89cb0995915060600160405180910390a25050505050565b611d58613d19565b60408051808201909152600281527f313900000000000000000000000000000000000000000000000000000000000060208201526127106fffffffffffffffffffffffffffffffff83161115611ddb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50603554604080517f074b2e43000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163074b2e439160048083019260209291908290030181865afa158015611e4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e7091906158b3565b603554604080517f6a99c036000000000000000000000000000000000000000000000000000000008152905192935073ffffffffffffffffffffffffffffffffffffffff9091169163bcb6e5229185918491636a99c0369160048083019260209291908290030181865afa158015611eec573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f1091906158b3565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b1681526fffffffffffffffffffffffffffffffff928316600482015291166024820152604401600060405180830381600087803b158015611f7957600080fd5b505af1158015611f8d573d6000803e3d6000fd5b5050604080516fffffffffffffffffffffffffffffffff8086168252861660208201527f71aba182c9d0529b516de7a78bed74d49c207ef7e152f52f7ea5d8730138f6439350019050610a84565b611fe3613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa158015612054573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061207891906156b5565b905081156120fe5780516704000000000000001615156040518060400160405280600281526020017f3330000000000000000000000000000000000000000000000000000000000000815250906120fc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b505b61210881836148ec565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b15801561217c57600080fd5b505af1158015612190573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff167f0b64d0941719acd363f1a6be3d8525d8ec9d71738f7445aabcd88d7939b472e783604051611163911515815260200190565b6121e6613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa158015612257573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061227b91906156b5565b90506122878183614931565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b1580156122fb57600080fd5b505af115801561230f573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff167f0c4443d258a350d27dc50c378b2ebf165e6469725f786d21b30cab16823f558783604051611163911515815260200190565b612365613a4e565b80156123745761237482614976565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa1580156123e5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061240991906156b5565b90506000612421825167400000000000000016151590565b905061242d8284614b0c565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b1580156124a157600080fd5b505af11580156124b5573d6000803e3d6000fd5b5050604080518415158152861515602082015273ffffffffffffffffffffffffffffffffffffffff881693507f842a280b07e8e502a9101f32a3b768ebaba3655556dd674f0831900861fc674b925001610661565b612512613d19565b6035546040517fb0f0935500000000000000000000000000000000000000000000000000000000815273__$3ddc574512022f331a6a4c7e4bbb5c67b6$__9163b0f09355916118419173ffffffffffffffffffffffffffffffffffffffff169085906004016158e9565b612584613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa1580156125f5573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061261991906156b5565b805190915060d41c64ffffffffff1680612636576126368461455e565b6126408284614b51565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b1580156126b457600080fd5b505af11580156126c8573d6000803e3d6000fd5b50505050826000141561275b576035546040517fe43e88a100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff86811660048301529091169063e43e88a190602401600060405180830381600087803b15801561274257600080fd5b505af1158015612756573d6000803e3d6000fd5b505050505b604080518281526020810185905273ffffffffffffffffffffffffffffffffffffffff8616917f6824a6c7fbc10d2979b1f1ccf2dd4ed0436541679a661dedb5c10bd4be8306829101610661565b6127b1613d19565b806127bf576127bf8261455e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa158015612830573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061285491906156b5565b90506128608183614bf5565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b1580156128d457600080fd5b505af11580156128e8573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff167fc36c7d11ba01a5869d52aa4a3781939dab851cbc9ee6e7fdcedc7d58898a3f1e83604051611163911515815260200190565b61293e613d19565b6035546040517fb13c96a800000000000000000000000000000000000000000000000000000000815273__$3ddc574512022f331a6a4c7e4bbb5c67b6$__9163b13c96a8916118419173ffffffffffffffffffffffffffffffffffffffff16908590600401615aa0565b6129b0613a4e565b60408051808201909152600281527f3231000000000000000000000000000000000000000000000000000000000000602082015261ffff8716612a20576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5060408051808201909152600281527f3231000000000000000000000000000000000000000000000000000000000000602082015261ffff8616612a91576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b508461ffff168661ffff1611156040518060400160405280600281526020017f323100000000000000000000000000000000000000000000000000000000000081525090612b0c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5060408051808201909152600281527f3231000000000000000000000000000000000000000000000000000000000000602082015261271061ffff861611612b81576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50612710612b9661ffff87811690871661451b565b11156040518060400160405280600281526020017f323100000000000000000000000000000000000000000000000000000000000081525090612c06576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50603554604080517fd1946dbc000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163d1946dbc91600480830192869291908290030181865afa158015612c76573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052612cbc91908101906159ee565b905060005b8151811015612ea457603554825160009173ffffffffffffffffffffffffffffffffffffffff169063c44b11f790859085908110612d0157612d01615178565b60200260200101516040518263ffffffff1660e01b8152600401612d41919073ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b602060405180830381865afa158015612d5e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d8291906156b5565b805190915060a81c60ff168a60ff161415612e9157805161ffff168961ffff16116040518060400160405280600281526020017f323100000000000000000000000000000000000000000000000000000000000081525090612e11576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50805160101c61ffff168861ffff16116040518060400160405280600281526020017f323100000000000000000000000000000000000000000000000000000000000081525090612e8f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b505b5080612e9c81615548565b915050612cc1565b50603560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d579ea7d896040518060a001604052808b61ffff1681526020018a61ffff1681526020018961ffff1681526020018873ffffffffffffffffffffffffffffffffffffffff16815260200187878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509152506040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b168152612f9b929190600401615bdf565b600060405180830381600087803b158015612fb557600080fd5b505af1158015612fc9573d6000803e3d6000fd5b505050508760ff167f0acf8b4a3cace10779798a89a206a0ae73a71b63acdd3be2801d39c2ef7ab3cb88888888888860405161300a96959493929190615c55565b60405180910390a25050505050505050565b6001805460ff168061302d5750303b155b80613039575060005481115b6130c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f436f6e747261637420696e7374616e63652068617320616c726561647920626560448201527f656e20696e697469616c697a6564000000000000000000000000000000000000606482015260840161087f565b60015460ff1615801561310257600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168117905560008290555b603480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8516908117909155604080517f026b1d5f000000000000000000000000000000000000000000000000000000008152905163026b1d5f916004808201926020929091908290030181865afa158015613199573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131bd9190615ca1565b603580547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9290921691909117905580156119c157600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055505050565b61323d613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa1580156132ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132d291906156b5565b805190915060501c640fffffffff166132eb8284614c3a565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b15801561335f57600080fd5b505af1158015613373573d6000803e3d6000fd5b5050604080518481526020810187905273ffffffffffffffffffffffffffffffffffffffff881693507fc51aca575985d521c5072ad11549bad77013bb786d57f30f94b40ed8f8dc9bc4925001610661565b6133cd613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa15801561343e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061346291906156b5565b905060ff8216156135ac576035546040517f6c6f6ae100000000000000000000000000000000000000000000000000000000815260ff8416600482015260009173ffffffffffffffffffffffffffffffffffffffff1690636c6f6ae190602401600060405180830381865afa1580156134df573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526135259190810190615cbe565b825190915060101c61ffff16816020015161ffff16116040518060400160405280600281526020017f3137000000000000000000000000000000000000000000000000000000000000815250906135a9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b50505b805160009060a81c60ff1690506135c68260ff8516614cde565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152845160248301529091169063f51e435b90604401600060405180830381600087803b15801561363a57600080fd5b505af115801561364e573d6000803e3d6000fd5b50506040805160ff80861682528716602082015273ffffffffffffffffffffffffffffffffffffffff881693507f5bb69795b6a2ea222d73a5f8939c23471a1f85a99c7ca43c207f1b71f10c6264925001610661565b6136ac613a4e565b6035546040517fc44b11f700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152600092169063c44b11f790602401602060405180830381865afa15801561371d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061374191906156b5565b905061374d8183614d7e565b6035546040517ff51e435b00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8581166004830152835160248301529091169063f51e435b90604401600060405180830381600087803b1580156137c157600080fd5b505af11580156137d5573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff167fc8ff3cc5b0fddaa3e6ebbbd7438f43393e4ea30e88b80ad016c1bc094655034d83604051611163911515815260200190565b603454604080517f707cd716000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163707cd7169160048083019260209291908290030181865afa158015613893573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138b79190615ca1565b6040517f13ee32e000000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff8216906313ee32e090602401602060405180830381865afa158015613924573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139489190615de9565b806139dc57506040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff821690637be53ca190602401602060405180830381865afa1580156139b8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139dc9190615de9565b6040518060400160405280600181526020017f350000000000000000000000000000000000000000000000000000000000000081525090613a4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5050565b603454604080517f707cd716000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163707cd7169160048083019260209291908290030181865afa158015613abe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ae29190615ca1565b6040517f674b5e4d00000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff82169063674b5e4d90602401602060405180830381865afa158015613b4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b739190615de9565b80613c0757506040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff821690637be53ca190602401602060405180830381865afa158015613be3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c079190615de9565b6040518060400160405280600181526020017f340000000000000000000000000000000000000000000000000000000000000081525090613a4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b60408051808201909152600281527f37320000000000000000000000000000000000000000000000000000000000006020820152640fffffffff821115613ce9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517ffffffffffff000000000ffffffffffffffffffffffffffffffffffffffffffff1660b09190911b179052565b603454604080517f707cd716000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163707cd7169160048083019260209291908290030181865afa158015613d89573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613dad9190615ca1565b6040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff821690637be53ca190602401602060405180830381865afa158015613e1a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e3e9190615de9565b6040518060400160405280600181526020017f310000000000000000000000000000000000000000000000000000000000000081525090613a4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b60408051808201909152600281527f3730000000000000000000000000000000000000000000000000000000000000602082015261ffff821115613f1d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517fffffffffffffffffffffff0000ffffffffffffffffffffffffffffffffffffff1660989190911b179052565b603d81613f5b576000613f5e565b60015b83517fffffffffffffffffffffffffffffffffffffffffffffffffdfffffffffffffff1660ff9190911690911b1790915250565b603454604080517f707cd716000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163707cd7169160048083019260209291908290030181865afa158015614002573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906140269190615ca1565b6040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff821690637be53ca190602401602060405180830381865afa158015614093573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906140b79190615de9565b8061414b57506040517f2500f2b600000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff821690632500f2b690602401602060405180830381865afa158015614127573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061414b9190615de9565b6040518060400160405280600181526020017f330000000000000000000000000000000000000000000000000000000000000081525090613a4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b603c816141c75760006141ca565b60015b83517fffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffff1660ff9190911690911b1790915250565b60408051808201909152600281527f3637000000000000000000000000000000000000000000000000000000000000602082015261ffff82111561426f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517fffffffffffffffffffffffffffffffffffffffffffff0000ffffffffffffffff1660409190911b179052565b60408051808201909152600281527f36390000000000000000000000000000000000000000000000000000000000006020820152640fffffffff821115614313576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517fffffffffffffffffffffffffff000000000fffffffffffffffffffffffffffff1660749190911b179052565b603a81614351576000614354565b60015b83517ffffffffffffffffffffffffffffffffffffffffffffffffffbffffffffffffff1660ff9190911690911b1790915250565b603454604080517f707cd716000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163707cd7169160048083019260209291908290030181865afa1580156143f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061441c9190615ca1565b6040517f2500f2b600000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff821690632500f2b690602401602060405180830381865afa158015614489573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144ad9190615de9565b6040518060400160405280600181526020017f320000000000000000000000000000000000000000000000000000000000000081525090613a4a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b600081157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec778390048411151761455057600080fd5b506127109102611388010490565b600080603460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e860accb6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156145ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906145f29190615ca1565b6040517f35ea6a7500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116600483015291909116906335ea6a759060240161018060405180830381865afa158015614661573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906146859190615e06565b50505050505050505092509250508060001480156146a1575081155b6040518060400160405280600281526020017f3138000000000000000000000000000000000000000000000000000000000000815250906104cf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b60408051808201909152600281527f3633000000000000000000000000000000000000000000000000000000000000602082015261ffff821115614780576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000016179052565b60408051808201909152600281527f3634000000000000000000000000000000000000000000000000000000000000602082015261ffff82111561481b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000ffff1660109190911b179052565b60408051808201909152600281527f3635000000000000000000000000000000000000000000000000000000000000602082015261ffff8211156148bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517fffffffffffffffffffffffffffffffffffffffffffffffffffff0000ffffffff1660209190911b179052565b603b816148fa5760006148fd565b60015b83517ffffffffffffffffffffffffffffffffffffffffffffffffff7ffffffffffffff1660ff9190911690911b1790915250565b60398161493f576000614942565b60015b83517ffffffffffffffffffffffffffffffffffffffffffffffffffdffffffffffffff1660ff9190911690911b1790915250565b603454604080517fe860accb000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163e860accb9160048083019260209291908290030181865afa1580156149e6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614a0a9190615ca1565b6040517f4d44ac4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84811660048301529190911690634d44ac4f90602401602060405180830381865afa158015614a78573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614a9c91906158d0565b60408051808201909152600281527f3930000000000000000000000000000000000000000000000000000000000000602082015290915081156119c1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b603e81614b1a576000614b1d565b60015b83517fffffffffffffffffffffffffffffffffffffffffffffffffbfffffffffffffff1660ff9190911690911b1790915250565b60408051808201909152600281527f3733000000000000000000000000000000000000000000000000000000000000602082015264ffffffffff821115614bc5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517ff0000000000fffffffffffffffffffffffffffffffffffffffffffffffffffff1660d49190911b179052565b603881614c03576000614c06565b60015b83517ffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffff1660ff9190911690911b1790915250565b60408051808201909152600281527f36380000000000000000000000000000000000000000000000000000000000006020820152640fffffffff821115614cae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517ffffffffffffffffffffffffffffffffffff000000000ffffffffffffffffffff1660509190911b179052565b60408051808201909152600281527f3731000000000000000000000000000000000000000000000000000000000000602082015260ff821115614d4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161087f91906158a0565b5081517fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff1660a89190911b179052565b603f81614d8c576000614d8f565b60015b83517fffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffff1660ff9190911690911b1790915250565b60008060208385031215614dd657600080fd5b823567ffffffffffffffff80821115614dee57600080fd5b818501915085601f830112614e0257600080fd5b813581811115614e1157600080fd5b8660208260051b8501011115614e2657600080fd5b60209290920196919550909350505050565b73ffffffffffffffffffffffffffffffffffffffff81168114614e5a57600080fd5b50565b8035614e6881614e38565b919050565b60008060408385031215614e8057600080fd5b8235614e8b81614e38565b946020939093013593505050565b60008060408385031215614eac57600080fd5b8235614eb781614e38565b91506020830135614ec781614e38565b809150509250929050565b6fffffffffffffffffffffffffffffffff81168114614e5a57600080fd5b600060208284031215614f0257600080fd5b8135614f0d81614ed2565b9392505050565b600060208284031215614f2657600080fd5b5035919050565b8015158114614e5a57600080fd5b60008060408385031215614f4e57600080fd5b8235614f5981614e38565b91506020830135614ec781614f2d565b600060208284031215614f7b57600080fd5b8135614f0d81614e38565b600060208284031215614f9857600080fd5b813567ffffffffffffffff811115614faf57600080fd5b820160c08185031215614f0d57600080fd5b600060208284031215614fd357600080fd5b8135614f0d81614f2d565b60008060008060808587031215614ff457600080fd5b8435614fff81614e38565b966020860135965060408601359560600135945092505050565b60006020828403121561502b57600080fd5b813567ffffffffffffffff81111561504257600080fd5b820160e08185031215614f0d57600080fd5b803560ff81168114614e6857600080fd5b61ffff81168114614e5a57600080fd5b600080600080600080600060c0888a03121561509057600080fd5b61509988615054565b965060208801356150a981615065565b955060408801356150b981615065565b945060608801356150c981615065565b935060808801356150d981614e38565b925060a088013567ffffffffffffffff808211156150f657600080fd5b818a0191508a601f83011261510a57600080fd5b81358181111561511957600080fd5b8b602082850101111561512b57600080fd5b60208301945080935050505092959891949750929550565b6000806040838503121561515657600080fd5b823561516181614e38565b915061516f60208401615054565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe218336030181126151db57600080fd5b9190910192915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261521a57600080fd5b830160208101925035905067ffffffffffffffff81111561523a57600080fd5b80360383131561524957600080fd5b9250929050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff83168152604060208201526152e3604082016152c984614e5d565b73ffffffffffffffffffffffffffffffffffffffff169052565b60006152f160208401614e5d565b73ffffffffffffffffffffffffffffffffffffffff16606083015261531860408401614e5d565b73ffffffffffffffffffffffffffffffffffffffff16608083015261533f60608401615054565b60ff1660a083015261535360808401614e5d565b73ffffffffffffffffffffffffffffffffffffffff1660c083015261537a60a08401614e5d565b73ffffffffffffffffffffffffffffffffffffffff1660e08301526153a160c08401614e5d565b6101006153c58185018373ffffffffffffffffffffffffffffffffffffffff169052565b6153d160e08601614e5d565b91506101206153f78186018473ffffffffffffffffffffffffffffffffffffffff169052565b615403828701876151e5565b935091506101e0610140818188015261542161022088018686615250565b945061542f838901896151e5565b945092507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc06101608189880301818a015261546b878787615250565b9650615479838b018b6151e5565b9650945061018092508189880301838a0152615496878787615250565b96506154a4818b018b6151e5565b96509450506101a08189880301818a01526154c0878787615250565b96506154ce838b018b6151e5565b965094506101c092508189880301838a01526154eb878787615250565b96506154f9818b018b6151e5565b9650945050808887030183890152615512868686615250565b9550615520828a018a6151e5565b95509350808887030161020089015250505061553d838383615250565b979650505050505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156155a1577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516101e0810167ffffffffffffffff811182821017156155fb576155fb6155a8565b60405290565b60405160a0810167ffffffffffffffff811182821017156155fb576155fb6155a8565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561566b5761566b6155a8565b604052919050565b60006020828403121561568557600080fd5b6040516020810181811067ffffffffffffffff821117156156a8576156a86155a8565b6040529151825250919050565b6000602082840312156156c757600080fd5b614f0d8383615673565b8051614e6881614ed2565b805164ffffffffff81168114614e6857600080fd5b8051614e6881615065565b8051614e6881614e38565b60006101e0828403121561571a57600080fd5b6157226155d7565b61572c8484615673565b815261573a602084016156d1565b602082015261574b604084016156d1565b604082015261575c606084016156d1565b606082015261576d608084016156d1565b608082015261577e60a084016156d1565b60a082015261578f60c084016156dc565b60c08201526157a060e084016156f1565b60e08201526101006157b38185016156fc565b908201526101206157c58482016156fc565b908201526101406157d78482016156fc565b908201526101606157e98482016156fc565b908201526101806157fb8482016156d1565b908201526101a061580d8482016156d1565b908201526101c061581f8482016156d1565b908201529392505050565b60005b8381101561584557818101518382015260200161582d565b838111156104cf5750506000910152565b6000815180845261586e81602086016020860161582a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000614f0d6020830184615856565b6000602082840312156158c557600080fd5b8151614f0d81614ed2565b6000602082840312156158e257600080fd5b5051919050565b600073ffffffffffffffffffffffffffffffffffffffff808516835260406020840152833561591781614e38565b81166040840152602084013561592c81614e38565b16606083015261593f60408401846151e5565b60c0608085015261595561010085018284615250565b91505061596560608501856151e5565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0808685030160a087015261599b848385615250565b93506159a960808801614e5d565b73ffffffffffffffffffffffffffffffffffffffff811660c088015292506159d460a08801886151e5565b93509150808685030160e08701525061553d838383615250565b60006020808385031215615a0157600080fd5b825167ffffffffffffffff80821115615a1957600080fd5b818501915085601f830112615a2d57600080fd5b815181811115615a3f57615a3f6155a8565b8060051b9150615a50848301615624565b8181529183018401918481019088841115615a6a57600080fd5b938501935b83851015615a945784519250615a8483614e38565b8282529385019390850190615a6f565b98975050505050505050565b600073ffffffffffffffffffffffffffffffffffffffff8085168352604060208401528335615ace81614e38565b166040830152615ae060208401614e5d565b73ffffffffffffffffffffffffffffffffffffffff166060830152615b0760408401614e5d565b73ffffffffffffffffffffffffffffffffffffffff166080830152615b2f60608401846151e5565b60e060a0850152615b4561012085018284615250565b915050615b5560808501856151e5565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0808685030160c0870152615b8b848385615250565b9350615b9960a08801614e5d565b73ffffffffffffffffffffffffffffffffffffffff811660e08801529250615bc460c08801886151e5565b9350915080868503016101008701525061553d838383615250565b60ff8316815260406020820152600061ffff8084511660408401528060208501511660608401528060408501511660808401525073ffffffffffffffffffffffffffffffffffffffff60608401511660a0830152608083015160a060c0840152615c4c60e0840182615856565b95945050505050565b600061ffff8089168352808816602084015280871660408401525073ffffffffffffffffffffffffffffffffffffffff8516606083015260a06080830152615a9460a083018486615250565b600060208284031215615cb357600080fd5b8151614f0d81614e38565b60006020808385031215615cd157600080fd5b825167ffffffffffffffff80821115615ce957600080fd5b9084019060a08287031215615cfd57600080fd5b615d05615601565b8251615d1081615065565b815282840151615d1f81615065565b818501526040830151615d3181615065565b60408201526060830151615d4481614e38565b6060820152608083015182811115615d5b57600080fd5b80840193505086601f840112615d7057600080fd5b825182811115615d8257615d826155a8565b615db2857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601615624565b92508083528785828601011115615dc857600080fd5b615dd78186850187870161582a565b50608081019190915295945050505050565b600060208284031215615dfb57600080fd5b8151614f0d81614f2d565b6000806000806000806000806000806000806101808d8f031215615e2957600080fd5b8c519b5060208d01519a5060408d0151995060608d0151985060808d0151975060a08d0151965060c08d0151955060e08d015194506101008d015193506101208d015192506101408d01519150615e836101608e016156dc565b90509295989b509295989b509295989b56fea264697066735822122001dac2e7ba1c03f36af015a9940a8249669c16fe14a56cc365e85fa9f7f0b8a764736f6c634300080a0033",
}

// PoolConfiguratorABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolConfiguratorMetaData.ABI instead.
var PoolConfiguratorABI = PoolConfiguratorMetaData.ABI

// PoolConfiguratorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PoolConfiguratorMetaData.Bin instead.
var PoolConfiguratorBin = PoolConfiguratorMetaData.Bin

// DeployPoolConfigurator deploys a new Ethereum contract, binding an instance of PoolConfigurator to it.
func DeployPoolConfigurator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PoolConfigurator, error) {
	parsed, err := PoolConfiguratorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PoolConfiguratorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PoolConfigurator{PoolConfiguratorCaller: PoolConfiguratorCaller{contract: contract}, PoolConfiguratorTransactor: PoolConfiguratorTransactor{contract: contract}, PoolConfiguratorFilterer: PoolConfiguratorFilterer{contract: contract}}, nil
}

// PoolConfigurator is an auto generated Go binding around an Ethereum contract.
type PoolConfigurator struct {
	PoolConfiguratorCaller     // Read-only binding to the contract
	PoolConfiguratorTransactor // Write-only binding to the contract
	PoolConfiguratorFilterer   // Log filterer for contract events
}

// PoolConfiguratorCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolConfiguratorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolConfiguratorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolConfiguratorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolConfiguratorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolConfiguratorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolConfiguratorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolConfiguratorSession struct {
	Contract     *PoolConfigurator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolConfiguratorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolConfiguratorCallerSession struct {
	Contract *PoolConfiguratorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// PoolConfiguratorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolConfiguratorTransactorSession struct {
	Contract     *PoolConfiguratorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// PoolConfiguratorRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolConfiguratorRaw struct {
	Contract *PoolConfigurator // Generic contract binding to access the raw methods on
}

// PoolConfiguratorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolConfiguratorCallerRaw struct {
	Contract *PoolConfiguratorCaller // Generic read-only contract binding to access the raw methods on
}

// PoolConfiguratorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolConfiguratorTransactorRaw struct {
	Contract *PoolConfiguratorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolConfigurator creates a new instance of PoolConfigurator, bound to a specific deployed contract.
func NewPoolConfigurator(address common.Address, backend bind.ContractBackend) (*PoolConfigurator, error) {
	contract, err := bindPoolConfigurator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolConfigurator{PoolConfiguratorCaller: PoolConfiguratorCaller{contract: contract}, PoolConfiguratorTransactor: PoolConfiguratorTransactor{contract: contract}, PoolConfiguratorFilterer: PoolConfiguratorFilterer{contract: contract}}, nil
}

// NewPoolConfiguratorCaller creates a new read-only instance of PoolConfigurator, bound to a specific deployed contract.
func NewPoolConfiguratorCaller(address common.Address, caller bind.ContractCaller) (*PoolConfiguratorCaller, error) {
	contract, err := bindPoolConfigurator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorCaller{contract: contract}, nil
}

// NewPoolConfiguratorTransactor creates a new write-only instance of PoolConfigurator, bound to a specific deployed contract.
func NewPoolConfiguratorTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolConfiguratorTransactor, error) {
	contract, err := bindPoolConfigurator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorTransactor{contract: contract}, nil
}

// NewPoolConfiguratorFilterer creates a new log filterer instance of PoolConfigurator, bound to a specific deployed contract.
func NewPoolConfiguratorFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolConfiguratorFilterer, error) {
	contract, err := bindPoolConfigurator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorFilterer{contract: contract}, nil
}

// bindPoolConfigurator binds a generic wrapper to an already deployed contract.
func bindPoolConfigurator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PoolConfiguratorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolConfigurator *PoolConfiguratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolConfigurator.Contract.PoolConfiguratorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolConfigurator *PoolConfiguratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.PoolConfiguratorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolConfigurator *PoolConfiguratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.PoolConfiguratorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolConfigurator *PoolConfiguratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolConfigurator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolConfigurator *PoolConfiguratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolConfigurator *PoolConfiguratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.contract.Transact(opts, method, params...)
}

// CONFIGURATORREVISION is a free data retrieval call binding the contract method 0x7af635a6.
//
// Solidity: function CONFIGURATOR_REVISION() view returns(uint256)
func (_PoolConfigurator *PoolConfiguratorCaller) CONFIGURATORREVISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolConfigurator.contract.Call(opts, &out, "CONFIGURATOR_REVISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CONFIGURATORREVISION is a free data retrieval call binding the contract method 0x7af635a6.
//
// Solidity: function CONFIGURATOR_REVISION() view returns(uint256)
func (_PoolConfigurator *PoolConfiguratorSession) CONFIGURATORREVISION() (*big.Int, error) {
	return _PoolConfigurator.Contract.CONFIGURATORREVISION(&_PoolConfigurator.CallOpts)
}

// CONFIGURATORREVISION is a free data retrieval call binding the contract method 0x7af635a6.
//
// Solidity: function CONFIGURATOR_REVISION() view returns(uint256)
func (_PoolConfigurator *PoolConfiguratorCallerSession) CONFIGURATORREVISION() (*big.Int, error) {
	return _PoolConfigurator.Contract.CONFIGURATORREVISION(&_PoolConfigurator.CallOpts)
}

// ConfigureReserveAsCollateral is a paid mutator transaction binding the contract method 0x7c4e560b.
//
// Solidity: function configureReserveAsCollateral(address asset, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) ConfigureReserveAsCollateral(opts *bind.TransactOpts, asset common.Address, ltv *big.Int, liquidationThreshold *big.Int, liquidationBonus *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "configureReserveAsCollateral", asset, ltv, liquidationThreshold, liquidationBonus)
}

// ConfigureReserveAsCollateral is a paid mutator transaction binding the contract method 0x7c4e560b.
//
// Solidity: function configureReserveAsCollateral(address asset, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus) returns()
func (_PoolConfigurator *PoolConfiguratorSession) ConfigureReserveAsCollateral(asset common.Address, ltv *big.Int, liquidationThreshold *big.Int, liquidationBonus *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.ConfigureReserveAsCollateral(&_PoolConfigurator.TransactOpts, asset, ltv, liquidationThreshold, liquidationBonus)
}

// ConfigureReserveAsCollateral is a paid mutator transaction binding the contract method 0x7c4e560b.
//
// Solidity: function configureReserveAsCollateral(address asset, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) ConfigureReserveAsCollateral(asset common.Address, ltv *big.Int, liquidationThreshold *big.Int, liquidationBonus *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.ConfigureReserveAsCollateral(&_PoolConfigurator.TransactOpts, asset, ltv, liquidationThreshold, liquidationBonus)
}

// DropReserve is a paid mutator transaction binding the contract method 0x63c9b860.
//
// Solidity: function dropReserve(address asset) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) DropReserve(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "dropReserve", asset)
}

// DropReserve is a paid mutator transaction binding the contract method 0x63c9b860.
//
// Solidity: function dropReserve(address asset) returns()
func (_PoolConfigurator *PoolConfiguratorSession) DropReserve(asset common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.DropReserve(&_PoolConfigurator.TransactOpts, asset)
}

// DropReserve is a paid mutator transaction binding the contract method 0x63c9b860.
//
// Solidity: function dropReserve(address asset) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) DropReserve(asset common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.DropReserve(&_PoolConfigurator.TransactOpts, asset)
}

// InitReserves is a paid mutator transaction binding the contract method 0x02fb45e6.
//
// Solidity: function initReserves((address,address,address,uint8,address,address,address,address,string,string,string,string,string,string,bytes)[] input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) InitReserves(opts *bind.TransactOpts, input []ConfiguratorInputTypesInitReserveInput) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "initReserves", input)
}

// InitReserves is a paid mutator transaction binding the contract method 0x02fb45e6.
//
// Solidity: function initReserves((address,address,address,uint8,address,address,address,address,string,string,string,string,string,string,bytes)[] input) returns()
func (_PoolConfigurator *PoolConfiguratorSession) InitReserves(input []ConfiguratorInputTypesInitReserveInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.InitReserves(&_PoolConfigurator.TransactOpts, input)
}

// InitReserves is a paid mutator transaction binding the contract method 0x02fb45e6.
//
// Solidity: function initReserves((address,address,address,uint8,address,address,address,address,string,string,string,string,string,string,bytes)[] input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) InitReserves(input []ConfiguratorInputTypesInitReserveInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.InitReserves(&_PoolConfigurator.TransactOpts, input)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address provider) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) Initialize(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "initialize", provider)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address provider) returns()
func (_PoolConfigurator *PoolConfiguratorSession) Initialize(provider common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.Initialize(&_PoolConfigurator.TransactOpts, provider)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address provider) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) Initialize(provider common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.Initialize(&_PoolConfigurator.TransactOpts, provider)
}

// SetAssetEModeCategory is a paid mutator transaction binding the contract method 0xd4fe3f99.
//
// Solidity: function setAssetEModeCategory(address asset, uint8 newCategoryId) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetAssetEModeCategory(opts *bind.TransactOpts, asset common.Address, newCategoryId uint8) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setAssetEModeCategory", asset, newCategoryId)
}

// SetAssetEModeCategory is a paid mutator transaction binding the contract method 0xd4fe3f99.
//
// Solidity: function setAssetEModeCategory(address asset, uint8 newCategoryId) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetAssetEModeCategory(asset common.Address, newCategoryId uint8) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetAssetEModeCategory(&_PoolConfigurator.TransactOpts, asset, newCategoryId)
}

// SetAssetEModeCategory is a paid mutator transaction binding the contract method 0xd4fe3f99.
//
// Solidity: function setAssetEModeCategory(address asset, uint8 newCategoryId) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetAssetEModeCategory(asset common.Address, newCategoryId uint8) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetAssetEModeCategory(&_PoolConfigurator.TransactOpts, asset, newCategoryId)
}

// SetBorrowCap is a paid mutator transaction binding the contract method 0xd14a0983.
//
// Solidity: function setBorrowCap(address asset, uint256 newBorrowCap) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetBorrowCap(opts *bind.TransactOpts, asset common.Address, newBorrowCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setBorrowCap", asset, newBorrowCap)
}

// SetBorrowCap is a paid mutator transaction binding the contract method 0xd14a0983.
//
// Solidity: function setBorrowCap(address asset, uint256 newBorrowCap) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetBorrowCap(asset common.Address, newBorrowCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetBorrowCap(&_PoolConfigurator.TransactOpts, asset, newBorrowCap)
}

// SetBorrowCap is a paid mutator transaction binding the contract method 0xd14a0983.
//
// Solidity: function setBorrowCap(address asset, uint256 newBorrowCap) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetBorrowCap(asset common.Address, newBorrowCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetBorrowCap(&_PoolConfigurator.TransactOpts, asset, newBorrowCap)
}

// SetBorrowableInIsolation is a paid mutator transaction binding the contract method 0x38ae0cc3.
//
// Solidity: function setBorrowableInIsolation(address asset, bool borrowable) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetBorrowableInIsolation(opts *bind.TransactOpts, asset common.Address, borrowable bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setBorrowableInIsolation", asset, borrowable)
}

// SetBorrowableInIsolation is a paid mutator transaction binding the contract method 0x38ae0cc3.
//
// Solidity: function setBorrowableInIsolation(address asset, bool borrowable) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetBorrowableInIsolation(asset common.Address, borrowable bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetBorrowableInIsolation(&_PoolConfigurator.TransactOpts, asset, borrowable)
}

// SetBorrowableInIsolation is a paid mutator transaction binding the contract method 0x38ae0cc3.
//
// Solidity: function setBorrowableInIsolation(address asset, bool borrowable) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetBorrowableInIsolation(asset common.Address, borrowable bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetBorrowableInIsolation(&_PoolConfigurator.TransactOpts, asset, borrowable)
}

// SetDebtCeiling is a paid mutator transaction binding the contract method 0xaeb4fcc1.
//
// Solidity: function setDebtCeiling(address asset, uint256 newDebtCeiling) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetDebtCeiling(opts *bind.TransactOpts, asset common.Address, newDebtCeiling *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setDebtCeiling", asset, newDebtCeiling)
}

// SetDebtCeiling is a paid mutator transaction binding the contract method 0xaeb4fcc1.
//
// Solidity: function setDebtCeiling(address asset, uint256 newDebtCeiling) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetDebtCeiling(asset common.Address, newDebtCeiling *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetDebtCeiling(&_PoolConfigurator.TransactOpts, asset, newDebtCeiling)
}

// SetDebtCeiling is a paid mutator transaction binding the contract method 0xaeb4fcc1.
//
// Solidity: function setDebtCeiling(address asset, uint256 newDebtCeiling) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetDebtCeiling(asset common.Address, newDebtCeiling *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetDebtCeiling(&_PoolConfigurator.TransactOpts, asset, newDebtCeiling)
}

// SetEModeCategory is a paid mutator transaction binding the contract method 0xc19d61e4.
//
// Solidity: function setEModeCategory(uint8 categoryId, uint16 ltv, uint16 liquidationThreshold, uint16 liquidationBonus, address oracle, string label) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetEModeCategory(opts *bind.TransactOpts, categoryId uint8, ltv uint16, liquidationThreshold uint16, liquidationBonus uint16, oracle common.Address, label string) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setEModeCategory", categoryId, ltv, liquidationThreshold, liquidationBonus, oracle, label)
}

// SetEModeCategory is a paid mutator transaction binding the contract method 0xc19d61e4.
//
// Solidity: function setEModeCategory(uint8 categoryId, uint16 ltv, uint16 liquidationThreshold, uint16 liquidationBonus, address oracle, string label) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetEModeCategory(categoryId uint8, ltv uint16, liquidationThreshold uint16, liquidationBonus uint16, oracle common.Address, label string) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetEModeCategory(&_PoolConfigurator.TransactOpts, categoryId, ltv, liquidationThreshold, liquidationBonus, oracle, label)
}

// SetEModeCategory is a paid mutator transaction binding the contract method 0xc19d61e4.
//
// Solidity: function setEModeCategory(uint8 categoryId, uint16 ltv, uint16 liquidationThreshold, uint16 liquidationBonus, address oracle, string label) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetEModeCategory(categoryId uint8, ltv uint16, liquidationThreshold uint16, liquidationBonus uint16, oracle common.Address, label string) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetEModeCategory(&_PoolConfigurator.TransactOpts, categoryId, ltv, liquidationThreshold, liquidationBonus, oracle, label)
}

// SetLiquidationProtocolFee is a paid mutator transaction binding the contract method 0x26d2cec2.
//
// Solidity: function setLiquidationProtocolFee(address asset, uint256 newFee) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetLiquidationProtocolFee(opts *bind.TransactOpts, asset common.Address, newFee *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setLiquidationProtocolFee", asset, newFee)
}

// SetLiquidationProtocolFee is a paid mutator transaction binding the contract method 0x26d2cec2.
//
// Solidity: function setLiquidationProtocolFee(address asset, uint256 newFee) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetLiquidationProtocolFee(asset common.Address, newFee *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetLiquidationProtocolFee(&_PoolConfigurator.TransactOpts, asset, newFee)
}

// SetLiquidationProtocolFee is a paid mutator transaction binding the contract method 0x26d2cec2.
//
// Solidity: function setLiquidationProtocolFee(address asset, uint256 newFee) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetLiquidationProtocolFee(asset common.Address, newFee *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetLiquidationProtocolFee(&_PoolConfigurator.TransactOpts, asset, newFee)
}

// SetPoolPause is a paid mutator transaction binding the contract method 0x7641f3d9.
//
// Solidity: function setPoolPause(bool paused) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetPoolPause(opts *bind.TransactOpts, paused bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setPoolPause", paused)
}

// SetPoolPause is a paid mutator transaction binding the contract method 0x7641f3d9.
//
// Solidity: function setPoolPause(bool paused) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetPoolPause(paused bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetPoolPause(&_PoolConfigurator.TransactOpts, paused)
}

// SetPoolPause is a paid mutator transaction binding the contract method 0x7641f3d9.
//
// Solidity: function setPoolPause(bool paused) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetPoolPause(paused bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetPoolPause(&_PoolConfigurator.TransactOpts, paused)
}

// SetReserveActive is a paid mutator transaction binding the contract method 0xb736aaeb.
//
// Solidity: function setReserveActive(address asset, bool active) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReserveActive(opts *bind.TransactOpts, asset common.Address, active bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReserveActive", asset, active)
}

// SetReserveActive is a paid mutator transaction binding the contract method 0xb736aaeb.
//
// Solidity: function setReserveActive(address asset, bool active) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReserveActive(asset common.Address, active bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveActive(&_PoolConfigurator.TransactOpts, asset, active)
}

// SetReserveActive is a paid mutator transaction binding the contract method 0xb736aaeb.
//
// Solidity: function setReserveActive(address asset, bool active) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReserveActive(asset common.Address, active bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveActive(&_PoolConfigurator.TransactOpts, asset, active)
}

// SetReserveBorrowing is a paid mutator transaction binding the contract method 0x682cf264.
//
// Solidity: function setReserveBorrowing(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReserveBorrowing(opts *bind.TransactOpts, asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReserveBorrowing", asset, enabled)
}

// SetReserveBorrowing is a paid mutator transaction binding the contract method 0x682cf264.
//
// Solidity: function setReserveBorrowing(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReserveBorrowing(asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveBorrowing(&_PoolConfigurator.TransactOpts, asset, enabled)
}

// SetReserveBorrowing is a paid mutator transaction binding the contract method 0x682cf264.
//
// Solidity: function setReserveBorrowing(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReserveBorrowing(asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveBorrowing(&_PoolConfigurator.TransactOpts, asset, enabled)
}

// SetReserveFactor is a paid mutator transaction binding the contract method 0x4b4e6753.
//
// Solidity: function setReserveFactor(address asset, uint256 newReserveFactor) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReserveFactor(opts *bind.TransactOpts, asset common.Address, newReserveFactor *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReserveFactor", asset, newReserveFactor)
}

// SetReserveFactor is a paid mutator transaction binding the contract method 0x4b4e6753.
//
// Solidity: function setReserveFactor(address asset, uint256 newReserveFactor) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReserveFactor(asset common.Address, newReserveFactor *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveFactor(&_PoolConfigurator.TransactOpts, asset, newReserveFactor)
}

// SetReserveFactor is a paid mutator transaction binding the contract method 0x4b4e6753.
//
// Solidity: function setReserveFactor(address asset, uint256 newReserveFactor) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReserveFactor(asset common.Address, newReserveFactor *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveFactor(&_PoolConfigurator.TransactOpts, asset, newReserveFactor)
}

// SetReserveFlashLoaning is a paid mutator transaction binding the contract method 0xf213ef0e.
//
// Solidity: function setReserveFlashLoaning(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReserveFlashLoaning(opts *bind.TransactOpts, asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReserveFlashLoaning", asset, enabled)
}

// SetReserveFlashLoaning is a paid mutator transaction binding the contract method 0xf213ef0e.
//
// Solidity: function setReserveFlashLoaning(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReserveFlashLoaning(asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveFlashLoaning(&_PoolConfigurator.TransactOpts, asset, enabled)
}

// SetReserveFlashLoaning is a paid mutator transaction binding the contract method 0xf213ef0e.
//
// Solidity: function setReserveFlashLoaning(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReserveFlashLoaning(asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveFlashLoaning(&_PoolConfigurator.TransactOpts, asset, enabled)
}

// SetReserveFreeze is a paid mutator transaction binding the contract method 0x96e957c4.
//
// Solidity: function setReserveFreeze(address asset, bool freeze) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReserveFreeze(opts *bind.TransactOpts, asset common.Address, freeze bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReserveFreeze", asset, freeze)
}

// SetReserveFreeze is a paid mutator transaction binding the contract method 0x96e957c4.
//
// Solidity: function setReserveFreeze(address asset, bool freeze) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReserveFreeze(asset common.Address, freeze bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveFreeze(&_PoolConfigurator.TransactOpts, asset, freeze)
}

// SetReserveFreeze is a paid mutator transaction binding the contract method 0x96e957c4.
//
// Solidity: function setReserveFreeze(address asset, bool freeze) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReserveFreeze(asset common.Address, freeze bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveFreeze(&_PoolConfigurator.TransactOpts, asset, freeze)
}

// SetReserveInterestRateStrategyAddress is a paid mutator transaction binding the contract method 0x1d2118f9.
//
// Solidity: function setReserveInterestRateStrategyAddress(address asset, address newRateStrategyAddress) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReserveInterestRateStrategyAddress(opts *bind.TransactOpts, asset common.Address, newRateStrategyAddress common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReserveInterestRateStrategyAddress", asset, newRateStrategyAddress)
}

// SetReserveInterestRateStrategyAddress is a paid mutator transaction binding the contract method 0x1d2118f9.
//
// Solidity: function setReserveInterestRateStrategyAddress(address asset, address newRateStrategyAddress) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReserveInterestRateStrategyAddress(asset common.Address, newRateStrategyAddress common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveInterestRateStrategyAddress(&_PoolConfigurator.TransactOpts, asset, newRateStrategyAddress)
}

// SetReserveInterestRateStrategyAddress is a paid mutator transaction binding the contract method 0x1d2118f9.
//
// Solidity: function setReserveInterestRateStrategyAddress(address asset, address newRateStrategyAddress) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReserveInterestRateStrategyAddress(asset common.Address, newRateStrategyAddress common.Address) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveInterestRateStrategyAddress(&_PoolConfigurator.TransactOpts, asset, newRateStrategyAddress)
}

// SetReservePause is a paid mutator transaction binding the contract method 0x48d9fba9.
//
// Solidity: function setReservePause(address asset, bool paused) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReservePause(opts *bind.TransactOpts, asset common.Address, paused bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReservePause", asset, paused)
}

// SetReservePause is a paid mutator transaction binding the contract method 0x48d9fba9.
//
// Solidity: function setReservePause(address asset, bool paused) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReservePause(asset common.Address, paused bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReservePause(&_PoolConfigurator.TransactOpts, asset, paused)
}

// SetReservePause is a paid mutator transaction binding the contract method 0x48d9fba9.
//
// Solidity: function setReservePause(address asset, bool paused) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReservePause(asset common.Address, paused bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReservePause(&_PoolConfigurator.TransactOpts, asset, paused)
}

// SetReserveStableRateBorrowing is a paid mutator transaction binding the contract method 0x8a751a60.
//
// Solidity: function setReserveStableRateBorrowing(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetReserveStableRateBorrowing(opts *bind.TransactOpts, asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setReserveStableRateBorrowing", asset, enabled)
}

// SetReserveStableRateBorrowing is a paid mutator transaction binding the contract method 0x8a751a60.
//
// Solidity: function setReserveStableRateBorrowing(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetReserveStableRateBorrowing(asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveStableRateBorrowing(&_PoolConfigurator.TransactOpts, asset, enabled)
}

// SetReserveStableRateBorrowing is a paid mutator transaction binding the contract method 0x8a751a60.
//
// Solidity: function setReserveStableRateBorrowing(address asset, bool enabled) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetReserveStableRateBorrowing(asset common.Address, enabled bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetReserveStableRateBorrowing(&_PoolConfigurator.TransactOpts, asset, enabled)
}

// SetSiloedBorrowing is a paid mutator transaction binding the contract method 0xa7fa83b7.
//
// Solidity: function setSiloedBorrowing(address asset, bool newSiloed) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetSiloedBorrowing(opts *bind.TransactOpts, asset common.Address, newSiloed bool) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setSiloedBorrowing", asset, newSiloed)
}

// SetSiloedBorrowing is a paid mutator transaction binding the contract method 0xa7fa83b7.
//
// Solidity: function setSiloedBorrowing(address asset, bool newSiloed) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetSiloedBorrowing(asset common.Address, newSiloed bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetSiloedBorrowing(&_PoolConfigurator.TransactOpts, asset, newSiloed)
}

// SetSiloedBorrowing is a paid mutator transaction binding the contract method 0xa7fa83b7.
//
// Solidity: function setSiloedBorrowing(address asset, bool newSiloed) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetSiloedBorrowing(asset common.Address, newSiloed bool) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetSiloedBorrowing(&_PoolConfigurator.TransactOpts, asset, newSiloed)
}

// SetSupplyCap is a paid mutator transaction binding the contract method 0x571f03e5.
//
// Solidity: function setSupplyCap(address asset, uint256 newSupplyCap) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetSupplyCap(opts *bind.TransactOpts, asset common.Address, newSupplyCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setSupplyCap", asset, newSupplyCap)
}

// SetSupplyCap is a paid mutator transaction binding the contract method 0x571f03e5.
//
// Solidity: function setSupplyCap(address asset, uint256 newSupplyCap) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetSupplyCap(asset common.Address, newSupplyCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetSupplyCap(&_PoolConfigurator.TransactOpts, asset, newSupplyCap)
}

// SetSupplyCap is a paid mutator transaction binding the contract method 0x571f03e5.
//
// Solidity: function setSupplyCap(address asset, uint256 newSupplyCap) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetSupplyCap(asset common.Address, newSupplyCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetSupplyCap(&_PoolConfigurator.TransactOpts, asset, newSupplyCap)
}

// SetUnbackedMintCap is a paid mutator transaction binding the contract method 0x145f5892.
//
// Solidity: function setUnbackedMintCap(address asset, uint256 newUnbackedMintCap) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) SetUnbackedMintCap(opts *bind.TransactOpts, asset common.Address, newUnbackedMintCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "setUnbackedMintCap", asset, newUnbackedMintCap)
}

// SetUnbackedMintCap is a paid mutator transaction binding the contract method 0x145f5892.
//
// Solidity: function setUnbackedMintCap(address asset, uint256 newUnbackedMintCap) returns()
func (_PoolConfigurator *PoolConfiguratorSession) SetUnbackedMintCap(asset common.Address, newUnbackedMintCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetUnbackedMintCap(&_PoolConfigurator.TransactOpts, asset, newUnbackedMintCap)
}

// SetUnbackedMintCap is a paid mutator transaction binding the contract method 0x145f5892.
//
// Solidity: function setUnbackedMintCap(address asset, uint256 newUnbackedMintCap) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) SetUnbackedMintCap(asset common.Address, newUnbackedMintCap *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.SetUnbackedMintCap(&_PoolConfigurator.TransactOpts, asset, newUnbackedMintCap)
}

// UpdateAToken is a paid mutator transaction binding the contract method 0xbb01c37c.
//
// Solidity: function updateAToken((address,address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) UpdateAToken(opts *bind.TransactOpts, input ConfiguratorInputTypesUpdateATokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "updateAToken", input)
}

// UpdateAToken is a paid mutator transaction binding the contract method 0xbb01c37c.
//
// Solidity: function updateAToken((address,address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorSession) UpdateAToken(input ConfiguratorInputTypesUpdateATokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateAToken(&_PoolConfigurator.TransactOpts, input)
}

// UpdateAToken is a paid mutator transaction binding the contract method 0xbb01c37c.
//
// Solidity: function updateAToken((address,address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) UpdateAToken(input ConfiguratorInputTypesUpdateATokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateAToken(&_PoolConfigurator.TransactOpts, input)
}

// UpdateBridgeProtocolFee is a paid mutator transaction binding the contract method 0x3036b439.
//
// Solidity: function updateBridgeProtocolFee(uint256 newBridgeProtocolFee) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) UpdateBridgeProtocolFee(opts *bind.TransactOpts, newBridgeProtocolFee *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "updateBridgeProtocolFee", newBridgeProtocolFee)
}

// UpdateBridgeProtocolFee is a paid mutator transaction binding the contract method 0x3036b439.
//
// Solidity: function updateBridgeProtocolFee(uint256 newBridgeProtocolFee) returns()
func (_PoolConfigurator *PoolConfiguratorSession) UpdateBridgeProtocolFee(newBridgeProtocolFee *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateBridgeProtocolFee(&_PoolConfigurator.TransactOpts, newBridgeProtocolFee)
}

// UpdateBridgeProtocolFee is a paid mutator transaction binding the contract method 0x3036b439.
//
// Solidity: function updateBridgeProtocolFee(uint256 newBridgeProtocolFee) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) UpdateBridgeProtocolFee(newBridgeProtocolFee *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateBridgeProtocolFee(&_PoolConfigurator.TransactOpts, newBridgeProtocolFee)
}

// UpdateFlashloanPremiumToProtocol is a paid mutator transaction binding the contract method 0x1df970bd.
//
// Solidity: function updateFlashloanPremiumToProtocol(uint128 newFlashloanPremiumToProtocol) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) UpdateFlashloanPremiumToProtocol(opts *bind.TransactOpts, newFlashloanPremiumToProtocol *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "updateFlashloanPremiumToProtocol", newFlashloanPremiumToProtocol)
}

// UpdateFlashloanPremiumToProtocol is a paid mutator transaction binding the contract method 0x1df970bd.
//
// Solidity: function updateFlashloanPremiumToProtocol(uint128 newFlashloanPremiumToProtocol) returns()
func (_PoolConfigurator *PoolConfiguratorSession) UpdateFlashloanPremiumToProtocol(newFlashloanPremiumToProtocol *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateFlashloanPremiumToProtocol(&_PoolConfigurator.TransactOpts, newFlashloanPremiumToProtocol)
}

// UpdateFlashloanPremiumToProtocol is a paid mutator transaction binding the contract method 0x1df970bd.
//
// Solidity: function updateFlashloanPremiumToProtocol(uint128 newFlashloanPremiumToProtocol) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) UpdateFlashloanPremiumToProtocol(newFlashloanPremiumToProtocol *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateFlashloanPremiumToProtocol(&_PoolConfigurator.TransactOpts, newFlashloanPremiumToProtocol)
}

// UpdateFlashloanPremiumTotal is a paid mutator transaction binding the contract method 0x8a493676.
//
// Solidity: function updateFlashloanPremiumTotal(uint128 newFlashloanPremiumTotal) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) UpdateFlashloanPremiumTotal(opts *bind.TransactOpts, newFlashloanPremiumTotal *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "updateFlashloanPremiumTotal", newFlashloanPremiumTotal)
}

// UpdateFlashloanPremiumTotal is a paid mutator transaction binding the contract method 0x8a493676.
//
// Solidity: function updateFlashloanPremiumTotal(uint128 newFlashloanPremiumTotal) returns()
func (_PoolConfigurator *PoolConfiguratorSession) UpdateFlashloanPremiumTotal(newFlashloanPremiumTotal *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateFlashloanPremiumTotal(&_PoolConfigurator.TransactOpts, newFlashloanPremiumTotal)
}

// UpdateFlashloanPremiumTotal is a paid mutator transaction binding the contract method 0x8a493676.
//
// Solidity: function updateFlashloanPremiumTotal(uint128 newFlashloanPremiumTotal) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) UpdateFlashloanPremiumTotal(newFlashloanPremiumTotal *big.Int) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateFlashloanPremiumTotal(&_PoolConfigurator.TransactOpts, newFlashloanPremiumTotal)
}

// UpdateStableDebtToken is a paid mutator transaction binding the contract method 0x7626cde3.
//
// Solidity: function updateStableDebtToken((address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) UpdateStableDebtToken(opts *bind.TransactOpts, input ConfiguratorInputTypesUpdateDebtTokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "updateStableDebtToken", input)
}

// UpdateStableDebtToken is a paid mutator transaction binding the contract method 0x7626cde3.
//
// Solidity: function updateStableDebtToken((address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorSession) UpdateStableDebtToken(input ConfiguratorInputTypesUpdateDebtTokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateStableDebtToken(&_PoolConfigurator.TransactOpts, input)
}

// UpdateStableDebtToken is a paid mutator transaction binding the contract method 0x7626cde3.
//
// Solidity: function updateStableDebtToken((address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) UpdateStableDebtToken(input ConfiguratorInputTypesUpdateDebtTokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateStableDebtToken(&_PoolConfigurator.TransactOpts, input)
}

// UpdateVariableDebtToken is a paid mutator transaction binding the contract method 0xad4e6432.
//
// Solidity: function updateVariableDebtToken((address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactor) UpdateVariableDebtToken(opts *bind.TransactOpts, input ConfiguratorInputTypesUpdateDebtTokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.contract.Transact(opts, "updateVariableDebtToken", input)
}

// UpdateVariableDebtToken is a paid mutator transaction binding the contract method 0xad4e6432.
//
// Solidity: function updateVariableDebtToken((address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorSession) UpdateVariableDebtToken(input ConfiguratorInputTypesUpdateDebtTokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateVariableDebtToken(&_PoolConfigurator.TransactOpts, input)
}

// UpdateVariableDebtToken is a paid mutator transaction binding the contract method 0xad4e6432.
//
// Solidity: function updateVariableDebtToken((address,address,string,string,address,bytes) input) returns()
func (_PoolConfigurator *PoolConfiguratorTransactorSession) UpdateVariableDebtToken(input ConfiguratorInputTypesUpdateDebtTokenInput) (*types.Transaction, error) {
	return _PoolConfigurator.Contract.UpdateVariableDebtToken(&_PoolConfigurator.TransactOpts, input)
}

// PoolConfiguratorATokenUpgradedIterator is returned from FilterATokenUpgraded and is used to iterate over the raw logs and unpacked data for ATokenUpgraded events raised by the PoolConfigurator contract.
type PoolConfiguratorATokenUpgradedIterator struct {
	Event *PoolConfiguratorATokenUpgraded // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorATokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorATokenUpgraded)
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
		it.Event = new(PoolConfiguratorATokenUpgraded)
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
func (it *PoolConfiguratorATokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorATokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorATokenUpgraded represents a ATokenUpgraded event raised by the PoolConfigurator contract.
type PoolConfiguratorATokenUpgraded struct {
	Asset          common.Address
	Proxy          common.Address
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterATokenUpgraded is a free log retrieval operation binding the contract event 0xa76f65411ec66a7fb6bc467432eb14767900449ae4469fa295e4441fe5e1cb73.
//
// Solidity: event ATokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterATokenUpgraded(opts *bind.FilterOpts, asset []common.Address, proxy []common.Address, implementation []common.Address) (*PoolConfiguratorATokenUpgradedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ATokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorATokenUpgradedIterator{contract: _PoolConfigurator.contract, event: "ATokenUpgraded", logs: logs, sub: sub}, nil
}

// WatchATokenUpgraded is a free log subscription operation binding the contract event 0xa76f65411ec66a7fb6bc467432eb14767900449ae4469fa295e4441fe5e1cb73.
//
// Solidity: event ATokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchATokenUpgraded(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorATokenUpgraded, asset []common.Address, proxy []common.Address, implementation []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ATokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorATokenUpgraded)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ATokenUpgraded", log); err != nil {
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

// ParseATokenUpgraded is a log parse operation binding the contract event 0xa76f65411ec66a7fb6bc467432eb14767900449ae4469fa295e4441fe5e1cb73.
//
// Solidity: event ATokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseATokenUpgraded(log types.Log) (*PoolConfiguratorATokenUpgraded, error) {
	event := new(PoolConfiguratorATokenUpgraded)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ATokenUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorBorrowCapChangedIterator is returned from FilterBorrowCapChanged and is used to iterate over the raw logs and unpacked data for BorrowCapChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorBorrowCapChangedIterator struct {
	Event *PoolConfiguratorBorrowCapChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorBorrowCapChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorBorrowCapChanged)
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
		it.Event = new(PoolConfiguratorBorrowCapChanged)
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
func (it *PoolConfiguratorBorrowCapChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorBorrowCapChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorBorrowCapChanged represents a BorrowCapChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorBorrowCapChanged struct {
	Asset        common.Address
	OldBorrowCap *big.Int
	NewBorrowCap *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBorrowCapChanged is a free log retrieval operation binding the contract event 0xc51aca575985d521c5072ad11549bad77013bb786d57f30f94b40ed8f8dc9bc4.
//
// Solidity: event BorrowCapChanged(address indexed asset, uint256 oldBorrowCap, uint256 newBorrowCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterBorrowCapChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorBorrowCapChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "BorrowCapChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorBorrowCapChangedIterator{contract: _PoolConfigurator.contract, event: "BorrowCapChanged", logs: logs, sub: sub}, nil
}

// WatchBorrowCapChanged is a free log subscription operation binding the contract event 0xc51aca575985d521c5072ad11549bad77013bb786d57f30f94b40ed8f8dc9bc4.
//
// Solidity: event BorrowCapChanged(address indexed asset, uint256 oldBorrowCap, uint256 newBorrowCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchBorrowCapChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorBorrowCapChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "BorrowCapChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorBorrowCapChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "BorrowCapChanged", log); err != nil {
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

// ParseBorrowCapChanged is a log parse operation binding the contract event 0xc51aca575985d521c5072ad11549bad77013bb786d57f30f94b40ed8f8dc9bc4.
//
// Solidity: event BorrowCapChanged(address indexed asset, uint256 oldBorrowCap, uint256 newBorrowCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseBorrowCapChanged(log types.Log) (*PoolConfiguratorBorrowCapChanged, error) {
	event := new(PoolConfiguratorBorrowCapChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "BorrowCapChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorBorrowableInIsolationChangedIterator is returned from FilterBorrowableInIsolationChanged and is used to iterate over the raw logs and unpacked data for BorrowableInIsolationChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorBorrowableInIsolationChangedIterator struct {
	Event *PoolConfiguratorBorrowableInIsolationChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorBorrowableInIsolationChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorBorrowableInIsolationChanged)
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
		it.Event = new(PoolConfiguratorBorrowableInIsolationChanged)
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
func (it *PoolConfiguratorBorrowableInIsolationChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorBorrowableInIsolationChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorBorrowableInIsolationChanged represents a BorrowableInIsolationChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorBorrowableInIsolationChanged struct {
	Asset      common.Address
	Borrowable bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBorrowableInIsolationChanged is a free log retrieval operation binding the contract event 0x74adf6aaf58c08bc4f993640385e136522375ea3d1589a10d02adbb906c67d1c.
//
// Solidity: event BorrowableInIsolationChanged(address asset, bool borrowable)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterBorrowableInIsolationChanged(opts *bind.FilterOpts) (*PoolConfiguratorBorrowableInIsolationChangedIterator, error) {

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "BorrowableInIsolationChanged")
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorBorrowableInIsolationChangedIterator{contract: _PoolConfigurator.contract, event: "BorrowableInIsolationChanged", logs: logs, sub: sub}, nil
}

// WatchBorrowableInIsolationChanged is a free log subscription operation binding the contract event 0x74adf6aaf58c08bc4f993640385e136522375ea3d1589a10d02adbb906c67d1c.
//
// Solidity: event BorrowableInIsolationChanged(address asset, bool borrowable)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchBorrowableInIsolationChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorBorrowableInIsolationChanged) (event.Subscription, error) {

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "BorrowableInIsolationChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorBorrowableInIsolationChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "BorrowableInIsolationChanged", log); err != nil {
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

// ParseBorrowableInIsolationChanged is a log parse operation binding the contract event 0x74adf6aaf58c08bc4f993640385e136522375ea3d1589a10d02adbb906c67d1c.
//
// Solidity: event BorrowableInIsolationChanged(address asset, bool borrowable)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseBorrowableInIsolationChanged(log types.Log) (*PoolConfiguratorBorrowableInIsolationChanged, error) {
	event := new(PoolConfiguratorBorrowableInIsolationChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "BorrowableInIsolationChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorBridgeProtocolFeeUpdatedIterator is returned from FilterBridgeProtocolFeeUpdated and is used to iterate over the raw logs and unpacked data for BridgeProtocolFeeUpdated events raised by the PoolConfigurator contract.
type PoolConfiguratorBridgeProtocolFeeUpdatedIterator struct {
	Event *PoolConfiguratorBridgeProtocolFeeUpdated // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorBridgeProtocolFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorBridgeProtocolFeeUpdated)
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
		it.Event = new(PoolConfiguratorBridgeProtocolFeeUpdated)
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
func (it *PoolConfiguratorBridgeProtocolFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorBridgeProtocolFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorBridgeProtocolFeeUpdated represents a BridgeProtocolFeeUpdated event raised by the PoolConfigurator contract.
type PoolConfiguratorBridgeProtocolFeeUpdated struct {
	OldBridgeProtocolFee *big.Int
	NewBridgeProtocolFee *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterBridgeProtocolFeeUpdated is a free log retrieval operation binding the contract event 0x30b17cb587a89089d003457c432f73e22aeee93de425e92224ba01080260ecd9.
//
// Solidity: event BridgeProtocolFeeUpdated(uint256 oldBridgeProtocolFee, uint256 newBridgeProtocolFee)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterBridgeProtocolFeeUpdated(opts *bind.FilterOpts) (*PoolConfiguratorBridgeProtocolFeeUpdatedIterator, error) {

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "BridgeProtocolFeeUpdated")
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorBridgeProtocolFeeUpdatedIterator{contract: _PoolConfigurator.contract, event: "BridgeProtocolFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchBridgeProtocolFeeUpdated is a free log subscription operation binding the contract event 0x30b17cb587a89089d003457c432f73e22aeee93de425e92224ba01080260ecd9.
//
// Solidity: event BridgeProtocolFeeUpdated(uint256 oldBridgeProtocolFee, uint256 newBridgeProtocolFee)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchBridgeProtocolFeeUpdated(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorBridgeProtocolFeeUpdated) (event.Subscription, error) {

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "BridgeProtocolFeeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorBridgeProtocolFeeUpdated)
				if err := _PoolConfigurator.contract.UnpackLog(event, "BridgeProtocolFeeUpdated", log); err != nil {
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

// ParseBridgeProtocolFeeUpdated is a log parse operation binding the contract event 0x30b17cb587a89089d003457c432f73e22aeee93de425e92224ba01080260ecd9.
//
// Solidity: event BridgeProtocolFeeUpdated(uint256 oldBridgeProtocolFee, uint256 newBridgeProtocolFee)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseBridgeProtocolFeeUpdated(log types.Log) (*PoolConfiguratorBridgeProtocolFeeUpdated, error) {
	event := new(PoolConfiguratorBridgeProtocolFeeUpdated)
	if err := _PoolConfigurator.contract.UnpackLog(event, "BridgeProtocolFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorCollateralConfigurationChangedIterator is returned from FilterCollateralConfigurationChanged and is used to iterate over the raw logs and unpacked data for CollateralConfigurationChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorCollateralConfigurationChangedIterator struct {
	Event *PoolConfiguratorCollateralConfigurationChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorCollateralConfigurationChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorCollateralConfigurationChanged)
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
		it.Event = new(PoolConfiguratorCollateralConfigurationChanged)
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
func (it *PoolConfiguratorCollateralConfigurationChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorCollateralConfigurationChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorCollateralConfigurationChanged represents a CollateralConfigurationChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorCollateralConfigurationChanged struct {
	Asset                common.Address
	Ltv                  *big.Int
	LiquidationThreshold *big.Int
	LiquidationBonus     *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterCollateralConfigurationChanged is a free log retrieval operation binding the contract event 0x637febbda9275aea2e85c0ff690444c8d87eb2e8339bbede9715abcc89cb0995.
//
// Solidity: event CollateralConfigurationChanged(address indexed asset, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterCollateralConfigurationChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorCollateralConfigurationChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "CollateralConfigurationChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorCollateralConfigurationChangedIterator{contract: _PoolConfigurator.contract, event: "CollateralConfigurationChanged", logs: logs, sub: sub}, nil
}

// WatchCollateralConfigurationChanged is a free log subscription operation binding the contract event 0x637febbda9275aea2e85c0ff690444c8d87eb2e8339bbede9715abcc89cb0995.
//
// Solidity: event CollateralConfigurationChanged(address indexed asset, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchCollateralConfigurationChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorCollateralConfigurationChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "CollateralConfigurationChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorCollateralConfigurationChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "CollateralConfigurationChanged", log); err != nil {
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

// ParseCollateralConfigurationChanged is a log parse operation binding the contract event 0x637febbda9275aea2e85c0ff690444c8d87eb2e8339bbede9715abcc89cb0995.
//
// Solidity: event CollateralConfigurationChanged(address indexed asset, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseCollateralConfigurationChanged(log types.Log) (*PoolConfiguratorCollateralConfigurationChanged, error) {
	event := new(PoolConfiguratorCollateralConfigurationChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "CollateralConfigurationChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorDebtCeilingChangedIterator is returned from FilterDebtCeilingChanged and is used to iterate over the raw logs and unpacked data for DebtCeilingChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorDebtCeilingChangedIterator struct {
	Event *PoolConfiguratorDebtCeilingChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorDebtCeilingChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorDebtCeilingChanged)
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
		it.Event = new(PoolConfiguratorDebtCeilingChanged)
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
func (it *PoolConfiguratorDebtCeilingChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorDebtCeilingChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorDebtCeilingChanged represents a DebtCeilingChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorDebtCeilingChanged struct {
	Asset          common.Address
	OldDebtCeiling *big.Int
	NewDebtCeiling *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDebtCeilingChanged is a free log retrieval operation binding the contract event 0x6824a6c7fbc10d2979b1f1ccf2dd4ed0436541679a661dedb5c10bd4be830682.
//
// Solidity: event DebtCeilingChanged(address indexed asset, uint256 oldDebtCeiling, uint256 newDebtCeiling)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterDebtCeilingChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorDebtCeilingChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "DebtCeilingChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorDebtCeilingChangedIterator{contract: _PoolConfigurator.contract, event: "DebtCeilingChanged", logs: logs, sub: sub}, nil
}

// WatchDebtCeilingChanged is a free log subscription operation binding the contract event 0x6824a6c7fbc10d2979b1f1ccf2dd4ed0436541679a661dedb5c10bd4be830682.
//
// Solidity: event DebtCeilingChanged(address indexed asset, uint256 oldDebtCeiling, uint256 newDebtCeiling)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchDebtCeilingChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorDebtCeilingChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "DebtCeilingChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorDebtCeilingChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "DebtCeilingChanged", log); err != nil {
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

// ParseDebtCeilingChanged is a log parse operation binding the contract event 0x6824a6c7fbc10d2979b1f1ccf2dd4ed0436541679a661dedb5c10bd4be830682.
//
// Solidity: event DebtCeilingChanged(address indexed asset, uint256 oldDebtCeiling, uint256 newDebtCeiling)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseDebtCeilingChanged(log types.Log) (*PoolConfiguratorDebtCeilingChanged, error) {
	event := new(PoolConfiguratorDebtCeilingChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "DebtCeilingChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorEModeAssetCategoryChangedIterator is returned from FilterEModeAssetCategoryChanged and is used to iterate over the raw logs and unpacked data for EModeAssetCategoryChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorEModeAssetCategoryChangedIterator struct {
	Event *PoolConfiguratorEModeAssetCategoryChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorEModeAssetCategoryChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorEModeAssetCategoryChanged)
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
		it.Event = new(PoolConfiguratorEModeAssetCategoryChanged)
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
func (it *PoolConfiguratorEModeAssetCategoryChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorEModeAssetCategoryChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorEModeAssetCategoryChanged represents a EModeAssetCategoryChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorEModeAssetCategoryChanged struct {
	Asset         common.Address
	OldCategoryId uint8
	NewCategoryId uint8
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterEModeAssetCategoryChanged is a free log retrieval operation binding the contract event 0x5bb69795b6a2ea222d73a5f8939c23471a1f85a99c7ca43c207f1b71f10c6264.
//
// Solidity: event EModeAssetCategoryChanged(address indexed asset, uint8 oldCategoryId, uint8 newCategoryId)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterEModeAssetCategoryChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorEModeAssetCategoryChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "EModeAssetCategoryChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorEModeAssetCategoryChangedIterator{contract: _PoolConfigurator.contract, event: "EModeAssetCategoryChanged", logs: logs, sub: sub}, nil
}

// WatchEModeAssetCategoryChanged is a free log subscription operation binding the contract event 0x5bb69795b6a2ea222d73a5f8939c23471a1f85a99c7ca43c207f1b71f10c6264.
//
// Solidity: event EModeAssetCategoryChanged(address indexed asset, uint8 oldCategoryId, uint8 newCategoryId)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchEModeAssetCategoryChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorEModeAssetCategoryChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "EModeAssetCategoryChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorEModeAssetCategoryChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "EModeAssetCategoryChanged", log); err != nil {
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

// ParseEModeAssetCategoryChanged is a log parse operation binding the contract event 0x5bb69795b6a2ea222d73a5f8939c23471a1f85a99c7ca43c207f1b71f10c6264.
//
// Solidity: event EModeAssetCategoryChanged(address indexed asset, uint8 oldCategoryId, uint8 newCategoryId)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseEModeAssetCategoryChanged(log types.Log) (*PoolConfiguratorEModeAssetCategoryChanged, error) {
	event := new(PoolConfiguratorEModeAssetCategoryChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "EModeAssetCategoryChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorEModeCategoryAddedIterator is returned from FilterEModeCategoryAdded and is used to iterate over the raw logs and unpacked data for EModeCategoryAdded events raised by the PoolConfigurator contract.
type PoolConfiguratorEModeCategoryAddedIterator struct {
	Event *PoolConfiguratorEModeCategoryAdded // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorEModeCategoryAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorEModeCategoryAdded)
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
		it.Event = new(PoolConfiguratorEModeCategoryAdded)
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
func (it *PoolConfiguratorEModeCategoryAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorEModeCategoryAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorEModeCategoryAdded represents a EModeCategoryAdded event raised by the PoolConfigurator contract.
type PoolConfiguratorEModeCategoryAdded struct {
	CategoryId           uint8
	Ltv                  *big.Int
	LiquidationThreshold *big.Int
	LiquidationBonus     *big.Int
	Oracle               common.Address
	Label                string
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterEModeCategoryAdded is a free log retrieval operation binding the contract event 0x0acf8b4a3cace10779798a89a206a0ae73a71b63acdd3be2801d39c2ef7ab3cb.
//
// Solidity: event EModeCategoryAdded(uint8 indexed categoryId, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus, address oracle, string label)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterEModeCategoryAdded(opts *bind.FilterOpts, categoryId []uint8) (*PoolConfiguratorEModeCategoryAddedIterator, error) {

	var categoryIdRule []interface{}
	for _, categoryIdItem := range categoryId {
		categoryIdRule = append(categoryIdRule, categoryIdItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "EModeCategoryAdded", categoryIdRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorEModeCategoryAddedIterator{contract: _PoolConfigurator.contract, event: "EModeCategoryAdded", logs: logs, sub: sub}, nil
}

// WatchEModeCategoryAdded is a free log subscription operation binding the contract event 0x0acf8b4a3cace10779798a89a206a0ae73a71b63acdd3be2801d39c2ef7ab3cb.
//
// Solidity: event EModeCategoryAdded(uint8 indexed categoryId, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus, address oracle, string label)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchEModeCategoryAdded(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorEModeCategoryAdded, categoryId []uint8) (event.Subscription, error) {

	var categoryIdRule []interface{}
	for _, categoryIdItem := range categoryId {
		categoryIdRule = append(categoryIdRule, categoryIdItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "EModeCategoryAdded", categoryIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorEModeCategoryAdded)
				if err := _PoolConfigurator.contract.UnpackLog(event, "EModeCategoryAdded", log); err != nil {
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

// ParseEModeCategoryAdded is a log parse operation binding the contract event 0x0acf8b4a3cace10779798a89a206a0ae73a71b63acdd3be2801d39c2ef7ab3cb.
//
// Solidity: event EModeCategoryAdded(uint8 indexed categoryId, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus, address oracle, string label)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseEModeCategoryAdded(log types.Log) (*PoolConfiguratorEModeCategoryAdded, error) {
	event := new(PoolConfiguratorEModeCategoryAdded)
	if err := _PoolConfigurator.contract.UnpackLog(event, "EModeCategoryAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorFlashloanPremiumToProtocolUpdatedIterator is returned from FilterFlashloanPremiumToProtocolUpdated and is used to iterate over the raw logs and unpacked data for FlashloanPremiumToProtocolUpdated events raised by the PoolConfigurator contract.
type PoolConfiguratorFlashloanPremiumToProtocolUpdatedIterator struct {
	Event *PoolConfiguratorFlashloanPremiumToProtocolUpdated // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorFlashloanPremiumToProtocolUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorFlashloanPremiumToProtocolUpdated)
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
		it.Event = new(PoolConfiguratorFlashloanPremiumToProtocolUpdated)
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
func (it *PoolConfiguratorFlashloanPremiumToProtocolUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorFlashloanPremiumToProtocolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorFlashloanPremiumToProtocolUpdated represents a FlashloanPremiumToProtocolUpdated event raised by the PoolConfigurator contract.
type PoolConfiguratorFlashloanPremiumToProtocolUpdated struct {
	OldFlashloanPremiumToProtocol *big.Int
	NewFlashloanPremiumToProtocol *big.Int
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterFlashloanPremiumToProtocolUpdated is a free log retrieval operation binding the contract event 0xe7e0c75e1fc2d0bd83dc85d59f085b3e763107c392fb368e85572b292f1f5576.
//
// Solidity: event FlashloanPremiumToProtocolUpdated(uint128 oldFlashloanPremiumToProtocol, uint128 newFlashloanPremiumToProtocol)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterFlashloanPremiumToProtocolUpdated(opts *bind.FilterOpts) (*PoolConfiguratorFlashloanPremiumToProtocolUpdatedIterator, error) {

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "FlashloanPremiumToProtocolUpdated")
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorFlashloanPremiumToProtocolUpdatedIterator{contract: _PoolConfigurator.contract, event: "FlashloanPremiumToProtocolUpdated", logs: logs, sub: sub}, nil
}

// WatchFlashloanPremiumToProtocolUpdated is a free log subscription operation binding the contract event 0xe7e0c75e1fc2d0bd83dc85d59f085b3e763107c392fb368e85572b292f1f5576.
//
// Solidity: event FlashloanPremiumToProtocolUpdated(uint128 oldFlashloanPremiumToProtocol, uint128 newFlashloanPremiumToProtocol)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchFlashloanPremiumToProtocolUpdated(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorFlashloanPremiumToProtocolUpdated) (event.Subscription, error) {

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "FlashloanPremiumToProtocolUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorFlashloanPremiumToProtocolUpdated)
				if err := _PoolConfigurator.contract.UnpackLog(event, "FlashloanPremiumToProtocolUpdated", log); err != nil {
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

// ParseFlashloanPremiumToProtocolUpdated is a log parse operation binding the contract event 0xe7e0c75e1fc2d0bd83dc85d59f085b3e763107c392fb368e85572b292f1f5576.
//
// Solidity: event FlashloanPremiumToProtocolUpdated(uint128 oldFlashloanPremiumToProtocol, uint128 newFlashloanPremiumToProtocol)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseFlashloanPremiumToProtocolUpdated(log types.Log) (*PoolConfiguratorFlashloanPremiumToProtocolUpdated, error) {
	event := new(PoolConfiguratorFlashloanPremiumToProtocolUpdated)
	if err := _PoolConfigurator.contract.UnpackLog(event, "FlashloanPremiumToProtocolUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorFlashloanPremiumTotalUpdatedIterator is returned from FilterFlashloanPremiumTotalUpdated and is used to iterate over the raw logs and unpacked data for FlashloanPremiumTotalUpdated events raised by the PoolConfigurator contract.
type PoolConfiguratorFlashloanPremiumTotalUpdatedIterator struct {
	Event *PoolConfiguratorFlashloanPremiumTotalUpdated // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorFlashloanPremiumTotalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorFlashloanPremiumTotalUpdated)
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
		it.Event = new(PoolConfiguratorFlashloanPremiumTotalUpdated)
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
func (it *PoolConfiguratorFlashloanPremiumTotalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorFlashloanPremiumTotalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorFlashloanPremiumTotalUpdated represents a FlashloanPremiumTotalUpdated event raised by the PoolConfigurator contract.
type PoolConfiguratorFlashloanPremiumTotalUpdated struct {
	OldFlashloanPremiumTotal *big.Int
	NewFlashloanPremiumTotal *big.Int
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterFlashloanPremiumTotalUpdated is a free log retrieval operation binding the contract event 0x71aba182c9d0529b516de7a78bed74d49c207ef7e152f52f7ea5d8730138f643.
//
// Solidity: event FlashloanPremiumTotalUpdated(uint128 oldFlashloanPremiumTotal, uint128 newFlashloanPremiumTotal)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterFlashloanPremiumTotalUpdated(opts *bind.FilterOpts) (*PoolConfiguratorFlashloanPremiumTotalUpdatedIterator, error) {

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "FlashloanPremiumTotalUpdated")
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorFlashloanPremiumTotalUpdatedIterator{contract: _PoolConfigurator.contract, event: "FlashloanPremiumTotalUpdated", logs: logs, sub: sub}, nil
}

// WatchFlashloanPremiumTotalUpdated is a free log subscription operation binding the contract event 0x71aba182c9d0529b516de7a78bed74d49c207ef7e152f52f7ea5d8730138f643.
//
// Solidity: event FlashloanPremiumTotalUpdated(uint128 oldFlashloanPremiumTotal, uint128 newFlashloanPremiumTotal)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchFlashloanPremiumTotalUpdated(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorFlashloanPremiumTotalUpdated) (event.Subscription, error) {

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "FlashloanPremiumTotalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorFlashloanPremiumTotalUpdated)
				if err := _PoolConfigurator.contract.UnpackLog(event, "FlashloanPremiumTotalUpdated", log); err != nil {
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

// ParseFlashloanPremiumTotalUpdated is a log parse operation binding the contract event 0x71aba182c9d0529b516de7a78bed74d49c207ef7e152f52f7ea5d8730138f643.
//
// Solidity: event FlashloanPremiumTotalUpdated(uint128 oldFlashloanPremiumTotal, uint128 newFlashloanPremiumTotal)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseFlashloanPremiumTotalUpdated(log types.Log) (*PoolConfiguratorFlashloanPremiumTotalUpdated, error) {
	event := new(PoolConfiguratorFlashloanPremiumTotalUpdated)
	if err := _PoolConfigurator.contract.UnpackLog(event, "FlashloanPremiumTotalUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorLiquidationProtocolFeeChangedIterator is returned from FilterLiquidationProtocolFeeChanged and is used to iterate over the raw logs and unpacked data for LiquidationProtocolFeeChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorLiquidationProtocolFeeChangedIterator struct {
	Event *PoolConfiguratorLiquidationProtocolFeeChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorLiquidationProtocolFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorLiquidationProtocolFeeChanged)
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
		it.Event = new(PoolConfiguratorLiquidationProtocolFeeChanged)
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
func (it *PoolConfiguratorLiquidationProtocolFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorLiquidationProtocolFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorLiquidationProtocolFeeChanged represents a LiquidationProtocolFeeChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorLiquidationProtocolFeeChanged struct {
	Asset  common.Address
	OldFee *big.Int
	NewFee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLiquidationProtocolFeeChanged is a free log retrieval operation binding the contract event 0xb5b0a963825337808b6e3154de8e98027595a5cad4219bb3a9bc55b192f4b391.
//
// Solidity: event LiquidationProtocolFeeChanged(address indexed asset, uint256 oldFee, uint256 newFee)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterLiquidationProtocolFeeChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorLiquidationProtocolFeeChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "LiquidationProtocolFeeChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorLiquidationProtocolFeeChangedIterator{contract: _PoolConfigurator.contract, event: "LiquidationProtocolFeeChanged", logs: logs, sub: sub}, nil
}

// WatchLiquidationProtocolFeeChanged is a free log subscription operation binding the contract event 0xb5b0a963825337808b6e3154de8e98027595a5cad4219bb3a9bc55b192f4b391.
//
// Solidity: event LiquidationProtocolFeeChanged(address indexed asset, uint256 oldFee, uint256 newFee)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchLiquidationProtocolFeeChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorLiquidationProtocolFeeChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "LiquidationProtocolFeeChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorLiquidationProtocolFeeChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "LiquidationProtocolFeeChanged", log); err != nil {
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

// ParseLiquidationProtocolFeeChanged is a log parse operation binding the contract event 0xb5b0a963825337808b6e3154de8e98027595a5cad4219bb3a9bc55b192f4b391.
//
// Solidity: event LiquidationProtocolFeeChanged(address indexed asset, uint256 oldFee, uint256 newFee)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseLiquidationProtocolFeeChanged(log types.Log) (*PoolConfiguratorLiquidationProtocolFeeChanged, error) {
	event := new(PoolConfiguratorLiquidationProtocolFeeChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "LiquidationProtocolFeeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveActiveIterator is returned from FilterReserveActive and is used to iterate over the raw logs and unpacked data for ReserveActive events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveActiveIterator struct {
	Event *PoolConfiguratorReserveActive // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveActiveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveActive)
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
		it.Event = new(PoolConfiguratorReserveActive)
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
func (it *PoolConfiguratorReserveActiveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveActiveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveActive represents a ReserveActive event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveActive struct {
	Asset  common.Address
	Active bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReserveActive is a free log retrieval operation binding the contract event 0xc36c7d11ba01a5869d52aa4a3781939dab851cbc9ee6e7fdcedc7d58898a3f1e.
//
// Solidity: event ReserveActive(address indexed asset, bool active)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveActive(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveActiveIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveActive", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveActiveIterator{contract: _PoolConfigurator.contract, event: "ReserveActive", logs: logs, sub: sub}, nil
}

// WatchReserveActive is a free log subscription operation binding the contract event 0xc36c7d11ba01a5869d52aa4a3781939dab851cbc9ee6e7fdcedc7d58898a3f1e.
//
// Solidity: event ReserveActive(address indexed asset, bool active)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveActive(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveActive, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveActive", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveActive)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveActive", log); err != nil {
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

// ParseReserveActive is a log parse operation binding the contract event 0xc36c7d11ba01a5869d52aa4a3781939dab851cbc9ee6e7fdcedc7d58898a3f1e.
//
// Solidity: event ReserveActive(address indexed asset, bool active)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveActive(log types.Log) (*PoolConfiguratorReserveActive, error) {
	event := new(PoolConfiguratorReserveActive)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveActive", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveBorrowingIterator is returned from FilterReserveBorrowing and is used to iterate over the raw logs and unpacked data for ReserveBorrowing events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveBorrowingIterator struct {
	Event *PoolConfiguratorReserveBorrowing // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveBorrowingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveBorrowing)
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
		it.Event = new(PoolConfiguratorReserveBorrowing)
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
func (it *PoolConfiguratorReserveBorrowingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveBorrowingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveBorrowing represents a ReserveBorrowing event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveBorrowing struct {
	Asset   common.Address
	Enabled bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveBorrowing is a free log retrieval operation binding the contract event 0x2443ba28e8d1d88d531a3d90b981816a4f3b3c7f1fd4085c6029e81d1b7a570d.
//
// Solidity: event ReserveBorrowing(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveBorrowing(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveBorrowingIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveBorrowing", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveBorrowingIterator{contract: _PoolConfigurator.contract, event: "ReserveBorrowing", logs: logs, sub: sub}, nil
}

// WatchReserveBorrowing is a free log subscription operation binding the contract event 0x2443ba28e8d1d88d531a3d90b981816a4f3b3c7f1fd4085c6029e81d1b7a570d.
//
// Solidity: event ReserveBorrowing(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveBorrowing(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveBorrowing, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveBorrowing", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveBorrowing)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveBorrowing", log); err != nil {
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

// ParseReserveBorrowing is a log parse operation binding the contract event 0x2443ba28e8d1d88d531a3d90b981816a4f3b3c7f1fd4085c6029e81d1b7a570d.
//
// Solidity: event ReserveBorrowing(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveBorrowing(log types.Log) (*PoolConfiguratorReserveBorrowing, error) {
	event := new(PoolConfiguratorReserveBorrowing)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveBorrowing", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveDroppedIterator is returned from FilterReserveDropped and is used to iterate over the raw logs and unpacked data for ReserveDropped events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveDroppedIterator struct {
	Event *PoolConfiguratorReserveDropped // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveDroppedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveDropped)
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
		it.Event = new(PoolConfiguratorReserveDropped)
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
func (it *PoolConfiguratorReserveDroppedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveDroppedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveDropped represents a ReserveDropped event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveDropped struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterReserveDropped is a free log retrieval operation binding the contract event 0xeeec4c06f7adad215cbdb4d2960896c83c26aedce02dde76d36fa28588d62da4.
//
// Solidity: event ReserveDropped(address indexed asset)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveDropped(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveDroppedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveDropped", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveDroppedIterator{contract: _PoolConfigurator.contract, event: "ReserveDropped", logs: logs, sub: sub}, nil
}

// WatchReserveDropped is a free log subscription operation binding the contract event 0xeeec4c06f7adad215cbdb4d2960896c83c26aedce02dde76d36fa28588d62da4.
//
// Solidity: event ReserveDropped(address indexed asset)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveDropped(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveDropped, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveDropped", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveDropped)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveDropped", log); err != nil {
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

// ParseReserveDropped is a log parse operation binding the contract event 0xeeec4c06f7adad215cbdb4d2960896c83c26aedce02dde76d36fa28588d62da4.
//
// Solidity: event ReserveDropped(address indexed asset)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveDropped(log types.Log) (*PoolConfiguratorReserveDropped, error) {
	event := new(PoolConfiguratorReserveDropped)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveDropped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveFactorChangedIterator is returned from FilterReserveFactorChanged and is used to iterate over the raw logs and unpacked data for ReserveFactorChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveFactorChangedIterator struct {
	Event *PoolConfiguratorReserveFactorChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveFactorChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveFactorChanged)
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
		it.Event = new(PoolConfiguratorReserveFactorChanged)
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
func (it *PoolConfiguratorReserveFactorChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveFactorChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveFactorChanged represents a ReserveFactorChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveFactorChanged struct {
	Asset            common.Address
	OldReserveFactor *big.Int
	NewReserveFactor *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterReserveFactorChanged is a free log retrieval operation binding the contract event 0xb46e2b82b0c2cf3d7d9dece53635e165c53e0eaa7a44f904d61a2b7174826aef.
//
// Solidity: event ReserveFactorChanged(address indexed asset, uint256 oldReserveFactor, uint256 newReserveFactor)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveFactorChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveFactorChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveFactorChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveFactorChangedIterator{contract: _PoolConfigurator.contract, event: "ReserveFactorChanged", logs: logs, sub: sub}, nil
}

// WatchReserveFactorChanged is a free log subscription operation binding the contract event 0xb46e2b82b0c2cf3d7d9dece53635e165c53e0eaa7a44f904d61a2b7174826aef.
//
// Solidity: event ReserveFactorChanged(address indexed asset, uint256 oldReserveFactor, uint256 newReserveFactor)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveFactorChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveFactorChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveFactorChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveFactorChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveFactorChanged", log); err != nil {
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

// ParseReserveFactorChanged is a log parse operation binding the contract event 0xb46e2b82b0c2cf3d7d9dece53635e165c53e0eaa7a44f904d61a2b7174826aef.
//
// Solidity: event ReserveFactorChanged(address indexed asset, uint256 oldReserveFactor, uint256 newReserveFactor)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveFactorChanged(log types.Log) (*PoolConfiguratorReserveFactorChanged, error) {
	event := new(PoolConfiguratorReserveFactorChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveFactorChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveFlashLoaningIterator is returned from FilterReserveFlashLoaning and is used to iterate over the raw logs and unpacked data for ReserveFlashLoaning events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveFlashLoaningIterator struct {
	Event *PoolConfiguratorReserveFlashLoaning // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveFlashLoaningIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveFlashLoaning)
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
		it.Event = new(PoolConfiguratorReserveFlashLoaning)
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
func (it *PoolConfiguratorReserveFlashLoaningIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveFlashLoaningIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveFlashLoaning represents a ReserveFlashLoaning event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveFlashLoaning struct {
	Asset   common.Address
	Enabled bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveFlashLoaning is a free log retrieval operation binding the contract event 0xc8ff3cc5b0fddaa3e6ebbbd7438f43393e4ea30e88b80ad016c1bc094655034d.
//
// Solidity: event ReserveFlashLoaning(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveFlashLoaning(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveFlashLoaningIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveFlashLoaning", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveFlashLoaningIterator{contract: _PoolConfigurator.contract, event: "ReserveFlashLoaning", logs: logs, sub: sub}, nil
}

// WatchReserveFlashLoaning is a free log subscription operation binding the contract event 0xc8ff3cc5b0fddaa3e6ebbbd7438f43393e4ea30e88b80ad016c1bc094655034d.
//
// Solidity: event ReserveFlashLoaning(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveFlashLoaning(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveFlashLoaning, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveFlashLoaning", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveFlashLoaning)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveFlashLoaning", log); err != nil {
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

// ParseReserveFlashLoaning is a log parse operation binding the contract event 0xc8ff3cc5b0fddaa3e6ebbbd7438f43393e4ea30e88b80ad016c1bc094655034d.
//
// Solidity: event ReserveFlashLoaning(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveFlashLoaning(log types.Log) (*PoolConfiguratorReserveFlashLoaning, error) {
	event := new(PoolConfiguratorReserveFlashLoaning)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveFlashLoaning", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveFrozenIterator is returned from FilterReserveFrozen and is used to iterate over the raw logs and unpacked data for ReserveFrozen events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveFrozenIterator struct {
	Event *PoolConfiguratorReserveFrozen // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveFrozen)
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
		it.Event = new(PoolConfiguratorReserveFrozen)
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
func (it *PoolConfiguratorReserveFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveFrozen represents a ReserveFrozen event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveFrozen struct {
	Asset  common.Address
	Frozen bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReserveFrozen is a free log retrieval operation binding the contract event 0x0c4443d258a350d27dc50c378b2ebf165e6469725f786d21b30cab16823f5587.
//
// Solidity: event ReserveFrozen(address indexed asset, bool frozen)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveFrozen(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveFrozenIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveFrozen", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveFrozenIterator{contract: _PoolConfigurator.contract, event: "ReserveFrozen", logs: logs, sub: sub}, nil
}

// WatchReserveFrozen is a free log subscription operation binding the contract event 0x0c4443d258a350d27dc50c378b2ebf165e6469725f786d21b30cab16823f5587.
//
// Solidity: event ReserveFrozen(address indexed asset, bool frozen)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveFrozen(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveFrozen, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveFrozen", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveFrozen)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveFrozen", log); err != nil {
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

// ParseReserveFrozen is a log parse operation binding the contract event 0x0c4443d258a350d27dc50c378b2ebf165e6469725f786d21b30cab16823f5587.
//
// Solidity: event ReserveFrozen(address indexed asset, bool frozen)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveFrozen(log types.Log) (*PoolConfiguratorReserveFrozen, error) {
	event := new(PoolConfiguratorReserveFrozen)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveInitializedIterator is returned from FilterReserveInitialized and is used to iterate over the raw logs and unpacked data for ReserveInitialized events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveInitializedIterator struct {
	Event *PoolConfiguratorReserveInitialized // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveInitialized)
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
		it.Event = new(PoolConfiguratorReserveInitialized)
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
func (it *PoolConfiguratorReserveInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveInitialized represents a ReserveInitialized event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveInitialized struct {
	Asset                       common.Address
	AToken                      common.Address
	StableDebtToken             common.Address
	VariableDebtToken           common.Address
	InterestRateStrategyAddress common.Address
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterReserveInitialized is a free log retrieval operation binding the contract event 0x3a0ca721fc364424566385a1aa271ed508cc2c0949c2272575fb3013a163a45f.
//
// Solidity: event ReserveInitialized(address indexed asset, address indexed aToken, address stableDebtToken, address variableDebtToken, address interestRateStrategyAddress)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveInitialized(opts *bind.FilterOpts, asset []common.Address, aToken []common.Address) (*PoolConfiguratorReserveInitializedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var aTokenRule []interface{}
	for _, aTokenItem := range aToken {
		aTokenRule = append(aTokenRule, aTokenItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveInitialized", assetRule, aTokenRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveInitializedIterator{contract: _PoolConfigurator.contract, event: "ReserveInitialized", logs: logs, sub: sub}, nil
}

// WatchReserveInitialized is a free log subscription operation binding the contract event 0x3a0ca721fc364424566385a1aa271ed508cc2c0949c2272575fb3013a163a45f.
//
// Solidity: event ReserveInitialized(address indexed asset, address indexed aToken, address stableDebtToken, address variableDebtToken, address interestRateStrategyAddress)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveInitialized(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveInitialized, asset []common.Address, aToken []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var aTokenRule []interface{}
	for _, aTokenItem := range aToken {
		aTokenRule = append(aTokenRule, aTokenItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveInitialized", assetRule, aTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveInitialized)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveInitialized", log); err != nil {
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

// ParseReserveInitialized is a log parse operation binding the contract event 0x3a0ca721fc364424566385a1aa271ed508cc2c0949c2272575fb3013a163a45f.
//
// Solidity: event ReserveInitialized(address indexed asset, address indexed aToken, address stableDebtToken, address variableDebtToken, address interestRateStrategyAddress)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveInitialized(log types.Log) (*PoolConfiguratorReserveInitialized, error) {
	event := new(PoolConfiguratorReserveInitialized)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveInterestRateStrategyChangedIterator is returned from FilterReserveInterestRateStrategyChanged and is used to iterate over the raw logs and unpacked data for ReserveInterestRateStrategyChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveInterestRateStrategyChangedIterator struct {
	Event *PoolConfiguratorReserveInterestRateStrategyChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveInterestRateStrategyChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveInterestRateStrategyChanged)
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
		it.Event = new(PoolConfiguratorReserveInterestRateStrategyChanged)
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
func (it *PoolConfiguratorReserveInterestRateStrategyChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveInterestRateStrategyChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveInterestRateStrategyChanged represents a ReserveInterestRateStrategyChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveInterestRateStrategyChanged struct {
	Asset       common.Address
	OldStrategy common.Address
	NewStrategy common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReserveInterestRateStrategyChanged is a free log retrieval operation binding the contract event 0xdb8dada53709ce4988154324196790c2e4a60c377e1256790946f83b87db3c33.
//
// Solidity: event ReserveInterestRateStrategyChanged(address indexed asset, address oldStrategy, address newStrategy)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveInterestRateStrategyChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveInterestRateStrategyChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveInterestRateStrategyChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveInterestRateStrategyChangedIterator{contract: _PoolConfigurator.contract, event: "ReserveInterestRateStrategyChanged", logs: logs, sub: sub}, nil
}

// WatchReserveInterestRateStrategyChanged is a free log subscription operation binding the contract event 0xdb8dada53709ce4988154324196790c2e4a60c377e1256790946f83b87db3c33.
//
// Solidity: event ReserveInterestRateStrategyChanged(address indexed asset, address oldStrategy, address newStrategy)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveInterestRateStrategyChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveInterestRateStrategyChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveInterestRateStrategyChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveInterestRateStrategyChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveInterestRateStrategyChanged", log); err != nil {
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

// ParseReserveInterestRateStrategyChanged is a log parse operation binding the contract event 0xdb8dada53709ce4988154324196790c2e4a60c377e1256790946f83b87db3c33.
//
// Solidity: event ReserveInterestRateStrategyChanged(address indexed asset, address oldStrategy, address newStrategy)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveInterestRateStrategyChanged(log types.Log) (*PoolConfiguratorReserveInterestRateStrategyChanged, error) {
	event := new(PoolConfiguratorReserveInterestRateStrategyChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveInterestRateStrategyChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReservePausedIterator is returned from FilterReservePaused and is used to iterate over the raw logs and unpacked data for ReservePaused events raised by the PoolConfigurator contract.
type PoolConfiguratorReservePausedIterator struct {
	Event *PoolConfiguratorReservePaused // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReservePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReservePaused)
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
		it.Event = new(PoolConfiguratorReservePaused)
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
func (it *PoolConfiguratorReservePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReservePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReservePaused represents a ReservePaused event raised by the PoolConfigurator contract.
type PoolConfiguratorReservePaused struct {
	Asset  common.Address
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReservePaused is a free log retrieval operation binding the contract event 0xe188d542a5f11925d3a3af33703cdd30a43cb3e8066a3cf68b1b57f61a5a94b5.
//
// Solidity: event ReservePaused(address indexed asset, bool paused)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReservePaused(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReservePausedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReservePaused", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReservePausedIterator{contract: _PoolConfigurator.contract, event: "ReservePaused", logs: logs, sub: sub}, nil
}

// WatchReservePaused is a free log subscription operation binding the contract event 0xe188d542a5f11925d3a3af33703cdd30a43cb3e8066a3cf68b1b57f61a5a94b5.
//
// Solidity: event ReservePaused(address indexed asset, bool paused)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReservePaused(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReservePaused, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReservePaused", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReservePaused)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReservePaused", log); err != nil {
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

// ParseReservePaused is a log parse operation binding the contract event 0xe188d542a5f11925d3a3af33703cdd30a43cb3e8066a3cf68b1b57f61a5a94b5.
//
// Solidity: event ReservePaused(address indexed asset, bool paused)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReservePaused(log types.Log) (*PoolConfiguratorReservePaused, error) {
	event := new(PoolConfiguratorReservePaused)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReservePaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorReserveStableRateBorrowingIterator is returned from FilterReserveStableRateBorrowing and is used to iterate over the raw logs and unpacked data for ReserveStableRateBorrowing events raised by the PoolConfigurator contract.
type PoolConfiguratorReserveStableRateBorrowingIterator struct {
	Event *PoolConfiguratorReserveStableRateBorrowing // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorReserveStableRateBorrowingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorReserveStableRateBorrowing)
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
		it.Event = new(PoolConfiguratorReserveStableRateBorrowing)
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
func (it *PoolConfiguratorReserveStableRateBorrowingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorReserveStableRateBorrowingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorReserveStableRateBorrowing represents a ReserveStableRateBorrowing event raised by the PoolConfigurator contract.
type PoolConfiguratorReserveStableRateBorrowing struct {
	Asset   common.Address
	Enabled bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveStableRateBorrowing is a free log retrieval operation binding the contract event 0x0b64d0941719acd363f1a6be3d8525d8ec9d71738f7445aabcd88d7939b472e7.
//
// Solidity: event ReserveStableRateBorrowing(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterReserveStableRateBorrowing(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorReserveStableRateBorrowingIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "ReserveStableRateBorrowing", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorReserveStableRateBorrowingIterator{contract: _PoolConfigurator.contract, event: "ReserveStableRateBorrowing", logs: logs, sub: sub}, nil
}

// WatchReserveStableRateBorrowing is a free log subscription operation binding the contract event 0x0b64d0941719acd363f1a6be3d8525d8ec9d71738f7445aabcd88d7939b472e7.
//
// Solidity: event ReserveStableRateBorrowing(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchReserveStableRateBorrowing(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorReserveStableRateBorrowing, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "ReserveStableRateBorrowing", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorReserveStableRateBorrowing)
				if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveStableRateBorrowing", log); err != nil {
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

// ParseReserveStableRateBorrowing is a log parse operation binding the contract event 0x0b64d0941719acd363f1a6be3d8525d8ec9d71738f7445aabcd88d7939b472e7.
//
// Solidity: event ReserveStableRateBorrowing(address indexed asset, bool enabled)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseReserveStableRateBorrowing(log types.Log) (*PoolConfiguratorReserveStableRateBorrowing, error) {
	event := new(PoolConfiguratorReserveStableRateBorrowing)
	if err := _PoolConfigurator.contract.UnpackLog(event, "ReserveStableRateBorrowing", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorSiloedBorrowingChangedIterator is returned from FilterSiloedBorrowingChanged and is used to iterate over the raw logs and unpacked data for SiloedBorrowingChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorSiloedBorrowingChangedIterator struct {
	Event *PoolConfiguratorSiloedBorrowingChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorSiloedBorrowingChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorSiloedBorrowingChanged)
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
		it.Event = new(PoolConfiguratorSiloedBorrowingChanged)
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
func (it *PoolConfiguratorSiloedBorrowingChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorSiloedBorrowingChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorSiloedBorrowingChanged represents a SiloedBorrowingChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorSiloedBorrowingChanged struct {
	Asset    common.Address
	OldState bool
	NewState bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSiloedBorrowingChanged is a free log retrieval operation binding the contract event 0x842a280b07e8e502a9101f32a3b768ebaba3655556dd674f0831900861fc674b.
//
// Solidity: event SiloedBorrowingChanged(address indexed asset, bool oldState, bool newState)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterSiloedBorrowingChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorSiloedBorrowingChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "SiloedBorrowingChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorSiloedBorrowingChangedIterator{contract: _PoolConfigurator.contract, event: "SiloedBorrowingChanged", logs: logs, sub: sub}, nil
}

// WatchSiloedBorrowingChanged is a free log subscription operation binding the contract event 0x842a280b07e8e502a9101f32a3b768ebaba3655556dd674f0831900861fc674b.
//
// Solidity: event SiloedBorrowingChanged(address indexed asset, bool oldState, bool newState)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchSiloedBorrowingChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorSiloedBorrowingChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "SiloedBorrowingChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorSiloedBorrowingChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "SiloedBorrowingChanged", log); err != nil {
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

// ParseSiloedBorrowingChanged is a log parse operation binding the contract event 0x842a280b07e8e502a9101f32a3b768ebaba3655556dd674f0831900861fc674b.
//
// Solidity: event SiloedBorrowingChanged(address indexed asset, bool oldState, bool newState)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseSiloedBorrowingChanged(log types.Log) (*PoolConfiguratorSiloedBorrowingChanged, error) {
	event := new(PoolConfiguratorSiloedBorrowingChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "SiloedBorrowingChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorStableDebtTokenUpgradedIterator is returned from FilterStableDebtTokenUpgraded and is used to iterate over the raw logs and unpacked data for StableDebtTokenUpgraded events raised by the PoolConfigurator contract.
type PoolConfiguratorStableDebtTokenUpgradedIterator struct {
	Event *PoolConfiguratorStableDebtTokenUpgraded // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorStableDebtTokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorStableDebtTokenUpgraded)
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
		it.Event = new(PoolConfiguratorStableDebtTokenUpgraded)
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
func (it *PoolConfiguratorStableDebtTokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorStableDebtTokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorStableDebtTokenUpgraded represents a StableDebtTokenUpgraded event raised by the PoolConfigurator contract.
type PoolConfiguratorStableDebtTokenUpgraded struct {
	Asset          common.Address
	Proxy          common.Address
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStableDebtTokenUpgraded is a free log retrieval operation binding the contract event 0x7a943a5b6c214bf7726c069a878b1e2a8e7371981d516048b84e03743e67bc28.
//
// Solidity: event StableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterStableDebtTokenUpgraded(opts *bind.FilterOpts, asset []common.Address, proxy []common.Address, implementation []common.Address) (*PoolConfiguratorStableDebtTokenUpgradedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "StableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorStableDebtTokenUpgradedIterator{contract: _PoolConfigurator.contract, event: "StableDebtTokenUpgraded", logs: logs, sub: sub}, nil
}

// WatchStableDebtTokenUpgraded is a free log subscription operation binding the contract event 0x7a943a5b6c214bf7726c069a878b1e2a8e7371981d516048b84e03743e67bc28.
//
// Solidity: event StableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchStableDebtTokenUpgraded(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorStableDebtTokenUpgraded, asset []common.Address, proxy []common.Address, implementation []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "StableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorStableDebtTokenUpgraded)
				if err := _PoolConfigurator.contract.UnpackLog(event, "StableDebtTokenUpgraded", log); err != nil {
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

// ParseStableDebtTokenUpgraded is a log parse operation binding the contract event 0x7a943a5b6c214bf7726c069a878b1e2a8e7371981d516048b84e03743e67bc28.
//
// Solidity: event StableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseStableDebtTokenUpgraded(log types.Log) (*PoolConfiguratorStableDebtTokenUpgraded, error) {
	event := new(PoolConfiguratorStableDebtTokenUpgraded)
	if err := _PoolConfigurator.contract.UnpackLog(event, "StableDebtTokenUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorSupplyCapChangedIterator is returned from FilterSupplyCapChanged and is used to iterate over the raw logs and unpacked data for SupplyCapChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorSupplyCapChangedIterator struct {
	Event *PoolConfiguratorSupplyCapChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorSupplyCapChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorSupplyCapChanged)
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
		it.Event = new(PoolConfiguratorSupplyCapChanged)
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
func (it *PoolConfiguratorSupplyCapChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorSupplyCapChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorSupplyCapChanged represents a SupplyCapChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorSupplyCapChanged struct {
	Asset        common.Address
	OldSupplyCap *big.Int
	NewSupplyCap *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSupplyCapChanged is a free log retrieval operation binding the contract event 0x0263602682188540a2d633561c0b4453b7d8566285e99f9f6018b8ef2facef49.
//
// Solidity: event SupplyCapChanged(address indexed asset, uint256 oldSupplyCap, uint256 newSupplyCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterSupplyCapChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorSupplyCapChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "SupplyCapChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorSupplyCapChangedIterator{contract: _PoolConfigurator.contract, event: "SupplyCapChanged", logs: logs, sub: sub}, nil
}

// WatchSupplyCapChanged is a free log subscription operation binding the contract event 0x0263602682188540a2d633561c0b4453b7d8566285e99f9f6018b8ef2facef49.
//
// Solidity: event SupplyCapChanged(address indexed asset, uint256 oldSupplyCap, uint256 newSupplyCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchSupplyCapChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorSupplyCapChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "SupplyCapChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorSupplyCapChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "SupplyCapChanged", log); err != nil {
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

// ParseSupplyCapChanged is a log parse operation binding the contract event 0x0263602682188540a2d633561c0b4453b7d8566285e99f9f6018b8ef2facef49.
//
// Solidity: event SupplyCapChanged(address indexed asset, uint256 oldSupplyCap, uint256 newSupplyCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseSupplyCapChanged(log types.Log) (*PoolConfiguratorSupplyCapChanged, error) {
	event := new(PoolConfiguratorSupplyCapChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "SupplyCapChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorUnbackedMintCapChangedIterator is returned from FilterUnbackedMintCapChanged and is used to iterate over the raw logs and unpacked data for UnbackedMintCapChanged events raised by the PoolConfigurator contract.
type PoolConfiguratorUnbackedMintCapChangedIterator struct {
	Event *PoolConfiguratorUnbackedMintCapChanged // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorUnbackedMintCapChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorUnbackedMintCapChanged)
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
		it.Event = new(PoolConfiguratorUnbackedMintCapChanged)
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
func (it *PoolConfiguratorUnbackedMintCapChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorUnbackedMintCapChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorUnbackedMintCapChanged represents a UnbackedMintCapChanged event raised by the PoolConfigurator contract.
type PoolConfiguratorUnbackedMintCapChanged struct {
	Asset              common.Address
	OldUnbackedMintCap *big.Int
	NewUnbackedMintCap *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUnbackedMintCapChanged is a free log retrieval operation binding the contract event 0x09808b1fc5abde94edf02fdde393bea0d2e4795999ba31695472848638b5c29f.
//
// Solidity: event UnbackedMintCapChanged(address indexed asset, uint256 oldUnbackedMintCap, uint256 newUnbackedMintCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterUnbackedMintCapChanged(opts *bind.FilterOpts, asset []common.Address) (*PoolConfiguratorUnbackedMintCapChangedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "UnbackedMintCapChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorUnbackedMintCapChangedIterator{contract: _PoolConfigurator.contract, event: "UnbackedMintCapChanged", logs: logs, sub: sub}, nil
}

// WatchUnbackedMintCapChanged is a free log subscription operation binding the contract event 0x09808b1fc5abde94edf02fdde393bea0d2e4795999ba31695472848638b5c29f.
//
// Solidity: event UnbackedMintCapChanged(address indexed asset, uint256 oldUnbackedMintCap, uint256 newUnbackedMintCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchUnbackedMintCapChanged(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorUnbackedMintCapChanged, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "UnbackedMintCapChanged", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorUnbackedMintCapChanged)
				if err := _PoolConfigurator.contract.UnpackLog(event, "UnbackedMintCapChanged", log); err != nil {
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

// ParseUnbackedMintCapChanged is a log parse operation binding the contract event 0x09808b1fc5abde94edf02fdde393bea0d2e4795999ba31695472848638b5c29f.
//
// Solidity: event UnbackedMintCapChanged(address indexed asset, uint256 oldUnbackedMintCap, uint256 newUnbackedMintCap)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseUnbackedMintCapChanged(log types.Log) (*PoolConfiguratorUnbackedMintCapChanged, error) {
	event := new(PoolConfiguratorUnbackedMintCapChanged)
	if err := _PoolConfigurator.contract.UnpackLog(event, "UnbackedMintCapChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolConfiguratorVariableDebtTokenUpgradedIterator is returned from FilterVariableDebtTokenUpgraded and is used to iterate over the raw logs and unpacked data for VariableDebtTokenUpgraded events raised by the PoolConfigurator contract.
type PoolConfiguratorVariableDebtTokenUpgradedIterator struct {
	Event *PoolConfiguratorVariableDebtTokenUpgraded // Event containing the contract specifics and raw log

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
func (it *PoolConfiguratorVariableDebtTokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolConfiguratorVariableDebtTokenUpgraded)
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
		it.Event = new(PoolConfiguratorVariableDebtTokenUpgraded)
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
func (it *PoolConfiguratorVariableDebtTokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolConfiguratorVariableDebtTokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolConfiguratorVariableDebtTokenUpgraded represents a VariableDebtTokenUpgraded event raised by the PoolConfigurator contract.
type PoolConfiguratorVariableDebtTokenUpgraded struct {
	Asset          common.Address
	Proxy          common.Address
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterVariableDebtTokenUpgraded is a free log retrieval operation binding the contract event 0x9439658a562a5c46b1173589df89cf001483d685bad28aedaff4a88656292d81.
//
// Solidity: event VariableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) FilterVariableDebtTokenUpgraded(opts *bind.FilterOpts, asset []common.Address, proxy []common.Address, implementation []common.Address) (*PoolConfiguratorVariableDebtTokenUpgradedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PoolConfigurator.contract.FilterLogs(opts, "VariableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &PoolConfiguratorVariableDebtTokenUpgradedIterator{contract: _PoolConfigurator.contract, event: "VariableDebtTokenUpgraded", logs: logs, sub: sub}, nil
}

// WatchVariableDebtTokenUpgraded is a free log subscription operation binding the contract event 0x9439658a562a5c46b1173589df89cf001483d685bad28aedaff4a88656292d81.
//
// Solidity: event VariableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) WatchVariableDebtTokenUpgraded(opts *bind.WatchOpts, sink chan<- *PoolConfiguratorVariableDebtTokenUpgraded, asset []common.Address, proxy []common.Address, implementation []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PoolConfigurator.contract.WatchLogs(opts, "VariableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolConfiguratorVariableDebtTokenUpgraded)
				if err := _PoolConfigurator.contract.UnpackLog(event, "VariableDebtTokenUpgraded", log); err != nil {
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

// ParseVariableDebtTokenUpgraded is a log parse operation binding the contract event 0x9439658a562a5c46b1173589df89cf001483d685bad28aedaff4a88656292d81.
//
// Solidity: event VariableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_PoolConfigurator *PoolConfiguratorFilterer) ParseVariableDebtTokenUpgraded(log types.Log) (*PoolConfiguratorVariableDebtTokenUpgraded, error) {
	event := new(PoolConfiguratorVariableDebtTokenUpgraded)
	if err := _PoolConfigurator.contract.UnpackLog(event, "VariableDebtTokenUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
