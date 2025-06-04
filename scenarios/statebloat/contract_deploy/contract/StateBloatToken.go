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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy10\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy11\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy12\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy13\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy14\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy15\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy16\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy17\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy18\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy19\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy20\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy21\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy22\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy23\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy24\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy25\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy26\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy27\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy28\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy29\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy3\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy30\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy31\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy32\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy33\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy34\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy35\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy36\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy37\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy38\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy39\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy4\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy40\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy41\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy42\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy43\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy44\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy45\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy46\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy47\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy48\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy49\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy5\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy50\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy51\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy52\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy53\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy54\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy55\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy56\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy57\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy58\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy59\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy6\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy60\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy61\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy62\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy63\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy64\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy65\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy7\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy8\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy9\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"salt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom1\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom10\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom11\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom12\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom13\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom14\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom15\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom16\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom17\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom18\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom19\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom3\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom4\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom5\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom6\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom7\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom8\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom9\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051615d6a380380615d6a833981810160405281019061003191906101f4565b6040518060400160405280601181526020017f537461746520426c6f617420546f6b656e0000000000000000000000000000008152505f90816100749190610453565b506040518060400160405280600381526020017f5342540000000000000000000000000000000000000000000000000000000000815250600190816100b99190610453565b50601260025f6101000a81548160ff021916908360ff160217905550806080818152505060025f9054906101000a900460ff16600a6100f8919061068a565b620f424061010691906106d4565b60038190555060035460045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055503373ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6003546040516101af9190610724565b60405180910390a35061073d565b5f5ffd5b5f819050919050565b6101d3816101c1565b81146101dd575f5ffd5b50565b5f815190506101ee816101ca565b92915050565b5f60208284031215610209576102086101bd565b5b5f610216848285016101e0565b91505092915050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061029a57607f821691505b6020821081036102ad576102ac610256565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261030f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826102d4565b61031986836102d4565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61035461034f61034a846101c1565b610331565b6101c1565b9050919050565b5f819050919050565b61036d8361033a565b6103816103798261035b565b8484546102e0565b825550505050565b5f5f905090565b610398610389565b6103a3818484610364565b505050565b5b818110156103c6576103bb5f82610390565b6001810190506103a9565b5050565b601f82111561040b576103dc816102b3565b6103e5846102c5565b810160208510156103f4578190505b610408610400856102c5565b8301826103a8565b50505b505050565b5f82821c905092915050565b5f61042b5f1984600802610410565b1980831691505092915050565b5f610443838361041c565b9150826002028217905092915050565b61045c8261021f565b67ffffffffffffffff81111561047557610474610229565b5b61047f8254610283565b61048a8282856103ca565b5f60209050601f8311600181146104bb575f84156104a9578287015190505b6104b38582610438565b86555061051a565b601f1984166104c9866102b3565b5f5b828110156104f0578489015182556001820191506020850194506020810190506104cb565b8683101561050d5784890151610509601f89168261041c565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f8160011c9050919050565b5f5f8291508390505b60018511156105a4578086048111156105805761057f610522565b5b600185161561058f5780820291505b808102905061059d8561054f565b9450610564565b94509492505050565b5f826105bc5760019050610677565b816105c9575f9050610677565b81600181146105df57600281146105e957610618565b6001915050610677565b60ff8411156105fb576105fa610522565b5b8360020a91508482111561061257610611610522565b5b50610677565b5060208310610133831016604e8410600b841016171561064d5782820a90508381111561064857610647610522565b5b610677565b61065a848484600161055b565b9250905081840481111561067157610670610522565b5b81810290505b9392505050565b5f60ff82169050919050565b5f610694826101c1565b915061069f8361067e565b92506106cc7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846105ad565b905092915050565b5f6106de826101c1565b91506106e9836101c1565b92508282026106f7816101c1565b9150828204841483151761070e5761070d610522565b5b5092915050565b61071e816101c1565b82525050565b5f6020820190506107375f830184610715565b92915050565b6080516156156107555f395f61482101526156155ff3fe608060405234801561000f575f5ffd5b5060043610610518575f3560e01c8063657b6ef7116102a2578063aaa7af7011610170578063d7419469116100d7578063f26c779b11610090578063f26c779b1461110a578063f5f5738114611128578063f8716f1414611146578063faf35ced14611164578063fe7d599614611194578063ffbf0469146111b257610518565b8063d741946914611032578063dc1d8a9b14611050578063dd62ed3e1461106e578063e2d275301461109e578063e8c927b3146110bc578063eb4329c8146110da57610518565b8063bb9bfe0611610129578063bb9bfe0614610f5a578063bfa0b13314610f8a578063c2be97e314610fa8578063c958d4bf14610fd8578063cfd6686314610ff6578063d101dcd01461101457610518565b8063aaa7af7014610e82578063acc5aee914610ea0578063ad8f422114610ed0578063b1802b9a14610eee578063b2bb360e14610f1e578063b66dd75014610f3c57610518565b80637c72ed0d116102145780638f4a8406116101cd5780638f4a840614610dbc57806395d89b4114610dda5780639c5dfe7314610df85780639df61a2514610e16578063a891d4d414610e34578063a9059cbb14610e5257610518565b80637c72ed0d14610cf65780637dffdc3214610d145780637e54493714610d325780637f34d94b14610d505780638619d60714610d6e5780638789ca6714610d8c57610518565b806374e73fd31161026657806374e73fd314610c3057806374f83d0214610c4e57806377c0209e14610c6c578063792c7f3e14610c8a5780637a319c1814610ca85780637c66673e14610cc657610518565b8063657b6ef714610b64578063672151fe14610b945780636abceacd14610bb25780636c12ed2814610bd057806370a0823114610c0057610518565b80633125f37a116103ea5780634a2e93c611610351578063552a1b561161030a578063552a1b5614610a9e57806358b6a9bd14610ace5780635af92c0514610aec57806361b970eb14610b0a578063639ec53a14610b285780636547317414610b4657610518565b80634a2e93c6146109d85780634b3c7f5f146109f65780634e1dbb8214610a145780634f5e555714610a325780634f7bd75a14610a5057806354c2792014610a8057610518565b80633ea117ce116103a35780633ea117ce146109125780634128a85d14610930578063418b18161461094e57806342937dbd1461096c57806343a6b92d1461098a57806344050a28146109a857610518565b80633125f37a1461083a578063313ce5671461085857806334517f0b1461087657806339e0bd12146108945780633a131990146108c45780633b6be459146108f457610518565b80631a97f18e1161048e5780631fd298ec116104475780631fd298ec1461076257806321ecd7a314610780578063239af2a51461079e57806323b872dd146107bc5780632545d8b7146107ec578063291c3bd71461080a57610518565b80631a97f18e1461068a5780631b17c65c146106a85780631bbffe6f146106d85780631d527cde146107085780631eaa7c52146107265780631f4495891461074457610518565b80631215a3ab116104e05780631215a3ab146105d657806312901b42146105f457806313ebb5ec1461061257806316a3045b1461063057806318160ddd1461064e57806319cf6a911461066c57610518565b80630460faf61461051c57806306fdde031461054c5780630717b1611461056a578063095ea7b3146105885780630cb7a9e7146105b8575b5f5ffd5b61053660048036038101906105319190615139565b6111d0565b60405161054391906151a3565b60405180910390f35b6105546114b0565b604051610561919061522c565b60405180910390f35b61057261153b565b60405161057f919061525b565b60405180910390f35b6105a2600480360381019061059d9190615274565b611543565b6040516105af91906151a3565b60405180910390f35b6105c0611630565b6040516105cd919061525b565b60405180910390f35b6105de611638565b6040516105eb919061525b565b60405180910390f35b6105fc611640565b604051610609919061525b565b60405180910390f35b61061a611648565b604051610627919061525b565b60405180910390f35b610638611650565b604051610645919061525b565b60405180910390f35b610656611658565b604051610663919061525b565b60405180910390f35b61067461165e565b604051610681919061525b565b60405180910390f35b610692611666565b60405161069f919061525b565b60405180910390f35b6106c260048036038101906106bd9190615139565b61166e565b6040516106cf91906151a3565b60405180910390f35b6106f260048036038101906106ed9190615139565b61194e565b6040516106ff91906151a3565b60405180910390f35b610710611c2e565b60405161071d919061525b565b60405180910390f35b61072e611c36565b60405161073b919061525b565b60405180910390f35b61074c611c3e565b604051610759919061525b565b60405180910390f35b61076a611c46565b604051610777919061525b565b60405180910390f35b610788611c4e565b604051610795919061525b565b60405180910390f35b6107a6611c56565b6040516107b3919061525b565b60405180910390f35b6107d660048036038101906107d19190615139565b611c5e565b6040516107e391906151a3565b60405180910390f35b6107f4611f3e565b604051610801919061525b565b60405180910390f35b610824600480360381019061081f9190615139565b611f46565b60405161083191906151a3565b60405180910390f35b610842612226565b60405161084f919061525b565b60405180910390f35b61086061222e565b60405161086d91906152cd565b60405180910390f35b61087e612240565b60405161088b919061525b565b60405180910390f35b6108ae60048036038101906108a99190615139565b612248565b6040516108bb91906151a3565b60405180910390f35b6108de60048036038101906108d99190615139565b612528565b6040516108eb91906151a3565b60405180910390f35b6108fc612808565b604051610909919061525b565b60405180910390f35b61091a612810565b604051610927919061525b565b60405180910390f35b610938612818565b604051610945919061525b565b60405180910390f35b610956612820565b604051610963919061525b565b60405180910390f35b610974612828565b604051610981919061525b565b60405180910390f35b610992612830565b60405161099f919061525b565b60405180910390f35b6109c260048036038101906109bd9190615139565b612838565b6040516109cf91906151a3565b60405180910390f35b6109e0612b18565b6040516109ed919061525b565b60405180910390f35b6109fe612b20565b604051610a0b919061525b565b60405180910390f35b610a1c612b28565b604051610a29919061525b565b60405180910390f35b610a3a612b30565b604051610a47919061525b565b60405180910390f35b610a6a6004803603810190610a659190615139565b612b38565b604051610a7791906151a3565b60405180910390f35b610a88612e18565b604051610a95919061525b565b60405180910390f35b610ab86004803603810190610ab39190615139565b612e20565b604051610ac591906151a3565b60405180910390f35b610ad6613100565b604051610ae3919061525b565b60405180910390f35b610af4613108565b604051610b01919061525b565b60405180910390f35b610b12613110565b604051610b1f919061525b565b60405180910390f35b610b30613118565b604051610b3d919061525b565b60405180910390f35b610b4e613120565b604051610b5b919061525b565b60405180910390f35b610b7e6004803603810190610b799190615139565b613128565b604051610b8b91906151a3565b60405180910390f35b610b9c613408565b604051610ba9919061525b565b60405180910390f35b610bba613410565b604051610bc7919061525b565b60405180910390f35b610bea6004803603810190610be59190615139565b613418565b604051610bf791906151a3565b60405180910390f35b610c1a6004803603810190610c1591906152e6565b6136f8565b604051610c27919061525b565b60405180910390f35b610c3861370d565b604051610c45919061525b565b60405180910390f35b610c56613715565b604051610c63919061525b565b60405180910390f35b610c7461371d565b604051610c81919061525b565b60405180910390f35b610c92613725565b604051610c9f919061525b565b60405180910390f35b610cb061372d565b604051610cbd919061525b565b60405180910390f35b610ce06004803603810190610cdb9190615139565b613735565b604051610ced91906151a3565b60405180910390f35b610cfe613a15565b604051610d0b919061525b565b60405180910390f35b610d1c613a1d565b604051610d29919061525b565b60405180910390f35b610d3a613a25565b604051610d47919061525b565b60405180910390f35b610d58613a2d565b604051610d65919061525b565b60405180910390f35b610d76613a35565b604051610d83919061525b565b60405180910390f35b610da66004803603810190610da19190615139565b613a3d565b604051610db391906151a3565b60405180910390f35b610dc4613d1d565b604051610dd1919061525b565b60405180910390f35b610de2613d25565b604051610def919061522c565b60405180910390f35b610e00613db1565b604051610e0d919061525b565b60405180910390f35b610e1e613db9565b604051610e2b919061525b565b60405180910390f35b610e3c613dc1565b604051610e49919061525b565b60405180910390f35b610e6c6004803603810190610e679190615274565b613dc9565b604051610e7991906151a3565b60405180910390f35b610e8a613f5f565b604051610e97919061525b565b60405180910390f35b610eba6004803603810190610eb59190615139565b613f67565b604051610ec791906151a3565b60405180910390f35b610ed8614247565b604051610ee5919061525b565b60405180910390f35b610f086004803603810190610f039190615139565b61424f565b604051610f1591906151a3565b60405180910390f35b610f2661452f565b604051610f33919061525b565b60405180910390f35b610f44614537565b604051610f51919061525b565b60405180910390f35b610f746004803603810190610f6f9190615139565b61453f565b604051610f8191906151a3565b60405180910390f35b610f9261481f565b604051610f9f919061525b565b60405180910390f35b610fc26004803603810190610fbd9190615139565b614843565b604051610fcf91906151a3565b60405180910390f35b610fe0614b23565b604051610fed919061525b565b60405180910390f35b610ffe614b2b565b60405161100b919061525b565b60405180910390f35b61101c614b33565b604051611029919061525b565b60405180910390f35b61103a614b3b565b604051611047919061525b565b60405180910390f35b611058614b43565b604051611065919061525b565b60405180910390f35b61108860048036038101906110839190615311565b614b4b565b604051611095919061525b565b60405180910390f35b6110a6614b6b565b6040516110b3919061525b565b60405180910390f35b6110c4614b73565b6040516110d1919061525b565b60405180910390f35b6110f460048036038101906110ef9190615139565b614b7b565b60405161110191906151a3565b60405180910390f35b611112614da0565b60405161111f919061525b565b60405180910390f35b611130614da8565b60405161113d919061525b565b60405180910390f35b61114e614db0565b60405161115b919061525b565b60405180910390f35b61117e60048036038101906111799190615139565b614db8565b60405161118b91906151a3565b60405180910390f35b61119c615098565b6040516111a9919061525b565b60405180910390f35b6111ba6150a0565b6040516111c7919061525b565b60405180910390f35b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611251576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161124890615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561130c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161130390615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611358919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546113ab919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611439919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161149d919061525b565b60405180910390a3600190509392505050565b5f80546114bc906154df565b80601f01602080910402602001604051908101604052809291908181526020018280546114e8906154df565b80156115335780601f1061150a57610100808354040283529160200191611533565b820191905f5260205f20905b81548152906001019060200180831161151657829003601f168201915b505050505081565b5f602f905090565b5f8160055f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161161e919061525b565b60405180910390a36001905092915050565b5f601b905090565b5f6005905090565b5f601a905090565b5f601d905090565b5f6033905090565b60035481565b5f6008905090565b5f6003905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156116ef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016116e690615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156117aa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016117a190615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546117f6919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611849919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546118d7919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161193b919061525b565b60405180910390a3600190509392505050565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156119cf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119c690615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611a8a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a8190615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611ad6919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611b29919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611bb7919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051611c1b919061525b565b60405180910390a3600190509392505050565b5f6002905090565b5f6010905090565b5f6041905090565b5f6035905090565b5f602e905090565b5f602c905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611cdf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611cd690615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611d9a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611d9190615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611de6919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611e39919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611ec7919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051611f2b919061525b565b60405180910390a3600190509392505050565b5f6001905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611fc7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611fbe90615559565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612082576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612079906155c1565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546120ce919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612121919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546121af919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612213919061525b565b60405180910390a3600190509392505050565b5f6034905090565b60025f9054906101000a900460ff1681565b5f603c905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156122c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016122c090615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612384576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161237b90615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546123d0919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612423919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546124b1919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612515919061525b565b60405180910390a3600190509392505050565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156125a9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016125a090615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612664576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161265b90615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546126b0919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612703919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612791919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516127f5919061525b565b60405180910390a3600190509392505050565b5f6004905090565b5f600c905090565b5f6006905090565b5f601c905090565b5f6032905090565b5f6023905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156128b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016128b090615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612974576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161296b90615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546129c0919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612a13919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612aa1919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612b05919061525b565b60405180910390a3600190509392505050565b5f6011905090565b5f6021905090565b5f601f905090565b5f6026905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612bb9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612bb090615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612c74576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612c6b90615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612cc0919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612d13919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612da1919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612e05919061525b565b60405180910390a3600190509392505050565b5f603a905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612ea1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612e9890615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612f5c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612f5390615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612fa8919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612ffb919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613089919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516130ed919061525b565b60405180910390a3600190509392505050565b5f6009905090565b5f602a905090565b5f6014905090565b5f6025905090565b5f6013905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156131a9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016131a090615559565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613264576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161325b906155c1565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546132b0919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613303919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613391919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516133f5919061525b565b60405180910390a3600190509392505050565b5f600d905090565b5f6017905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613499576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161349090615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613554576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161354b90615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546135a0919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546135f3919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613681919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516136e5919061525b565b60405180910390a3600190509392505050565b6004602052805f5260405f205f915090505481565b5f6016905090565b5f6039905090565b5f603b905090565b5f6036905090565b5f603f905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156137b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016137ad90615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613871576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161386890615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546138bd919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613910919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461399e919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051613a02919061525b565b60405180910390a3600190509392505050565b5f602b905090565b5f6037905090565b5f600a905090565b5f602d905090565b5f6030905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613abe576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613ab590615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613b79576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613b7090615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613bc5919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613c18919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613ca6919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051613d0a919061525b565b60405180910390a3600190509392505050565b5f600b905090565b60018054613d32906154df565b80601f0160208091040260200160405190810160405280929190818152602001828054613d5e906154df565b8015613da95780601f10613d8057610100808354040283529160200191613da9565b820191905f5260205f20905b815481529060010190602001808311613d8c57829003601f168201915b505050505081565b5f6015905090565b5f6029905090565b5f6038905090565b5f8160045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613e4a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613e4190615399565b60405180910390fd5b8160045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613e96919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613ee9919061547f565b925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051613f4d919061525b565b60405180910390a36001905092915050565b5f6028905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613fe8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613fdf90615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156140a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161409a90615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546140ef919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614142919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546141d0919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051614234919061525b565b60405180910390a3600190509392505050565b5f6040905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156142d0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016142c790615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561438b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161438290615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546143d7919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461442a919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546144b8919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161451c919061525b565b60405180910390a3600190509392505050565b5f6018905090565b5f600f905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156145c0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016145b790615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561467b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161467290615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546146c7919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461471a919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546147a8919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161480c919061525b565b60405180910390a3600190509392505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156148c4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016148bb90615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561497f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161497690615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546149cb919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614a1e919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614aac919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051614b10919061525b565b60405180910390a3600190509392505050565b5f6024905090565b5f601e905090565b5f6031905090565b5f6019905090565b5f6027905090565b6005602052815f5260405f20602052805f5260405f205f91509150505481565b5f603d905090565b5f6022905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614bfc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614bf390615559565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614c48919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614c9b919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614d29919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051614d8d919061525b565b60405180910390a3600190509392505050565b5f6007905090565b5f603e905090565b5f6020905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614e39576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614e3090615399565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614ef4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614eeb90615401565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614f40919061544c565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614f93919061547f565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254615021919061544c565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051615085919061525b565b60405180910390a3600190509392505050565b5f600e905090565b5f6012905090565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6150d5826150ac565b9050919050565b6150e5816150cb565b81146150ef575f5ffd5b50565b5f81359050615100816150dc565b92915050565b5f819050919050565b61511881615106565b8114615122575f5ffd5b50565b5f813590506151338161510f565b92915050565b5f5f5f606084860312156151505761514f6150a8565b5b5f61515d868287016150f2565b935050602061516e868287016150f2565b925050604061517f86828701615125565b9150509250925092565b5f8115159050919050565b61519d81615189565b82525050565b5f6020820190506151b65f830184615194565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6151fe826151bc565b61520881856151c6565b93506152188185602086016151d6565b615221816151e4565b840191505092915050565b5f6020820190508181035f83015261524481846151f4565b905092915050565b61525581615106565b82525050565b5f60208201905061526e5f83018461524c565b92915050565b5f5f6040838503121561528a576152896150a8565b5b5f615297858286016150f2565b92505060206152a885828601615125565b9150509250929050565b5f60ff82169050919050565b6152c7816152b2565b82525050565b5f6020820190506152e05f8301846152be565b92915050565b5f602082840312156152fb576152fa6150a8565b5b5f615308848285016150f2565b91505092915050565b5f5f60408385031215615327576153266150a8565b5b5f615334858286016150f2565b9250506020615345858286016150f2565b9150509250929050565b7f496e73756666696369656e742062616c616e63650000000000000000000000005f82015250565b5f6153836014836151c6565b915061538e8261534f565b602082019050919050565b5f6020820190508181035f8301526153b081615377565b9050919050565b7f496e73756666696369656e7420616c6c6f77616e6365000000000000000000005f82015250565b5f6153eb6016836151c6565b91506153f6826153b7565b602082019050919050565b5f6020820190508181035f830152615418816153df565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61545682615106565b915061546183615106565b92508282039050818111156154795761547861541f565b5b92915050565b5f61548982615106565b915061549483615106565b92508282019050808211156154ac576154ab61541f565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806154f657607f821691505b602082108103615509576155086154b2565b5b50919050565b7f41000000000000000000000000000000000000000000000000000000000000005f82015250565b5f6155436001836151c6565b915061554e8261550f565b602082019050919050565b5f6020820190508181035f83015261557081615537565b9050919050565b7f42000000000000000000000000000000000000000000000000000000000000005f82015250565b5f6155ab6001836151c6565b91506155b682615577565b602082019050919050565b5f6020820190508181035f8301526155d88161559f565b905091905056fea26469706673582212206fbf384dcd110eb4456b07e83a2c8a392433c7852a66dec9f8ebf3989cad1f6764736f6c634300081e0033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend, _salt *big.Int) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend, _salt)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Contract *ContractCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Contract *ContractSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Contract.Contract.Allowance(&_Contract.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Contract *ContractCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Contract.Contract.Allowance(&_Contract.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Contract *ContractCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Contract *ContractSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.BalanceOf(&_Contract.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Contract *ContractCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.BalanceOf(&_Contract.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Contract *ContractCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Contract *ContractSession) Decimals() (uint8, error) {
	return _Contract.Contract.Decimals(&_Contract.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Contract *ContractCallerSession) Decimals() (uint8, error) {
	return _Contract.Contract.Decimals(&_Contract.CallOpts)
}

