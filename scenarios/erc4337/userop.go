package erc4337

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethpandaops/spamoor/scenarios/erc4337/contract"
)

// ABI primitive types reused when ABI-encoding UserOperation fields and the
// userOpHash preimage. Parsed once at init - the inputs are constant so the
// errors can never trigger.
var (
	abiAddressTy = mustNewABIType("address")
	abiUint256Ty = mustNewABIType("uint256")
	abiBytes32Ty = mustNewABIType("bytes32")
	abiBytesTy   = mustNewABIType("bytes")
)

func mustNewABIType(t string) abi.Type {
	ty, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(fmt.Sprintf("erc4337: invalid abi type %q: %v", t, err))
	}
	return ty
}

// userOpGasConfig holds the per-UserOperation gas parameters shared across a run.
type userOpGasConfig struct {
	VerificationGasLimit uint64
	CallGasLimit         uint64
	PreVerificationGas   uint64
	PaymasterVerifGas    uint64
	PaymasterPostOpGas   uint64
}

// packUint128Pair packs two values into a single bytes32 as
// (high << 128) | low, matching the ERC-4337 v0.7 UserOperationLib encoding of
// accountGasLimits and gasFees. Each value must fit in 128 bits.
func packUint128Pair(high, low *big.Int) [32]byte {
	var out [32]byte
	high.FillBytes(out[0:16])
	low.FillBytes(out[16:32])
	return out
}

// buildPaymasterAndData builds the v0.7 paymasterAndData field:
// paymaster(20) || paymasterVerificationGasLimit(16) || paymasterPostOpGasLimit(16).
// No paymaster-specific data is appended (the accept-all paymaster needs none).
func buildPaymasterAndData(paymaster common.Address, verifGas, postOpGas uint64) []byte {
	out := make([]byte, 52)
	copy(out[0:20], paymaster.Bytes())
	new(big.Int).SetUint64(verifGas).FillBytes(out[20:36])
	new(big.Int).SetUint64(postOpGas).FillBytes(out[36:52])
	return out
}

// encodeInitCode builds the UserOperation initCode that deploys a fresh
// SimpleAccount for owner at the given salt: factory(20) || createAccount(owner, salt).
func encodeInitCode(factory, owner common.Address, salt *big.Int) ([]byte, error) {
	factoryABI, err := contract.SimpleAccountFactoryMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("could not load factory abi: %w", err)
	}
	call, err := factoryABI.Pack("createAccount", owner, salt)
	if err != nil {
		return nil, fmt.Errorf("could not pack createAccount: %w", err)
	}
	return append(factory.Bytes(), call...), nil
}

// encodeCounterIncrementCall builds the SimpleAccount.execute(...) calldata that
// invokes Counter.increment() with no value - the inner action of each op.
func encodeCounterIncrementCall(counter common.Address) ([]byte, error) {
	counterABI, err := contract.CounterMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("could not load counter abi: %w", err)
	}
	innerCall, err := counterABI.Pack("increment")
	if err != nil {
		return nil, fmt.Errorf("could not pack increment: %w", err)
	}

	executeArgs := abi.Arguments{{Type: abiAddressTy}, {Type: abiUint256Ty}, {Type: abiBytesTy}}
	packed, err := executeArgs.Pack(counter, big.NewInt(0), innerCall)
	if err != nil {
		return nil, fmt.Errorf("could not pack execute args: %w", err)
	}

	// execute(address dest, uint256 value, bytes func)
	selector := crypto.Keccak256([]byte("execute(address,uint256,bytes)"))[:4]
	return append(selector, packed...), nil
}

