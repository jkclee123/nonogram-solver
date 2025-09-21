package types

import (
	"fmt"
	"math/big"
)

type Line struct {
	Blocks []Block
	Size   int
}

func IsSameColor(l *Line, index uint8, index2 uint8) bool {
	return l.Blocks[index].ColorID == l.Blocks[index2].ColorID
}

func (l *Line) CheckCommonCombinations() map[uint8]*big.Int {
	result := l.CheckCommonFillCell()
	result[0] = l.CheckCommonEmptyCell()
	return result
}

// CheckCommonFillCell finds cells that must be filled with a specific color across all valid combinations.
// For each color, it performs bitwise AND on all combinations within each block to find positions
// that are consistently filled in every combination of that block, then ORs these results together.
// Returns a map where keys are color IDs and values are big.Int bitmasks where 1 bits represent
// positions that must be filled with that color in every valid combination.
func (l *Line) CheckCommonFillCell() map[uint8]*big.Int {
	result := make(map[uint8]*big.Int)

	for _, block := range l.Blocks {
		commonFill := block.BitwiseAnd()
		if commonFill.Cmp(big.NewInt(0)) != 0 {
			if result[block.ColorID] == nil {
				result[block.ColorID] = big.NewInt(0)
			}
			result[block.ColorID].Or(result[block.ColorID], commonFill)
		}
	}

	return result
}

// CheckCommonEmptyCell finds cells that must remain empty across all valid combinations.
// If there are no blocks, all cells must be empty.
// Otherwise, for each block it performs bitwise OR on all combinations within that block,
// then ORs these results together across all blocks, and finally inverts the result.
// Returns a big.Int bitmask where 1 bits represent positions that must be empty
// in every valid combination (cells that are never filled in any combination).
func (l *Line) CheckCommonEmptyCell() *big.Int {
	mask := big.NewInt(1)
	mask.Lsh(mask, uint(l.Size))
	mask.Sub(mask, big.NewInt(1))

	if len(l.Blocks) == 0 {
		// Return a big int with bits from 0 to size-1 set to 1
		return mask
	}

	commonEmpty := big.NewInt(0)
	for _, block := range l.Blocks {
		commonEmpty.Or(commonEmpty, block.BitwiseOr())
	}

	commonEmpty.Xor(commonEmpty, mask)
	return commonEmpty
}

// Print prints the line information in a readable format
func (l *Line) Print() {
	fmt.Printf("Line{Size: %d, Blocks: [", l.Size)
	for i, block := range l.Blocks {
		if i > 0 {
			fmt.Print(", ")
		}
		// Delegate combination rendering to Block
		block.PrintWithWidth(l.Size)
	}
	fmt.Println("]}")
}