// Dummy1 is a free data retrieval call binding the contract method 0x2545d8b7.
//
// Solidity: function dummy1() pure returns(uint256)
func (_Contract *ContractCaller) Dummy1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy1 is a free data retrieval call binding the contract method 0x2545d8b7.
//
// Solidity: function dummy1() pure returns(uint256)
func (_Contract *ContractSession) Dummy1() (*big.Int, error) {
	return _Contract.Contract.Dummy1(&_Contract.CallOpts)
}

// Dummy1 is a free data retrieval call binding the contract method 0x2545d8b7.
//
// Solidity: function dummy1() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy1() (*big.Int, error) {
	return _Contract.Contract.Dummy1(&_Contract.CallOpts)
}

// Dummy10 is a free data retrieval call binding the contract method 0x7e544937.
//
// Solidity: function dummy10() pure returns(uint256)
func (_Contract *ContractCaller) Dummy10(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy10")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy10 is a free data retrieval call binding the contract method 0x7e544937.
//
// Solidity: function dummy10() pure returns(uint256)
func (_Contract *ContractSession) Dummy10() (*big.Int, error) {
	return _Contract.Contract.Dummy10(&_Contract.CallOpts)
}

// Dummy10 is a free data retrieval call binding the contract method 0x7e544937.
//
// Solidity: function dummy10() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy10() (*big.Int, error) {
	return _Contract.Contract.Dummy10(&_Contract.CallOpts)
}