// computeUserOpHash reproduces EntryPoint.getUserOpHash for v0.7:
//
//	inner = keccak256(abi.encode(sender, nonce, keccak(initCode), keccak(callData),
//	                             accountGasLimits, preVerificationGas, gasFees,
//	                             keccak(paymasterAndData)))
//	hash  = keccak256(abi.encode(inner, entryPoint, chainId))
func computeUserOpHash(op *contract.PackedUserOperation, entryPoint common.Address, chainID *big.Int) ([32]byte, error) {
	innerArgs := abi.Arguments{
		{Type: abiAddressTy}, {Type: abiUint256Ty}, {Type: abiBytes32Ty}, {Type: abiBytes32Ty},
		{Type: abiBytes32Ty}, {Type: abiUint256Ty}, {Type: abiBytes32Ty}, {Type: abiBytes32Ty},
	}
	inner, err := innerArgs.Pack(
		op.Sender,
		op.Nonce,
		[32]byte(crypto.Keccak256Hash(op.InitCode)),
		[32]byte(crypto.Keccak256Hash(op.CallData)),
		op.AccountGasLimits,
		op.PreVerificationGas,
		op.GasFees,
		[32]byte(crypto.Keccak256Hash(op.PaymasterAndData)),
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not pack userOp preimage: %w", err)
	}
	innerHash := crypto.Keccak256Hash(inner)

	outerArgs := abi.Arguments{{Type: abiBytes32Ty}, {Type: abiAddressTy}, {Type: abiUint256Ty}}
	outer, err := outerArgs.Pack([32]byte(innerHash), entryPoint, chainID)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not pack userOpHash preimage: %w", err)
	}
	return crypto.Keccak256Hash(outer), nil
}

// signUserOp signs the userOpHash the way SimpleAccount expects: an EIP-191
// personal_sign over the 32-byte userOpHash, with the recovery id normalized to
// the Ethereum 27/28 convention.
func signUserOp(userOpHash [32]byte, key *ecdsa.PrivateKey) ([]byte, error) {
	if key == nil {
		return nil, fmt.Errorf("owner wallet has no private key")
	}
	sig, err := crypto.Sign(accounts.TextHash(userOpHash[:]), key)
	if err != nil {
		return nil, fmt.Errorf("could not sign userOp: %w", err)
	}
	sig[64] += 27
	return sig, nil
}

// buildSignedUserOp assembles a fully-populated, signed UserOperation that calls
// Counter.increment() through the account at sender (owned by owner) and signs
// it with ownerKey.
//
// When withInitCode is true the op carries the factory initCode that deploys the
// account (used for the first op against a counterfactual address). When false
// the account is assumed already deployed and is reused - no new contract is
// created. nonce is the account's EntryPoint nonce (key 0); reused accounts use
// sequential nonces while a fresh-deployment op uses 0.
func buildSignedUserOp(
	cfg *userOpGasConfig,
	entryPoint, factory, paymaster, counter common.Address,
	owner common.Address,
	ownerKey *ecdsa.PrivateKey,
	sender common.Address,
	salt *big.Int,
	nonce uint64,
	withInitCode bool,
	chainID *big.Int,
	maxFeePerGas, maxPriorityFeePerGas *big.Int,
) (contract.PackedUserOperation, error) {
	var initCode []byte
	if withInitCode {
		ic, err := encodeInitCode(factory, owner, salt)
		if err != nil {
			return contract.PackedUserOperation{}, err
		}
		initCode = ic
	}
	callData, err := encodeCounterIncrementCall(counter)
	if err != nil {
		return contract.PackedUserOperation{}, err
	}

	op := contract.PackedUserOperation{
		Sender:             sender,
		Nonce:              new(big.Int).SetUint64(nonce),
		InitCode:           initCode,
		CallData:           callData,
		AccountGasLimits:   packUint128Pair(new(big.Int).SetUint64(cfg.VerificationGasLimit), new(big.Int).SetUint64(cfg.CallGasLimit)),
		PreVerificationGas: new(big.Int).SetUint64(cfg.PreVerificationGas),
		GasFees:            packUint128Pair(maxPriorityFeePerGas, maxFeePerGas),
		PaymasterAndData:   buildPaymasterAndData(paymaster, cfg.PaymasterVerifGas, cfg.PaymasterPostOpGas),
	}

	hash, err := computeUserOpHash(&op, entryPoint, chainID)
	if err != nil {
		return contract.PackedUserOperation{}, err
	}
	sig, err := signUserOp(hash, ownerKey)
	if err != nil {
		return contract.PackedUserOperation{}, err
	}
	op.Signature = sig

	return op, nil
}
