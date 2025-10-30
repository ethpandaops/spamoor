package evmfuzz

import "math/big"

// safeFillBytes safely fills a byte slice with a big.Int value, handling cases where the value is too large
func safeFillBytes(value *big.Int, buf []byte) {
	// Get the bytes representation of the value
	valueBytes := value.Bytes()

	// Clear the buffer first
	for i := range buf {
		buf[i] = 0
	}

	// If the value is larger than the buffer, take only the least significant bytes
	if len(valueBytes) > len(buf) {
		copy(buf, valueBytes[len(valueBytes)-len(buf):])
	} else {
		// Normal case: copy to the end of the buffer (big-endian)
		copy(buf[len(buf)-len(valueBytes):], valueBytes)
	}
}
