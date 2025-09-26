package types

import (
	"math/big"
	"nonogram-solver/internal/combinatorics"
)

// Line represents a single row or column in the nonogram
type Line struct {
	ID           LineID
	Direction    Direction
	Length       int
	Clues        []ClueItem
	Facts        *Facts
	Combinations combinatorics.CombinationsProvider
}

// Bitset is an alias for combinatorics.Bitset for backward compatibility
type Bitset = combinatorics.Bitset

// NewBitset creates a new Bitset with the given value
func NewBitset(value *big.Int) *Bitset {
	return combinatorics.NewBitset(value)
}
