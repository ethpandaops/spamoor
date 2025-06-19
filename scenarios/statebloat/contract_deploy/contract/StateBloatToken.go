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

// StateBloatTokenMetaData contains all meta data concerning the StateBloatToken contract.
var StateBloatTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PADDING_DATA\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy10\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy11\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy12\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy13\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy14\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy15\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy16\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy17\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy18\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy19\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy20\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy21\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy22\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy23\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy24\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy25\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy26\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy27\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy28\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy29\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy3\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy30\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy31\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy32\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy33\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy34\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy35\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy36\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy37\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy38\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy39\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy4\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy40\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy41\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy42\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy43\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy44\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy45\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy46\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy47\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy48\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy49\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy5\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy50\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy51\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy52\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy53\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy54\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy55\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy56\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy57\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy58\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy59\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy6\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy60\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy61\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy62\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy63\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy64\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy65\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy66\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy67\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy68\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy69\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy7\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy8\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dummy9\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"salt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom1\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom10\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom11\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom12\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom13\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom14\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom15\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom16\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom17\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom18\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom19\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom3\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom4\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom5\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom6\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom7\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom8\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom9\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801562000010575f80fd5b50604051620016f1380380620016f1833981016040819052620000339162000118565b60408051808201909152601181527029ba30ba3290213637b0ba102a37b5b2b760791b60208201525f90620000699082620001ce565b5060408051808201909152600381526214d09560ea1b6020820152600190620000939082620001ce565b506002805460ff191660129081179091556080829052620000b690600a620003a9565b620000c590620f4240620003c0565b6003819055335f81815260046020908152604080832085905551938452919290917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a350620003da565b5f6020828403121562000129575f80fd5b5051919050565b634e487b7160e01b5f52604160045260245ffd5b600181811c908216806200015957607f821691505b6020821081036200017857634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620001c957805f5260205f20601f840160051c81016020851015620001a55750805b601f840160051c820191505b81811015620001c6575f8155600101620001b1565b50505b505050565b81516001600160401b03811115620001ea57620001ea62000130565b6200020281620001fb845462000144565b846200017e565b602080601f83116001811462000238575f8415620002205750858301515b5f19600386901b1c1916600185901b17855562000292565b5f85815260208120601f198616915b82811015620002685788860151825594840194600190910190840162000247565b50858210156200028657878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b634e487b7160e01b5f52601160045260245ffd5b600181815b80851115620002ee57815f1904821115620002d257620002d26200029a565b80851615620002e057918102915b93841c9390800290620002b3565b509250929050565b5f826200030657506001620003a3565b816200031457505f620003a3565b81600181146200032d5760028114620003385762000358565b6001915050620003a3565b60ff8411156200034c576200034c6200029a565b50506001821b620003a3565b5060208310610133831016604e8410600b84101617156200037d575081810a620003a3565b620003898383620002ae565b805f19048211156200039f576200039f6200029a565b0290505b92915050565b5f620003b960ff841683620002f6565b9392505050565b8082028115828204841417620003a357620003a36200029a565b6080516112fe620003f35f395f61089001526112fe5ff3fe608060405234801561000f575f80fd5b5060043610610630575f3560e01c8063657b6ef711610333578063ad8f4221116101b3578063d9bb3174116100fe578063ee1682b6116100a9578063f8716f1411610084578063f8716f141461093d578063faf35ced14610634578063fe7d599614610944578063ffbf04691461094b575f80fd5b8063ee1682b614610927578063f26c779b1461092f578063f5f5738114610936575f80fd5b8063e2d27530116100d9578063e2d2753014610906578063e8c927b31461090d578063eb4329c814610914575f80fd5b8063d9bb3174146108ce578063dc1d8a9b146108d5578063dd62ed3e146108dc575f80fd5b8063bfa0b1331161015e578063cfd6686311610139578063cfd66863146108b9578063d101dcd0146108c0578063d7419469146108c7575f80fd5b8063bfa0b1331461088b578063c2be97e314610634578063c958d4bf146108b2575f80fd5b8063b66dd7501161018e578063b66dd7501461087d578063b9a6d64514610884578063bb9bfe0614610634575f80fd5b8063ad8f42211461086f578063b1802b9a14610634578063b2bb360e14610876575f80fd5b80637dffdc321161027e57806395d89b4111610229578063a891d4d411610204578063a891d4d41461084e578063a9059cbb14610855578063aaa7af7014610868578063acc5aee914610634575f80fd5b806395d89b41146108385780639c5dfe73146108405780639df61a2514610847575f80fd5b80638619d607116102595780638619d6071461082a5780638789ca67146106345780638f4a840614610831575f80fd5b80637dffdc32146108155780637e5449371461081c5780637f34d94b14610823575f80fd5b806374f83d02116102de5780637a319c18116102b95780637a319c18146108075780637c66673e146106345780637c72ed0d1461080e575f80fd5b806374f83d02146107f257806377c0209e146107f9578063792c7f3e14610800575f80fd5b80636c12ed281161030e5780636c12ed281461063457806370a08231146107cc57806374e73fd3146107eb575f80fd5b8063657b6ef714610707578063672151fe146107be5780636abceacd146107c5575f80fd5b80633125f37a116104be5780634a2e93c611610409578063552a1b56116103b457806361b970eb1161038f57806361b970eb146107a2578063639ec53a146107a957806365473174146107b05780636578534c146107b7575f80fd5b8063552a1b561461063457806358b6a9bd146107945780635af92c051461079b575f80fd5b80634f5e5557116103e45780634f5e5557146107865780634f7bd75a1461063457806354c279201461078d575f80fd5b80634a2e93c6146107715780634b3c7f5f146107785780634e1dbb821461077f575f80fd5b80633ea117ce1161046957806342937dbd1161044457806342937dbd1461076357806343a6b92d1461076a57806344050a2814610634575f80fd5b80633ea117ce1461074e5780634128a85d14610755578063418b18161461075c575f80fd5b806339e0bd121161049957806339e0bd12146106345780633a131990146106345780633b6be45914610747575f80fd5b80633125f37a1461071a578063313ce5671461072157806334517f0b14610740575f80fd5b80631b17c65c1161057e57806321ecd7a3116105295780632545d8b7116105045780632545d8b7146106f95780632787325b14610700578063291c3bd714610707575f80fd5b806321ecd7a3146106eb578063239af2a5146106f257806323b872dd14610634575f80fd5b80631eaa7c52116105595780631eaa7c52146106d65780631f449589146106dd5780631fd298ec146106e4575f80fd5b80631b17c65c146106345780631bbffe6f146106345780631d527cde146106cf575f80fd5b806312901b42116105de57806318160ddd116105b957806318160ddd146106b857806319cf6a91146106c15780631a97f18e146106c8575f80fd5b806312901b42146106a357806313ebb5ec146106aa57806316a3045b146106b1575f80fd5b8063095ea7b31161060e578063095ea7b3146106825780630cb7a9e7146106955780631215a3ab1461069c575f80fd5b80630460faf61461063457806306fdde031461065c5780630717b16114610671575b5f80fd5b610647610642366004610fd2565b610952565b60405190151581526020015b60405180910390f35b610664610ba7565b604051610653919061100b565b602f5b604051908152602001610653565b610647610690366004611075565b610c32565b601b610674565b6005610674565b601a610674565b601d610674565b6033610674565b61067460035481565b6008610674565b6003610674565b6002610674565b6010610674565b6041610674565b6035610674565b602e610674565b602c610674565b6001610674565b6045610674565b610647610715366004610fd2565b610cab565b6034610674565b60025461072e9060ff1681565b60405160ff9091168152602001610653565b603c610674565b6004610674565b600c610674565b6006610674565b601c610674565b6032610674565b6023610674565b6011610674565b6021610674565b601f610674565b6026610674565b603a610674565b6009610674565b602a610674565b6014610674565b6025610674565b6013610674565b6043610674565b600d610674565b6017610674565b6106746107da36600461109d565b60046020525f908152604090205481565b6016610674565b6039610674565b603b610674565b6036610674565b603f610674565b602b610674565b6037610674565b600a610674565b602d610674565b6030610674565b600b610674565b610664610dd2565b6015610674565b6029610674565b6038610674565b610647610863366004611075565b610ddf565b6028610674565b6040610674565b6018610674565b600f610674565b6044610674565b6106747f000000000000000000000000000000000000000000000000000000000000000081565b6024610674565b601e610674565b6031610674565b6019610674565b6042610674565b6027610674565b6106746108ea3660046110bd565b600560209081525f928352604080842090915290825290205481565b603d610674565b6022610674565b610647610922366004610fd2565b610efd565b610664610f8b565b6007610674565b603e610674565b6020610674565b600e610674565b6012610674565b73ffffffffffffffffffffffffffffffffffffffff83165f908152600460205260408120548211156109e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e73756666696369656e742062616c616e636500000000000000000000000060448201526064015b60405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff84165f908152600560209081526040808320338452909152902054821115610a7e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e73756666696369656e7420616c6c6f77616e63650000000000000000000060448201526064016109dc565b73ffffffffffffffffffffffffffffffffffffffff84165f9081526004602052604081208054849290610ab290849061111b565b909155505073ffffffffffffffffffffffffffffffffffffffff83165f9081526004602052604081208054849290610aeb90849061112e565b909155505073ffffffffffffffffffffffffffffffffffffffff84165f90815260056020908152604080832033845290915281208054849290610b2f90849061111b565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610b9591815260200190565b60405180910390a35060019392505050565b5f8054610bb390611141565b80601f0160208091040260200160405190810160405280929190818152602001828054610bdf90611141565b8015610c2a5780601f10610c0157610100808354040283529160200191610c2a565b820191905f5260205f20905b815481529060010190602001808311610c0d57829003601f168201915b505050505081565b335f81815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716808552925280832085905551919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92590610c999086815260200190565b60405180910390a35060015b92915050565b73ffffffffffffffffffffffffffffffffffffffff83165f90815260046020526040812054821115610d39576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600160248201527f410000000000000000000000000000000000000000000000000000000000000060448201526064016109dc565b73ffffffffffffffffffffffffffffffffffffffff84165f908152600560209081526040808320338452909152902054821115610a7e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600160248201527f420000000000000000000000000000000000000000000000000000000000000060448201526064016109dc565b60018054610bb390611141565b335f90815260046020526040812054821115610e57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e73756666696369656e742062616c616e636500000000000000000000000060448201526064016109dc565b335f9081526004602052604081208054849290610e7590849061111b565b909155505073ffffffffffffffffffffffffffffffffffffffff83165f9081526004602052604081208054849290610eae90849061112e565b909155505060405182815273ffffffffffffffffffffffffffffffffffffffff84169033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610c99565b73ffffffffffffffffffffffffffffffffffffffff83165f90815260046020526040812054821115610a7e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600160248201527f410000000000000000000000000000000000000000000000000000000000000060448201526064016109dc565b6040518061016001604052806101368152602001611193610136913981565b803573ffffffffffffffffffffffffffffffffffffffff81168114610fcd575f80fd5b919050565b5f805f60608486031215610fe4575f80fd5b610fed84610faa565b9250610ffb60208501610faa565b9150604084013590509250925092565b5f602080835283518060208501525f5b818110156110375785810183015185820160400152820161101b565b505f6040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b5f8060408385031215611086575f80fd5b61108f83610faa565b946020939093013593505050565b5f602082840312156110ad575f80fd5b6110b682610faa565b9392505050565b5f80604083850312156110ce575f80fd5b6110d783610faa565b91506110e560208401610faa565b90509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b81810381811115610ca557610ca56110ee565b80820180821115610ca557610ca56110ee565b600181811c9082168061115557607f821691505b60208210810361118c577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5091905056fe41414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141a2646970667358221220c70035169e4d38ee509e85ab2f5566c500caae43b66c105f4552e44745217d5f64736f6c63430008160033",
}

// StateBloatTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use StateBloatTokenMetaData.ABI instead.
var StateBloatTokenABI = StateBloatTokenMetaData.ABI

// StateBloatTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StateBloatTokenMetaData.Bin instead.
var StateBloatTokenBin = StateBloatTokenMetaData.Bin

// DeployStateBloatToken deploys a new Ethereum contract, binding an instance of StateBloatToken to it.
func DeployStateBloatToken(auth *bind.TransactOpts, backend bind.ContractBackend, _salt *big.Int) (common.Address, *types.Transaction, *StateBloatToken, error) {
	parsed, err := StateBloatTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StateBloatTokenBin), backend, _salt)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StateBloatToken{StateBloatTokenCaller: StateBloatTokenCaller{contract: contract}, StateBloatTokenTransactor: StateBloatTokenTransactor{contract: contract}, StateBloatTokenFilterer: StateBloatTokenFilterer{contract: contract}}, nil
}

// StateBloatToken is an auto generated Go binding around an Ethereum contract.
type StateBloatToken struct {
	StateBloatTokenCaller     // Read-only binding to the contract
	StateBloatTokenTransactor // Write-only binding to the contract
	StateBloatTokenFilterer   // Log filterer for contract events
}

// StateBloatTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type StateBloatTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateBloatTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StateBloatTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateBloatTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StateBloatTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateBloatTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StateBloatTokenSession struct {
	Contract     *StateBloatToken  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateBloatTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StateBloatTokenCallerSession struct {
	Contract *StateBloatTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// StateBloatTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StateBloatTokenTransactorSession struct {
	Contract     *StateBloatTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// StateBloatTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type StateBloatTokenRaw struct {
	Contract *StateBloatToken // Generic contract binding to access the raw methods on
}

// StateBloatTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StateBloatTokenCallerRaw struct {
	Contract *StateBloatTokenCaller // Generic read-only contract binding to access the raw methods on
}

// StateBloatTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StateBloatTokenTransactorRaw struct {
	Contract *StateBloatTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStateBloatToken creates a new instance of StateBloatToken, bound to a specific deployed contract.
func NewStateBloatToken(address common.Address, backend bind.ContractBackend) (*StateBloatToken, error) {
	contract, err := bindStateBloatToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StateBloatToken{StateBloatTokenCaller: StateBloatTokenCaller{contract: contract}, StateBloatTokenTransactor: StateBloatTokenTransactor{contract: contract}, StateBloatTokenFilterer: StateBloatTokenFilterer{contract: contract}}, nil
}

// NewStateBloatTokenCaller creates a new read-only instance of StateBloatToken, bound to a specific deployed contract.
func NewStateBloatTokenCaller(address common.Address, caller bind.ContractCaller) (*StateBloatTokenCaller, error) {
	contract, err := bindStateBloatToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StateBloatTokenCaller{contract: contract}, nil
}

// NewStateBloatTokenTransactor creates a new write-only instance of StateBloatToken, bound to a specific deployed contract.
func NewStateBloatTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*StateBloatTokenTransactor, error) {
	contract, err := bindStateBloatToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StateBloatTokenTransactor{contract: contract}, nil
}

// NewStateBloatTokenFilterer creates a new log filterer instance of StateBloatToken, bound to a specific deployed contract.
func NewStateBloatTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*StateBloatTokenFilterer, error) {
	contract, err := bindStateBloatToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StateBloatTokenFilterer{contract: contract}, nil
}

// bindStateBloatToken binds a generic wrapper to an already deployed contract.
func bindStateBloatToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StateBloatTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateBloatToken *StateBloatTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateBloatToken.Contract.StateBloatTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateBloatToken *StateBloatTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateBloatToken.Contract.StateBloatTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateBloatToken *StateBloatTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateBloatToken.Contract.StateBloatTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateBloatToken *StateBloatTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateBloatToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateBloatToken *StateBloatTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateBloatToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateBloatToken *StateBloatTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateBloatToken.Contract.contract.Transact(opts, method, params...)
}