// Dummy11 is a free data retrieval call binding the contract method 0x8f4a8406.
//
// Solidity: function dummy11() pure returns(uint256)
func (_Contract *ContractCaller) Dummy11(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy11")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy11 is a free data retrieval call binding the contract method 0x8f4a8406.
//
// Solidity: function dummy11() pure returns(uint256)
func (_Contract *ContractSession) Dummy11() (*big.Int, error) {
	return _Contract.Contract.Dummy11(&_Contract.CallOpts)
}

// Dummy11 is a free data retrieval call binding the contract method 0x8f4a8406.
//
// Solidity: function dummy11() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy11() (*big.Int, error) {
	return _Contract.Contract.Dummy11(&_Contract.CallOpts)
}

// Dummy12 is a free data retrieval call binding the contract method 0x3ea117ce.
//
// Solidity: function dummy12() pure returns(uint256)
func (_Contract *ContractCaller) Dummy12(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy12")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy12 is a free data retrieval call binding the contract method 0x3ea117ce.
//
// Solidity: function dummy12() pure returns(uint256)
func (_Contract *ContractSession) Dummy12() (*big.Int, error) {
	return _Contract.Contract.Dummy12(&_Contract.CallOpts)
}

// Dummy12 is a free data retrieval call binding the contract method 0x3ea117ce.
//
// Solidity: function dummy12() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy12() (*big.Int, error) {
	return _Contract.Contract.Dummy12(&_Contract.CallOpts)
}

// Dummy13 is a free data retrieval call binding the contract method 0x672151fe.
//
// Solidity: function dummy13() pure returns(uint256)
func (_Contract *ContractCaller) Dummy13(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy13")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy13 is a free data retrieval call binding the contract method 0x672151fe.
//
// Solidity: function dummy13() pure returns(uint256)
func (_Contract *ContractSession) Dummy13() (*big.Int, error) {
	return _Contract.Contract.Dummy13(&_Contract.CallOpts)
}

// Dummy13 is a free data retrieval call binding the contract method 0x672151fe.
//
// Solidity: function dummy13() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy13() (*big.Int, error) {
	return _Contract.Contract.Dummy13(&_Contract.CallOpts)
}

// Dummy14 is a free data retrieval call binding the contract method 0xfe7d5996.
//
// Solidity: function dummy14() pure returns(uint256)
func (_Contract *ContractCaller) Dummy14(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy14")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy14 is a free data retrieval call binding the contract method 0xfe7d5996.
//
// Solidity: function dummy14() pure returns(uint256)
func (_Contract *ContractSession) Dummy14() (*big.Int, error) {
	return _Contract.Contract.Dummy14(&_Contract.CallOpts)
}

// Dummy14 is a free data retrieval call binding the contract method 0xfe7d5996.
//
// Solidity: function dummy14() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy14() (*big.Int, error) {
	return _Contract.Contract.Dummy14(&_Contract.CallOpts)
}

// Dummy15 is a free data retrieval call binding the contract method 0xb66dd750.
//
// Solidity: function dummy15() pure returns(uint256)
func (_Contract *ContractCaller) Dummy15(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy15")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy15 is a free data retrieval call binding the contract method 0xb66dd750.
//
// Solidity: function dummy15() pure returns(uint256)
func (_Contract *ContractSession) Dummy15() (*big.Int, error) {
	return _Contract.Contract.Dummy15(&_Contract.CallOpts)
}

// Dummy15 is a free data retrieval call binding the contract method 0xb66dd750.
//
// Solidity: function dummy15() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy15() (*big.Int, error) {
	return _Contract.Contract.Dummy15(&_Contract.CallOpts)
}

// Dummy16 is a free data retrieval call binding the contract method 0x1eaa7c52.
//
// Solidity: function dummy16() pure returns(uint256)
func (_Contract *ContractCaller) Dummy16(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy16")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy16 is a free data retrieval call binding the contract method 0x1eaa7c52.
//
// Solidity: function dummy16() pure returns(uint256)
func (_Contract *ContractSession) Dummy16() (*big.Int, error) {
	return _Contract.Contract.Dummy16(&_Contract.CallOpts)
}

// Dummy16 is a free data retrieval call binding the contract method 0x1eaa7c52.
//
// Solidity: function dummy16() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy16() (*big.Int, error) {
	return _Contract.Contract.Dummy16(&_Contract.CallOpts)
}

// Dummy17 is a free data retrieval call binding the contract method 0x4a2e93c6.
//
// Solidity: function dummy17() pure returns(uint256)
func (_Contract *ContractCaller) Dummy17(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy17")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy17 is a free data retrieval call binding the contract method 0x4a2e93c6.
//
// Solidity: function dummy17() pure returns(uint256)
func (_Contract *ContractSession) Dummy17() (*big.Int, error) {
	return _Contract.Contract.Dummy17(&_Contract.CallOpts)
}

// Dummy17 is a free data retrieval call binding the contract method 0x4a2e93c6.
//
// Solidity: function dummy17() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy17() (*big.Int, error) {
	return _Contract.Contract.Dummy17(&_Contract.CallOpts)
}

// Dummy18 is a free data retrieval call binding the contract method 0xffbf0469.
//
// Solidity: function dummy18() pure returns(uint256)
func (_Contract *ContractCaller) Dummy18(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy18")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy18 is a free data retrieval call binding the contract method 0xffbf0469.
//
// Solidity: function dummy18() pure returns(uint256)
func (_Contract *ContractSession) Dummy18() (*big.Int, error) {
	return _Contract.Contract.Dummy18(&_Contract.CallOpts)
}

// Dummy18 is a free data retrieval call binding the contract method 0xffbf0469.
//
// Solidity: function dummy18() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy18() (*big.Int, error) {
	return _Contract.Contract.Dummy18(&_Contract.CallOpts)
}

