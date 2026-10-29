package codeadopt

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	opSetCodeFrom = 0xa7
	opSendAll     = 0xff
	opApprove     = 0xaa
)

var valueSourceCode = common.FromHex("604260005260206000f3")

var approverSourceCode = []byte{
	0x60, 0x00, 0x35,
	0x60, 0x00,
	0x60, 0x00,
	opApprove,
	0x00,
}

var adopterCode = []byte{
	0x60, 0x00, 0x35,
	opSetCodeFrom,
	0x60, 0x00, 0x52,
	0x60, 0x20, 0x60, 0x00, 0xf3,
}

func push20(address common.Address) []byte {
	return append([]byte{0x73}, address.Bytes()...)
}

func deployCode(runtime []byte) []byte {
	return append(deployPrefix(len(runtime), 0x0b), runtime...)
}

func deployPrefix(length int, offset byte) []byte {
	return []byte{0x60, byte(length), 0x80, 0x60, offset, 0x60, 0x00, 0x39, 0x60, 0x00, 0xf3}
}

func staticCallerCode(target common.Address) []byte {
	code := []byte{0x60, 0x20, 0x60, 0x00, 0x60, 0x00, 0x37, 0x60, 0x20, 0x60, 0x20, 0x60, 0x20, 0x60, 0x00}
	code = append(code, push20(target)...)
	code = append(code, 0x5a, 0xfa, 0x60, 0x00, 0x52, 0x60, 0x40, 0x60, 0x00, 0xf3)

	return code
}

func sendAllCode(beneficiary common.Address) []byte {
	return append(push20(beneficiary), opSendAll)
}

func sendAllDeployCode(beneficiary common.Address, marker byte) []byte {
	runtime := sendAllCode(beneficiary)
	code := []byte{0x60, marker, 0x60, 0x00, 0x55}
	code = append(code, deployPrefix(len(runtime), byte(len(code)+11))...)

	return append(code, runtime...)
}

func sendAllInitCode(beneficiary common.Address, marker byte) []byte {
	code := []byte{0x60, marker, 0x60, 0x00, 0x55}

	return append(code, sendAllCode(beneficiary)...)
}

func adoptInInitCode(source common.Address) []byte {
	code := push20(source)

	return append(code, opSetCodeFrom, 0x50, 0x60, 0xef, 0x60, 0x00, 0x53, 0x60, 0x01, 0x60, 0x00, 0xf3)
}

func adoptSameTxSourceInitCode(source common.Address) []byte {
	code := common.FromHex("6460016000f3600052600560" + "1b6000f0")
	code = append(code, opSetCodeFrom, 0x60, 0x00, 0x55)
	code = append(code, push20(source)...)

	return append(code, opSetCodeFrom, 0x60, 0x01, 0x55)
}