// PADDINGDATA is a free data retrieval call binding the contract method 0xee1682b6.
//
// Solidity: function PADDING_DATA() view returns(string)
func (_StateBloatToken *StateBloatTokenCaller) PADDINGDATA(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "PADDING_DATA")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PADDINGDATA is a free data retrieval call binding the contract method 0xee1682b6.
//
// Solidity: function PADDING_DATA() view returns(string)
func (_StateBloatToken *StateBloatTokenSession) PADDINGDATA() (string, error) {
	return _StateBloatToken.Contract.PADDINGDATA(&_StateBloatToken.CallOpts)
}

// PADDINGDATA is a free data retrieval call binding the contract method 0xee1682b6.
//
// Solidity: function PADDING_DATA() view returns(string)
func (_StateBloatToken *StateBloatTokenCallerSession) PADDINGDATA() (string, error) {
	return _StateBloatToken.Contract.PADDINGDATA(&_StateBloatToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _StateBloatToken.Contract.Allowance(&_StateBloatToken.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _StateBloatToken.Contract.Allowance(&_StateBloatToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _StateBloatToken.Contract.BalanceOf(&_StateBloatToken.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _StateBloatToken.Contract.BalanceOf(&_StateBloatToken.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_StateBloatToken *StateBloatTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_StateBloatToken *StateBloatTokenSession) Decimals() (uint8, error) {
	return _StateBloatToken.Contract.Decimals(&_StateBloatToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_StateBloatToken *StateBloatTokenCallerSession) Decimals() (uint8, error) {
	return _StateBloatToken.Contract.Decimals(&_StateBloatToken.CallOpts)
}

// Dummy1 is a free data retrieval call binding the contract method 0x2545d8b7.
//
// Solidity: function dummy1() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy1 is a free data retrieval call binding the contract method 0x2545d8b7.
//
// Solidity: function dummy1() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy1() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy1(&_StateBloatToken.CallOpts)
}

// Dummy1 is a free data retrieval call binding the contract method 0x2545d8b7.
//
// Solidity: function dummy1() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy1() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy1(&_StateBloatToken.CallOpts)
}

// Dummy10 is a free data retrieval call binding the contract method 0x7e544937.
//
// Solidity: function dummy10() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy10(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy10")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy10 is a free data retrieval call binding the contract method 0x7e544937.
//
// Solidity: function dummy10() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy10() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy10(&_StateBloatToken.CallOpts)
}

// Dummy10 is a free data retrieval call binding the contract method 0x7e544937.
//
// Solidity: function dummy10() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy10() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy10(&_StateBloatToken.CallOpts)
}

// Dummy11 is a free data retrieval call binding the contract method 0x8f4a8406.
//
// Solidity: function dummy11() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy11(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy11")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy11 is a free data retrieval call binding the contract method 0x8f4a8406.
//
// Solidity: function dummy11() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy11() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy11(&_StateBloatToken.CallOpts)
}

// Dummy11 is a free data retrieval call binding the contract method 0x8f4a8406.
//
// Solidity: function dummy11() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy11() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy11(&_StateBloatToken.CallOpts)
}

// Dummy12 is a free data retrieval call binding the contract method 0x3ea117ce.
//
// Solidity: function dummy12() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy12(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy12")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy12 is a free data retrieval call binding the contract method 0x3ea117ce.
//
// Solidity: function dummy12() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy12() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy12(&_StateBloatToken.CallOpts)
}

// Dummy12 is a free data retrieval call binding the contract method 0x3ea117ce.
//
// Solidity: function dummy12() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy12() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy12(&_StateBloatToken.CallOpts)
}

// Dummy13 is a free data retrieval call binding the contract method 0x672151fe.
//
// Solidity: function dummy13() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy13(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy13")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy13 is a free data retrieval call binding the contract method 0x672151fe.
//
// Solidity: function dummy13() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy13() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy13(&_StateBloatToken.CallOpts)
}

// Dummy13 is a free data retrieval call binding the contract method 0x672151fe.
//
// Solidity: function dummy13() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy13() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy13(&_StateBloatToken.CallOpts)
}

// Dummy14 is a free data retrieval call binding the contract method 0xfe7d5996.
//
// Solidity: function dummy14() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy14(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy14")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy14 is a free data retrieval call binding the contract method 0xfe7d5996.
//
// Solidity: function dummy14() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy14() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy14(&_StateBloatToken.CallOpts)
}

// Dummy14 is a free data retrieval call binding the contract method 0xfe7d5996.
//
// Solidity: function dummy14() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy14() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy14(&_StateBloatToken.CallOpts)
}

// Dummy15 is a free data retrieval call binding the contract method 0xb66dd750.
//
// Solidity: function dummy15() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy15(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy15")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy15 is a free data retrieval call binding the contract method 0xb66dd750.
//
// Solidity: function dummy15() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy15() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy15(&_StateBloatToken.CallOpts)
}

// Dummy15 is a free data retrieval call binding the contract method 0xb66dd750.
//
// Solidity: function dummy15() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy15() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy15(&_StateBloatToken.CallOpts)
}

// Dummy16 is a free data retrieval call binding the contract method 0x1eaa7c52.
//
// Solidity: function dummy16() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy16(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy16")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy16 is a free data retrieval call binding the contract method 0x1eaa7c52.
//
// Solidity: function dummy16() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy16() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy16(&_StateBloatToken.CallOpts)
}

// Dummy16 is a free data retrieval call binding the contract method 0x1eaa7c52.
//
// Solidity: function dummy16() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy16() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy16(&_StateBloatToken.CallOpts)
}

// Dummy17 is a free data retrieval call binding the contract method 0x4a2e93c6.
//
// Solidity: function dummy17() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy17(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy17")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy17 is a free data retrieval call binding the contract method 0x4a2e93c6.
//
// Solidity: function dummy17() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy17() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy17(&_StateBloatToken.CallOpts)
}

// Dummy17 is a free data retrieval call binding the contract method 0x4a2e93c6.
//
// Solidity: function dummy17() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy17() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy17(&_StateBloatToken.CallOpts)
}

// Dummy18 is a free data retrieval call binding the contract method 0xffbf0469.
//
// Solidity: function dummy18() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy18(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy18")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy18 is a free data retrieval call binding the contract method 0xffbf0469.
//
// Solidity: function dummy18() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy18() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy18(&_StateBloatToken.CallOpts)
}

// Dummy18 is a free data retrieval call binding the contract method 0xffbf0469.
//
// Solidity: function dummy18() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy18() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy18(&_StateBloatToken.CallOpts)
}

// Dummy19 is a free data retrieval call binding the contract method 0x65473174.
//
// Solidity: function dummy19() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy19(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy19")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy19 is a free data retrieval call binding the contract method 0x65473174.
//
// Solidity: function dummy19() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy19() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy19(&_StateBloatToken.CallOpts)
}

// Dummy19 is a free data retrieval call binding the contract method 0x65473174.
//
// Solidity: function dummy19() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy19() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy19(&_StateBloatToken.CallOpts)
}

// Dummy2 is a free data retrieval call binding the contract method 0x1d527cde.
//
// Solidity: function dummy2() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy2 is a free data retrieval call binding the contract method 0x1d527cde.
//
// Solidity: function dummy2() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy2() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy2(&_StateBloatToken.CallOpts)
}

// Dummy2 is a free data retrieval call binding the contract method 0x1d527cde.
//
// Solidity: function dummy2() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy2() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy2(&_StateBloatToken.CallOpts)
}

// Dummy20 is a free data retrieval call binding the contract method 0x61b970eb.
//
// Solidity: function dummy20() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy20(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy20")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy20 is a free data retrieval call binding the contract method 0x61b970eb.
//
// Solidity: function dummy20() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy20() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy20(&_StateBloatToken.CallOpts)
}

// Dummy20 is a free data retrieval call binding the contract method 0x61b970eb.
//
// Solidity: function dummy20() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy20() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy20(&_StateBloatToken.CallOpts)
}