// Dummy19 is a free data retrieval call binding the contract method 0x65473174.
//
// Solidity: function dummy19() pure returns(uint256)
func (_Contract *ContractCaller) Dummy19(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy19")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy19 is a free data retrieval call binding the contract method 0x65473174.
//
// Solidity: function dummy19() pure returns(uint256)
func (_Contract *ContractSession) Dummy19() (*big.Int, error) {
	return _Contract.Contract.Dummy19(&_Contract.CallOpts)
}

// Dummy19 is a free data retrieval call binding the contract method 0x65473174.
//
// Solidity: function dummy19() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy19() (*big.Int, error) {
	return _Contract.Contract.Dummy19(&_Contract.CallOpts)
}

// Dummy2 is a free data retrieval call binding the contract method 0x1d527cde.
//
// Solidity: function dummy2() pure returns(uint256)
func (_Contract *ContractCaller) Dummy2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy2 is a free data retrieval call binding the contract method 0x1d527cde.
//
// Solidity: function dummy2() pure returns(uint256)
func (_Contract *ContractSession) Dummy2() (*big.Int, error) {
	return _Contract.Contract.Dummy2(&_Contract.CallOpts)
}

// Dummy2 is a free data retrieval call binding the contract method 0x1d527cde.
//
// Solidity: function dummy2() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy2() (*big.Int, error) {
	return _Contract.Contract.Dummy2(&_Contract.CallOpts)
}

// Dummy20 is a free data retrieval call binding the contract method 0x61b970eb.
//
// Solidity: function dummy20() pure returns(uint256)
func (_Contract *ContractCaller) Dummy20(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy20")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy20 is a free data retrieval call binding the contract method 0x61b970eb.
//
// Solidity: function dummy20() pure returns(uint256)
func (_Contract *ContractSession) Dummy20() (*big.Int, error) {
	return _Contract.Contract.Dummy20(&_Contract.CallOpts)
}

// Dummy20 is a free data retrieval call binding the contract method 0x61b970eb.
//
// Solidity: function dummy20() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy20() (*big.Int, error) {
	return _Contract.Contract.Dummy20(&_Contract.CallOpts)
}

// Dummy21 is a free data retrieval call binding the contract method 0x9c5dfe73.
//
// Solidity: function dummy21() pure returns(uint256)
func (_Contract *ContractCaller) Dummy21(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy21")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy21 is a free data retrieval call binding the contract method 0x9c5dfe73.
//
// Solidity: function dummy21() pure returns(uint256)
func (_Contract *ContractSession) Dummy21() (*big.Int, error) {
	return _Contract.Contract.Dummy21(&_Contract.CallOpts)
}

// Dummy21 is a free data retrieval call binding the contract method 0x9c5dfe73.
//
// Solidity: function dummy21() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy21() (*big.Int, error) {
	return _Contract.Contract.Dummy21(&_Contract.CallOpts)
}

// Dummy22 is a free data retrieval call binding the contract method 0x74e73fd3.
//
// Solidity: function dummy22() pure returns(uint256)
func (_Contract *ContractCaller) Dummy22(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy22")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy22 is a free data retrieval call binding the contract method 0x74e73fd3.
//
// Solidity: function dummy22() pure returns(uint256)
func (_Contract *ContractSession) Dummy22() (*big.Int, error) {
	return _Contract.Contract.Dummy22(&_Contract.CallOpts)
}

// Dummy22 is a free data retrieval call binding the contract method 0x74e73fd3.
//
// Solidity: function dummy22() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy22() (*big.Int, error) {
	return _Contract.Contract.Dummy22(&_Contract.CallOpts)
}

// Dummy23 is a free data retrieval call binding the contract method 0x6abceacd.
//
// Solidity: function dummy23() pure returns(uint256)
func (_Contract *ContractCaller) Dummy23(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy23")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy23 is a free data retrieval call binding the contract method 0x6abceacd.
//
// Solidity: function dummy23() pure returns(uint256)
func (_Contract *ContractSession) Dummy23() (*big.Int, error) {
	return _Contract.Contract.Dummy23(&_Contract.CallOpts)
}

// Dummy23 is a free data retrieval call binding the contract method 0x6abceacd.
//
// Solidity: function dummy23() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy23() (*big.Int, error) {
	return _Contract.Contract.Dummy23(&_Contract.CallOpts)
}

// Dummy24 is a free data retrieval call binding the contract method 0xb2bb360e.
//
// Solidity: function dummy24() pure returns(uint256)
func (_Contract *ContractCaller) Dummy24(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy24")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy24 is a free data retrieval call binding the contract method 0xb2bb360e.
//
// Solidity: function dummy24() pure returns(uint256)
func (_Contract *ContractSession) Dummy24() (*big.Int, error) {
	return _Contract.Contract.Dummy24(&_Contract.CallOpts)
}

// Dummy24 is a free data retrieval call binding the contract method 0xb2bb360e.
//
// Solidity: function dummy24() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy24() (*big.Int, error) {
	return _Contract.Contract.Dummy24(&_Contract.CallOpts)
}

// Dummy25 is a free data retrieval call binding the contract method 0xd7419469.
//
// Solidity: function dummy25() pure returns(uint256)
func (_Contract *ContractCaller) Dummy25(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy25")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy25 is a free data retrieval call binding the contract method 0xd7419469.
//
// Solidity: function dummy25() pure returns(uint256)
func (_Contract *ContractSession) Dummy25() (*big.Int, error) {
	return _Contract.Contract.Dummy25(&_Contract.CallOpts)
}

// Dummy25 is a free data retrieval call binding the contract method 0xd7419469.
//
// Solidity: function dummy25() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy25() (*big.Int, error) {
	return _Contract.Contract.Dummy25(&_Contract.CallOpts)
}

// Dummy26 is a free data retrieval call binding the contract method 0x12901b42.
//
// Solidity: function dummy26() pure returns(uint256)
func (_Contract *ContractCaller) Dummy26(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy26")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy26 is a free data retrieval call binding the contract method 0x12901b42.
//
// Solidity: function dummy26() pure returns(uint256)
func (_Contract *ContractSession) Dummy26() (*big.Int, error) {
	return _Contract.Contract.Dummy26(&_Contract.CallOpts)
}

// Dummy26 is a free data retrieval call binding the contract method 0x12901b42.
//
// Solidity: function dummy26() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy26() (*big.Int, error) {
	return _Contract.Contract.Dummy26(&_Contract.CallOpts)
}

// Dummy27 is a free data retrieval call binding the contract method 0x0cb7a9e7.
//
// Solidity: function dummy27() pure returns(uint256)
func (_Contract *ContractCaller) Dummy27(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy27")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy27 is a free data retrieval call binding the contract method 0x0cb7a9e7.
//
// Solidity: function dummy27() pure returns(uint256)
func (_Contract *ContractSession) Dummy27() (*big.Int, error) {
	return _Contract.Contract.Dummy27(&_Contract.CallOpts)
}

// Dummy27 is a free data retrieval call binding the contract method 0x0cb7a9e7.
//
// Solidity: function dummy27() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy27() (*big.Int, error) {
	return _Contract.Contract.Dummy27(&_Contract.CallOpts)
}

// Dummy28 is a free data retrieval call binding the contract method 0x418b1816.
//
// Solidity: function dummy28() pure returns(uint256)
func (_Contract *ContractCaller) Dummy28(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy28")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy28 is a free data retrieval call binding the contract method 0x418b1816.
//
// Solidity: function dummy28() pure returns(uint256)
func (_Contract *ContractSession) Dummy28() (*big.Int, error) {
	return _Contract.Contract.Dummy28(&_Contract.CallOpts)
}

// Dummy28 is a free data retrieval call binding the contract method 0x418b1816.
//
// Solidity: function dummy28() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy28() (*big.Int, error) {
	return _Contract.Contract.Dummy28(&_Contract.CallOpts)
}

// Dummy29 is a free data retrieval call binding the contract method 0x13ebb5ec.
//
// Solidity: function dummy29() pure returns(uint256)
func (_Contract *ContractCaller) Dummy29(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy29")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy29 is a free data retrieval call binding the contract method 0x13ebb5ec.
//
// Solidity: function dummy29() pure returns(uint256)
func (_Contract *ContractSession) Dummy29() (*big.Int, error) {
	return _Contract.Contract.Dummy29(&_Contract.CallOpts)
}

// Dummy29 is a free data retrieval call binding the contract method 0x13ebb5ec.
//
// Solidity: function dummy29() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy29() (*big.Int, error) {
	return _Contract.Contract.Dummy29(&_Contract.CallOpts)
}

// Dummy3 is a free data retrieval call binding the contract method 0x1a97f18e.
//
// Solidity: function dummy3() pure returns(uint256)
func (_Contract *ContractCaller) Dummy3(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy3")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy3 is a free data retrieval call binding the contract method 0x1a97f18e.
//
// Solidity: function dummy3() pure returns(uint256)
func (_Contract *ContractSession) Dummy3() (*big.Int, error) {
	return _Contract.Contract.Dummy3(&_Contract.CallOpts)
}

// Dummy3 is a free data retrieval call binding the contract method 0x1a97f18e.
//
// Solidity: function dummy3() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy3() (*big.Int, error) {
	return _Contract.Contract.Dummy3(&_Contract.CallOpts)
}

// Dummy30 is a free data retrieval call binding the contract method 0xcfd66863.
//
// Solidity: function dummy30() pure returns(uint256)
func (_Contract *ContractCaller) Dummy30(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy30")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy30 is a free data retrieval call binding the contract method 0xcfd66863.
//
// Solidity: function dummy30() pure returns(uint256)
func (_Contract *ContractSession) Dummy30() (*big.Int, error) {
	return _Contract.Contract.Dummy30(&_Contract.CallOpts)
}

// Dummy30 is a free data retrieval call binding the contract method 0xcfd66863.
//
// Solidity: function dummy30() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy30() (*big.Int, error) {
	return _Contract.Contract.Dummy30(&_Contract.CallOpts)
}

// Dummy31 is a free data retrieval call binding the contract method 0x4e1dbb82.
//
// Solidity: function dummy31() pure returns(uint256)
func (_Contract *ContractCaller) Dummy31(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy31")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy31 is a free data retrieval call binding the contract method 0x4e1dbb82.
//
// Solidity: function dummy31() pure returns(uint256)
func (_Contract *ContractSession) Dummy31() (*big.Int, error) {
	return _Contract.Contract.Dummy31(&_Contract.CallOpts)
}

// Dummy31 is a free data retrieval call binding the contract method 0x4e1dbb82.
//
// Solidity: function dummy31() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy31() (*big.Int, error) {
	return _Contract.Contract.Dummy31(&_Contract.CallOpts)
}

