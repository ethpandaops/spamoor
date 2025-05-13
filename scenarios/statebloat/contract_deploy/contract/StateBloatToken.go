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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy10\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy11\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy12\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy13\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy14\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy15\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy16\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy17\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy18\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy19\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy20\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy21\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy22\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy23\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy24\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy25\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy26\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy27\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy28\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy29\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy3\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy30\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy31\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy32\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy33\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy34\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy35\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy36\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy37\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy38\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy39\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy4\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy40\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy41\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy42\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy43\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy44\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy45\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy46\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy47\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy48\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy49\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy5\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy50\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy51\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy52\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy53\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy54\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy55\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy56\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy57\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy58\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy59\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy6\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy60\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy61\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy62\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy63\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy64\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy65\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy66\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy67\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy68\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy69\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy7\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy70\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy71\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy72\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy73\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy74\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy75\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy8\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy9\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"salt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom1\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom10\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom11\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom12\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom13\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom14\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom15\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom16\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom17\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom18\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom19\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom3\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom4\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom5\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom6\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom7\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom8\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom9\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051615fd4380380615fd4833981810160405281019061003191906101f4565b6040518060400160405280601181526020017f537461746520426c6f617420546f6b656e0000000000000000000000000000008152505f90816100749190610453565b506040518060400160405280600381526020017f5342540000000000000000000000000000000000000000000000000000000000815250600190816100b99190610453565b50601260025f6101000a81548160ff021916908360ff160217905550806080818152505060025f9054906101000a900460ff16600a6100f8919061068a565b620f424061010691906106d4565b60038190555060035460045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055503373ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6003546040516101af9190610724565b60405180910390a35061073d565b5f5ffd5b5f819050919050565b6101d3816101c1565b81146101dd575f5ffd5b50565b5f815190506101ee816101ca565b92915050565b5f60208284031215610209576102086101bd565b5b5f610216848285016101e0565b91505092915050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061029a57607f821691505b6020821081036102ad576102ac610256565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261030f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826102d4565b61031986836102d4565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61035461034f61034a846101c1565b610331565b6101c1565b9050919050565b5f819050919050565b61036d8361033a565b6103816103798261035b565b8484546102e0565b825550505050565b5f5f905090565b610398610389565b6103a3818484610364565b505050565b5b818110156103c6576103bb5f82610390565b6001810190506103a9565b5050565b601f82111561040b576103dc816102b3565b6103e5846102c5565b810160208510156103f4578190505b610408610400856102c5565b8301826103a8565b50505b505050565b5f82821c905092915050565b5f61042b5f1984600802610410565b1980831691505092915050565b5f610443838361041c565b9150826002028217905092915050565b61045c8261021f565b67ffffffffffffffff81111561047557610474610229565b5b61047f8254610283565b61048a8282856103ca565b5f60209050601f8311600181146104bb575f84156104a9578287015190505b6104b38582610438565b86555061051a565b601f1984166104c9866102b3565b5f5b828110156104f0578489015182556001820191506020850194506020810190506104cb565b8683101561050d5784890151610509601f89168261041c565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f8160011c9050919050565b5f5f8291508390505b60018511156105a4578086048111156105805761057f610522565b5b600185161561058f5780820291505b808102905061059d8561054f565b9450610564565b94509492505050565b5f826105bc5760019050610677565b816105c9575f9050610677565b81600181146105df57600281146105e957610618565b6001915050610677565b60ff8411156105fb576105fa610522565b5b8360020a91508482111561061257610611610522565b5b50610677565b5060208310610133831016604e8410600b841016171561064d5782820a90508381111561064857610647610522565b5b610677565b61065a848484600161055b565b9250905081840481111561067157610670610522565b5b81810290505b9392505050565b5f60ff82169050919050565b5f610694826101c1565b915061069f8361067e565b92506106cc7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846105ad565b905092915050565b5f6106de826101c1565b91506106e9836101c1565b92508282026106f7816101c1565b9150828204841483151761070e5761070d610522565b5b5092915050565b61071e816101c1565b82525050565b5f6020820190506107375f830184610715565b92915050565b60805161587f6107555f395f614a7b015261587f5ff3fe608060405234801561000f575f5ffd5b5060043610610606575f3560e01c8063657b6ef711610319578063b1802b9a116101a6578063d9bb3174116100f2578063f26c779b116100ab578063f8716f1411610085578063f8716f1414611360578063faf35ced1461137e578063fe7d5996146113ae578063ffbf0469146113cc57610606565b8063f26c779b14611306578063f4978b0414611324578063f5f573811461134257610606565b8063d9bb31741461122e578063dc1d8a9b1461124c578063dd62ed3e1461126a578063e2d275301461129a578063e8c927b3146112b8578063eb4329c8146112d657610606565b8063bbe231321161015f578063c958d4bf11610139578063c958d4bf146111b6578063cfd66863146111d4578063d101dcd0146111f2578063d74194691461121057610606565b8063bbe231321461114a578063bfa0b13314611168578063c2be97e31461118657610606565b8063b1802b9a14611072578063b2bb360e146110a2578063b4c48668146110c0578063b66dd750146110de578063b9a6d645146110fc578063bb9bfe061461111a57610606565b80637e544937116102655780639c5dfe731161021e578063a9059cbb116101f8578063a9059cbb14610fd6578063aaa7af7014611006578063acc5aee914611024578063ad8f42211461105457610606565b80639c5dfe7314610f7c5780639df61a2514610f9a578063a891d4d414610fb857610606565b80637e54493714610eb65780637f34d94b14610ed45780638619d60714610ef25780638789ca6714610f105780638f4a840614610f4057806395d89b4114610f5e57610606565b806374f83d02116102d25780637a319c18116102ac5780637a319c1814610e2c5780637c66673e14610e4a5780637c72ed0d14610e7a5780637dffdc3214610e9857610606565b806374f83d0214610dd257806377c0209e14610df0578063792c7f3e14610e0e57610606565b8063657b6ef714610ce8578063672151fe14610d185780636abceacd14610d365780636c12ed2814610d5457806370a0823114610d8457806374e73fd314610db457610606565b8063313ce567116104975780634b3c7f5f116103e357806358b6a9bd1161039c57806361b970eb1161037657806361b970eb14610c70578063639ec53a14610c8e5780636547317414610cac5780636578534c14610cca57610606565b806358b6a9bd14610c16578063595471b814610c345780635af92c0514610c5257610606565b80634b3c7f5f14610b3e5780634e1dbb8214610b5c5780634f5e555714610b7a5780634f7bd75a14610b9857806354c2792014610bc8578063552a1b5614610be657610606565b80633ea117ce1161045057806342937dbd1161042a57806342937dbd14610ab457806343a6b92d14610ad257806344050a2814610af05780634a2e93c614610b2057610606565b80633ea117ce14610a5a5780634128a85d14610a78578063418b181614610a9657610606565b8063313ce5671461098257806334517f0b146109a0578063378c5382146109be57806339e0bd12146109dc5780633a13199014610a0c5780633b6be45914610a3c57610606565b80631bbffe6f1161055657806321ecd7a31161050f5780632545d8b7116104e95780632545d8b7146108f85780632787325b14610916578063291c3bd7146109345780633125f37a1461096457610606565b806321ecd7a31461088c578063239af2a5146108aa57806323b872dd146108c857610606565b80631bbffe6f146107c65780631d527cde146107f65780631eaa7c52146108145780631f29783f146108325780631f449589146108505780631fd298ec1461086e57610606565b806312901b42116105c357806318160ddd1161059d57806318160ddd1461073c57806319cf6a911461075a5780631a97f18e146107785780631b17c65c1461079657610606565b806312901b42146106e257806313ebb5ec1461070057806316a3045b1461071e57610606565b80630460faf61461060a57806306fdde031461063a5780630717b16114610658578063095ea7b3146106765780630cb7a9e7146106a65780631215a3ab146106c4575b5f5ffd5b610624600480360381019061061f91906153a3565b6113ea565b604051610631919061540d565b60405180910390f35b6106426116ca565b60405161064f9190615496565b60405180910390f35b610660611755565b60405161066d91906154c5565b60405180910390f35b610690600480360381019061068b91906154de565b61175d565b60405161069d919061540d565b60405180910390f35b6106ae61184a565b6040516106bb91906154c5565b60405180910390f35b6106cc611852565b6040516106d991906154c5565b60405180910390f35b6106ea61185a565b6040516106f791906154c5565b60405180910390f35b610708611862565b60405161071591906154c5565b60405180910390f35b61072661186a565b60405161073391906154c5565b60405180910390f35b610744611872565b60405161075191906154c5565b60405180910390f35b610762611878565b60405161076f91906154c5565b60405180910390f35b610780611880565b60405161078d91906154c5565b60405180910390f35b6107b060048036038101906107ab91906153a3565b611888565b6040516107bd919061540d565b60405180910390f35b6107e060048036038101906107db91906153a3565b611b68565b6040516107ed919061540d565b60405180910390f35b6107fe611e48565b60405161080b91906154c5565b60405180910390f35b61081c611e50565b60405161082991906154c5565b60405180910390f35b61083a611e58565b60405161084791906154c5565b60405180910390f35b610858611e60565b60405161086591906154c5565b60405180910390f35b610876611e68565b60405161088391906154c5565b60405180910390f35b610894611e70565b6040516108a191906154c5565b60405180910390f35b6108b2611e78565b6040516108bf91906154c5565b60405180910390f35b6108e260048036038101906108dd91906153a3565b611e80565b6040516108ef919061540d565b60405180910390f35b610900612160565b60405161090d91906154c5565b60405180910390f35b61091e612168565b60405161092b91906154c5565b60405180910390f35b61094e600480360381019061094991906153a3565b612170565b60405161095b919061540d565b60405180910390f35b61096c612450565b60405161097991906154c5565b60405180910390f35b61098a612458565b6040516109979190615537565b60405180910390f35b6109a861246a565b6040516109b591906154c5565b60405180910390f35b6109c6612472565b6040516109d391906154c5565b60405180910390f35b6109f660048036038101906109f191906153a3565b61247a565b604051610a03919061540d565b60405180910390f35b610a266004803603810190610a2191906153a3565b61275a565b604051610a33919061540d565b60405180910390f35b610a44612a3a565b604051610a5191906154c5565b60405180910390f35b610a62612a42565b604051610a6f91906154c5565b60405180910390f35b610a80612a4a565b604051610a8d91906154c5565b60405180910390f35b610a9e612a52565b604051610aab91906154c5565b60405180910390f35b610abc612a5a565b604051610ac991906154c5565b60405180910390f35b610ada612a62565b604051610ae791906154c5565b60405180910390f35b610b0a6004803603810190610b0591906153a3565b612a6a565b604051610b17919061540d565b60405180910390f35b610b28612d4a565b604051610b3591906154c5565b60405180910390f35b610b46612d52565b604051610b5391906154c5565b60405180910390f35b610b64612d5a565b604051610b7191906154c5565b60405180910390f35b610b82612d62565b604051610b8f91906154c5565b60405180910390f35b610bb26004803603810190610bad91906153a3565b612d6a565b604051610bbf919061540d565b60405180910390f35b610bd061304a565b604051610bdd91906154c5565b60405180910390f35b610c006004803603810190610bfb91906153a3565b613052565b604051610c0d919061540d565b60405180910390f35b610c1e613332565b604051610c2b91906154c5565b60405180910390f35b610c3c61333a565b604051610c4991906154c5565b60405180910390f35b610c5a613342565b604051610c6791906154c5565b60405180910390f35b610c7861334a565b604051610c8591906154c5565b60405180910390f35b610c96613352565b604051610ca391906154c5565b60405180910390f35b610cb461335a565b604051610cc191906154c5565b60405180910390f35b610cd2613362565b604051610cdf91906154c5565b60405180910390f35b610d026004803603810190610cfd91906153a3565b61336a565b604051610d0f919061540d565b60405180910390f35b610d2061364a565b604051610d2d91906154c5565b60405180910390f35b610d3e613652565b604051610d4b91906154c5565b60405180910390f35b610d6e6004803603810190610d6991906153a3565b61365a565b604051610d7b919061540d565b60405180910390f35b610d9e6004803603810190610d999190615550565b61393a565b604051610dab91906154c5565b60405180910390f35b610dbc61394f565b604051610dc991906154c5565b60405180910390f35b610dda613957565b604051610de791906154c5565b60405180910390f35b610df861395f565b604051610e0591906154c5565b60405180910390f35b610e16613967565b604051610e2391906154c5565b60405180910390f35b610e3461396f565b604051610e4191906154c5565b60405180910390f35b610e646004803603810190610e5f91906153a3565b613977565b604051610e71919061540d565b60405180910390f35b610e82613c57565b604051610e8f91906154c5565b60405180910390f35b610ea0613c5f565b604051610ead91906154c5565b60405180910390f35b610ebe613c67565b604051610ecb91906154c5565b60405180910390f35b610edc613c6f565b604051610ee991906154c5565b60405180910390f35b610efa613c77565b604051610f0791906154c5565b60405180910390f35b610f2a6004803603810190610f2591906153a3565b613c7f565b604051610f37919061540d565b60405180910390f35b610f48613f5f565b604051610f5591906154c5565b60405180910390f35b610f66613f67565b604051610f739190615496565b60405180910390f35b610f84613ff3565b604051610f9191906154c5565b60405180910390f35b610fa2613ffb565b604051610faf91906154c5565b60405180910390f35b610fc0614003565b604051610fcd91906154c5565b60405180910390f35b610ff06004803603810190610feb91906154de565b61400b565b604051610ffd919061540d565b60405180910390f35b61100e6141a1565b60405161101b91906154c5565b60405180910390f35b61103e600480360381019061103991906153a3565b6141a9565b60405161104b919061540d565b60405180910390f35b61105c614489565b60405161106991906154c5565b60405180910390f35b61108c600480360381019061108791906153a3565b614491565b604051611099919061540d565b60405180910390f35b6110aa614771565b6040516110b791906154c5565b60405180910390f35b6110c8614779565b6040516110d591906154c5565b60405180910390f35b6110e6614781565b6040516110f391906154c5565b60405180910390f35b611104614789565b60405161111191906154c5565b60405180910390f35b611134600480360381019061112f91906153a3565b614791565b604051611141919061540d565b60405180910390f35b611152614a71565b60405161115f91906154c5565b60405180910390f35b611170614a79565b60405161117d91906154c5565b60405180910390f35b6111a0600480360381019061119b91906153a3565b614a9d565b6040516111ad919061540d565b60405180910390f35b6111be614d7d565b6040516111cb91906154c5565b60405180910390f35b6111dc614d85565b6040516111e991906154c5565b60405180910390f35b6111fa614d8d565b60405161120791906154c5565b60405180910390f35b611218614d95565b60405161122591906154c5565b60405180910390f35b611236614d9d565b60405161124391906154c5565b60405180910390f35b611254614da5565b60405161126191906154c5565b60405180910390f35b611284600480360381019061127f919061557b565b614dad565b60405161129191906154c5565b60405180910390f35b6112a2614dcd565b6040516112af91906154c5565b60405180910390f35b6112c0614dd5565b6040516112cd91906154c5565b60405180910390f35b6112f060048036038101906112eb91906153a3565b614ddd565b6040516112fd919061540d565b60405180910390f35b61130e615002565b60405161131b91906154c5565b60405180910390f35b61132c61500a565b60405161133991906154c5565b60405180910390f35b61134a615012565b60405161135791906154c5565b60405180910390f35b61136861501a565b60405161137591906154c5565b60405180910390f35b611398600480360381019061139391906153a3565b615022565b6040516113a5919061540d565b60405180910390f35b6113b6615302565b6040516113c391906154c5565b60405180910390f35b6113d461530a565b6040516113e191906154c5565b60405180910390f35b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561146b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161146290615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611526576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161151d9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461157291906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546115c591906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461165391906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516116b791906154c5565b60405180910390a3600190509392505050565b5f80546116d690615749565b80601f016020809104026020016040519081016040528092919081815260200182805461170290615749565b801561174d5780601f106117245761010080835404028352916020019161174d565b820191905f5260205f20905b81548152906001019060200180831161173057829003601f168201915b505050505081565b5f602f905090565b5f8160055f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161183891906154c5565b60405180910390a36001905092915050565b5f601b905090565b5f6005905090565b5f601a905090565b5f601d905090565b5f6033905090565b60035481565b5f6008905090565b5f6003905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611909576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161190090615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156119c4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119bb9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611a1091906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611a6391906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611af191906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051611b5591906154c5565b60405180910390a3600190509392505050565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611be9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611be090615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611ca4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611c9b9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611cf091906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611d4391906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254611dd191906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051611e3591906154c5565b60405180910390a3600190509392505050565b5f6002905090565b5f6010905090565b5f6048905090565b5f6041905090565b5f6035905090565b5f602e905090565b5f602c905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611f01576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611ef890615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015611fbc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611fb39061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461200891906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461205b91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546120e991906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161214d91906154c5565b60405180910390a3600190509392505050565b5f6001905090565b5f6045905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156121f1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016121e8906157c3565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156122ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016122a39061582b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546122f891906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461234b91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546123d991906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161243d91906154c5565b60405180910390a3600190509392505050565b5f6034905090565b60025f9054906101000a900460ff1681565b5f603c905090565b5f604b905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156124fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016124f290615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156125b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016125ad9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461260291906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461265591906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546126e391906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161274791906154c5565b60405180910390a3600190509392505050565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156127db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016127d290615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612896576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161288d9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546128e291906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461293591906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546129c391906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612a2791906154c5565b60405180910390a3600190509392505050565b5f6004905090565b5f600c905090565b5f6006905090565b5f601c905090565b5f6032905090565b5f6023905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612aeb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612ae290615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612ba6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612b9d9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612bf291906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612c4591906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612cd391906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051612d3791906154c5565b60405180910390a3600190509392505050565b5f6011905090565b5f6021905090565b5f601f905090565b5f6026905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612deb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612de290615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015612ea6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612e9d9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612ef291906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612f4591906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254612fd391906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161303791906154c5565b60405180910390a3600190509392505050565b5f603a905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156130d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016130ca90615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561318e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016131859061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546131da91906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461322d91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546132bb91906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161331f91906154c5565b60405180910390a3600190509392505050565b5f6009905090565b5f6049905090565b5f602a905090565b5f6014905090565b5f6025905090565b5f6013905090565b5f6043905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156133eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016133e2906157c3565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156134a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161349d9061582b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546134f291906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461354591906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546135d391906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161363791906154c5565b60405180910390a3600190509392505050565b5f600d905090565b5f6017905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156136db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016136d290615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613796576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161378d9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546137e291906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461383591906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546138c391906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161392791906154c5565b60405180910390a3600190509392505050565b6004602052805f5260405f205f915090505481565b5f6016905090565b5f6039905090565b5f603b905090565b5f6036905090565b5f603f905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156139f8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016139ef90615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613ab3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613aaa9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613aff91906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613b5291906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613be091906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051613c4491906154c5565b60405180910390a3600190509392505050565b5f602b905090565b5f6037905090565b5f600a905090565b5f602d905090565b5f6030905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613d00576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613cf790615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015613dbb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613db29061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613e0791906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613e5a91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254613ee891906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051613f4c91906154c5565b60405180910390a3600190509392505050565b5f600b905090565b60018054613f7490615749565b80601f0160208091040260200160405190810160405280929190818152602001828054613fa090615749565b8015613feb5780601f10613fc257610100808354040283529160200191613feb565b820191905f5260205f20905b815481529060010190602001808311613fce57829003601f168201915b505050505081565b5f6015905090565b5f6029905090565b5f6038905090565b5f8160045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561408c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161408390615603565b60405180910390fd5b8160045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546140d891906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461412b91906156e9565b925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161418f91906154c5565b60405180910390a36001905092915050565b5f6028905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561422a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161422190615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156142e5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016142dc9061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461433191906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461438491906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461441291906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161447691906154c5565b60405180910390a3600190509392505050565b5f6040905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614512576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161450990615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156145cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016145c49061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461461991906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461466c91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546146fa91906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161475e91906154c5565b60405180910390a3600190509392505050565b5f6018905090565b5f604a905090565b5f600f905090565b5f6044905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614812576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161480990615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156148cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016148c49061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461491991906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461496c91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546149fa91906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051614a5e91906154c5565b60405180910390a3600190509392505050565b5f6046905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614b1e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614b1590615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614bd9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614bd09061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614c2591906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614c7891906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614d0691906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051614d6a91906154c5565b60405180910390a3600190509392505050565b5f6024905090565b5f601e905090565b5f6031905090565b5f6019905090565b5f6042905090565b5f6027905090565b6005602052815f5260405f20602052805f5260405f205f91509150505481565b5f603d905090565b5f6022905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541015614e5e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614e55906157c3565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614eaa91906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614efd91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f828254614f8b91906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051614fef91906154c5565b60405180910390a3600190509392505050565b5f6007905090565b5f6047905090565b5f603e905090565b5f6020905090565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156150a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161509a90615603565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561515e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016151559061566b565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546151aa91906156b6565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546151fd91906156e9565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461528b91906156b6565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516152ef91906154c5565b60405180910390a3600190509392505050565b5f600e905090565b5f6012905090565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61533f82615316565b9050919050565b61534f81615335565b8114615359575f5ffd5b50565b5f8135905061536a81615346565b92915050565b5f819050919050565b61538281615370565b811461538c575f5ffd5b50565b5f8135905061539d81615379565b92915050565b5f5f5f606084860312156153ba576153b9615312565b5b5f6153c78682870161535c565b93505060206153d88682870161535c565b92505060406153e98682870161538f565b9150509250925092565b5f8115159050919050565b615407816153f3565b82525050565b5f6020820190506154205f8301846153fe565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61546882615426565b6154728185615430565b9350615482818560208601615440565b61548b8161544e565b840191505092915050565b5f6020820190508181035f8301526154ae818461545e565b905092915050565b6154bf81615370565b82525050565b5f6020820190506154d85f8301846154b6565b92915050565b5f5f604083850312156154f4576154f3615312565b5b5f6155018582860161535c565b92505060206155128582860161538f565b9150509250929050565b5f60ff82169050919050565b6155318161551c565b82525050565b5f60208201905061554a5f830184615528565b92915050565b5f6020828403121561556557615564615312565b5b5f6155728482850161535c565b91505092915050565b5f5f6040838503121561559157615590615312565b5b5f61559e8582860161535c565b92505060206155af8582860161535c565b9150509250929050565b7f496e73756666696369656e742062616c616e63650000000000000000000000005f82015250565b5f6155ed601483615430565b91506155f8826155b9565b602082019050919050565b5f6020820190508181035f83015261561a816155e1565b9050919050565b7f496e73756666696369656e7420616c6c6f77616e6365000000000000000000005f82015250565b5f615655601683615430565b915061566082615621565b602082019050919050565b5f6020820190508181035f83015261568281615649565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6156c082615370565b91506156cb83615370565b92508282039050818111156156e3576156e2615689565b5b92915050565b5f6156f382615370565b91506156fe83615370565b925082820190508082111561571657615715615689565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061576057607f821691505b6020821081036157735761577261571c565b5b50919050565b7f41000000000000000000000000000000000000000000000000000000000000005f82015250565b5f6157ad600183615430565b91506157b882615779565b602082019050919050565b5f6020820190508181035f8301526157da816157a1565b9050919050565b7f42000000000000000000000000000000000000000000000000000000000000005f82015250565b5f615815600183615430565b9150615820826157e1565b602082019050919050565b5f6020820190508181035f83015261584281615809565b905091905056fea2646970667358221220ce03f2ea0f8ec341a07b9c0512b41f44d879b740c55b3572c5707216663e8a3f64736f6c634300081e0033",
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

