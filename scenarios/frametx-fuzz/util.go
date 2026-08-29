package frametxfuzz

import (
	"encoding/binary"
	"math/big"

	"github.com/ethpandaops/spamoor/txtypes"
)

// big0 is the genesis block number, named so the header lookups read as intent.
var big0 = new(big.Int)

// bigFromUint64 renders a block number for the header lookups.
func bigFromUint64(v uint64) *big.Int { return new(big.Int).SetUint64(v) }

// expiryData encodes the expiry verifier's 8-byte deadline.
func expiryData(deadline uint64) []byte {
	data := make([]byte, txtypes.ExpiryDataLength)
	binary.BigEndian.PutUint64(data, deadline)

	return data
}