// Dummy32 is a free data retrieval call binding the contract method 0xf8716f14.
//
// Solidity: function dummy32() pure returns(uint256)
func (_Contract *ContractCaller) Dummy32(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy32")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy32 is a free data retrieval call binding the contract method 0xf8716f14.
//
// Solidity: function dummy32() pure returns(uint256)
func (_Contract *ContractSession) Dummy32() (*big.Int, error) {
	return _Contract.Contract.Dummy32(&_Contract.CallOpts)
}

// Dummy32 is a free data retrieval call binding the contract method 0xf8716f14.
//
// Solidity: function dummy32() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy32() (*big.Int, error) {
	return _Contract.Contract.Dummy32(&_Contract.CallOpts)
}

// Dummy33 is a free data retrieval call binding the contract method 0x4b3c7f5f.
//
// Solidity: function dummy33() pure returns(uint256)
func (_Contract *ContractCaller) Dummy33(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy33")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy33 is a free data retrieval call binding the contract method 0x4b3c7f5f.
//
// Solidity: function dummy33() pure returns(uint256)
func (_Contract *ContractSession) Dummy33() (*big.Int, error) {
	return _Contract.Contract.Dummy33(&_Contract.CallOpts)
}

// Dummy33 is a free data retrieval call binding the contract method 0x4b3c7f5f.
//
// Solidity: function dummy33() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy33() (*big.Int, error) {
	return _Contract.Contract.Dummy33(&_Contract.CallOpts)
}

// Dummy34 is a free data retrieval call binding the contract method 0xe8c927b3.
//
// Solidity: function dummy34() pure returns(uint256)
func (_Contract *ContractCaller) Dummy34(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy34")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy34 is a free data retrieval call binding the contract method 0xe8c927b3.
//
// Solidity: function dummy34() pure returns(uint256)
func (_Contract *ContractSession) Dummy34() (*big.Int, error) {
	return _Contract.Contract.Dummy34(&_Contract.CallOpts)
}

// Dummy34 is a free data retrieval call binding the contract method 0xe8c927b3.
//
// Solidity: function dummy34() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy34() (*big.Int, error) {
	return _Contract.Contract.Dummy34(&_Contract.CallOpts)
}

// Dummy35 is a free data retrieval call binding the contract method 0x43a6b92d.
//
// Solidity: function dummy35() pure returns(uint256)
func (_Contract *ContractCaller) Dummy35(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy35")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy35 is a free data retrieval call binding the contract method 0x43a6b92d.
//
// Solidity: function dummy35() pure returns(uint256)
func (_Contract *ContractSession) Dummy35() (*big.Int, error) {
	return _Contract.Contract.Dummy35(&_Contract.CallOpts)
}

// Dummy35 is a free data retrieval call binding the contract method 0x43a6b92d.
//
// Solidity: function dummy35() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy35() (*big.Int, error) {
	return _Contract.Contract.Dummy35(&_Contract.CallOpts)
}

// Dummy36 is a free data retrieval call binding the contract method 0xc958d4bf.
//
// Solidity: function dummy36() pure returns(uint256)
func (_Contract *ContractCaller) Dummy36(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy36")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy36 is a free data retrieval call binding the contract method 0xc958d4bf.
//
// Solidity: function dummy36() pure returns(uint256)
func (_Contract *ContractSession) Dummy36() (*big.Int, error) {
	return _Contract.Contract.Dummy36(&_Contract.CallOpts)
}

// Dummy36 is a free data retrieval call binding the contract method 0xc958d4bf.
//
// Solidity: function dummy36() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy36() (*big.Int, error) {
	return _Contract.Contract.Dummy36(&_Contract.CallOpts)
}

// Dummy37 is a free data retrieval call binding the contract method 0x639ec53a.
//
// Solidity: function dummy37() pure returns(uint256)
func (_Contract *ContractCaller) Dummy37(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy37")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy37 is a free data retrieval call binding the contract method 0x639ec53a.
//
// Solidity: function dummy37() pure returns(uint256)
func (_Contract *ContractSession) Dummy37() (*big.Int, error) {
	return _Contract.Contract.Dummy37(&_Contract.CallOpts)
}

// Dummy37 is a free data retrieval call binding the contract method 0x639ec53a.
//
// Solidity: function dummy37() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy37() (*big.Int, error) {
	return _Contract.Contract.Dummy37(&_Contract.CallOpts)
}

// Dummy38 is a free data retrieval call binding the contract method 0x4f5e5557.
//
// Solidity: function dummy38() pure returns(uint256)
func (_Contract *ContractCaller) Dummy38(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy38")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy38 is a free data retrieval call binding the contract method 0x4f5e5557.
//
// Solidity: function dummy38() pure returns(uint256)
func (_Contract *ContractSession) Dummy38() (*big.Int, error) {
	return _Contract.Contract.Dummy38(&_Contract.CallOpts)
}

// Dummy38 is a free data retrieval call binding the contract method 0x4f5e5557.
//
// Solidity: function dummy38() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy38() (*big.Int, error) {
	return _Contract.Contract.Dummy38(&_Contract.CallOpts)
}

// Dummy39 is a free data retrieval call binding the contract method 0xdc1d8a9b.
//
// Solidity: function dummy39() pure returns(uint256)
func (_Contract *ContractCaller) Dummy39(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy39")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy39 is a free data retrieval call binding the contract method 0xdc1d8a9b.
//
// Solidity: function dummy39() pure returns(uint256)
func (_Contract *ContractSession) Dummy39() (*big.Int, error) {
	return _Contract.Contract.Dummy39(&_Contract.CallOpts)
}

// Dummy39 is a free data retrieval call binding the contract method 0xdc1d8a9b.
//
// Solidity: function dummy39() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy39() (*big.Int, error) {
	return _Contract.Contract.Dummy39(&_Contract.CallOpts)
}

// Dummy4 is a free data retrieval call binding the contract method 0x3b6be459.
//
// Solidity: function dummy4() pure returns(uint256)
func (_Contract *ContractCaller) Dummy4(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy4")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy4 is a free data retrieval call binding the contract method 0x3b6be459.
//
// Solidity: function dummy4() pure returns(uint256)
func (_Contract *ContractSession) Dummy4() (*big.Int, error) {
	return _Contract.Contract.Dummy4(&_Contract.CallOpts)
}

// Dummy4 is a free data retrieval call binding the contract method 0x3b6be459.
//
// Solidity: function dummy4() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy4() (*big.Int, error) {
	return _Contract.Contract.Dummy4(&_Contract.CallOpts)
}

// Dummy40 is a free data retrieval call binding the contract method 0xaaa7af70.
//
// Solidity: function dummy40() pure returns(uint256)
func (_Contract *ContractCaller) Dummy40(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy40")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy40 is a free data retrieval call binding the contract method 0xaaa7af70.
//
// Solidity: function dummy40() pure returns(uint256)
func (_Contract *ContractSession) Dummy40() (*big.Int, error) {
	return _Contract.Contract.Dummy40(&_Contract.CallOpts)
}

// Dummy40 is a free data retrieval call binding the contract method 0xaaa7af70.
//
// Solidity: function dummy40() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy40() (*big.Int, error) {
	return _Contract.Contract.Dummy40(&_Contract.CallOpts)
}

// Dummy41 is a free data retrieval call binding the contract method 0x9df61a25.
//
// Solidity: function dummy41() pure returns(uint256)
func (_Contract *ContractCaller) Dummy41(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy41")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy41 is a free data retrieval call binding the contract method 0x9df61a25.
//
// Solidity: function dummy41() pure returns(uint256)
func (_Contract *ContractSession) Dummy41() (*big.Int, error) {
	return _Contract.Contract.Dummy41(&_Contract.CallOpts)
}

// Dummy41 is a free data retrieval call binding the contract method 0x9df61a25.
//
// Solidity: function dummy41() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy41() (*big.Int, error) {
	return _Contract.Contract.Dummy41(&_Contract.CallOpts)
}

// Dummy42 is a free data retrieval call binding the contract method 0x5af92c05.
//
// Solidity: function dummy42() pure returns(uint256)
func (_Contract *ContractCaller) Dummy42(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy42")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy42 is a free data retrieval call binding the contract method 0x5af92c05.
//
// Solidity: function dummy42() pure returns(uint256)
func (_Contract *ContractSession) Dummy42() (*big.Int, error) {
	return _Contract.Contract.Dummy42(&_Contract.CallOpts)
}

// Dummy42 is a free data retrieval call binding the contract method 0x5af92c05.
//
// Solidity: function dummy42() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy42() (*big.Int, error) {
	return _Contract.Contract.Dummy42(&_Contract.CallOpts)
}

// Dummy43 is a free data retrieval call binding the contract method 0x7c72ed0d.
//
// Solidity: function dummy43() pure returns(uint256)
func (_Contract *ContractCaller) Dummy43(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy43")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy43 is a free data retrieval call binding the contract method 0x7c72ed0d.
//
// Solidity: function dummy43() pure returns(uint256)
func (_Contract *ContractSession) Dummy43() (*big.Int, error) {
	return _Contract.Contract.Dummy43(&_Contract.CallOpts)
}

// Dummy43 is a free data retrieval call binding the contract method 0x7c72ed0d.
//
// Solidity: function dummy43() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy43() (*big.Int, error) {
	return _Contract.Contract.Dummy43(&_Contract.CallOpts)
}

// Dummy44 is a free data retrieval call binding the contract method 0x239af2a5.
//
// Solidity: function dummy44() pure returns(uint256)
func (_Contract *ContractCaller) Dummy44(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy44")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy44 is a free data retrieval call binding the contract method 0x239af2a5.
//
// Solidity: function dummy44() pure returns(uint256)
func (_Contract *ContractSession) Dummy44() (*big.Int, error) {
	return _Contract.Contract.Dummy44(&_Contract.CallOpts)
}

// Dummy44 is a free data retrieval call binding the contract method 0x239af2a5.
//
// Solidity: function dummy44() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy44() (*big.Int, error) {
	return _Contract.Contract.Dummy44(&_Contract.CallOpts)
}

