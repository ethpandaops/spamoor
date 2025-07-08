package aitx

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoor"
)

type PayloadTemplate struct {
	Type         string `json:"type"`
	Description  string `json:"description"`
	InitCode     string `json:"init_code"`
	RunCode      string `json:"run_code"`
	PostCode     string `json:"post_code,omitempty"`     // Optional - executes once at end
	GasRemainder string `json:"gas_remainder,omitempty"` // Optional - defaults to 10000
	Calldata     string `json:"calldata,omitempty"`      // Optional - calldata for contract call
}

type PayloadInstance struct {
	Type         string
	Description  string
	InitCode     string
	RunCode      string
	PostCode     string
	GasRemainder uint64
	Calldata     []byte
}

type PlaceholderSubstituter struct {
	walletPool *spamoor.WalletPool
	client     *spamoor.Client
	logger     logrus.FieldLogger
}

func NewPlaceholderSubstituter(walletPool *spamoor.WalletPool, client *spamoor.Client, logger logrus.FieldLogger) *PlaceholderSubstituter {
	return &PlaceholderSubstituter{
		walletPool: walletPool,
		client:     client,
		logger:     logger,
	}
}

func (pt *PayloadTemplate) Substitute(substituter *PlaceholderSubstituter) (*PayloadInstance, error) {
	instance := &PayloadInstance{
		Type:        pt.Type,
		Description: pt.Description,
	}

	var err error
	instance.InitCode, err = substituter.SubstitutePlaceholders(pt.InitCode)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute init code placeholders: %w", err)
	}

	instance.RunCode, err = substituter.SubstitutePlaceholders(pt.RunCode)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute run code placeholders: %w", err)
	}

	instance.PostCode, err = substituter.SubstitutePlaceholders(pt.PostCode)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute post code placeholders: %w", err)
	}

	gasRemainderStr, err := substituter.SubstitutePlaceholders(pt.GasRemainder)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute gas remainder placeholders: %w", err)
	}

	if gasRemainderStr != "" {
		gasRemainder, err := strconv.ParseUint(gasRemainderStr, 10, 64)
		if err != nil {
			instance.GasRemainder = 10000 // Default value
		} else {
			instance.GasRemainder = gasRemainder
		}
	} else {
		instance.GasRemainder = 10000 // Default value
	}

	// Handle calldata
	if pt.Calldata != "" {
		calldataStr, err := substituter.SubstitutePlaceholders(pt.Calldata)
		if err != nil {
			return nil, fmt.Errorf("failed to substitute calldata placeholders: %w", err)
		}

		// Parse hex calldata
		calldataStr = strings.TrimPrefix(calldataStr, "0x")
		instance.Calldata, err = hex.DecodeString(calldataStr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode calldata: %w", err)
		}
	}

	return instance, instance.Validate()
}

func (pi *PayloadInstance) Validate() error {
	if pi.Type != "geas" {
		return fmt.Errorf("only 'geas' type is supported, got: %s", pi.Type)
	}

	if pi.Description == "" {
		return fmt.Errorf("payload description is required")
	}

	if pi.InitCode == "" {
		return fmt.Errorf("geas payload requires init_code")
	}

	if pi.RunCode == "" {
		return fmt.Errorf("geas payload requires run_code")
	}

	return nil
}

func (ps *PlaceholderSubstituter) SubstitutePlaceholders(input string) (string, error) {
	result := input

	substitutions := map[string]func() string{
		"${WALLET_ADDRESS}": func() string {
			if ps.walletPool.GetWalletCount() == 0 {
				return "0x0000000000000000000000000000000000000000"
			}
			walletIdx := time.Now().UnixNano() % int64(ps.walletPool.GetWalletCount())
			wallet := ps.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(walletIdx))
			return wallet.GetAddress().Hex()
		},
		"${RANDOM_ADDRESS}": func() string {
			bytes := make([]byte, 20)
			rand.Read(bytes)
			return common.BytesToAddress(bytes).Hex()
		},
		"${ETH_AMOUNT_SMALL}": func() string {
			min := big.NewInt(1000000000000000)  // 0.001 ETH
			max := big.NewInt(10000000000000000) // 0.01 ETH
			diff := new(big.Int).Sub(max, min)
			randInt, _ := rand.Int(rand.Reader, diff)
			return new(big.Int).Add(min, randInt).String()
		},
		"${ETH_AMOUNT_MEDIUM}": func() string {
			min := big.NewInt(10000000000000000)  // 0.01 ETH
			max := big.NewInt(100000000000000000) // 0.1 ETH
			diff := new(big.Int).Sub(max, min)
			randInt, _ := rand.Int(rand.Reader, diff)
			return new(big.Int).Add(min, randInt).String()
		},
		"${ETH_AMOUNT_LARGE}": func() string {
			min := big.NewInt(100000000000000000)  // 0.1 ETH
			max := big.NewInt(1000000000000000000) // 1.0 ETH
			diff := new(big.Int).Sub(max, min)
			randInt, _ := rand.Int(rand.Reader, diff)
			return new(big.Int).Add(min, randInt).String()
		},
		"${GAS_LIMIT_LOW}": func() string {
			randBytes := make([]byte, 4)
			rand.Read(randBytes)
			randVal := int64(randBytes[0])<<24 | int64(randBytes[1])<<16 | int64(randBytes[2])<<8 | int64(randBytes[3])
			if randVal < 0 {
				randVal = -randVal
			}
			return fmt.Sprintf("%d", 21000+(randVal%29000))
		},
		"${GAS_LIMIT_MEDIUM}": func() string {
			randBytes := make([]byte, 4)
			rand.Read(randBytes)
			randVal := int64(randBytes[0])<<24 | int64(randBytes[1])<<16 | int64(randBytes[2])<<8 | int64(randBytes[3])
			if randVal < 0 {
				randVal = -randVal
			}
			return fmt.Sprintf("%d", 50000+(randVal%150000))
		},
		"${GAS_LIMIT_HIGH}": func() string {
			randBytes := make([]byte, 4)
			rand.Read(randBytes)
			randVal := int64(randBytes[0])<<24 | int64(randBytes[1])<<16 | int64(randBytes[2])<<8 | int64(randBytes[3])
			if randVal < 0 {
				randVal = -randVal
			}
			return fmt.Sprintf("%d", 200000+(randVal%800000))
		},
		"${RANDOM_UINT256}": func() string {
			bytes := make([]byte, 32)
			rand.Read(bytes)
			return new(big.Int).SetBytes(bytes).String()
		},
		"${RANDOM_BYTES32}": func() string {
			bytes := make([]byte, 32)
			rand.Read(bytes)
			return "0x" + hex.EncodeToString(bytes)
		},
		"${CURRENT_BLOCK}": func() string {
			return "0"
		},
		"${LOOP_COUNT_SMALL}": func() string {
			randBytes := make([]byte, 1)
			rand.Read(randBytes)
			return fmt.Sprintf("%d", 1+(int(randBytes[0])%10))
		},
		"${LOOP_COUNT_MEDIUM}": func() string {
			randBytes := make([]byte, 1)
			rand.Read(randBytes)
			return fmt.Sprintf("%d", 10+(int(randBytes[0])%90))
		},
		"${LOOP_COUNT_LARGE}": func() string {
			randBytes := make([]byte, 2)
			rand.Read(randBytes)
			randVal := int(randBytes[0])<<8 | int(randBytes[1])
			return fmt.Sprintf("%d", 100+(randVal%900))
		},
	}

	for placeholder, substituteFn := range substitutions {
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, substituteFn())
		}
	}

	return result, nil
}
