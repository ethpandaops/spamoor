package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// PlaceholderFn resolves a `{name}` or `{name:arg}` placeholder; arg is empty for `{name}`.
type PlaceholderFn func(arg string) (string, error)

// Placeholders maps placeholder names to their resolvers.
type Placeholders map[string]PlaceholderFn

var placeholderRegex = regexp.MustCompile(`\{([a-z0-9_]+)(?::([^{}]*))?\}`)

// Resolve substitutes placeholders innermost-first, so arguments can nest
// (`{create2:{factory_address}:0x..:{txid}}`). Unknown names are left untouched.
func (p Placeholders) Resolve(s string) (string, error) {
	for {
		var resolveErr error
		changed := false
		out := placeholderRegex.ReplaceAllStringFunc(s, func(match string) string {
			sub := placeholderRegex.FindStringSubmatch(match)
			fn, ok := p[sub[1]]
			if !ok || resolveErr != nil {
				return match
			}
			value, err := fn(sub[2])
			if err != nil {
				resolveErr = fmt.Errorf("placeholder %s: %w", match, err)
				return match
			}
			changed = true
			return value
		})
		if resolveErr != nil {
			return "", resolveErr
		}
		if !changed {
			return out, nil
		}
		s = out
	}
}

// TxPlaceholders returns the placeholders every scenario supports for transaction txIdx:
// {txid}, {random}, {random:N} (0 <= x < N), {randomaddr} and {create2:factory:initcodehash:salt}.
func TxPlaceholders(txIdx uint64) Placeholders {
	return Placeholders{
		"txid": func(string) (string, error) {
			return strconv.FormatUint(txIdx, 10), nil
		},
		"random": func(arg string) (string, error) {
			if arg == "" {
				return generateRandomUint256()
			}
			maxVal, err := strconv.ParseUint(arg, 10, 64)
			if err != nil {
				return "", fmt.Errorf("invalid random range: %s", arg)
			}
			value, err := generateRandomInRange(maxVal)
			if err != nil {
				return "", err
			}
			return strconv.FormatUint(value, 10), nil
		},
		"randomaddr": func(string) (string, error) {
			return generateRandomAddress()
		},
		"create2": func(arg string) (string, error) {
			addr, err := create2Address(arg)
			if err != nil {
				return "", err
			}
			return addr.Hex(), nil
		},
	}
}

// create2Address derives a CREATE2 address from "factory:initcodehash:salt"; salt is a
// decimal or 0x-prefixed uint256, matching the salt encoding of factorydeploytx.
func create2Address(arg string) (common.Address, error) {
	parts := strings.Split(arg, ":")
	if len(parts) != 3 {
		return common.Address{}, fmt.Errorf("expected factory:initcodehash:salt")
	}
	if !common.IsHexAddress(parts[0]) {
		return common.Address{}, fmt.Errorf("invalid factory address %q", parts[0])
	}
	initCodeHash, err := hexToHash(parts[1])
	if err != nil {
		return common.Address{}, err
	}
	salt, ok := new(big.Int).SetString(parts[2], 0)
	if !ok || salt.Sign() < 0 || salt.BitLen() > 256 {
		return common.Address{}, fmt.Errorf("invalid salt %q", parts[2])
	}
	var saltBytes [32]byte
	salt.FillBytes(saltBytes[:])
	return crypto.CreateAddress2(common.HexToAddress(parts[0]), saltBytes, initCodeHash[:]), nil
}

func hexToHash(s string) (common.Hash, error) {
	b, err := hexutil.Decode("0x" + strings.TrimPrefix(s, "0x"))
	if err != nil || len(b) != common.HashLength {
		return common.Hash{}, fmt.Errorf("invalid init code hash %q", s)
	}
	return common.BytesToHash(b), nil
}

func generateRandomUint256() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return new(big.Int).SetBytes(bytes).String(), nil
}

func generateRandomInRange(max uint64) (uint64, error) {
	if max == 0 {
		return 0, nil
	}
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return 0, err
	}
	return new(big.Int).SetBytes(bytes).Uint64() % max, nil
}

func generateRandomAddress() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return common.BytesToAddress(bytes).Hex(), nil
}
