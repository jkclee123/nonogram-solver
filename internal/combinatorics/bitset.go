package combinatorics

import "math/big"

// Bitset is a wrapper around math/big.Int for bit operations
type Bitset struct {
	*big.Int
}

// NewBitset creates a new Bitset with the given value
func NewBitset(value *big.Int) *Bitset {
	return &Bitset{Int: value}
}