// Dummy21 is a free data retrieval call binding the contract method 0x9c5dfe73.
//
// Solidity: function dummy21() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy21(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy21")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy21 is a free data retrieval call binding the contract method 0x9c5dfe73.
//
// Solidity: function dummy21() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy21() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy21(&_StateBloatToken.CallOpts)
}

// Dummy21 is a free data retrieval call binding the contract method 0x9c5dfe73.
//
// Solidity: function dummy21() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy21() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy21(&_StateBloatToken.CallOpts)
}

// Dummy22 is a free data retrieval call binding the contract method 0x74e73fd3.
//
// Solidity: function dummy22() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy22(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy22")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy22 is a free data retrieval call binding the contract method 0x74e73fd3.
//
// Solidity: function dummy22() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy22() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy22(&_StateBloatToken.CallOpts)
}

// Dummy22 is a free data retrieval call binding the contract method 0x74e73fd3.
//
// Solidity: function dummy22() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy22() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy22(&_StateBloatToken.CallOpts)
}

// Dummy23 is a free data retrieval call binding the contract method 0x6abceacd.
//
// Solidity: function dummy23() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy23(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy23")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy23 is a free data retrieval call binding the contract method 0x6abceacd.
//
// Solidity: function dummy23() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy23() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy23(&_StateBloatToken.CallOpts)
}

// Dummy23 is a free data retrieval call binding the contract method 0x6abceacd.
//
// Solidity: function dummy23() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy23() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy23(&_StateBloatToken.CallOpts)
}

// Dummy24 is a free data retrieval call binding the contract method 0xb2bb360e.
//
// Solidity: function dummy24() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy24(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy24")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy24 is a free data retrieval call binding the contract method 0xb2bb360e.
//
// Solidity: function dummy24() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy24() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy24(&_StateBloatToken.CallOpts)
}

// Dummy24 is a free data retrieval call binding the contract method 0xb2bb360e.
//
// Solidity: function dummy24() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy24() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy24(&_StateBloatToken.CallOpts)
}

// Dummy25 is a free data retrieval call binding the contract method 0xd7419469.
//
// Solidity: function dummy25() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy25(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy25")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy25 is a free data retrieval call binding the contract method 0xd7419469.
//
// Solidity: function dummy25() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy25() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy25(&_StateBloatToken.CallOpts)
}

// Dummy25 is a free data retrieval call binding the contract method 0xd7419469.
//
// Solidity: function dummy25() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy25() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy25(&_StateBloatToken.CallOpts)
}

// Dummy26 is a free data retrieval call binding the contract method 0x12901b42.
//
// Solidity: function dummy26() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy26(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy26")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy26 is a free data retrieval call binding the contract method 0x12901b42.
//
// Solidity: function dummy26() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy26() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy26(&_StateBloatToken.CallOpts)
}

// Dummy26 is a free data retrieval call binding the contract method 0x12901b42.
//
// Solidity: function dummy26() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy26() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy26(&_StateBloatToken.CallOpts)
}

// Dummy27 is a free data retrieval call binding the contract method 0x0cb7a9e7.
//
// Solidity: function dummy27() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy27(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy27")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy27 is a free data retrieval call binding the contract method 0x0cb7a9e7.
//
// Solidity: function dummy27() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy27() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy27(&_StateBloatToken.CallOpts)
}

// Dummy27 is a free data retrieval call binding the contract method 0x0cb7a9e7.
//
// Solidity: function dummy27() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy27() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy27(&_StateBloatToken.CallOpts)
}

// Dummy28 is a free data retrieval call binding the contract method 0x418b1816.
//
// Solidity: function dummy28() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy28(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy28")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy28 is a free data retrieval call binding the contract method 0x418b1816.
//
// Solidity: function dummy28() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy28() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy28(&_StateBloatToken.CallOpts)
}

// Dummy28 is a free data retrieval call binding the contract method 0x418b1816.
//
// Solidity: function dummy28() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy28() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy28(&_StateBloatToken.CallOpts)
}

// Dummy29 is a free data retrieval call binding the contract method 0x13ebb5ec.
//
// Solidity: function dummy29() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy29(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy29")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy29 is a free data retrieval call binding the contract method 0x13ebb5ec.
//
// Solidity: function dummy29() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy29() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy29(&_StateBloatToken.CallOpts)
}

// Dummy29 is a free data retrieval call binding the contract method 0x13ebb5ec.
//
// Solidity: function dummy29() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy29() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy29(&_StateBloatToken.CallOpts)
}

// Dummy3 is a free data retrieval call binding the contract method 0x1a97f18e.
//
// Solidity: function dummy3() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy3(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy3")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy3 is a free data retrieval call binding the contract method 0x1a97f18e.
//
// Solidity: function dummy3() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy3() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy3(&_StateBloatToken.CallOpts)
}

// Dummy3 is a free data retrieval call binding the contract method 0x1a97f18e.
//
// Solidity: function dummy3() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy3() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy3(&_StateBloatToken.CallOpts)
}

// Dummy30 is a free data retrieval call binding the contract method 0xcfd66863.
//
// Solidity: function dummy30() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy30(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy30")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy30 is a free data retrieval call binding the contract method 0xcfd66863.
//
// Solidity: function dummy30() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy30() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy30(&_StateBloatToken.CallOpts)
}

// Dummy30 is a free data retrieval call binding the contract method 0xcfd66863.
//
// Solidity: function dummy30() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy30() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy30(&_StateBloatToken.CallOpts)
}

// Dummy31 is a free data retrieval call binding the contract method 0x4e1dbb82.
//
// Solidity: function dummy31() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy31(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy31")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy31 is a free data retrieval call binding the contract method 0x4e1dbb82.
//
// Solidity: function dummy31() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy31() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy31(&_StateBloatToken.CallOpts)
}

// Dummy31 is a free data retrieval call binding the contract method 0x4e1dbb82.
//
// Solidity: function dummy31() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy31() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy31(&_StateBloatToken.CallOpts)
}

// Dummy32 is a free data retrieval call binding the contract method 0xf8716f14.
//
// Solidity: function dummy32() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy32(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy32")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy32 is a free data retrieval call binding the contract method 0xf8716f14.
//
// Solidity: function dummy32() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy32() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy32(&_StateBloatToken.CallOpts)
}

// Dummy32 is a free data retrieval call binding the contract method 0xf8716f14.
//
// Solidity: function dummy32() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy32() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy32(&_StateBloatToken.CallOpts)
}

// Dummy33 is a free data retrieval call binding the contract method 0x4b3c7f5f.
//
// Solidity: function dummy33() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy33(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy33")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy33 is a free data retrieval call binding the contract method 0x4b3c7f5f.
//
// Solidity: function dummy33() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy33() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy33(&_StateBloatToken.CallOpts)
}

// Dummy33 is a free data retrieval call binding the contract method 0x4b3c7f5f.
//
// Solidity: function dummy33() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy33() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy33(&_StateBloatToken.CallOpts)
}

// Dummy34 is a free data retrieval call binding the contract method 0xe8c927b3.
//
// Solidity: function dummy34() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy34(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy34")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy34 is a free data retrieval call binding the contract method 0xe8c927b3.
//
// Solidity: function dummy34() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy34() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy34(&_StateBloatToken.CallOpts)
}

// Dummy34 is a free data retrieval call binding the contract method 0xe8c927b3.
//
// Solidity: function dummy34() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy34() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy34(&_StateBloatToken.CallOpts)
}

// Dummy35 is a free data retrieval call binding the contract method 0x43a6b92d.
//
// Solidity: function dummy35() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy35(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy35")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy35 is a free data retrieval call binding the contract method 0x43a6b92d.
//
// Solidity: function dummy35() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy35() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy35(&_StateBloatToken.CallOpts)
}

// Dummy35 is a free data retrieval call binding the contract method 0x43a6b92d.
//
// Solidity: function dummy35() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy35() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy35(&_StateBloatToken.CallOpts)
}

