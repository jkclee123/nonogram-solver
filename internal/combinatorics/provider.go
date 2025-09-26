package combinatorics

import (
	"math/big"
)

// Provider exposes lazy combination retrieval per color.
type Provider interface {
	Get(colorID int) []*big.Int
}