// Dummy66 is a free data retrieval call binding the contract method 0xd9bb3174.
//
// Solidity: function dummy66() pure returns(uint256)
func (_Contract *ContractCaller) Dummy66(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy66")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy66 is a free data retrieval call binding the contract method 0xd9bb3174.
//
// Solidity: function dummy66() pure returns(uint256)
func (_Contract *ContractSession) Dummy66() (*big.Int, error) {
	return _Contract.Contract.Dummy66(&_Contract.CallOpts)
}

// Dummy66 is a free data retrieval call binding the contract method 0xd9bb3174.
//
// Solidity: function dummy66() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy66() (*big.Int, error) {
	return _Contract.Contract.Dummy66(&_Contract.CallOpts)
}

// Dummy67 is a free data retrieval call binding the contract method 0x6578534c.
//
// Solidity: function dummy67() pure returns(uint256)
func (_Contract *ContractCaller) Dummy67(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy67")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy67 is a free data retrieval call binding the contract method 0x6578534c.
//
// Solidity: function dummy67() pure returns(uint256)
func (_Contract *ContractSession) Dummy67() (*big.Int, error) {
	return _Contract.Contract.Dummy67(&_Contract.CallOpts)
}

// Dummy67 is a free data retrieval call binding the contract method 0x6578534c.
//
// Solidity: function dummy67() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy67() (*big.Int, error) {
	return _Contract.Contract.Dummy67(&_Contract.CallOpts)
}