// Dummy36 is a free data retrieval call binding the contract method 0xc958d4bf.
//
// Solidity: function dummy36() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy36(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy36")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy36 is a free data retrieval call binding the contract method 0xc958d4bf.
//
// Solidity: function dummy36() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy36() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy36(&_StateBloatToken.CallOpts)
}

// Dummy36 is a free data retrieval call binding the contract method 0xc958d4bf.
//
// Solidity: function dummy36() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy36() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy36(&_StateBloatToken.CallOpts)
}

// Dummy37 is a free data retrieval call binding the contract method 0x639ec53a.
//
// Solidity: function dummy37() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy37(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy37")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy37 is a free data retrieval call binding the contract method 0x639ec53a.
//
// Solidity: function dummy37() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy37() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy37(&_StateBloatToken.CallOpts)
}

// Dummy37 is a free data retrieval call binding the contract method 0x639ec53a.
//
// Solidity: function dummy37() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy37() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy37(&_StateBloatToken.CallOpts)
}

// Dummy38 is a free data retrieval call binding the contract method 0x4f5e5557.
//
// Solidity: function dummy38() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy38(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy38")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy38 is a free data retrieval call binding the contract method 0x4f5e5557.
//
// Solidity: function dummy38() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy38() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy38(&_StateBloatToken.CallOpts)
}

// Dummy38 is a free data retrieval call binding the contract method 0x4f5e5557.
//
// Solidity: function dummy38() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy38() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy38(&_StateBloatToken.CallOpts)
}

// Dummy39 is a free data retrieval call binding the contract method 0xdc1d8a9b.
//
// Solidity: function dummy39() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy39(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy39")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy39 is a free data retrieval call binding the contract method 0xdc1d8a9b.
//
// Solidity: function dummy39() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy39() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy39(&_StateBloatToken.CallOpts)
}

// Dummy39 is a free data retrieval call binding the contract method 0xdc1d8a9b.
//
// Solidity: function dummy39() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy39() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy39(&_StateBloatToken.CallOpts)
}

// Dummy4 is a free data retrieval call binding the contract method 0x3b6be459.
//
// Solidity: function dummy4() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy4(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy4")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy4 is a free data retrieval call binding the contract method 0x3b6be459.
//
// Solidity: function dummy4() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy4() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy4(&_StateBloatToken.CallOpts)
}

// Dummy4 is a free data retrieval call binding the contract method 0x3b6be459.
//
// Solidity: function dummy4() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy4() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy4(&_StateBloatToken.CallOpts)
}

// Dummy40 is a free data retrieval call binding the contract method 0xaaa7af70.
//
// Solidity: function dummy40() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy40(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy40")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy40 is a free data retrieval call binding the contract method 0xaaa7af70.
//
// Solidity: function dummy40() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy40() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy40(&_StateBloatToken.CallOpts)
}

// Dummy40 is a free data retrieval call binding the contract method 0xaaa7af70.
//
// Solidity: function dummy40() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy40() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy40(&_StateBloatToken.CallOpts)
}

// Dummy41 is a free data retrieval call binding the contract method 0x9df61a25.
//
// Solidity: function dummy41() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy41(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy41")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy41 is a free data retrieval call binding the contract method 0x9df61a25.
//
// Solidity: function dummy41() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy41() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy41(&_StateBloatToken.CallOpts)
}

// Dummy41 is a free data retrieval call binding the contract method 0x9df61a25.
//
// Solidity: function dummy41() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy41() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy41(&_StateBloatToken.CallOpts)
}

// Dummy42 is a free data retrieval call binding the contract method 0x5af92c05.
//
// Solidity: function dummy42() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy42(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy42")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy42 is a free data retrieval call binding the contract method 0x5af92c05.
//
// Solidity: function dummy42() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy42() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy42(&_StateBloatToken.CallOpts)
}

// Dummy42 is a free data retrieval call binding the contract method 0x5af92c05.
//
// Solidity: function dummy42() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy42() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy42(&_StateBloatToken.CallOpts)
}

// Dummy43 is a free data retrieval call binding the contract method 0x7c72ed0d.
//
// Solidity: function dummy43() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy43(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy43")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy43 is a free data retrieval call binding the contract method 0x7c72ed0d.
//
// Solidity: function dummy43() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy43() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy43(&_StateBloatToken.CallOpts)
}

// Dummy43 is a free data retrieval call binding the contract method 0x7c72ed0d.
//
// Solidity: function dummy43() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy43() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy43(&_StateBloatToken.CallOpts)
}

// Dummy44 is a free data retrieval call binding the contract method 0x239af2a5.
//
// Solidity: function dummy44() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy44(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy44")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy44 is a free data retrieval call binding the contract method 0x239af2a5.
//
// Solidity: function dummy44() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy44() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy44(&_StateBloatToken.CallOpts)
}

// Dummy44 is a free data retrieval call binding the contract method 0x239af2a5.
//
// Solidity: function dummy44() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy44() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy44(&_StateBloatToken.CallOpts)
}

// Dummy45 is a free data retrieval call binding the contract method 0x7f34d94b.
//
// Solidity: function dummy45() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy45(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy45")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy45 is a free data retrieval call binding the contract method 0x7f34d94b.
//
// Solidity: function dummy45() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy45() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy45(&_StateBloatToken.CallOpts)
}

// Dummy45 is a free data retrieval call binding the contract method 0x7f34d94b.
//
// Solidity: function dummy45() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy45() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy45(&_StateBloatToken.CallOpts)
}

// Dummy46 is a free data retrieval call binding the contract method 0x21ecd7a3.
//
// Solidity: function dummy46() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy46(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy46")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy46 is a free data retrieval call binding the contract method 0x21ecd7a3.
//
// Solidity: function dummy46() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy46() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy46(&_StateBloatToken.CallOpts)
}

// Dummy46 is a free data retrieval call binding the contract method 0x21ecd7a3.
//
// Solidity: function dummy46() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy46() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy46(&_StateBloatToken.CallOpts)
}

// Dummy47 is a free data retrieval call binding the contract method 0x0717b161.
//
// Solidity: function dummy47() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy47(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy47")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy47 is a free data retrieval call binding the contract method 0x0717b161.
//
// Solidity: function dummy47() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy47() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy47(&_StateBloatToken.CallOpts)
}

// Dummy47 is a free data retrieval call binding the contract method 0x0717b161.
//
// Solidity: function dummy47() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy47() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy47(&_StateBloatToken.CallOpts)
}

// Dummy48 is a free data retrieval call binding the contract method 0x8619d607.
//
// Solidity: function dummy48() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy48(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy48")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy48 is a free data retrieval call binding the contract method 0x8619d607.
//
// Solidity: function dummy48() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy48() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy48(&_StateBloatToken.CallOpts)
}

// Dummy48 is a free data retrieval call binding the contract method 0x8619d607.
//
// Solidity: function dummy48() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy48() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy48(&_StateBloatToken.CallOpts)
}

// Dummy49 is a free data retrieval call binding the contract method 0xd101dcd0.
//
// Solidity: function dummy49() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy49(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy49")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy49 is a free data retrieval call binding the contract method 0xd101dcd0.
//
// Solidity: function dummy49() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy49() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy49(&_StateBloatToken.CallOpts)
}

// Dummy49 is a free data retrieval call binding the contract method 0xd101dcd0.
//
// Solidity: function dummy49() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy49() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy49(&_StateBloatToken.CallOpts)
}

// Dummy5 is a free data retrieval call binding the contract method 0x1215a3ab.
//
// Solidity: function dummy5() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy5(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy5")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy5 is a free data retrieval call binding the contract method 0x1215a3ab.
//
// Solidity: function dummy5() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy5() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy5(&_StateBloatToken.CallOpts)
}

// Dummy5 is a free data retrieval call binding the contract method 0x1215a3ab.
//
// Solidity: function dummy5() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy5() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy5(&_StateBloatToken.CallOpts)
}

// Dummy50 is a free data retrieval call binding the contract method 0x42937dbd.
//
// Solidity: function dummy50() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy50(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy50")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy50 is a free data retrieval call binding the contract method 0x42937dbd.
//
// Solidity: function dummy50() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy50() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy50(&_StateBloatToken.CallOpts)
}