// Dummy45 is a free data retrieval call binding the contract method 0x7f34d94b.
//
// Solidity: function dummy45() pure returns(uint256)
func (_Contract *ContractCaller) Dummy45(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy45")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy45 is a free data retrieval call binding the contract method 0x7f34d94b.
//
// Solidity: function dummy45() pure returns(uint256)
func (_Contract *ContractSession) Dummy45() (*big.Int, error) {
	return _Contract.Contract.Dummy45(&_Contract.CallOpts)
}

// Dummy45 is a free data retrieval call binding the contract method 0x7f34d94b.
//
// Solidity: function dummy45() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy45() (*big.Int, error) {
	return _Contract.Contract.Dummy45(&_Contract.CallOpts)
}

// Dummy46 is a free data retrieval call binding the contract method 0x21ecd7a3.
//
// Solidity: function dummy46() pure returns(uint256)
func (_Contract *ContractCaller) Dummy46(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy46")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy46 is a free data retrieval call binding the contract method 0x21ecd7a3.
//
// Solidity: function dummy46() pure returns(uint256)
func (_Contract *ContractSession) Dummy46() (*big.Int, error) {
	return _Contract.Contract.Dummy46(&_Contract.CallOpts)
}

// Dummy46 is a free data retrieval call binding the contract method 0x21ecd7a3.
//
// Solidity: function dummy46() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy46() (*big.Int, error) {
	return _Contract.Contract.Dummy46(&_Contract.CallOpts)
}

// Dummy47 is a free data retrieval call binding the contract method 0x0717b161.
//
// Solidity: function dummy47() pure returns(uint256)
func (_Contract *ContractCaller) Dummy47(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy47")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy47 is a free data retrieval call binding the contract method 0x0717b161.
//
// Solidity: function dummy47() pure returns(uint256)
func (_Contract *ContractSession) Dummy47() (*big.Int, error) {
	return _Contract.Contract.Dummy47(&_Contract.CallOpts)
}

// Dummy47 is a free data retrieval call binding the contract method 0x0717b161.
//
// Solidity: function dummy47() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy47() (*big.Int, error) {
	return _Contract.Contract.Dummy47(&_Contract.CallOpts)
}

// Dummy48 is a free data retrieval call binding the contract method 0x8619d607.
//
// Solidity: function dummy48() pure returns(uint256)
func (_Contract *ContractCaller) Dummy48(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy48")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy48 is a free data retrieval call binding the contract method 0x8619d607.
//
// Solidity: function dummy48() pure returns(uint256)
func (_Contract *ContractSession) Dummy48() (*big.Int, error) {
	return _Contract.Contract.Dummy48(&_Contract.CallOpts)
}

// Dummy48 is a free data retrieval call binding the contract method 0x8619d607.
//
// Solidity: function dummy48() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy48() (*big.Int, error) {
	return _Contract.Contract.Dummy48(&_Contract.CallOpts)
}

// Dummy49 is a free data retrieval call binding the contract method 0xd101dcd0.
//
// Solidity: function dummy49() pure returns(uint256)
func (_Contract *ContractCaller) Dummy49(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy49")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy49 is a free data retrieval call binding the contract method 0xd101dcd0.
//
// Solidity: function dummy49() pure returns(uint256)
func (_Contract *ContractSession) Dummy49() (*big.Int, error) {
	return _Contract.Contract.Dummy49(&_Contract.CallOpts)
}

// Dummy49 is a free data retrieval call binding the contract method 0xd101dcd0.
//
// Solidity: function dummy49() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy49() (*big.Int, error) {
	return _Contract.Contract.Dummy49(&_Contract.CallOpts)
}

// Dummy5 is a free data retrieval call binding the contract method 0x1215a3ab.
//
// Solidity: function dummy5() pure returns(uint256)
func (_Contract *ContractCaller) Dummy5(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy5")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy5 is a free data retrieval call binding the contract method 0x1215a3ab.
//
// Solidity: function dummy5() pure returns(uint256)
func (_Contract *ContractSession) Dummy5() (*big.Int, error) {
	return _Contract.Contract.Dummy5(&_Contract.CallOpts)
}

// Dummy5 is a free data retrieval call binding the contract method 0x1215a3ab.
//
// Solidity: function dummy5() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy5() (*big.Int, error) {
	return _Contract.Contract.Dummy5(&_Contract.CallOpts)
}

// Dummy50 is a free data retrieval call binding the contract method 0x42937dbd.
//
// Solidity: function dummy50() pure returns(uint256)
func (_Contract *ContractCaller) Dummy50(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy50")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy50 is a free data retrieval call binding the contract method 0x42937dbd.
//
// Solidity: function dummy50() pure returns(uint256)
func (_Contract *ContractSession) Dummy50() (*big.Int, error) {
	return _Contract.Contract.Dummy50(&_Contract.CallOpts)
}

// Dummy50 is a free data retrieval call binding the contract method 0x42937dbd.
//
// Solidity: function dummy50() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy50() (*big.Int, error) {
	return _Contract.Contract.Dummy50(&_Contract.CallOpts)
}

// Dummy51 is a free data retrieval call binding the contract method 0x16a3045b.
//
// Solidity: function dummy51() pure returns(uint256)
func (_Contract *ContractCaller) Dummy51(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy51")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy51 is a free data retrieval call binding the contract method 0x16a3045b.
//
// Solidity: function dummy51() pure returns(uint256)
func (_Contract *ContractSession) Dummy51() (*big.Int, error) {
	return _Contract.Contract.Dummy51(&_Contract.CallOpts)
}

// Dummy51 is a free data retrieval call binding the contract method 0x16a3045b.
//
// Solidity: function dummy51() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy51() (*big.Int, error) {
	return _Contract.Contract.Dummy51(&_Contract.CallOpts)
}

// Dummy52 is a free data retrieval call binding the contract method 0x3125f37a.
//
// Solidity: function dummy52() pure returns(uint256)
func (_Contract *ContractCaller) Dummy52(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy52")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy52 is a free data retrieval call binding the contract method 0x3125f37a.
//
// Solidity: function dummy52() pure returns(uint256)
func (_Contract *ContractSession) Dummy52() (*big.Int, error) {
	return _Contract.Contract.Dummy52(&_Contract.CallOpts)
}

// Dummy52 is a free data retrieval call binding the contract method 0x3125f37a.
//
// Solidity: function dummy52() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy52() (*big.Int, error) {
	return _Contract.Contract.Dummy52(&_Contract.CallOpts)
}

// Dummy53 is a free data retrieval call binding the contract method 0x1fd298ec.
//
// Solidity: function dummy53() pure returns(uint256)
func (_Contract *ContractCaller) Dummy53(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy53")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy53 is a free data retrieval call binding the contract method 0x1fd298ec.
//
// Solidity: function dummy53() pure returns(uint256)
func (_Contract *ContractSession) Dummy53() (*big.Int, error) {
	return _Contract.Contract.Dummy53(&_Contract.CallOpts)
}

// Dummy53 is a free data retrieval call binding the contract method 0x1fd298ec.
//
// Solidity: function dummy53() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy53() (*big.Int, error) {
	return _Contract.Contract.Dummy53(&_Contract.CallOpts)
}

// Dummy54 is a free data retrieval call binding the contract method 0x792c7f3e.
//
// Solidity: function dummy54() pure returns(uint256)
func (_Contract *ContractCaller) Dummy54(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy54")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy54 is a free data retrieval call binding the contract method 0x792c7f3e.
//
// Solidity: function dummy54() pure returns(uint256)
func (_Contract *ContractSession) Dummy54() (*big.Int, error) {
	return _Contract.Contract.Dummy54(&_Contract.CallOpts)
}

// Dummy54 is a free data retrieval call binding the contract method 0x792c7f3e.
//
// Solidity: function dummy54() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy54() (*big.Int, error) {
	return _Contract.Contract.Dummy54(&_Contract.CallOpts)
}

// Dummy55 is a free data retrieval call binding the contract method 0x7dffdc32.
//
// Solidity: function dummy55() pure returns(uint256)
func (_Contract *ContractCaller) Dummy55(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy55")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy55 is a free data retrieval call binding the contract method 0x7dffdc32.
//
// Solidity: function dummy55() pure returns(uint256)
func (_Contract *ContractSession) Dummy55() (*big.Int, error) {
	return _Contract.Contract.Dummy55(&_Contract.CallOpts)
}

// Dummy55 is a free data retrieval call binding the contract method 0x7dffdc32.
//
// Solidity: function dummy55() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy55() (*big.Int, error) {
	return _Contract.Contract.Dummy55(&_Contract.CallOpts)
}

// Dummy56 is a free data retrieval call binding the contract method 0xa891d4d4.
//
// Solidity: function dummy56() pure returns(uint256)
func (_Contract *ContractCaller) Dummy56(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy56")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy56 is a free data retrieval call binding the contract method 0xa891d4d4.
//
// Solidity: function dummy56() pure returns(uint256)
func (_Contract *ContractSession) Dummy56() (*big.Int, error) {
	return _Contract.Contract.Dummy56(&_Contract.CallOpts)
}

// Dummy56 is a free data retrieval call binding the contract method 0xa891d4d4.
//
// Solidity: function dummy56() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy56() (*big.Int, error) {
	return _Contract.Contract.Dummy56(&_Contract.CallOpts)
}

// Dummy57 is a free data retrieval call binding the contract method 0x74f83d02.
//
// Solidity: function dummy57() pure returns(uint256)
func (_Contract *ContractCaller) Dummy57(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy57")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy57 is a free data retrieval call binding the contract method 0x74f83d02.
//
// Solidity: function dummy57() pure returns(uint256)
func (_Contract *ContractSession) Dummy57() (*big.Int, error) {
	return _Contract.Contract.Dummy57(&_Contract.CallOpts)
}

// Dummy57 is a free data retrieval call binding the contract method 0x74f83d02.
//
// Solidity: function dummy57() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy57() (*big.Int, error) {
	return _Contract.Contract.Dummy57(&_Contract.CallOpts)
}

