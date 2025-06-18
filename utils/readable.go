package utils

import (
	"fmt"

	"github.com/holiman/uint256"
)

func ReadableAmount(amount *uint256.Int) string {
	if amount == nil {
		return "0 Wei"
	}
	val := new(uint256.Int).Set(amount)
	if val.Cmp(uint256.NewInt(1e15)) >= 0 { // >= 0.001 ETH
		eth := new(uint256.Int).Div(val, uint256.NewInt(1e18))
		dec := new(uint256.Int).Mod(val, uint256.NewInt(1e18))
		if dec.IsZero() {
			return fmt.Sprintf("%d ETH", eth.Uint64())
		}
		return fmt.Sprintf("%.3f ETH", float64(eth.Uint64())+(float64(dec.Uint64())/1e18))
	} else if val.Cmp(uint256.NewInt(1e9)) >= 0 { // >= 1 Gwei
		return fmt.Sprintf("%.2f Gwei", float64(val.Uint64())/1e9)
	}
	return fmt.Sprintf("%d Wei", val.Uint64())
}