// Dummy50 is a free data retrieval call binding the contract method 0x42937dbd.
//
// Solidity: function dummy50() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy50() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy50(&_StateBloatToken.CallOpts)
}

// Dummy51 is a free data retrieval call binding the contract method 0x16a3045b.
//
// Solidity: function dummy51() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy51(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy51")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy51 is a free data retrieval call binding the contract method 0x16a3045b.
//
// Solidity: function dummy51() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy51() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy51(&_StateBloatToken.CallOpts)
}

// Dummy51 is a free data retrieval call binding the contract method 0x16a3045b.
//
// Solidity: function dummy51() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy51() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy51(&_StateBloatToken.CallOpts)
}

// Dummy52 is a free data retrieval call binding the contract method 0x3125f37a.
//
// Solidity: function dummy52() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy52(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy52")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy52 is a free data retrieval call binding the contract method 0x3125f37a.
//
// Solidity: function dummy52() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy52() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy52(&_StateBloatToken.CallOpts)
}

// Dummy52 is a free data retrieval call binding the contract method 0x3125f37a.
//
// Solidity: function dummy52() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy52() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy52(&_StateBloatToken.CallOpts)
}

// Dummy53 is a free data retrieval call binding the contract method 0x1fd298ec.
//
// Solidity: function dummy53() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy53(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy53")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy53 is a free data retrieval call binding the contract method 0x1fd298ec.
//
// Solidity: function dummy53() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy53() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy53(&_StateBloatToken.CallOpts)
}

// Dummy53 is a free data retrieval call binding the contract method 0x1fd298ec.
//
// Solidity: function dummy53() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy53() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy53(&_StateBloatToken.CallOpts)
}

// Dummy54 is a free data retrieval call binding the contract method 0x792c7f3e.
//
// Solidity: function dummy54() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy54(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy54")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy54 is a free data retrieval call binding the contract method 0x792c7f3e.
//
// Solidity: function dummy54() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy54() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy54(&_StateBloatToken.CallOpts)
}

// Dummy54 is a free data retrieval call binding the contract method 0x792c7f3e.
//
// Solidity: function dummy54() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy54() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy54(&_StateBloatToken.CallOpts)
}

// Dummy55 is a free data retrieval call binding the contract method 0x7dffdc32.
//
// Solidity: function dummy55() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy55(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy55")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy55 is a free data retrieval call binding the contract method 0x7dffdc32.
//
// Solidity: function dummy55() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy55() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy55(&_StateBloatToken.CallOpts)
}

// Dummy55 is a free data retrieval call binding the contract method 0x7dffdc32.
//
// Solidity: function dummy55() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy55() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy55(&_StateBloatToken.CallOpts)
}

// Dummy56 is a free data retrieval call binding the contract method 0xa891d4d4.
//
// Solidity: function dummy56() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy56(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy56")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy56 is a free data retrieval call binding the contract method 0xa891d4d4.
//
// Solidity: function dummy56() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy56() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy56(&_StateBloatToken.CallOpts)
}

// Dummy56 is a free data retrieval call binding the contract method 0xa891d4d4.
//
// Solidity: function dummy56() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy56() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy56(&_StateBloatToken.CallOpts)
}

// Dummy57 is a free data retrieval call binding the contract method 0x74f83d02.
//
// Solidity: function dummy57() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy57(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy57")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy57 is a free data retrieval call binding the contract method 0x74f83d02.
//
// Solidity: function dummy57() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy57() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy57(&_StateBloatToken.CallOpts)
}

// Dummy57 is a free data retrieval call binding the contract method 0x74f83d02.
//
// Solidity: function dummy57() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy57() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy57(&_StateBloatToken.CallOpts)
}

// Dummy58 is a free data retrieval call binding the contract method 0x54c27920.
//
// Solidity: function dummy58() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy58(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy58")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy58 is a free data retrieval call binding the contract method 0x54c27920.
//
// Solidity: function dummy58() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy58() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy58(&_StateBloatToken.CallOpts)
}

// Dummy58 is a free data retrieval call binding the contract method 0x54c27920.
//
// Solidity: function dummy58() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy58() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy58(&_StateBloatToken.CallOpts)
}

// Dummy59 is a free data retrieval call binding the contract method 0x77c0209e.
//
// Solidity: function dummy59() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy59(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy59")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy59 is a free data retrieval call binding the contract method 0x77c0209e.
//
// Solidity: function dummy59() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy59() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy59(&_StateBloatToken.CallOpts)
}

// Dummy59 is a free data retrieval call binding the contract method 0x77c0209e.
//
// Solidity: function dummy59() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy59() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy59(&_StateBloatToken.CallOpts)
}

// Dummy6 is a free data retrieval call binding the contract method 0x4128a85d.
//
// Solidity: function dummy6() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy6(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy6")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy6 is a free data retrieval call binding the contract method 0x4128a85d.
//
// Solidity: function dummy6() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy6() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy6(&_StateBloatToken.CallOpts)
}

// Dummy6 is a free data retrieval call binding the contract method 0x4128a85d.
//
// Solidity: function dummy6() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy6() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy6(&_StateBloatToken.CallOpts)
}

// Dummy60 is a free data retrieval call binding the contract method 0x34517f0b.
//
// Solidity: function dummy60() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy60(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy60")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy60 is a free data retrieval call binding the contract method 0x34517f0b.
//
// Solidity: function dummy60() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy60() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy60(&_StateBloatToken.CallOpts)
}

// Dummy60 is a free data retrieval call binding the contract method 0x34517f0b.
//
// Solidity: function dummy60() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy60() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy60(&_StateBloatToken.CallOpts)
}

// Dummy61 is a free data retrieval call binding the contract method 0xe2d27530.
//
// Solidity: function dummy61() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy61(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy61")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy61 is a free data retrieval call binding the contract method 0xe2d27530.
//
// Solidity: function dummy61() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy61() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy61(&_StateBloatToken.CallOpts)
}

// Dummy61 is a free data retrieval call binding the contract method 0xe2d27530.
//
// Solidity: function dummy61() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy61() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy61(&_StateBloatToken.CallOpts)
}

// Dummy62 is a free data retrieval call binding the contract method 0xf5f57381.
//
// Solidity: function dummy62() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy62(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy62")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy62 is a free data retrieval call binding the contract method 0xf5f57381.
//
// Solidity: function dummy62() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy62() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy62(&_StateBloatToken.CallOpts)
}

// Dummy62 is a free data retrieval call binding the contract method 0xf5f57381.
//
// Solidity: function dummy62() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy62() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy62(&_StateBloatToken.CallOpts)
}

// Dummy63 is a free data retrieval call binding the contract method 0x7a319c18.
//
// Solidity: function dummy63() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy63(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy63")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy63 is a free data retrieval call binding the contract method 0x7a319c18.
//
// Solidity: function dummy63() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy63() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy63(&_StateBloatToken.CallOpts)
}

// Dummy63 is a free data retrieval call binding the contract method 0x7a319c18.
//
// Solidity: function dummy63() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy63() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy63(&_StateBloatToken.CallOpts)
}

// Dummy64 is a free data retrieval call binding the contract method 0xad8f4221.
//
// Solidity: function dummy64() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy64(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy64")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy64 is a free data retrieval call binding the contract method 0xad8f4221.
//
// Solidity: function dummy64() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy64() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy64(&_StateBloatToken.CallOpts)
}

// Dummy64 is a free data retrieval call binding the contract method 0xad8f4221.
//
// Solidity: function dummy64() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy64() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy64(&_StateBloatToken.CallOpts)
}

// Dummy65 is a free data retrieval call binding the contract method 0x1f449589.
//
// Solidity: function dummy65() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy65(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy65")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy65 is a free data retrieval call binding the contract method 0x1f449589.
//
// Solidity: function dummy65() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy65() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy65(&_StateBloatToken.CallOpts)
}

// Dummy65 is a free data retrieval call binding the contract method 0x1f449589.
//
// Solidity: function dummy65() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy65() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy65(&_StateBloatToken.CallOpts)
}