// Dummy58 is a free data retrieval call binding the contract method 0x54c27920.
//
// Solidity: function dummy58() pure returns(uint256)
func (_Contract *ContractCaller) Dummy58(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy58")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy58 is a free data retrieval call binding the contract method 0x54c27920.
//
// Solidity: function dummy58() pure returns(uint256)
func (_Contract *ContractSession) Dummy58() (*big.Int, error) {
	return _Contract.Contract.Dummy58(&_Contract.CallOpts)
}

// Dummy58 is a free data retrieval call binding the contract method 0x54c27920.
//
// Solidity: function dummy58() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy58() (*big.Int, error) {
	return _Contract.Contract.Dummy58(&_Contract.CallOpts)
}

// Dummy59 is a free data retrieval call binding the contract method 0x77c0209e.
//
// Solidity: function dummy59() pure returns(uint256)
func (_Contract *ContractCaller) Dummy59(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy59")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy59 is a free data retrieval call binding the contract method 0x77c0209e.
//
// Solidity: function dummy59() pure returns(uint256)
func (_Contract *ContractSession) Dummy59() (*big.Int, error) {
	return _Contract.Contract.Dummy59(&_Contract.CallOpts)
}

// Dummy59 is a free data retrieval call binding the contract method 0x77c0209e.
//
// Solidity: function dummy59() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy59() (*big.Int, error) {
	return _Contract.Contract.Dummy59(&_Contract.CallOpts)
}

// Dummy6 is a free data retrieval call binding the contract method 0x4128a85d.
//
// Solidity: function dummy6() pure returns(uint256)
func (_Contract *ContractCaller) Dummy6(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy6")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy6 is a free data retrieval call binding the contract method 0x4128a85d.
//
// Solidity: function dummy6() pure returns(uint256)
func (_Contract *ContractSession) Dummy6() (*big.Int, error) {
	return _Contract.Contract.Dummy6(&_Contract.CallOpts)
}

// Dummy6 is a free data retrieval call binding the contract method 0x4128a85d.
//
// Solidity: function dummy6() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy6() (*big.Int, error) {
	return _Contract.Contract.Dummy6(&_Contract.CallOpts)
}

// Dummy60 is a free data retrieval call binding the contract method 0x34517f0b.
//
// Solidity: function dummy60() pure returns(uint256)
func (_Contract *ContractCaller) Dummy60(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy60")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy60 is a free data retrieval call binding the contract method 0x34517f0b.
//
// Solidity: function dummy60() pure returns(uint256)
func (_Contract *ContractSession) Dummy60() (*big.Int, error) {
	return _Contract.Contract.Dummy60(&_Contract.CallOpts)
}

// Dummy60 is a free data retrieval call binding the contract method 0x34517f0b.
//
// Solidity: function dummy60() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy60() (*big.Int, error) {
	return _Contract.Contract.Dummy60(&_Contract.CallOpts)
}

// Dummy61 is a free data retrieval call binding the contract method 0xe2d27530.
//
// Solidity: function dummy61() pure returns(uint256)
func (_Contract *ContractCaller) Dummy61(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy61")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy61 is a free data retrieval call binding the contract method 0xe2d27530.
//
// Solidity: function dummy61() pure returns(uint256)
func (_Contract *ContractSession) Dummy61() (*big.Int, error) {
	return _Contract.Contract.Dummy61(&_Contract.CallOpts)
}

// Dummy61 is a free data retrieval call binding the contract method 0xe2d27530.
//
// Solidity: function dummy61() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy61() (*big.Int, error) {
	return _Contract.Contract.Dummy61(&_Contract.CallOpts)
}

// Dummy62 is a free data retrieval call binding the contract method 0xf5f57381.
//
// Solidity: function dummy62() pure returns(uint256)
func (_Contract *ContractCaller) Dummy62(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy62")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy62 is a free data retrieval call binding the contract method 0xf5f57381.
//
// Solidity: function dummy62() pure returns(uint256)
func (_Contract *ContractSession) Dummy62() (*big.Int, error) {
	return _Contract.Contract.Dummy62(&_Contract.CallOpts)
}

// Dummy62 is a free data retrieval call binding the contract method 0xf5f57381.
//
// Solidity: function dummy62() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy62() (*big.Int, error) {
	return _Contract.Contract.Dummy62(&_Contract.CallOpts)
}

// Dummy63 is a free data retrieval call binding the contract method 0x7a319c18.
//
// Solidity: function dummy63() pure returns(uint256)
func (_Contract *ContractCaller) Dummy63(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy63")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy63 is a free data retrieval call binding the contract method 0x7a319c18.
//
// Solidity: function dummy63() pure returns(uint256)
func (_Contract *ContractSession) Dummy63() (*big.Int, error) {
	return _Contract.Contract.Dummy63(&_Contract.CallOpts)
}

// Dummy63 is a free data retrieval call binding the contract method 0x7a319c18.
//
// Solidity: function dummy63() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy63() (*big.Int, error) {
	return _Contract.Contract.Dummy63(&_Contract.CallOpts)
}

// Dummy64 is a free data retrieval call binding the contract method 0xad8f4221.
//
// Solidity: function dummy64() pure returns(uint256)
func (_Contract *ContractCaller) Dummy64(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy64")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy64 is a free data retrieval call binding the contract method 0xad8f4221.
//
// Solidity: function dummy64() pure returns(uint256)
func (_Contract *ContractSession) Dummy64() (*big.Int, error) {
	return _Contract.Contract.Dummy64(&_Contract.CallOpts)
}

// Dummy64 is a free data retrieval call binding the contract method 0xad8f4221.
//
// Solidity: function dummy64() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy64() (*big.Int, error) {
	return _Contract.Contract.Dummy64(&_Contract.CallOpts)
}

// Dummy65 is a free data retrieval call binding the contract method 0x1f449589.
//
// Solidity: function dummy65() pure returns(uint256)
func (_Contract *ContractCaller) Dummy65(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy65")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy65 is a free data retrieval call binding the contract method 0x1f449589.
//
// Solidity: function dummy65() pure returns(uint256)
func (_Contract *ContractSession) Dummy65() (*big.Int, error) {
	return _Contract.Contract.Dummy65(&_Contract.CallOpts)
}

// Dummy65 is a free data retrieval call binding the contract method 0x1f449589.
//
// Solidity: function dummy65() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy65() (*big.Int, error) {
	return _Contract.Contract.Dummy65(&_Contract.CallOpts)
}

// Dummy7 is a free data retrieval call binding the contract method 0xf26c779b.
//
// Solidity: function dummy7() pure returns(uint256)
func (_Contract *ContractCaller) Dummy7(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy7")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy7 is a free data retrieval call binding the contract method 0xf26c779b.
//
// Solidity: function dummy7() pure returns(uint256)
func (_Contract *ContractSession) Dummy7() (*big.Int, error) {
	return _Contract.Contract.Dummy7(&_Contract.CallOpts)
}

// Dummy7 is a free data retrieval call binding the contract method 0xf26c779b.
//
// Solidity: function dummy7() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy7() (*big.Int, error) {
	return _Contract.Contract.Dummy7(&_Contract.CallOpts)
}

// Dummy8 is a free data retrieval call binding the contract method 0x19cf6a91.
//
// Solidity: function dummy8() pure returns(uint256)
func (_Contract *ContractCaller) Dummy8(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy8")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy8 is a free data retrieval call binding the contract method 0x19cf6a91.
//
// Solidity: function dummy8() pure returns(uint256)
func (_Contract *ContractSession) Dummy8() (*big.Int, error) {
	return _Contract.Contract.Dummy8(&_Contract.CallOpts)
}

// Dummy8 is a free data retrieval call binding the contract method 0x19cf6a91.
//
// Solidity: function dummy8() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy8() (*big.Int, error) {
	return _Contract.Contract.Dummy8(&_Contract.CallOpts)
}