// Dummy68 is a free data retrieval call binding the contract method 0xb9a6d645.
//
// Solidity: function dummy68() pure returns(uint256)
func (_Contract *ContractCaller) Dummy68(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy68")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy68 is a free data retrieval call binding the contract method 0xb9a6d645.
//
// Solidity: function dummy68() pure returns(uint256)
func (_Contract *ContractSession) Dummy68() (*big.Int, error) {
	return _Contract.Contract.Dummy68(&_Contract.CallOpts)
}

// Dummy68 is a free data retrieval call binding the contract method 0xb9a6d645.
//
// Solidity: function dummy68() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy68() (*big.Int, error) {
	return _Contract.Contract.Dummy68(&_Contract.CallOpts)
}

// Dummy69 is a free data retrieval call binding the contract method 0x2787325b.
//
// Solidity: function dummy69() pure returns(uint256)
func (_Contract *ContractCaller) Dummy69(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy69")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy69 is a free data retrieval call binding the contract method 0x2787325b.
//
// Solidity: function dummy69() pure returns(uint256)
func (_Contract *ContractSession) Dummy69() (*big.Int, error) {
	return _Contract.Contract.Dummy69(&_Contract.CallOpts)
}

// Dummy69 is a free data retrieval call binding the contract method 0x2787325b.
//
// Solidity: function dummy69() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy69() (*big.Int, error) {
	return _Contract.Contract.Dummy69(&_Contract.CallOpts)
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

// Dummy70 is a free data retrieval call binding the contract method 0xbbe23132.
//
// Solidity: function dummy70() pure returns(uint256)
func (_Contract *ContractCaller) Dummy70(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy70")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy70 is a free data retrieval call binding the contract method 0xbbe23132.
//
// Solidity: function dummy70() pure returns(uint256)
func (_Contract *ContractSession) Dummy70() (*big.Int, error) {
	return _Contract.Contract.Dummy70(&_Contract.CallOpts)
}

// Dummy70 is a free data retrieval call binding the contract method 0xbbe23132.
//
// Solidity: function dummy70() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy70() (*big.Int, error) {
	return _Contract.Contract.Dummy70(&_Contract.CallOpts)
}

// Dummy71 is a free data retrieval call binding the contract method 0xf4978b04.
//
// Solidity: function dummy71() pure returns(uint256)
func (_Contract *ContractCaller) Dummy71(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy71")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy71 is a free data retrieval call binding the contract method 0xf4978b04.
//
// Solidity: function dummy71() pure returns(uint256)
func (_Contract *ContractSession) Dummy71() (*big.Int, error) {
	return _Contract.Contract.Dummy71(&_Contract.CallOpts)
}

// Dummy71 is a free data retrieval call binding the contract method 0xf4978b04.
//
// Solidity: function dummy71() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy71() (*big.Int, error) {
	return _Contract.Contract.Dummy71(&_Contract.CallOpts)
}

// Dummy72 is a free data retrieval call binding the contract method 0x1f29783f.
//
// Solidity: function dummy72() pure returns(uint256)
func (_Contract *ContractCaller) Dummy72(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy72")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy72 is a free data retrieval call binding the contract method 0x1f29783f.
//
// Solidity: function dummy72() pure returns(uint256)
func (_Contract *ContractSession) Dummy72() (*big.Int, error) {
	return _Contract.Contract.Dummy72(&_Contract.CallOpts)
}

// Dummy72 is a free data retrieval call binding the contract method 0x1f29783f.
//
// Solidity: function dummy72() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy72() (*big.Int, error) {
	return _Contract.Contract.Dummy72(&_Contract.CallOpts)
}

// Dummy73 is a free data retrieval call binding the contract method 0x595471b8.
//
// Solidity: function dummy73() pure returns(uint256)
func (_Contract *ContractCaller) Dummy73(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy73")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy73 is a free data retrieval call binding the contract method 0x595471b8.
//
// Solidity: function dummy73() pure returns(uint256)
func (_Contract *ContractSession) Dummy73() (*big.Int, error) {
	return _Contract.Contract.Dummy73(&_Contract.CallOpts)
}

// Dummy73 is a free data retrieval call binding the contract method 0x595471b8.
//
// Solidity: function dummy73() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy73() (*big.Int, error) {
	return _Contract.Contract.Dummy73(&_Contract.CallOpts)
}

// Dummy74 is a free data retrieval call binding the contract method 0xb4c48668.
//
// Solidity: function dummy74() pure returns(uint256)
func (_Contract *ContractCaller) Dummy74(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy74")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy74 is a free data retrieval call binding the contract method 0xb4c48668.
//
// Solidity: function dummy74() pure returns(uint256)
func (_Contract *ContractSession) Dummy74() (*big.Int, error) {
	return _Contract.Contract.Dummy74(&_Contract.CallOpts)
}

// Dummy74 is a free data retrieval call binding the contract method 0xb4c48668.
//
// Solidity: function dummy74() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy74() (*big.Int, error) {
	return _Contract.Contract.Dummy74(&_Contract.CallOpts)
}

// Dummy75 is a free data retrieval call binding the contract method 0x378c5382.
//
// Solidity: function dummy75() pure returns(uint256)
func (_Contract *ContractCaller) Dummy75(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "dummy75")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy75 is a free data retrieval call binding the contract method 0x378c5382.
//
// Solidity: function dummy75() pure returns(uint256)
func (_Contract *ContractSession) Dummy75() (*big.Int, error) {
	return _Contract.Contract.Dummy75(&_Contract.CallOpts)
}

// Dummy75 is a free data retrieval call binding the contract method 0x378c5382.
//
// Solidity: function dummy75() pure returns(uint256)
func (_Contract *ContractCallerSession) Dummy75() (*big.Int, error) {
	return _Contract.Contract.Dummy75(&_Contract.CallOpts)
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