// Dummy66 is a free data retrieval call binding the contract method 0xd9bb3174.
//
// Solidity: function dummy66() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy66(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy66")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy66 is a free data retrieval call binding the contract method 0xd9bb3174.
//
// Solidity: function dummy66() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy66() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy66(&_StateBloatToken.CallOpts)
}

// Dummy66 is a free data retrieval call binding the contract method 0xd9bb3174.
//
// Solidity: function dummy66() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy66() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy66(&_StateBloatToken.CallOpts)
}

// Dummy67 is a free data retrieval call binding the contract method 0x6578534c.
//
// Solidity: function dummy67() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy67(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy67")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy67 is a free data retrieval call binding the contract method 0x6578534c.
//
// Solidity: function dummy67() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy67() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy67(&_StateBloatToken.CallOpts)
}

// Dummy67 is a free data retrieval call binding the contract method 0x6578534c.
//
// Solidity: function dummy67() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy67() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy67(&_StateBloatToken.CallOpts)
}

// Dummy68 is a free data retrieval call binding the contract method 0xb9a6d645.
//
// Solidity: function dummy68() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy68(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy68")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy68 is a free data retrieval call binding the contract method 0xb9a6d645.
//
// Solidity: function dummy68() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy68() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy68(&_StateBloatToken.CallOpts)
}

// Dummy68 is a free data retrieval call binding the contract method 0xb9a6d645.
//
// Solidity: function dummy68() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy68() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy68(&_StateBloatToken.CallOpts)
}

// Dummy69 is a free data retrieval call binding the contract method 0x2787325b.
//
// Solidity: function dummy69() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy69(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy69")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy69 is a free data retrieval call binding the contract method 0x2787325b.
//
// Solidity: function dummy69() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy69() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy69(&_StateBloatToken.CallOpts)
}

// Dummy69 is a free data retrieval call binding the contract method 0x2787325b.
//
// Solidity: function dummy69() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy69() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy69(&_StateBloatToken.CallOpts)
}

// Dummy7 is a free data retrieval call binding the contract method 0xf26c779b.
//
// Solidity: function dummy7() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy7(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy7")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy7 is a free data retrieval call binding the contract method 0xf26c779b.
//
// Solidity: function dummy7() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy7() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy7(&_StateBloatToken.CallOpts)
}

// Dummy7 is a free data retrieval call binding the contract method 0xf26c779b.
//
// Solidity: function dummy7() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy7() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy7(&_StateBloatToken.CallOpts)
}

// Dummy8 is a free data retrieval call binding the contract method 0x19cf6a91.
//
// Solidity: function dummy8() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy8(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy8")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy8 is a free data retrieval call binding the contract method 0x19cf6a91.
//
// Solidity: function dummy8() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy8() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy8(&_StateBloatToken.CallOpts)
}

// Dummy8 is a free data retrieval call binding the contract method 0x19cf6a91.
//
// Solidity: function dummy8() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy8() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy8(&_StateBloatToken.CallOpts)
}

// Dummy9 is a free data retrieval call binding the contract method 0x58b6a9bd.
//
// Solidity: function dummy9() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Dummy9(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "dummy9")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Dummy9 is a free data retrieval call binding the contract method 0x58b6a9bd.
//
// Solidity: function dummy9() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Dummy9() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy9(&_StateBloatToken.CallOpts)
}

// Dummy9 is a free data retrieval call binding the contract method 0x58b6a9bd.
//
// Solidity: function dummy9() pure returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Dummy9() (*big.Int, error) {
	return _StateBloatToken.Contract.Dummy9(&_StateBloatToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StateBloatToken *StateBloatTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StateBloatToken *StateBloatTokenSession) Name() (string, error) {
	return _StateBloatToken.Contract.Name(&_StateBloatToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StateBloatToken *StateBloatTokenCallerSession) Name() (string, error) {
	return _StateBloatToken.Contract.Name(&_StateBloatToken.CallOpts)
}

// Salt is a free data retrieval call binding the contract method 0xbfa0b133.
//
// Solidity: function salt() view returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) Salt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "salt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Salt is a free data retrieval call binding the contract method 0xbfa0b133.
//
// Solidity: function salt() view returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) Salt() (*big.Int, error) {
	return _StateBloatToken.Contract.Salt(&_StateBloatToken.CallOpts)
}

// Salt is a free data retrieval call binding the contract method 0xbfa0b133.
//
// Solidity: function salt() view returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) Salt() (*big.Int, error) {
	return _StateBloatToken.Contract.Salt(&_StateBloatToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StateBloatToken *StateBloatTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StateBloatToken *StateBloatTokenSession) Symbol() (string, error) {
	return _StateBloatToken.Contract.Symbol(&_StateBloatToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StateBloatToken *StateBloatTokenCallerSession) Symbol() (string, error) {
	return _StateBloatToken.Contract.Symbol(&_StateBloatToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_StateBloatToken *StateBloatTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateBloatToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_StateBloatToken *StateBloatTokenSession) TotalSupply() (*big.Int, error) {
	return _StateBloatToken.Contract.TotalSupply(&_StateBloatToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_StateBloatToken *StateBloatTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _StateBloatToken.Contract.TotalSupply(&_StateBloatToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.Approve(&_StateBloatToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.Approve(&_StateBloatToken.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.Transfer(&_StateBloatToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.Transfer(&_StateBloatToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0xbb9bfe06.
//
// Solidity: function transferFrom1(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom1(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom1", from, to, value)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0xbb9bfe06.
//
// Solidity: function transferFrom1(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom1(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom1(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom1 is a paid mutator transaction binding the contract method 0xbb9bfe06.
//
// Solidity: function transferFrom1(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom1(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom1(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom10 is a paid mutator transaction binding the contract method 0xb1802b9a.
//
// Solidity: function transferFrom10(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom10(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom10", from, to, value)
}

// TransferFrom10 is a paid mutator transaction binding the contract method 0xb1802b9a.
//
// Solidity: function transferFrom10(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom10(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom10(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom10 is a paid mutator transaction binding the contract method 0xb1802b9a.
//
// Solidity: function transferFrom10(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom10(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom10(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom11 is a paid mutator transaction binding the contract method 0xc2be97e3.
//
// Solidity: function transferFrom11(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom11(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom11", from, to, value)
}

// TransferFrom11 is a paid mutator transaction binding the contract method 0xc2be97e3.
//
// Solidity: function transferFrom11(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom11(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom11(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom11 is a paid mutator transaction binding the contract method 0xc2be97e3.
//
// Solidity: function transferFrom11(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom11(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom11(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom12 is a paid mutator transaction binding the contract method 0x44050a28.
//
// Solidity: function transferFrom12(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom12(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom12", from, to, value)
}

// TransferFrom12 is a paid mutator transaction binding the contract method 0x44050a28.
//
// Solidity: function transferFrom12(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom12(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom12(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom12 is a paid mutator transaction binding the contract method 0x44050a28.
//
// Solidity: function transferFrom12(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom12(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom12(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom13 is a paid mutator transaction binding the contract method 0xacc5aee9.
//
// Solidity: function transferFrom13(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom13(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom13", from, to, value)
}

// TransferFrom13 is a paid mutator transaction binding the contract method 0xacc5aee9.
//
// Solidity: function transferFrom13(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom13(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom13(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom13 is a paid mutator transaction binding the contract method 0xacc5aee9.
//
// Solidity: function transferFrom13(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom13(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom13(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom14 is a paid mutator transaction binding the contract method 0x1bbffe6f.
//
// Solidity: function transferFrom14(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom14(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom14", from, to, value)
}

// TransferFrom14 is a paid mutator transaction binding the contract method 0x1bbffe6f.
//
// Solidity: function transferFrom14(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom14(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom14(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom14 is a paid mutator transaction binding the contract method 0x1bbffe6f.
//
// Solidity: function transferFrom14(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom14(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom14(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom15 is a paid mutator transaction binding the contract method 0x8789ca67.
//
// Solidity: function transferFrom15(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom15(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom15", from, to, value)
}

// TransferFrom15 is a paid mutator transaction binding the contract method 0x8789ca67.
//
// Solidity: function transferFrom15(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom15(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom15(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom15 is a paid mutator transaction binding the contract method 0x8789ca67.
//
// Solidity: function transferFrom15(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom15(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom15(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom16 is a paid mutator transaction binding the contract method 0x39e0bd12.
//
// Solidity: function transferFrom16(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom16(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom16", from, to, value)
}

// TransferFrom16 is a paid mutator transaction binding the contract method 0x39e0bd12.
//
// Solidity: function transferFrom16(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom16(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom16(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom16 is a paid mutator transaction binding the contract method 0x39e0bd12.
//
// Solidity: function transferFrom16(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom16(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom16(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom17 is a paid mutator transaction binding the contract method 0x291c3bd7.
//
// Solidity: function transferFrom17(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom17(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom17", from, to, value)
}

// TransferFrom17 is a paid mutator transaction binding the contract method 0x291c3bd7.
//
// Solidity: function transferFrom17(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom17(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom17(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom17 is a paid mutator transaction binding the contract method 0x291c3bd7.
//
// Solidity: function transferFrom17(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom17(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom17(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom18 is a paid mutator transaction binding the contract method 0x657b6ef7.
//
// Solidity: function transferFrom18(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom18(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom18", from, to, value)
}