// Dummy9 is a free data retrieval call binding the contract method 0x58b6a9bd.
//
// Solidity: function dummy9() pure returns(uint256)
func (_Contract *ContractCaller) Dummy9(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy9")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy9 is a free data retrieval call binding the contract method 0x58b6a9bd.
//
// Solidity: function dummy9() pure returns(uint256)
func (_Contract *ContractSession) Dummy9() (*big.Int, error) {
	return _Contract.Contract.Dummy9(&_Contract.CallOpts)
}

// Dummy9 is a free data retrieval call binding the contract method 0x58b6a9bd.
//
// Solidity: function dummy9() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy9() (*big.Int, error) {
	return _Contract.Contract.Dummy9(&_Contract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Contract *ContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Contract *ContractSession) Name() (string, error) {
	return _Contract.Contract.Name(&_Contract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Contract *ContractCallerSession) Name() (string, error) {
	return _Contract.Contract.Name(&_Contract.CallOpts)
}

// Salt is a free data retrieval call binding the contract method 0xbfa0b133.
//
// Solidity: function salt() view returns(uint256)
func (_Contract *ContractCaller) Salt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "salt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Salt is a free data retrieval call binding the contract method 0xbfa0b133.
//
// Solidity: function salt() view returns(uint256)
func (_Contract *ContractSession) Salt() (*big.Int, error) {
	return _Contract.Contract.Salt(&_Contract.CallOpts)
}

// Salt is a free data retrieval call binding the contract method 0xbfa0b133.
//
// Solidity: function salt() view returns(uint256)
func (_Contract *ContractCallerSession) Salt() (*big.Int, error) {
	return _Contract.Contract.Salt(&_Contract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Contract *ContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Contract *ContractSession) Symbol() (string, error) {
	return _Contract.Contract.Symbol(&_Contract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Contract *ContractCallerSession) Symbol() (string, error) {
	return _Contract.Contract.Symbol(&_Contract.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Contract *ContractCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Contract *ContractSession) TotalSupply() (*big.Int, error) {
	return _Contract.Contract.TotalSupply(&_Contract.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Contract *ContractCallerSession) TotalSupply() (*big.Int, error) {
	return _Contract.Contract.TotalSupply(&_Contract.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Contract *ContractTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Contract *ContractSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Approve(&_Contract.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Approve(&_Contract.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Contract *ContractSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Transfer(&_Contract.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Transfer(&_Contract.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0xbb9bfe06.
//
// Solidity: function transferFrom1(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom1(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom1", from, to, value)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0xbb9bfe06.
//
// Solidity: function transferFrom1(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom1(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom1(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0xbb9bfe06.
//
// Solidity: function transferFrom1(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom1(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom1(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom10 is a paid mutator transaction binding the contract method 0xb1802b9a.
//
// Solidity: function transferFrom10(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom10(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom10", from, to, value)
}

// TransferFrom10 is a paid mutator transaction binding the contract method 0xb1802b9a.
//
// Solidity: function transferFrom10(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom10(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom10(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom10 is a paid mutator transaction binding the contract method 0xb1802b9a.
//
// Solidity: function transferFrom10(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom10(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom10(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom11 is a paid mutator transaction binding the contract method 0xc2be97e3.
//
// Solidity: function transferFrom11(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom11(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom11", from, to, value)
}

// TransferFrom11 is a paid mutator transaction binding the contract method 0xc2be97e3.
//
// Solidity: function transferFrom11(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom11(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom11(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom11 is a paid mutator transaction binding the contract method 0xc2be97e3.
//
// Solidity: function transferFrom11(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom11(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom11(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom12 is a paid mutator transaction binding the contract method 0x44050a28.
//
// Solidity: function transferFrom12(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom12(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom12", from, to, value)
}

// TransferFrom12 is a paid mutator transaction binding the contract method 0x44050a28.
//
// Solidity: function transferFrom12(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom12(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom12(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom12 is a paid mutator transaction binding the contract method 0x44050a28.
//
// Solidity: function transferFrom12(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom12(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom12(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom13 is a paid mutator transaction binding the contract method 0xacc5aee9.
//
// Solidity: function transferFrom13(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom13(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom13", from, to, value)
}

// TransferFrom13 is a paid mutator transaction binding the contract method 0xacc5aee9.
//
// Solidity: function transferFrom13(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom13(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom13(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom13 is a paid mutator transaction binding the contract method 0xacc5aee9.
//
// Solidity: function transferFrom13(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom13(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom13(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom14 is a paid mutator transaction binding the contract method 0x1bbffe6f.
//
// Solidity: function transferFrom14(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom14(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom14", from, to, value)
}

// TransferFrom14 is a paid mutator transaction binding the contract method 0x1bbffe6f.
//
// Solidity: function transferFrom14(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom14(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom14(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom14 is a paid mutator transaction binding the contract method 0x1bbffe6f.
//
// Solidity: function transferFrom14(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom14(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom14(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom15 is a paid mutator transaction binding the contract method 0x8789ca67.
//
// Solidity: function transferFrom15(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom15(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom15", from, to, value)
}

// TransferFrom15 is a paid mutator transaction binding the contract method 0x8789ca67.
//
// Solidity: function transferFrom15(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom15(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom15(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom15 is a paid mutator transaction binding the contract method 0x8789ca67.
//
// Solidity: function transferFrom15(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom15(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom15(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom16 is a paid mutator transaction binding the contract method 0x39e0bd12.
//
// Solidity: function transferFrom16(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom16(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom16", from, to, value)
}

// TransferFrom16 is a paid mutator transaction binding the contract method 0x39e0bd12.
//
// Solidity: function transferFrom16(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom16(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom16(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom16 is a paid mutator transaction binding the contract method 0x39e0bd12.
//
// Solidity: function transferFrom16(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom16(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom16(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom17 is a paid mutator transaction binding the contract method 0x291c3bd7.
//
// Solidity: function transferFrom17(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom17(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom17", from, to, value)
}

// TransferFrom17 is a paid mutator transaction binding the contract method 0x291c3bd7.
//
// Solidity: function transferFrom17(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom17(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom17(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom17 is a paid mutator transaction binding the contract method 0x291c3bd7.
//
// Solidity: function transferFrom17(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom17(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom17(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom18 is a paid mutator transaction binding the contract method 0x657b6ef7.
//
// Solidity: function transferFrom18(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom18(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom18", from, to, value)
}

// TransferFrom18 is a paid mutator transaction binding the contract method 0x657b6ef7.
//
// Solidity: function transferFrom18(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom18(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom18(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom18 is a paid mutator transaction binding the contract method 0x657b6ef7.
//
// Solidity: function transferFrom18(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom18(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom18(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom19 is a paid mutator transaction binding the contract method 0xeb4329c8.
//
// Solidity: function transferFrom19(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom19(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom19", from, to, value)
}

// TransferFrom19 is a paid mutator transaction binding the contract method 0xeb4329c8.
//
// Solidity: function transferFrom19(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom19(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom19(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom19 is a paid mutator transaction binding the contract method 0xeb4329c8.
//
// Solidity: function transferFrom19(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom19(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom19(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom2 is a paid mutator transaction binding the contract method 0x6c12ed28.
//
// Solidity: function transferFrom2(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom2(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom2", from, to, value)
}

// TransferFrom2 is a paid mutator transaction binding the contract method 0x6c12ed28.
//
// Solidity: function transferFrom2(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom2(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom2(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom2 is a paid mutator transaction binding the contract method 0x6c12ed28.
//
// Solidity: function transferFrom2(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom2(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom2(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom3 is a paid mutator transaction binding the contract method 0x1b17c65c.
//
// Solidity: function transferFrom3(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom3(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom3", from, to, value)
}

// TransferFrom3 is a paid mutator transaction binding the contract method 0x1b17c65c.
//
// Solidity: function transferFrom3(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom3(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom3(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom3 is a paid mutator transaction binding the contract method 0x1b17c65c.
//
// Solidity: function transferFrom3(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom3(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom3(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom4 is a paid mutator transaction binding the contract method 0x3a131990.
//
// Solidity: function transferFrom4(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom4(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom4", from, to, value)
}

// TransferFrom4 is a paid mutator transaction binding the contract method 0x3a131990.
//
// Solidity: function transferFrom4(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom4(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom4(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom4 is a paid mutator transaction binding the contract method 0x3a131990.
//
// Solidity: function transferFrom4(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom4(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom4(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom5 is a paid mutator transaction binding the contract method 0x0460faf6.
//
// Solidity: function transferFrom5(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom5(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom5", from, to, value)
}

// TransferFrom5 is a paid mutator transaction binding the contract method 0x0460faf6.
//
// Solidity: function transferFrom5(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom5(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom5(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom5 is a paid mutator transaction binding the contract method 0x0460faf6.
//
// Solidity: function transferFrom5(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom5(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom5(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom6 is a paid mutator transaction binding the contract method 0x7c66673e.
//
// Solidity: function transferFrom6(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom6(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom6", from, to, value)
}

// TransferFrom6 is a paid mutator transaction binding the contract method 0x7c66673e.
//
// Solidity: function transferFrom6(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom6(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom6(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom6 is a paid mutator transaction binding the contract method 0x7c66673e.
//
// Solidity: function transferFrom6(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom6(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom6(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom7 is a paid mutator transaction binding the contract method 0xfaf35ced.
//
// Solidity: function transferFrom7(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom7(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom7", from, to, value)
}

// TransferFrom7 is a paid mutator transaction binding the contract method 0xfaf35ced.
//
// Solidity: function transferFrom7(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom7(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom7(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom7 is a paid mutator transaction binding the contract method 0xfaf35ced.
//
// Solidity: function transferFrom7(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom7(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom7(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom8 is a paid mutator transaction binding the contract method 0x552a1b56.
//
// Solidity: function transferFrom8(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom8(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom8", from, to, value)
}

// TransferFrom8 is a paid mutator transaction binding the contract method 0x552a1b56.
//
// Solidity: function transferFrom8(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom8(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom8(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom8 is a paid mutator transaction binding the contract method 0x552a1b56.
//
// Solidity: function transferFrom8(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom8(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom8(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom9 is a paid mutator transaction binding the contract method 0x4f7bd75a.
//
// Solidity: function transferFrom9(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactor) TransferFrom9(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferFrom9", from, to, value)
}

// TransferFrom9 is a paid mutator transaction binding the contract method 0x4f7bd75a.
//
// Solidity: function transferFrom9(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractSession) TransferFrom9(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom9(&_Contract.TransactOpts, from, to, value)
}

// TransferFrom9 is a paid mutator transaction binding the contract method 0x4f7bd75a.
//
// Solidity: function transferFrom9(address from, address to, uint256 value) returns(bool)
func (_Contract *ContractTransactorSession) TransferFrom9(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.TransferFrom9(&_Contract.TransactOpts, from, to, value)
}

// ContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Contract contract.
type ContractApprovalIterator struct {
	Event *ContractApproval // Event containing the contract specifics and raw log

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
func (it *ContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractApproval)
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
		it.Event = new(ContractApproval)
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
func (it *ContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractApproval represents a Approval event raised by the Contract contract.
type ContractApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Contract *ContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ContractApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ContractApprovalIterator{contract: _Contract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Contract *ContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ContractApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractApproval)
				if err := _Contract.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Contract *ContractFilterer) ParseApproval(log types.Log) (*ContractApproval, error) {
	event := new(ContractApproval)
	if err := _Contract.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Contract contract.
type ContractTransferIterator struct {
	Event *ContractTransfer // Event containing the contract specifics and raw log

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
func (it *ContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractTransfer)
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
		it.Event = new(ContractTransfer)
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
func (it *ContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractTransfer represents a Transfer event raised by the Contract contract.
type ContractTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Contract *ContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ContractTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ContractTransferIterator{contract: _Contract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Contract *ContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ContractTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractTransfer)
				if err := _Contract.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Contract *ContractFilterer) ParseTransfer(log types.Log) (*ContractTransfer, error) {
	event := new(ContractTransfer)
	if err := _Contract.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
