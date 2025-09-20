package types

import (
	"fmt"
	"math/big"
)

type Line struct {
	Blocks []Block
	Size   uint8
}

func IsSameColor(l *Line, index uint8, index2 uint8) bool {
	return l.Blocks[index].ColorID == l.Blocks[index2].ColorID
}

func (l *Line) CheckCommonCombinations() map[uint8][]uint8 {
	result := l.CheckCommonFillCell()
	result[0] = l.CheckCommonEmptyCell()
	return result
}

// CheckCommonFillCell finds cells that must be filled with a specific color across all valid combinations.
// For each color, it performs bitwise AND on all combinations to find positions that are consistently filled.
// Returns a map where keys are color IDs and values are slices of cell positions that must be that color.
func (l *Line) CheckCommonFillCell() map[uint8][]uint8 {
	result := make(map[uint8][]uint8)
	commonColorMap := make(map[uint8]*big.Int)

	for _, block := range l.Blocks {
		commonColor, exists := commonColorMap[block.ColorID]
		if !exists {
			commonColor = big.NewInt(0)
			for i := 0; i < int(l.Size); i++ {
				commonColor.SetBit(commonColor, i, 1)
			}
		}

		for e := block.Combinations.Front(); e != nil; e = e.Next() {
			if combination, ok := e.Value.(*big.Int); ok {
				commonColor.And(commonColor, combination)
			}
		}
		commonColorMap[block.ColorID] = commonColor
	}

	for colorID, commonColor := range commonColorMap {
		for i := 0; i < int(l.Size); i++ {
			if commonColor.Bit(i) == 1 {
				result[colorID] = append(result[colorID], uint8(i))
			}
		}
	}

	return result
}

// CheckCommonEmptyCell finds cells that must remain empty across all valid combinations.
// If there are no blocks, all cells must be empty.
// Otherwise, it performs bitwise OR on all combinations from all blocks
// positions with 0 bits are cells that must be empty in every valid combination.
// Returns a slice of cell positions that must be empty.
func (l *Line) CheckCommonEmptyCell() []uint8 {
	var result []uint8
	if len(l.Blocks) == 0 {
		// return range from 0 to size - 1
		for i := uint8(0); i < l.Size; i++ {
			result = append(result, i)
		}
		return result
	}

	commonEmpty := big.NewInt(0)
	for _, block := range l.Blocks {
		for e := block.Combinations.Front(); e != nil; e = e.Next() {
			if combination, ok := e.Value.(*big.Int); ok {
				commonEmpty.Or(commonEmpty, combination)
			}
		}
	}

	// Find positions of 0 bits (common empty cells)
	for i := 0; i < int(l.Size); i++ {
		if commonEmpty.Bit(i) == 0 {
			result = append(result, uint8(i))
		}
	}

	return result
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