// TransferFrom18 is a paid mutator transaction binding the contract method 0x657b6ef7.
//
// Solidity: function transferFrom18(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom18(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom18(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom18 is a paid mutator transaction binding the contract method 0x657b6ef7.
//
// Solidity: function transferFrom18(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom18(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom18(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom19 is a paid mutator transaction binding the contract method 0xeb4329c8.
//
// Solidity: function transferFrom19(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom19(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom19", from, to, value)
}

// TransferFrom19 is a paid mutator transaction binding the contract method 0xeb4329c8.
//
// Solidity: function transferFrom19(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom19(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom19(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom19 is a paid mutator transaction binding the contract method 0xeb4329c8.
//
// Solidity: function transferFrom19(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom19(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom19(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom2 is a paid mutator transaction binding the contract method 0x6c12ed28.
//
// Solidity: function transferFrom2(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom2(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom2", from, to, value)
}

// TransferFrom2 is a paid mutator transaction binding the contract method 0x6c12ed28.
//
// Solidity: function transferFrom2(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom2(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom2(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom2 is a paid mutator transaction binding the contract method 0x6c12ed28.
//
// Solidity: function transferFrom2(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom2(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom2(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom3 is a paid mutator transaction binding the contract method 0x1b17c65c.
//
// Solidity: function transferFrom3(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom3(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom3", from, to, value)
}

// TransferFrom3 is a paid mutator transaction binding the contract method 0x1b17c65c.
//
// Solidity: function transferFrom3(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom3(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom3(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom3 is a paid mutator transaction binding the contract method 0x1b17c65c.
//
// Solidity: function transferFrom3(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom3(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom3(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom4 is a paid mutator transaction binding the contract method 0x3a131990.
//
// Solidity: function transferFrom4(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom4(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom4", from, to, value)
}

// TransferFrom4 is a paid mutator transaction binding the contract method 0x3a131990.
//
// Solidity: function transferFrom4(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom4(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom4(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom4 is a paid mutator transaction binding the contract method 0x3a131990.
//
// Solidity: function transferFrom4(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom4(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom4(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom5 is a paid mutator transaction binding the contract method 0x0460faf6.
//
// Solidity: function transferFrom5(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom5(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom5", from, to, value)
}

// TransferFrom5 is a paid mutator transaction binding the contract method 0x0460faf6.
//
// Solidity: function transferFrom5(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom5(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom5(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom5 is a paid mutator transaction binding the contract method 0x0460faf6.
//
// Solidity: function transferFrom5(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom5(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom5(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom6 is a paid mutator transaction binding the contract method 0x7c66673e.
//
// Solidity: function transferFrom6(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom6(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom6", from, to, value)
}

// TransferFrom6 is a paid mutator transaction binding the contract method 0x7c66673e.
//
// Solidity: function transferFrom6(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom6(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom6(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom6 is a paid mutator transaction binding the contract method 0x7c66673e.
//
// Solidity: function transferFrom6(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom6(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom6(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom7 is a paid mutator transaction binding the contract method 0xfaf35ced.
//
// Solidity: function transferFrom7(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom7(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom7", from, to, value)
}

// TransferFrom7 is a paid mutator transaction binding the contract method 0xfaf35ced.
//
// Solidity: function transferFrom7(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom7(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom7(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom7 is a paid mutator transaction binding the contract method 0xfaf35ced.
//
// Solidity: function transferFrom7(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom7(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom7(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom8 is a paid mutator transaction binding the contract method 0x552a1b56.
//
// Solidity: function transferFrom8(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom8(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom8", from, to, value)
}

// TransferFrom8 is a paid mutator transaction binding the contract method 0x552a1b56.
//
// Solidity: function transferFrom8(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom8(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom8(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom8 is a paid mutator transaction binding the contract method 0x552a1b56.
//
// Solidity: function transferFrom8(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom8(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom8(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom9 is a paid mutator transaction binding the contract method 0x4f7bd75a.
//
// Solidity: function transferFrom9(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactor) TransferFrom9(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.contract.Transact(opts, "transferFrom9", from, to, value)
}

// TransferFrom9 is a paid mutator transaction binding the contract method 0x4f7bd75a.
//
// Solidity: function transferFrom9(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenSession) TransferFrom9(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom9(&_StateBloatToken.TransactOpts, from, to, value)
}

// TransferFrom9 is a paid mutator transaction binding the contract method 0x4f7bd75a.
//
// Solidity: function transferFrom9(address from, address to, uint256 value) returns(bool)
func (_StateBloatToken *StateBloatTokenTransactorSession) TransferFrom9(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _StateBloatToken.Contract.TransferFrom9(&_StateBloatToken.TransactOpts, from, to, value)
}

// StateBloatTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the StateBloatToken contract.
type StateBloatTokenApprovalIterator struct {
	Event *StateBloatTokenApproval // Event containing the contract specifics and raw log

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
func (it *StateBloatTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateBloatTokenApproval)
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
		it.Event = new(StateBloatTokenApproval)
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
func (it *StateBloatTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateBloatTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateBloatTokenApproval represents a Approval event raised by the StateBloatToken contract.
type StateBloatTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StateBloatToken *StateBloatTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StateBloatTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StateBloatToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StateBloatTokenApprovalIterator{contract: _StateBloatToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StateBloatToken *StateBloatTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StateBloatTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StateBloatToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateBloatTokenApproval)
				if err := _StateBloatToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_StateBloatToken *StateBloatTokenFilterer) ParseApproval(log types.Log) (*StateBloatTokenApproval, error) {
	event := new(StateBloatTokenApproval)
	if err := _StateBloatToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateBloatTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StateBloatToken contract.
type StateBloatTokenTransferIterator struct {
	Event *StateBloatTokenTransfer // Event containing the contract specifics and raw log

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
func (it *StateBloatTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateBloatTokenTransfer)
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
		it.Event = new(StateBloatTokenTransfer)
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
func (it *StateBloatTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateBloatTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateBloatTokenTransfer represents a Transfer event raised by the StateBloatToken contract.
type StateBloatTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StateBloatToken *StateBloatTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StateBloatTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StateBloatToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StateBloatTokenTransferIterator{contract: _StateBloatToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StateBloatToken *StateBloatTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StateBloatTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StateBloatToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateBloatTokenTransfer)
				if err := _StateBloatToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_StateBloatToken *StateBloatTokenFilterer) ParseTransfer(log types.Log) (*StateBloatTokenTransfer, error) {
	event := new(StateBloatTokenTransfer)
	if err := _StateBloatToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
